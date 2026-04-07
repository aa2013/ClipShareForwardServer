package types

import (
	"clipshare/ratelimiter"
	"clipshare/utils"
	"net"
	"sync"
	"time"
)

type SocketInfo struct {
	mu            sync.RWMutex
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
	skt.mu.Lock()
	defer skt.mu.Unlock()
	skt.speed = skt.TotalBytes - skt.LastBytes
	skt.LastBytes = skt.TotalBytes
}
func (skt *SocketInfo) AddTotalBytes(cnt int64) {
	skt.mu.Lock()
	defer skt.mu.Unlock()
	skt.TotalBytes += cnt
}
func (skt *SocketInfo) SetTotalBytes(total int64) {
	skt.mu.Lock()
	defer skt.mu.Unlock()
	skt.TotalBytes = total
}
func (skt *SocketInfo) BytesSnapshot() (last int64, total int64) {
	skt.mu.RLock()
	defer skt.mu.RUnlock()
	return skt.LastBytes, skt.TotalBytes
}
func (skt *SocketInfo) UpdateRateLimit(limit int) {
	skt.mu.RLock()
	writer := skt.LimitedWriter
	skt.mu.RUnlock()
	if writer == nil {
		return
	}
	writer.UpdateLimit(limit)
}

func (skt *SocketInfo) ToDto() ConnectionStatusDto {
	skt.mu.RLock()
	defer skt.mu.RUnlock()
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
