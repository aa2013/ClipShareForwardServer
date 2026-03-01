package types

import (
	"clipshare/ratelimiter"
	"clipshare/utils"
	"net"
	"time"
)

type SocketInfo struct {
	Self          DevInfo
	Target        *DevInfo
	ConnType      string
	CreateTime    time.Time
	Conn          net.Conn
	speed         int64
	LastBytes     int64
	TotalBytes    int64
	Unlimited     bool
	LimitedWriter *ratelimiter.RateLimitWriter
	//only base conn
	AccessKey *string
	//only base conn
	KeyFirstUseAt *time.Time
	//only base conn
	PlanType *PlanTypeDto
}

func (skt *SocketInfo) UpdateSnapshot() {
	skt.speed = skt.TotalBytes - skt.LastBytes
	skt.LastBytes = skt.TotalBytes
}
func (skt *SocketInfo) UpdateRateLimit(limit int) {
	if skt.LimitedWriter == nil {
		return
	}
	skt.LimitedWriter.UpdateLimit(limit)
}

func (skt *SocketInfo) ToDto() ConnectionStatusDto {
	return ConnectionStatusDto{
		Self:             skt.Self,
		Target:           skt.Target,
		ConnType:         skt.ConnType,
		CreateTime:       skt.CreateTime.Format("2006-01-02 15:04:05"),
		Speed:            utils.IntToSizeStr(skt.speed) + "/s",
		Unlimited:        skt.Unlimited,
		TransferredBytes: utils.IntToSizeStr(skt.TotalBytes),
	}
}
