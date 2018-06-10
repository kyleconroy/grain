package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kyleconroy/grain/twitter"
	toml "github.com/pelletier/go-toml"
)

func run() error {
	config, err := toml.LoadFile("config.toml")
	if err != nil {
		return err
	}

	if config.Has("twitter") {
		config, ok := config.Get("twitter").(*toml.Tree)
		if !ok {
			return fmt.Errorf("Config file should contain a [twitter] section")
		}
		a, err := twitter.NewArchiver(config)
		if err != nil {
			return err
		}
		if err := a.Sync(context.TODO()); err != nil {
			return err
		}
	}

	if config.Has("facebook") {
		// config := config.Get("facebook").(*toml.Tree)
		// a := facebook.NewArchiver(config, db, fs)

		// if err := a.Sync(context.TODO()); err != nil {
		// return err
		// }

		// if err := a.Parse(context.TODO()); err != nil {
		// return err
		// }
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
