package foxpass

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"encoding/json"
	"html"
	"net/url"
	"strings"
	"errors"
)

const (
	BaseURLV1 = "https://api.foxpass.com/v1"
)

type FoxpassClient struct {
	BaseURL    string
	apiToken   string
	HTTPClient *http.Client
}

func NewClient(apiToken string) *FoxpassClient {
	return &FoxpassClient{
		BaseURL: BaseURLV1,
		apiToken:  apiToken,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (c *FoxpassClient) foxpassRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.apiToken))
	req.Header.Set("User-Agent", "foxpass-client-go")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

type GetMacEntryPrefix struct {
	Status string `json:"status"`
	Exists bool `json:"data"`
}

func (c *FoxpassClient) GetMacEntryPrefix(name string, prefix string) (bool, error) {
	resource := fmt.Sprintf("/mac_entries/%s/prefixes/%s/", name, html.EscapeString(prefix))

	u, _ := url.Parse(c.BaseURL)
	u.Path += resource
	urlStr := u.String()

	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return false, err
	}

	body, err := c.foxpassRequest(req)
	if err != nil {
		return false, err
	}

	var Prefix GetMacEntryPrefix
	err = json.Unmarshal(body, &Prefix)
	if err != nil {
		return false, err
	}
	if Prefix.Status != "ok" {
		return false, errors.New("FoxPass Status Not Okay")
	}

	return Prefix.Exists, nil
}

type AddMacEntryPrefix struct {
	Status string `json:"status"`
}

func (c *FoxpassClient) AddMacEntryPrefix(name string, prefix string) error {
	resource := fmt.Sprintf("/mac_entries/%s/prefixes/", name)

	payload := strings.NewReader("{\"prefix\":\"" + prefix + "\"}")

	u, _ := url.Parse(c.BaseURL)
	u.Path += resource
	urlStr := u.String() 

	req, err := http.NewRequest(http.MethodPut, urlStr, payload) 
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	body, err := c.foxpassRequest(req)
	if err != nil {
		return err
	}

	var Prefix AddMacEntryPrefix
	err = json.Unmarshal(body, &Prefix)
	if err != nil {
		return err
	}
	if Prefix.Status != "ok" {
		return errors.New("FoxPass Status Not Okay")
	}

	return nil
}

type DeleteMacEntryPrefix struct {
	Status string `json:"status"`
}

func (c *FoxpassClient) DeleteMacEntryPrefix(name string, prefix string) error {
	resource := fmt.Sprintf("/mac_entries/%s/prefixes/%s/", name, html.EscapeString(prefix))

	u, _ := url.Parse(c.BaseURL)
	u.Path += resource
	urlStr := u.String() 

	req, err := http.NewRequest(http.MethodDelete, urlStr, nil)
	if err != nil {
		return err
	}

	body, err := c.foxpassRequest(req)
	if err != nil {
		return err
	}

	var Prefix DeleteMacEntryPrefix
	err = json.Unmarshal(body, &Prefix)
	if err != nil {
		return err
	}
	if Prefix.Status != "ok" {
		return errors.New("FoxPass Status Not Okay")
	}

	return nil
}
