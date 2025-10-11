package main

import (
	"flag"
	"fmt"
	"log"

	"zq-xu/gotools/bricks/cryptokit"
	"zq-xu/gotools/setup"
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

	err := setup.Setup()
	if err != nil {
		log.Fatalf("failed to setup. %s", err)
	}

	str, err := cryptokit.Crypto.Encrypt([]byte(*input))
	if err != nil {
		log.Fatal("invaild str", err)
	}

	fmt.Println(str)
}
