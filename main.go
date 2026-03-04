package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/evbruno/arbot-go/factory"
	"github.com/evbruno/arbot-go/pair"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Dex struct {
	Name    string
	Factory common.Address
}

type Token struct {
	Symbol   string
	Address  common.Address
	Decimals uint8
}

var (
	USDC = Token{
		Symbol:   "USDC",
		Address:  common.HexToAddress("0xA0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"),
		Decimals: 6,
	}

	WETH = Token{
		Symbol:   "WETH",
		Address:  common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
		Decimals: 18,
	}

	WBTC = Token{
		Symbol:   "WBTC",
		Address:  common.HexToAddress("0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599"),
		Decimals: 8,
	}
)

func main() {
	rpc := os.Getenv("MAINNET_RPC")

	if rpc == "" {
		log.Fatal("Set MAINNET_RPC env variable")
	}

	client, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatal(err)
	}

	dexes := []Dex{
		{
			Name:    "UniswapV2",
			Factory: common.HexToAddress("0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"),
		},
		{
			Name:    "SushiSwap",
			Factory: common.HexToAddress("0xC0AEe478e3658e2610c5F7A4A2E1777ce9e4f2Ac"),
		},
		{
			Name:    "ShibaSwap",
			Factory: common.HexToAddress("0x115934131916C8b277DD010Ee02de363c09d037c"),
		},
	}

	pairs := [][2]Token{
		{USDC, WETH},
		{WBTC, WETH},
	}

	for _, dex := range dexes {
		f, err := factory.NewFactory(dex.Factory, client)
		if err != nil {
			log.Fatalf("Factory init error: %v", err)
		}

		for _, tokens := range pairs {
			printPairInfo(client, f, dex.Name, tokens[0], tokens[1])
		}
	}
}

func printPairInfo(
	client *ethclient.Client,
	f *factory.Factory,
	dexName string,
	tokenA, tokenB Token,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// pairAddr, err := f.GetPair(&bindCallOpts(ctx), tokenA, tokenB)
	opts := bind.CallOpts{
		Context: ctx,
	}

	pairAddr, err := f.GetPair(&opts, tokenA.Address, tokenB.Address)

	if err != nil || pairAddr == (common.Address{}) {
		fmt.Printf("%s | %s/%s | Pair not found\n",
			dexName, tokenA.Address.Hex(), tokenB.Address.Hex())
		return
	}

	p, err := pair.NewPair(pairAddr, client)
	if err != nil {
		log.Printf("Pair init error: %v\n", err)
		return
	}

	reserves, err := p.GetReserves(&opts)
	if err != nil {
		log.Printf("Reserve read error: %v\n", err)
		return
	}

	fmt.Printf(
		"%s | %s/%s | reserve0: %s | reserve1: %s\n",
		dexName,
		//tokenA.Address.Hex(),
		//tokenB.Address.Hex(),
		tokenA.Symbol,
		tokenB.Symbol,
		reserves.Reserve0.String(),
		reserves.Reserve1.String(),
	)
}
