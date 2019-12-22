// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package stakingManager

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// StakingABI is the input ABI used to generate the binding from.
const StakingABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"transcoderAddr\",\"type\":\"address\"},{\"name\":\"delegAddr\",\"type\":\"address\"}],\"name\":\"getTranscodersDelegatorStake\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawTranscoderStake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"withdrawalTranscoderProposal\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"getTranscodersBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minRegStake\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_transcoderAddr\",\"type\":\"address\"}],\"name\":\"delegateToTranscoder\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"isTranscoder\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minStake\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"transcoderApprovalPeriod\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"getTranscodersCurrentAverageStake\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"isTranscoderBelowStake\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"getTranscodersStake\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"delegatorApprovalPeriod\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"transcoderAddr\",\"type\":\"address\"}],\"name\":\"withdrawDelegatorStake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"delegators\",\"outputs\":[{\"name\":\"timestamp\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"transcoderRegister\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"transcoderAddr\",\"type\":\"address\"}],\"name\":\"slashing\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"transcoderAddr\",\"type\":\"address\"}],\"name\":\"ban\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"selfStake\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"transcoderAddr\",\"type\":\"address\"}],\"name\":\"getLatestStakeTime\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"withdrawalPeriod\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"delegatorRegister\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"transcoderAddr\",\"type\":\"address\"}],\"name\":\"withdrawalDelegatorProposal\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"transcoders\",\"outputs\":[{\"name\":\"total\",\"type\":\"uint256\"},{\"name\":\"timestamp\",\"type\":\"uint256\"},{\"name\":\"stakeTimestamp\",\"type\":\"uint256\"},{\"name\":\"withdrawalProposition\",\"type\":\"uint256\"},{\"name\":\"banned\",\"type\":\"bool\"},{\"name\":\"belowStake\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"getTranscodersTotalAverageStake\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"isDelegator\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_minStake\",\"type\":\"uint256\"},{\"name\":\"_minRegStake\",\"type\":\"uint256\"},{\"name\":\"_delegatorApprovalPeriod\",\"type\":\"uint256\"},{\"name\":\"_transcoderApprovalPeriod\",\"type\":\"uint256\"},{\"name\":\"_withdrawalPeriod\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"transcoder\",\"type\":\"address\"}],\"name\":\"TranscoderRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"delegator\",\"type\":\"address\"}],\"name\":\"DelegatorRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"transcoder\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Staked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"transcoder\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Delegated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"transcoder\",\"type\":\"address\"}],\"name\":\"Slashed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"transcoder\",\"type\":\"address\"}],\"name\":\"Banned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"transcoder\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TranscoderStakeWithdrawal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"transcoder\",\"type\":\"address\"}],\"name\":\"TranscoderStakeWithdrawalProposition\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"transcoder\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"DelegatorStakeWithdrawal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"transcoder\",\"type\":\"address\"}],\"name\":\"DelegatorStakeWithdrawalProposition\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// Staking is an auto generated Go binding around an Ethereum contract.
type Staking struct {
	StakingCaller     // Read-only binding to the contract
	StakingTransactor // Write-only binding to the contract
	StakingFilterer   // Log filterer for contract events
}

// StakingCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakingSession struct {
	Contract     *Staking          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakingCallerSession struct {
	Contract *StakingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// StakingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakingTransactorSession struct {
	Contract     *StakingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// StakingRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakingRaw struct {
	Contract *Staking // Generic contract binding to access the raw methods on
}

// StakingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakingCallerRaw struct {
	Contract *StakingCaller // Generic read-only contract binding to access the raw methods on
}

// StakingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakingTransactorRaw struct {
	Contract *StakingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStaking creates a new instance of Staking, bound to a specific deployed contract.
func NewStaking(address common.Address, backend bind.ContractBackend) (*Staking, error) {
	contract, err := bindStaking(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Staking{StakingCaller: StakingCaller{contract: contract}, StakingTransactor: StakingTransactor{contract: contract}, StakingFilterer: StakingFilterer{contract: contract}}, nil
}

// NewStakingCaller creates a new read-only instance of Staking, bound to a specific deployed contract.
func NewStakingCaller(address common.Address, caller bind.ContractCaller) (*StakingCaller, error) {
	contract, err := bindStaking(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakingCaller{contract: contract}, nil
}

// NewStakingTransactor creates a new write-only instance of Staking, bound to a specific deployed contract.
func NewStakingTransactor(address common.Address, transactor bind.ContractTransactor) (*StakingTransactor, error) {
	contract, err := bindStaking(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakingTransactor{contract: contract}, nil
}

// NewStakingFilterer creates a new log filterer instance of Staking, bound to a specific deployed contract.
func NewStakingFilterer(address common.Address, filterer bind.ContractFilterer) (*StakingFilterer, error) {
	contract, err := bindStaking(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakingFilterer{contract: contract}, nil
}

// bindStaking binds a generic wrapper to an already deployed contract.
func bindStaking(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StakingABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Staking *StakingRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Staking.Contract.StakingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Staking *StakingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.Contract.StakingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Staking *StakingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Staking.Contract.StakingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Staking *StakingCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Staking.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Staking *StakingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Staking *StakingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Staking.Contract.contract.Transact(opts, method, params...)
}

// DelegatorApprovalPeriod is a free data retrieval call binding the contract method 0x7d36c792.
//
// Solidity: function delegatorApprovalPeriod() constant returns(uint256)
func (_Staking *StakingCaller) DelegatorApprovalPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Staking.contract.Call(opts, out, "delegatorApprovalPeriod")
	return *ret0, err
}

// DelegatorApprovalPeriod is a free data retrieval call binding the contract method 0x7d36c792.
//
// Solidity: function delegatorApprovalPeriod() constant returns(uint256)
func (_Staking *StakingSession) DelegatorApprovalPeriod() (*big.Int, error) {
	return _Staking.Contract.DelegatorApprovalPeriod(&_Staking.CallOpts)
}

// DelegatorApprovalPeriod is a free data retrieval call binding the contract method 0x7d36c792.
//
// Solidity: function delegatorApprovalPeriod() constant returns(uint256)
func (_Staking *StakingCallerSession) DelegatorApprovalPeriod() (*big.Int, error) {
	return _Staking.Contract.DelegatorApprovalPeriod(&_Staking.CallOpts)
}

// Delegators is a free data retrieval call binding the contract method 0x8d23fc61.
//
// Solidity: function delegators(address ) constant returns(uint256 timestamp)
func (_Staking *StakingCaller) Delegators(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Staking.contract.Call(opts, out, "delegators", arg0)
	return *ret0, err
}

// Delegators is a free data retrieval call binding the contract method 0x8d23fc61.
//
// Solidity: function delegators(address ) constant returns(uint256 timestamp)
func (_Staking *StakingSession) Delegators(arg0 common.Address) (*big.Int, error) {
	return _Staking.Contract.Delegators(&_Staking.CallOpts, arg0)
}

// Delegators is a free data retrieval call binding the contract method 0x8d23fc61.
//
// Solidity: function delegators(address ) constant returns(uint256 timestamp)
func (_Staking *StakingCallerSession) Delegators(arg0 common.Address) (*big.Int, error) {
	return _Staking.Contract.Delegators(&_Staking.CallOpts, arg0)
}

// GetLatestStakeTime is a free data retrieval call binding the contract method 0x9d4f4f3c.
//
// Solidity: function getLatestStakeTime(address transcoderAddr) constant returns(uint256)
func (_Staking *StakingCaller) GetLatestStakeTime(opts *bind.CallOpts, transcoderAddr common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Staking.contract.Call(opts, out, "getLatestStakeTime", transcoderAddr)
	return *ret0, err
}

// GetLatestStakeTime is a free data retrieval call binding the contract method 0x9d4f4f3c.
//
// Solidity: function getLatestStakeTime(address transcoderAddr) constant returns(uint256)
func (_Staking *StakingSession) GetLatestStakeTime(transcoderAddr common.Address) (*big.Int, error) {
	return _Staking.Contract.GetLatestStakeTime(&_Staking.CallOpts, transcoderAddr)
}

// GetLatestStakeTime is a free data retrieval call binding the contract method 0x9d4f4f3c.
//
// Solidity: function getLatestStakeTime(address transcoderAddr) constant returns(uint256)
func (_Staking *StakingCallerSession) GetLatestStakeTime(transcoderAddr common.Address) (*big.Int, error) {
	return _Staking.Contract.GetLatestStakeTime(&_Staking.CallOpts, transcoderAddr)
}

// GetTranscodersBalance is a free data retrieval call binding the contract method 0x2cfe7261.
//
// Solidity: function getTranscodersBalance(address _addr) constant returns(uint256)
func (_Staking *StakingCaller) GetTranscodersBalance(opts *bind.CallOpts, _addr common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Staking.contract.Call(opts, out, "getTranscodersBalance", _addr)
	return *ret0, err
}

// GetTranscodersBalance is a free data retrieval call binding the contract method 0x2cfe7261.
//
// Solidity: function getTranscodersBalance(address _addr) constant returns(uint256)
func (_Staking *StakingSession) GetTranscodersBalance(_addr common.Address) (*big.Int, error) {
	return _Staking.Contract.GetTranscodersBalance(&_Staking.CallOpts, _addr)
}

// GetTranscodersBalance is a free data retrieval call binding the contract method 0x2cfe7261.
//
// Solidity: function getTranscodersBalance(address _addr) constant returns(uint256)
func (_Staking *StakingCallerSession) GetTranscodersBalance(_addr common.Address) (*big.Int, error) {
	return _Staking.Contract.GetTranscodersBalance(&_Staking.CallOpts, _addr)
}

// GetTranscodersCurrentAverageStake is a free data retrieval call binding the contract method 0x6a061e1f.
//
// Solidity: function getTranscodersCurrentAverageStake(address _addr) constant returns(uint256)
func (_Staking *StakingCaller) GetTranscodersCurrentAverageStake(opts *bind.CallOpts, _addr common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Staking.contract.Call(opts, out, "getTranscodersCurrentAverageStake", _addr)
	return *ret0, err
}

// GetTranscodersCurrentAverageStake is a free data retrieval call binding the contract method 0x6a061e1f.
//
// Solidity: function getTranscodersCurrentAverageStake(address _addr) constant returns(uint256)
func (_Staking *StakingSession) GetTranscodersCurrentAverageStake(_addr common.Address) (*big.Int, error) {
	return _Staking.Contract.GetTranscodersCurrentAverageStake(&_Staking.CallOpts, _addr)
}

// GetTranscodersCurrentAverageStake is a free data retrieval call binding the contract method 0x6a061e1f.
//
// Solidity: function getTranscodersCurrentAverageStake(address _addr) constant returns(uint256)
func (_Staking *StakingCallerSession) GetTranscodersCurrentAverageStake(_addr common.Address) (*big.Int, error) {
	return _Staking.Contract.GetTranscodersCurrentAverageStake(&_Staking.CallOpts, _addr)
}

// GetTranscodersDelegatorStake is a free data retrieval call binding the contract method 0x006215b1.
//
// Solidity: function getTranscodersDelegatorStake(address transcoderAddr, address delegAddr) constant returns(uint256)
func (_Staking *StakingCaller) GetTranscodersDelegatorStake(opts *bind.CallOpts, transcoderAddr common.Address, delegAddr common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Staking.contract.Call(opts, out, "getTranscodersDelegatorStake", transcoderAddr, delegAddr)
	return *ret0, err
}

// GetTranscodersDelegatorStake is a free data retrieval call binding the contract method 0x006215b1.
//
// Solidity: function getTranscodersDelegatorStake(address transcoderAddr, address delegAddr) constant returns(uint256)
func (_Staking *StakingSession) GetTranscodersDelegatorStake(transcoderAddr common.Address, delegAddr common.Address) (*big.Int, error) {
	return _Staking.Contract.GetTranscodersDelegatorStake(&_Staking.CallOpts, transcoderAddr, delegAddr)
}

// GetTranscodersDelegatorStake is a free data retrieval call binding the contract method 0x006215b1.
//
// Solidity: function getTranscodersDelegatorStake(address transcoderAddr, address delegAddr) constant returns(uint256)
func (_Staking *StakingCallerSession) GetTranscodersDelegatorStake(transcoderAddr common.Address, delegAddr common.Address) (*big.Int, error) {
	return _Staking.Contract.GetTranscodersDelegatorStake(&_Staking.CallOpts, transcoderAddr, delegAddr)
}

// GetTranscodersStake is a free data retrieval call binding the contract method 0x7af96ab0.
//
// Solidity: function getTranscodersStake(address _addr) constant returns(uint256)
func (_Staking *StakingCaller) GetTranscodersStake(opts *bind.CallOpts, _addr common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Staking.contract.Call(opts, out, "getTranscodersStake", _addr)
	return *ret0, err
}

// GetTranscodersStake is a free data retrieval call binding the contract method 0x7af96ab0.
//
// Solidity: function getTranscodersStake(address _addr) constant returns(uint256)
func (_Staking *StakingSession) GetTranscodersStake(_addr common.Address) (*big.Int, error) {
	return _Staking.Contract.GetTranscodersStake(&_Staking.CallOpts, _addr)
}

// GetTranscodersStake is a free data retrieval call binding the contract method 0x7af96ab0.
//
// Solidity: function getTranscodersStake(address _addr) constant returns(uint256)
func (_Staking *StakingCallerSession) GetTranscodersStake(_addr common.Address) (*big.Int, error) {
	return _Staking.Contract.GetTranscodersStake(&_Staking.CallOpts, _addr)
}

// GetTranscodersTotalAverageStake is a free data retrieval call binding the contract method 0xf01c1438.
//
// Solidity: function getTranscodersTotalAverageStake(address _addr) constant returns(uint256)
func (_Staking *StakingCaller) GetTranscodersTotalAverageStake(opts *bind.CallOpts, _addr common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Staking.contract.Call(opts, out, "getTranscodersTotalAverageStake", _addr)
	return *ret0, err
}

// GetTranscodersTotalAverageStake is a free data retrieval call binding the contract method 0xf01c1438.
//
// Solidity: function getTranscodersTotalAverageStake(address _addr) constant returns(uint256)
func (_Staking *StakingSession) GetTranscodersTotalAverageStake(_addr common.Address) (*big.Int, error) {
	return _Staking.Contract.GetTranscodersTotalAverageStake(&_Staking.CallOpts, _addr)
}

// GetTranscodersTotalAverageStake is a free data retrieval call binding the contract method 0xf01c1438.
//
// Solidity: function getTranscodersTotalAverageStake(address _addr) constant returns(uint256)
func (_Staking *StakingCallerSession) GetTranscodersTotalAverageStake(_addr common.Address) (*big.Int, error) {
	return _Staking.Contract.GetTranscodersTotalAverageStake(&_Staking.CallOpts, _addr)
}

// IsDelegator is a free data retrieval call binding the contract method 0xfd8ab482.
//
// Solidity: function isDelegator(address _addr) constant returns(bool)
func (_Staking *StakingCaller) IsDelegator(opts *bind.CallOpts, _addr common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Staking.contract.Call(opts, out, "isDelegator", _addr)
	return *ret0, err
}

// IsDelegator is a free data retrieval call binding the contract method 0xfd8ab482.
//
// Solidity: function isDelegator(address _addr) constant returns(bool)
func (_Staking *StakingSession) IsDelegator(_addr common.Address) (bool, error) {
	return _Staking.Contract.IsDelegator(&_Staking.CallOpts, _addr)
}

// IsDelegator is a free data retrieval call binding the contract method 0xfd8ab482.
//
// Solidity: function isDelegator(address _addr) constant returns(bool)
func (_Staking *StakingCallerSession) IsDelegator(_addr common.Address) (bool, error) {
	return _Staking.Contract.IsDelegator(&_Staking.CallOpts, _addr)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Staking *StakingCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Staking.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Staking *StakingSession) IsOwner() (bool, error) {
	return _Staking.Contract.IsOwner(&_Staking.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Staking *StakingCallerSession) IsOwner() (bool, error) {
	return _Staking.Contract.IsOwner(&_Staking.CallOpts)
}

// IsTranscoder is a free data retrieval call binding the contract method 0x34df9dc4.
//
// Solidity: function isTranscoder(address _addr) constant returns(bool)
func (_Staking *StakingCaller) IsTranscoder(opts *bind.CallOpts, _addr common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Staking.contract.Call(opts, out, "isTranscoder", _addr)
	return *ret0, err
}

// IsTranscoder is a free data retrieval call binding the contract method 0x34df9dc4.
//
// Solidity: function isTranscoder(address _addr) constant returns(bool)
func (_Staking *StakingSession) IsTranscoder(_addr common.Address) (bool, error) {
	return _Staking.Contract.IsTranscoder(&_Staking.CallOpts, _addr)
}

// IsTranscoder is a free data retrieval call binding the contract method 0x34df9dc4.
//
// Solidity: function isTranscoder(address _addr) constant returns(bool)
func (_Staking *StakingCallerSession) IsTranscoder(_addr common.Address) (bool, error) {
	return _Staking.Contract.IsTranscoder(&_Staking.CallOpts, _addr)
}

// IsTranscoderBelowStake is a free data retrieval call binding the contract method 0x794c5425.
//
// Solidity: function isTranscoderBelowStake(address _addr) constant returns(bool)
func (_Staking *StakingCaller) IsTranscoderBelowStake(opts *bind.CallOpts, _addr common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Staking.contract.Call(opts, out, "isTranscoderBelowStake", _addr)
	return *ret0, err
}

// IsTranscoderBelowStake is a free data retrieval call binding the contract method 0x794c5425.
//
// Solidity: function isTranscoderBelowStake(address _addr) constant returns(bool)
func (_Staking *StakingSession) IsTranscoderBelowStake(_addr common.Address) (bool, error) {
	return _Staking.Contract.IsTranscoderBelowStake(&_Staking.CallOpts, _addr)
}

// IsTranscoderBelowStake is a free data retrieval call binding the contract method 0x794c5425.
//
// Solidity: function isTranscoderBelowStake(address _addr) constant returns(bool)
func (_Staking *StakingCallerSession) IsTranscoderBelowStake(_addr common.Address) (bool, error) {
	return _Staking.Contract.IsTranscoderBelowStake(&_Staking.CallOpts, _addr)
}

// MinRegStake is a free data retrieval call binding the contract method 0x2e8071dd.
//
// Solidity: function minRegStake() constant returns(uint256)
func (_Staking *StakingCaller) MinRegStake(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Staking.contract.Call(opts, out, "minRegStake")
	return *ret0, err
}

// MinRegStake is a free data retrieval call binding the contract method 0x2e8071dd.
//
// Solidity: function minRegStake() constant returns(uint256)
func (_Staking *StakingSession) MinRegStake() (*big.Int, error) {
	return _Staking.Contract.MinRegStake(&_Staking.CallOpts)
}

// MinRegStake is a free data retrieval call binding the contract method 0x2e8071dd.
//
// Solidity: function minRegStake() constant returns(uint256)
func (_Staking *StakingCallerSession) MinRegStake() (*big.Int, error) {
	return _Staking.Contract.MinRegStake(&_Staking.CallOpts)
}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() constant returns(uint256)
func (_Staking *StakingCaller) MinStake(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Staking.contract.Call(opts, out, "minStake")
	return *ret0, err
}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() constant returns(uint256)
func (_Staking *StakingSession) MinStake() (*big.Int, error) {
	return _Staking.Contract.MinStake(&_Staking.CallOpts)
}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() constant returns(uint256)
func (_Staking *StakingCallerSession) MinStake() (*big.Int, error) {
	return _Staking.Contract.MinStake(&_Staking.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Staking *StakingCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Staking.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Staking *StakingSession) Owner() (common.Address, error) {
	return _Staking.Contract.Owner(&_Staking.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Staking *StakingCallerSession) Owner() (common.Address, error) {
	return _Staking.Contract.Owner(&_Staking.CallOpts)
}

// TranscoderApprovalPeriod is a free data retrieval call binding the contract method 0x3939e608.
//
// Solidity: function transcoderApprovalPeriod() constant returns(uint256)
func (_Staking *StakingCaller) TranscoderApprovalPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Staking.contract.Call(opts, out, "transcoderApprovalPeriod")
	return *ret0, err
}

// TranscoderApprovalPeriod is a free data retrieval call binding the contract method 0x3939e608.
//
// Solidity: function transcoderApprovalPeriod() constant returns(uint256)
func (_Staking *StakingSession) TranscoderApprovalPeriod() (*big.Int, error) {
	return _Staking.Contract.TranscoderApprovalPeriod(&_Staking.CallOpts)
}

// TranscoderApprovalPeriod is a free data retrieval call binding the contract method 0x3939e608.
//
// Solidity: function transcoderApprovalPeriod() constant returns(uint256)
func (_Staking *StakingCallerSession) TranscoderApprovalPeriod() (*big.Int, error) {
	return _Staking.Contract.TranscoderApprovalPeriod(&_Staking.CallOpts)
}

// Transcoders is a free data retrieval call binding the contract method 0xe2dc17f6.
//
// Solidity: function transcoders(address ) constant returns(uint256 total, uint256 timestamp, uint256 stakeTimestamp, uint256 withdrawalProposition, bool banned, bool belowStake)
func (_Staking *StakingCaller) Transcoders(opts *bind.CallOpts, arg0 common.Address) (struct {
	Total                 *big.Int
	Timestamp             *big.Int
	StakeTimestamp        *big.Int
	WithdrawalProposition *big.Int
	Banned                bool
	BelowStake            bool
}, error) {
	ret := new(struct {
		Total                 *big.Int
		Timestamp             *big.Int
		StakeTimestamp        *big.Int
		WithdrawalProposition *big.Int
		Banned                bool
		BelowStake            bool
	})
	out := ret
	err := _Staking.contract.Call(opts, out, "transcoders", arg0)
	return *ret, err
}

// Transcoders is a free data retrieval call binding the contract method 0xe2dc17f6.
//
// Solidity: function transcoders(address ) constant returns(uint256 total, uint256 timestamp, uint256 stakeTimestamp, uint256 withdrawalProposition, bool banned, bool belowStake)
func (_Staking *StakingSession) Transcoders(arg0 common.Address) (struct {
	Total                 *big.Int
	Timestamp             *big.Int
	StakeTimestamp        *big.Int
	WithdrawalProposition *big.Int
	Banned                bool
	BelowStake            bool
}, error) {
	return _Staking.Contract.Transcoders(&_Staking.CallOpts, arg0)
}

// Transcoders is a free data retrieval call binding the contract method 0xe2dc17f6.
//
// Solidity: function transcoders(address ) constant returns(uint256 total, uint256 timestamp, uint256 stakeTimestamp, uint256 withdrawalProposition, bool banned, bool belowStake)
func (_Staking *StakingCallerSession) Transcoders(arg0 common.Address) (struct {
	Total                 *big.Int
	Timestamp             *big.Int
	StakeTimestamp        *big.Int
	WithdrawalProposition *big.Int
	Banned                bool
	BelowStake            bool
}, error) {
	return _Staking.Contract.Transcoders(&_Staking.CallOpts, arg0)
}

// WithdrawalPeriod is a free data retrieval call binding the contract method 0xbca7093d.
//
// Solidity: function withdrawalPeriod() constant returns(uint256)
func (_Staking *StakingCaller) WithdrawalPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Staking.contract.Call(opts, out, "withdrawalPeriod")
	return *ret0, err
}

// WithdrawalPeriod is a free data retrieval call binding the contract method 0xbca7093d.
//
// Solidity: function withdrawalPeriod() constant returns(uint256)
func (_Staking *StakingSession) WithdrawalPeriod() (*big.Int, error) {
	return _Staking.Contract.WithdrawalPeriod(&_Staking.CallOpts)
}

// WithdrawalPeriod is a free data retrieval call binding the contract method 0xbca7093d.
//
// Solidity: function withdrawalPeriod() constant returns(uint256)
func (_Staking *StakingCallerSession) WithdrawalPeriod() (*big.Int, error) {
	return _Staking.Contract.WithdrawalPeriod(&_Staking.CallOpts)
}

// Ban is a paid mutator transaction binding the contract method 0x97c3ccd8.
//
// Solidity: function ban(address transcoderAddr) returns()
func (_Staking *StakingTransactor) Ban(opts *bind.TransactOpts, transcoderAddr common.Address) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "ban", transcoderAddr)
}

// Ban is a paid mutator transaction binding the contract method 0x97c3ccd8.
//
// Solidity: function ban(address transcoderAddr) returns()
func (_Staking *StakingSession) Ban(transcoderAddr common.Address) (*types.Transaction, error) {
	return _Staking.Contract.Ban(&_Staking.TransactOpts, transcoderAddr)
}

// Ban is a paid mutator transaction binding the contract method 0x97c3ccd8.
//
// Solidity: function ban(address transcoderAddr) returns()
func (_Staking *StakingTransactorSession) Ban(transcoderAddr common.Address) (*types.Transaction, error) {
	return _Staking.Contract.Ban(&_Staking.TransactOpts, transcoderAddr)
}

// DelegateToTranscoder is a paid mutator transaction binding the contract method 0x34814020.
//
// Solidity: function delegateToTranscoder(address _transcoderAddr) returns()
func (_Staking *StakingTransactor) DelegateToTranscoder(opts *bind.TransactOpts, _transcoderAddr common.Address) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "delegateToTranscoder", _transcoderAddr)
}

// DelegateToTranscoder is a paid mutator transaction binding the contract method 0x34814020.
//
// Solidity: function delegateToTranscoder(address _transcoderAddr) returns()
func (_Staking *StakingSession) DelegateToTranscoder(_transcoderAddr common.Address) (*types.Transaction, error) {
	return _Staking.Contract.DelegateToTranscoder(&_Staking.TransactOpts, _transcoderAddr)
}

// DelegateToTranscoder is a paid mutator transaction binding the contract method 0x34814020.
//
// Solidity: function delegateToTranscoder(address _transcoderAddr) returns()
func (_Staking *StakingTransactorSession) DelegateToTranscoder(_transcoderAddr common.Address) (*types.Transaction, error) {
	return _Staking.Contract.DelegateToTranscoder(&_Staking.TransactOpts, _transcoderAddr)
}

// DelegatorRegister is a paid mutator transaction binding the contract method 0xca7abdd4.
//
// Solidity: function delegatorRegister() returns()
func (_Staking *StakingTransactor) DelegatorRegister(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "delegatorRegister")
}

// DelegatorRegister is a paid mutator transaction binding the contract method 0xca7abdd4.
//
// Solidity: function delegatorRegister() returns()
func (_Staking *StakingSession) DelegatorRegister() (*types.Transaction, error) {
	return _Staking.Contract.DelegatorRegister(&_Staking.TransactOpts)
}

// DelegatorRegister is a paid mutator transaction binding the contract method 0xca7abdd4.
//
// Solidity: function delegatorRegister() returns()
func (_Staking *StakingTransactorSession) DelegatorRegister() (*types.Transaction, error) {
	return _Staking.Contract.DelegatorRegister(&_Staking.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Staking *StakingTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Staking *StakingSession) RenounceOwnership() (*types.Transaction, error) {
	return _Staking.Contract.RenounceOwnership(&_Staking.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Staking *StakingTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Staking.Contract.RenounceOwnership(&_Staking.TransactOpts)
}

// SelfStake is a paid mutator transaction binding the contract method 0x9bdd9480.
//
// Solidity: function selfStake() returns()
func (_Staking *StakingTransactor) SelfStake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "selfStake")
}

// SelfStake is a paid mutator transaction binding the contract method 0x9bdd9480.
//
// Solidity: function selfStake() returns()
func (_Staking *StakingSession) SelfStake() (*types.Transaction, error) {
	return _Staking.Contract.SelfStake(&_Staking.TransactOpts)
}

// SelfStake is a paid mutator transaction binding the contract method 0x9bdd9480.
//
// Solidity: function selfStake() returns()
func (_Staking *StakingTransactorSession) SelfStake() (*types.Transaction, error) {
	return _Staking.Contract.SelfStake(&_Staking.TransactOpts)
}

// Slashing is a paid mutator transaction binding the contract method 0x90d4582f.
//
// Solidity: function slashing(address transcoderAddr) returns()
func (_Staking *StakingTransactor) Slashing(opts *bind.TransactOpts, transcoderAddr common.Address) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "slashing", transcoderAddr)
}

// Slashing is a paid mutator transaction binding the contract method 0x90d4582f.
//
// Solidity: function slashing(address transcoderAddr) returns()
func (_Staking *StakingSession) Slashing(transcoderAddr common.Address) (*types.Transaction, error) {
	return _Staking.Contract.Slashing(&_Staking.TransactOpts, transcoderAddr)
}

// Slashing is a paid mutator transaction binding the contract method 0x90d4582f.
//
// Solidity: function slashing(address transcoderAddr) returns()
func (_Staking *StakingTransactorSession) Slashing(transcoderAddr common.Address) (*types.Transaction, error) {
	return _Staking.Contract.Slashing(&_Staking.TransactOpts, transcoderAddr)
}

// TranscoderRegister is a paid mutator transaction binding the contract method 0x8fe0e81d.
//
// Solidity: function transcoderRegister() returns()
func (_Staking *StakingTransactor) TranscoderRegister(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "transcoderRegister")
}

// TranscoderRegister is a paid mutator transaction binding the contract method 0x8fe0e81d.
//
// Solidity: function transcoderRegister() returns()
func (_Staking *StakingSession) TranscoderRegister() (*types.Transaction, error) {
	return _Staking.Contract.TranscoderRegister(&_Staking.TransactOpts)
}

// TranscoderRegister is a paid mutator transaction binding the contract method 0x8fe0e81d.
//
// Solidity: function transcoderRegister() returns()
func (_Staking *StakingTransactorSession) TranscoderRegister() (*types.Transaction, error) {
	return _Staking.Contract.TranscoderRegister(&_Staking.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Staking *StakingTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Staking *StakingSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Staking.Contract.TransferOwnership(&_Staking.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Staking *StakingTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Staking.Contract.TransferOwnership(&_Staking.TransactOpts, newOwner)
}

// WithdrawDelegatorStake is a paid mutator transaction binding the contract method 0x85bffd2f.
//
// Solidity: function withdrawDelegatorStake(uint256 amount, address transcoderAddr) returns()
func (_Staking *StakingTransactor) WithdrawDelegatorStake(opts *bind.TransactOpts, amount *big.Int, transcoderAddr common.Address) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "withdrawDelegatorStake", amount, transcoderAddr)
}

// WithdrawDelegatorStake is a paid mutator transaction binding the contract method 0x85bffd2f.
//
// Solidity: function withdrawDelegatorStake(uint256 amount, address transcoderAddr) returns()
func (_Staking *StakingSession) WithdrawDelegatorStake(amount *big.Int, transcoderAddr common.Address) (*types.Transaction, error) {
	return _Staking.Contract.WithdrawDelegatorStake(&_Staking.TransactOpts, amount, transcoderAddr)
}

// WithdrawDelegatorStake is a paid mutator transaction binding the contract method 0x85bffd2f.
//
// Solidity: function withdrawDelegatorStake(uint256 amount, address transcoderAddr) returns()
func (_Staking *StakingTransactorSession) WithdrawDelegatorStake(amount *big.Int, transcoderAddr common.Address) (*types.Transaction, error) {
	return _Staking.Contract.WithdrawDelegatorStake(&_Staking.TransactOpts, amount, transcoderAddr)
}

// WithdrawTranscoderStake is a paid mutator transaction binding the contract method 0x0d8877e8.
//
// Solidity: function withdrawTranscoderStake(uint256 amount) returns()
func (_Staking *StakingTransactor) WithdrawTranscoderStake(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "withdrawTranscoderStake", amount)
}

// WithdrawTranscoderStake is a paid mutator transaction binding the contract method 0x0d8877e8.
//
// Solidity: function withdrawTranscoderStake(uint256 amount) returns()
func (_Staking *StakingSession) WithdrawTranscoderStake(amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.WithdrawTranscoderStake(&_Staking.TransactOpts, amount)
}

// WithdrawTranscoderStake is a paid mutator transaction binding the contract method 0x0d8877e8.
//
// Solidity: function withdrawTranscoderStake(uint256 amount) returns()
func (_Staking *StakingTransactorSession) WithdrawTranscoderStake(amount *big.Int) (*types.Transaction, error) {
	return _Staking.Contract.WithdrawTranscoderStake(&_Staking.TransactOpts, amount)
}

// WithdrawalDelegatorProposal is a paid mutator transaction binding the contract method 0xd81df838.
//
// Solidity: function withdrawalDelegatorProposal(address transcoderAddr) returns()
func (_Staking *StakingTransactor) WithdrawalDelegatorProposal(opts *bind.TransactOpts, transcoderAddr common.Address) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "withdrawalDelegatorProposal", transcoderAddr)
}

// WithdrawalDelegatorProposal is a paid mutator transaction binding the contract method 0xd81df838.
//
// Solidity: function withdrawalDelegatorProposal(address transcoderAddr) returns()
func (_Staking *StakingSession) WithdrawalDelegatorProposal(transcoderAddr common.Address) (*types.Transaction, error) {
	return _Staking.Contract.WithdrawalDelegatorProposal(&_Staking.TransactOpts, transcoderAddr)
}

// WithdrawalDelegatorProposal is a paid mutator transaction binding the contract method 0xd81df838.
//
// Solidity: function withdrawalDelegatorProposal(address transcoderAddr) returns()
func (_Staking *StakingTransactorSession) WithdrawalDelegatorProposal(transcoderAddr common.Address) (*types.Transaction, error) {
	return _Staking.Contract.WithdrawalDelegatorProposal(&_Staking.TransactOpts, transcoderAddr)
}

// WithdrawalTranscoderProposal is a paid mutator transaction binding the contract method 0x10a8715a.
//
// Solidity: function withdrawalTranscoderProposal() returns()
func (_Staking *StakingTransactor) WithdrawalTranscoderProposal(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Staking.contract.Transact(opts, "withdrawalTranscoderProposal")
}

// WithdrawalTranscoderProposal is a paid mutator transaction binding the contract method 0x10a8715a.
//
// Solidity: function withdrawalTranscoderProposal() returns()
func (_Staking *StakingSession) WithdrawalTranscoderProposal() (*types.Transaction, error) {
	return _Staking.Contract.WithdrawalTranscoderProposal(&_Staking.TransactOpts)
}

// WithdrawalTranscoderProposal is a paid mutator transaction binding the contract method 0x10a8715a.
//
// Solidity: function withdrawalTranscoderProposal() returns()
func (_Staking *StakingTransactorSession) WithdrawalTranscoderProposal() (*types.Transaction, error) {
	return _Staking.Contract.WithdrawalTranscoderProposal(&_Staking.TransactOpts)
}

// StakingBannedIterator is returned from FilterBanned and is used to iterate over the raw logs and unpacked data for Banned events raised by the Staking contract.
type StakingBannedIterator struct {
	Event *StakingBanned // Event containing the contract specifics and raw log

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
func (it *StakingBannedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingBanned)
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
		it.Event = new(StakingBanned)
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
func (it *StakingBannedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingBannedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingBanned represents a Banned event raised by the Staking contract.
type StakingBanned struct {
	Transcoder common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBanned is a free log retrieval operation binding the contract event 0x30d1df1214d91553408ca5384ce29e10e5866af8423c628be22860e41fb81005.
//
// Solidity: event Banned(address indexed transcoder)
func (_Staking *StakingFilterer) FilterBanned(opts *bind.FilterOpts, transcoder []common.Address) (*StakingBannedIterator, error) {

	var transcoderRule []interface{}
	for _, transcoderItem := range transcoder {
		transcoderRule = append(transcoderRule, transcoderItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Banned", transcoderRule)
	if err != nil {
		return nil, err
	}
	return &StakingBannedIterator{contract: _Staking.contract, event: "Banned", logs: logs, sub: sub}, nil
}

// WatchBanned is a free log subscription operation binding the contract event 0x30d1df1214d91553408ca5384ce29e10e5866af8423c628be22860e41fb81005.
//
// Solidity: event Banned(address indexed transcoder)
func (_Staking *StakingFilterer) WatchBanned(opts *bind.WatchOpts, sink chan<- *StakingBanned, transcoder []common.Address) (event.Subscription, error) {

	var transcoderRule []interface{}
	for _, transcoderItem := range transcoder {
		transcoderRule = append(transcoderRule, transcoderItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Banned", transcoderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingBanned)
				if err := _Staking.contract.UnpackLog(event, "Banned", log); err != nil {
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

// ParseBanned is a log parse operation binding the contract event 0x30d1df1214d91553408ca5384ce29e10e5866af8423c628be22860e41fb81005.
//
// Solidity: event Banned(address indexed transcoder)
func (_Staking *StakingFilterer) ParseBanned(log types.Log) (*StakingBanned, error) {
	event := new(StakingBanned)
	if err := _Staking.contract.UnpackLog(event, "Banned", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakingDelegatedIterator is returned from FilterDelegated and is used to iterate over the raw logs and unpacked data for Delegated events raised by the Staking contract.
type StakingDelegatedIterator struct {
	Event *StakingDelegated // Event containing the contract specifics and raw log

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
func (it *StakingDelegatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingDelegated)
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
		it.Event = new(StakingDelegated)
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
func (it *StakingDelegatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingDelegatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingDelegated represents a Delegated event raised by the Staking contract.
type StakingDelegated struct {
	Transcoder common.Address
	Delegator  common.Address
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterDelegated is a free log retrieval operation binding the contract event 0xe5541a6b6103d4fa7e021ed54fad39c66f27a76bd13d374cf6240ae6bd0bb72b.
//
// Solidity: event Delegated(address indexed transcoder, address indexed delegator, uint256 indexed amount)
func (_Staking *StakingFilterer) FilterDelegated(opts *bind.FilterOpts, transcoder []common.Address, delegator []common.Address, amount []*big.Int) (*StakingDelegatedIterator, error) {

	var transcoderRule []interface{}
	for _, transcoderItem := range transcoder {
		transcoderRule = append(transcoderRule, transcoderItem)
	}
	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Delegated", transcoderRule, delegatorRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &StakingDelegatedIterator{contract: _Staking.contract, event: "Delegated", logs: logs, sub: sub}, nil
}

// WatchDelegated is a free log subscription operation binding the contract event 0xe5541a6b6103d4fa7e021ed54fad39c66f27a76bd13d374cf6240ae6bd0bb72b.
//
// Solidity: event Delegated(address indexed transcoder, address indexed delegator, uint256 indexed amount)
func (_Staking *StakingFilterer) WatchDelegated(opts *bind.WatchOpts, sink chan<- *StakingDelegated, transcoder []common.Address, delegator []common.Address, amount []*big.Int) (event.Subscription, error) {

	var transcoderRule []interface{}
	for _, transcoderItem := range transcoder {
		transcoderRule = append(transcoderRule, transcoderItem)
	}
	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Delegated", transcoderRule, delegatorRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingDelegated)
				if err := _Staking.contract.UnpackLog(event, "Delegated", log); err != nil {
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

// ParseDelegated is a log parse operation binding the contract event 0xe5541a6b6103d4fa7e021ed54fad39c66f27a76bd13d374cf6240ae6bd0bb72b.
//
// Solidity: event Delegated(address indexed transcoder, address indexed delegator, uint256 indexed amount)
func (_Staking *StakingFilterer) ParseDelegated(log types.Log) (*StakingDelegated, error) {
	event := new(StakingDelegated)
	if err := _Staking.contract.UnpackLog(event, "Delegated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakingDelegatorRegisteredIterator is returned from FilterDelegatorRegistered and is used to iterate over the raw logs and unpacked data for DelegatorRegistered events raised by the Staking contract.
type StakingDelegatorRegisteredIterator struct {
	Event *StakingDelegatorRegistered // Event containing the contract specifics and raw log

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
func (it *StakingDelegatorRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingDelegatorRegistered)
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
		it.Event = new(StakingDelegatorRegistered)
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
func (it *StakingDelegatorRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingDelegatorRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingDelegatorRegistered represents a DelegatorRegistered event raised by the Staking contract.
type StakingDelegatorRegistered struct {
	Delegator common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDelegatorRegistered is a free log retrieval operation binding the contract event 0xde8d1acb0dd7d4b197ffcdeef216c2d299a44d709010af2156d0371bec0c286c.
//
// Solidity: event DelegatorRegistered(address indexed delegator)
func (_Staking *StakingFilterer) FilterDelegatorRegistered(opts *bind.FilterOpts, delegator []common.Address) (*StakingDelegatorRegisteredIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "DelegatorRegistered", delegatorRule)
	if err != nil {
		return nil, err
	}
	return &StakingDelegatorRegisteredIterator{contract: _Staking.contract, event: "DelegatorRegistered", logs: logs, sub: sub}, nil
}

// WatchDelegatorRegistered is a free log subscription operation binding the contract event 0xde8d1acb0dd7d4b197ffcdeef216c2d299a44d709010af2156d0371bec0c286c.
//
// Solidity: event DelegatorRegistered(address indexed delegator)
func (_Staking *StakingFilterer) WatchDelegatorRegistered(opts *bind.WatchOpts, sink chan<- *StakingDelegatorRegistered, delegator []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "DelegatorRegistered", delegatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingDelegatorRegistered)
				if err := _Staking.contract.UnpackLog(event, "DelegatorRegistered", log); err != nil {
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

// ParseDelegatorRegistered is a log parse operation binding the contract event 0xde8d1acb0dd7d4b197ffcdeef216c2d299a44d709010af2156d0371bec0c286c.
//
// Solidity: event DelegatorRegistered(address indexed delegator)
func (_Staking *StakingFilterer) ParseDelegatorRegistered(log types.Log) (*StakingDelegatorRegistered, error) {
	event := new(StakingDelegatorRegistered)
	if err := _Staking.contract.UnpackLog(event, "DelegatorRegistered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakingDelegatorStakeWithdrawalIterator is returned from FilterDelegatorStakeWithdrawal and is used to iterate over the raw logs and unpacked data for DelegatorStakeWithdrawal events raised by the Staking contract.
type StakingDelegatorStakeWithdrawalIterator struct {
	Event *StakingDelegatorStakeWithdrawal // Event containing the contract specifics and raw log

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
func (it *StakingDelegatorStakeWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingDelegatorStakeWithdrawal)
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
		it.Event = new(StakingDelegatorStakeWithdrawal)
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
func (it *StakingDelegatorStakeWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingDelegatorStakeWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingDelegatorStakeWithdrawal represents a DelegatorStakeWithdrawal event raised by the Staking contract.
type StakingDelegatorStakeWithdrawal struct {
	Delegator  common.Address
	Transcoder common.Address
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterDelegatorStakeWithdrawal is a free log retrieval operation binding the contract event 0xc1df1162e4be363deb10badc7a5e3d22ec7eb7f445fe5c76bea8c0e8b726840f.
//
// Solidity: event DelegatorStakeWithdrawal(address indexed delegator, address indexed transcoder, uint256 indexed amount)
func (_Staking *StakingFilterer) FilterDelegatorStakeWithdrawal(opts *bind.FilterOpts, delegator []common.Address, transcoder []common.Address, amount []*big.Int) (*StakingDelegatorStakeWithdrawalIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var transcoderRule []interface{}
	for _, transcoderItem := range transcoder {
		transcoderRule = append(transcoderRule, transcoderItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "DelegatorStakeWithdrawal", delegatorRule, transcoderRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &StakingDelegatorStakeWithdrawalIterator{contract: _Staking.contract, event: "DelegatorStakeWithdrawal", logs: logs, sub: sub}, nil
}

// WatchDelegatorStakeWithdrawal is a free log subscription operation binding the contract event 0xc1df1162e4be363deb10badc7a5e3d22ec7eb7f445fe5c76bea8c0e8b726840f.
//
// Solidity: event DelegatorStakeWithdrawal(address indexed delegator, address indexed transcoder, uint256 indexed amount)
func (_Staking *StakingFilterer) WatchDelegatorStakeWithdrawal(opts *bind.WatchOpts, sink chan<- *StakingDelegatorStakeWithdrawal, delegator []common.Address, transcoder []common.Address, amount []*big.Int) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var transcoderRule []interface{}
	for _, transcoderItem := range transcoder {
		transcoderRule = append(transcoderRule, transcoderItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "DelegatorStakeWithdrawal", delegatorRule, transcoderRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingDelegatorStakeWithdrawal)
				if err := _Staking.contract.UnpackLog(event, "DelegatorStakeWithdrawal", log); err != nil {
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

// ParseDelegatorStakeWithdrawal is a log parse operation binding the contract event 0xc1df1162e4be363deb10badc7a5e3d22ec7eb7f445fe5c76bea8c0e8b726840f.
//
// Solidity: event DelegatorStakeWithdrawal(address indexed delegator, address indexed transcoder, uint256 indexed amount)
func (_Staking *StakingFilterer) ParseDelegatorStakeWithdrawal(log types.Log) (*StakingDelegatorStakeWithdrawal, error) {
	event := new(StakingDelegatorStakeWithdrawal)
	if err := _Staking.contract.UnpackLog(event, "DelegatorStakeWithdrawal", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakingDelegatorStakeWithdrawalPropositionIterator is returned from FilterDelegatorStakeWithdrawalProposition and is used to iterate over the raw logs and unpacked data for DelegatorStakeWithdrawalProposition events raised by the Staking contract.
type StakingDelegatorStakeWithdrawalPropositionIterator struct {
	Event *StakingDelegatorStakeWithdrawalProposition // Event containing the contract specifics and raw log

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
func (it *StakingDelegatorStakeWithdrawalPropositionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingDelegatorStakeWithdrawalProposition)
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
		it.Event = new(StakingDelegatorStakeWithdrawalProposition)
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
func (it *StakingDelegatorStakeWithdrawalPropositionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingDelegatorStakeWithdrawalPropositionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingDelegatorStakeWithdrawalProposition represents a DelegatorStakeWithdrawalProposition event raised by the Staking contract.
type StakingDelegatorStakeWithdrawalProposition struct {
	Delegator  common.Address
	Transcoder common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterDelegatorStakeWithdrawalProposition is a free log retrieval operation binding the contract event 0xe81331619774d34a5a3cc97ec5040d55a8f59a8b8f068f0488d28c99d10f5502.
//
// Solidity: event DelegatorStakeWithdrawalProposition(address indexed delegator, address indexed transcoder)
func (_Staking *StakingFilterer) FilterDelegatorStakeWithdrawalProposition(opts *bind.FilterOpts, delegator []common.Address, transcoder []common.Address) (*StakingDelegatorStakeWithdrawalPropositionIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var transcoderRule []interface{}
	for _, transcoderItem := range transcoder {
		transcoderRule = append(transcoderRule, transcoderItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "DelegatorStakeWithdrawalProposition", delegatorRule, transcoderRule)
	if err != nil {
		return nil, err
	}
	return &StakingDelegatorStakeWithdrawalPropositionIterator{contract: _Staking.contract, event: "DelegatorStakeWithdrawalProposition", logs: logs, sub: sub}, nil
}

// WatchDelegatorStakeWithdrawalProposition is a free log subscription operation binding the contract event 0xe81331619774d34a5a3cc97ec5040d55a8f59a8b8f068f0488d28c99d10f5502.
//
// Solidity: event DelegatorStakeWithdrawalProposition(address indexed delegator, address indexed transcoder)
func (_Staking *StakingFilterer) WatchDelegatorStakeWithdrawalProposition(opts *bind.WatchOpts, sink chan<- *StakingDelegatorStakeWithdrawalProposition, delegator []common.Address, transcoder []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var transcoderRule []interface{}
	for _, transcoderItem := range transcoder {
		transcoderRule = append(transcoderRule, transcoderItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "DelegatorStakeWithdrawalProposition", delegatorRule, transcoderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingDelegatorStakeWithdrawalProposition)
				if err := _Staking.contract.UnpackLog(event, "DelegatorStakeWithdrawalProposition", log); err != nil {
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

// ParseDelegatorStakeWithdrawalProposition is a log parse operation binding the contract event 0xe81331619774d34a5a3cc97ec5040d55a8f59a8b8f068f0488d28c99d10f5502.
//
// Solidity: event DelegatorStakeWithdrawalProposition(address indexed delegator, address indexed transcoder)
func (_Staking *StakingFilterer) ParseDelegatorStakeWithdrawalProposition(log types.Log) (*StakingDelegatorStakeWithdrawalProposition, error) {
	event := new(StakingDelegatorStakeWithdrawalProposition)
	if err := _Staking.contract.UnpackLog(event, "DelegatorStakeWithdrawalProposition", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakingOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Staking contract.
type StakingOwnershipTransferredIterator struct {
	Event *StakingOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *StakingOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingOwnershipTransferred)
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
		it.Event = new(StakingOwnershipTransferred)
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
func (it *StakingOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingOwnershipTransferred represents a OwnershipTransferred event raised by the Staking contract.
type StakingOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Staking *StakingFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*StakingOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &StakingOwnershipTransferredIterator{contract: _Staking.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Staking *StakingFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *StakingOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingOwnershipTransferred)
				if err := _Staking.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Staking *StakingFilterer) ParseOwnershipTransferred(log types.Log) (*StakingOwnershipTransferred, error) {
	event := new(StakingOwnershipTransferred)
	if err := _Staking.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakingSlashedIterator is returned from FilterSlashed and is used to iterate over the raw logs and unpacked data for Slashed events raised by the Staking contract.
type StakingSlashedIterator struct {
	Event *StakingSlashed // Event containing the contract specifics and raw log

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
func (it *StakingSlashedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingSlashed)
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
		it.Event = new(StakingSlashed)
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
func (it *StakingSlashedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingSlashedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingSlashed represents a Slashed event raised by the Staking contract.
type StakingSlashed struct {
	Transcoder common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSlashed is a free log retrieval operation binding the contract event 0x975ad74f3e1160b20a620ec57c636ffa56c6b6679a0e2fb5689b4f95e65a6946.
//
// Solidity: event Slashed(address indexed transcoder)
func (_Staking *StakingFilterer) FilterSlashed(opts *bind.FilterOpts, transcoder []common.Address) (*StakingSlashedIterator, error) {

	var transcoderRule []interface{}
	for _, transcoderItem := range transcoder {
		transcoderRule = append(transcoderRule, transcoderItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Slashed", transcoderRule)
	if err != nil {
		return nil, err
	}
	return &StakingSlashedIterator{contract: _Staking.contract, event: "Slashed", logs: logs, sub: sub}, nil
}

// WatchSlashed is a free log subscription operation binding the contract event 0x975ad74f3e1160b20a620ec57c636ffa56c6b6679a0e2fb5689b4f95e65a6946.
//
// Solidity: event Slashed(address indexed transcoder)
func (_Staking *StakingFilterer) WatchSlashed(opts *bind.WatchOpts, sink chan<- *StakingSlashed, transcoder []common.Address) (event.Subscription, error) {

	var transcoderRule []interface{}
	for _, transcoderItem := range transcoder {
		transcoderRule = append(transcoderRule, transcoderItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Slashed", transcoderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingSlashed)
				if err := _Staking.contract.UnpackLog(event, "Slashed", log); err != nil {
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

// ParseSlashed is a log parse operation binding the contract event 0x975ad74f3e1160b20a620ec57c636ffa56c6b6679a0e2fb5689b4f95e65a6946.
//
// Solidity: event Slashed(address indexed transcoder)
func (_Staking *StakingFilterer) ParseSlashed(log types.Log) (*StakingSlashed, error) {
	event := new(StakingSlashed)
	if err := _Staking.contract.UnpackLog(event, "Slashed", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakingStakedIterator is returned from FilterStaked and is used to iterate over the raw logs and unpacked data for Staked events raised by the Staking contract.
type StakingStakedIterator struct {
	Event *StakingStaked // Event containing the contract specifics and raw log

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
func (it *StakingStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingStaked)
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
		it.Event = new(StakingStaked)
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
func (it *StakingStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingStaked represents a Staked event raised by the Staking contract.
type StakingStaked struct {
	Transcoder common.Address
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterStaked is a free log retrieval operation binding the contract event 0x9e71bc8eea02a63969f509818f2dafb9254532904319f9dbda79b67bd34a5f3d.
//
// Solidity: event Staked(address indexed transcoder, uint256 indexed amount)
func (_Staking *StakingFilterer) FilterStaked(opts *bind.FilterOpts, transcoder []common.Address, amount []*big.Int) (*StakingStakedIterator, error) {

	var transcoderRule []interface{}
	for _, transcoderItem := range transcoder {
		transcoderRule = append(transcoderRule, transcoderItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "Staked", transcoderRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &StakingStakedIterator{contract: _Staking.contract, event: "Staked", logs: logs, sub: sub}, nil
}

// WatchStaked is a free log subscription operation binding the contract event 0x9e71bc8eea02a63969f509818f2dafb9254532904319f9dbda79b67bd34a5f3d.
//
// Solidity: event Staked(address indexed transcoder, uint256 indexed amount)
func (_Staking *StakingFilterer) WatchStaked(opts *bind.WatchOpts, sink chan<- *StakingStaked, transcoder []common.Address, amount []*big.Int) (event.Subscription, error) {

	var transcoderRule []interface{}
	for _, transcoderItem := range transcoder {
		transcoderRule = append(transcoderRule, transcoderItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "Staked", transcoderRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingStaked)
				if err := _Staking.contract.UnpackLog(event, "Staked", log); err != nil {
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

// ParseStaked is a log parse operation binding the contract event 0x9e71bc8eea02a63969f509818f2dafb9254532904319f9dbda79b67bd34a5f3d.
//
// Solidity: event Staked(address indexed transcoder, uint256 indexed amount)
func (_Staking *StakingFilterer) ParseStaked(log types.Log) (*StakingStaked, error) {
	event := new(StakingStaked)
	if err := _Staking.contract.UnpackLog(event, "Staked", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakingTranscoderRegisteredIterator is returned from FilterTranscoderRegistered and is used to iterate over the raw logs and unpacked data for TranscoderRegistered events raised by the Staking contract.
type StakingTranscoderRegisteredIterator struct {
	Event *StakingTranscoderRegistered // Event containing the contract specifics and raw log

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
func (it *StakingTranscoderRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingTranscoderRegistered)
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
		it.Event = new(StakingTranscoderRegistered)
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
func (it *StakingTranscoderRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingTranscoderRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingTranscoderRegistered represents a TranscoderRegistered event raised by the Staking contract.
type StakingTranscoderRegistered struct {
	Transcoder common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTranscoderRegistered is a free log retrieval operation binding the contract event 0x6fbcf0f12b438f90175bebf725f86a4a74d12525d5d2c144a68e400696bce58b.
//
// Solidity: event TranscoderRegistered(address indexed transcoder)
func (_Staking *StakingFilterer) FilterTranscoderRegistered(opts *bind.FilterOpts, transcoder []common.Address) (*StakingTranscoderRegisteredIterator, error) {

	var transcoderRule []interface{}
	for _, transcoderItem := range transcoder {
		transcoderRule = append(transcoderRule, transcoderItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "TranscoderRegistered", transcoderRule)
	if err != nil {
		return nil, err
	}
	return &StakingTranscoderRegisteredIterator{contract: _Staking.contract, event: "TranscoderRegistered", logs: logs, sub: sub}, nil
}

// WatchTranscoderRegistered is a free log subscription operation binding the contract event 0x6fbcf0f12b438f90175bebf725f86a4a74d12525d5d2c144a68e400696bce58b.
//
// Solidity: event TranscoderRegistered(address indexed transcoder)
func (_Staking *StakingFilterer) WatchTranscoderRegistered(opts *bind.WatchOpts, sink chan<- *StakingTranscoderRegistered, transcoder []common.Address) (event.Subscription, error) {

	var transcoderRule []interface{}
	for _, transcoderItem := range transcoder {
		transcoderRule = append(transcoderRule, transcoderItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "TranscoderRegistered", transcoderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingTranscoderRegistered)
				if err := _Staking.contract.UnpackLog(event, "TranscoderRegistered", log); err != nil {
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

// ParseTranscoderRegistered is a log parse operation binding the contract event 0x6fbcf0f12b438f90175bebf725f86a4a74d12525d5d2c144a68e400696bce58b.
//
// Solidity: event TranscoderRegistered(address indexed transcoder)
func (_Staking *StakingFilterer) ParseTranscoderRegistered(log types.Log) (*StakingTranscoderRegistered, error) {
	event := new(StakingTranscoderRegistered)
	if err := _Staking.contract.UnpackLog(event, "TranscoderRegistered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakingTranscoderStakeWithdrawalIterator is returned from FilterTranscoderStakeWithdrawal and is used to iterate over the raw logs and unpacked data for TranscoderStakeWithdrawal events raised by the Staking contract.
type StakingTranscoderStakeWithdrawalIterator struct {
	Event *StakingTranscoderStakeWithdrawal // Event containing the contract specifics and raw log

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
func (it *StakingTranscoderStakeWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingTranscoderStakeWithdrawal)
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
		it.Event = new(StakingTranscoderStakeWithdrawal)
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
func (it *StakingTranscoderStakeWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingTranscoderStakeWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingTranscoderStakeWithdrawal represents a TranscoderStakeWithdrawal event raised by the Staking contract.
type StakingTranscoderStakeWithdrawal struct {
	Transcoder common.Address
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTranscoderStakeWithdrawal is a free log retrieval operation binding the contract event 0x1fb30b4ae3be530c8ba1440cab77c7adf646d5dea2d15b33b3f1127a344b09c6.
//
// Solidity: event TranscoderStakeWithdrawal(address indexed transcoder, uint256 indexed amount)
func (_Staking *StakingFilterer) FilterTranscoderStakeWithdrawal(opts *bind.FilterOpts, transcoder []common.Address, amount []*big.Int) (*StakingTranscoderStakeWithdrawalIterator, error) {

	var transcoderRule []interface{}
	for _, transcoderItem := range transcoder {
		transcoderRule = append(transcoderRule, transcoderItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "TranscoderStakeWithdrawal", transcoderRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &StakingTranscoderStakeWithdrawalIterator{contract: _Staking.contract, event: "TranscoderStakeWithdrawal", logs: logs, sub: sub}, nil
}

// WatchTranscoderStakeWithdrawal is a free log subscription operation binding the contract event 0x1fb30b4ae3be530c8ba1440cab77c7adf646d5dea2d15b33b3f1127a344b09c6.
//
// Solidity: event TranscoderStakeWithdrawal(address indexed transcoder, uint256 indexed amount)
func (_Staking *StakingFilterer) WatchTranscoderStakeWithdrawal(opts *bind.WatchOpts, sink chan<- *StakingTranscoderStakeWithdrawal, transcoder []common.Address, amount []*big.Int) (event.Subscription, error) {

	var transcoderRule []interface{}
	for _, transcoderItem := range transcoder {
		transcoderRule = append(transcoderRule, transcoderItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "TranscoderStakeWithdrawal", transcoderRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingTranscoderStakeWithdrawal)
				if err := _Staking.contract.UnpackLog(event, "TranscoderStakeWithdrawal", log); err != nil {
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

// ParseTranscoderStakeWithdrawal is a log parse operation binding the contract event 0x1fb30b4ae3be530c8ba1440cab77c7adf646d5dea2d15b33b3f1127a344b09c6.
//
// Solidity: event TranscoderStakeWithdrawal(address indexed transcoder, uint256 indexed amount)
func (_Staking *StakingFilterer) ParseTranscoderStakeWithdrawal(log types.Log) (*StakingTranscoderStakeWithdrawal, error) {
	event := new(StakingTranscoderStakeWithdrawal)
	if err := _Staking.contract.UnpackLog(event, "TranscoderStakeWithdrawal", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakingTranscoderStakeWithdrawalPropositionIterator is returned from FilterTranscoderStakeWithdrawalProposition and is used to iterate over the raw logs and unpacked data for TranscoderStakeWithdrawalProposition events raised by the Staking contract.
type StakingTranscoderStakeWithdrawalPropositionIterator struct {
	Event *StakingTranscoderStakeWithdrawalProposition // Event containing the contract specifics and raw log

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
func (it *StakingTranscoderStakeWithdrawalPropositionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingTranscoderStakeWithdrawalProposition)
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
		it.Event = new(StakingTranscoderStakeWithdrawalProposition)
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
func (it *StakingTranscoderStakeWithdrawalPropositionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingTranscoderStakeWithdrawalPropositionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingTranscoderStakeWithdrawalProposition represents a TranscoderStakeWithdrawalProposition event raised by the Staking contract.
type StakingTranscoderStakeWithdrawalProposition struct {
	Transcoder common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTranscoderStakeWithdrawalProposition is a free log retrieval operation binding the contract event 0x8ad5f6f9fb6f7ed1ec2f8f13c7ccff6e4626e35e8b2ded7771b831283625324f.
//
// Solidity: event TranscoderStakeWithdrawalProposition(address indexed transcoder)
func (_Staking *StakingFilterer) FilterTranscoderStakeWithdrawalProposition(opts *bind.FilterOpts, transcoder []common.Address) (*StakingTranscoderStakeWithdrawalPropositionIterator, error) {

	var transcoderRule []interface{}
	for _, transcoderItem := range transcoder {
		transcoderRule = append(transcoderRule, transcoderItem)
	}

	logs, sub, err := _Staking.contract.FilterLogs(opts, "TranscoderStakeWithdrawalProposition", transcoderRule)
	if err != nil {
		return nil, err
	}
	return &StakingTranscoderStakeWithdrawalPropositionIterator{contract: _Staking.contract, event: "TranscoderStakeWithdrawalProposition", logs: logs, sub: sub}, nil
}

// WatchTranscoderStakeWithdrawalProposition is a free log subscription operation binding the contract event 0x8ad5f6f9fb6f7ed1ec2f8f13c7ccff6e4626e35e8b2ded7771b831283625324f.
//
// Solidity: event TranscoderStakeWithdrawalProposition(address indexed transcoder)
func (_Staking *StakingFilterer) WatchTranscoderStakeWithdrawalProposition(opts *bind.WatchOpts, sink chan<- *StakingTranscoderStakeWithdrawalProposition, transcoder []common.Address) (event.Subscription, error) {

	var transcoderRule []interface{}
	for _, transcoderItem := range transcoder {
		transcoderRule = append(transcoderRule, transcoderItem)
	}

	logs, sub, err := _Staking.contract.WatchLogs(opts, "TranscoderStakeWithdrawalProposition", transcoderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingTranscoderStakeWithdrawalProposition)
				if err := _Staking.contract.UnpackLog(event, "TranscoderStakeWithdrawalProposition", log); err != nil {
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

// ParseTranscoderStakeWithdrawalProposition is a log parse operation binding the contract event 0x8ad5f6f9fb6f7ed1ec2f8f13c7ccff6e4626e35e8b2ded7771b831283625324f.
//
// Solidity: event TranscoderStakeWithdrawalProposition(address indexed transcoder)
func (_Staking *StakingFilterer) ParseTranscoderStakeWithdrawalProposition(log types.Log) (*StakingTranscoderStakeWithdrawalProposition, error) {
	event := new(StakingTranscoderStakeWithdrawalProposition)
	if err := _Staking.contract.UnpackLog(event, "TranscoderStakeWithdrawalProposition", log); err != nil {
		return nil, err
	}
	return event, nil
}
