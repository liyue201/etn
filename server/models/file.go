package models

import (
	"github.com/liyue201/go-logger"
	"time"
	"usermgr/common/status"
)

//登录记录
type File struct {
	Model
	AccountId uint      `json:"account_id" gorm:"not null;index"`
	LoginAt   time.Time `json:"login_at"`
	Ip        string    `json:"ip"`
	Location  string    `json:"location"` //登录的地理位置
}

func GetAccountLoginLogs(accountId uint, offset, limit int) ([]*LoginLog, int) {
	logs := []*LoginLog{}

	ret := db.Model(&LoginLog{}).Where("account_id=?", accountId).Order("login_at DESC").Limit(limit).Offset(offset).Find(&logs)
	if ret.Error != nil {
		logger.Error("[GetAccountLoginLog] error,", ret.Error)
		return logs, status.InternalServerError
	}
	return logs, status.OK
}
