package model

type AdminOnline struct {
	SessionID      string `json:"session_id"`
	LoginName      string `json:"login_name"`
	DeptName       string `json:"dept_name"`
	IpAddr         string `json:"ip_addr"`
	LoginLocation  string `json:"login_location"`
	Browser        string `json:"browser"`
	Os             string `json:"os"`
	Status         string `json:"status"`
	StartTimestamp int64  `json:"start_timestamp"`
	LastAccessTime string `json:"last_access_time"`
	ExpireTime     int    `json:"expire_time"`
}
