package task

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"os"
	"time"

	"github.com/blocto/solana-go-sdk/client"
	solCom "github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/system"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	SolanaHttp = "SOLANA_HTTP"
	SolanaWss  = "SOLANA_WSS"
	PrivateKey = "PRIVATE_KEY"
	ToAddress  = "TO_ADDRESS"
)

func query() {

	ctx := context.Background()

	clientSolana := client.NewClient(rpc.TestnetRPCEndpoint)

	block, err := clientSolana.GetBlock(ctx, 0) // 0表示最新区块
	if err != nil {
		panic(err)
	}

	fmt.Println("block info:", block.BlockHeight)
	fmt.Println("block info:", block.BlockTime)
	fmt.Println("block info:", block.ParentSlot)
	fmt.Println("block info:", len(block.Transactions))

}

func execTransaction() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientSolana := client.NewClient(rpc.TestnetRPCEndpoint)

	privateKey, _, err := buildPriKeyAndFromAddr()
	if err != nil {
		return fmt.Errorf("failed to build private key: %w", err)
	}

	// 注意：这里需要将 Ethereum ECDSA 密钥转换为 Solana Ed25519 密钥
	// 当前实现可能会失败，需要根据实际情况调整密钥处理逻辑
	privateKeyBytes := crypto.FromECDSA(privateKey)
	payer, err := types.AccountFromBytes(privateKeyBytes)
	if err != nil {
		return fmt.Errorf("failed to create account from bytes: %w", err)
	}

	recentBlockHash, err := clientSolana.GetLatestBlockhash(ctx)
	if err != nil {
		return fmt.Errorf("failed to get recent blockhash: %w", err)
	}

	toPubKey := solCom.PublicKeyFromString(os.Getenv(ToAddress))
	// 交易指令:特别重要！！！！！！！核心
	transferInstruction := system.Transfer(system.TransferParam{
		From: payer.PublicKey,
		To:   toPubKey,
		//To:     types.NewAccount().PublicKey,
		Amount: 1_000_000, // 0.01 SOL (lamport单位)
	})

	//构建交易
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{payer},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        payer.PublicKey,
			RecentBlockhash: recentBlockHash.Blockhash,
			Instructions:    []types.Instruction{transferInstruction},
		}),
	})
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	//发送交易
	sig, err := clientSolana.SendTransactionWithConfig(
		ctx,
		tx,
		client.SendTransactionConfig{
			SkipPreflight:       true,
			PreflightCommitment: rpc.CommitmentFinalized,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	fmt.Println("tx sent:", sig)
	return nil

}

func buildPriKeyAndFromAddr() (*ecdsa.PrivateKey, common.Address, error) {
	privateKeyHex := os.Getenv(PrivateKey)
	if privateKeyHex == "" {
		return nil, common.Address{}, fmt.Errorf("PRIVATE_KEY environment variable is not set")
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("invalid private key: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, common.Address{}, fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return privateKey, fromAddress, nil
}
