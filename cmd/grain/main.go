package main

import (
	"context"
	"log"

	"github.com/kyleconroy/grain/twitter"
	toml "github.com/pelletier/go-toml"
)

func main() {
	config, err := toml.LoadFile("config.toml")
	if err != nil {
		log.Fatal(err)
	}

	if true {
		config := config.Get("twitter").(*toml.Tree)
		a := twitter.NewArchiver(config)

		if err := a.Sync(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}

	if false {
		// config := config.Get("facebook").(*toml.Tree)
		// a := facebook.NewArchiver(config, db, fs)

		// if err := a.Sync(context.TODO()); err != nil {
		// 	log.Fatal(err)
		// }

		// if err := a.Parse(context.TODO()); err != nil {
		// 	log.Fatal(err)
		// }
	}
}
