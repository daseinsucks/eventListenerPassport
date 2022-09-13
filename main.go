package main

import (
	"context"

	"log"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	passport "github.com/MoonSHRD/IKY-telegram-bot/artifacts/TGPassport"
	"github.com/ethereum/go-ethereum/core/types"
)

func main() {
	contractAddress := common.HexToAddress("contractAddress")
	Client, err := ethclient.Dial("wss://rinkeby.infura.io/ws/v3/")
	if err != nil {
		log.Fatal(err)
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	var ch = make(chan types.Log)
	ctx := context.Background()

	sub, err := Client.SubscribeFilterLogs(ctx, query, ch)

	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(passport.PassportABI)))

	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case eventLog := <-ch:

			event, err := contractAbi.Unpack("passportApplied", eventLog.Data)

			if err != nil {
				log.Println("Failed to unpack")
				continue
			}

			applyerTg := event[0]
			applyerAddress := event[1]

			log.Println("TelegramID:", applyerTg)
			log.Println("Applyer address:", applyerAddress)

		}
	}
}
