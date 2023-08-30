package model

import "time"

type Task struct {
	TaskCommon
	CreateTime time.Time
}

type TaskCommon struct {
	Id           int    `json:"id" zh:"ID" form:"id"`
	SourceType   int    `json:"source_type" zh:"源类型" binding:"required" form:"source_type"`
	SourceServer int    `json:"source_server" zh:"源服务器" form:"source_server"`
	SourcePath   string `json:"source_path" zh:"源地址" binding:"required" form:"source_path"`
	DstType      int    `json:"dst_type" zh:"目标类型" binding:"required" form:"dst_type"`
	DstServer    int    `json:"dst_server" zh:"目标服务器" form:"dst_server"`
	DstPath      string `json:"dst_path" zh:"目标地址" binding:"required" form:"dst_path"`
	ExcludeType  string `json:"exclude_type" zh:"排除类型" form:"exclude_type"`
	TaskFileNum  int    `json:"task_file_num"`
}

type TaskResp struct {
	TaskCommon
	TaskLogNum int    `json:"task_log_num"`
	CreateTime string `json:"create_time"`
}
