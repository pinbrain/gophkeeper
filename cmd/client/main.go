package main

import (
	"github.com/pinbrain/gophkeeper/internal/client"
)

func main() {
	if err := client.Run(); err != nil {
		panic(err)
	}
}
