package clientadminback

type Account struct {
	Launched    string `json:"launched"`
	Created     string `json:"created"`
	Id          int    `json:"id"`
	Phone       string `json:"phone"`
	IsActive    bool   `json:"is_active"`
	GroupId     int    `json:"group_id"`
	Type        string `json:"type"`
	ServiceInfo struct {
		Token    string `json:"token"`
		Password string `json:"password"`
		Username string `json:"username"`
	} `json:"service_info"`
	AdditionalInfo interface{} `json:"additional_info"`
	Backup         string      `json:"backup"`
	Requests       int         `json:"requests"`
	RequestsLimit  int         `json:"requests_limit"`
	CreatedBy      string      `json:"created_by"`
	Edited         interface{} `json:"edited"`
	EditedBy       interface{} `json:"edited_by"`
}
