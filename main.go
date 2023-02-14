package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"

	//"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://eth-mainnet.g.alchemy.com/v2/NuT8WSwUgOJlbibA3l6OaK2vVo2COvl8")
	if err != nil {
		log.Fatal(err)
	}

	// Get the balance of an account
	account := common.HexToAddress("0x9c5083dd4838e120dbeac44c052179692aa5dac5")
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Account balance(WEI): ", balance) // 25893180161173005034
	balanceInEth := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
	fmt.Printf("Account balance(ETH):  %v ETH\n", balanceInEth)
	//block, err := client.BlockByNumber(context.Background(), nil)
	// get the current Ethereum price in USD from Coingecko API
	response, err := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd")
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var priceMap map[string]map[string]float64
	err = json.Unmarshal(body, &priceMap)
	if err != nil {
		log.Fatal(err)
	}

	// get the current Ethereum price in USD
	ethPriceInUSD := priceMap["ethereum"]["usd"]

	// convert the balance from Ether to USD
	balanceInUSD := new(big.Float).Mul(balanceInEth, big.NewFloat(ethPriceInUSD))

	// print the balance in Ether and USD
	fmt.Printf("Balance in Ether: %v ETH\n", balanceInEth)
	fmt.Printf("Balance in USD: $%v USD\n", balanceInUSD.Text('f', 2))
}
