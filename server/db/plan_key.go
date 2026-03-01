package db

import (
	"clipshare/types"
	"clipshare/utils"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

// PlanKey 套餐的密钥
type PlanKey struct {
	Uuid   string     `gorm:"index"`        //密钥的原始id
	Key    string     `gorm:"primaryKey"`   //密钥
	PlanId string     `gorm:"index"`        //套餐id
	UseAt  *time.Time `gorm:"null"`         //开始使用时间
	Enable bool       `gorm:"default:true"` //是否启用
	Remark string     `gorm:"null"`         //备注
	gorm.Model
}

func GetPlanKeysPageData(params types.PlanKeysSearchDto) (int64, []PlanKey, error) {
	checkDb()
	planId := params.PlanId
	pageNum := params.PageNum
	pageSize := params.PageSize
	var total int64
	query := AppDb.Model(PlanKey{})
	if planId != "" {
		query = query.Where("plan_id = ?", planId)
	}
	if err := query.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	var keys []PlanKey
	offset := (pageNum - 1) * pageSize
	if params.Sorts != nil && len(params.Sorts) > 0 {
		length := len(params.Sorts)
		orders := ""
		for i, sort := range params.Sorts {
			orders += sort.Key + " " + sort.Order
			if i != length-1 {
				orders += ","
			}
		}
		query = query.Order(orders)
	}
	result := query.Offset(offset).Limit(pageSize).Find(&keys)
	return total, keys, result.Error
}

func AddBatchPlanKeys(keys []*PlanKey) error {
	checkDb()
	result := AppDb.Create(&keys)
	return result.Error
}
func UpdatePlanKeyStatus(id int, enable bool) error {
	checkDb()
	result := AppDb.Model(&PlanKey{}).Where("id = ?", id).Update("enable", enable)
	return result.Error
}
func GetPlanKey(key string) *PlanKey {
	checkDb()
	pk := &PlanKey{}
	result := AppDb.Model(&PlanKey{}).First(pk, "key = ?", key)
	if result.Error != nil {
		return nil
	}
	return pk
}
func GetFistUseTime(key string) *time.Time {
	checkDb()
	pk := &PlanKey{}
	result := AppDb.Model(&PlanKey{}).First(pk, "key = ?", key)
	if result.Error != nil {
		return nil
	}
	return pk.UseAt
}
func GetPlanKeyById(id int) *PlanKey {
	checkDb()
	pk := &PlanKey{}
	result := AppDb.Model(&PlanKey{}).First(pk, "id = ?", id)
	if result.Error != nil {
		return nil
	}
	return pk
}
func getPlanKeyByUuid(keyUuid string) *PlanKey {
	var key PlanKey
	res := AppDb.Model(&PlanKey{}).First(&key, "uuid = ?", keyUuid)
	if res.Error != nil {
		return nil
	}
	return &key
}
func VerifyKey(key string) *PlanType {
	parts := strings.Split(key, "--")
	if len(parts) != 3 {
		return nil
	}
	//校验是否存在该Plan
	id := parts[0]
	plan := getPlanInfo(id)
	if plan == nil {
		return nil
	}
	//校验检验和
	encryptedKey := parts[1]
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedKey)
	if err != nil {
		return nil
	}
	checksum := fmt.Sprintf("%x", sha256.Sum256(encryptedBytes))[:8]
	if checksum != parts[2] {
		return nil
	}
	//校验加密数据
	decryptedBytes := utils.AesDecryptCBC(encryptedBytes, []byte(plan.AesKey))
	keyUuid := string(decryptedBytes)
	pk := getPlanKeyByUuid(keyUuid)
	if pk == nil || pk.Enable == false {
		return nil
	}
	return plan
}
func GetKeyUUID(key string) *string {
	plan := VerifyKey(key)
	if plan == nil {
		return nil
	}
	parts := strings.Split(key, "--")
	encryptedKey := parts[1]
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedKey)
	if err != nil {
		return nil
	}
	//校验加密数据
	decryptedBytes := utils.AesDecryptCBC(encryptedBytes, []byte(plan.AesKey))
	keyUuid := string(decryptedBytes)
	return &keyUuid
}
func GetKeyRemainingSeconds(key string) *uint64 {
	exhausted := uint64(0)
	keyUuid := GetKeyUUID(key)
	if keyUuid == nil {
		return nil
	}
	var pk PlanKey
	result := AppDb.First(&pk, "uuid=?", keyUuid)
	if result.Error != nil {
		return nil
	}
	plan := getPlanInfo(pk.PlanId)
	if plan == nil {
		return nil
	}
	if plan.Lifespan == nil {
		return nil
	}
	remainSeconds := *plan.Lifespan
	if pk.UseAt == nil {
		return nil
	}
	now := time.Now()
	//时间回溯，数据有问题，返回0
	if now.Before(*pk.UseAt) {
		return &exhausted
	}
	duration := now.Sub(*pk.UseAt)
	//计算从第一次使用到现在的秒数
	offsetSeconds := uint64(duration.Seconds())
	//超过plan的使用时间，返回0
	if offsetSeconds >= remainSeconds {
		return &exhausted
	}
	//返回剩余使用时间
	remainSeconds -= offsetSeconds
	return &remainSeconds
}
func InitializeKeyFirstUseTime(key string) {
	AppDb.Model(&PlanKey{}).
		Where("key = ? and use_at is null", key).
		Update("use_at", time.Now())
}
func ToPlanKeyDto(pk PlanKey, planName string) types.PlanKeyDto {
	timeLayout := "2006-01-02 15:04:05.000"
	var useTimeStr *string
	if pk.UseAt != nil {
		str := pk.UseAt.Format(timeLayout)
		useTimeStr = &str
	}
	return types.PlanKeyDto{
		Id:        pk.ID,
		Key:       pk.Key,
		PlanId:    pk.PlanId,
		PlanName:  planName,
		UseAt:     useTimeStr,
		Enable:    pk.Enable,
		Remark:    pk.Remark,
		CreatedAt: pk.CreatedAt.Format(timeLayout),
	}
}
