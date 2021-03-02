package modules

import (
	"context"
	"fmt"
	"github.com/3hajk/ethereum-cli/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetTxByHash(client *ethclient.Client, hash common.Hash) *models.Transaction {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	tx, pending, err := client.TransactionByHash(context.Background(), hash)
	if err != nil {
		fmt.Println(err)
	}
	return &models.Transaction{
		Hash:     tx.Hash().String(),
		Value:    tx.Value().String(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice().Uint64(),
		Nonce:    tx.Nonce(),
		To:       tx.To().String(),
		Pending:  pending,
	}
}
