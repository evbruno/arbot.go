package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"time"

	"github.com/evbruno/arbot-go/factory"
	"github.com/evbruno/arbot-go/pair"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
		//log.Fatal("Set MAINNET_RPC env variable")
		rpc = "http://127.0.0.1:8545/"
		log.Println("MAINNET_RPC env variable not set, using local anvil...")
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

	getRPCInfo(rpc, client)

	for _, tokens := range pairs {

		prices := make(map[string]float64)

		for _, dex := range dexes {
			// f, err := factory.NewFactory(dex.Factory, client)

			// if err != nil {
			// 	log.Fatalf("Factory init error: %v", err)
			// }

			// printPairInfo(client, f, dex.Name, tokens[0], tokens[1])
			fmt.Printf("\n====== %s/%s ======\n",
				tokens[0].Symbol,
				tokens[1].Symbol)

			price, err := fetchSpotPrice(client, dex, tokens[0], tokens[1])
			if err != nil {
				log.Printf("Error fetching spot price: %v\n", err)
				continue
			}

			prices[dex.Name] = price

			fmt.Printf("%s | Price %s per %s: %.6f\n",
				dex.Name,
				tokens[0].Symbol,
				tokens[1].Symbol,
				price)
		}

		printSpread(prices)
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

func fetchSpotPrice(client *ethclient.Client, dex Dex, tokenA, tokenB Token) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := &bind.CallOpts{Context: ctx}

	f, err := factory.NewFactory(dex.Factory, client)
	if err != nil {
		return 0, err
	}

	pairAddr, err := f.GetPair(opts, tokenA.Address, tokenB.Address)
	if err != nil || pairAddr == (common.Address{}) {
		return 0, fmt.Errorf("pair not found")
	}

	p, err := pair.NewPair(pairAddr, client)
	if err != nil {
		return 0, err
	}

	// 1️⃣ read token0 and token1
	token0, err := p.Token0(opts)
	if err != nil {
		return 0, err
	}

	reserves, err := p.GetReserves(opts)
	if err != nil {
		return 0, err
	}

	// 2️⃣ map reserves properly
	var reserveA, reserveB *big.Int

	if token0 == tokenA.Address {
		reserveA = reserves.Reserve0
		reserveB = reserves.Reserve1
	} else {
		reserveA = reserves.Reserve1
		reserveB = reserves.Reserve0
	}

	// 3️⃣ normalize by decimals
	normA := normalize(reserveA, tokenA.Decimals)
	normB := normalize(reserveB, tokenB.Decimals)

	// 4️⃣ compute spot price
	price := normA / normB

	return price, nil
}

func normalize(amount *big.Int, decimals uint8) float64 {
	f := new(big.Float).SetInt(amount)
	denom := new(big.Float).SetFloat64(math.Pow10(int(decimals)))
	result, _ := new(big.Float).Quo(f, denom).Float64()
	return result
}

func printSpread(prices map[string]float64) {
	if len(prices) < 2 {
		return
	}

	var minDex, maxDex string
	minPrice := math.MaxFloat64
	maxPrice := 0.0

	for dex, price := range prices {
		if price < minPrice {
			minPrice = price
			minDex = dex
		}
		if price > maxPrice {
			maxPrice = price
			maxDex = dex
		}
	}

	spread := ((maxPrice - minPrice) / minPrice) * 100

	fmt.Printf("Spread: Buy on %s, Sell on %s | %.4f%%\n",
		minDex, maxDex, spread)
}

func getRPCInfo(rpc string, client *ethclient.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Println("💣💣💣")
		log.Fatal("Error getting header, is it running? `anvil --fork-url https://eth-mainnet.g.alchemy.com/v2/$ALCHEMY_API_KEY`\nErr: ", err)
	}

	printBlockInfo(rpc, header)
}

func printBlockInfo(network string, header *types.Header) {
	blockNumber := header.Number
	blockHash := header.Hash()
	timestamp := time.Unix(int64(header.Time), 0)
	baseFee := header.BaseFee

	fmt.Println("====================================")
	fmt.Printf("Network: %s\n", network)
	fmt.Printf("Block Number: %s\n", blockNumber.String())
	fmt.Printf("Block Hash: %s\n", blockHash.Hex())
	fmt.Printf("Timestamp: %s\n", timestamp.UTC().Format(time.RFC3339))

	if baseFee != nil {
		fmt.Printf("Base Fee (Wei): %s\n", baseFee.String())
		fmt.Printf("Base Fee (Gwei): %s\n", weiToGwei(baseFee))
	} else {
		fmt.Println("Base Fee: nil (pre-EIP1559 chain?)")
	}
	fmt.Println("====================================")
}

func weiToGwei(wei *big.Int) string {
	gwei := new(big.Float).Quo(
		new(big.Float).SetInt(wei),
		big.NewFloat(1e9),
	)
	return gwei.Text('f', 6)
}
