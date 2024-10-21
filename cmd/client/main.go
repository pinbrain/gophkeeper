package main

import (
	"fmt"

	"github.com/pinbrain/gophkeeper/internal/client"
	"github.com/pinbrain/gophkeeper/internal/client/config"
)

func main() {
	cfg, err := config.InitConfig()
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
