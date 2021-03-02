package modules

import (
	"context"
	"fmt"
	"github.com/3hajk/ethereum-cli/models"
	"github.com/3hajk/ethereum-cli/store"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

func GetLatestBlock(client *ethclient.Client) *models.Block {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	header, _ := client.HeaderByNumber(context.Background(), nil)
	blockNumber := big.NewInt(header.Number.Int64())

	b, _ := GetBlockByNumber(client, blockNumber)

	return b
}

func GetBlocks(client *ethclient.Client, from, to int64) []*models.Block {
	blocks := make([]*models.Block, 0)

	for i := from; i <= to; i++ {
		blockNumber := big.NewInt(i)
		block, _ := GetBlockByNumber(client, blockNumber)
		if block != nil {
			blocks = append(blocks, block)
		}
	}
	return blocks
}

func GetBlocksAsync(client *ethclient.Client, store *store.TransactionData, from, to int64) {
	blocks := make([]*models.Block, 0)
	for i := from; i <= to; i++ {
		blockNumber := big.NewInt(i)
		block, _ := GetBlockByNumber(client, blockNumber)
		if block != nil {
			blocks = append(blocks, block)
		}
	}
	txCount := 0
	for _, block := range blocks {
		for _, tx := range block.Transactions {
			txCount++
			//store.Set(block.Timestamp, &tx)
			store.Add(&tx)
		}
	}
	if txCount > 0 {
		store.SetReady()
	}
}

func GetBlockByNumber(client *ethclient.Client, blockNumber *big.Int) (*models.Block, error) {
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		return nil, err
	}
	b := &models.Block{
		BlockNumber:      block.Number().Int64(),
		Timestamp:        block.Time(),
		Difficulty:       block.Difficulty().Uint64(),
		Hash:             block.Hash().String(),
		TransactionCount: len(block.Transactions()),
		Transactions:     []models.Transaction{},
	}
	for _, tx := range block.Transactions() {
		if tx == nil {
			continue
		}
		to := ""
		if tx.To() != nil {
			to = tx.To().String()
		}

		b.Transactions = append(b.Transactions,
			models.Transaction{
				Hash:     tx.Hash().String(),
				Value:    tx.Value().String(),
				Gas:      tx.Gas(),
				GasPrice: tx.GasPrice().Uint64(),
				Nonce:    tx.Nonce(),
				To:       to,
			})
	}
	return b, nil
}
