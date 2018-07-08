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

func LoadOfficialArchive(base string) (*twitterpb.Archive, error) {
	m := jsonpb.Unmarshaler{}
	a := twitterpb.Archive{}
	paths := []string{
		"tweet.js",
	}

	for _, filename := range paths {
		path := filepath.Join(base, filename)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		input, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		// TODO: Make this a regular expression
		output := bytes.Replace(input,
			[]byte("window.YTD.tweet.part0 = "),
			[]byte("{\"timeline\": "),
			-1)
		output = append(output, "}"...)

		if err := m.Unmarshal(bytes.NewReader(output), &a); err != nil {
			fmt.Println(err)
			return nil, err
		}
	}
	return &a, nil
}
