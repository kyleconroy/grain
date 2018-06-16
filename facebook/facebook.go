package facebook

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/golang/protobuf/jsonpb"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/kyleconroy/grain/archive"
	"github.com/kyleconroy/grain/gen/facebook"
	toml "github.com/pelletier/go-toml"
)

const (
	FilenameEvents  = "events.json"
	FilenameFriends = "friends.json"
	FilenameMe      = "me.json"
	FilenamePhotos  = "photos.json"
)

type fieldschema struct {
	Fields []string `json:"fields"`
	Conns  []string `json:"connections"`
}

type kindCache struct {
	albums map[string]struct{}
	me     string
	photos map[string]struct{}
	users  map[string]struct{}
}

type Archiver struct {
	accessToken string

	count     int
	seen      map[string]struct{}
	api       *Client
	archive   *facebookpb.Archive
	kindCache kindCache
	metadata  map[string]fieldschema
}

func NewArchiver(c *toml.Tree) *Archiver {
	fields := map[string]fieldschema{}
	if blob, err := ioutil.ReadFile("lookup.json"); err == nil {
		if err := json.Unmarshal(blob, &fields); err != nil {
			panic(err)
		}
	}

	transport := httpcache.NewTransport(diskcache.New("httpcache"))
	fapi := NewClient(c.Get("access-token").(string))
	fapi.HttpClient = transport.Client()

	return &Archiver{
		accessToken: c.Get("access-token").(string),
		archive:     &facebookpb.Archive{},
		metadata:    fields,
		api:         fapi,
		seen:        map[string]struct{}{},
	}
}

func LoadArchive(base string) (*facebookpb.Archive, error) {
	m := jsonpb.Unmarshaler{}
	a := facebookpb.Archive{}
	paths := []string{
		FilenameEvents,
		FilenameFriends,
		FilenameMe,
		FilenamePhotos,
	}

	for _, filename := range paths {
		path := filepath.Join(base, filename)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		r, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		if err := m.Unmarshal(r, &a); err != nil {
			return nil, err
		}
	}
	return &a, nil
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
	fmt.Println("save")
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
	{
		path := filepath.Join("archive", "facebook", "albums.json")
		b, err := json.MarshalIndent(facebookpb.Archive{Albums: a.archive.Albums}, "", "  ")
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(path, b, 0644); err != nil {
			return err
		}
	}
	{
		path := filepath.Join("archive", "facebook", "friends.json")
		b, err := json.MarshalIndent(facebookpb.Archive{Friends: a.archive.Friends}, "", "  ")
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(path, b, 0644); err != nil {
			return err
		}
	}
	return nil
}

func (a *Archiver) cachedKind(id string) string {
	if a.kindCache.me == id {
		return "me"
	}
	if _, ok := a.kindCache.photos[id]; ok {
		return "photo"
	}
	if _, ok := a.kindCache.albums[id]; ok {
		return "album"
	}
	if _, ok := a.kindCache.users[id]; ok {
		return "user"
	}
	return ""
}

func (a *Archiver) visitNode(ctx context.Context, id string) (string, error) {
	// If in the cache, return the type with no request
	if kind := a.cachedKind(id); kind != "" {
		return kind, nil
	}

	// Fetch the node
	kind, err := a.buildMetadata(ctx, id)
	if err != nil {
		return "", err
	}

	node, err := a.api.GetNode(id, Fields(a.metadata[kind].Fields...))
	if err != nil {
		return "", err
	}

	switch v := node.Message.(type) {
	case *facebookpb.Photo:
		a.archive.Photos = append(a.archive.Photos, v)
		sort.Slice(v.Images, func(i, j int) bool {
			return (v.Images[i].Width + v.Images[i].Height) > (v.Images[j].Width + v.Images[j].Height)
		})
		if err := archive.ArchiveURL(ctx, "facebook", "media", v.Images[0].Source); err != nil {
			return "", err
		}
	case *facebookpb.Album:
		a.archive.Albums = append(a.archive.Albums, v)
	case *facebookpb.User:
		if id == "me" {
			a.archive.Me = v
		} else {
			a.archive.Friends = append(a.archive.Friends, v)
		}
	}

	a.count += 1
	if a.count > 100 {
		if err := a.save(); err != nil {
			return "", err
		}
		a.count = 0
	}
	return kind, nil
}

func (a *Archiver) archiveNode(ctx context.Context, id string) error {
	// Don't process a node twice
	if _, ok := a.seen[id]; ok {
		return nil
	}

	kind, err := a.visitNode(ctx, id)
	if err != nil {
		return err
	}

	a.seen[id] = struct{}{}
	fmt.Printf("node:%s kind:%s\n", id, kind)

	// Archive connections
	for _, connection := range a.metadata[kind].Conns {
		if a.CanFollow(id, kind, connection) {
			if err := a.archiveConnection(ctx, id, connection); err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *Archiver) CanFollow(id, kind, connection string) bool {
	if id == "me" {
		return connection == "photos" || connection == "albums"
	}
	switch kind {
	case "album":
		return connection == "photos"
	default:
		return false
	}
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
			// TODO: Pass kind to archiveNode so we don't have to make two requests
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
	arc, err := LoadArchive(filepath.Join("archive", "facebook"))
	if err != nil {
		return err
	}
	a.archive = arc

	// Build out the kind cache
	a.kindCache = kindCache{
		photos: map[string]struct{}{},
		albums: map[string]struct{}{},
		users:  map[string]struct{}{},
	}
	for _, node := range a.archive.Photos {
		a.kindCache.photos[node.Id] = struct{}{}
	}
	for _, node := range a.archive.Albums {
		a.kindCache.albums[node.Id] = struct{}{}
	}
	for _, node := range a.archive.Friends {
		a.kindCache.users[node.Id] = struct{}{}
	}
	if a.archive.Me != nil {
		a.kindCache.me = a.archive.Me.Id
	}

	err = a.archiveNode(ctx, "me")
	// best effort to save
	a.save()
	return err
}
