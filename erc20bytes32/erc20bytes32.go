// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package erc20bytes32

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

// Erc20Bytes32MetaData contains all meta data concerning the Erc20Bytes32 contract.
var Erc20Bytes32MetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// Erc20Bytes32ABI is the input ABI used to generate the binding from.
// Deprecated: Use Erc20Bytes32MetaData.ABI instead.
var Erc20Bytes32ABI = Erc20Bytes32MetaData.ABI

// Erc20Bytes32 is an auto generated Go binding around an Ethereum contract.
type Erc20Bytes32 struct {
	Erc20Bytes32Caller     // Read-only binding to the contract
	Erc20Bytes32Transactor // Write-only binding to the contract
	Erc20Bytes32Filterer   // Log filterer for contract events
}

// Erc20Bytes32Caller is an auto generated read-only Go binding around an Ethereum contract.
type Erc20Bytes32Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc20Bytes32Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Erc20Bytes32Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc20Bytes32Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Erc20Bytes32Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc20Bytes32Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Erc20Bytes32Session struct {
	Contract     *Erc20Bytes32     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Erc20Bytes32CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Erc20Bytes32CallerSession struct {
	Contract *Erc20Bytes32Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// Erc20Bytes32TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Erc20Bytes32TransactorSession struct {
	Contract     *Erc20Bytes32Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// Erc20Bytes32Raw is an auto generated low-level Go binding around an Ethereum contract.
type Erc20Bytes32Raw struct {
	Contract *Erc20Bytes32 // Generic contract binding to access the raw methods on
}

// Erc20Bytes32CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Erc20Bytes32CallerRaw struct {
	Contract *Erc20Bytes32Caller // Generic read-only contract binding to access the raw methods on
}

// Erc20Bytes32TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Erc20Bytes32TransactorRaw struct {
	Contract *Erc20Bytes32Transactor // Generic write-only contract binding to access the raw methods on
}

// NewErc20Bytes32 creates a new instance of Erc20Bytes32, bound to a specific deployed contract.
func NewErc20Bytes32(address common.Address, backend bind.ContractBackend) (*Erc20Bytes32, error) {
	contract, err := bindErc20Bytes32(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Erc20Bytes32{Erc20Bytes32Caller: Erc20Bytes32Caller{contract: contract}, Erc20Bytes32Transactor: Erc20Bytes32Transactor{contract: contract}, Erc20Bytes32Filterer: Erc20Bytes32Filterer{contract: contract}}, nil
}

// NewErc20Bytes32Caller creates a new read-only instance of Erc20Bytes32, bound to a specific deployed contract.
func NewErc20Bytes32Caller(address common.Address, caller bind.ContractCaller) (*Erc20Bytes32Caller, error) {
	contract, err := bindErc20Bytes32(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Erc20Bytes32Caller{contract: contract}, nil
}

// NewErc20Bytes32Transactor creates a new write-only instance of Erc20Bytes32, bound to a specific deployed contract.
func NewErc20Bytes32Transactor(address common.Address, transactor bind.ContractTransactor) (*Erc20Bytes32Transactor, error) {
	contract, err := bindErc20Bytes32(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Erc20Bytes32Transactor{contract: contract}, nil
}

// NewErc20Bytes32Filterer creates a new log filterer instance of Erc20Bytes32, bound to a specific deployed contract.
func NewErc20Bytes32Filterer(address common.Address, filterer bind.ContractFilterer) (*Erc20Bytes32Filterer, error) {
	contract, err := bindErc20Bytes32(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Erc20Bytes32Filterer{contract: contract}, nil
}

// bindErc20Bytes32 binds a generic wrapper to an already deployed contract.
func bindErc20Bytes32(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := Erc20Bytes32MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Erc20Bytes32 *Erc20Bytes32Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Erc20Bytes32.Contract.Erc20Bytes32Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Erc20Bytes32 *Erc20Bytes32Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc20Bytes32.Contract.Erc20Bytes32Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Erc20Bytes32 *Erc20Bytes32Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Erc20Bytes32.Contract.Erc20Bytes32Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Erc20Bytes32 *Erc20Bytes32CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Erc20Bytes32.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Erc20Bytes32 *Erc20Bytes32TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc20Bytes32.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Erc20Bytes32 *Erc20Bytes32TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Erc20Bytes32.Contract.contract.Transact(opts, method, params...)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(bytes32)
func (_Erc20Bytes32 *Erc20Bytes32Caller) Name(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Erc20Bytes32.contract.Call(opts, &out, "name")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(bytes32)
func (_Erc20Bytes32 *Erc20Bytes32Session) Name() ([32]byte, error) {
	return _Erc20Bytes32.Contract.Name(&_Erc20Bytes32.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(bytes32)
func (_Erc20Bytes32 *Erc20Bytes32CallerSession) Name() ([32]byte, error) {
	return _Erc20Bytes32.Contract.Name(&_Erc20Bytes32.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(bytes32)
func (_Erc20Bytes32 *Erc20Bytes32Caller) Symbol(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Erc20Bytes32.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(bytes32)
func (_Erc20Bytes32 *Erc20Bytes32Session) Symbol() ([32]byte, error) {
	return _Erc20Bytes32.Contract.Symbol(&_Erc20Bytes32.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(bytes32)
func (_Erc20Bytes32 *Erc20Bytes32CallerSession) Symbol() ([32]byte, error) {
	return _Erc20Bytes32.Contract.Symbol(&_Erc20Bytes32.CallOpts)
}
