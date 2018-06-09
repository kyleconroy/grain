package archive

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

func ArchiveURL(ctx context.Context, service, folder, uri string) error {
	parsed, err := url.Parse(uri)
	if err != nil {
		return err
	}
	if parsed.Host == "" {
		panic("All URLs must have attached host names: " + uri)
	}

	fullpath := path.Join("archive", service, folder, parsed.Host, parsed.Path)

	if _, err := os.Stat(fullpath); err == nil {
		// path/to/whatever exists
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(fullpath), 0755); err != nil {
		return err
	}

	resp, err := http.Get(parsed.String())
	if err != nil {
		return err
	}

	// For now, ignore everything else that isn't 200
	if resp.StatusCode != http.StatusOK {
		return nil
	}

	defer resp.Body.Close()
	blob, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fullpath, blob, 0644)
}
