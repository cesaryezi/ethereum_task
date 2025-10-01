package task2

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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	SepoliaHttp_2  = "SEPOLIA_HTTP"
	PrivateKey_2   = "PRIVATE_KEY"
	ContractAddr_2 = "CONTRACT_ADDR"
)

func deployCounterContract() {

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//1 创建ETH客户端
	client, err := ethclient.DialContext(ctx, os.Getenv(SepoliaHttp_2))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// 2 创建部署者
	// 2.1 创建私钥和地址
	privateKey, fromAddress := buildPriKeyAndFromAddr()

	// 2.2 获取nonce
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		log.Fatalf("Failed to get pending nonce: %v", err)
	}

	//2.3 获取建议的 gasPrice
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}

	// 2.4 获取chainID
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatalf("Failed to get network ID: %v", err)
	}

	// 创建部署合约 的 交易参数
	// 2.5 在特定的chain构建交易参数
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	// 设置nonce
	auth.Nonce = big.NewInt(int64(nonce))
	// 设置值
	auth.Value = big.NewInt(0) // in wei
	//  设置gasLimit
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	//3 部署合约
	address, tx, instance, err := DeployTask(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	//4 打印部署信息
	fmt.Println(address.Hex())
	fmt.Println(tx.Hash().Hex())

	_ = instance

}

func execCounterContract() {

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//1 创建ETH客户端
	client, err := ethclient.DialContext(ctx, os.Getenv(SepoliaHttp_2))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	//2 创建合约实例
	task, err := NewTask(common.HexToAddress(os.Getenv(ContractAddr_2)), client)
	if err != nil {
		log.Fatal(err)
	}

	//3 创建调用者
	privateKey, err := crypto.HexToECDSA(os.Getenv(PrivateKey_2))
	if err != nil {
		log.Fatal(err)
	}

	//4 创建TransactOpts交易参数
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatalf("Failed to get network ID: %v", err)
	}
	//TransactOpts 会消耗 gas。改变合约状态
	opt, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}

	//5 发送交易
	tx, err := task.IncreaseOne(opt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("tx hash:", tx.Hash().Hex())

	//6 绑定  获取参数类型  bind.CallOpts此种方法 不会 消耗 gas，不会改变合约状态，只是读取pure,view 方法
	callOpt := bind.CallOpts{Context: context.Background()}
	//7 获取合约中的值
	valueInContract, err := task.GetCount(&callOpt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("value in contract:", valueInContract)

}

func buildPriKeyAndFromAddr() (*ecdsa.PrivateKey, common.Address) {

	privateKeyHex := os.Getenv(PrivateKey_2)
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
	return privateKey, fromAddress
}
