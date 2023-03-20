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
	"strings"

	"github.com/HappyKod/clientadminback"
)

// ErrorAuthFailed error occurs when it was not possible to get account data.
var ErrorAuthFailed = errors.New("auth failed")

// Option represents the client options.
type Option func(*ClientAdmin) error

// WithJWTToken authorization occurs using the tokenJWT token received from the user.
func WithJWTToken(token string) Option {
	return func(c *ClientAdmin) error {
		path, err := url.JoinPath(c.ServiceURL, "v1/users/me")
		if err != nil {
			return err
		}
		request, err := http.NewRequest(http.MethodGet, path, nil)
		if err != nil {
			return err
		}
		//remove Bearer if the user entered this format
		token = strings.ReplaceAll(token, "Bearer ", "")
		token = "Bearer " + token
		request.Header.Add("Authorization", token)
		request.Header.Add("content-type", "application/json")
		client := &http.Client{}
		do, err := client.Do(request)
		if err != nil {
			return fmt.Errorf("%w: %s", err, token)
		}
		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				log.Printf("error closing response body: %v", err)
			}
		}(do.Body)
		all, err := io.ReadAll(do.Body)
		if err != nil {
			return err
		}
		var clientInfo clientadminback.ClientInfo
		err = json.Unmarshal(all, &clientInfo)
		if err != nil {
			return err
		}
		if !clientInfo.IsActive {
			return fmt.Errorf("%w: %s", ErrorAuthFailed, token)
		}
		c.tokenJWT = token
		return nil
	}
}

// WithLoginAndPassword authorization occurs using username and password received from the user.
func WithLoginAndPassword(loginAndPassword clientadminback.LoginAndPassword) Option {
	return func(c *ClientAdmin) error {
		marshal, err := json.Marshal(loginAndPassword)
		if err != nil {
			return fmt.Errorf("%w: %s", err, loginAndPassword.Email)
		}
		path, err := url.JoinPath(c.ServiceURL, "v1/auth/login")
		if err != nil {
			return err
		}
		post, err := http.Post(path, "application/json", bytes.NewReader(marshal))
		if err != nil {
			return fmt.Errorf("%w: %s", err, loginAndPassword.Email)
		}
		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				log.Printf("%s: %v", ErrorCloseBody, err)
			}
		}(post.Body)
		all, err := io.ReadAll(post.Body)
		if err != nil {
			return err
		}
		var clientInfo clientadminback.ClientInfo
		err = json.Unmarshal(all, &clientInfo)
		if err != nil {
			return err
		}
		if !clientInfo.IsActive {
			return fmt.Errorf("%w: %s", ErrorAuthFailed, loginAndPassword.Email)
		}
		c.tokenJWT = clientInfo.Token
		return nil
	}
}
