package facebook

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	toml "github.com/pelletier/go-toml"
)

type fieldschema struct {
	Fields []string `json:"fields"`
	Conns  []string `json:"connections"`
}

type Archiver struct {
	accessToken string

	metadata map[string]fieldschema
}

func NewArchiver(c *toml.Tree) *Archiver {
	fields := map[string]fieldschema{}
	if blob, err := ioutil.ReadFile("lookup.json"); err == nil {
		if err := json.Unmarshal(blob, &fields); err != nil {
			panic(err)
		}
	}

	return &Archiver{
		accessToken: c.Get("access-token").(string),
		metadata:    fields,
	}
}

type field struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type metadata struct {
	Fields      []field           `json:"fields"`
	Connections map[string]string `json:"connections"`
	Type        string            `json:"type"`
}

type node struct {
	Metadata metadata `json:"metadata"`
	ID       string   `json:"id"`
}

func (a *Archiver) buildMetadata(ctx context.Context, id string) (string, error) {
	// Get the metadata
	form := url.Values{}
	form.Set("access_token", a.accessToken)
	form.Set("method", "GET")
	form.Set("metadata", "1")
	req, _ := http.NewRequest("POST", "https://graph.facebook.com/v2.12/"+id, strings.NewReader(form.Encode()))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	var m node
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&m); err != nil {
		return "", err
	}

	kind := m.Metadata.Type

	_, known := a.metadata[kind]
	if known {
		return kind, nil
	}

	form.Del("metadata") // Clean up from previous request

	fields := []string{}
	for _, field := range m.Metadata.Fields {
		fmt.Println("Checking field name: ", field.Name)
		form.Set("fields", field.Name)
		req, _ := http.NewRequest("POST", "https://graph.facebook.com/v2.12/"+id, strings.NewReader(form.Encode()))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}
		if resp.StatusCode == http.StatusOK {
			fields = append(fields, field.Name)
		}
	}

	form.Del("fields") // Clean up from previous request

	conns := []string{}
	for edge, _ := range m.Metadata.Connections {
		fmt.Println("Checking edge: ", edge)
		req, _ := http.NewRequest("POST", "https://graph.facebook.com/v2.12/"+id+"/"+edge, strings.NewReader(form.Encode()))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}
		if resp.StatusCode == http.StatusOK {
			conns = append(conns, edge)
		}
	}

	schema := fieldschema{
		Fields: fields,
		Conns:  conns,
	}

	a.metadata[kind] = schema

	blob, _ := json.Marshal(a.metadata)
	ioutil.WriteFile("lookup.json", blob, 0644)
	return kind, nil
}

func (a *Archiver) archiveNode(ctx context.Context, id string) error {
	// Fetch the node
	kind, err := a.buildMetadata(ctx, id)
	if err != nil {
		return err
	}

	fmt.Println("Node:", id, kind)

	form := url.Values{}
	form.Set("access_token", a.accessToken)
	form.Set("method", "GET")
	form.Set("fields", strings.Join(a.metadata[kind].Fields, ","))

	req, _ := http.NewRequest("POST", "https://graph.facebook.com/v2.12/"+id, strings.NewReader(form.Encode()))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		blob, _ := httputil.DumpResponse(resp, true)
		fmt.Println(string(blob))
		return fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	blob, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// if err := q.PutNode(ctx, a.db, id, kind, blob); err != nil {
	// 	return err
	// }

	// Archive connections
	for _, connection := range a.metadata[kind].Conns {
		if connection == "albums" || connection == "photos" {
			if err := a.archiveConnection(ctx, id, connection); err != nil {
				return err
			}
		}
	}

	return nil
}

type paging struct {
	Previous string `json:"previous"`
	Next     string `json:"next"`
}

type datalist struct {
	Data   []node
	Paging paging
}

func (a *Archiver) archiveConnection(ctx context.Context, id, connection string) error {
	base := url.Values{}
	base.Set("access_token", a.accessToken)
	base.Set("method", "GET")
	base.Set("limit", "2000")
	u := "https://graph.facebook.com/v2.12/" + id + "/" + connection

	fmt.Println("Connection:", id, connection)

	for {
		if u == "" {
			return nil
		}

		req, _ := http.NewRequest("POST", u, strings.NewReader(base.Encode()))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			blob, _ := httputil.DumpResponse(resp, true)
			fmt.Println(string(blob))
			return fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
		}
		defer resp.Body.Close()

		var data datalist
		dec := json.NewDecoder(resp.Body)

		if err := dec.Decode(&data); err != nil {
			return err
		}

		for _, node := range data.Data {
			// if err := q.PutEdge(ctx, a.db, connection, id, node.ID); err != nil {
			// 	return err
			// }
			if err := a.archiveNode(ctx, node.ID); err != nil {
				return err
			}
		}

		u = data.Paging.Next
	}
	return nil
}

func (a *Archiver) Sync(ctx context.Context) error {
	// hardcoded ID for now
	id := "1058550007"
	return a.archiveNode(ctx, id)
}
