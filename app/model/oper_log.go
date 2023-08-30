package model

type OperLog struct {
	ID            uint64
	Title         string `json:"title"`
	BusinessType  int    `json:"business_type"`
	Method        string `json:"method"`
	RequestMethod string `json:"request_method"`
	OperatorType  int    `json:"operator_type"`
	OperName      string `json:"oper_name"`
	DeptName      string `json:"dept_name"`
	OperUrl       string `json:"oper_url"`
	OperIp        string `json:"oper_ip"`
	OperLocation  string `json:"oper_location"`
	OperParam     string `json:"oper_param"`
	JsonResult    string `json:"json_result"`
	Status        int    `json:"status"`
	ErrorMsg      string `json:"error_msg"`
	OperTime      string `json:"oper_time"`
}

type OperForm struct {
	Title      string
	InContent  string
	ErrorMsg   string
	OutContent *CommonResp
}
