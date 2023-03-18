package client_admin_back

import "github.com/HappyKod/clientadminback/clientadminmodels"

type Clienter interface {
	GetAccounts(srcs string, active bool, groupID, limit int) ([]models.Account, error)
	DeleteAccounts(accountID string) error
}
