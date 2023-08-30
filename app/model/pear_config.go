package model

type PearConfig struct {
	ID           uint           `gorm:"primary_key"`
	ConfigType   PearConfigType `json:"config_type"`
	ConfigData   string         `json:"config_data" gorm:"type:varchar(3000);"`
	ConfigStatus int            `json:"config_status"`
	CreatedAt    string         `json:"created_at"`
	UpdatedAt    string         `json:"updated_at"`
}

type PearConfigType string

const (
	PearSiteConfig PearConfigType = "pear-config"
)

type PearConfigForm struct {
	Colors []PearColors `json:"colors"`
	Header PearHeader   `json:"header"`
	Links  []PearLinks  `json:"links"`
	Logo   PearLogo     `json:"logo"`
	Menu   PearMenu     `json:"menu"`
	Other  PearOther    `json:"other"`
	Tab    PearTab      `json:"tab"`
	Theme  PearTheme    `json:"theme"`
}

type PearMenu struct {
	Accordion bool   `json:"accordion"`
	Control   bool   `json:"control"`
	Data      string `json:"data"`
	Method    string `json:"method"`
	Select    string `json:"select"`
}

type PearTheme struct {
	AllowCustom  bool   `json:"allowCustom"`
	DefaultColor string `json:"defaultColor"`
	DefaultMenu  string `json:"defaultMenu"`
}

type PearOther struct {
	AutoHead bool  `json:"autoHead"`
	KeepLoad int64 `json:"keepLoad"`
}

type PearColors struct {
	Color string `json:"color"`
	ID    string `json:"id"`
}

type PearLinks struct {
	Href  string `json:"href"`
	Icon  string `json:"icon"`
	Title string `json:"title"`
}

type PearLogo struct {
	Image string `json:"image"`
	Title string `json:"title"`
}

type PearTab struct {
	Index     PearTabIndex `json:"index"`
	KeepState bool         `json:"keepState"`
	MuiltTab  bool         `json:"muiltTab"`
	TabMax    int64        `json:"tabMax"`
}

type PearTabIndex struct {
	Href  string `json:"href"`
	ID    string `json:"id"`
	Title string `json:"title"`
}

type PearHeader struct {
	Message string `json:"message"`
}

type MenuResp struct {
	ID       int                `json:"id"`
	Title    string             `json:"title"`
	Type     int                `json:"type"`
	Icon     string             `json:"icon"`
	Href     string             `json:"href"`
	Children []MenuChildrenResp `json:"children"`
}

type MenuChildrenResp struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Icon     string `json:"icon"`
	Type     int    `json:"type"`
	OpenType string `json:"openType"`
	Href     string `json:"href"`
}
