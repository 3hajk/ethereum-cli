package handler

import (
	"encoding/json"
	"fmt"
	"github.com/3hajk/ethereum-cli/models"
	"github.com/3hajk/ethereum-cli/modules"
	"github.com/3hajk/ethereum-cli/store"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ClientHandler struct {
	Client *ethclient.Client
	Data   *store.TransactionData
}

func (c ClientHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	module := vars["module"]
	hash := r.URL.Query().Get("hash")
	address := r.URL.Query().Get("address")
	from := r.URL.Query().Get("fromBlock")
	to := r.URL.Query().Get("toBlock")
	cnt := r.URL.Query().Get("count")
	async := r.URL.Query().Get("async")
	thread := r.URL.Query().Get("thread")
	w.Header().Set("Content-Type", "application/json")
	switch module {
	case "latest-block":
		block := modules.GetLatestBlock(c.Client)
		json.NewEncoder(w).Encode(block)
	case "get-block":
		if from == "" || to == "" {
			json.NewEncoder(w).Encode(&models.Error{
				Code:    404,
				Message: "Malformed request",
			})
		}
		fromBlock, _ := strconv.ParseInt(from, 10, 64)
		toBlock, _ := strconv.ParseInt(to, 10, 64)
		if c.Data.DataIsReady() {
			json.NewEncoder(w).Encode(&models.Error{
				Code:    400,
				Message: fmt.Sprintf("Transactions is loaded. Clean first"),
			})
			return
		}
		if async == "true" {
			if c.Data.IsProcessing() {
				json.NewEncoder(w).Encode(&models.Error{
					Code:    400,
					Message: fmt.Sprintf("Data is processing! wait"),
				})
				return
			}

			c.Data.SetProcessing()
			run := 1
			if thread != "" {
				if i, err := strconv.Atoi(thread); err == nil {
					run = i
				}
			}
			if run > 1 {
				blocks := int(toBlock - fromBlock)
				last := blocks % run
				step := (blocks - last) / run
				for i := 0; i < blocks; i = i + step {
					go modules.GetBlocksAsync(c.Client, c.Data, fromBlock+int64(i), fromBlock+int64(i+step))
				}
				if last > 0 {
					go modules.GetBlocksAsync(c.Client, c.Data, toBlock-int64(last), toBlock)
				}
			} else {
				go modules.GetBlocksAsync(c.Client, c.Data, fromBlock, toBlock)
			}
			json.NewEncoder(w).Encode(&models.Error{
				Code:    200,
				Message: fmt.Sprintf("is processing"),
			})
			return
		}
		blocks := modules.GetBlocks(c.Client, fromBlock, toBlock)
		txCount := 0
		for _, block := range blocks {
			for _, tx := range block.Transactions {
				txCount++
				c.Data.Set(block.Timestamp, &tx)
			}
		}
		if txCount > 0 {
			c.Data.SetReady()
		}
		json.NewEncoder(w).Encode(&models.Error{
			Code:    200,
			Message: fmt.Sprintf("complete get %d blocks, %d transaction", len(blocks), txCount),
		})
	case "get-tx":
		if hash == "" {
			json.NewEncoder(w).Encode(&models.Error{
				Code:    400,
				Message: "Malformed request",
			})
		}
		thHash := common.HexToHash(hash)
		tx := modules.GetTxByHash(c.Client, thHash)
		if tx != nil {
			json.NewEncoder(w).Encode(tx)
		}
		json.NewEncoder(w).Encode(&models.Error{
			Code:    404,
			Message: "Tx Not Found",
		})
	case "top-address":
		count := 10
		if cnt != "" {
			i, err := strconv.Atoi(cnt)
			if err == nil {
				count = i
			}
		}
		data := c.Data.GetTopAddress(count)
		if len(data) > 0 {
			for k, v := range data {
				json.NewEncoder(w).Encode(&models.TopAddress{Address: k, Count: v})
			}
		} else {
			json.NewEncoder(w).Encode(&models.Error{
				Code:    400,
				Message: "Transaction not loaded",
			})
		}
	case "get-tx-by-address":
		if address == "" {
			json.NewEncoder(w).Encode(&models.Error{
				Code:    400,
				Message: "Malformed request",
			})
			return
		}
		if !c.Data.DataIsReady() {
			json.NewEncoder(w).Encode(&models.Error{
				Code:    404,
				Message: "Tx Not Loaded",
			})
			return
		}
		txs, err := c.Data.GetTmByAddress(address)
		if err != nil {
			json.NewEncoder(w).Encode(&models.Error{
				Code:    404,
				Message: "Tx Not Found",
			})
			return
		}
		json.NewEncoder(w).Encode(&models.Error{
			Code:    200,
			Message: "Ok",
		})
		json.NewEncoder(w).Encode(txs)
	case "tx-count":
		count := c.Data.Count()
		json.NewEncoder(w).Encode(&models.Error{
			Code:    200,
			Message: fmt.Sprintf("Cached %d tx", count),
		})
	case "clean":
		c.Data.Clean()
		json.NewEncoder(w).Encode(&models.Error{
			Code:    200,
			Message: "Ok",
		})
	default:
		json.NewEncoder(w).Encode(&models.Error{
			Code:    400,
			Message: "Unknown request! use: [latest-block scan-block get-tx top clean]",
		})
	}

}
