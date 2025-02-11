package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	infuraURL     = "https://sepolia.infura.io/v3/YOUR_INFURA_PROJECT_ID"
	privateKeyHex = "YOUR_PRIVATE_KEY"
	contractAddr  = "YOUR_DEPLOYED_CONTRACT_ADDRESS"
)

func main() {
	// Connect to Ethereum network
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum: %v", err)
	}
	defer client.Close()

	// Load private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	// Get sender address
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKey)

	// Get nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	// Set gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to get gas price: %v", err)
	}

	// Create transaction options
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(11155111)) // Sepolia Chain ID
	if err != nil {
		log.Fatalf("Failed to create auth: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(3000000)

	// Load contract instance
	voting, err := NewVoting(common.HexToAddress(contractAddr), client)
	if err != nil {
		log.Fatalf("Failed to load contract: %v", err)
	}

	// Vote for a candidate
	candidate := "Alice"
	tx, err := voting.Vote(auth, candidate)
	if err != nil {
		log.Fatalf("Failed to vote: %v", err)
	}

	fmt.Printf("Voted for %s! Transaction Hash: %s\n", candidate, tx.Hash().Hex())

	// Retrieve votes
	voteCount, err := voting.GetVotes(nil, candidate)
	if err != nil {
		log.Fatalf("Failed to get votes: %v", err)
	}

	fmt.Printf("%s has %d votes!\n", candidate, voteCount)
}
