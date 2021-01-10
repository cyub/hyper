package main

import (
	"log"
	"os"

	"github.com/cyub/hyper/cli"
)

func main() {
	err := cli.Run()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
