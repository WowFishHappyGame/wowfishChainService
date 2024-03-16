// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package wowfish

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

// WowFishBankMetaData contains all meta data concerning the WowFishBank contract.
var WowFishBankMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"_dayWithdrawAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_totalAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// WowFishBankABI is the input ABI used to generate the binding from.
// Deprecated: Use WowFishBankMetaData.ABI instead.
var WowFishBankABI = WowFishBankMetaData.ABI

// WowFishBank is an auto generated Go binding around an Ethereum contract.
type WowFishBank struct {
	WowFishBankCaller     // Read-only binding to the contract
	WowFishBankTransactor // Write-only binding to the contract
	WowFishBankFilterer   // Log filterer for contract events
}

// WowFishBankCaller is an auto generated read-only Go binding around an Ethereum contract.
type WowFishBankCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WowFishBankTransactor is an auto generated write-only Go binding around an Ethereum contract.
type WowFishBankTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WowFishBankFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WowFishBankFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WowFishBankSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WowFishBankSession struct {
	Contract     *WowFishBank      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WowFishBankCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WowFishBankCallerSession struct {
	Contract *WowFishBankCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// WowFishBankTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WowFishBankTransactorSession struct {
	Contract     *WowFishBankTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// WowFishBankRaw is an auto generated low-level Go binding around an Ethereum contract.
type WowFishBankRaw struct {
	Contract *WowFishBank // Generic contract binding to access the raw methods on
}

// WowFishBankCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WowFishBankCallerRaw struct {
	Contract *WowFishBankCaller // Generic read-only contract binding to access the raw methods on
}

// WowFishBankTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WowFishBankTransactorRaw struct {
	Contract *WowFishBankTransactor // Generic write-only contract binding to access the raw methods on
}

// NewWowFishBank creates a new instance of WowFishBank, bound to a specific deployed contract.
func NewWowFishBank(address common.Address, backend bind.ContractBackend) (*WowFishBank, error) {
	contract, err := bindWowFishBank(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WowFishBank{WowFishBankCaller: WowFishBankCaller{contract: contract}, WowFishBankTransactor: WowFishBankTransactor{contract: contract}, WowFishBankFilterer: WowFishBankFilterer{contract: contract}}, nil
}

// NewWowFishBankCaller creates a new read-only instance of WowFishBank, bound to a specific deployed contract.
func NewWowFishBankCaller(address common.Address, caller bind.ContractCaller) (*WowFishBankCaller, error) {
	contract, err := bindWowFishBank(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WowFishBankCaller{contract: contract}, nil
}

// NewWowFishBankTransactor creates a new write-only instance of WowFishBank, bound to a specific deployed contract.
func NewWowFishBankTransactor(address common.Address, transactor bind.ContractTransactor) (*WowFishBankTransactor, error) {
	contract, err := bindWowFishBank(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WowFishBankTransactor{contract: contract}, nil
}

// NewWowFishBankFilterer creates a new log filterer instance of WowFishBank, bound to a specific deployed contract.
func NewWowFishBankFilterer(address common.Address, filterer bind.ContractFilterer) (*WowFishBankFilterer, error) {
	contract, err := bindWowFishBank(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WowFishBankFilterer{contract: contract}, nil
}

// bindWowFishBank binds a generic wrapper to an already deployed contract.
func bindWowFishBank(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := WowFishBankMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WowFishBank *WowFishBankRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WowFishBank.Contract.WowFishBankCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WowFishBank *WowFishBankRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WowFishBank.Contract.WowFishBankTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WowFishBank *WowFishBankRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WowFishBank.Contract.WowFishBankTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WowFishBank *WowFishBankCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WowFishBank.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WowFishBank *WowFishBankTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WowFishBank.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WowFishBank *WowFishBankTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WowFishBank.Contract.contract.Transact(opts, method, params...)
}

// DayWithdrawAmount is a free data retrieval call binding the contract method 0xb6700305.
//
// Solidity: function _dayWithdrawAmount() view returns(uint256)
func (_WowFishBank *WowFishBankCaller) DayWithdrawAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _WowFishBank.contract.Call(opts, &out, "_dayWithdrawAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DayWithdrawAmount is a free data retrieval call binding the contract method 0xb6700305.
//
// Solidity: function _dayWithdrawAmount() view returns(uint256)
func (_WowFishBank *WowFishBankSession) DayWithdrawAmount() (*big.Int, error) {
	return _WowFishBank.Contract.DayWithdrawAmount(&_WowFishBank.CallOpts)
}

// DayWithdrawAmount is a free data retrieval call binding the contract method 0xb6700305.
//
// Solidity: function _dayWithdrawAmount() view returns(uint256)
func (_WowFishBank *WowFishBankCallerSession) DayWithdrawAmount() (*big.Int, error) {
	return _WowFishBank.Contract.DayWithdrawAmount(&_WowFishBank.CallOpts)
}

// TotalAmount is a free data retrieval call binding the contract method 0x3bbeaab5.
//
// Solidity: function _totalAmount() view returns(uint256)
func (_WowFishBank *WowFishBankCaller) TotalAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _WowFishBank.contract.Call(opts, &out, "_totalAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalAmount is a free data retrieval call binding the contract method 0x3bbeaab5.
//
// Solidity: function _totalAmount() view returns(uint256)
func (_WowFishBank *WowFishBankSession) TotalAmount() (*big.Int, error) {
	return _WowFishBank.Contract.TotalAmount(&_WowFishBank.CallOpts)
}

// TotalAmount is a free data retrieval call binding the contract method 0x3bbeaab5.
//
// Solidity: function _totalAmount() view returns(uint256)
func (_WowFishBank *WowFishBankCallerSession) TotalAmount() (*big.Int, error) {
	return _WowFishBank.Contract.TotalAmount(&_WowFishBank.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_WowFishBank *WowFishBankCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _WowFishBank.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_WowFishBank *WowFishBankSession) Owner() (common.Address, error) {
	return _WowFishBank.Contract.Owner(&_WowFishBank.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_WowFishBank *WowFishBankCallerSession) Owner() (common.Address, error) {
	return _WowFishBank.Contract.Owner(&_WowFishBank.CallOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_WowFishBank *WowFishBankTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WowFishBank.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_WowFishBank *WowFishBankSession) RenounceOwnership() (*types.Transaction, error) {
	return _WowFishBank.Contract.RenounceOwnership(&_WowFishBank.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_WowFishBank *WowFishBankTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _WowFishBank.Contract.RenounceOwnership(&_WowFishBank.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_WowFishBank *WowFishBankTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _WowFishBank.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_WowFishBank *WowFishBankSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _WowFishBank.Contract.TransferOwnership(&_WowFishBank.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_WowFishBank *WowFishBankTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _WowFishBank.Contract.TransferOwnership(&_WowFishBank.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0x51cff8d9.
//
// Solidity: function withdraw(address to) returns()
func (_WowFishBank *WowFishBankTransactor) Withdraw(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _WowFishBank.contract.Transact(opts, "withdraw", to)
}

// Withdraw is a paid mutator transaction binding the contract method 0x51cff8d9.
//
// Solidity: function withdraw(address to) returns()
func (_WowFishBank *WowFishBankSession) Withdraw(to common.Address) (*types.Transaction, error) {
	return _WowFishBank.Contract.Withdraw(&_WowFishBank.TransactOpts, to)
}

// Withdraw is a paid mutator transaction binding the contract method 0x51cff8d9.
//
// Solidity: function withdraw(address to) returns()
func (_WowFishBank *WowFishBankTransactorSession) Withdraw(to common.Address) (*types.Transaction, error) {
	return _WowFishBank.Contract.Withdraw(&_WowFishBank.TransactOpts, to)
}

// WowFishBankOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the WowFishBank contract.
type WowFishBankOwnershipTransferredIterator struct {
	Event *WowFishBankOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *WowFishBankOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WowFishBankOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(WowFishBankOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *WowFishBankOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WowFishBankOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WowFishBankOwnershipTransferred represents a OwnershipTransferred event raised by the WowFishBank contract.
type WowFishBankOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_WowFishBank *WowFishBankFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*WowFishBankOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _WowFishBank.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &WowFishBankOwnershipTransferredIterator{contract: _WowFishBank.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_WowFishBank *WowFishBankFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *WowFishBankOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _WowFishBank.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WowFishBankOwnershipTransferred)
				if err := _WowFishBank.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_WowFishBank *WowFishBankFilterer) ParseOwnershipTransferred(log types.Log) (*WowFishBankOwnershipTransferred, error) {
	event := new(WowFishBankOwnershipTransferred)
	if err := _WowFishBank.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
