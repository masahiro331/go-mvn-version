package main

import (
	"fmt"
	"log"

	"github.com/masahiro331/go-mvn-version/pkg/version"
)

func main() {
	v1, err := version.NewVersion("1.alpha")
	if err != nil {
		log.Fatal(err)
	}
	v2, err := version.NewVersion("1")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(v1.GreaterThan(*v2))
}
