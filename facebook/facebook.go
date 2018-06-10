package facebook

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	toml "github.com/pelletier/go-toml"
)

type fieldschema struct {
	Fields []string `json:"fields"`
	Conns  []string `json:"connections"`
}

type Archiver struct {
	accessToken string
	id          string

	api      *Client
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
		id:          c.Get("id").(string),
		api:         NewClient(c.Get("access-token").(string)),
		metadata:    fields,
	}
}

func (a *Archiver) buildMetadata(ctx context.Context, id string) (string, error) {
	// Get the metadata
	m, err := a.api.GetNode(id, Meta)
	if err != nil {
		return "", err
	}

	kind := m.Metadata.Type

	fmt.Printf("%#v\n", m)

	_, known := a.metadata[kind]
	if known {
		return kind, nil
	}

	fields := []string{}
	for _, field := range m.Metadata.Fields {
		fmt.Println("Checking field name: ", field.Name)
		_, err := a.api.GetNode(id, Fields(field.Name))
		if err == nil {
			fields = append(fields, field.Name)
		}
	}

	conns := []string{}
	for edge, _ := range m.Metadata.Connections {
		fmt.Println("Checking edge: ", edge)
		if _, err := a.api.GetEdge(id, edge, nil); err == nil {
			conns = append(conns, edge)
		}
	}

	schema := fieldschema{
		Fields: fields,
		Conns:  conns,
	}

	a.metadata[kind] = schema

	blob, _ := json.MarshalIndent(a.metadata, "", "  ")
	ioutil.WriteFile("lookup.json", blob, 0644)
	return kind, nil
}

func (a *Archiver) archiveNode(ctx context.Context, id string) error {
	// Fetch the node
	kind, err := a.buildMetadata(ctx, id)
	if err != nil {
		return err
	}

	fmt.Printf("node:%s kind:%s\n", id, kind)

	_, err = a.api.GetNode(id, Meta, Fields(a.metadata[kind].Fields...))
	if err != nil {
		return err
	}

	return nil

	fmt.Println(a.metadata[kind].Conns)

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

func (a *Archiver) archiveConnection(ctx context.Context, id, connection string) error {
	fmt.Println("Connection:", id, connection)

	paging := Paging{}

	for {
		data, err := a.api.GetEdge(id, connection, &paging)
		if err != nil {
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

		paging = data.Paging
		if paging.Next == "" {
			return nil
		}
	}
	return nil
}

func (a *Archiver) Sync(ctx context.Context) error {
	return a.archiveNode(ctx, a.id)
}
