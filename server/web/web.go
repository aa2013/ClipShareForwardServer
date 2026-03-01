package web

import (
	"clipshare/db"
	"clipshare/forward"
	"clipshare/types"
	"clipshare/utils"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"time"
	"embed"
	"io/fs"
)

const authKey = "Authorization"
const version = "1.1.1"
//go:embed all:dist/*
var embeddedFiles embed.FS
var lastOperationTime *time.Time = nil
var loginToken string

func StartWebServer() {
	r := gin.Default()
	r.Use(GlobalExceptionMiddleware(), Cors())

    distFS, _ := fs.Sub(embeddedFiles, "dist")
    fileServer := http.FileServer(http.FS(distFS))
    r.NoRoute(func(c *gin.Context) {
        fileServer.ServeHTTP(c.Writer, c.Request)
    })

    api := r.Group("/api")
	api.POST("/login", login)
	api.GET("/version", getVersion)

	authorized := api.Group("/admin", AuthMiddleware(false))
	authorized.POST("/logout", logout)
	authorized.GET("/connectionStatus", getConnectionStatus)
	authorized.GET("/charts", getChartsData)
	authorized.POST("/forcedDisconnection", forcedDisconnection)
	authorized.GET("/config", getConfig)
	authorized.POST("/config/update", updateConfig)
	authorized.GET("/logs", getLogs)
	authorized.POST("/plan/add", addPlan)
	authorized.POST("/plan/edit", editPlan)
	authorized.GET("/plan/list", getPlans)
	authorized.POST("/plan/updateStatus", updatePlanStatus)
	authorized.POST("/plan/generatePlanKeys", generatePlanKeys)
	authorized.POST("/planKeys/list", getPlanKeys)
	authorized.POST("/planKeys/updateStatus", updatePlanKeyStatus)
	authorized.GET("/planKeys/verify", verifyKey)
	// Start web forwardServer
	_ = r.Run(fmt.Sprintf(":%d", *types.AppConfig.Web.Port))
}
func login(c *gin.Context) {
	var loginUser types.LoginDto
	types.BindJsonData(c, &loginUser)
	var adminUser = types.AppConfig.Web.Admin
	if loginUser.Username != *adminUser.Username || loginUser.Password != *adminUser.Password {
		errorResult(c, http.StatusUnauthorized, "Incorrect username or password", nil)
		return
	}
	now := time.Now()
	lastOperationTime = &now
	rand.Seed(time.Now().UnixNano())
	number := rand.Int63()
	t := loginUser.Username + strconv.FormatInt(number, 10) + loginUser.Password
	loginToken = fmt.Sprintf("%x", md5.Sum([]byte(t)))
	successResult(c, gin.H{"token": loginToken})
}

func getVersion(c *gin.Context) {
	successResult(c, gin.H{"version": version})
}
func logout(c *gin.Context) {
	lastOperationTime = nil
	loginToken = ""
	successResult(c, nil)
}
func getConnectionStatus(c *gin.Context) {
	baseConnections := make([]types.ConnectionStatusDto, 0, len(forward.BaseSocketsMap))
	for _, v := range forward.BaseSocketsMap {
		baseConnections = append(baseConnections, v.ToDto())
	}
	dataSyncConnections := make([]types.ConnectionStatusDto, 0, len(forward.DataSyncSocketsMap))
	for _, v := range forward.DataSyncSocketsMap {
		dataSyncConnections = append(dataSyncConnections, v.ToDto())
	}
	fileSyncConnections := make([]types.ConnectionStatusDto, 0, len(forward.SendFileConnMap))
	for _, v := range forward.SendFileConnMap {
		fileSyncConnections = append(fileSyncConnections, v.ToDto())
	}
	//createTime desc
	sort.Slice(baseConnections, func(i, j int) bool {
		return baseConnections[i].CreateTime >= baseConnections[j].CreateTime
	})
	sort.Slice(dataSyncConnections, func(i, j int) bool {
		return dataSyncConnections[i].CreateTime >= dataSyncConnections[j].CreateTime
	})
	sort.Slice(fileSyncConnections, func(i, j int) bool {
		return fileSyncConnections[i].CreateTime >= fileSyncConnections[j].CreateTime
	})
	successResult(c, gin.H{
		forward.Base:     baseConnections,
		forward.DataSync: dataSyncConnections,
		forward.FileSync: fileSyncConnections,
	})
}
func forcedDisconnection(c *gin.Context) {
	var item types.ForcedDisconnectionDto
	types.BindJsonData(c, &item)
	var skt *types.SocketInfo
	var hasConn bool
	switch item.ConnType {
	case forward.Base:
		skt, hasConn = forward.BaseSocketsMap[item.Key]
	case forward.FileSync:
		skt, hasConn = forward.SendFileConnMap[item.Key]
	case forward.DataSync:
		skt, hasConn = forward.DataSyncSocketsMap[item.Key]
	default:
		panic("UnSupport connType: " + item.ConnType)
	}
	if !hasConn {
		successResult(c, false)
	} else {
		_ = skt.Conn.Close()
		successResult(c, true)
	}
}
func getConfig(c *gin.Context) {
	successResult(c, types.AppConfig.ToDto())
}
func updateConfig(c *gin.Context) {
	var newCfg types.ConfigDto
	types.BindJsonData(c, &newCfg)
	if newCfg.LoginExpiredSeconds != nil && *newCfg.LoginExpiredSeconds <= 0 {
		panic(fmt.Sprintf("loginExpiredSeconds can not <= 0"))
	}
	if newCfg.FileTransferRateLimit != nil && *newCfg.FileTransferRateLimit <= 0 {
		panic(fmt.Sprintf("fileTransferRateLimit can not <= 0"))
	}
	if newCfg.FileTransferEnabled != nil {
		types.AppConfig.Forward.FileTransferLimit.Enabled = *newCfg.FileTransferEnabled
	}
	if newCfg.FileTransferRateLimit != nil {
		types.AppConfig.Forward.FileTransferLimit.Rate = newCfg.FileTransferRateLimit
	}
	if newCfg.LoginExpiredSeconds != nil {
		types.AppConfig.Web.LoginExpiredSeconds = *newCfg.LoginExpiredSeconds
	}
	if newCfg.UnlimitedDevices != nil {
		types.AppConfig.Forward.UnlimitedDevices = *newCfg.UnlimitedDevices
	}
	toPrivateMode := false
	if newCfg.PublicMode != nil {
		if types.AppConfig.PublicMode && *newCfg.PublicMode == false {
			toPrivateMode = true
		}
	}
	if newCfg.PublicMode != nil {
		types.AppConfig.PublicMode = *newCfg.PublicMode
	}
	updateLogBuffer := false
	if newCfg.Log != nil {
		if newCfg.Log.MemoryBufferSize < 10 {
			panic("MemoryBufferSize not be less than 10")
		} else {
			types.AppConfig.Log.MemoryBufferSize = newCfg.Log.MemoryBufferSize
			updateLogBuffer = true
		}
	}
	err := types.AppConfig.Save("./data/config.yaml")
	if err != nil {
		panic(fmt.Sprintf("Save to memory successful, but error message failed to save to file: %v, ", err))
	}
	if toPrivateMode {
		forward.SwitchToPrivateMode()
	}
	if updateLogBuffer {
		utils.LogUtil.ResizeBuffer(newCfg.Log.MemoryBufferSize)
	}
	successResult(c, true)
}
func getChartsData(c *gin.Context) {
	successResult(c, forward.ChartData)
}
func getLogs(c *gin.Context) {
	timeLayout := "2006-01-02 15:04:05.000"
	memoryLogs := utils.LogUtil.GetAllMemoryLogs()
	beginParam := c.Query("begin")
	if beginParam != "" {
		beginTime, err := time.Parse(timeLayout, beginParam)
		if err != nil {
			panic(err)
		}
		//println(beginTime.Format("2006-01-02 15:04:05.000"))
		var filterLogs []utils.Log
		for _, v := range memoryLogs {
			logTime, _ := time.Parse(timeLayout, v.Time.Format(timeLayout))
			if logTime.After(beginTime) {
				filterLogs = append(filterLogs, v)
			}
		}
		memoryLogs = filterLogs
	}
	logs := make([]types.LogDto, len(memoryLogs))
	for i, v := range memoryLogs {
		logs[i] = types.LogDto{
			Log:  v.Log,
			Time: v.Time.Format(timeLayout),
		}
	}
	successResult(c, logs)
}
func addPlan(c *gin.Context) {
	var planDto types.PlanTypeDto
	types.BindJsonData(c, &planDto)
	err := db.AddPlan(db.ToPlanType(planDto))
	if err != nil {
		panic(err)
	}
	forward.UpdateConnPlanTypeCache(planDto)
	successResult(c, true)
}
func editPlan(c *gin.Context) {
	var planDto types.PlanTypeDto
	types.BindJsonData(c, &planDto)
	err := db.UpdatePlan(db.ToPlanType(planDto))
	if err != nil {
		panic(err)
	}
	successResult(c, true)
}
func getPlans(c *gin.Context) {
	plans, err := db.GetAllPlans()
	if err != nil {
		panic(err)
	}
	res := make([]types.PlanTypeDto, 0, len(plans))
	for _, v := range plans {
		res = append(res, v.ToPlanTypeDto())
	}
	successResult(c, res)
}
func updatePlanStatus(c *gin.Context) {
	var ps types.PlanStatusDto
	types.BindJsonData(c, &ps)
	err := db.UpdatePlanStatus(*ps.Id, *ps.Status)
	if err != nil {
		panic(err)
	}
	successResult(c, true)
}
func generatePlanKeys(c *gin.Context) {
	var ps types.GeneratePlanKeysDto
	types.BindJsonData(c, &ps)
	err := db.GeneratePlanKeys(*ps.Id, *ps.Size)
	if err != nil {
		panic(err)
	}
	successResult(c, true)
}
func getPlanKeys(c *gin.Context) {
	var params types.PlanKeysSearchDto
	types.BindJsonData(c, &params)
	plans, err := db.GetAllPlans()
	planNamesMap := map[string]string{}
	for _, plan := range plans {
		planNamesMap[plan.Id] = plan.Name
	}
	if err != nil {
		panic(err)
	}

	total, keys, err := db.GetPlanKeysPageData(params)
	if err != nil {
		panic(err)
	}
	if total == 0 {
		params.PageNum = 0
	}
	rows := make([]types.PlanKeyDto, 0, len(keys))
	for _, v := range keys {
		rows = append(rows, db.ToPlanKeyDto(v, planNamesMap[v.PlanId]))
	}
	successResult(c, types.PageDataDto[[]types.PlanKeyDto]{
		Rows:       rows,
		Total:      total,
		PageParams: params.PageParams,
	})
}

func updatePlanKeyStatus(c *gin.Context) {
	var pks types.PlanKeyStatusDto
	types.BindJsonData(c, &pks)
	err := db.UpdatePlanKeyStatus(*pks.Id, *pks.Status)
	if err != nil {
		panic(err)
	}
	pk := db.GetPlanKeyById(*pks.Id)
	if *pks.Status == false {
		forward.StopPlanKeyConn(pk.Key)
	}
	successResult(c, true)
}
func verifyKey(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		panic("Key does not exist")
	}
	pt := db.VerifyKey(key)
	if pt == nil {
		successResult(c, nil)
	}
	successResult(c, types.PlanTypeDto{
		Id:          pt.Id,
		Name:        pt.Name,
		Rate:        pt.Rate,
		Lifespan:    pt.Lifespan,
		DeviceLimit: pt.DeviceLimit,
		Remark:      pt.Remark,
		Enable:      pt.Enable,
	})
}
