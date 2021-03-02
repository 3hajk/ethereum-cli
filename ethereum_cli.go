package main

import (
	"flag"
	"fmt"
	Handlers "github.com/3hajk/ethereum-cli/handler"
	"github.com/3hajk/ethereum-cli/store"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var (
	eth  = flag.String("eht", "https://cloudflare-eth.com/", "")
	port = flag.Int("port", 8080, "")
)

func main() {
	flag.Parse()
	client, err := ethclient.Dial(*eth)
	if err != nil {
		log.Fatal(err)
	}
	r := mux.NewRouter()
	d := store.InitData()
	r.Handle("/api/eth/{module}", Handlers.ClientHandler{
		Client: client,
		Data:   d,
	})
	log.Printf("Listening on port %d", *port)
	log.Printf("Open http://localhost:%d/api/eht/{cmd}?data in the browser", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), r))
}
