package client

import (
	"bytes"
	"devture-matrix-corporal/corporal/httpapi/handler"
	"devture-matrix-corporal/corporal/httphelp"
	"devture-matrix-corporal/corporal/policy"
	"encoding/json"
	"fmt"
	"matrix-corporal-multitenant-proxy/config"
	"net/http"

	"github.com/sirupsen/logrus"
)

const (
	POLICY_ENDPOINT = "/_matrix/corporal/policy"
)

type Client struct {
	config     config.Instance
	httpClient *http.Client
}

func NewClient(config config.Instance) *Client {
	return &Client{
		config:     config,
		httpClient: &http.Client{},
	}
}

func (me *Client) UpdatePolicy(policy *policy.Policy) error {
	body, err := json.Marshal(policy)
	if err != nil {
		logrus.Warnf(`Failed to convert outgoing policy to json.`)
		return err
	}

	req, err := http.NewRequest("PUT", me.config.CorporalApiUrl+POLICY_ENDPOINT, bytes.NewReader(body))
	if err != nil {
		logrus.Warnf(`Policy PUT: failed creating request: %s`, err.Error())
		return err
	}

	req.Header.Add("Authorization", "Bearer "+me.config.CorporalAccessToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := me.httpClient.Do(req)
	if err != nil {
		logrus.Warnf(`Policy PUT: request failed: %s`, err.Error())
		return err
	}
	defer res.Body.Close()

	var apiError handler.ApiResponseError
	if err := httphelp.GetJsonFromResponseBody(res, &apiError); err != nil {
		logrus.Infof(`Policy PUT: failed read json from response. statusCode %d error: %s`, res.StatusCode, err.Error())
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Putting policy failed with stats code %d", res.StatusCode)
	}

	if apiError.ErrorCode != "" {
		return fmt.Errorf(`Corporal returned error: %s, %s`, apiError.ErrorCode, apiError.ErrorMessage)
	}

	return nil
}
