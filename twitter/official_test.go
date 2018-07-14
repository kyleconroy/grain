package twitter

import (
	"fmt"
	"testing"
)

func TestLoadArchive(t *testing.T) {
	a, err := LoadOfficialArchive("/Users/kyle/Documents/Archives/Twitter/twitter-2018-07-04-74651e7cc75c9f3cf916df6cd7706a2f9bb4ba963cc2aba7edb300fa6c368439")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(a)
}
