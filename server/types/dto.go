package types

import (
	"github.com/gin-gonic/gin"
)

func BindJsonData[T any](c *gin.Context, dto *T) {
	// 使用 ShouldBind 进行绑定，自动绑定到结构体字段
	if err := c.ShouldBindJSON(&dto); err != nil {
		panic(err)
	}
}
func BindQueryData[T any](c *gin.Context, dto *T) {
	// 使用 ShouldBind 进行绑定，自动绑定到结构体字段
	if err := c.ShouldBindQuery(&dto); err != nil {
		panic(err)
	}
}

type LoginDto struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type ConnectionStatusDto struct {
	Self             DevInfo  `json:"self"`
	Target           *DevInfo `json:"target"`
	ConnType         string   `json:"connType"`
	CreateTime       string   `json:"createTime"`
	Speed            string   `json:"speed"`
	Unlimited        bool     `json:"unlimited"`
	TransferredBytes string   `json:"transferredBytes"`
}

type ForcedDisconnectionDto struct {
	ConnType string `form:"connType"`
	Key      string `form:"key"`
}
type ConfigDto struct {
	LoginExpiredSeconds   *int              `json:"loginExpiredSeconds"`
	UnlimitedDevices      *[]DeviceBaseInfo `json:"unlimitedDevices" binding:"omitempty,dive"`
	FileTransferEnabled   *bool             `json:"fileTransferEnabled"`
	FileTransferRateLimit *int              `json:"fileTransferRateLimit"`
	Log                   *LogConfig        `json:"log"`
	PublicMode            *bool             `json:"publicMode"`
}

// ConnectionCnt 连接数记录
type ConnectionCnt struct {
	BaseCnt     uint
	DataSyncCnt uint
	FileSyncCnt uint
	Time        string
}

// NetSpeed 网速记录
type NetSpeed struct {
	FileSyncSpeed uint64
	DataSyncSpeed uint64
	Time          string
}

// DevTraffic 设备流量记录
type DevTraffic struct {
	DevId   string
	Traffic uint64
	Time    string
}
type NetWorkChartData struct {
	DataSync uint64 `json:"dataSync"`
	FileSync uint64 `json:"fileSync"`
}
type ConnectionChartData struct {
	BaseCnt     uint64 `json:"baseCnt"`
	DataSyncCnt uint64 `json:"dataSyncCnt"`
	FileSyncCnt uint64 `json:"fileSyncCnt"`
}
type TrafficChartData struct {
	BaseCnt     uint64 `json:"baseCnt"`
	DataSyncCnt uint64 `json:"dataSyncCnt"`
	FileSyncCnt uint64 `json:"fileSyncCnt"`
}
type ChartData struct {
	NetSpeed      NetWorkChartData    `json:"netSpeed"`
	Traffic       uint64              `json:"traffic"`
	ConnectionCnt ConnectionChartData `json:"connectionCnt"`
}
type LogDto struct {
	Log  string `json:"log"`
	Time string `json:"time"`
}

type PlanTypeDto struct {
	Id          string  `json:"id"`
	Name        string  `json:"name" binding:"required"`
	Rate        *uint   `json:"rate" binding:"omitempty,min=1"`
	Lifespan    *uint64 `json:"lifespan" binding:"omitempty,min=1"`
	DeviceLimit *uint   `json:"deviceLimit" binding:"omitempty,min=1"`
	Remark      *string `json:"remark"`
	Enable      bool    `json:"enable"`
}

type PlanStatusDto struct {
	Id     *string `json:"id" binding:"required"`
	Status *bool   `json:"status" binding:"required"`
}
type PlanKeyStatusDto struct {
	Id     *int  `json:"id" binding:"required"`
	Status *bool `json:"status" binding:"required"`
}
type GeneratePlanKeysDto struct {
	Id   *string `json:"id" binding:"required"`
	Size *uint   `json:"size" binding:"required,min=1"`
}
type PageDataDto[T interface{}] struct {
	Rows  T     `json:"rows"`
	Total int64 `json:"total"`
	PageParams
}
type PageParams struct {
	PageNum  int `json:"pageNum" form:"pageNum" binding:"required,min=1"`
	PageSize int `json:"pageSize" form:"pageSize" binding:"required,min=1"`
}
type PlanKeysSearchDto struct {
	PlanId string        `json:"planId" form:"planId"`
	Sorts  []SortOptions `json:"sorts" form:"sorts"`
	PageParams
}
type SortOptions struct {
	Key   string `json:"key" form:"key"`     //column name
	Order string `json:"order" form:"order"` //desc or asc
}
type PlanKeyDto struct {
	Id        uint    `json:"id"`        //id
	Key       string  `json:"key"`       //密钥
	PlanId    string  `json:"planId"`    //套餐id
	PlanName  string  `json:"planName"`  //套餐名称
	UseAt     *string `json:"useAt"`     //开始使用时间
	Enable    bool    `json:"enable"`    //是否启用
	Remark    string  `json:"remark"`    //备注
	CreatedAt string  `json:"createdAt"` //创建时间
}
