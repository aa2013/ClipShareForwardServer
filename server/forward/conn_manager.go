package forward

import (
	"clipshare/db"
	"clipshare/ratelimiter"
	"clipshare/types"
	"clipshare/utils"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

type connListenerInfo struct {
	msg         map[string]string
	selfId      string
	selfDevName string
	platform    string
	appVersion  string
	targetId    string
}

func isUnlimitedDevice(devId string) bool {
	return slices.Contains(types.AppConfig.GetUnlimitedDeviceIds(), devId)
}
func GetDataSyncTransferLimited(devId string) int {
	sktInfo := BaseSocketsMap[devId]
	if sktInfo.Unlimited {
		return 0
	}
	if types.AppConfig.PublicMode {
		return 0
	}
	return int(*sktInfo.PlanType.Rate)
}
func GetFileSyncTransferLimited(devId string) int {
	forwardCfg := types.AppConfig.Forward
	rate := 0
	if forwardCfg.FileTransferLimit.Enabled {
		rate = (*forwardCfg.FileTransferLimit.Rate) * 1024
	}
	return min(GetDataSyncTransferLimited(devId), rate)
}
func checkKey(conn net.Conn, msgInfo connListenerInfo) bool {
	if isUnlimitedDevice(msgInfo.selfId) {
		sendData, _ := json.Marshal(map[string]string{
			"type":      check,
			"result":    "success",
			"unlimited": "true",
		})
		err := sendPacket(conn, sendData)
		if err != nil {
			logs.Error("checkKey send data failed", err)
			_ = conn.Close()
			return false
		}
		return true
	}
	key := msgInfo.msg["key"]
	if types.AppConfig.PublicMode {
		mapData := map[string]string{
			"type":   check,
			"result": "success",
		}
		forwardCfg := types.AppConfig.Forward
		if forwardCfg.FileTransferLimit.Enabled {
			mapData["fileSyncRate"] = strconv.Itoa(*forwardCfg.FileTransferLimit.Rate)
		} else {
			mapData[FileSyncNotAllowed] = "true"
		}
		sendData, _ := json.Marshal(mapData)
		err := sendPacket(conn, sendData)
		if err != nil {
			logs.Error("onCheckKeyTypeConnected send failed", err)
			_ = conn.Close()
			return false
		}
		return true
	}
	if key == "" {
		//no key
		sendData, _ := json.Marshal(map[string]string{
			"type":   check,
			"result": "missing key",
		})
		err := sendPacket(conn, sendData)
		if err != nil {
			logs.Error("checkKey send data failed", err)
			_ = conn.Close()
			return false
		}
		return false
	}
	plan := db.VerifyKey(key)
	if plan == nil {
		sendData, _ := json.Marshal(map[string]string{
			"type":   check,
			"result": "invalid key",
		})
		err := sendPacket(conn, sendData)
		if err != nil {
			logs.Error("checkKey send data failed", err)
			_ = conn.Close()
			return false
		}
		return false
	}
	var deviceLimit string
	if plan.DeviceLimit != nil {
		deviceLimit = strconv.Itoa(int(*plan.DeviceLimit))
	} else {
		deviceLimit = "∞"
	}
	var lifeSpan string
	if plan.Lifespan != nil {
		days := *plan.Lifespan / (24 * 60 * 60)
		lifeSpan = strconv.Itoa(int(days))
	} else {
		lifeSpan = "∞"
	}
	var rate string
	if plan.Rate != nil {
		bytes := *plan.Rate / 1024
		rate = strconv.Itoa(int(bytes))
	} else {
		rate = "∞"
	}
	remainSeconds := db.GetKeyRemainingSeconds(key)
	remaining := "-1"
	if remainSeconds != nil {
		remaining = strconv.FormatInt(int64(*remainSeconds), 10)
	}
	remark := ""
	if plan.Remark != nil {
		remark = *plan.Remark
	}
	sendData, _ := json.Marshal(map[string]string{
		"type":        check,
		"result":      "success",
		"deviceLimit": deviceLimit,
		"lifeSpan":    lifeSpan,
		"remaining":   remaining,
		"rate":        rate,
		"remark":      remark,
	})
	err := sendPacket(conn, sendData)
	if err != nil {
		logs.Error("checkKey send data failed", err)
		_ = conn.Close()
		return false
	}
	if remainSeconds != nil && *remainSeconds == 0 {
		return false
	}
	return true
}
func onBaseTypeConnected(conn net.Conn, msgInfo connListenerInfo, packetReader *ratelimiter.PacketReader) {
	selfId := msgInfo.selfId
	platform := msgInfo.platform
	appVersion := msgInfo.appVersion
	isValid := checkKey(conn, msgInfo)
	if !isValid {
		logs.Warn("onBaseTypeConnected key is invalid")
		sendData, _ := json.Marshal(map[string]string{
			"type":   check,
			"result": "The key is invalid or has expired",
		})
		_ = sendPacket(conn, sendData)
		_ = conn.Close()
		return
	}
	var planType *types.PlanTypeDto
	var firstUseTime *time.Time
	key := msgInfo.msg["key"]
	if key != "" && !types.AppConfig.PublicMode {
		pt := db.VerifyKey(key)
		if pt != nil {
			tmp := pt.ToPlanTypeDto()
			planType = &tmp
		}
		if pt.DeviceLimit != nil {
			//判断是否超出设备连接限制
			cnt := uint(0)
			for _, baseConn := range BaseSocketsMap {
				if baseConn.AccessKey == nil {
					continue
				}
				if *baseConn.AccessKey == key {
					cnt++
				}
				if cnt >= *pt.DeviceLimit {
					sendData, _ := json.Marshal(map[string]string{
						"type":   check,
						"result": "Exceeding the limit of device connections",
					})
					err := sendPacket(conn, sendData)
					if err != nil {
						logs.Error("checkKey send data failed", err)
						_ = conn.Close()
					}
					return
				}
			}
		}
		db.InitializeKeyFirstUseTime(key)
		firstUseTime = db.GetFistUseTime(key)
	}
	oldSkt, connected := BaseSocketsMap[selfId]
	if connected {
		logs.Warn("already connected:", selfId)
		_ = oldSkt.Conn.Close()
	}
	BaseSocketsMap[selfId] = &types.SocketInfo{
		Self:          types.DevInfo{DevId: selfId, Platform: platform, AppVersion: appVersion, DevName: msgInfo.selfDevName},
		ConnType:      Base,
		CreateTime:    time.Now(),
		Conn:          conn,
		Unlimited:     isUnlimitedDevice(selfId),
		AccessKey:     utils.SimpleIf[*string](key == "", nil, &key),
		KeyFirstUseAt: firstUseTime,
		PlanType:      planType,
	}
	//保持连接
	for {
		baseRecData, err := packetReader.ReadPacket()
		if err == nil {
			onBaseConnReceived(selfId, baseRecData)
			continue
		}
		logs.Error("read data failed: ", err, ".self: ", selfId)
		//断开基础连接
		delete(BaseSocketsMap, selfId)
		_ = conn.Close()
		//断开转发连接
		for key := range DataSyncSocketsMap {
			if strings.HasSuffix(key, selfId) || strings.HasPrefix(key, selfId) {
				skt := DataSyncSocketsMap[key]
				_ = skt.Conn.Close()
				delete(DataSyncSocketsMap, key)
			}
		}
		//断开文件发送连接
		for key := range SendFileConnMap {
			if strings.HasSuffix(key, selfId) || strings.HasPrefix(key, selfId) {
				sfi := SendFileConnMap[key]
				_ = sfi.Conn.Close()
				delete(SendFileConnMap, key)
			}
		}
		return
	}
}
func onSendFileTypeConnected(conn net.Conn, msgInfo connListenerInfo) {
	msg := msgInfo.msg
	selfId := msgInfo.selfId
	platform := msgInfo.platform
	appVersion := msgInfo.appVersion
	targetId := msgInfo.targetId
	logs.Info("connType", sendFile, ", self:", selfId)
	//不是白名单且不允许文件传输
	if !isUnlimitedDevice(selfId) && !types.AppConfig.Forward.FileTransferLimit.Enabled {
		sendData, _ := json.Marshal(map[string]string{
			"type": FileSyncNotAllowed,
		})
		err := sendPacket(BaseSocketsMap[selfId].Conn, sendData)
		if err != nil {
			logs.Error("send data failed: ", err, "self: ", selfId)
			_ = conn.Close()
		}
		return
	}
	//region 检查参数
	targetId, hasTarget := msg["target"]
	if !hasTarget {
		logs.Error("targetId not found")
		_ = conn.Close()
		return
	}
	fileId, hasFileId := msg["fileId"]
	if !hasFileId {
		logs.Error("fileId not found")
		_ = conn.Close()
		return
	}
	strSize, hasSize := msg["size"]
	if !hasSize {
		logs.Error("size not found")
		_ = conn.Close()
		return
	}
	fileName, hasFileName := msg["fileName"]
	if !hasFileName {
		logs.Error("fileName not found")
		_ = conn.Close()
		return
	}
	userId, hasUserId := msg["userId"]
	if !hasUserId {
		logs.Error("userId not found")
	}
	targetConn, hasTargetConn := BaseSocketsMap[targetId]
	if !hasTargetConn {
		logs.Error("The file recipient is not online:", targetId)
		_ = conn.Close()
		return
	}
	//endregion
	SendFileConnMap[selfId] = &types.SocketInfo{
		Self:       types.DevInfo{DevId: selfId, Platform: platform, AppVersion: appVersion, DevName: msgInfo.selfDevName},
		Target:     &targetConn.Self,
		ConnType:   FileSync,
		CreateTime: time.Now(),
		Conn:       conn,
		Unlimited:  isUnlimitedDevice(selfId),
	}
	sendData, _ := json.Marshal(map[string]string{
		"type":     sendFile,
		"sender":   selfId,
		"size":     strSize,
		"fileName": fileName,
		"fileId":   fileId,
		"userId":   userId,
	})
	err := sendPacket(targetConn.Conn, sendData)
	if err != nil {
		logs.Error("send data failed: ", err, "self: ", selfId)
		delete(SendFileConnMap, targetId)
		_ = conn.Close()
		return
	}
}
func onRecFileTypeConnected(conn net.Conn, msgInfo connListenerInfo) {
	msg := msgInfo.msg
	selfId := msgInfo.selfId
	logs.Info("connType", recFile, ", self:", selfId)
	target, hasTarget := msg["target"]
	if !hasTarget {
		logs.Error("want to receive file but target devId not exist:", target)
		return
	}
	fileId, hasFileId := msg["fileId"]
	if !hasFileId {
		logs.Error("want to receive file but fileId not exist:", target)
		return
	}
	senderBaseConn, _ := BaseSocketsMap[target]
	sendData, _ := json.Marshal(map[string]string{
		"type":   fileReceiverConnected,
		"fileId": fileId,
	})
	err := sendPacket(senderBaseConn.Conn, sendData)
	if err != nil {
		logs.Error("send data failed: ", err, ", self: ", selfId, ", msg:", string(sendData))
		_ = conn.Close()
	} else {
		//转发连接
		senderConn, hasSenderConn := SendFileConnMap[target]
		if !hasSenderConn {
			logs.Warn("want to receive file but sender not exist, path: ", selfId+" -> "+target)
			return
		}
		rate := GetFileSyncTransferLimited(senderConn.Self.DevId)
		writer := ratelimiter.NewLimitedWriterCallBack(conn, rate, func(cnt int) {
			senderConn.TotalBytes += int64(cnt)
		})
		senderConn.LimitedWriter = writer
		_, err = io.Copy(writer, senderConn.Conn)
		_ = conn.Close()
		_ = senderConn.Conn.Close()
		if err == nil {
			logs.Info("sync file finished ", target+" -> "+selfId)
		} else {
			logs.Info("sync file failed: ", err, ", ", target+" -> "+selfId)
		}
		delete(SendFileConnMap, target)
		delete(SendFileConnMap, selfId)
	}
}
func onDataSyncTypeConnected(conn net.Conn, msgInfo connListenerInfo) {
	selfId := msgInfo.selfId
	selfSkt, exists := BaseSocketsMap[selfId]
	if !exists {
		_ = conn.Close()
		logs.Warn(selfId, "base conn not exists")
		return
	}
	targetId := msgInfo.targetId
	//数据同步
	var syncTargetKey = fmt.Sprintf("%s->%s", selfId, targetId)
	var syncSelfKey = fmt.Sprintf("%s->%s", targetId, selfId)
	sktInfo := types.SocketInfo{
		Self:       selfSkt.Self,
		ConnType:   DataSync,
		CreateTime: time.Now(),
		Conn:       conn,
		Unlimited:  isUnlimitedDevice(selfId),
	}
	DataSyncSocketsMap[syncTargetKey] = &sktInfo
	//检查对向连接
	syncSelfConn, hasConn := DataSyncSocketsMap[syncSelfKey]
	if hasConn {
		syncSelfConn.Target = &sktInfo.Self
		sktInfo.Target = &syncSelfConn.Self
		var wg sync.WaitGroup
		wg.Add(2)
		//region 比较两个速率限制的最小值
		rate1 := GetDataSyncTransferLimited(selfId)
		rate2 := GetDataSyncTransferLimited(targetId)
		rate := min(rate1, rate2)
		//endregion
		go forward(conn, syncSelfConn.Conn, selfId, targetId, &wg, rate)
		go forward(syncSelfConn.Conn, conn, targetId, selfId, &wg, rate)
		sendData, _ := json.Marshal(map[string]string{
			"type":   bothConnected,
			"sender": targetId,
		})
		err := sendPacket(conn, sendData)
		if err != nil {
			logs.Error("send data failed: ", err, ", self: ", selfId)
			_ = conn.Close()
			delete(DataSyncSocketsMap, selfId+"->"+targetId)
			delete(DataSyncSocketsMap, targetId+"->"+selfId)
			return
		}
		err = sendPacket(syncSelfConn.Conn, sendData)
		if err != nil {
			logs.Error("send data failed: ", err, ", self: ", selfId)
			_ = syncSelfConn.Conn.Close()
			delete(DataSyncSocketsMap, selfId+"->"+targetId)
			delete(DataSyncSocketsMap, targetId+"->"+selfId)
			return
		}
		wg.Wait()
		logs.Info("wait finished")
		sendData, _ = json.Marshal(map[string]string{
			"type":   forwardReady,
			"sender": targetId,
		})
		_ = sendPacket(syncSelfConn.Conn, sendData)
	} else {
		//若没有对向连接，则请求
		sendData, _ := json.Marshal(map[string]string{
			"type":   requestConnect,
			"sender": selfId,
		})
		targetAliveSocket, targetIsOnline := BaseSocketsMap[targetId]
		if !targetIsOnline {
			delete(DataSyncSocketsMap, syncTargetKey)
			_ = conn.Close()
			return
		}
		_ = sendPacket(targetAliveSocket.Conn, sendData)
	}
}

func forward(source net.Conn, target net.Conn, sourceId string, targetId string, wg *sync.WaitGroup, rate int) {
	//region 等待回应
	packetReader := ratelimiter.NewPacketReader(source)
	data, err := packetReader.ReadPacket()
	if err != nil {
		//接收失败
		logs.Error("read forward wait response client data failed:", err)
		_ = source.Close()
		_ = target.Close()
		return
	}
	logs.Info("received:", string(data))
	var msg map[string]string
	{
	}
	err = json.Unmarshal(data, &msg)
	//解析失败
	if err != nil {
		logs.Error("forward wait response data unmarshal failed:", err)
		_ = source.Close()
		_ = target.Close()
		return
	}
	//endregion
	dataType, _ := msg["type"]
	if dataType == bothConnected {
		wg.Done()
	} else {
		logs.Error("invalid forward response data type:", dataType)
		_ = source.Close()
		_ = target.Close()
	}
	//双方已经建立双向转发
	logs.Info("forward ready:", "device1:", sourceId, "device2:", targetId)
	_, err = 0, nil
	targetWriter := ratelimiter.NewLimitedWriterCallBack(target, rate, func(cnt int) {
		DataSyncSocketsMap[sourceId+"->"+targetId].TotalBytes += int64(cnt)
	})
	DataSyncSocketsMap[sourceId+"->"+targetId].LimitedWriter = targetWriter
	_, err = io.Copy(targetWriter, source)
	if err == nil {
		logs.Info("forward finished.", "source:", sourceId, "target:", targetId)
	}
	logs.Error("forward failed:", err, ".", "source:", sourceId, "target:", targetId)
	//断开连接
	_ = source.Close()
	_ = target.Close()
	//清除连接
	delete(DataSyncSocketsMap, sourceId+"->"+targetId)
	delete(DataSyncSocketsMap, targetId+"->"+sourceId)
}
func onBaseConnReceived(selfId string, data []byte) {
	var msg map[string]string
	{
	}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		logs.Error("forward wait response data unmarshal failed:", err)
		return
	}
	msgType := msg["type"]
	switch msgType {
	case cancelSendFile:
		targetId := msg["targetId"]
		delete(SendFileConnMap, targetId)
	}
}
