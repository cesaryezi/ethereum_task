package task

import (
	"context"
	"os"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

const (
	SolanaWss = "SOLANA_WSS"
	PublicKey = "Public_KEY"
)

func subscribeBlockLog() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, err := ws.Connect(ctx, os.Getenv(SolanaWss))
	if err != nil {
		panic(err)
	}
	defer c.Close()

	pub := []byte(os.Getenv(PublicKey))
	logs, err := c.LogsSubscribeMentions(
		solana.PublicKeyFromBytes(pub),
		rpc.CommitmentConfirmed,
	)
	if err != nil {
		panic(err)
	}
	defer logs.Unsubscribe()
	for {
		select {
		case <-ctx.Done():
			return
		case log, ok := <-logs.Response():
			if !ok {
				return
			}
			println(log.Value.Logs)
		case err := <-logs.Err():
			panic(err)
		}
	}
}
