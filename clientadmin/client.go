package clientadmin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	client_admin_back "github.com/HappyKod/clientadminback"
	"github.com/HappyKod/clientadminback/models"
)

type ClientAdmin struct {
	ServiceURL string
	tokenJWT   string
}

var _ client_admin_back.Clienter = (*ClientAdmin)(nil)
var ErrorCloseBody = errors.New("error closing response body: ")
var ErrorUnexpectedAnswer = errors.New("unexpected answer: ")

func NewClientAdmin(serviceURL string, opts ...Option) (*ClientAdmin, error) {
	clientAdmin := ClientAdmin{
		ServiceURL: serviceURL,
	}
	for _, opt := range opts {
		err := opt(&clientAdmin)
		if err != nil {
			return nil, err
		}
	}
	return &clientAdmin, nil
}

func (c ClientAdmin) GetAccounts(srcs string, active bool, groupID, limit int) ([]models.Account, error) {
	path, err := url.JoinPath(c.ServiceURL, "v1/accounts")
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.tokenJWT)
	q := req.URL.Query()
	q.Add("type", srcs)
	q.Add("group_id", fmt.Sprint(groupID))
	q.Add("limit", fmt.Sprint(limit))
	q.Add("is_active", fmt.Sprint(active))
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	do, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("%s: %v", ErrorCloseBody, err)
		}
	}(do.Body)
	all, err := io.ReadAll(do.Body)
	if err != nil {
		return nil, err
	}
	var accounts []models.Account
	err = json.Unmarshal(all, &accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (c ClientAdmin) DeleteAccounts(accountID string) error {
	path, err := url.JoinPath(c.ServiceURL, "v1/accounts", accountID)
	req, err := http.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.tokenJWT)
	client := &http.Client{}
	do, err := client.Do(req)
	if err != nil {
		return err
	}
	if do.StatusCode != http.StatusOK {
		all, errReadAll := io.ReadAll(do.Body)
		if errReadAll != nil {
			return errReadAll
		}
		return fmt.Errorf("%w: message %s", ErrorUnexpectedAnswer, string(all))
	}
	return nil
}
