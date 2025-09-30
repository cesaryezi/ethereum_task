package task1

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	SepoliaHttp = "SEPOLIA_HTTP"
	PrivateKey  = "PRIVATE_KEY"
	toAddress   = "TO_ADDRESS"
)

func queryBlock() {

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, os.Getenv(SepoliaHttp))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// get block
	blockNumber := big.NewInt(5671744)
	block, err := client.BlockByNumber(ctx, blockNumber)
	if err != nil {
		log.Fatalf("Failed to get block: %v", err)
	}

	fmt.Println("block info:", block.Number().Uint64())
	fmt.Println("block info:", block.Hash().Hex())
	fmt.Println("block info:", block.Time())
	fmt.Println("block info:", len(block.Transactions()))

	fmt.Println("-------------------------")
	fmt.Println("block info:", block.Number())
	fmt.Println("block info:", block.Hash())
	fmt.Println("block info:", block.Time())
	fmt.Println("block info:", len(block.Transactions()))

	//获取单个区块内交易数量
	count, err := client.TransactionCount(ctx, block.Hash())
	if err != nil {
		log.Fatalf("Failed to get transaction count: %v", err)
	}

	fmt.Println("block info:", count)

}

func doTransaction() {
	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, os.Getenv(SepoliaHttp))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	privateKeyHex := os.Getenv(PrivateKey)
	if privateKeyHex == "" {
		log.Fatal("PRIVATE_KEY environment variable is not set")
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Invalid private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		log.Fatalf("Failed to get pending nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}

	toAddressStr := os.Getenv(toAddress)
	toAddress_ := common.HexToAddress(toAddressStr)
	if toAddress_ == (common.Address{}) {
		log.Fatal("invalid receiver address")
	}

	gasLimit := uint64(21000)
	value := big.NewInt(1000000000000000000) // 1 ETH
	data := []byte("")

	//tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress_,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatalf("Failed to get network ID: %v", err)
	}

	signer := types.NewEIP155Signer(chainID)
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	txHash := signedTx.Hash().Hex()
	fmt.Println("tx sent:", txHash)

	// 等待交易被打包并查看结果
	receipt, err := bind.WaitMined(ctx, client, signedTx)
	if err != nil {
		log.Printf("Failed to wait for transaction to be mined: %v", err)
		return
	}
	if receipt.Status == types.ReceiptStatusSuccessful {
		fmt.Println("Transaction succeeded!")
	} else {
		fmt.Println("Transaction failed.")
	}
}
