package main

import (
	"flag"
	"log"

	"github.com/zq-xu/gotools"
)

func main() {
	input := flag.String("s", "", "the string to encrypt")
	flag.Parse()

	if *input == "" {
		log.Fatal("please provide the string to encrypt with the flag -s")
	}

	if len(*input) >= 20 {
		log.Fatal("The length of the input should be less than 20")
	}

	err := gotools.Setup("config.yaml")
	if err != nil {
		log.Fatalf("failed to setup. %s", err)
	}

	str, err := gotools.Encrypt([]byte(*input))
	if err != nil {
		log.Fatal("invaild str", err)
	}

	gotools.Logger.Infoln(str)
}
