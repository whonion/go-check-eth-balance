package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Read addresses from file
	addresses, err := readAddressesFromFile("addresses.txt")
	if err != nil {
		log.Fatal(err)
	}

	client, err := ethclient.Dial("https://mainnet.infura.io/v3/{YOUR_API_KEY}")
	if err != nil {
		log.Fatal(err)
	}

	// Process each address
	for _, address := range addresses {
		// Get the balance of an account
		account := common.HexToAddress(address)
		balance, err := client.BalanceAt(context.Background(), account, nil)
		if err != nil {
			log.Printf("Error retrieving balance for address %s: %s\n", address, err)
			continue
		}

		balanceInEth := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))

		// Get the current Ethereum price in USD from Coingecko API
		response, err := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd")
		if err != nil {
			log.Printf("Error retrieving Ethereum price: %s\n", err)
			continue
		}

		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Printf("Error reading Coingecko API response: %s\n", err)
			continue
		}

		var priceMap map[string]map[string]float64
		err = json.Unmarshal(body, &priceMap)
		if err != nil {
			log.Printf("Error parsing Coingecko API response: %s\n", err)
			continue
		}

		// Get the current Ethereum price in USD
		ethPriceInUSD := priceMap["ethereum"]["usd"]

		// Convert the balance from Ether to USD
		balanceInUSD := new(big.Float).Mul(balanceInEth, big.NewFloat(ethPriceInUSD))

		// Print the balance in Ether and USD
		fmt.Printf("Address: %s\n", address)
		fmt.Printf("Balance in Ether: %v ETH\n", balanceInEth)
		fmt.Printf("Balance in USD: $%v USD\n", balanceInUSD.Text('f', 2))
		fmt.Println("------------------------------------------------")
	}
}

// Read addresses from a file
func readAddressesFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var addresses []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		addresses = append(addresses, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
}
