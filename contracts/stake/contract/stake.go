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

// StakeInterfaceABI is the input ABI used to generate the binding from.
const StakeInterfaceABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"hasStake\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// StakeInterfaceBin is the compiled bytecode used for deploying new contracts.
const StakeInterfaceBin = `0x`

// DeployStakeInterface deploys a new Ethereum contract, binding an instance of StakeInterface to it.
func DeployStakeInterface(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *StakeInterface, error) {
	parsed, err := abi.JSON(strings.NewReader(StakeInterfaceABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(StakeInterfaceBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StakeInterface{StakeInterfaceCaller: StakeInterfaceCaller{contract: contract}, StakeInterfaceTransactor: StakeInterfaceTransactor{contract: contract}, StakeInterfaceFilterer: StakeInterfaceFilterer{contract: contract}}, nil
}

// StakeInterface is an auto generated Go binding around an Ethereum contract.
type StakeInterface struct {
	StakeInterfaceCaller     // Read-only binding to the contract
	StakeInterfaceTransactor // Write-only binding to the contract
	StakeInterfaceFilterer   // Log filterer for contract events
}

// StakeInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakeInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakeInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakeInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakeInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakeInterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakeInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakeInterfaceSession struct {
	Contract     *StakeInterface   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakeInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakeInterfaceCallerSession struct {
	Contract *StakeInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// StakeInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakeInterfaceTransactorSession struct {
	Contract     *StakeInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// StakeInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakeInterfaceRaw struct {
	Contract *StakeInterface // Generic contract binding to access the raw methods on
}

// StakeInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakeInterfaceCallerRaw struct {
	Contract *StakeInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// StakeInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakeInterfaceTransactorRaw struct {
	Contract *StakeInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStakeInterface creates a new instance of StakeInterface, bound to a specific deployed contract.
func NewStakeInterface(address common.Address, backend bind.ContractBackend) (*StakeInterface, error) {
	contract, err := bindStakeInterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StakeInterface{StakeInterfaceCaller: StakeInterfaceCaller{contract: contract}, StakeInterfaceTransactor: StakeInterfaceTransactor{contract: contract}, StakeInterfaceFilterer: StakeInterfaceFilterer{contract: contract}}, nil
}

// NewStakeInterfaceCaller creates a new read-only instance of StakeInterface, bound to a specific deployed contract.
func NewStakeInterfaceCaller(address common.Address, caller bind.ContractCaller) (*StakeInterfaceCaller, error) {
	contract, err := bindStakeInterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakeInterfaceCaller{contract: contract}, nil
}

// NewStakeInterfaceTransactor creates a new write-only instance of StakeInterface, bound to a specific deployed contract.
func NewStakeInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*StakeInterfaceTransactor, error) {
	contract, err := bindStakeInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakeInterfaceTransactor{contract: contract}, nil
}

// NewStakeInterfaceFilterer creates a new log filterer instance of StakeInterface, bound to a specific deployed contract.
func NewStakeInterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*StakeInterfaceFilterer, error) {
	contract, err := bindStakeInterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakeInterfaceFilterer{contract: contract}, nil
}

// bindStakeInterface binds a generic wrapper to an already deployed contract.
func bindStakeInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StakeInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakeInterface *StakeInterfaceRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _StakeInterface.Contract.StakeInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakeInterface *StakeInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakeInterface.Contract.StakeInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakeInterface *StakeInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakeInterface.Contract.StakeInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakeInterface *StakeInterfaceCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _StakeInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakeInterface *StakeInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakeInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakeInterface *StakeInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakeInterface.Contract.contract.Transact(opts, method, params...)
}

// HasStake is a free data retrieval call binding the contract method 0xe73e14bf.
//
// Solidity: function hasStake(_address address) constant returns(bool)
func (_StakeInterface *StakeInterfaceCaller) HasStake(opts *bind.CallOpts, _address common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _StakeInterface.contract.Call(opts, out, "hasStake", _address)
	return *ret0, err
}

// HasStake is a free data retrieval call binding the contract method 0xe73e14bf.
//
// Solidity: function hasStake(_address address) constant returns(bool)
func (_StakeInterface *StakeInterfaceSession) HasStake(_address common.Address) (bool, error) {
	return _StakeInterface.Contract.HasStake(&_StakeInterface.CallOpts, _address)
}

// HasStake is a free data retrieval call binding the contract method 0xe73e14bf.
//
// Solidity: function hasStake(_address address) constant returns(bool)
func (_StakeInterface *StakeInterfaceCallerSession) HasStake(_address common.Address) (bool, error) {
	return _StakeInterface.Contract.HasStake(&_StakeInterface.CallOpts, _address)
}
