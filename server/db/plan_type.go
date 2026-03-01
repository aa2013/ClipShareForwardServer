package db

import (
	"clipshare/types"
	"clipshare/utils"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PlanType 套餐类型
type PlanType struct {
	Id          string  `gorm:"type:text;primaryKey"`
	Name        string  `gorm:"not null"` //名称
	AesKey      string  `gorm:"not null"`
	Rate        *uint   `gorm:"null"`         //限速速率（KB/s）
	Lifespan    *uint64 `gorm:"null"`         //有效期，秒数，第一次连接时开始计时
	DeviceLimit *uint   `gorm:"null"`         //同时使用设备数量限制
	Remark      *string `gorm:"null"`         //备注
	Enable      bool    `gorm:"default:true"` //是否启用
	gorm.Model
}

func ToPlanType(p types.PlanTypeDto) PlanType {
	return PlanType{
		Id:          p.Id,
		Name:        p.Name,
		Rate:        p.Rate,
		Lifespan:    p.Lifespan,
		DeviceLimit: p.DeviceLimit,
		Remark:      p.Remark,
	}
}
func (p *PlanType) ToPlanTypeDto() types.PlanTypeDto {
	return types.PlanTypeDto{
		Id:          p.Id,
		Name:        p.Name,
		Rate:        p.Rate,
		Lifespan:    p.Lifespan,
		DeviceLimit: p.DeviceLimit,
		Remark:      p.Remark,
		Enable:      p.Enable,
	}
}
func GetAllPlans() ([]PlanType, error) {
	checkDb()
	var plans []PlanType
	result := AppDb.Order("enable desc").Find(&plans)
	return plans, result.Error
}
func getPlanInfo(id string) *PlanType {
	checkDb()
	var plan PlanType
	result := AppDb.First(&plan, "id = ?", id)
	if result.Error != nil {
		return nil
	}
	return &plan
}
func getPlanInfoByKey(key string) *PlanType {
	checkDb()
	pk := GetPlanKey(key)
	if pk == nil {
		return nil
	}
	return getPlanInfo(pk.PlanId)
}
func AddPlan(plan PlanType) error {
	checkDb()
	str, err := utils.GenRandomString(32)
	if err != nil {
		return err
	}
	plan.Id = uuid.New().String()
	plan.AesKey = str
	result := AppDb.Create(&plan)
	return result.Error
}
func UpdatePlanStatus(id string, enable bool) error {
	checkDb()
	result := AppDb.Model(&PlanType{}).Where("id = ?", id).Update("enable", enable)
	return result.Error
}
func UpdatePlan(plan PlanType) error {
	checkDb()
	result := AppDb.Model(&PlanType{}).Where("id = ?", plan.Id).Updates(plan)
	return result.Error
}

// GeneratePlanKeys planid--aes(uuid)--sha256(uuid)[:8]
func GeneratePlanKeys(planId string, size uint) error {
	checkDb()
	var plan PlanType
	result := AppDb.First(&plan, "id = ?", planId)
	if result.Error != nil {
		return result.Error
	}
	if plan.Enable == false {
		return errors.New("this plan has been deactivate")
	}
	keyList := make([]*PlanKey, 0, size)
	aesKey := plan.AesKey
	for i := uint(0); i < size; i++ {
		id := uuid.New().String()
		encryptedBytes := utils.AesEncryptCBC([]byte(id), []byte(aesKey))
		encrypted := base64.StdEncoding.EncodeToString(encryptedBytes)
		hash := sha256.Sum256(encryptedBytes)
		checksum := fmt.Sprintf("%x", hash)[:8]
		newKey := fmt.Sprintf("%s--%s--%s", planId, encrypted, checksum)
		pk := PlanKey{
			Uuid:   id,
			Key:    newKey,
			PlanId: plan.Id,
		}
		keyList = append(keyList, &pk)
	}
	return AddBatchPlanKeys(keyList)
}
