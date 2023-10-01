package clientadminback

type ClientAdmin interface {
	GetAccounts(srcs string, active bool, groupID, limit int) ([]Account, error)
	PatchAccount(ID int, account Account) error
	DeleteAccounts(accountID int) error
	GetProxies() ([]Proxy, error)
}
