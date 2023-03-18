package clientadminback

type ClientAdmin interface {
	GetAccounts(srcs string, active bool, groupID, limit int) ([]Account, error)
	DeleteAccounts(accountID string) error
}
