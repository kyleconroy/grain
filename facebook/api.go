package facebook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/kyleconroy/grain/gen/facebook"
)

type Client struct {
	httpClient  *http.Client
	baseURL     string
	accessToken string
}

func NewClient(token string) *Client {
	return &Client{
		accessToken: token,
		httpClient:  http.DefaultClient,
		baseURL:     "https://graph.facebook.com/v3.0/",
	}
}

type Field struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type Metadata struct {
	Fields      []Field           `json:"fields"`
	Connections map[string]string `json:"connections"`
	Type        string            `json:"type"`
}

type Node struct {
	Metadata Metadata `json:"metadata"`
	ID       string   `json:"id"`

	Message proto.Message
}

type Option func(*url.Values)

func Fields(args ...string) func(form *url.Values) {
	return func(form *url.Values) {
		form.Set("fields", strings.Join(args, ","))
	}
}

func (c *Client) GetNode(id string, options ...Option) (*Node, error) {
	form := url.Values{}
	form.Set("access_token", c.accessToken)
	form.Set("metadata", "1")
	form.Set("method", "GET")
	for _, opt := range options {
		opt(&form)
	}

	req, _ := http.NewRequest("POST", c.baseURL+id, strings.NewReader(form.Encode()))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		blob, _ := httputil.DumpResponse(resp, true)
		fmt.Println(string(blob))
		return nil, fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	var n Node
	blob, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(blob, &n); err != nil {
		return nil, err
	}

	// fmt.Println(string(blob))
	var m proto.Message
	switch n.Metadata.Type {
	case "user":
		m = &facebookpb.User{}
	case "photo":
		m = &facebookpb.Photo{}
	}

	if err := json.Unmarshal(blob, &m); err != nil {
		return nil, err
	}

	n.Message = m
	return &n, nil
}

type Paging struct {
	Previous string `json:"previous"`
	Next     string `json:"next"`
}

type Datalist struct {
	Data   []Node
	Paging Paging
}

func (c *Client) GetEdge(id, connection string, paging *Paging, options ...Option) (*Datalist, error) {
	form := url.Values{}
	form.Set("access_token", c.accessToken)
	form.Set("method", "GET")
	form.Set("limit", "2000")
	for _, opt := range options {
		opt(&form)
	}

	u := c.baseURL + id + "/" + connection
	if paging != nil && paging.Next != "" {
		u = paging.Next
	}

	req, _ := http.NewRequest("POST", u, strings.NewReader(form.Encode()))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		blob, _ := httputil.DumpResponse(resp, true)
		fmt.Println(string(blob))
		return nil, fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	var dl Datalist
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&dl); err != nil {
		return nil, err
	}
	return &dl, nil
}
