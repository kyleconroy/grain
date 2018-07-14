package twitter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/golang/protobuf/jsonpb"
	"github.com/kyleconroy/grain/gen/twitter"
)

// TODO: Load from a zip file
type archiveEntry struct {
	filename    string
	windowName  string
	archiveName string
}

func LoadOfficialArchive(base string) (*twitterpb.Archive, error) {
	m := jsonpb.Unmarshaler{}
	a := twitterpb.Archive{}

	paths := []archiveEntry{
		{"account-creation-ip.js", "account_creation_ip", "account_creation_ips"},
		{"account.js", "account", "accounts"},
		{"ad-engagements.js", "ad_engagements", ""},
		{"ad-impressions.js", "ad_impressions", ""},
		{"ageinfo.js", "ageinfo", "ageinfos"},
	}

	for _, e := range paths {
		path := filepath.Join(base, e.filename)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		input, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		name := e.archiveName
		if name == "" {
			name = e.windowName
		}

		// TODO: Make this a regular expression
		output := bytes.Replace(input,
			[]byte("window.YTD."+e.windowName+".part0 = "),
			[]byte("{\""+name+"\": "),
			-1)
		output = append(output, "}"...)

		if err := m.Unmarshal(bytes.NewReader(output), &a); err != nil {
			fmt.Println(err)
			return nil, err
		}
	}
	return &a, nil
}
