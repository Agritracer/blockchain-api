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
	ApiPort    string
	WebPort    string
	APIKey     string
	TOTPSecret string
	JWTToken   string
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
		ApiPort:    os.Getenv("API_PORT"),
		WebPort:    os.Getenv("WEB_PORT"),
		EthscanAPI: os.Getenv("ETHERSCAN_API_KEY"),
		Wallet:     os.Getenv("TRACE_WALLET_ADDRESS"),
		APIKey:     os.Getenv("API_KEY"),
		TOTPSecret: os.Getenv("TOTP_SECRET"),
		JWTToken:   os.Getenv("JWT_TOKEN"),
	}
}
