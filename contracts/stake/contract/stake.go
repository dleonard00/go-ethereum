// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// HasStakeABI is the input ABI used to generate the binding from.
const HasStakeABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"hasStake\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// HasStakeBin is the compiled bytecode used for deploying new contracts.
const HasStakeBin = `0x`

// DeployHasStake deploys a new Ethereum contract, binding an instance of HasStake to it.
func DeployHasStake(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *HasStake, error) {
	parsed, err := abi.JSON(strings.NewReader(HasStakeABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(HasStakeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &HasStake{HasStakeCaller: HasStakeCaller{contract: contract}, HasStakeTransactor: HasStakeTransactor{contract: contract}, HasStakeFilterer: HasStakeFilterer{contract: contract}}, nil
}

// HasStake is an auto generated Go binding around an Ethereum contract.
type HasStake struct {
	HasStakeCaller     // Read-only binding to the contract
	HasStakeTransactor // Write-only binding to the contract
	HasStakeFilterer   // Log filterer for contract events
}

// HasStakeCaller is an auto generated read-only Go binding around an Ethereum contract.
type HasStakeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HasStakeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type HasStakeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HasStakeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type HasStakeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HasStakeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type HasStakeSession struct {
	Contract     *HasStake         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HasStakeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type HasStakeCallerSession struct {
	Contract *HasStakeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// HasStakeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type HasStakeTransactorSession struct {
	Contract     *HasStakeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// HasStakeRaw is an auto generated low-level Go binding around an Ethereum contract.
type HasStakeRaw struct {
	Contract *HasStake // Generic contract binding to access the raw methods on
}

// HasStakeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type HasStakeCallerRaw struct {
	Contract *HasStakeCaller // Generic read-only contract binding to access the raw methods on
}

// HasStakeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type HasStakeTransactorRaw struct {
	Contract *HasStakeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewHasStake creates a new instance of HasStake, bound to a specific deployed contract.
func NewHasStake(address common.Address, backend bind.ContractBackend) (*HasStake, error) {
	contract, err := bindHasStake(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &HasStake{HasStakeCaller: HasStakeCaller{contract: contract}, HasStakeTransactor: HasStakeTransactor{contract: contract}, HasStakeFilterer: HasStakeFilterer{contract: contract}}, nil
}

// NewHasStakeCaller creates a new read-only instance of HasStake, bound to a specific deployed contract.
func NewHasStakeCaller(address common.Address, caller bind.ContractCaller) (*HasStakeCaller, error) {
	contract, err := bindHasStake(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HasStakeCaller{contract: contract}, nil
}

// NewHasStakeTransactor creates a new write-only instance of HasStake, bound to a specific deployed contract.
func NewHasStakeTransactor(address common.Address, transactor bind.ContractTransactor) (*HasStakeTransactor, error) {
	contract, err := bindHasStake(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HasStakeTransactor{contract: contract}, nil
}

// NewHasStakeFilterer creates a new log filterer instance of HasStake, bound to a specific deployed contract.
func NewHasStakeFilterer(address common.Address, filterer bind.ContractFilterer) (*HasStakeFilterer, error) {
	contract, err := bindHasStake(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HasStakeFilterer{contract: contract}, nil
}

// bindHasStake binds a generic wrapper to an already deployed contract.
func bindHasStake(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(HasStakeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HasStake *HasStakeRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _HasStake.Contract.HasStakeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HasStake *HasStakeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HasStake.Contract.HasStakeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HasStake *HasStakeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HasStake.Contract.HasStakeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HasStake *HasStakeCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _HasStake.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HasStake *HasStakeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HasStake.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HasStake *HasStakeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HasStake.Contract.contract.Transact(opts, method, params...)
}

// HasStake is a free data retrieval call binding the contract method 0xe73e14bf.
//
// Solidity: function hasStake(_address address) constant returns(bool)
func (_HasStake *HasStakeCaller) HasStake(opts *bind.CallOpts, _address common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _HasStake.contract.Call(opts, out, "hasStake", _address)
	return *ret0, err
}

// HasStake is a free data retrieval call binding the contract method 0xe73e14bf.
//
// Solidity: function hasStake(_address address) constant returns(bool)
func (_HasStake *HasStakeSession) HasStake(_address common.Address) (bool, error) {
	return _HasStake.Contract.HasStake(&_HasStake.CallOpts, _address)
}

// HasStake is a free data retrieval call binding the contract method 0xe73e14bf.
//
// Solidity: function hasStake(_address address) constant returns(bool)
func (_HasStake *HasStakeCallerSession) HasStake(_address common.Address) (bool, error) {
	return _HasStake.Contract.HasStake(&_HasStake.CallOpts, _address)
}
