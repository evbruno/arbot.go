// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package pair

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// PairMetaData contains all meta data concerning the Pair contract.
var PairMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getReserves\",\"outputs\":[{\"internalType\":\"uint112\",\"name\":\"reserve0\",\"type\":\"uint112\"},{\"internalType\":\"uint112\",\"name\":\"reserve1\",\"type\":\"uint112\"},{\"internalType\":\"uint32\",\"name\":\"blockTimestampLast\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// PairABI is the input ABI used to generate the binding from.
// Deprecated: Use PairMetaData.ABI instead.
var PairABI = PairMetaData.ABI

// Pair is an auto generated Go binding around an Ethereum contract.
type Pair struct {
	PairCaller     // Read-only binding to the contract
	PairTransactor // Write-only binding to the contract
	PairFilterer   // Log filterer for contract events
}

// PairCaller is an auto generated read-only Go binding around an Ethereum contract.
type PairCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PairTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PairTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PairFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PairFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PairSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PairSession struct {
	Contract     *Pair             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PairCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PairCallerSession struct {
	Contract *PairCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// PairTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PairTransactorSession struct {
	Contract     *PairTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PairRaw is an auto generated low-level Go binding around an Ethereum contract.
type PairRaw struct {
	Contract *Pair // Generic contract binding to access the raw methods on
}

// PairCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PairCallerRaw struct {
	Contract *PairCaller // Generic read-only contract binding to access the raw methods on
}

// PairTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PairTransactorRaw struct {
	Contract *PairTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPair creates a new instance of Pair, bound to a specific deployed contract.
func NewPair(address common.Address, backend bind.ContractBackend) (*Pair, error) {
	contract, err := bindPair(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Pair{PairCaller: PairCaller{contract: contract}, PairTransactor: PairTransactor{contract: contract}, PairFilterer: PairFilterer{contract: contract}}, nil
}

// NewPairCaller creates a new read-only instance of Pair, bound to a specific deployed contract.
func NewPairCaller(address common.Address, caller bind.ContractCaller) (*PairCaller, error) {
	contract, err := bindPair(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PairCaller{contract: contract}, nil
}

// NewPairTransactor creates a new write-only instance of Pair, bound to a specific deployed contract.
func NewPairTransactor(address common.Address, transactor bind.ContractTransactor) (*PairTransactor, error) {
	contract, err := bindPair(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PairTransactor{contract: contract}, nil
}

// NewPairFilterer creates a new log filterer instance of Pair, bound to a specific deployed contract.
func NewPairFilterer(address common.Address, filterer bind.ContractFilterer) (*PairFilterer, error) {
	contract, err := bindPair(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PairFilterer{contract: contract}, nil
}

// bindPair binds a generic wrapper to an already deployed contract.
func bindPair(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PairMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Pair *PairRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Pair.Contract.PairCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Pair *PairRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Pair.Contract.PairTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Pair *PairRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Pair.Contract.PairTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Pair *PairCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Pair.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Pair *PairTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Pair.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Pair *PairTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Pair.Contract.contract.Transact(opts, method, params...)
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint112 reserve0, uint112 reserve1, uint32 blockTimestampLast)
func (_Pair *PairCaller) GetReserves(opts *bind.CallOpts) (struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
}, error) {
	var out []interface{}
	err := _Pair.contract.Call(opts, &out, "getReserves")

	outstruct := new(struct {
		Reserve0           *big.Int
		Reserve1           *big.Int
		BlockTimestampLast uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Reserve0 = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Reserve1 = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.BlockTimestampLast = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err

}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint112 reserve0, uint112 reserve1, uint32 blockTimestampLast)
func (_Pair *PairSession) GetReserves() (struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
}, error) {
	return _Pair.Contract.GetReserves(&_Pair.CallOpts)
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint112 reserve0, uint112 reserve1, uint32 blockTimestampLast)
func (_Pair *PairCallerSession) GetReserves() (struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
}, error) {
	return _Pair.Contract.GetReserves(&_Pair.CallOpts)
}
