package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/evbruno/arbot-go/erc20"
	"github.com/evbruno/arbot-go/erc20bytes32"
	"github.com/evbruno/arbot-go/factory"
	"github.com/evbruno/arbot-go/internal/ethereum"
	"github.com/evbruno/arbot-go/pair"

	"gopkg.in/yaml.v3"
)

type PairInfo struct {
	Pair               string
	Price01            *big.Float
	PairAddress        common.Address
	Reserve0           *big.Int
	Reserve1           *big.Int
	Token0Symbol       string
	Token0Name         string
	Token0Decimals     uint8
	Token0             common.Address
	Token1             common.Address
	Token1Symbol       string
	Token1Name         string
	Token1Decimals     uint8
	BlockTimestampLast uint32
}

const (
	UniswapV2FactoryAddress = "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"
	Multicall2Address       = "0x5BA1e12693Dc8F9c48aAD8770482f4739bEeD696"
	SushiSwapFactoryAddress = "0xC0AEe478e3658e2610c5F7A4A2E1777ce9e4f2Ac"
)

var (
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
)

func main() {
	ctx := context.Background()

	uniswapAddress := common.HexToAddress(UniswapV2FactoryAddress)

	//$ anvil --fork-url https://eth-mainnet.g.alchemy.com/v2/your-api-key
	client, err := ethereum.Dial("http://127.0.0.1:8545")
	if err != nil {
		logger.Error("rpc dial failed", "error", err)
		os.Exit(1)
	}

	blockNumber, _ := client.BlockNumber(ctx)
	logger.Info("current block number", "block", blockNumber)

	factoryContract, err := factory.NewFactory(uniswapAddress, client)
	if err != nil {
		logger.Error("failed to create factory instance", "error", err)
		os.Exit(1)
	}

	opts := &bind.CallOpts{Context: ctx}
	allPairsLength, err := factoryContract.AllPairsLength(opts)
	if err != nil {
		logger.Error("failed to get all pairs length", "error", err)
		os.Exit(1)
	}
	logger.Info("all pairs length", "length", allPairsLength)

	uniswapPairs, err := collectUniswapV2Pairs(ctx, logger, client)
	if err != nil {
		logger.Error("failed to collect pairs", "error", err)
		os.Exit(1)
	}

	sushiPairs, err := collectSushiSwapPairs(ctx, logger, client)
	if err != nil {
		logger.Error("failed to collect pairs", "error", err)
		os.Exit(1)
	}

	fileName := "sushiPairs.csv"
	logger.Info("collected pairs", "count", len(sushiPairs), "writing to file", fileName)
	if len(sushiPairs) == 0 {
		logger.Warn("no pairs collected, skipping file write")
		return
	}

	if err := writePairsCSV(fileName, sushiPairs); err != nil {
		logger.Error("failed to write pairs csv", "error", err)
		os.Exit(1)
	}

	if err := writeUniswapYAML("uniswapv2.yaml", uniswapPairs); err != nil {
		logger.Error("failed to write pairs yaml", "error", err)
		os.Exit(1)
	}

	if err := writeSushiYAML("sushiswapv2.yaml", sushiPairs); err != nil {
		logger.Error("failed to write pairs yaml", "error", err)
		os.Exit(1)
	}

	logger.Info("done")
}

func collectUniswapV2Pairs(ctx context.Context, logger *slog.Logger, client bind.ContractBackend) ([]PairInfo, error) {
	pairs := []string{
		"0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc", // USDC/ETH
		"0x0d4a11d5EEaaC28EC3F61d100daF4d40471f1852", // ETH/USDT
		"0xBb2b8038a1640196FbE3e38816F3e67Cba72D940", // WBTC/ETH
		"0xA478c2975Ab1Ea89e8196811F51A7B7Ade33eB11", // DAI/ETH
		"0xA43fe16908251ee70EF74718545e4FE6C5cCEc9f", // PEPE/ETH
		"0x9C4Fe5FFD9A9fC5678cFBd93Aa2D4FD684b67C4C", // PAXG/ETH
		"0x52c77b0CB827aFbAD022E6d6CAF2C44452eDbc39", // ETH/SPX
		"0xA478c2975Ab1Ea89e8196811F51A7B7Ade33eB11", // DAI/ETH
		"0x69C7bd26512f52bF6F76faB834140D13Dda673Ca", // NPC/ETH
		"0xc2eaB7d33d3cB97692eCB231A5D0e4A649Cb539d", // Mog/ETH
		"0x43DE4318b6EB91a7cF37975dBB574396A7b5B5c6", // Banana/ETH
		"0x859f7092f56c43BB48bb46dE7119d9c799716CDF", // XCN/ETH
		"0xFFf8D5fFF6Ee3226fa2F5d7D5D8C3Ff785be9C74", // KEKIUS/ETH
		"0x25647E01Bd0967C1B9599FA3521939871D1d0888", // ETH/Super
		"0x180EFC1349A69390aDE25667487a826164C9c6E4", // Bean/ETH
		"0x2a6c340bCbb0a79D3deecD3bc5cBc2605ea9259f", // $PAAL/ETH
		"0xC555D55279023E732CcD32D812114cAF5838fD46", // Neiro/ETH
		"0x470e8de2eBaef52014A47Cb5E6aF86884947F08c",
		"0xca7c2771D248dCBe09EABE0CE57A62e18dA178c0",
		"0xC730EF0f4973DA9cC0aB8Ab291890D3e77f58F79",
	}

	return collectPairsFromList(ctx, logger, client, pairs)
}

func collectSushiSwapPairs(ctx context.Context, logger *slog.Logger, client bind.ContractBackend) ([]PairInfo, error) {
	pairs := []string{
		"0x06da0fd433c1a5d7a4faa01111c044910a184553", // ETH/USDT
		"0xc3f279090a47e80990fe3a9c30d24cb117ef91a8",
		"0x001b6450083e531a5a7bf310bd2c1af4247e23d4",
		"0x58dc5a51fe44589beb22e8ce67720b5bc5378009",
		"0x397ff1542f962076d0bfe58ea045ffa2d347aca0",
		"0x6a091a3406e0073c3cd6340122143009adac0eda",
		"0xbdc120fef90fb185a49ad8fa62c7bc0ed0516cc7",
		"0x62ccb80f72cc5c975c5bc7fb4433d3c336ce5ceb",
		"0xc3d03e4f041fd4cd388c549ee2a29a9e5075882f",
		"0x15e86e6f65ef7ea1dbb72a5e51a07926fb1c82e3",
		"0x795065dcc9f64b5614c407a6efdc400da6221fb0",
		"0x99b42f2b49c395d2a77d973f6009abb5d67da343",
		"0x4a86c01d67965f8cb3d0aaa2c655705e64097c31",
		"0xe4ebd836832f1a8a81641111a5b081a2f90b9430",
		"0x05767d9ef41dc40689678ffca0608878fb3de906",
		"0x611cde65dea90918c0078ac0400a72b0d25b9bb1",
		"0xec8c342bc3e07f05b9a782bc34e7f04fb9b44502",
		"0x0c365789dbbb94a29f8720dc465554c587e897db",
		"0x7825de5586e4d2fd04459091bbe783fa243e1bf3",
		"0xceff51756c56ceffca006cd410b03ffc46dd3a58",
	}

	return collectPairsFromList(ctx, logger, client, pairs)
}

func collectPairsFromList(ctx context.Context, logger *slog.Logger, client bind.ContractBackend, pairAddresses []string) ([]PairInfo, error) {
	multicall, err := newMulticall(client)
	if err != nil {
		return nil, err
	}

	results := make([]PairInfo, 0, len(pairAddresses))
	for _, pairAddr := range pairAddresses {
		address := common.HexToAddress(pairAddr)
		info, err := fetchPairInfoMulticall(ctx, client, multicall, address)
		if err != nil {
			logger.Error("failed to collect pair info", "error", err, "pair", pairAddr)
			continue
		}
		results = append(results, *info)
	}

	return results, nil
}

func newMulticall(client bind.ContractBackend) (*ethereum.Multicall, error) {
	ethClient, ok := client.(*ethclient.Client)
	if !ok {
		return nil, fmt.Errorf("multicall requires *ethclient.Client")
	}

	multicall := ethereum.NewMulticall(ethClient, common.HexToAddress(Multicall2Address))
	if multicall == nil || multicall.IsZero() {
		return nil, fmt.Errorf("multicall not configured")
	}

	return multicall, nil
}

func collectPairs(ctx context.Context, logger *slog.Logger, client bind.ContractBackend, factoryContract *factory.Factory, limit int) ([]PairInfo, error) {
	if limit <= 0 {
		return nil, nil
	}

	if limit > 0 {
		logger.Info("collecting pairs", "limit", limit)
	}

	results := make([]PairInfo, 0, limit)
	callOpts := &bind.CallOpts{Context: ctx}
	for i := range limit {
		pairAddress, err := factoryContract.AllPairs(callOpts, big.NewInt(int64(i)))
		if err != nil {
			logger.Error("failed to get pair address", "error", err, "index", i)
			continue
		}

		info, err := fetchPairInfo(callOpts, client, pairAddress)
		if err != nil {
			logger.Error("failed to collect pair info", "error", err, "pair", pairAddress.Hex())
			continue
		}

		if i%25 == 0 {
			logger.Info("collected pair", "index", i)
			time.Sleep(1 * time.Second) // throttle to avoid overwhelming the node
		}

		results = append(results, *info)
	}

	return results, nil
}

func fetchPairInfo(opts *bind.CallOpts, client bind.ContractBackend, pairAddress common.Address) (*PairInfo, error) {
	p, err := pair.NewPair(pairAddress, client)
	if err != nil {
		return nil, fmt.Errorf("create pair instance: %w", err)
	}

	token0, err := p.Token0(opts)
	if err != nil {
		return nil, fmt.Errorf("get token0: %w", err)
	}

	token1, err := p.Token1(opts)
	if err != nil {
		return nil, fmt.Errorf("get token1: %w", err)
	}

	token0Contract, err := erc20.NewErc20(token0, client)
	if err != nil {
		return nil, fmt.Errorf("create token0 instance: %w", err)
	}

	token1Contract, err := erc20.NewErc20(token1, client)
	if err != nil {
		return nil, fmt.Errorf("create token1 instance: %w", err)
	}

	token0Symbol, err := readTokenSymbol(opts, client, token0)
	if err != nil {
		return nil, fmt.Errorf("get token0 symbol: %w", err)
	}

	token0Name, err := readTokenName(opts, client, token0)
	if err != nil {
		return nil, fmt.Errorf("get token0 name: %w", err)
	}

	token0Decimals, err := token0Contract.Decimals(opts)
	if err != nil {
		return nil, fmt.Errorf("get token0 decimals: %w", err)
	}

	token1Symbol, err := readTokenSymbol(opts, client, token1)
	if err != nil {
		return nil, fmt.Errorf("get token1 symbol: %w", err)
	}

	token1Name, err := readTokenName(opts, client, token1)
	if err != nil {
		return nil, fmt.Errorf("get token1 name: %w", err)
	}

	token1Decimals, err := token1Contract.Decimals(opts)
	if err != nil {
		return nil, fmt.Errorf("get token1 decimals: %w", err)
	}

	reserves, err := p.GetReserves(opts)
	if err != nil {
		return nil, fmt.Errorf("get reserves: %w", err)
	}

	pi := &PairInfo{
		Pair:               fmt.Sprintf("%s/%s", token0Symbol, token1Symbol),
		Price01:            nil, // will be calculated later
		PairAddress:        pairAddress,
		Reserve0:           reserves.Reserve0,
		Reserve1:           reserves.Reserve1,
		Token0:             token0,
		Token0Symbol:       token0Symbol,
		Token0Name:         token0Name,
		Token0Decimals:     token0Decimals,
		Token1:             token1,
		Token1Symbol:       token1Symbol,
		Token1Name:         token1Name,
		Token1Decimals:     token1Decimals,
		BlockTimestampLast: reserves.BlockTimestampLast,
	}
	pi.Price01 = priceToken0InToken1(pi)
	return pi, nil
}

func fetchPairInfoMulticall(ctx context.Context, client bind.ContractBackend, multicall *ethereum.Multicall, pairAddress common.Address) (*PairInfo, error) {
	callOpts := &bind.CallOpts{Context: ctx}

	token0, token1, reserves, err := fetchPairBasicsMulticall(ctx, multicall, pairAddress)
	if err != nil {
		return nil, err
	}

	token0Contract, err := erc20.NewErc20(token0, client)
	if err != nil {
		return nil, fmt.Errorf("create token0 instance: %w", err)
	}

	token1Contract, err := erc20.NewErc20(token1, client)
	if err != nil {
		return nil, fmt.Errorf("create token1 instance: %w", err)
	}

	token0Symbol, err := readTokenSymbol(callOpts, client, token0)
	if err != nil {
		return nil, fmt.Errorf("get token0 symbol: %w", err)
	}

	token0Name, err := readTokenName(callOpts, client, token0)
	if err != nil {
		return nil, fmt.Errorf("get token0 name: %w", err)
	}

	token0Decimals, err := token0Contract.Decimals(callOpts)
	if err != nil {
		return nil, fmt.Errorf("get token0 decimals: %w", err)
	}

	token1Symbol, err := readTokenSymbol(callOpts, client, token1)
	if err != nil {
		return nil, fmt.Errorf("get token1 symbol: %w", err)
	}

	token1Name, err := readTokenName(callOpts, client, token1)
	if err != nil {
		return nil, fmt.Errorf("get token1 name: %w", err)
	}

	token1Decimals, err := token1Contract.Decimals(callOpts)
	if err != nil {
		return nil, fmt.Errorf("get token1 decimals: %w", err)
	}

	pi := &PairInfo{
		Pair:               fmt.Sprintf("%s/%s", token0Symbol, token1Symbol),
		Price01:            nil,
		PairAddress:        pairAddress,
		Reserve0:           reserves.Reserve0,
		Reserve1:           reserves.Reserve1,
		Token0:             token0,
		Token0Symbol:       token0Symbol,
		Token0Name:         token0Name,
		Token0Decimals:     token0Decimals,
		Token1:             token1,
		Token1Symbol:       token1Symbol,
		Token1Name:         token1Name,
		Token1Decimals:     token1Decimals,
		BlockTimestampLast: reserves.BlockTimestampLast,
	}
	pi.Price01 = priceToken0InToken1(pi)
	return pi, nil
}

func fetchPairBasicsMulticall(ctx context.Context, multicall *ethereum.Multicall, pairAddress common.Address) (common.Address, common.Address, struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
}, error) {

	var emptyReserves struct {
		Reserve0           *big.Int
		Reserve1           *big.Int
		BlockTimestampLast uint32
	}

	if multicall == nil || multicall.IsZero() {
		return common.Address{}, common.Address{}, emptyReserves, fmt.Errorf("multicall not configured")
	}

	pairABI, err := pair.PairMetaData.GetAbi()
	if err != nil {
		return common.Address{}, common.Address{}, emptyReserves, fmt.Errorf("load pair abi: %w", err)
	}

	token0Call, err := pairABI.Pack("token0")
	if err != nil {
		return common.Address{}, common.Address{}, emptyReserves, fmt.Errorf("pack token0 call: %w", err)
	}

	token1Call, err := pairABI.Pack("token1")
	if err != nil {
		return common.Address{}, common.Address{}, emptyReserves, fmt.Errorf("pack token1 call: %w", err)
	}

	reservesCall, err := pairABI.Pack("getReserves")
	if err != nil {
		return common.Address{}, common.Address{}, emptyReserves, fmt.Errorf("pack getReserves call: %w", err)
	}

	_, returnData, err := multicall.Aggregate(ctx, []ethereum.MulticallCall{
		{Target: pairAddress, CallData: token0Call},
		{Target: pairAddress, CallData: token1Call},
		{Target: pairAddress, CallData: reservesCall},
	})
	if err != nil {
		return common.Address{}, common.Address{}, emptyReserves, fmt.Errorf("multicall aggregate: %w", err)
	}

	if len(returnData) != 3 {
		return common.Address{}, common.Address{}, emptyReserves, fmt.Errorf("unexpected multicall return count: %d", len(returnData))
	}

	token0Out, err := pairABI.Unpack("token0", returnData[0])
	if err != nil {
		return common.Address{}, common.Address{}, emptyReserves, fmt.Errorf("unpack token0: %w", err)
	}
	token1Out, err := pairABI.Unpack("token1", returnData[1])
	if err != nil {
		return common.Address{}, common.Address{}, emptyReserves, fmt.Errorf("unpack token1: %w", err)
	}
	reservesOut, err := pairABI.Unpack("getReserves", returnData[2])
	if err != nil {
		return common.Address{}, common.Address{}, emptyReserves, fmt.Errorf("unpack getReserves: %w", err)
	}

	if len(token0Out) != 1 || len(token1Out) != 1 || len(reservesOut) != 3 {
		return common.Address{}, common.Address{}, emptyReserves, fmt.Errorf("unexpected multicall output lengths")
	}

	token0, ok := token0Out[0].(common.Address)
	if !ok {
		return common.Address{}, common.Address{}, emptyReserves, fmt.Errorf("token0 has unexpected type")
	}
	token1, ok := token1Out[0].(common.Address)
	if !ok {
		return common.Address{}, common.Address{}, emptyReserves, fmt.Errorf("token1 has unexpected type")
	}

	reserve0, ok := reservesOut[0].(*big.Int)
	if !ok {
		return common.Address{}, common.Address{}, emptyReserves, fmt.Errorf("reserve0 has unexpected type")
	}
	reserve1, ok := reservesOut[1].(*big.Int)
	if !ok {
		return common.Address{}, common.Address{}, emptyReserves, fmt.Errorf("reserve1 has unexpected type")
	}
	blockTimestampLast, ok := reservesOut[2].(uint32)
	if !ok {
		return common.Address{}, common.Address{}, emptyReserves, fmt.Errorf("blockTimestampLast has unexpected type")
	}

	return token0, token1, struct {
		Reserve0           *big.Int
		Reserve1           *big.Int
		BlockTimestampLast uint32
	}{
		Reserve0:           reserve0,
		Reserve1:           reserve1,
		BlockTimestampLast: blockTimestampLast,
	}, nil
}

func readTokenSymbol(opts *bind.CallOpts, backend bind.ContractBackend, token common.Address) (string, error) {
	if symbol, err := readTokenString(opts, backend, token, "symbol"); err == nil {
		return symbol, nil
	}

	return readTokenBytes32String(opts, backend, token, "symbol")
}

func readTokenName(opts *bind.CallOpts, backend bind.ContractBackend, token common.Address) (string, error) {
	if name, err := readTokenString(opts, backend, token, "name"); err == nil {
		return name, nil
	}

	return readTokenBytes32String(opts, backend, token, "name")
}

func readTokenString(opts *bind.CallOpts, backend bind.ContractBackend, token common.Address, method string) (string, error) {
	contract, err := erc20.NewErc20(token, backend)
	if err != nil {
		return "", err
	}

	switch method {
	case "symbol":
		return contract.Symbol(opts)
	case "name":
		return contract.Name(opts)
	default:
		return "", fmt.Errorf("unsupported token method: %s", method)
	}
}

func readTokenBytes32String(opts *bind.CallOpts, backend bind.ContractBackend, token common.Address, method string) (string, error) {
	contract, err := erc20bytes32.NewErc20Bytes32(token, backend)
	if err != nil {
		return "", err
	}

	switch method {
	case "symbol":
		value, err := contract.Symbol(opts)
		if err != nil {
			return "", err
		}
		return bytes32ToString(value), nil
	case "name":
		value, err := contract.Name(opts)
		if err != nil {
			return "", err
		}
		return bytes32ToString(value), nil
	default:
		return "", fmt.Errorf("unsupported token method: %s", method)
	}
}

func bytes32ToString(value [32]byte) string {
	return string(bytes.TrimRight(value[:], "\x00"))
}

func writePairsCSV(path string, pairs []PairInfo) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	writer := csv.NewWriter(file)
	if err := writer.Write([]string{
		"pair",
		"price01",
		"pair_address",
		"reserve0",
		"reserve1",
		"token0",
		"token0_symbol",
		"token0_name",
		"token0_decimals",
		"token1",
		"token1_symbol",
		"token1_name",
		"token1_decimals",
		"block_timestamp_last",
	}); err != nil {
		return err
	}

	for _, info := range pairs {
		price := info.Price01
		if price == nil {
			price = priceToken0InToken1(&info)
		}

		record := []string{
			info.Pair,
			price.String(),
			info.PairAddress.Hex(),
			info.Reserve0.String(),
			info.Reserve1.String(),
			info.Token0.Hex(),
			info.Token0Symbol,
			info.Token0Name,
			strconv.FormatUint(uint64(info.Token0Decimals), 10),
			info.Token1.Hex(),
			info.Token1Symbol,
			info.Token1Name,
			strconv.FormatUint(uint64(info.Token1Decimals), 10),
			strconv.FormatUint(uint64(info.BlockTimestampLast), 10),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	writer.Flush()
	return writer.Error()
}

func priceToken0InToken1(info *PairInfo) *big.Float {
	// price of token0 in terms of token1, adjusted for decimals
	precision := uint(256)
	reserve0 := new(big.Float).SetPrec(precision).SetInt(info.Reserve0)
	reserve1 := new(big.Float).SetPrec(precision).SetInt(info.Reserve1)

	adjustedReserve0 := new(big.Float).SetPrec(precision).Quo(reserve0, pow10Float(info.Token0Decimals))
	adjustedReserve1 := new(big.Float).SetPrec(precision).Quo(reserve1, pow10Float(info.Token1Decimals))

	if adjustedReserve0.Sign() == 0 {
		return new(big.Float).SetPrec(precision)
	}

	return new(big.Float).SetPrec(precision).Quo(adjustedReserve1, adjustedReserve0)
}

func pow10Float(decimals uint8) *big.Float {
	pow := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	return new(big.Float).SetPrec(256).SetInt(pow)
}

type YAMLPair struct {
	Pair    string `yaml:"pair"`
	Address string `yaml:"address"`
	Token0  string `yaml:"token0"`
	Token1  string `yaml:"token1"`
}

type YAMLToken struct {
	Symbol   string `yaml:"symbol"`
	Name     string `yaml:"name"`
	Address  string `yaml:"address"`
	Decimals uint8  `yaml:"decimals"`
}

type YamlDex struct {
	Name             string      `yaml:"name"`
	FactoryAddress   string      `yaml:"factory_address"`
	MulticallAddress string      `yaml:"multicall_address"`
	Pairs            []YAMLPair  `yaml:"pairs"`
	Tokens           []YAMLToken `yaml:"tokens"`
}

func writeUniswapYAML(path string, pairs []PairInfo) error {
	return writeDexYAML(path, "Uniswap V2", UniswapV2FactoryAddress, Multicall2Address, pairs)
}

func writeSushiYAML(path string, pairs []PairInfo) error {
	return writeDexYAML(path, "SushiSwap", SushiSwapFactoryAddress, Multicall2Address, pairs)
}

func writeDexYAML(path, name, factoryAddress, multicallAddress string, pairs []PairInfo) error {
	dex := YamlDex{
		Name:             name,
		FactoryAddress:   factoryAddress,
		MulticallAddress: multicallAddress,
		Pairs:            buildPairsYAML(pairs),
		Tokens:           buildTokensYAML(pairs),
	}

	out, err := yaml.Marshal(&dex)
	if err != nil {
		return err
	}

	return os.WriteFile(path, out, 0o644)
}

func buildPairsYAML(pairs []PairInfo) []YAMLPair {
	pairsYAML := make([]YAMLPair, len(pairs))
	for i, info := range pairs {
		pairsYAML[i] = YAMLPair{
			Pair:    info.Pair,
			Token0:  info.Token0.Hex(),
			Token1:  info.Token1.Hex(),
			Address: info.PairAddress.Hex(),
		}
	}

	slices.SortFunc(pairsYAML, func(a, b YAMLPair) int {
		return strings.Compare(a.Pair, b.Pair)
	})

	return pairsYAML
}

func buildTokensYAML(pairs []PairInfo) []YAMLToken {
	tokensMap := make(map[string]YAMLToken)
	for _, info := range pairs {
		if _, exists := tokensMap[info.Token0.Hex()]; !exists {
			tokensMap[info.Token0.Hex()] = YAMLToken{
				Symbol:   info.Token0Symbol,
				Name:     info.Token0Name,
				Address:  info.Token0.Hex(),
				Decimals: info.Token0Decimals,
			}
		}
		if _, exists := tokensMap[info.Token1.Hex()]; !exists {
			tokensMap[info.Token1.Hex()] = YAMLToken{
				Symbol:   info.Token1Symbol,
				Name:     info.Token1Name,
				Address:  info.Token1.Hex(),
				Decimals: info.Token1Decimals,
			}
		}
	}

	tokensYAML := make([]YAMLToken, 0, len(tokensMap))
	for _, token := range tokensMap {
		tokensYAML = append(tokensYAML, token)
	}

	slices.SortFunc(tokensYAML, func(a, b YAMLToken) int {
		return strings.Compare(a.Symbol, b.Symbol)
	})

	return tokensYAML
}
