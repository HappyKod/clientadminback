package clientadmin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/HappyKod/clientadminback"
)

type ClientAdmin struct {
	ServiceURL string
	tokenJWT   string
}

var _ clientadminback.ClientAdmin = (*ClientAdmin)(nil)
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

func (c *ClientAdmin) GetAccounts(srcs string, active bool, groupID, limit int) ([]clientadminback.Account, error) {
	path, err := url.JoinPath(c.ServiceURL, "v1/accounts")
	if err != nil {
		return nil, err
	}
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
	var accounts []clientadminback.Account
	err = json.Unmarshal(all, &accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (c *ClientAdmin) DeleteAccounts(accountID int) error {
	path, err := url.JoinPath(c.ServiceURL, "v1/accounts", fmt.Sprint(accountID))
	if err != nil {
		return err
	}
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
	defer func() {
		err = do.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	if do.StatusCode != http.StatusOK {
		all, errReadAll := io.ReadAll(do.Body)
		if errReadAll != nil {
			return errReadAll
		}
		return fmt.Errorf("%w: message %s", ErrorUnexpectedAnswer, string(all))
	}
	return nil
}

func (c *ClientAdmin) GetProxies() ([]clientadminback.Proxy, error) {
	path, err := url.JoinPath(c.ServiceURL, "v1/proxies")
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.tokenJWT)
	client := &http.Client{}
	do, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = do.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	all, err := io.ReadAll(do.Body)
	if err != nil {
		return nil, err
	}
	var proxies []clientadminback.Proxy
	err = json.Unmarshal(all, &proxies)
	if err != nil {
		return nil, err
	}
	return proxies, nil
}

func (c *ClientAdmin) PatchAccount(ID int, account clientadminback.Account) error {
	path, err := url.JoinPath(c.ServiceURL, "v1/accounts/", fmt.Sprint(ID))
	if err != nil {
		return err
	}
	marshal, err := json.Marshal(account)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPatch, path, bytes.NewReader(marshal))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.tokenJWT)
	client := &http.Client{}
	do, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		err = do.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	if do.StatusCode != http.StatusCreated {
		return errors.New(fmt.Sprintf("error %d", do.StatusCode))
	}
	return nil
}
