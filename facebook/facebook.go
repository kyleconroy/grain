package facebook

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"

	"github.com/kyleconroy/grain/archive"
	"github.com/kyleconroy/grain/gen/facebook"
	toml "github.com/pelletier/go-toml"
)

type fieldschema struct {
	Fields []string `json:"fields"`
	Conns  []string `json:"connections"`
}

type Archiver struct {
	accessToken string
	id          string

	count    int
	api      *Client
	archive  *facebookpb.Archive
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
		archive:     &facebookpb.Archive{},
		metadata:    fields,
	}
}

func (a *Archiver) buildMetadata(ctx context.Context, id string) (string, error) {
	// Get the metadata
	m, err := a.api.GetNode(id)
	if err != nil {
		return "", err
	}

	kind := m.Metadata.Type

	_, known := a.metadata[kind]
	if known {
		return kind, nil
	}

	fields := []string{}
	for _, field := range m.Metadata.Fields {
		fmt.Println("Checking field: ", field.Name)
		_, err := a.api.GetNode(id, Fields(field.Name))
		if err == nil {
			fields = append(fields, field.Name)
		}
	}

	conns := []string{}
	for edge, _ := range m.Metadata.Connections {
		fmt.Println("Checking connection: ", edge)
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

func (a *Archiver) save() error {
	{
		path := filepath.Join("archive", "facebook", "me.json")
		b, err := json.MarshalIndent(facebookpb.Archive{Me: a.archive.Me}, "", "  ")
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(path, b, 0644); err != nil {
			return err
		}
	}
	{
		path := filepath.Join("archive", "facebook", "photos.json")
		b, err := json.MarshalIndent(facebookpb.Archive{Photos: a.archive.Photos}, "", "  ")
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(path, b, 0644); err != nil {
			return err
		}
	}
	return nil
}

func (a *Archiver) archiveNode(ctx context.Context, id string) error {
	// Fetch the node
	kind, err := a.buildMetadata(ctx, id)
	if err != nil {
		return err
	}

	fmt.Printf("node:%s kind:%s\n", id, kind)

	node, err := a.api.GetNode(id, Fields(a.metadata[kind].Fields...))
	if err != nil {
		return err
	}

	switch v := node.Message.(type) {
	case *facebookpb.Photo:
		a.archive.Photos = append(a.archive.Photos, v)
		sort.Slice(v.Images, func(i, j int) bool {
			return (v.Images[i].Width + v.Images[i].Height) > (v.Images[j].Width + v.Images[j].Height)
		})
		if err := archive.ArchiveURL(ctx, "facebook", "media", v.Images[0].Source); err != nil {
			return err
		}
	case *facebookpb.User:
		if v.Id == a.id {
			a.archive.Me = v
		} else {
		}
	}

	a.count += 1
	if a.count > 100 {
		if err := a.save(); err != nil {
			return err
		}
	}

	// fmt.Println(a.metadata[kind].Conns)

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
