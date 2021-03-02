# ethereum-cli


Run:
go run ethereum_cli.go

Use:
open url: http://localhost:8080/api/eth/{cmd}
example: 
Getting tx from ethereum
sync
http://localhost:8080/api/eth/get-block?fromBlock=11953428&toBlock=11953484
async 1 thread
http://localhost:8080/api/eth/get-block?fromBlock=11953428&toBlock=11953484&async=true
async some thread
http://localhost:8080/api/eth/get-block?fromBlock=11953428&toBlock=11953484&async=true&thread=2



Get top address
http://localhost:8080/api/eth/top-address?count=8
Get tx by address
http://localhost:8080/api/eth/get-tx-by-address?address=0xd9e1cE17f2641f24aE83637ab66a2cca9C378B9F
Count getting tx
http://localhost:8080/api/eth/tx-count
Clean data
http://localhost:8080/api/eth/clean