package eth

import (
	"context"
	"fmt"
	"math/big"

	"agritrace/internal/config"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func SendToEthereum(data string) (string, error) {
	config.LoadConfig()
	rpcURL := config.Cfg.RPCUrl
	privateKeyHex := config.Cfg.PrivateKey

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return "", err
	}
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", err
	}
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	nonce, _ := client.PendingNonceAt(context.Background(), fromAddress)
	gasPrice, _ := client.SuggestGasPrice(context.Background())
	gasLimit := uint64(50000)

	tx := types.NewTransaction(nonce, fromAddress, big.NewInt(0), gasLimit, gasPrice, []byte(data))
	chainID, _ := client.NetworkID(context.Background())
	signedTx, _ := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}
	return signedTx.Hash().Hex(), nil
}

func GetDataFromTransaction(txHash string) (string, error) {
	config.LoadConfig()
	rpcURL := config.Cfg.RPCUrl

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return "", err
	}

	tx, _, err := client.TransactionByHash(context.Background(), common.HexToHash(txHash))
	if err != nil {
		return "", err
	}

	data := tx.Data()
	if len(data) == 0 {
		return "", fmt.Errorf("giao dịch không chứa dữ liệu")
	}

	return string(data), nil
}
