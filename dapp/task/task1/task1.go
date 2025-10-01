package task1

import (
	"context"
	"crypto/ecdsa"
	_ "embed"
	"ethereum_task/dapp/task/task2"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/crypto/sha3"
)

const (
	SepoliaHttp  = "SEPOLIA_HTTP"
	SepoliaWss   = "SEPOLIA_WSS"
	PrivateKey   = "PRIVATE_KEY"
	ToAddress    = "TO_ADDRESS"
	TokenAddress = "TOKEN_ADDRESS"
	ContractAddr = "CONTRACT_ADDR"
)

func queryBlock() {

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//1 创建ETH客户端
	client, err := ethclient.DialContext(ctx, os.Getenv(SepoliaHttp))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	//2 通过区块高度获取区块信息
	//blockNumber := big.NewInt(5671744)
	blockNumber := big.NewInt(int64(rpc.LatestBlockNumber))
	block, err := client.BlockByNumber(ctx, blockNumber) //HeaderByNumber
	if err != nil {
		log.Fatalf("Failed to get block: %v", err)
	}

	fmt.Println("block info:", block.Number().Uint64()) //对象转换为具体数值
	fmt.Println("block info:", block.Hash().Hex())
	fmt.Println("block info:", block.Time())
	fmt.Println("block info:", len(block.Transactions()))

	fmt.Println("-------------------------")
	/*fmt.Println("block info:", block.Number())
	fmt.Println("block info:", block.Hash())
	fmt.Println("block info:", block.Time())
	fmt.Println("block info:", len(block.Transactions()))*/

	//第二种获取交易数量  获取单个区块内交易数量：  len(block.Transactions())
	count, err := client.TransactionCount(ctx, block.Hash())
	if err != nil {
		log.Fatalf("Failed to get transaction count: %v", err)
	}

	fmt.Println("block info:", count)

}

func doTransaction() {
	// 创建上下文
	ctx := context.Background()

	//1 创建ETH客户端
	client, err := ethclient.DialContext(ctx, os.Getenv(SepoliaHttp))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	//2 构造交易参数：types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	//2.1 获取私钥 解密 出 公钥：注意使用   ECDSA
	privateKeyHex := os.Getenv(PrivateKey)
	if privateKeyHex == "" {
		log.Fatal("PRIVATE_KEY environment variable is not set")
	}
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Invalid private key: %v", err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey) //断言 转化为ecdsa.PublicKey
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	// 2.2 使用公钥获取发送方地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	//2.3 获取发送方获取账户在待处理状态下的 nonce 值  client.NonceAt用于查看历史的
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		log.Fatalf("Failed to get pending nonce: %v", err)
	}
	//2.4 获取建议的 gasPrice
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}
	//2.5 交易接收者
	toAddressStr := os.Getenv(ToAddress)
	toAddress_ := common.HexToAddress(toAddressStr)
	if toAddress_ == (common.Address{}) {
		log.Fatal("invalid receiver address")
	}
	//2.6 设置可接受的gas最大值 gasLimit
	gasLimit := uint64(21000)
	//2.7 设置交易金额
	//value := big.NewInt(1000000000000000000) // 1 ETH
	//value := big.NewInt(1e18) // 1 ETH
	value := big.NewInt(100) // 100 wei
	//2.8 添加数据
	data := []byte("")

	// 2.9 组装交易参数
	//tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress_,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})

	// 3 签名交易
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatalf("Failed to get network ID: %v", err)
	}
	signer := types.NewEIP155Signer(chainID)
	signedTx, err := types.SignTx(tx, signer, privateKey) // 签名交易：交易参数数据 + chain + 私钥
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	//4 广播交易：提交交易
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	//5 打印交易Hash
	txHash := signedTx.Hash().Hex()
	fmt.Println("tx sent:", txHash)

	//6  等待交易被打包，被确认交易，并查看结果（之前SendTransaction只保证这笔交易提交成功，但是不保证打包 上链 ）
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

func buildPriKeyAndFromAddr() (*ecdsa.PrivateKey, common.Address) {
	//2.1 获取私钥 解密 出 公钥：注意使用   ECDSA
	privateKeyHex := os.Getenv(PrivateKey)
	if privateKeyHex == "" {
		log.Fatal("PRIVATE_KEY environment variable is not set")
	}
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Invalid private key: %v", err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey) //断言 转化为ecdsa.PublicKey
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	// 2.2 使用公钥获取发送方地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return privateKey, fromAddress
}

func createWallAddr() {

	//公私钥 都需要 ecdsa

	privateKey, _ := crypto.GenerateKey()
	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println(hexutil.Encode(privateKeyBytes)[2:]) // 去掉'0x'

	//派生公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	addr := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("address:", addr.Hex())

	//以下是自己手动使用pub生成 地址
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println("from pubKey:", hexutil.Encode(publicKeyBytes)[4:]) // 去掉'0x04'
	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println("manu address:", hexutil.Encode(hash.Sum(nil)[12:]))

}

func transferToken() {

	// 创建上下文
	ctx := context.Background()

	//1 创建ETH客户端
	client, err := ethclient.DialContext(ctx, os.Getenv(SepoliaHttp))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	privateKey, fromAddress := buildPriKeyAndFromAddr()

	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		log.Fatalf("Failed to get pending nonce: %v", err)
	}
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}

	toAddressStr := os.Getenv(ToAddress)
	toAddress_ := common.HexToAddress(toAddressStr)

	tokenAddressStr := os.Getenv(TokenAddress)
	tokenAddress_ := common.HexToAddress(tokenAddressStr)

	//代币注意通用的ABI编码：
	signatureStr := []byte("transfer(address,uint256)")
	//交易使用Keccak256进行hash
	hash := sha3.NewLegacyKeccak256()
	hash.Write(signatureStr)
	methodID := hash.Sum(nil)[:4] //111   transfer
	//太坊 ABI 编码的要求，所有参数都需要是 32 字节的倍数
	data := append(methodID, common.LeftPadBytes(tokenAddress_.Bytes(), 32)...) //111  address
	amount := new(big.Int)
	amount.SetString("1000000000000000000000", 10) // 1000 tokens
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	data = append(data, paddedAmount...) //111 uint256

	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
		To:   &toAddress_,
		Data: data,
	})
	if err != nil {
		log.Fatalf("Failed to estimate gas: %v", err)
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress_,
		Value:    big.NewInt(0), // 0 wei
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

}

func getEthBalance() {

	// 创建上下文
	ctx := context.Background()

	//1 创建ETH客户端
	client, err := ethclient.DialContext(ctx, os.Getenv(SepoliaHttp))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	_, fromAddress := buildPriKeyAndFromAddr()

	balance, err := client.BalanceAt(ctx, fromAddress, nil)
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}
	fmt.Println("balance:", balance)

	//要读取 ETH 值，您必须做计算 wei/10^18,应为 balance 的单位是wei
	fbalance := new(big.Float)
	//防止精度丢失 ，需要 SetString
	fbalance.SetString(balance.String())
	// fbalance  /   big.NewFloat(1e18)
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(1e18))
	fmt.Println("eth:", ethValue)

}

func getTokenBalance() {

	// 创建上下文
	ctx := context.Background()

	//1 创建ETH客户端
	client, err := ethclient.DialContext(ctx, os.Getenv(SepoliaHttp))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	_, fromAddress := buildPriKeyAndFromAddr()

	instance, err := task2.NewTask(common.HexToAddress(os.Getenv(ContractAddr)), client)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("instance:", instance)
	fmt.Println("fromAddress:", fromAddress)

	//bal, err := instance.BalanceOf(&bind.CallOpts{}, fromAddress)//实现ERC20的balanceOf方法

}

func SubscribeBlock() {

	// 创建上下文
	ctx := context.Background()

	//1 创建ETH客户端
	client, err := ethclient.DialContext(ctx, os.Getenv(SepoliaWss))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	ch := make(chan *types.Header)

	sub, err := client.SubscribeNewHead(ctx, ch)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-ch:
			hash := header.Hash()
			fmt.Println("block:", header.Number.Int64(), "hash:", hash.Hex())
			block, err := client.BlockByHash(ctx, hash)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("block info:", block.Number().Uint64()) //对象转换为具体数值
			fmt.Println("block info:", block.Hash().Hex())
			fmt.Println("block info:", block.Time())
			fmt.Println("block info:", len(block.Transactions()))
		}
	}

}

//go:embed "../task2/Counter_sol_Counter.abi"
var counterAbi string

func queryContractEvent() {

	// 创建上下文
	ctx := context.Background()

	//1 创建ETH客户端
	client, err := ethclient.DialContext(ctx, os.Getenv(SepoliaWss))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	query := ethereum.FilterQuery{
		//FromBlock: big.NewInt(rpc.LatestBlockNumber.Int64()),
		Addresses: []common.Address{
			common.HexToAddress(os.Getenv(ContractAddr)),
		},
		/*Topics: [][]common.Hash{
			{
				common.HexToHash(""),
			},
		},*/
	}

	logs, err := client.FilterLogs(ctx, query)
	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(counterAbi))
	if err != nil {
		log.Fatal(err)
	}
	for _, vLog := range logs {
		fmt.Println("block:", vLog.BlockNumber, "tx:", vLog.TxHash.Hex())
		fmt.Println("block:", vLog.Address.Hex())
		fmt.Println("block:", vLog.Topics)
		fmt.Println("block:", vLog.Data)

		//解析事件
		event := struct{}{}
		err := contractAbi.UnpackIntoInterface(event, "EventName", vLog.Data)
		if err != nil {
			log.Fatal(err)
		}

	}

}

func SubscribeEvents() {

	// 创建上下文
	ctx := context.Background()

	//1 创建ETH客户端
	client, err := ethclient.DialContext(ctx, os.Getenv(SepoliaWss))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(counterAbi))
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan types.Log)

	query := ethereum.FilterQuery{
		//FromBlock: big.NewInt(rpc.LatestBlockNumber.Int64()),
		Addresses: []common.Address{
			common.HexToAddress(os.Getenv(ContractAddr)),
		},
		/*Topics: [][]common.Hash{
			{
				common.HexToHash(""),
			},
		},*/
	}

	subs, err := client.SubscribeFilterLogs(ctx, query, ch)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-subs.Err():
			log.Fatal(err)
		case vLog := <-ch:
			fmt.Println("block:", vLog.BlockNumber, "tx:", vLog.TxHash.Hex())
			fmt.Println("block:", vLog.Address.Hex())
			fmt.Println("block:", vLog.Topics)
			fmt.Println("block:", vLog.Data)

			//解析事件
			event := struct {
				Key   [32]byte
				Value []byte
			}{}
			err := contractAbi.UnpackIntoInterface(event, "EventName", vLog.Data)
			//contractAbi.Pack(methodName, key, value)//使用ABI进行编码，进行transation的data数据，进行执行
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("event:", event)

			fmt.Println(common.Bytes2Hex(event.Key[:]))
			fmt.Println(common.Bytes2Hex(event.Value[:]))
			var topics []string
			for i := range vLog.Topics {
				topics = append(topics, vLog.Topics[i].Hex())
			}
			fmt.Println("topics[0]=", topics[0])
			if len(topics) > 1 {
				fmt.Println("index topic:", topics[1:])
			}
		}
	}

}
