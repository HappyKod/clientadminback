package clientadminback

type Proxy struct {
	Start    string `json:"start"`
	End      string `json:"end"`
	Created  string `json:"created"`
	Id       int    `json:"id"`
	Host     string `json:"host"`
	Http     int    `json:"http"`
	Socks    int    `json:"socks"`
	Username string `json:"username"`
	Passw    string `json:"passw"`
	Type     string `json:"type"`
	Country  string `json:"country"`
	Active   bool   `json:"active"`
	Status   struct {
		Socks5 struct {
			Status  bool `json:"status"`
			Checked struct {
				Time  string `json:"time"`
				Valid bool   `json:"valid"`
			} `json:"checked"`
		} `json:"socks5"`
		Https struct {
			Status  bool `json:"status"`
			Checked struct {
				Time  string `json:"time"`
				Valid bool   `json:"valid"`
			} `json:"checked"`
		} `json:"https"`
	} `json:"status"`
	DaysToEnd int         `json:"days_to_end"`
	CreatedBy string      `json:"created_by"`
	Edited    interface{} `json:"edited"`
	EditedBy  interface{} `json:"edited_by"`
}
