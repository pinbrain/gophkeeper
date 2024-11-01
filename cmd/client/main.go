package main

import (
	"fmt"

	"github.com/pinbrain/gophkeeper/internal/client"
	"github.com/pinbrain/gophkeeper/internal/client/config"
)

var (
	Version   = "N/A"
	BuildTime = "N/A"
)

func main() {
	cfg, err := config.InitConfig(Version, BuildTime)
	if err != nil {
		panic(err)
	}
	client, err := client.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	if err = client.Execute(); err != nil {
		fmt.Println(err)
	}
}
