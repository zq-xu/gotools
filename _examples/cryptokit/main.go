package main

import (
	"flag"
	"log"

	"github.com/zq-xu/gotools/configx"
	"github.com/zq-xu/gotools/cryptox"
	"github.com/zq-xu/gotools/logx"
)

// cryptox 不再自己读取配置，需要把 "crypto" 配置段接上：
// 引用 storex/account 时会自动注册，这里是用 cryptox 单独时的写法。
func init() {
	configx.Register("crypto", &cryptox.CryptoConfig, func() error {
		return cryptox.Init(cryptox.CryptoConfig)
	})
}

func main() {
	input := flag.String("s", "", "the string to encrypt")
	flag.Parse()

	if *input == "" {
		log.Fatal("please provide the string to encrypt with the flag -s")
	}

	if len(*input) >= 20 {
		log.Fatal("The length of the input should be less than 20")
	}

	err := configx.Setup("config.yaml")
	if err != nil {
		log.Fatalf("failed to setup. %s", err)
	}

	str, err := cryptox.Encrypt([]byte(*input))
	if err != nil {
		log.Fatal("invalid str: ", err)
	}

	logx.Logger.Infoln(str)
}
