package model

import (
	"time"
)

type TaskLog struct {
	Id         int       `json:"id" zh:"ID" form:"id"`
	TaskId     int       `json:"task_id"`
	ServerId   int       `json:"server_id"`
	SourcePath string    `json:"source_path"`
	DstPath    string    `json:"dst_path"`
	Size       int64     `json:"size"`
	CreateTime time.Time `json:"create_time"`
}
