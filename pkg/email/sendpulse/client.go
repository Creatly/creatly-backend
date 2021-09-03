package sendpulse

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/zhashkevych/creatly-backend/pkg/cache"
	"github.com/zhashkevych/creatly-backend/pkg/email"
	"github.com/zhashkevych/creatly-backend/pkg/logger"
)

// Documentation https://sendpulse.com/integrations/api
// Note: The request limit is 10 requests per second.

const (
	endpoint          = "https://api.sendpulse.com"
	authorizeEndpoint = "/oauth/access_token"
	addToListEndpoint = "/addressbooks/%s/emails" // addressbooks/{id}/emails

	grantType = "client_credentials"

	cacheTTL = 3600 // In seconds. SendPulse access tokens are valid for 1 hour
)

type authRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type authResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type addToListRequest struct {
	Emails []emailInfo `json:"emails"`
}

type emailInfo struct {
	Email     string            `json:"email"`
	Variables map[string]string `json:"variables"`
}

// Client is SendPulse API client implementation.
type Client struct {
	id     string
	secret string

	cache cache.Cache
}

func NewClient(id, secret string, cache cache.Cache) *Client {
	return &Client{id: id, secret: secret, cache: cache}
}

// AddEmailToList adds lead to provided email list with specific variables.
func (c *Client) AddEmailToList(input email.AddEmailInput) error {
	token, err := c.getToken()
	if err != nil {
		return err
	}

	reqData := addToListRequest{
		Emails: []emailInfo{
			{
				Email:     input.Email,
				Variables: input.Variables,
			},
		},
	}

	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return err
	}

	path := fmt.Sprintf(addToListEndpoint, input.ListID)

	req, err := http.NewRequest(http.MethodPost, endpoint+path, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	logger.Infof("SendPulse response: %s", string(body))

	if resp.StatusCode != 200 {
		return errors.New("status code is not OK")
	}

	return nil
}

func (c *Client) getToken() (string, error) {
	// todo set unique key (by schoolId)
	token, err := c.cache.Get("t")
	if err == nil {
		return token.(string), nil
	}

	token, err = c.authenticate()
	if err != nil {
		return "", err
	}

	err = c.cache.Set("t", token, cacheTTL)

	return token.(string), err
}

func (c *Client) authenticate() (string, error) {
	reqData := authRequest{
		GrantType:    grantType,
		ClientID:     c.id,
		ClientSecret: c.secret,
	}

	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(endpoint+authorizeEndpoint, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("status code is not OK")
	}

	var respData authResponse

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(respBody, &respData)
	if err != nil {
		return "", err
	}

	return respData.AccessToken, nil
}
