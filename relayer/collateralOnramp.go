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

// CollateralOnrampMetaData contains all meta data concerning the CollateralOnramp contract.
var CollateralOnrampMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_collateralToken\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AlreadyInitialized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpiredDeadline\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAsset\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidNonce\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewOwnerIsZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoHandoverRequest\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyUnpaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"pendingOwner\",\"type\":\"address\"}],\"name\":\"OwnershipHandoverCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"pendingOwner\",\"type\":\"address\"}],\"name\":\"OwnershipHandoverRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roles\",\"type\":\"uint256\"}],\"name\":\"RolesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COLLATERAL_TOKEN\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"}],\"name\":\"addAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"cancelOwnershipHandover\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"pendingOwner\",\"type\":\"address\"}],\"name\":\"completeOwnershipHandover\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"roles\",\"type\":\"uint256\"}],\"name\":\"grantRoles\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"roles\",\"type\":\"uint256\"}],\"name\":\"hasAllRoles\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"roles\",\"type\":\"uint256\"}],\"name\":\"hasAnyRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"result\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"pendingOwner\",\"type\":\"address\"}],\"name\":\"ownershipHandoverExpiresAt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"result\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_asset\",\"type\":\"address\"}],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"}],\"name\":\"removeAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roles\",\"type\":\"uint256\"}],\"name\":\"renounceRoles\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestOwnershipHandover\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"roles\",\"type\":\"uint256\"}],\"name\":\"revokeRoles\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"rolesOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"roles\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_asset\",\"type\":\"address\"}],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_asset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"wrap\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561000f575f5ffd5b50604051610bd5380380610bd583398101604081905261002e91610114565b6001600160a01b03811660805261004483610057565b61004f826001610092565b505050610154565b6001600160a01b0316638b78c6d819819055805f7f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08180a350565b61009e828260016100a2565b5050565b638b78c6d8600c52825f526020600c208054838117836100c3575080841681185b80835580600c5160601c7f715ad5ce61fc9595c7b415289d59cf203f23a94fa06f04af7e489a0a76e1fe265f5fa3505050505050565b80516001600160a01b038116811461010f575f5ffd5b919050565b5f5f5f60608486031215610126575f5ffd5b61012f846100f9565b925061013d602085016100f9565b915061014b604085016100f9565b90509250925092565b608051610a5b61017a5f395f81816103c0015281816105e1015261066e0152610a5b5ff3fe608060405260043610610157575f3560e01c806357b001f9116100bb5780638da5cb5b11610071578063f2fde38b11610057578063f2fde38b1461039c578063f5f1f1a7146103af578063fee81cf4146103e2575f5ffd5b80638da5cb5b14610335578063f04e283e14610389575f5ffd5b806370480275116100a157806370480275146102ef578063715018a61461030e57806376a67a5114610316575f5ffd5b806357b001f9146102b157806362355638146102d0575f5ffd5b80632de94807116101105780634a4ee7b1116100f65780634a4ee7b114610261578063514e62fc1461027457806354d1f13d146102a9575f5ffd5b80632de94807146101f45780632e48152c14610233575f5ffd5b80631c10893f116101405780631c10893f1461018f5780631cd64df4146101a257806325692962146101ec575f5ffd5b80631785f53c1461015b578063183a4f6e1461017c575b5f5ffd5b348015610166575f5ffd5b5061017a61017536600461098c565b610413565b005b61017a61018a3660046109ac565b61042d565b61017a61019d3660046109c3565b61043a565b3480156101ad575f5ffd5b506101d76101bc3660046109c3565b638b78c6d8600c9081525f9290925260209091205481161490565b60405190151581526020015b60405180910390f35b61017a61044c565b3480156101ff575f5ffd5b5061022561020e36600461098c565b638b78c6d8600c9081525f91909152602090205490565b6040519081526020016101e3565b34801561023e575f5ffd5b506101d761024d36600461098c565b5f6020819052908152604090205460ff1681565b61017a61026f3660046109c3565b610499565b34801561027f575f5ffd5b506101d761028e3660046109c3565b638b78c6d8600c9081525f9290925260209091205416151590565b61017a6104ab565b3480156102bc575f5ffd5b5061017a6102cb36600461098c565b6104e4565b3480156102db575f5ffd5b5061017a6102ea3660046109eb565b610563565b3480156102fa575f5ffd5b5061017a61030936600461098c565b6106cb565b61017a6106e1565b348015610321575f5ffd5b5061017a61033036600461098c565b6106f4565b348015610340575f5ffd5b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffff74873927545b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016101e3565b61017a61039736600461098c565b610776565b61017a6103aa36600461098c565b6107b0565b3480156103ba575f5ffd5b506103647f000000000000000000000000000000000000000000000000000000000000000081565b3480156103ed575f5ffd5b506102256103fc36600461098c565b63389a75e1600c9081525f91909152602090205490565b600161041e816107d6565b6104298260016107fa565b5050565b61043733826107fa565b50565b610442610805565b610429828261083a565b5f6202a30067ffffffffffffffff164201905063389a75e1600c52335f52806020600c2055337fdbf36a107da19e49527a7176a1babf963b4b0ff8cde35ee35d6cd8f1f9ac7e1d5f5fa250565b6104a1610805565b61042982826107fa565b63389a75e1600c52335f525f6020600c2055337ffa7b8eab7da67f412cc9575ed43464468f9bfbae89d1675917346ca6d8fe3c925f5fa2565b60016104ef816107d6565b73ffffffffffffffffffffffffffffffffffffffff82165f8181526020819052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055517f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa9190a25050565b73ffffffffffffffffffffffffffffffffffffffff83165f90815260208190526040902054839060ff16156105c4576040517f49b8b3ac00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61060673ffffffffffffffffffffffffffffffffffffffff8516337f000000000000000000000000000000000000000000000000000000000000000085610846565b6040517fb97b57c700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff85811660048301528481166024830152604482018490525f6064830181905260a0608484015260a48301527f0000000000000000000000000000000000000000000000000000000000000000169063b97b57c79060c4015f604051808303815f87803b1580156106af575f5ffd5b505af11580156106c1573d5f5f3e3d5ffd5b5050505050505050565b60016106d6816107d6565b61042982600161083a565b6106e9610805565b6106f25f6108a8565b565b60016106ff816107d6565b73ffffffffffffffffffffffffffffffffffffffff82165f8181526020819052604080822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055517f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2589190a25050565b61077e610805565b63389a75e1600c52805f526020600c2080544211156107a457636f5e88185f526004601cfd5b5f9055610437816108a8565b6107b8610805565b8060601b6107cd57637448fbae5f526004601cfd5b610437816108a8565b638b78c6d8600c52335f52806020600c205416610437576382b429005f526004601cfd5b61042982825f61090d565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffff748739275433146106f2576382b429005f526004601cfd5b6104298282600161090d565b60405181606052826040528360601b602c526f23b872dd000000000000000000000000600c5260205f6064601c5f895af18060015f51141661089a57803d873b15171061089a57637939f4245f526004601cfd5b505f60605260405250505050565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffff74873927805473ffffffffffffffffffffffffffffffffffffffff9092169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e05f80a355565b638b78c6d8600c52825f526020600c2080548381178361092e575080841681185b80835580600c5160601c7f715ad5ce61fc9595c7b415289d59cf203f23a94fa06f04af7e489a0a76e1fe265f5fa3505050505050565b803573ffffffffffffffffffffffffffffffffffffffff81168114610987575f5ffd5b919050565b5f6020828403121561099c575f5ffd5b6109a582610964565b9392505050565b5f602082840312156109bc575f5ffd5b5035919050565b5f5f604083850312156109d4575f5ffd5b6109dd83610964565b946020939093013593505050565b5f5f5f606084860312156109fd575f5ffd5b610a0684610964565b9250610a1460208501610964565b92959294505050604091909101359056fea26469706673582212202144bd45778f76447f845face3e8ff8a1af84301236eb8215154cc7229f0dbda64736f6c6343000822003300000000000000000000000047ebfac3353314c788b96cdcbf41daadfe03629c0000000000000000000000003dce0a29139a851da1dfca56af8e8a6440b4d952000000000000000000000000c011a7e12a19f7b1f670d46f03b03f3342e82dfb",
}

// CollateralOnrampABI is the input ABI used to generate the binding from.
// Deprecated: Use CollateralOnrampMetaData.ABI instead.
var CollateralOnrampABI = CollateralOnrampMetaData.ABI

// CollateralOnrampBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CollateralOnrampMetaData.Bin instead.
var CollateralOnrampBin = CollateralOnrampMetaData.Bin

// DeployCollateralOnramp deploys a new Ethereum contract, binding an instance of CollateralOnramp to it.
func DeployCollateralOnramp(auth *bind.TransactOpts, backend bind.ContractBackend, _owner common.Address, _admin common.Address, _collateralToken common.Address) (common.Address, *types.Transaction, *CollateralOnramp, error) {
	parsed, err := CollateralOnrampMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CollateralOnrampBin), backend, _owner, _admin, _collateralToken)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CollateralOnramp{CollateralOnrampCaller: CollateralOnrampCaller{contract: contract}, CollateralOnrampTransactor: CollateralOnrampTransactor{contract: contract}, CollateralOnrampFilterer: CollateralOnrampFilterer{contract: contract}}, nil
}

// CollateralOnramp is an auto generated Go binding around an Ethereum contract.
type CollateralOnramp struct {
	CollateralOnrampCaller     // Read-only binding to the contract
	CollateralOnrampTransactor // Write-only binding to the contract
	CollateralOnrampFilterer   // Log filterer for contract events
}

// CollateralOnrampCaller is an auto generated read-only Go binding around an Ethereum contract.
type CollateralOnrampCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CollateralOnrampTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CollateralOnrampTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CollateralOnrampFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CollateralOnrampFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CollateralOnrampSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CollateralOnrampSession struct {
	Contract     *CollateralOnramp // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CollateralOnrampCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CollateralOnrampCallerSession struct {
	Contract *CollateralOnrampCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// CollateralOnrampTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CollateralOnrampTransactorSession struct {
	Contract     *CollateralOnrampTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// CollateralOnrampRaw is an auto generated low-level Go binding around an Ethereum contract.
type CollateralOnrampRaw struct {
	Contract *CollateralOnramp // Generic contract binding to access the raw methods on
}

// CollateralOnrampCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CollateralOnrampCallerRaw struct {
	Contract *CollateralOnrampCaller // Generic read-only contract binding to access the raw methods on
}

// CollateralOnrampTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CollateralOnrampTransactorRaw struct {
	Contract *CollateralOnrampTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCollateralOnramp creates a new instance of CollateralOnramp, bound to a specific deployed contract.
func NewCollateralOnramp(address common.Address, backend bind.ContractBackend) (*CollateralOnramp, error) {
	contract, err := bindCollateralOnramp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CollateralOnramp{CollateralOnrampCaller: CollateralOnrampCaller{contract: contract}, CollateralOnrampTransactor: CollateralOnrampTransactor{contract: contract}, CollateralOnrampFilterer: CollateralOnrampFilterer{contract: contract}}, nil
}

// NewCollateralOnrampCaller creates a new read-only instance of CollateralOnramp, bound to a specific deployed contract.
func NewCollateralOnrampCaller(address common.Address, caller bind.ContractCaller) (*CollateralOnrampCaller, error) {
	contract, err := bindCollateralOnramp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CollateralOnrampCaller{contract: contract}, nil
}

// NewCollateralOnrampTransactor creates a new write-only instance of CollateralOnramp, bound to a specific deployed contract.
func NewCollateralOnrampTransactor(address common.Address, transactor bind.ContractTransactor) (*CollateralOnrampTransactor, error) {
	contract, err := bindCollateralOnramp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CollateralOnrampTransactor{contract: contract}, nil
}

// NewCollateralOnrampFilterer creates a new log filterer instance of CollateralOnramp, bound to a specific deployed contract.
func NewCollateralOnrampFilterer(address common.Address, filterer bind.ContractFilterer) (*CollateralOnrampFilterer, error) {
	contract, err := bindCollateralOnramp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CollateralOnrampFilterer{contract: contract}, nil
}

// bindCollateralOnramp binds a generic wrapper to an already deployed contract.
func bindCollateralOnramp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CollateralOnrampMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CollateralOnramp *CollateralOnrampRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CollateralOnramp.Contract.CollateralOnrampCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CollateralOnramp *CollateralOnrampRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CollateralOnramp.Contract.CollateralOnrampTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CollateralOnramp *CollateralOnrampRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CollateralOnramp.Contract.CollateralOnrampTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CollateralOnramp *CollateralOnrampCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CollateralOnramp.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CollateralOnramp *CollateralOnrampTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CollateralOnramp.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CollateralOnramp *CollateralOnrampTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CollateralOnramp.Contract.contract.Transact(opts, method, params...)
}

// COLLATERALTOKEN is a free data retrieval call binding the contract method 0xf5f1f1a7.
//
// Solidity: function COLLATERAL_TOKEN() view returns(address)
func (_CollateralOnramp *CollateralOnrampCaller) COLLATERALTOKEN(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CollateralOnramp.contract.Call(opts, &out, "COLLATERAL_TOKEN")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// COLLATERALTOKEN is a free data retrieval call binding the contract method 0xf5f1f1a7.
//
// Solidity: function COLLATERAL_TOKEN() view returns(address)
func (_CollateralOnramp *CollateralOnrampSession) COLLATERALTOKEN() (common.Address, error) {
	return _CollateralOnramp.Contract.COLLATERALTOKEN(&_CollateralOnramp.CallOpts)
}

// COLLATERALTOKEN is a free data retrieval call binding the contract method 0xf5f1f1a7.
//
// Solidity: function COLLATERAL_TOKEN() view returns(address)
func (_CollateralOnramp *CollateralOnrampCallerSession) COLLATERALTOKEN() (common.Address, error) {
	return _CollateralOnramp.Contract.COLLATERALTOKEN(&_CollateralOnramp.CallOpts)
}

// HasAllRoles is a free data retrieval call binding the contract method 0x1cd64df4.
//
// Solidity: function hasAllRoles(address user, uint256 roles) view returns(bool)
func (_CollateralOnramp *CollateralOnrampCaller) HasAllRoles(opts *bind.CallOpts, user common.Address, roles *big.Int) (bool, error) {
	var out []interface{}
	err := _CollateralOnramp.contract.Call(opts, &out, "hasAllRoles", user, roles)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasAllRoles is a free data retrieval call binding the contract method 0x1cd64df4.
//
// Solidity: function hasAllRoles(address user, uint256 roles) view returns(bool)
func (_CollateralOnramp *CollateralOnrampSession) HasAllRoles(user common.Address, roles *big.Int) (bool, error) {
	return _CollateralOnramp.Contract.HasAllRoles(&_CollateralOnramp.CallOpts, user, roles)
}

// HasAllRoles is a free data retrieval call binding the contract method 0x1cd64df4.
//
// Solidity: function hasAllRoles(address user, uint256 roles) view returns(bool)
func (_CollateralOnramp *CollateralOnrampCallerSession) HasAllRoles(user common.Address, roles *big.Int) (bool, error) {
	return _CollateralOnramp.Contract.HasAllRoles(&_CollateralOnramp.CallOpts, user, roles)
}

// HasAnyRole is a free data retrieval call binding the contract method 0x514e62fc.
//
// Solidity: function hasAnyRole(address user, uint256 roles) view returns(bool)
func (_CollateralOnramp *CollateralOnrampCaller) HasAnyRole(opts *bind.CallOpts, user common.Address, roles *big.Int) (bool, error) {
	var out []interface{}
	err := _CollateralOnramp.contract.Call(opts, &out, "hasAnyRole", user, roles)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasAnyRole is a free data retrieval call binding the contract method 0x514e62fc.
//
// Solidity: function hasAnyRole(address user, uint256 roles) view returns(bool)
func (_CollateralOnramp *CollateralOnrampSession) HasAnyRole(user common.Address, roles *big.Int) (bool, error) {
	return _CollateralOnramp.Contract.HasAnyRole(&_CollateralOnramp.CallOpts, user, roles)
}

// HasAnyRole is a free data retrieval call binding the contract method 0x514e62fc.
//
// Solidity: function hasAnyRole(address user, uint256 roles) view returns(bool)
func (_CollateralOnramp *CollateralOnrampCallerSession) HasAnyRole(user common.Address, roles *big.Int) (bool, error) {
	return _CollateralOnramp.Contract.HasAnyRole(&_CollateralOnramp.CallOpts, user, roles)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address result)
func (_CollateralOnramp *CollateralOnrampCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CollateralOnramp.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address result)
func (_CollateralOnramp *CollateralOnrampSession) Owner() (common.Address, error) {
	return _CollateralOnramp.Contract.Owner(&_CollateralOnramp.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address result)
func (_CollateralOnramp *CollateralOnrampCallerSession) Owner() (common.Address, error) {
	return _CollateralOnramp.Contract.Owner(&_CollateralOnramp.CallOpts)
}

// OwnershipHandoverExpiresAt is a free data retrieval call binding the contract method 0xfee81cf4.
//
// Solidity: function ownershipHandoverExpiresAt(address pendingOwner) view returns(uint256 result)
func (_CollateralOnramp *CollateralOnrampCaller) OwnershipHandoverExpiresAt(opts *bind.CallOpts, pendingOwner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CollateralOnramp.contract.Call(opts, &out, "ownershipHandoverExpiresAt", pendingOwner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OwnershipHandoverExpiresAt is a free data retrieval call binding the contract method 0xfee81cf4.
//
// Solidity: function ownershipHandoverExpiresAt(address pendingOwner) view returns(uint256 result)
func (_CollateralOnramp *CollateralOnrampSession) OwnershipHandoverExpiresAt(pendingOwner common.Address) (*big.Int, error) {
	return _CollateralOnramp.Contract.OwnershipHandoverExpiresAt(&_CollateralOnramp.CallOpts, pendingOwner)
}

// OwnershipHandoverExpiresAt is a free data retrieval call binding the contract method 0xfee81cf4.
//
// Solidity: function ownershipHandoverExpiresAt(address pendingOwner) view returns(uint256 result)
func (_CollateralOnramp *CollateralOnrampCallerSession) OwnershipHandoverExpiresAt(pendingOwner common.Address) (*big.Int, error) {
	return _CollateralOnramp.Contract.OwnershipHandoverExpiresAt(&_CollateralOnramp.CallOpts, pendingOwner)
}

// Paused is a free data retrieval call binding the contract method 0x2e48152c.
//
// Solidity: function paused(address ) view returns(bool)
func (_CollateralOnramp *CollateralOnrampCaller) Paused(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _CollateralOnramp.contract.Call(opts, &out, "paused", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x2e48152c.
//
// Solidity: function paused(address ) view returns(bool)
func (_CollateralOnramp *CollateralOnrampSession) Paused(arg0 common.Address) (bool, error) {
	return _CollateralOnramp.Contract.Paused(&_CollateralOnramp.CallOpts, arg0)
}

// Paused is a free data retrieval call binding the contract method 0x2e48152c.
//
// Solidity: function paused(address ) view returns(bool)
func (_CollateralOnramp *CollateralOnrampCallerSession) Paused(arg0 common.Address) (bool, error) {
	return _CollateralOnramp.Contract.Paused(&_CollateralOnramp.CallOpts, arg0)
}

// RolesOf is a free data retrieval call binding the contract method 0x2de94807.
//
// Solidity: function rolesOf(address user) view returns(uint256 roles)
func (_CollateralOnramp *CollateralOnrampCaller) RolesOf(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CollateralOnramp.contract.Call(opts, &out, "rolesOf", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RolesOf is a free data retrieval call binding the contract method 0x2de94807.
//
// Solidity: function rolesOf(address user) view returns(uint256 roles)
func (_CollateralOnramp *CollateralOnrampSession) RolesOf(user common.Address) (*big.Int, error) {
	return _CollateralOnramp.Contract.RolesOf(&_CollateralOnramp.CallOpts, user)
}

// RolesOf is a free data retrieval call binding the contract method 0x2de94807.
//
// Solidity: function rolesOf(address user) view returns(uint256 roles)
func (_CollateralOnramp *CollateralOnrampCallerSession) RolesOf(user common.Address) (*big.Int, error) {
	return _CollateralOnramp.Contract.RolesOf(&_CollateralOnramp.CallOpts, user)
}

// AddAdmin is a paid mutator transaction binding the contract method 0x70480275.
//
// Solidity: function addAdmin(address _admin) returns()
func (_CollateralOnramp *CollateralOnrampTransactor) AddAdmin(opts *bind.TransactOpts, _admin common.Address) (*types.Transaction, error) {
	return _CollateralOnramp.contract.Transact(opts, "addAdmin", _admin)
}

// AddAdmin is a paid mutator transaction binding the contract method 0x70480275.
//
// Solidity: function addAdmin(address _admin) returns()
func (_CollateralOnramp *CollateralOnrampSession) AddAdmin(_admin common.Address) (*types.Transaction, error) {
	return _CollateralOnramp.Contract.AddAdmin(&_CollateralOnramp.TransactOpts, _admin)
}

// AddAdmin is a paid mutator transaction binding the contract method 0x70480275.
//
// Solidity: function addAdmin(address _admin) returns()
func (_CollateralOnramp *CollateralOnrampTransactorSession) AddAdmin(_admin common.Address) (*types.Transaction, error) {
	return _CollateralOnramp.Contract.AddAdmin(&_CollateralOnramp.TransactOpts, _admin)
}

// CancelOwnershipHandover is a paid mutator transaction binding the contract method 0x54d1f13d.
//
// Solidity: function cancelOwnershipHandover() payable returns()
func (_CollateralOnramp *CollateralOnrampTransactor) CancelOwnershipHandover(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CollateralOnramp.contract.Transact(opts, "cancelOwnershipHandover")
}

// CancelOwnershipHandover is a paid mutator transaction binding the contract method 0x54d1f13d.
//
// Solidity: function cancelOwnershipHandover() payable returns()
func (_CollateralOnramp *CollateralOnrampSession) CancelOwnershipHandover() (*types.Transaction, error) {
	return _CollateralOnramp.Contract.CancelOwnershipHandover(&_CollateralOnramp.TransactOpts)
}

// CancelOwnershipHandover is a paid mutator transaction binding the contract method 0x54d1f13d.
//
// Solidity: function cancelOwnershipHandover() payable returns()
func (_CollateralOnramp *CollateralOnrampTransactorSession) CancelOwnershipHandover() (*types.Transaction, error) {
	return _CollateralOnramp.Contract.CancelOwnershipHandover(&_CollateralOnramp.TransactOpts)
}

// CompleteOwnershipHandover is a paid mutator transaction binding the contract method 0xf04e283e.
//
// Solidity: function completeOwnershipHandover(address pendingOwner) payable returns()
func (_CollateralOnramp *CollateralOnrampTransactor) CompleteOwnershipHandover(opts *bind.TransactOpts, pendingOwner common.Address) (*types.Transaction, error) {
	return _CollateralOnramp.contract.Transact(opts, "completeOwnershipHandover", pendingOwner)
}

// CompleteOwnershipHandover is a paid mutator transaction binding the contract method 0xf04e283e.
//
// Solidity: function completeOwnershipHandover(address pendingOwner) payable returns()
func (_CollateralOnramp *CollateralOnrampSession) CompleteOwnershipHandover(pendingOwner common.Address) (*types.Transaction, error) {
	return _CollateralOnramp.Contract.CompleteOwnershipHandover(&_CollateralOnramp.TransactOpts, pendingOwner)
}

// CompleteOwnershipHandover is a paid mutator transaction binding the contract method 0xf04e283e.
//
// Solidity: function completeOwnershipHandover(address pendingOwner) payable returns()
func (_CollateralOnramp *CollateralOnrampTransactorSession) CompleteOwnershipHandover(pendingOwner common.Address) (*types.Transaction, error) {
	return _CollateralOnramp.Contract.CompleteOwnershipHandover(&_CollateralOnramp.TransactOpts, pendingOwner)
}

// GrantRoles is a paid mutator transaction binding the contract method 0x1c10893f.
//
// Solidity: function grantRoles(address user, uint256 roles) payable returns()
func (_CollateralOnramp *CollateralOnrampTransactor) GrantRoles(opts *bind.TransactOpts, user common.Address, roles *big.Int) (*types.Transaction, error) {
	return _CollateralOnramp.contract.Transact(opts, "grantRoles", user, roles)
}

// GrantRoles is a paid mutator transaction binding the contract method 0x1c10893f.
//
// Solidity: function grantRoles(address user, uint256 roles) payable returns()
func (_CollateralOnramp *CollateralOnrampSession) GrantRoles(user common.Address, roles *big.Int) (*types.Transaction, error) {
	return _CollateralOnramp.Contract.GrantRoles(&_CollateralOnramp.TransactOpts, user, roles)
}

// GrantRoles is a paid mutator transaction binding the contract method 0x1c10893f.
//
// Solidity: function grantRoles(address user, uint256 roles) payable returns()
func (_CollateralOnramp *CollateralOnrampTransactorSession) GrantRoles(user common.Address, roles *big.Int) (*types.Transaction, error) {
	return _CollateralOnramp.Contract.GrantRoles(&_CollateralOnramp.TransactOpts, user, roles)
}

// Pause is a paid mutator transaction binding the contract method 0x76a67a51.
//
// Solidity: function pause(address _asset) returns()
func (_CollateralOnramp *CollateralOnrampTransactor) Pause(opts *bind.TransactOpts, _asset common.Address) (*types.Transaction, error) {
	return _CollateralOnramp.contract.Transact(opts, "pause", _asset)
}

// Pause is a paid mutator transaction binding the contract method 0x76a67a51.
//
// Solidity: function pause(address _asset) returns()
func (_CollateralOnramp *CollateralOnrampSession) Pause(_asset common.Address) (*types.Transaction, error) {
	return _CollateralOnramp.Contract.Pause(&_CollateralOnramp.TransactOpts, _asset)
}

// Pause is a paid mutator transaction binding the contract method 0x76a67a51.
//
// Solidity: function pause(address _asset) returns()
func (_CollateralOnramp *CollateralOnrampTransactorSession) Pause(_asset common.Address) (*types.Transaction, error) {
	return _CollateralOnramp.Contract.Pause(&_CollateralOnramp.TransactOpts, _asset)
}

// RemoveAdmin is a paid mutator transaction binding the contract method 0x1785f53c.
//
// Solidity: function removeAdmin(address _admin) returns()
func (_CollateralOnramp *CollateralOnrampTransactor) RemoveAdmin(opts *bind.TransactOpts, _admin common.Address) (*types.Transaction, error) {
	return _CollateralOnramp.contract.Transact(opts, "removeAdmin", _admin)
}

// RemoveAdmin is a paid mutator transaction binding the contract method 0x1785f53c.
//
// Solidity: function removeAdmin(address _admin) returns()
func (_CollateralOnramp *CollateralOnrampSession) RemoveAdmin(_admin common.Address) (*types.Transaction, error) {
	return _CollateralOnramp.Contract.RemoveAdmin(&_CollateralOnramp.TransactOpts, _admin)
}

// RemoveAdmin is a paid mutator transaction binding the contract method 0x1785f53c.
//
// Solidity: function removeAdmin(address _admin) returns()
func (_CollateralOnramp *CollateralOnrampTransactorSession) RemoveAdmin(_admin common.Address) (*types.Transaction, error) {
	return _CollateralOnramp.Contract.RemoveAdmin(&_CollateralOnramp.TransactOpts, _admin)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() payable returns()
func (_CollateralOnramp *CollateralOnrampTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CollateralOnramp.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() payable returns()
func (_CollateralOnramp *CollateralOnrampSession) RenounceOwnership() (*types.Transaction, error) {
	return _CollateralOnramp.Contract.RenounceOwnership(&_CollateralOnramp.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() payable returns()
func (_CollateralOnramp *CollateralOnrampTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _CollateralOnramp.Contract.RenounceOwnership(&_CollateralOnramp.TransactOpts)
}

// RenounceRoles is a paid mutator transaction binding the contract method 0x183a4f6e.
//
// Solidity: function renounceRoles(uint256 roles) payable returns()
func (_CollateralOnramp *CollateralOnrampTransactor) RenounceRoles(opts *bind.TransactOpts, roles *big.Int) (*types.Transaction, error) {
	return _CollateralOnramp.contract.Transact(opts, "renounceRoles", roles)
}

// RenounceRoles is a paid mutator transaction binding the contract method 0x183a4f6e.
//
// Solidity: function renounceRoles(uint256 roles) payable returns()
func (_CollateralOnramp *CollateralOnrampSession) RenounceRoles(roles *big.Int) (*types.Transaction, error) {
	return _CollateralOnramp.Contract.RenounceRoles(&_CollateralOnramp.TransactOpts, roles)
}

// RenounceRoles is a paid mutator transaction binding the contract method 0x183a4f6e.
//
// Solidity: function renounceRoles(uint256 roles) payable returns()
func (_CollateralOnramp *CollateralOnrampTransactorSession) RenounceRoles(roles *big.Int) (*types.Transaction, error) {
	return _CollateralOnramp.Contract.RenounceRoles(&_CollateralOnramp.TransactOpts, roles)
}

// RequestOwnershipHandover is a paid mutator transaction binding the contract method 0x25692962.
//
// Solidity: function requestOwnershipHandover() payable returns()
func (_CollateralOnramp *CollateralOnrampTransactor) RequestOwnershipHandover(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CollateralOnramp.contract.Transact(opts, "requestOwnershipHandover")
}

// RequestOwnershipHandover is a paid mutator transaction binding the contract method 0x25692962.
//
// Solidity: function requestOwnershipHandover() payable returns()
func (_CollateralOnramp *CollateralOnrampSession) RequestOwnershipHandover() (*types.Transaction, error) {
	return _CollateralOnramp.Contract.RequestOwnershipHandover(&_CollateralOnramp.TransactOpts)
}

// RequestOwnershipHandover is a paid mutator transaction binding the contract method 0x25692962.
//
// Solidity: function requestOwnershipHandover() payable returns()
func (_CollateralOnramp *CollateralOnrampTransactorSession) RequestOwnershipHandover() (*types.Transaction, error) {
	return _CollateralOnramp.Contract.RequestOwnershipHandover(&_CollateralOnramp.TransactOpts)
}

// RevokeRoles is a paid mutator transaction binding the contract method 0x4a4ee7b1.
//
// Solidity: function revokeRoles(address user, uint256 roles) payable returns()
func (_CollateralOnramp *CollateralOnrampTransactor) RevokeRoles(opts *bind.TransactOpts, user common.Address, roles *big.Int) (*types.Transaction, error) {
	return _CollateralOnramp.contract.Transact(opts, "revokeRoles", user, roles)
}

// RevokeRoles is a paid mutator transaction binding the contract method 0x4a4ee7b1.
//
// Solidity: function revokeRoles(address user, uint256 roles) payable returns()
func (_CollateralOnramp *CollateralOnrampSession) RevokeRoles(user common.Address, roles *big.Int) (*types.Transaction, error) {
	return _CollateralOnramp.Contract.RevokeRoles(&_CollateralOnramp.TransactOpts, user, roles)
}

// RevokeRoles is a paid mutator transaction binding the contract method 0x4a4ee7b1.
//
// Solidity: function revokeRoles(address user, uint256 roles) payable returns()
func (_CollateralOnramp *CollateralOnrampTransactorSession) RevokeRoles(user common.Address, roles *big.Int) (*types.Transaction, error) {
	return _CollateralOnramp.Contract.RevokeRoles(&_CollateralOnramp.TransactOpts, user, roles)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) payable returns()
func (_CollateralOnramp *CollateralOnrampTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _CollateralOnramp.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) payable returns()
func (_CollateralOnramp *CollateralOnrampSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CollateralOnramp.Contract.TransferOwnership(&_CollateralOnramp.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) payable returns()
func (_CollateralOnramp *CollateralOnrampTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CollateralOnramp.Contract.TransferOwnership(&_CollateralOnramp.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x57b001f9.
//
// Solidity: function unpause(address _asset) returns()
func (_CollateralOnramp *CollateralOnrampTransactor) Unpause(opts *bind.TransactOpts, _asset common.Address) (*types.Transaction, error) {
	return _CollateralOnramp.contract.Transact(opts, "unpause", _asset)
}

// Unpause is a paid mutator transaction binding the contract method 0x57b001f9.
//
// Solidity: function unpause(address _asset) returns()
func (_CollateralOnramp *CollateralOnrampSession) Unpause(_asset common.Address) (*types.Transaction, error) {
	return _CollateralOnramp.Contract.Unpause(&_CollateralOnramp.TransactOpts, _asset)
}

// Unpause is a paid mutator transaction binding the contract method 0x57b001f9.
//
// Solidity: function unpause(address _asset) returns()
func (_CollateralOnramp *CollateralOnrampTransactorSession) Unpause(_asset common.Address) (*types.Transaction, error) {
	return _CollateralOnramp.Contract.Unpause(&_CollateralOnramp.TransactOpts, _asset)
}

// Wrap is a paid mutator transaction binding the contract method 0x62355638.
//
// Solidity: function wrap(address _asset, address _to, uint256 _amount) returns()
func (_CollateralOnramp *CollateralOnrampTransactor) Wrap(opts *bind.TransactOpts, _asset common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CollateralOnramp.contract.Transact(opts, "wrap", _asset, _to, _amount)
}

// Wrap is a paid mutator transaction binding the contract method 0x62355638.
//
// Solidity: function wrap(address _asset, address _to, uint256 _amount) returns()
func (_CollateralOnramp *CollateralOnrampSession) Wrap(_asset common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CollateralOnramp.Contract.Wrap(&_CollateralOnramp.TransactOpts, _asset, _to, _amount)
}

// Wrap is a paid mutator transaction binding the contract method 0x62355638.
//
// Solidity: function wrap(address _asset, address _to, uint256 _amount) returns()
func (_CollateralOnramp *CollateralOnrampTransactorSession) Wrap(_asset common.Address, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _CollateralOnramp.Contract.Wrap(&_CollateralOnramp.TransactOpts, _asset, _to, _amount)
}

// CollateralOnrampOwnershipHandoverCanceledIterator is returned from FilterOwnershipHandoverCanceled and is used to iterate over the raw logs and unpacked data for OwnershipHandoverCanceled events raised by the CollateralOnramp contract.
type CollateralOnrampOwnershipHandoverCanceledIterator struct {
	Event *CollateralOnrampOwnershipHandoverCanceled // Event containing the contract specifics and raw log

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
func (it *CollateralOnrampOwnershipHandoverCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollateralOnrampOwnershipHandoverCanceled)
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
		it.Event = new(CollateralOnrampOwnershipHandoverCanceled)
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
func (it *CollateralOnrampOwnershipHandoverCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollateralOnrampOwnershipHandoverCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollateralOnrampOwnershipHandoverCanceled represents a OwnershipHandoverCanceled event raised by the CollateralOnramp contract.
type CollateralOnrampOwnershipHandoverCanceled struct {
	PendingOwner common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterOwnershipHandoverCanceled is a free log retrieval operation binding the contract event 0xfa7b8eab7da67f412cc9575ed43464468f9bfbae89d1675917346ca6d8fe3c92.
//
// Solidity: event OwnershipHandoverCanceled(address indexed pendingOwner)
func (_CollateralOnramp *CollateralOnrampFilterer) FilterOwnershipHandoverCanceled(opts *bind.FilterOpts, pendingOwner []common.Address) (*CollateralOnrampOwnershipHandoverCanceledIterator, error) {

	var pendingOwnerRule []interface{}
	for _, pendingOwnerItem := range pendingOwner {
		pendingOwnerRule = append(pendingOwnerRule, pendingOwnerItem)
	}

	logs, sub, err := _CollateralOnramp.contract.FilterLogs(opts, "OwnershipHandoverCanceled", pendingOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CollateralOnrampOwnershipHandoverCanceledIterator{contract: _CollateralOnramp.contract, event: "OwnershipHandoverCanceled", logs: logs, sub: sub}, nil
}

// WatchOwnershipHandoverCanceled is a free log subscription operation binding the contract event 0xfa7b8eab7da67f412cc9575ed43464468f9bfbae89d1675917346ca6d8fe3c92.
//
// Solidity: event OwnershipHandoverCanceled(address indexed pendingOwner)
func (_CollateralOnramp *CollateralOnrampFilterer) WatchOwnershipHandoverCanceled(opts *bind.WatchOpts, sink chan<- *CollateralOnrampOwnershipHandoverCanceled, pendingOwner []common.Address) (event.Subscription, error) {

	var pendingOwnerRule []interface{}
	for _, pendingOwnerItem := range pendingOwner {
		pendingOwnerRule = append(pendingOwnerRule, pendingOwnerItem)
	}

	logs, sub, err := _CollateralOnramp.contract.WatchLogs(opts, "OwnershipHandoverCanceled", pendingOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollateralOnrampOwnershipHandoverCanceled)
				if err := _CollateralOnramp.contract.UnpackLog(event, "OwnershipHandoverCanceled", log); err != nil {
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
func (_CollateralOnramp *CollateralOnrampFilterer) ParseOwnershipHandoverCanceled(log types.Log) (*CollateralOnrampOwnershipHandoverCanceled, error) {
	event := new(CollateralOnrampOwnershipHandoverCanceled)
	if err := _CollateralOnramp.contract.UnpackLog(event, "OwnershipHandoverCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CollateralOnrampOwnershipHandoverRequestedIterator is returned from FilterOwnershipHandoverRequested and is used to iterate over the raw logs and unpacked data for OwnershipHandoverRequested events raised by the CollateralOnramp contract.
type CollateralOnrampOwnershipHandoverRequestedIterator struct {
	Event *CollateralOnrampOwnershipHandoverRequested // Event containing the contract specifics and raw log

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
func (it *CollateralOnrampOwnershipHandoverRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollateralOnrampOwnershipHandoverRequested)
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
		it.Event = new(CollateralOnrampOwnershipHandoverRequested)
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
func (it *CollateralOnrampOwnershipHandoverRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollateralOnrampOwnershipHandoverRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollateralOnrampOwnershipHandoverRequested represents a OwnershipHandoverRequested event raised by the CollateralOnramp contract.
type CollateralOnrampOwnershipHandoverRequested struct {
	PendingOwner common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterOwnershipHandoverRequested is a free log retrieval operation binding the contract event 0xdbf36a107da19e49527a7176a1babf963b4b0ff8cde35ee35d6cd8f1f9ac7e1d.
//
// Solidity: event OwnershipHandoverRequested(address indexed pendingOwner)
func (_CollateralOnramp *CollateralOnrampFilterer) FilterOwnershipHandoverRequested(opts *bind.FilterOpts, pendingOwner []common.Address) (*CollateralOnrampOwnershipHandoverRequestedIterator, error) {

	var pendingOwnerRule []interface{}
	for _, pendingOwnerItem := range pendingOwner {
		pendingOwnerRule = append(pendingOwnerRule, pendingOwnerItem)
	}

	logs, sub, err := _CollateralOnramp.contract.FilterLogs(opts, "OwnershipHandoverRequested", pendingOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CollateralOnrampOwnershipHandoverRequestedIterator{contract: _CollateralOnramp.contract, event: "OwnershipHandoverRequested", logs: logs, sub: sub}, nil
}

// WatchOwnershipHandoverRequested is a free log subscription operation binding the contract event 0xdbf36a107da19e49527a7176a1babf963b4b0ff8cde35ee35d6cd8f1f9ac7e1d.
//
// Solidity: event OwnershipHandoverRequested(address indexed pendingOwner)
func (_CollateralOnramp *CollateralOnrampFilterer) WatchOwnershipHandoverRequested(opts *bind.WatchOpts, sink chan<- *CollateralOnrampOwnershipHandoverRequested, pendingOwner []common.Address) (event.Subscription, error) {

	var pendingOwnerRule []interface{}
	for _, pendingOwnerItem := range pendingOwner {
		pendingOwnerRule = append(pendingOwnerRule, pendingOwnerItem)
	}

	logs, sub, err := _CollateralOnramp.contract.WatchLogs(opts, "OwnershipHandoverRequested", pendingOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollateralOnrampOwnershipHandoverRequested)
				if err := _CollateralOnramp.contract.UnpackLog(event, "OwnershipHandoverRequested", log); err != nil {
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
func (_CollateralOnramp *CollateralOnrampFilterer) ParseOwnershipHandoverRequested(log types.Log) (*CollateralOnrampOwnershipHandoverRequested, error) {
	event := new(CollateralOnrampOwnershipHandoverRequested)
	if err := _CollateralOnramp.contract.UnpackLog(event, "OwnershipHandoverRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CollateralOnrampOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the CollateralOnramp contract.
type CollateralOnrampOwnershipTransferredIterator struct {
	Event *CollateralOnrampOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *CollateralOnrampOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollateralOnrampOwnershipTransferred)
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
		it.Event = new(CollateralOnrampOwnershipTransferred)
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
func (it *CollateralOnrampOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollateralOnrampOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollateralOnrampOwnershipTransferred represents a OwnershipTransferred event raised by the CollateralOnramp contract.
type CollateralOnrampOwnershipTransferred struct {
	OldOwner common.Address
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed oldOwner, address indexed newOwner)
func (_CollateralOnramp *CollateralOnrampFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, oldOwner []common.Address, newOwner []common.Address) (*CollateralOnrampOwnershipTransferredIterator, error) {

	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CollateralOnramp.contract.FilterLogs(opts, "OwnershipTransferred", oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CollateralOnrampOwnershipTransferredIterator{contract: _CollateralOnramp.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed oldOwner, address indexed newOwner)
func (_CollateralOnramp *CollateralOnrampFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CollateralOnrampOwnershipTransferred, oldOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CollateralOnramp.contract.WatchLogs(opts, "OwnershipTransferred", oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollateralOnrampOwnershipTransferred)
				if err := _CollateralOnramp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_CollateralOnramp *CollateralOnrampFilterer) ParseOwnershipTransferred(log types.Log) (*CollateralOnrampOwnershipTransferred, error) {
	event := new(CollateralOnrampOwnershipTransferred)
	if err := _CollateralOnramp.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CollateralOnrampPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the CollateralOnramp contract.
type CollateralOnrampPausedIterator struct {
	Event *CollateralOnrampPaused // Event containing the contract specifics and raw log

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
func (it *CollateralOnrampPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollateralOnrampPaused)
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
		it.Event = new(CollateralOnrampPaused)
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
func (it *CollateralOnrampPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollateralOnrampPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollateralOnrampPaused represents a Paused event raised by the CollateralOnramp contract.
type CollateralOnrampPaused struct {
	Asset common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address indexed asset)
func (_CollateralOnramp *CollateralOnrampFilterer) FilterPaused(opts *bind.FilterOpts, asset []common.Address) (*CollateralOnrampPausedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _CollateralOnramp.contract.FilterLogs(opts, "Paused", assetRule)
	if err != nil {
		return nil, err
	}
	return &CollateralOnrampPausedIterator{contract: _CollateralOnramp.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address indexed asset)
func (_CollateralOnramp *CollateralOnrampFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *CollateralOnrampPaused, asset []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _CollateralOnramp.contract.WatchLogs(opts, "Paused", assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollateralOnrampPaused)
				if err := _CollateralOnramp.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_CollateralOnramp *CollateralOnrampFilterer) ParsePaused(log types.Log) (*CollateralOnrampPaused, error) {
	event := new(CollateralOnrampPaused)
	if err := _CollateralOnramp.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CollateralOnrampRolesUpdatedIterator is returned from FilterRolesUpdated and is used to iterate over the raw logs and unpacked data for RolesUpdated events raised by the CollateralOnramp contract.
type CollateralOnrampRolesUpdatedIterator struct {
	Event *CollateralOnrampRolesUpdated // Event containing the contract specifics and raw log

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
func (it *CollateralOnrampRolesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollateralOnrampRolesUpdated)
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
		it.Event = new(CollateralOnrampRolesUpdated)
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
func (it *CollateralOnrampRolesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollateralOnrampRolesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollateralOnrampRolesUpdated represents a RolesUpdated event raised by the CollateralOnramp contract.
type CollateralOnrampRolesUpdated struct {
	User  common.Address
	Roles *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterRolesUpdated is a free log retrieval operation binding the contract event 0x715ad5ce61fc9595c7b415289d59cf203f23a94fa06f04af7e489a0a76e1fe26.
//
// Solidity: event RolesUpdated(address indexed user, uint256 indexed roles)
func (_CollateralOnramp *CollateralOnrampFilterer) FilterRolesUpdated(opts *bind.FilterOpts, user []common.Address, roles []*big.Int) (*CollateralOnrampRolesUpdatedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var rolesRule []interface{}
	for _, rolesItem := range roles {
		rolesRule = append(rolesRule, rolesItem)
	}

	logs, sub, err := _CollateralOnramp.contract.FilterLogs(opts, "RolesUpdated", userRule, rolesRule)
	if err != nil {
		return nil, err
	}
	return &CollateralOnrampRolesUpdatedIterator{contract: _CollateralOnramp.contract, event: "RolesUpdated", logs: logs, sub: sub}, nil
}

// WatchRolesUpdated is a free log subscription operation binding the contract event 0x715ad5ce61fc9595c7b415289d59cf203f23a94fa06f04af7e489a0a76e1fe26.
//
// Solidity: event RolesUpdated(address indexed user, uint256 indexed roles)
func (_CollateralOnramp *CollateralOnrampFilterer) WatchRolesUpdated(opts *bind.WatchOpts, sink chan<- *CollateralOnrampRolesUpdated, user []common.Address, roles []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var rolesRule []interface{}
	for _, rolesItem := range roles {
		rolesRule = append(rolesRule, rolesItem)
	}

	logs, sub, err := _CollateralOnramp.contract.WatchLogs(opts, "RolesUpdated", userRule, rolesRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollateralOnrampRolesUpdated)
				if err := _CollateralOnramp.contract.UnpackLog(event, "RolesUpdated", log); err != nil {
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
func (_CollateralOnramp *CollateralOnrampFilterer) ParseRolesUpdated(log types.Log) (*CollateralOnrampRolesUpdated, error) {
	event := new(CollateralOnrampRolesUpdated)
	if err := _CollateralOnramp.contract.UnpackLog(event, "RolesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CollateralOnrampUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the CollateralOnramp contract.
type CollateralOnrampUnpausedIterator struct {
	Event *CollateralOnrampUnpaused // Event containing the contract specifics and raw log

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
func (it *CollateralOnrampUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CollateralOnrampUnpaused)
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
		it.Event = new(CollateralOnrampUnpaused)
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
func (it *CollateralOnrampUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CollateralOnrampUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CollateralOnrampUnpaused represents a Unpaused event raised by the CollateralOnramp contract.
type CollateralOnrampUnpaused struct {
	Asset common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address indexed asset)
func (_CollateralOnramp *CollateralOnrampFilterer) FilterUnpaused(opts *bind.FilterOpts, asset []common.Address) (*CollateralOnrampUnpausedIterator, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _CollateralOnramp.contract.FilterLogs(opts, "Unpaused", assetRule)
	if err != nil {
		return nil, err
	}
	return &CollateralOnrampUnpausedIterator{contract: _CollateralOnramp.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address indexed asset)
func (_CollateralOnramp *CollateralOnrampFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *CollateralOnrampUnpaused, asset []common.Address) (event.Subscription, error) {

	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _CollateralOnramp.contract.WatchLogs(opts, "Unpaused", assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CollateralOnrampUnpaused)
				if err := _CollateralOnramp.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_CollateralOnramp *CollateralOnrampFilterer) ParseUnpaused(log types.Log) (*CollateralOnrampUnpaused, error) {
	event := new(CollateralOnrampUnpaused)
	if err := _CollateralOnramp.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
