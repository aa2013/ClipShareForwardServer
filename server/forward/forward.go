package forward

import (
	"bytes"
	"clipshare/ratelimiter"
	"clipshare/types"
	"clipshare/utils"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"slices"
	"time"
)

// 基础连接 selfId->SocketInfo
var lastBaseConnCnt = 0
var BaseSocketsMap = make(map[string]*types.SocketInfo)

// 数据同步连接
var lastDataSyncConnCnt = 0
var DataSyncSocketsMap = make(map[string]*types.SocketInfo)

// 文件传输连接
var lastSendFileConnCnt = 0
var SendFileConnMap = make(map[string]*types.SocketInfo)

var ChartData = types.ChartData{
	NetSpeed: types.NetWorkChartData{
		DataSync: 0,
		FileSync: 0,
	},
	Traffic: 0,
	ConnectionCnt: types.ConnectionChartData{
		DataSyncCnt: 0,
		FileSyncCnt: 0,
	},
}

const (
	alreadyConnected      = "alreadyConnected"
	forwardReady          = "forwardReady"
	bothConnected         = "bothConnected"
	requestConnect        = "requestConnect"
	sendFile              = "sendFile"
	ping                  = "ping"
	cancelSendFile        = "cancelSendFile"
	recFile               = "recFile"
	fileReceiverConnected = "fileReceiverConnected"
	Base                  = "base"
	check                 = "check"
	DataSync              = "dataSync"
	FileSync              = "fileSync"
	FileSyncNotAllowed    = "fileSyncNotAllowed"
)

var logs *utils.LogManager

// StartForwardServer 启动服务
func StartForwardServer() {
	defer utils.OnDeferRecover()
	logs = utils.LogUtil
	netListen, err := net.Listen("tcp", fmt.Sprintf(":%d", *types.AppConfig.Forward.Port))
	if err != nil {
		logs.Error("startForwardServer failed:", err)
	}
	defer func(netListen net.Listener) {
		err := netListen.Close()
		if err != nil {
			logs.Error("forwardServer close failed:", err)
		}
	}(netListen)
	go startSocketTrafficTicker()
	go startTTLSticker()
	go sendHeartbeatPeriod()
	//等待客户端连接
	for {
		conn, err := netListen.Accept()
		if err != nil {
			logs.Error("socket accept failed:", err)
			continue
		}
		go clientListener(conn)
	}
}
func clientListener(conn net.Conn) {
	defer utils.OnDeferRecover()
	logs.Info("client connected:", conn.RemoteAddr().String())
	//region 解析初始数据
	packetReader := ratelimiter.NewPacketReader(conn)
	data, err := packetReader.ReadPacket()
	if err != nil {
		//接收失败
		logs.Error("read client data failed:", err)
		return
	}
	logs.Info("received:", string(data))
	var msg map[string]string
	{
	}
	err = json.Unmarshal(data, &msg)
	//解析失败
	if err != nil {
		logs.Error("data unmarshal failed:", err, data)
		return
	}
	//endregion

	//处理初始数据
	connType := msg["connType"]
	selfId := msg["self"]
	platform := msg["platform"]
	appVersion := msg["appVersion"]
	targetId := msg["target"]
	selfDevName := msg["devName"]
	if selfId == "" {
		logs.Error("selfId not found")
		_ = conn.Close()
	}
	msgInfo := connListenerInfo{
		msg:         msg,
		selfId:      selfId,
		platform:    platform,
		appVersion:  appVersion,
		targetId:    targetId,
		selfDevName: selfDevName,
	}
	switch connType {
	case check:
		{
			checkKey(conn, msgInfo)
			_ = conn.Close()
		}
	case Base:
		onBaseTypeConnected(conn, msgInfo, packetReader)
	case sendFile:
		onSendFileTypeConnected(conn, msgInfo)
	case recFile:
		onRecFileTypeConnected(conn, msgInfo)
	default:
		onDataSyncTypeConnected(conn, msgInfo)
	}
}
func startSocketTrafficTicker() {
	scanSockets()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop() // 程序结束前停止 Ticker
	for {
		<-ticker.C    // 等待下一个 1 秒间隔
		scanSockets() //扫描 sockets
	}
}
func scanSockets() {
	var dataSyncSpeed int64
	var fileSyncSpeed int64
	lastDataSyncBytesMap := make(map[string]int64)
	lastFileSyncBytesMap := make(map[string]int64)
	curDataSyncBytesMap := make(map[string]int64)
	curFileSyncBytesMap := make(map[string]int64)
	for _, v := range DataSyncSocketsMap {
		//记录上次的流量
		if _, hasData := lastDataSyncBytesMap[v.Self.DevId]; hasData {
			lastDataSyncBytesMap[v.Self.DevId] += v.LastBytes
		} else {
			lastDataSyncBytesMap[v.Self.DevId] = v.LastBytes
		}
		//更新当前流量
		v.UpdateSnapshot()
		//记录当前流量
		if _, hasData := curDataSyncBytesMap[v.Self.DevId]; hasData {
			curDataSyncBytesMap[v.Self.DevId] += v.TotalBytes
		} else {
			curDataSyncBytesMap[v.Self.DevId] = v.TotalBytes
		}
		dataSyncSpeed += curDataSyncBytesMap[v.Self.DevId] - lastDataSyncBytesMap[v.Self.DevId]
	}
	for _, v := range SendFileConnMap {
		//记录上次的流量
		if _, hasData := lastFileSyncBytesMap[v.Self.DevId]; hasData {
			lastFileSyncBytesMap[v.Self.DevId] += v.LastBytes
		} else {
			lastFileSyncBytesMap[v.Self.DevId] = v.LastBytes
		}
		v.UpdateSnapshot()
		//记录当前流量
		if _, hasData := curFileSyncBytesMap[v.Self.DevId]; hasData {
			curFileSyncBytesMap[v.Self.DevId] += v.TotalBytes
		} else {
			curFileSyncBytesMap[v.Self.DevId] = v.TotalBytes
		}
		fileSyncSpeed += curFileSyncBytesMap[v.Self.DevId] - lastFileSyncBytesMap[v.Self.DevId]
	}
	//计算总的
	var traffic int64
	for _, v := range BaseSocketsMap {
		totalBytes := curFileSyncBytesMap[v.Self.DevId] + curDataSyncBytesMap[v.Self.DevId]
		v.TotalBytes = totalBytes
		traffic += totalBytes
		v.UpdateSnapshot()
	}
	ChartData.Traffic = uint64(traffic)
	if len(BaseSocketsMap) != 0 {
		ChartData.NetSpeed.FileSync = uint64(fileSyncSpeed)
		ChartData.NetSpeed.DataSync = uint64(dataSyncSpeed)
	}
	curBaseConnCnt := len(BaseSocketsMap)
	curDataSyncConnCnt := len(DataSyncSocketsMap)
	curSendFileConnCnt := len(SendFileConnMap)

	if lastBaseConnCnt != curBaseConnCnt || lastDataSyncConnCnt != curDataSyncConnCnt || lastSendFileConnCnt != curSendFileConnCnt {
		// 执行代码
		lastDataSyncConnCnt = curDataSyncConnCnt
		lastBaseConnCnt = curBaseConnCnt
		lastSendFileConnCnt = curSendFileConnCnt
		ChartData.ConnectionCnt.DataSyncCnt = uint64(curDataSyncConnCnt)
		ChartData.ConnectionCnt.FileSyncCnt = uint64(curSendFileConnCnt)
		ChartData.ConnectionCnt.BaseCnt = uint64(curBaseConnCnt)
	}
}
func startTTLSticker() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop() // 程序结束前停止 Ticker
	for {
		<-ticker.C // 等待下一个间隔
		scanKeyTTL()
	}
}
func sendHeartbeatPeriod() {
	sendData, _ := json.Marshal(map[string]string{
		"type": ping,
	})
	ticker := time.NewTicker(30 * time.Second)

	defer ticker.Stop() // 程序结束前停止 Ticker
	for {
		<-ticker.C // 等待下一个间隔
		baseConns := make([]*types.SocketInfo, 0, len(BaseSocketsMap))
		for _, skt := range BaseSocketsMap {
			baseConns = append(baseConns, skt)
		}
		for _, skt := range baseConns {
			if skt == nil {
				continue
			}
			_ = sendPacket(skt.Conn, sendData)
		}
	}
}
func scanKeyTTL() {
	now := time.Now()
	sendData, _ := json.Marshal(map[string]string{
		"type":   check,
		"result": "The current key has expired",
	})
	for _, baseConn := range BaseSocketsMap {
		if baseConn.KeyFirstUseAt == nil {
			continue
		}
		if baseConn.PlanType.Lifespan == nil {
			continue
		}
		lifeSpanSeconds := *baseConn.PlanType.Lifespan
		duration := now.Sub(*baseConn.KeyFirstUseAt)
		//计算从第一次使用到现在的秒数
		offsetSeconds := uint64(duration.Seconds())
		if offsetSeconds >= lifeSpanSeconds {
			//超时
			_ = sendPacket(baseConn.Conn, sendData)
			_ = baseConn.Conn.Close()
		}
	}
}
func SwitchToPrivateMode() {
	sendData, _ := json.Marshal(map[string]string{
		"type":   check,
		"result": "Service switched to private mode, please reconnect using key",
	})
	for _, baseConn := range BaseSocketsMap {
		if !baseConn.Unlimited {
			_ = sendPacket(baseConn.Conn, sendData)
			_ = baseConn.Conn.Close()
		}
	}
}
func UpdateConnPlanTypeCache(pt types.PlanTypeDto) {
	for _, baseConn := range BaseSocketsMap {
		if baseConn.PlanType == nil || baseConn.PlanType.Id != pt.Id {
			continue
		}
		baseConn.PlanType = &pt
	}
}
func StopPlanKeyConn(key string) {
	sendData, _ := json.Marshal(map[string]string{
		"type":   check,
		"result": "The current key has been deactivated",
	})
	for _, baseConn := range BaseSocketsMap {
		if baseConn.AccessKey == nil {
			continue
		}
		if *baseConn.AccessKey == key {
			_ = sendPacket(baseConn.Conn, sendData)
			_ = baseConn.Conn.Close()
		}
	}
}
func UpdateRateLimitConfig() {
	config := types.AppConfig.Forward

	sktList := make([]*types.SocketInfo, 0, len(DataSyncSocketsMap)+len(SendFileConnMap))
	for _, skt := range DataSyncSocketsMap {
		sktList = append(sktList, skt)
	}
	for _, skt := range SendFileConnMap {
		sktList = append(sktList, skt)
	}
	//白名单设备
	unlimitedDevIds := types.AppConfig.GetUnlimitedDeviceIds()
	for _, item := range config.UnlimitedDevices {
		devId := item.Id
		unlimitedDevIds = append(unlimitedDevIds, devId)
	}
	newRate := *config.FileTransferLimit.Rate
	for _, skt := range BaseSocketsMap {
		selfId := skt.Self.DevId
		skt.Unlimited = slices.Contains(unlimitedDevIds, selfId)
	}
	for _, skt := range sktList {
		selfId := skt.Self.DevId
		targetId := skt.Target.DevId
		//更新白名单状态
		skt.Unlimited = slices.Contains(unlimitedDevIds, selfId)
		//判断是否有一边是白名单
		unlimited := slices.Contains(unlimitedDevIds, selfId) || slices.Contains(unlimitedDevIds, targetId)
		//判断是否需要限速: 双方都不是白名单且启用了速率限制
		unlimited = unlimited || !config.FileTransferLimit.Enabled
		if unlimited {
			if skt.LimitedWriter != nil {
				skt.LimitedWriter.UpdateLimit(0)
			}
			if skt.LimitedWriter != nil {
				skt.LimitedWriter.UpdateLimit(0)
			}
		} else {
			if skt.LimitedWriter != nil {
				skt.LimitedWriter.UpdateLimit(newRate * 1024)
			}
			if skt.LimitedWriter != nil {
				skt.LimitedWriter.UpdateLimit(newRate * 1024)
			}
		}
	}
}
func sendPacket(conn net.Conn, data []byte) error {

	// region 计算包头信息
	totalSize := uint32(len(data))
	packetSize := uint16(len(data))
	//中转程序发送量很小，一个包就够了
	totalPackets := uint16(1)
	seq := uint16(1)
	//endregion

	//region 构造包头
	var headerBuf bytes.Buffer
	if err := binary.Write(&headerBuf, binary.BigEndian, totalSize); err != nil {
		return err
	}
	if err := binary.Write(&headerBuf, binary.BigEndian, packetSize); err != nil {
		return err
	}
	if err := binary.Write(&headerBuf, binary.BigEndian, totalPackets); err != nil {
		return err
	}
	if err := binary.Write(&headerBuf, binary.BigEndian, seq); err != nil {
		return err
	}
	//endregion

	// 将包头和包体拼接成完整的数据包
	fullPacket := append(headerBuf.Bytes(), data...)

	// 发送数据包
	_, err := conn.Write(fullPacket)
	return err
}
