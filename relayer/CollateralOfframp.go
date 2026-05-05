// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package relayer

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

// CollateralOfframpMetaData contains all meta data concerning the CollateralOfframp contract.
var CollateralOfframpMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_collateralToken\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AlreadyInitialized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpiredDeadline\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAsset\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidNonce\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewOwnerIsZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoHandoverRequest\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"pendingOwner\",\"type\":\"address\"}],\"name\":\"OwnershipHandoverCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"pendingOwner\",\"type\":\"address\"}],\"name\":\"OwnershipHandoverRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roles\",\"type\":\"uint256\"}],\"name\":\"RolesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COLLATERAL_TOKEN\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"}],\"name\":\"addAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"cancelOwnershipHandover\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"pendingOwner\",\"type\":\"address\"}],\"name\":\"completeOwnershipHandover\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"roles\",\"type\":\"uint256\"}],\"name\":\"grantRoles\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"roles\",\"type\":\"uint256\"}],\"name\":\"hasAllRoles\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"roles\",\"type\":\"uint256\"}],\"name\":\"hasAnyRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"result\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"pendingOwner\",\"type\":\"address\"}],\"name\":\"ownershipHandoverExpiresAt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"result\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_asset\",\"type\":\"address\"}],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"}],\"name\":\"removeAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roles\",\"type\":\"uint256\"}],\"name\":\"renounceRoles\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestOwnershipHandover\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"roles\",\"type\":\"uint256\"}],\"name\":\"revokeRoles\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"rolesOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"roles\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_asset\",\"type\":\"address\"}],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_asset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"unwrap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561000f575f5ffd5b50604051610bd7380380610bd783398101604081905261002e91610114565b6001600160a01b03811660805261004483610057565b61004f826001610092565b505050610154565b6001600160a01b0316638b78c6d819819055805f7f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08180a350565b61009e828260016100a2565b5050565b638b78c6d8600c52825f526020600c208054838117836100c3575080841681185b80835580600c5160601c7f715ad5ce61fc9595c7b415289d59cf203f23a94fa06f04af7e489a0a76e1fe265f5fa3505050505050565b80516001600160a01b038116811461010f575f5ffd5b919050565b5f5f5f60608486031215610126575f5ffd5b61012f846100f9565b925061013d602085016100f9565b915061014b604085016100f9565b90509250925092565b608051610a5d61017a5f395f81816103c001528181610674015261071b0152610a5d5ff3fe608060405260043610610157575f3560e01c806357b001f9116100bb5780638da5cb5b11610071578063f2fde38b11610057578063f2fde38b1461039c578063f5f1f1a7146103af578063fee81cf4146103e2575f5ffd5b80638da5cb5b14610335578063f04e283e14610389575f5ffd5b8063715018a6116100a1578063715018a6146102ef57806376a67a51146102f75780638cc7104f14610316575f5ffd5b806357b001f9146102b157806370480275146102d0575f5ffd5b80632de94807116101105780634a4ee7b1116100f65780634a4ee7b114610261578063514e62fc1461027457806354d1f13d146102a9575f5ffd5b80632de94807146101f45780632e48152c14610233575f5ffd5b80631c10893f116101405780631c10893f1461018f5780631cd64df4146101a257806325692962146101ec575f5ffd5b80631785f53c1461015b578063183a4f6e1461017c575b5f5ffd5b348015610166575f5ffd5b5061017a61017536600461098e565b610413565b005b61017a61018a3660046109ae565b61042d565b61017a61019d3660046109c5565b61043a565b3480156101ad575f5ffd5b506101d76101bc3660046109c5565b638b78c6d8600c9081525f9290925260209091205481161490565b60405190151581526020015b60405180910390f35b61017a61044c565b3480156101ff575f5ffd5b5061022561020e36600461098e565b638b78c6d8600c9081525f91909152602090205490565b6040519081526020016101e3565b34801561023e575f5ffd5b506101d761024d36600461098e565b5f6020819052908152604090205460ff1681565b61017a61026f3660046109c5565b610499565b34801561027f575f5ffd5b506101d761028e3660046109c5565b638b78c6d8600c9081525f9290925260209091205416151590565b61017a6104ab565b3480156102bc575f5ffd5b5061017a6102cb36600461098e565b6104e4565b3480156102db575f5ffd5b5061017a6102ea36600461098e565b610563565b61017a610579565b348015610302575f5ffd5b5061017a61031136600461098e565b61058c565b348015610321575f5ffd5b5061017a6103303660046109ed565b61060e565b348015610340575f5ffd5b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffff74873927545b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101e3565b61017a61039736600461098e565b610778565b61017a6103aa36600461098e565b6107b2565b3480156103ba575f5ffd5b506103647f000000000000000000000000000000000000000000000000000000000000000081565b3480156103ed575f5ffd5b506102256103fc36600461098e565b63389a75e1600c9081525f91909152602090205490565b600161041e816107d8565b6104298260016107fc565b5050565b61043733826107fc565b50565b610442610807565b610429828261083c565b5f6202a30067ffffffffffffffff164201905063389a75e1600c52335f52806020600c2055337fdbf36a107da19e49527a7176a1babf963b4b0ff8cde35ee35d6cd8f1f9ac7e1d5f5fa250565b6104a1610807565b61042982826107fc565b63389a75e1600c52335f525f6020600c2055337ffa7b8eab7da67f412cc9575ed43464468f9bfbae89d1675917346ca6d8fe3c925f5fa2565b60016104ef816107d8565b73ffffffffffffffffffffffffffffffffffffffff82165f8181526020819052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055517f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa9190a25050565b600161056e816107d8565b61042982600161083c565b610581610807565b61058a5f610848565b565b6001610597816107d8565b73ffffffffffffffffffffffffffffffffffffffff82165f8181526020819052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055517f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2589190a25050565b73ffffffffffffffffffffffffffffffffffffffff83165f90815260208190526040902054839060ff161561066f576040517f49b8b3ac00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6106b37f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff8116903390856108ad565b6040517fd600875d00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff85811660048301528481166024830152604482018490525f6064830181905260a0608484015260a48301527f0000000000000000000000000000000000000000000000000000000000000000169063d600875d9060c4015f604051808303815f87803b15801561075c575f5ffd5b505af115801561076e573d5f5f3e3d5ffd5b5050505050505050565b610780610807565b63389a75e1600c52805f526020600c2080544211156107a657636f5e88185f526004601cfd5b5f905561043781610848565b6107ba610807565b8060601b6107cf57637448fbae5f526004601cfd5b61043781610848565b638b78c6d8600c52335f52806020600c205416610437576382b429005f526004601cfd5b61042982825f61090f565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffff7487392754331461058a576382b429005f526004601cfd5b6104298282600161090f565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffff74873927805473ffffffffffffffffffffffffffffffffffffffff9092169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a355565b60405181606052826040528360601b602c526f23b872dd000000000000000000000000600c5260205f6064601c5f895af18060015f51141661090157803d873b15171061090157637939f4245f526004601cfd5b505f60605260405250505050565b638b78c6d8600c52825f526020600c20805483811783610930575080841681185b80835580600c5160601c7f715ad5ce61fc9595c7b415289d59cf203f23a94fa06f04af7e489a0a76e1fe265f5fa3505050505050565b803573ffffffffffffffffffffffffffffffffffffffff81168114610989575f5ffd5b919050565b5f6020828403121561099e575f5ffd5b6109a782610966565b9392505050565b5f602082840312156109be575f5ffd5b5035919050565b5f5f604083850312156109d6575f5ffd5b6109df83610966565b946020939093013593505050565b5f5f5f606084860312156109ff575f5ffd5b610a0884610966565b9250610a1660208501610966565b92959294505050604091909101359056fea2646970667358221220a45f1e4f53db87fb26e9c6f9cc8a913514b29536571165af2cdd0c5c5f4968d564736f6c6343000822003300000000000000000000000047ebfac3353314c788b96cdcbf41daadfe03629c0000000000000000000000003dce0a29139a851da1dfca56af8e8a6440b4d952000000000000000000000000c011a7e12a19f7b1f670d46f03b03f3342e82dfb",
}

// CollateralOfframpABI is the input ABI used to generate the binding from.
// Deprecated: Use CollateralOfframpMetaData.ABI instead.
var CollateralOfframpABI = CollateralOfframpMetaData.ABI

// CollateralOfframpBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CollateralOfframpMetaData.Bin instead.
var CollateralOfframpBin = CollateralOfframpMetaData.Bin

// DeployCollateralOfframp deploys a new Ethereum contract, binding an instance of CollateralOfframp to it.
func DeployCollateralOfframp(auth *bind.TransactOpts, backend bind.ContractBackend, _owner common.Address, _admin common.Address, _collateralToken common.Address) (common.Address, *types.Transaction, *CollateralOfframp, error) {
	parsed, err := CollateralOfframpMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CollateralOfframpBin), backend, _owner, _admin, _collateralToken)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CollateralOfframp{CollateralOfframpCaller: CollateralOfframpCaller{contract: contract}, CollateralOfframpTransactor: CollateralOfframpTransactor{contract: contract}, CollateralOfframpFilterer: CollateralOfframpFilterer{contract: contract}}, nil
}

// CollateralOfframp is an auto generated Go binding around an Ethereum contract.
type CollateralOfframp struct {
	CollateralOfframpCaller     // Read-only binding to the contract
	CollateralOfframpTransactor // Write-only binding to the contract
	CollateralOfframpFilterer   // Log filterer for contract events
}

// CollateralOfframpCaller is an auto generated read-only Go binding around an Ethereum contract.
type CollateralOfframpCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CollateralOfframpTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CollateralOfframpTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CollateralOfframpFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CollateralOfframpFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CollateralOfframpSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CollateralOfframpSession struct {
	Contract     *CollateralOfframp // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// CollateralOfframpCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CollateralOfframpCallerSession struct {
	Contract *CollateralOfframpCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// CollateralOfframpTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CollateralOfframpTransactorSession struct {
	Contract     *CollateralOfframpTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// CollateralOfframpRaw is an auto generated low-level Go binding around an Ethereum contract.
type CollateralOfframpRaw struct {
	Contract *CollateralOfframp // Generic contract binding to access the raw methods on
}

// CollateralOfframpCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CollateralOfframpCallerRaw struct {
	Contract *CollateralOfframpCaller // Generic read-only contract binding to access the raw methods on
}

// CollateralOfframpTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CollateralOfframpTransactorRaw struct {
	Contract *CollateralOfframpTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCollateralOfframp creates a new instance of CollateralOfframp, bound to a specific deployed contract.
func NewCollateralOfframp(address common.Address, backend bind.ContractBackend) (*CollateralOfframp, error) {
	contract, err := bindCollateralOfframp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CollateralOfframp{CollateralOfframpCaller: CollateralOfframpCaller{contract: contract}, CollateralOfframpTransactor: CollateralOfframpTransactor{contract: contract}, CollateralOfframpFilterer: CollateralOfframpFilterer{contract: contract}}, nil
}

// NewCollateralOfframpCaller creates a new read-only instance of CollateralOfframp, bound to a specific deployed contract.
func NewCollateralOfframpCaller(address common.Address, caller bind.ContractCaller) (*CollateralOfframpCaller, error) {
	contract, err := bindCollateralOfframp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CollateralOfframpCaller{contract: contract}, nil
}

// NewCollateralOfframpTransactor creates a new write-only instance of CollateralOfframp, bound to a specific deployed contract.
func NewCollateralOfframpTransactor(address common.Address, transactor bind.ContractTransactor) (*CollateralOfframpTransactor, error) {
	contract, err := bindCollateralOfframp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CollateralOfframpTransactor{contract: contract}, nil
}

// NewCollateralOfframpFilterer creates a new log filterer instance of CollateralOfframp, bound to a specific deployed contract.
func NewCollateralOfframpFilterer(address common.Address, filterer bind.ContractFilterer) (*CollateralOfframpFilterer, error) {
	contract, err := bindCollateralOfframp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CollateralOfframpFilterer{contract: contract}, nil
}

// bindCollateralOfframp binds a generic wrapper to an already deployed contract.
func bindCollateralOfframp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CollateralOfframpMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CollateralOfframp *CollateralOfframpRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CollateralOfframp.Contract.CollateralOfframpCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CollateralOfframp *CollateralOfframpRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CollateralOfframp.Contract.CollateralOfframpTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CollateralOfframp *CollateralOfframpRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CollateralOfframp.Contract.CollateralOfframpTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CollateralOfframp *CollateralOfframpCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CollateralOfframp.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CollateralOfframp *CollateralOfframpTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CollateralOfframp.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CollateralOfframp *CollateralOfframpTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CollateralOfframp.Contract.contract.Transact(opts, method, params...)
}

// COLLATERALTOKEN is a free data retrieval call binding the contract method 0xf5f1f1a7.
//
// Solidity: function COLLATERAL_TOKEN() view returns(address)
func (_CollateralOfframp *CollateralOfframpCaller) COLLATERALTOKEN(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CollateralOfframp.contract.Call(opts, &out, "COLLATERAL_TOKEN")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// COLLATERALTOKEN is a free data retrieval call binding the contract method 0xf5f1f1a7.
//
// Solidity: function COLLATERAL_TOKEN() view returns(address)
func (_CollateralOfframp *CollateralOfframpSession) COLLATERALTOKEN() (common.Address, error) {
	return _CollateralOfframp.Contract.COLLATERALTOKEN(&_CollateralOfframp.CallOpts)
}

// COLLATERALTOKEN is a free data retrieval call binding the contract method 0xf5f1f1a7.
//
// Solidity: function COLLATERAL_TOKEN() view returns(address)
func (_CollateralOfframp *CollateralOfframpCallerSession) COLLATERALTOKEN() (common.Address, error) {
	return _CollateralOfframp.Contract.COLLATERALTOKEN(&_CollateralOfframp.CallOpts)
}

// HasAllRoles is a free data retrieval call binding the contract method 0x1cd64df4.
//
// Solidity: function hasAllRoles(address user, uint256 roles) view returns(bool)
func (_CollateralOfframp *CollateralOfframpCaller) HasAllRoles(opts *bind.CallOpts, user common.Address, roles *big.Int) (bool, error) {
	var out []interface{}
	err := _CollateralOfframp.contract.Call(opts, &out, "hasAllRoles", user, roles)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasAllRoles is a free data retrieval call binding the contract method 0x1cd64df4.
//
// Solidity: function hasAllRoles(address user, uint256 roles) view returns(bool)
func (_CollateralOfframp *CollateralOfframpSession) HasAllRoles(user common.Address, roles *big.Int) (bool, error) {
	return _CollateralOfframp.Contract.HasAllRoles(&_CollateralOfframp.CallOpts, user, roles)
}

// HasAllRoles is a free data retrieval call binding the contract method 0x1cd64df4.
//
// Solidity: function hasAllRoles(address user, uint256 roles) view returns(bool)
func (_CollateralOfframp *CollateralOfframpCallerSession) HasAllRoles(user common.Address, roles *big.Int) (bool, error) {
	return _CollateralOfframp.Contract.HasAllRoles(&_CollateralOfframp.CallOpts, user, roles)
}

// HasAnyRole is a free data retrieval call binding the contract method 0x514e62fc.
//
// Solidity: function hasAnyRole(address user, uint256 roles) view returns(bool)
func (_CollateralOfframp *CollateralOfframpCaller) HasAnyRole(opts *bind.CallOpts, user common.Address, roles *big.Int) (bool, error) {
	var out []interface{}
	err := _CollateralOfframp.contract.Call(opts, &out, "hasAnyRole", user, roles)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasAnyRole is a free data retrieval call binding the contract method 0x514e62fc.
//
// Solidity: function hasAnyRole(address user, uint256 roles) view returns(bool)
func (_CollateralOfframp *CollateralOfframpSession) HasAnyRole(user common.Address, roles *big.Int) (bool, error) {
	return _CollateralOfframp.Contract.HasAnyRole(&_CollateralOfframp.CallOpts, user, roles)
}

// HasAnyRole is a free data retrieval call binding the contract method 0x514e62fc.
//
// Solidity: function hasAnyRole(address user, uint256 roles) view returns(bool)
func (_CollateralOfframp *CollateralOfframpCallerSession) HasAnyRole(user common.Address, roles *big.Int) (bool, error) {
	return _CollateralOfframp.Contract.HasAnyRole(&_CollateralOfframp.CallOpts, user, roles)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address result)
func (_CollateralOfframp *CollateralOfframpCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CollateralOfframp.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address result)
func (_CollateralOfframp *CollateralOfframpSession) Owner() (common.Address, error) {
	return _CollateralOfframp.Contract.Owner(&_CollateralOfframp.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address result)
func (_CollateralOfframp *CollateralOfframpCallerSession) Owner() (common.Address, error) {
	return _CollateralOfframp.Contract.Owner(&_CollateralOfframp.CallOpts)
}

// OwnershipHandoverExpiresAt is a free data retrieval call binding the contract method 0xfee81cf4.
//
// Solidity: function ownershipHandoverExpiresAt(address pendingOwner) view returns(uint256 result)
func (_CollateralOfframp *CollateralOfframpCaller) OwnershipHandoverExpiresAt(opts *bind.CallOpts, pendingOwner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CollateralOfframp.contract.Call(opts, &out, "ownershipHandoverExpiresAt", pendingOwner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OwnershipHandoverExpiresAt is a free data retrieval call binding the contract method 0xfee81cf4.
//
// Solidity: function ownershipHandoverExpiresAt(address pendingOwner) view returns(uint256 result)
func (_CollateralOfframp *CollateralOfframpSession) OwnershipHandoverExpiresAt(pendingOwner common.Address) (*big.Int, error) {
	return _CollateralOfframp.Contract.OwnershipHandoverExpiresAt(&_CollateralOfframp.CallOpts, pendingOwner)
}

// OwnershipHandoverExpiresAt is a free data retrieval call binding the contract method 0xfee81cf4.
//
// Solidity: function ownershipHandoverExpiresAt(address pendingOwner) view returns(uint256 result)
func (_CollateralOfframp *CollateralOfframpCallerSession) OwnershipHandoverExpiresAt(pendingOwner common.Address) (*big.Int, error) {
	return _CollateralOfframp.Contract.OwnershipHandoverExpiresAt(&_CollateralOfframp.CallOpts, pendingOwner)
}

// Paused is a free data retrieval call binding the contract method 0x2e48152c.
//
// Solidity: function paused(address ) view returns(bool)
func (_CollateralOfframp *CollateralOfframpCaller) Paused(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _CollateralOfframp.contract.Call(opts, &out, "paused", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x2e48152c.
//
// Solidity: function paused(address ) view returns(bool)
func (_CollateralOfframp *CollateralOfframpSession) Paused(arg0 common.Address) (bool, error) {
	return _CollateralOfframp.Contract.Paused(&_CollateralOfframp.CallOpts, arg0)
}

// Paused is a free data retrieval call binding the contract method 0x2e48152c.
//
// Solidity: function paused(address ) view returns(bool)
func (_CollateralOfframp *CollateralOfframpCallerSession) Paused(arg0 common.Address) (bool, error) {
	return _CollateralOfframp.Contract.Paused(&_CollateralOfframp.CallOpts, arg0)
}

// RolesOf is a free data retrieval call binding the contract method 0x2de94807.
//
// Solidity: function rolesOf(address user) view returns(uint256 roles)
func (_CollateralOfframp *CollateralOfframpCaller) RolesOf(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CollateralOfframp.contract.Call(opts, &out, "rolesOf", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RolesOf is a free data retrieval call binding the contract method 0x2de94807.
//
// Solidity: function rolesOf(address user) view returns(uint256 roles)
func (_CollateralOfframp *CollateralOfframpSession) RolesOf(user common.Address) (*big.Int, error) {
	return _CollateralOfframp.Contract.RolesOf(&_CollateralOfframp.CallOpts, user)
}

// RolesOf is a free data retrieval call binding the contract method 0x2de94807.
//
// Solidity: function rolesOf(address user) view returns(uint256 roles)
func (_CollateralOfframp *CollateralOfframpCallerSession) RolesOf(user common.Address) (*big.Int, error) {
	return _CollateralOfframp.Contract.RolesOf(&_CollateralOfframp.CallOpts, user)
}

// AddAdmin is a paid mutator transaction binding the contract method 0x70480275.
//
// Solidity: function addAdmin(address _admin) returns()
func (_CollateralOfframp *CollateralOfframpTransactor) AddAdmin(opts *bind.TransactOpts, _admin common.Address) (*types.Transaction, error) {
	return _CollateralOfframp.contract.Transact(opts, "addAdmin", _admin)
}

// AddAdmin is a paid mutator transaction binding the contract method 0x70480275.
//
// Solidity: function addAdmin(address _admin) returns()
func (_CollateralOfframp *CollateralOfframpSession) AddAdmin(_admin common.Address) (*types.Transaction, error) {
	return _CollateralOfframp.Contract.AddAdmin(&_CollateralOfframp.TransactOpts, _admin)
}

// AddAdmin is a paid mutator transaction binding the contract method 0x70480275.
//
// Solidity: function addAdmin(address _admin) returns()
func (_CollateralOfframp *CollateralOfframpTransactorSession) AddAdmin(_admin common.Address) (*types.Transaction, error) {
	return _CollateralOfframp.Contract.AddAdmin(&_CollateralOfframp.TransactOpts, _admin)
}

// CancelOwnershipHandover is a paid mutator transaction binding the contract method 0x54d1f13d.
//
// Solidity: function cancelOwnershipHandover() payable returns()
func (_CollateralOfframp *CollateralOfframpTransactor) CancelOwnershipHandover(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CollateralOfframp.contract.Transact(opts, "cancelOwnershipHandover")
}

// CancelOwnershipHandover is a paid mutator transaction binding the contract method 0x54d1f13d.
//
// Solidity: function cancelOwnershipHandover() payable returns()
func (_CollateralOfframp *CollateralOfframpSession) CancelOwnershipHandover() (*types.Transaction, error) {
	return _CollateralOfframp.Contract.CancelOwnershipHandover(&_CollateralOfframp.TransactOpts)
}

// CancelOwnershipHandover is a paid mutator transaction binding the contract method 0x54d1f13d.
//
// Solidity: function cancelOwnershipHandover() payable returns()
func (_CollateralOfframp *CollateralOfframpTransactorSession) CancelOwnershipHandover() (*types.Transaction, error) {
	return _CollateralOfframp.Contract.CancelOwnershipHandover(&_CollateralOfframp.TransactOpts)
}

// CompleteOwnershipHandover is a paid mutator transaction binding the contract method 0xf04e283e.
//
// Solidity: function completeOwnershipHandover(address pendingOwner) payable returns()
func (_CollateralOfframp *CollateralOfframpTransactor) CompleteOwnershipHandover(opts *bind.TransactOpts, pendingOwner common.Address) (*types.Transaction, error) {
	return _CollateralOfframp.contract.Transact(opts, "completeOwnershipHandover", pendingOwner)
}

// CompleteOwnershipHandover is a paid mutator transaction binding the contract method 0xf04e283e.
//
// Solidity: function completeOwnershipHandover(address pendingOwner) payable returns()
func (_CollateralOfframp *CollateralOfframpSession) CompleteOwnershipHandover(pendingOwner common.Address) (*types.Transaction, error) {
	return _CollateralOfframp.Contract.CompleteOwnershipHandover(&_CollateralOfframp.TransactOpts, pendingOwner)
}

// CompleteOwnershipHandover is a paid mutator transaction binding the contract method 0xf04e283e.
//
// Solidity: function completeOwnershipHandover(address pendingOwner) payable returns()
func (_CollateralOfframp *CollateralOfframpTransactorSession) CompleteOwnershipHandover(pendingOwner common.Address) (*types.Transaction, error) {
	return _CollateralOfframp.Contract.CompleteOwnershipHandover(&_CollateralOfframp.TransactOpts, pendingOwner)
}

// GrantRoles is a paid mutator transaction binding the contract method 0x1c10893f.
//
// Solidity: function grantRoles(address user, uint256 roles) payable returns()
func (_CollateralOfframp *CollateralOfframpTransactor) GrantRoles(opts *bind.TransactOpts, user common.Address, roles *big.Int) (*types.Transaction, error) {
	return _CollateralOfframp.contract.Transact(opts, "grantRoles", user, roles)
}

// GrantRoles is a paid mutator transaction binding the contract method 0x1c10893f.
//
// Solidity: function grantRoles(address user, uint256 roles) payable returns()
func (_CollateralOfframp *CollateralOfframpSession) GrantRoles(user common.Address, roles *big.Int) (*types.Transaction, error) {
	return _CollateralOfframp.Contract.GrantRoles(&_CollateralOfframp.TransactOpts, user, roles)
}

// GrantRoles is a paid mutator transaction binding the contract method 0x1c10893f.
//
// Solidity: function grantRoles(address user, uint256 roles) payable returns()
func (_CollateralOfframp *CollateralOfframpTransactorSession) GrantRoles(user common.Address, roles *big.Int) (*types.Transaction, error) {
	return _CollateralOfframp.Contract.GrantRoles(&_CollateralOfframp.TransactOpts, user, roles)
}

// Pause is a paid mutator transaction binding the contract method 0x76a67a51.
//
// Solidity: function pause(address _asset) returns()
func (_CollateralOfframp *CollateralOfframpTransactor) Pause(opts *bind.TransactOpts, _asset common.Address) (*types.Transaction, error) {
	return _CollateralOfframp.contract.Transact(opts, "pause", _asset)
}

// Pause is a paid mutator transaction binding the contract method 0x76a67a51.
//
// Solidity: function pause(address _asset) returns()
func (_CollateralOfframp *CollateralOfframpSession) Pause(_asset common.Address) (*types.Transaction, error) {
	return _CollateralOfframp.Contract.Pause(&_CollateralOfframp.TransactOpts, _asset)
}

// Pause is a paid mutator transaction binding the contract method 0x76a67a51.
//
// Solidity: function pause(address _asset) returns()
func (_CollateralOfframp *CollateralOfframpTransactorSession) Pause(_asset common.Address) (*types.Transaction, error) {
	return _CollateralOfframp.Contract.Pause(&_CollateralOfframp.TransactOpts, _asset)
}

// RemoveAdmin is a paid mutator transaction binding the contract method 0x1785f53c.
//
// Solidity: function removeAdmin(address _admin) returns()
func (_CollateralOfframp *CollateralOfframpTransactor) RemoveAdmin(opts *bind.TransactOpts, _admin common.Address) (*types.Transaction, error) {
	return _CollateralOfframp.contract.Transact(opts, "removeAdmin", _admin)
}

// RemoveAdmin is a paid mutator transaction binding the contract method 0x1785f53c.
//
// Solidity: function removeAdmin(address _admin) returns()
func (_CollateralOfframp *CollateralOfframpSession) RemoveAdmin(_admin common.Address) (*types.Transaction, error) {
	return _CollateralOfframp.Contract.RemoveAdmin(&_CollateralOfframp.TransactOpts, _admin)
}

// RemoveAdmin is a paid mutator transaction binding the contract method 0x1785f53c.
//
// Solidity: function removeAdmin(address _admin) returns()
func (_CollateralOfframp *CollateralOfframpTransactorSession) RemoveAdmin(_admin common.Address) (*types.Transaction, error) {
	return _CollateralOfframp.Contract.RemoveAdmin(&_CollateralOfframp.TransactOpts, _admin)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() payable returns()
func (_CollateralOfframp *CollateralOfframpTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CollateralOfframp.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() payable returns()
func (_CollateralOfframp *CollateralOfframpSession) RenounceOwnership() (*types.Transaction, error) {
	return _CollateralOfframp.Contract.RenounceOwnership(&_CollateralOfframp.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() payable returns()
func (_CollateralOfframp *CollateralOfframpTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _CollateralOfframp.Contract.RenounceOwnership(&_CollateralOfframp.TransactOpts)
}

// RenounceRoles is a paid mutator transaction binding the contract method 0x183a4f6e.
//
// Solidity: function renounceRoles(uint256 roles) payable returns()
func (_CollateralOfframp *CollateralOfframpTransactor) RenounceRoles(opts *bind.TransactOpts, roles *big.Int) (*types.Transaction, error) {
	return _CollateralOfframp.contract.Transact(opts, "renounceRoles", roles)
}

// RenounceRoles is a paid mutator transaction binding the contract method 0x183a4f6e.
//
// Solidity: function renounceRoles(uint256 roles) payable returns()
func (_CollateralOfframp *CollateralOfframpSession) RenounceRoles(roles *big.Int) (*types.Transaction, error) {
	return _CollateralOfframp.Contract.RenounceRoles(&_CollateralOfframp.TransactOpts, roles)
}

// RenounceRoles is a paid mutator transaction binding the contract method 0x183a4f6e.
//
// Solidity: function renounceRoles(uint256 roles) payable returns()
func (_CollateralOfframp *CollateralOfframpTransactorSession) RenounceRoles(roles *big.Int) (*types.Transaction, error) {
	return _CollateralOfframp.Contract.RenounceRoles(&_CollateralOfframp.TransactOpts, roles)
}

// RequestOwnershipHandover is a paid mutator transaction binding the contract method 0x25692962.
//
// Solidity: function requestOwnershipHandover() payable returns()
func (_CollateralOfframp *CollateralOfframpTransactor) RequestOwnershipHandover(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CollateralOfframp.contract.Transact(opts, "requestOwnershipHandover")
}

// RequestOwnershipHandover is a paid mutator transaction binding the contract method 0x25692962.
//
// Solidity: function requestOwnershipHandover() payable returns()
func (_CollateralOfframp *CollateralOfframpSession) RequestOwnershipHandover() (*types.Transaction, error) {
	return _CollateralOfframp.Contract.RequestOwnershipHandover(&_CollateralOfframp.TransactOpts)
}

// RequestOwnershipHandover is a paid mutator transaction binding the contract method 0x25692962.
//
// Solidity: function requestOwnershipHandover() payable returns()
func (_CollateralOfframp *CollateralOfframpTransactorSession) RequestOwnershipHandover() (*types.Transaction, error) {
	return _CollateralOfframp.Contract.RequestOwnershipHandover(&_CollateralOfframp.TransactOpts)
}

// RevokeRoles is a paid mutator transaction binding the contract method 0x4a4ee7b1.
//
// Solidity: function revokeRoles(address user, uint256 roles) payable returns()
func (_CollateralOfframp *CollateralOfframpTransactor) RevokeRoles(opts *bind.TransactOpts, user common.Address, roles *big.Int) (*types.Transaction, error) {
	return _CollateralOfframp.contract.Transact(opts, "revokeRoles", user, roles)
}

// RevokeRoles is a paid mutator transaction binding the contract method 0x4a4ee7b1.
//
// Solidity: function revokeRoles(address user, uint256 roles) payable returns()
func (_CollateralOfframp *CollateralOfframpSession) RevokeRoles(user common.Address, roles *big.Int) (*types.Transaction, error) {
	return _CollateralOfframp.Contract.RevokeRoles(&_CollateralOfframp.TransactOpts, user, roles)
}

// RevokeRoles is a paid mutator transaction binding the contract method 0x4a4ee7b1.
//
// Solidity: function revokeRoles(address user, uint256 roles) payable returns()
func (_CollateralOfframp *CollateralOfframpTransactorSession) RevokeRoles(user common.Address, roles *big.Int) (*types.Transaction, error) {
	return _CollateralOfframp.Contract.RevokeRoles(&_CollateralOfframp.TransactOpts, user, roles)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) payable returns()
func (_CollateralOfframp *CollateralOfframpTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _CollateralOfframp.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) payable returns()
func (_CollateralOfframp *CollateralOfframpSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CollateralOfframp.Contract.TransferOwnership(&_CollateralOfframp.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) payable returns()
func (_CollateralOfframp *CollateralOfframpTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CollateralOfframp.Contract.TransferOwnership(&_CollateralOfframp.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x57b001f9.
//
// Solidity: function unpause(address _asset) returns()
func (_CollateralOfframp *CollateralOfframpTransactor) Unpause(opts *bind.TransactOpts, _asset common.Address) (*types.Transaction, error) {
	return _CollateralOfframp.contract.Transact(opts, "unpause", _asset)
}

// Unpause is a paid mutator transaction binding the contract method 0x57b001f9.
//
// Solidity: function unpause(address _asset) returns()
func (_CollateralOfframp *CollateralOfframpSession) Unpause(_asset common.Address) (*types.Transaction, error) {
	return _CollateralOfframp.Contract.Unpause(&_CollateralOfframp.TransactOpts, _asset)
}

// Unpause is a paid mutator transaction binding the contract method 0x57b001f9.
//
// Solidity: function unpause(address _asset) returns()
func (_CollateralOfframp *CollateralOfframpTransactorSession) Unpause(_asset common.Address) (*types.Transaction, error) {
	return _CollateralOfframp.Contract.Unpause(&_CollateralOfframp.TransactOpts, _asset)
}

// Unwrap is a paid mutator transaction binding the contract method 0x8cc7104f.
//
// Solidity: function unwrap(address _asset, address _to, uint256 _amount) returns()
func (_CollateralOfframp *CollateralOfframpTransactor) Unwrap(opts *bind.TransactOpts, _asset common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CollateralOfframp.contract.Transact(opts, "unwrap", _asset, _to, _amount)
}

// Unwrap is a paid mutator transaction binding the contract method 0x8cc7104f.
//
// Solidity: function unwrap(address _asset, address _to, uint256 _amount) returns()
func (_CollateralOfframp *CollateralOfframpSession) Unwrap(_asset common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CollateralOfframp.Contract.Unwrap(&_CollateralOfframp.TransactOpts, _asset, _to, _amount)
}

// Unwrap is a paid mutator transaction binding the contract method 0x8cc7104f.
//
// Solidity: function unwrap(address _asset, address _to, uint256 _amount) returns()
func (_CollateralOfframp *CollateralOfframpTransactorSession) Unwrap(_asset common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CollateralOfframp.Contract.Unwrap(&_CollateralOfframp.TransactOpts, _asset, _to, _amount)
}

// CollateralOfframpOwnershipHandoverCanceledIterator is returned from FilterOwnershipHandoverCanceled and is used to iterate over the raw logs and unpacked data for OwnershipHandoverCanceled events raised by the CollateralOfframp contract.
type CollateralOfframpOwnershipHandoverCanceledIterator struct {
	Event *CollateralOfframpOwnershipHandoverCanceled // Event containing the contract specifics and raw log

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
func (it *CollateralOfframpOwnershipHandoverCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollateralOfframpOwnershipHandoverCanceled)
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
		it.Event = new(CollateralOfframpOwnershipHandoverCanceled)
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
func (it *CollateralOfframpOwnershipHandoverCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollateralOfframpOwnershipHandoverCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollateralOfframpOwnershipHandoverCanceled represents a OwnershipHandoverCanceled event raised by the CollateralOfframp contract.
type CollateralOfframpOwnershipHandoverCanceled struct {
	PendingOwner common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterOwnershipHandoverCanceled is a free log retrieval operation binding the contract event 0xfa7b8eab7da67f412cc9575ed43464468f9bfbae89d1675917346ca6d8fe3c92.
//
// Solidity: event OwnershipHandoverCanceled(address indexed pendingOwner)
func (_CollateralOfframp *CollateralOfframpFilterer) FilterOwnershipHandoverCanceled(opts *bind.FilterOpts, pendingOwner []common.Address) (*CollateralOfframpOwnershipHandoverCanceledIterator, error) {

	var pendingOwnerRule []interface{}
	for _, pendingOwnerItem := range pendingOwner {
		pendingOwnerRule = append(pendingOwnerRule, pendingOwnerItem)
	}

	logs, sub, err := _CollateralOfframp.contract.FilterLogs(opts, "OwnershipHandoverCanceled", pendingOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CollateralOfframpOwnershipHandoverCanceledIterator{contract: _CollateralOfframp.contract, event: "OwnershipHandoverCanceled", logs: logs, sub: sub}, nil
}

// WatchOwnershipHandoverCanceled is a free log subscription operation binding the contract event 0xfa7b8eab7da67f412cc9575ed43464468f9bfbae89d1675917346ca6d8fe3c92.
//
// Solidity: event OwnershipHandoverCanceled(address indexed pendingOwner)
func (_CollateralOfframp *CollateralOfframpFilterer) WatchOwnershipHandoverCanceled(opts *bind.WatchOpts, sink chan<- *CollateralOfframpOwnershipHandoverCanceled, pendingOwner []common.Address) (event.Subscription, error) {

	var pendingOwnerRule []interface{}
	for _, pendingOwnerItem := range pendingOwner {
		pendingOwnerRule = append(pendingOwnerRule, pendingOwnerItem)
	}

	logs, sub, err := _CollateralOfframp.contract.WatchLogs(opts, "OwnershipHandoverCanceled", pendingOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollateralOfframpOwnershipHandoverCanceled)
				if err := _CollateralOfframp.contract.UnpackLog(event, "OwnershipHandoverCanceled", log); err != nil {
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

// ParseOwnershipHandoverCanceled is a log parse operation binding the contract event 0xfa7b8eab7da67f412cc9575ed43464468f9bfbae89d1675917346ca6d8fe3c92.
//
// Solidity: event OwnershipHandoverCanceled(address indexed pendingOwner)
func (_CollateralOfframp *CollateralOfframpFilterer) ParseOwnershipHandoverCanceled(log types.Log) (*CollateralOfframpOwnershipHandoverCanceled, error) {
	event := new(CollateralOfframpOwnershipHandoverCanceled)
	if err := _CollateralOfframp.contract.UnpackLog(event, "OwnershipHandoverCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CollateralOfframpOwnershipHandoverRequestedIterator is returned from FilterOwnershipHandoverRequested and is used to iterate over the raw logs and unpacked data for OwnershipHandoverRequested events raised by the CollateralOfframp contract.
type CollateralOfframpOwnershipHandoverRequestedIterator struct {
	Event *CollateralOfframpOwnershipHandoverRequested // Event containing the contract specifics and raw log

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
func (it *CollateralOfframpOwnershipHandoverRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollateralOfframpOwnershipHandoverRequested)
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
		it.Event = new(CollateralOfframpOwnershipHandoverRequested)
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
func (it *CollateralOfframpOwnershipHandoverRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollateralOfframpOwnershipHandoverRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollateralOfframpOwnershipHandoverRequested represents a OwnershipHandoverRequested event raised by the CollateralOfframp contract.
type CollateralOfframpOwnershipHandoverRequested struct {
	PendingOwner common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterOwnershipHandoverRequested is a free log retrieval operation binding the contract event 0xdbf36a107da19e49527a7176a1babf963b4b0ff8cde35ee35d6cd8f1f9ac7e1d.
//
// Solidity: event OwnershipHandoverRequested(address indexed pendingOwner)
func (_CollateralOfframp *CollateralOfframpFilterer) FilterOwnershipHandoverRequested(opts *bind.FilterOpts, pendingOwner []common.Address) (*CollateralOfframpOwnershipHandoverRequestedIterator, error) {

	var pendingOwnerRule []interface{}
	for _, pendingOwnerItem := range pendingOwner {
		pendingOwnerRule = append(pendingOwnerRule, pendingOwnerItem)
	}

	logs, sub, err := _CollateralOfframp.contract.FilterLogs(opts, "OwnershipHandoverRequested", pendingOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CollateralOfframpOwnershipHandoverRequestedIterator{contract: _CollateralOfframp.contract, event: "OwnershipHandoverRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipHandoverRequested is a free log subscription operation binding the contract event 0xdbf36a107da19e49527a7176a1babf963b4b0ff8cde35ee35d6cd8f1f9ac7e1d.
//
// Solidity: event OwnershipHandoverRequested(address indexed pendingOwner)
func (_CollateralOfframp *CollateralOfframpFilterer) WatchOwnershipHandoverRequested(opts *bind.WatchOpts, sink chan<- *CollateralOfframpOwnershipHandoverRequested, pendingOwner []common.Address) (event.Subscription, error) {

	var pendingOwnerRule []interface{}
	for _, pendingOwnerItem := range pendingOwner {
		pendingOwnerRule = append(pendingOwnerRule, pendingOwnerItem)
	}

	logs, sub, err := _CollateralOfframp.contract.WatchLogs(opts, "OwnershipHandoverRequested", pendingOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollateralOfframpOwnershipHandoverRequested)
				if err := _CollateralOfframp.contract.UnpackLog(event, "OwnershipHandoverRequested", log); err != nil {
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

// ParseOwnershipHandoverRequested is a log parse operation binding the contract event 0xdbf36a107da19e49527a7176a1babf963b4b0ff8cde35ee35d6cd8f1f9ac7e1d.
//
// Solidity: event OwnershipHandoverRequested(address indexed pendingOwner)
func (_CollateralOfframp *CollateralOfframpFilterer) ParseOwnershipHandoverRequested(log types.Log) (*CollateralOfframpOwnershipHandoverRequested, error) {
	event := new(CollateralOfframpOwnershipHandoverRequested)
	if err := _CollateralOfframp.contract.UnpackLog(event, "OwnershipHandoverRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CollateralOfframpOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the CollateralOfframp contract.
type CollateralOfframpOwnershipTransferredIterator struct {
	Event *CollateralOfframpOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *CollateralOfframpOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollateralOfframpOwnershipTransferred)
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
		it.Event = new(CollateralOfframpOwnershipTransferred)
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
func (it *CollateralOfframpOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollateralOfframpOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollateralOfframpOwnershipTransferred represents a OwnershipTransferred event raised by the CollateralOfframp contract.
type CollateralOfframpOwnershipTransferred struct {
	OldOwner common.Address
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed oldOwner, address indexed newOwner)
func (_CollateralOfframp *CollateralOfframpFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, oldOwner []common.Address, newOwner []common.Address) (*CollateralOfframpOwnershipTransferredIterator, error) {

	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CollateralOfframp.contract.FilterLogs(opts, "OwnershipTransferred", oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CollateralOfframpOwnershipTransferredIterator{contract: _CollateralOfframp.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed oldOwner, address indexed newOwner)
func (_CollateralOfframp *CollateralOfframpFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CollateralOfframpOwnershipTransferred, oldOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CollateralOfframp.contract.WatchLogs(opts, "OwnershipTransferred", oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollateralOfframpOwnershipTransferred)
				if err := _CollateralOfframp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
// Solidity: event OwnershipTransferred(address indexed oldOwner, address indexed newOwner)
func (_CollateralOfframp *CollateralOfframpFilterer) ParseOwnershipTransferred(log types.Log) (*CollateralOfframpOwnershipTransferred, error) {
	event := new(CollateralOfframpOwnershipTransferred)
	if err := _CollateralOfframp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CollateralOfframpPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the CollateralOfframp contract.
type CollateralOfframpPausedIterator struct {
	Event *CollateralOfframpPaused // Event containing the contract specifics and raw log

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
func (it *CollateralOfframpPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollateralOfframpPaused)
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
		it.Event = new(CollateralOfframpPaused)
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
func (it *CollateralOfframpPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollateralOfframpPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollateralOfframpPaused represents a Paused event raised by the CollateralOfframp contract.
type CollateralOfframpPaused struct {
	Asset common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address indexed asset)
func (_CollateralOfframp *CollateralOfframpFilterer) FilterPaused(opts *bind.FilterOpts, asset []common.Address) (*CollateralOfframpPausedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _CollateralOfframp.contract.FilterLogs(opts, "Paused", assetRule)
	if err != nil {
		return nil, err
	}
	return &CollateralOfframpPausedIterator{contract: _CollateralOfframp.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address indexed asset)
func (_CollateralOfframp *CollateralOfframpFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *CollateralOfframpPaused, asset []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _CollateralOfframp.contract.WatchLogs(opts, "Paused", assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollateralOfframpPaused)
				if err := _CollateralOfframp.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address indexed asset)
func (_CollateralOfframp *CollateralOfframpFilterer) ParsePaused(log types.Log) (*CollateralOfframpPaused, error) {
	event := new(CollateralOfframpPaused)
	if err := _CollateralOfframp.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CollateralOfframpRolesUpdatedIterator is returned from FilterRolesUpdated and is used to iterate over the raw logs and unpacked data for RolesUpdated events raised by the CollateralOfframp contract.
type CollateralOfframpRolesUpdatedIterator struct {
	Event *CollateralOfframpRolesUpdated // Event containing the contract specifics and raw log

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
func (it *CollateralOfframpRolesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollateralOfframpRolesUpdated)
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
		it.Event = new(CollateralOfframpRolesUpdated)
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
func (it *CollateralOfframpRolesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollateralOfframpRolesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollateralOfframpRolesUpdated represents a RolesUpdated event raised by the CollateralOfframp contract.
type CollateralOfframpRolesUpdated struct {
	User  common.Address
	Roles *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterRolesUpdated is a free log retrieval operation binding the contract event 0x715ad5ce61fc9595c7b415289d59cf203f23a94fa06f04af7e489a0a76e1fe26.
//
// Solidity: event RolesUpdated(address indexed user, uint256 indexed roles)
func (_CollateralOfframp *CollateralOfframpFilterer) FilterRolesUpdated(opts *bind.FilterOpts, user []common.Address, roles []*big.Int) (*CollateralOfframpRolesUpdatedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var rolesRule []interface{}
	for _, rolesItem := range roles {
		rolesRule = append(rolesRule, rolesItem)
	}

	logs, sub, err := _CollateralOfframp.contract.FilterLogs(opts, "RolesUpdated", userRule, rolesRule)
	if err != nil {
		return nil, err
	}
	return &CollateralOfframpRolesUpdatedIterator{contract: _CollateralOfframp.contract, event: "RolesUpdated", logs: logs, sub: sub}, nil
}

// WatchRolesUpdated is a free log subscription operation binding the contract event 0x715ad5ce61fc9595c7b415289d59cf203f23a94fa06f04af7e489a0a76e1fe26.
//
// Solidity: event RolesUpdated(address indexed user, uint256 indexed roles)
func (_CollateralOfframp *CollateralOfframpFilterer) WatchRolesUpdated(opts *bind.WatchOpts, sink chan<- *CollateralOfframpRolesUpdated, user []common.Address, roles []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var rolesRule []interface{}
	for _, rolesItem := range roles {
		rolesRule = append(rolesRule, rolesItem)
	}

	logs, sub, err := _CollateralOfframp.contract.WatchLogs(opts, "RolesUpdated", userRule, rolesRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollateralOfframpRolesUpdated)
				if err := _CollateralOfframp.contract.UnpackLog(event, "RolesUpdated", log); err != nil {
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

// ParseRolesUpdated is a log parse operation binding the contract event 0x715ad5ce61fc9595c7b415289d59cf203f23a94fa06f04af7e489a0a76e1fe26.
//
// Solidity: event RolesUpdated(address indexed user, uint256 indexed roles)
func (_CollateralOfframp *CollateralOfframpFilterer) ParseRolesUpdated(log types.Log) (*CollateralOfframpRolesUpdated, error) {
	event := new(CollateralOfframpRolesUpdated)
	if err := _CollateralOfframp.contract.UnpackLog(event, "RolesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CollateralOfframpUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the CollateralOfframp contract.
type CollateralOfframpUnpausedIterator struct {
	Event *CollateralOfframpUnpaused // Event containing the contract specifics and raw log

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
func (it *CollateralOfframpUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollateralOfframpUnpaused)
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
		it.Event = new(CollateralOfframpUnpaused)
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
func (it *CollateralOfframpUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollateralOfframpUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollateralOfframpUnpaused represents a Unpaused event raised by the CollateralOfframp contract.
type CollateralOfframpUnpaused struct {
	Asset common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address indexed asset)
func (_CollateralOfframp *CollateralOfframpFilterer) FilterUnpaused(opts *bind.FilterOpts, asset []common.Address) (*CollateralOfframpUnpausedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _CollateralOfframp.contract.FilterLogs(opts, "Unpaused", assetRule)
	if err != nil {
		return nil, err
	}
	return &CollateralOfframpUnpausedIterator{contract: _CollateralOfframp.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address indexed asset)
func (_CollateralOfframp *CollateralOfframpFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *CollateralOfframpUnpaused, asset []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _CollateralOfframp.contract.WatchLogs(opts, "Unpaused", assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollateralOfframpUnpaused)
				if err := _CollateralOfframp.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address indexed asset)
func (_CollateralOfframp *CollateralOfframpFilterer) ParseUnpaused(log types.Log) (*CollateralOfframpUnpaused, error) {
	event := new(CollateralOfframpUnpaused)
	if err := _CollateralOfframp.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
