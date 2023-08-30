package service

import (
	"pear-admin-go/app/core/cache"
	e2 "pear-admin-go/app/global/e"
	"pear-admin-go/app/util/gconv"
	"time"
)

func Lock(loginName string) {
	cache.Instance().Set(e2.UserLock+loginName, true, time.Minute*5)
}

func UnLock(loginName string) {
	cache.Instance().Delete(e2.UserLock + loginName)
}

func CheckLock(loginName string) bool {
	res, ok := cache.Instance().Get(e2.UserLock + loginName)
	if ok && res == true {
		return true
	}
	return false
}

func SetPwdErrNum(loginName string) int {
	countNum := 0
	errNum, _ := cache.Instance().Get(e2.UserLoginErr + loginName)
	if errNum != nil {
		countNum = gconv.Int(errNum)
	}
	countNum = countNum + 1
	cache.Instance().Set(e2.UserLoginErr+loginName, countNum, time.Minute*1)
	if countNum >= 5 {
		Lock(loginName)
	}
	return countNum
}

func RemovePwdErrNum(loginName string) {
	cache.Instance().Delete(e2.UserLoginErr + loginName)
}
