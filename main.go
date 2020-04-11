package main

import (
	"fmt"
	"log"

	"github.com/masahiro331/go-mvn-version/pkg/version"
)

func main() {
	v1, err := version.NewVersion("5.0.1.RELEASE")
	if err != nil {
		log.Fatal(err)
	}
	v2, err := version.NewVersion("5.0.0.RELEASE")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(v1.GreaterThan(*v2))
}
