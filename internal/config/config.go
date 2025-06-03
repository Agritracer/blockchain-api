package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RPCUrl     string
	PrivateKey string
	EthscanAPI string
	Wallet     string
	Port       string
}

var Cfg Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Không tìm thấy .env, đang dùng biến môi trường hệ thống (nếu có)")
	}

	Cfg = Config{
		RPCUrl:     os.Getenv("RPC_URL"),
		PrivateKey: os.Getenv("PRIVATE_KEY"),
		Port:       os.Getenv("PORT"),
		EthscanAPI: os.Getenv("ETHERSCAN_API_KEY"),
		Wallet:     os.Getenv("TRACE_WALLET_ADDRESS"),
	}
}
