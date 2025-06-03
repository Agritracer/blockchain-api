package eth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const rpcURL = "https://sepolia.infura.io/v3/XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
const privateKeyHex = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

func SendToEthereum(data string) (string, error) {
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
