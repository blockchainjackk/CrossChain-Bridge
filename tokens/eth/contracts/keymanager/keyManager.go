// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package keymanager

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
)

// KeymanagerMetaData contains all meta data concerning the Keymanager contract.
var KeymanagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_priKeys\",\"type\":\"string[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_priKey\",\"type\":\"string\"}],\"name\":\"addPriKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAll\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"num\",\"type\":\"uint8\"}],\"name\":\"getPriKeys\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"_priKeys\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"setOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"oldKey\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"newKey\",\"type\":\"string\"}],\"name\":\"updatePriKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b506040516200182a3803806200182a833981810160405281019062000037919062000464565b60008151116200007e576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401620000759062000516565b60405180910390fd5b336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508060019080519060200190620000d6929190620000de565b50506200086a565b8280548282559060005260206000209081019282156200012b579160200282015b828111156200012a57825182908162000119919062000783565b5091602001919060010190620000ff565b5b5090506200013a91906200013e565b5090565b5b8082111562000162576000818162000158919062000166565b506001016200013f565b5090565b508054620001749062000572565b6000825580601f10620001885750620001a9565b601f016020900490600052602060002090810190620001a89190620001ac565b5b50565b5b80821115620001c7576000816000905550600101620001ad565b5090565b6000604051905090565b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6200022f82620001e4565b810181811067ffffffffffffffff82111715620002515762000250620001f5565b5b80604052505050565b600062000266620001cb565b905062000274828262000224565b919050565b600067ffffffffffffffff821115620002975762000296620001f5565b5b602082029050602081019050919050565b600080fd5b600080fd5b600067ffffffffffffffff821115620002d057620002cf620001f5565b5b620002db82620001e4565b9050602081019050919050565b60005b8381101562000308578082015181840152602081019050620002eb565b8381111562000318576000848401525b50505050565b6000620003356200032f84620002b2565b6200025a565b905082815260208101848484011115620003545762000353620002ad565b5b62000361848285620002e8565b509392505050565b600082601f830112620003815762000380620001df565b5b8151620003938482602086016200031e565b91505092915050565b6000620003b3620003ad8462000279565b6200025a565b90508083825260208201905060208402830185811115620003d957620003d8620002a8565b5b835b818110156200042757805167ffffffffffffffff811115620004025762000401620001df565b5b80860162000411898262000369565b85526020850194505050602081019050620003db565b5050509392505050565b600082601f830112620004495762000448620001df565b5b81516200045b8482602086016200039c565b91505092915050565b6000602082840312156200047d576200047c620001d5565b5b600082015167ffffffffffffffff8111156200049e576200049d620001da565b5b620004ac8482850162000431565b91505092915050565b600082825260208201905092915050565b7f70726976617465204b6579207265717569726564000000000000000000000000600082015250565b6000620004fe601483620004b5565b91506200050b82620004c6565b602082019050919050565b600060208201905081810360008301526200053181620004ef565b9050919050565b600081519050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600060028204905060018216806200058b57607f821691505b602082108103620005a157620005a062000543565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b6000600883026200060b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82620005cc565b620006178683620005cc565b95508019841693508086168417925050509392505050565b6000819050919050565b6000819050919050565b6000620006646200065e62000658846200062f565b62000639565b6200062f565b9050919050565b6000819050919050565b620006808362000643565b620006986200068f826200066b565b848454620005d9565b825550505050565b600090565b620006af620006a0565b620006bc81848462000675565b505050565b5b81811015620006e457620006d8600082620006a5565b600181019050620006c2565b5050565b601f8211156200073357620006fd81620005a7565b6200070884620005bc565b8101602085101562000718578190505b620007306200072785620005bc565b830182620006c1565b50505b505050565b600082821c905092915050565b6000620007586000198460080262000738565b1980831691505092915050565b600062000773838362000745565b9150826002028217905092915050565b6200078e8262000538565b67ffffffffffffffff811115620007aa57620007a9620001f5565b5b620007b6825462000572565b620007c3828285620006e8565b600060209050601f831160018114620007fb5760008415620007e6578287015190505b620007f2858262000765565b86555062000862565b601f1984166200080b86620005a7565b60005b8281101562000835578489015182556001820191506020850194506020810190506200080e565b8683101562000855578489015162000851601f89168262000745565b8355505b6001600288020188555050505b505050505050565b610fb0806200087a6000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c806313af403514610067578063394680891461008357806353ed51431461009f5780636557fa87146100bd5780638da5cb5b146100ed578063d08c143b1461010b575b600080fd5b610081600480360381019061007c91906107f2565b610127565b005b61009d600480360381019061009891906109c5565b610267565b005b6100a76102fa565b6040516100b49190610b8b565b60405180910390f35b6100d760048036038101906100d29190610be6565b610461565b6040516100e49190610b8b565b60405180910390f35b6100f56106cb565b6040516101029190610c22565b60405180910390f35b61012560048036038101906101209190610c3d565b6106ef565b005b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146101b5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101ac90610ce3565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610224576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161021b90610d4f565b60405180910390fd5b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146102f5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102ec90610ce3565b60405180910390fd5b505050565b606060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461038a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161038190610ce3565b60405180910390fd5b6001805480602002602001604051908101604052809291908181526020016000905b828210156104585783829060005260206000200180546103cb90610d9e565b80601f01602080910402602001604051908101604052809291908181526020018280546103f790610d9e565b80156104445780601f1061041957610100808354040283529160200191610444565b820191906000526020600020905b81548152906001019060200180831161042757829003601f168201915b5050505050815260200190600101906103ac565b50505050905090565b606060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146104f1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104e890610ce3565b60405180910390fd5b6001805490508260ff16111561053c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161053390610e41565b60405180910390fd5b60ff8260ff161115610583576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161057a90610ead565b60405180910390fd5b8160ff1667ffffffffffffffff8111156105a05761059f61089a565b5b6040519080825280602002602001820160405280156105d357816020015b60608152602001906001900390816105be5790505b50905060005b8260ff168160ff1610156106c55760018160ff16815481106105fe576105fd610ecd565b5b90600052602060002001805461061390610d9e565b80601f016020809104026020016040519081016040528092919081815260200182805461063f90610d9e565b801561068c5780601f106106615761010080835404028352916020019161068c565b820191906000526020600020905b81548152906001019060200180831161066f57829003601f168201915b5050505050828260ff16815181106106a7576106a6610ecd565b5b602002602001018190525080806106bd90610f2b565b9150506105d9565b50919050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461077d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161077490610ce3565b60405180910390fd5b50565b6000604051905090565b600080fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006107bf82610794565b9050919050565b6107cf816107b4565b81146107da57600080fd5b50565b6000813590506107ec816107c6565b92915050565b6000602082840312156108085761080761078a565b5b6000610816848285016107dd565b91505092915050565b600080fd5b600080fd5b600080fd5b60008083601f8401126108445761084361081f565b5b8235905067ffffffffffffffff81111561086157610860610824565b5b60208301915083600182028301111561087d5761087c610829565b5b9250929050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6108d282610889565b810181811067ffffffffffffffff821117156108f1576108f061089a565b5b80604052505050565b6000610904610780565b905061091082826108c9565b919050565b600067ffffffffffffffff8211156109305761092f61089a565b5b61093982610889565b9050602081019050919050565b82818337600083830152505050565b600061096861096384610915565b6108fa565b90508281526020810184848401111561098457610983610884565b5b61098f848285610946565b509392505050565b600082601f8301126109ac576109ab61081f565b5b81356109bc848260208601610955565b91505092915050565b6000806000604084860312156109de576109dd61078a565b5b600084013567ffffffffffffffff8111156109fc576109fb61078f565b5b610a088682870161082e565b9350935050602084013567ffffffffffffffff811115610a2b57610a2a61078f565b5b610a3786828701610997565b9150509250925092565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610aa7578082015181840152602081019050610a8c565b83811115610ab6576000848401525b50505050565b6000610ac782610a6d565b610ad18185610a78565b9350610ae1818560208601610a89565b610aea81610889565b840191505092915050565b6000610b018383610abc565b905092915050565b6000602082019050919050565b6000610b2182610a41565b610b2b8185610a4c565b935083602082028501610b3d85610a5d565b8060005b85811015610b795784840389528151610b5a8582610af5565b9450610b6583610b09565b925060208a01995050600181019050610b41565b50829750879550505050505092915050565b60006020820190508181036000830152610ba58184610b16565b905092915050565b600060ff82169050919050565b610bc381610bad565b8114610bce57600080fd5b50565b600081359050610be081610bba565b92915050565b600060208284031215610bfc57610bfb61078a565b5b6000610c0a84828501610bd1565b91505092915050565b610c1c816107b4565b82525050565b6000602082019050610c376000830184610c13565b92915050565b600060208284031215610c5357610c5261078a565b5b600082013567ffffffffffffffff811115610c7157610c7061078f565b5b610c7d84828501610997565b91505092915050565b600082825260208201905092915050565b7f6e6f74206f776e65720000000000000000000000000000000000000000000000600082015250565b6000610ccd600983610c86565b9150610cd882610c97565b602082019050919050565b60006020820190508181036000830152610cfc81610cc0565b9050919050565b7f696e76616c696420616464726573730000000000000000000000000000000000600082015250565b6000610d39600f83610c86565b9150610d4482610d03565b602082019050919050565b60006020820190508181036000830152610d6881610d2c565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680610db657607f821691505b602082108103610dc957610dc8610d6f565b5b50919050565b7f6e756d2069732067726561746572207468616e20746865206e756d626572206f60008201527f662070726976617465206b657973000000000000000000000000000000000000602082015250565b6000610e2b602e83610c86565b9150610e3682610dcf565b604082019050919050565b60006020820190508181036000830152610e5a81610e1e565b9050919050565b7f6e756d2069732067726561746572207468616e20323535000000000000000000600082015250565b6000610e97601783610c86565b9150610ea282610e61565b602082019050919050565b60006020820190508181036000830152610ec681610e8a565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610f3682610bad565b915060ff8203610f4957610f48610efc565b5b60018201905091905056fea264697066735822122020a1099cc8eb8dfe9ff4dbb0e1bf9d5471759ddf52e23d3dcb81d5c1b239da9d64736f6c637828302e382e31362d646576656c6f702e323032322e372e31362b636f6d6d69742e38303030383865330059",
}

// KeymanagerABI is the input ABI used to generate the binding from.
// Deprecated: Use KeymanagerMetaData.ABI instead.
var KeymanagerABI = KeymanagerMetaData.ABI

// KeymanagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use KeymanagerMetaData.Bin instead.
var KeymanagerBin = KeymanagerMetaData.Bin

// DeployKeymanager deploys a new Ethereum contract, binding an instance of Keymanager to it.
func DeployKeymanager(auth *bind.TransactOpts, backend bind.ContractBackend, _priKeys []string) (common.Address, *types.Transaction, *Keymanager, error) {
	parsed, err := KeymanagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(KeymanagerBin), backend, _priKeys)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Keymanager{KeymanagerCaller: KeymanagerCaller{contract: contract}, KeymanagerTransactor: KeymanagerTransactor{contract: contract}, KeymanagerFilterer: KeymanagerFilterer{contract: contract}}, nil
}

// Keymanager is an auto generated Go binding around an Ethereum contract.
type Keymanager struct {
	KeymanagerCaller     // Read-only binding to the contract
	KeymanagerTransactor // Write-only binding to the contract
	KeymanagerFilterer   // Log filterer for contract events
}

// KeymanagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type KeymanagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeymanagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type KeymanagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeymanagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KeymanagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KeymanagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KeymanagerSession struct {
	Contract     *Keymanager       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KeymanagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KeymanagerCallerSession struct {
	Contract *KeymanagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// KeymanagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KeymanagerTransactorSession struct {
	Contract     *KeymanagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// KeymanagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type KeymanagerRaw struct {
	Contract *Keymanager // Generic contract binding to access the raw methods on
}

// KeymanagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KeymanagerCallerRaw struct {
	Contract *KeymanagerCaller // Generic read-only contract binding to access the raw methods on
}

// KeymanagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KeymanagerTransactorRaw struct {
	Contract *KeymanagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewKeymanager creates a new instance of Keymanager, bound to a specific deployed contract.
func NewKeymanager(address common.Address, backend bind.ContractBackend) (*Keymanager, error) {
	contract, err := bindKeymanager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Keymanager{KeymanagerCaller: KeymanagerCaller{contract: contract}, KeymanagerTransactor: KeymanagerTransactor{contract: contract}, KeymanagerFilterer: KeymanagerFilterer{contract: contract}}, nil
}

// NewKeymanagerCaller creates a new read-only instance of Keymanager, bound to a specific deployed contract.
func NewKeymanagerCaller(address common.Address, caller bind.ContractCaller) (*KeymanagerCaller, error) {
	contract, err := bindKeymanager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KeymanagerCaller{contract: contract}, nil
}

// NewKeymanagerTransactor creates a new write-only instance of Keymanager, bound to a specific deployed contract.
func NewKeymanagerTransactor(address common.Address, transactor bind.ContractTransactor) (*KeymanagerTransactor, error) {
	contract, err := bindKeymanager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KeymanagerTransactor{contract: contract}, nil
}

// NewKeymanagerFilterer creates a new log filterer instance of Keymanager, bound to a specific deployed contract.
func NewKeymanagerFilterer(address common.Address, filterer bind.ContractFilterer) (*KeymanagerFilterer, error) {
	contract, err := bindKeymanager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KeymanagerFilterer{contract: contract}, nil
}

// bindKeymanager binds a generic wrapper to an already deployed contract.
func bindKeymanager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(KeymanagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Keymanager *KeymanagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Keymanager.Contract.KeymanagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Keymanager *KeymanagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Keymanager.Contract.KeymanagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Keymanager *KeymanagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Keymanager.Contract.KeymanagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Keymanager *KeymanagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Keymanager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Keymanager *KeymanagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Keymanager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Keymanager *KeymanagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Keymanager.Contract.contract.Transact(opts, method, params...)
}

// GetAll is a free data retrieval call binding the contract method 0x53ed5143.
//
// Solidity: function getAll() view returns(string[])
func (_Keymanager *KeymanagerCaller) GetAll(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _Keymanager.contract.Call(opts, &out, "getAll")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetAll is a free data retrieval call binding the contract method 0x53ed5143.
//
// Solidity: function getAll() view returns(string[])
func (_Keymanager *KeymanagerSession) GetAll() ([]string, error) {
	return _Keymanager.Contract.GetAll(&_Keymanager.CallOpts)
}

// GetAll is a free data retrieval call binding the contract method 0x53ed5143.
//
// Solidity: function getAll() view returns(string[])
func (_Keymanager *KeymanagerCallerSession) GetAll() ([]string, error) {
	return _Keymanager.Contract.GetAll(&_Keymanager.CallOpts)
}

// GetPriKeys is a free data retrieval call binding the contract method 0x6557fa87.
//
// Solidity: function getPriKeys(uint8 num) view returns(string[] _priKeys)
func (_Keymanager *KeymanagerCaller) GetPriKeys(opts *bind.CallOpts, num uint8) ([]string, error) {
	var out []interface{}
	err := _Keymanager.contract.Call(opts, &out, "getPriKeys", num)

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetPriKeys is a free data retrieval call binding the contract method 0x6557fa87.
//
// Solidity: function getPriKeys(uint8 num) view returns(string[] _priKeys)
func (_Keymanager *KeymanagerSession) GetPriKeys(num uint8) ([]string, error) {
	return _Keymanager.Contract.GetPriKeys(&_Keymanager.CallOpts, num)
}

// GetPriKeys is a free data retrieval call binding the contract method 0x6557fa87.
//
// Solidity: function getPriKeys(uint8 num) view returns(string[] _priKeys)
func (_Keymanager *KeymanagerCallerSession) GetPriKeys(num uint8) ([]string, error) {
	return _Keymanager.Contract.GetPriKeys(&_Keymanager.CallOpts, num)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Keymanager *KeymanagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Keymanager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Keymanager *KeymanagerSession) Owner() (common.Address, error) {
	return _Keymanager.Contract.Owner(&_Keymanager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Keymanager *KeymanagerCallerSession) Owner() (common.Address, error) {
	return _Keymanager.Contract.Owner(&_Keymanager.CallOpts)
}

// AddPriKey is a paid mutator transaction binding the contract method 0xd08c143b.
//
// Solidity: function addPriKey(string _priKey) returns()
func (_Keymanager *KeymanagerTransactor) AddPriKey(opts *bind.TransactOpts, _priKey string) (*types.Transaction, error) {
	return _Keymanager.contract.Transact(opts, "addPriKey", _priKey)
}

// AddPriKey is a paid mutator transaction binding the contract method 0xd08c143b.
//
// Solidity: function addPriKey(string _priKey) returns()
func (_Keymanager *KeymanagerSession) AddPriKey(_priKey string) (*types.Transaction, error) {
	return _Keymanager.Contract.AddPriKey(&_Keymanager.TransactOpts, _priKey)
}

// AddPriKey is a paid mutator transaction binding the contract method 0xd08c143b.
//
// Solidity: function addPriKey(string _priKey) returns()
func (_Keymanager *KeymanagerTransactorSession) AddPriKey(_priKey string) (*types.Transaction, error) {
	return _Keymanager.Contract.AddPriKey(&_Keymanager.TransactOpts, _priKey)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(address _newOwner) returns()
func (_Keymanager *KeymanagerTransactor) SetOwner(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _Keymanager.contract.Transact(opts, "setOwner", _newOwner)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(address _newOwner) returns()
func (_Keymanager *KeymanagerSession) SetOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _Keymanager.Contract.SetOwner(&_Keymanager.TransactOpts, _newOwner)
}

// SetOwner is a paid mutator transaction binding the contract method 0x13af4035.
//
// Solidity: function setOwner(address _newOwner) returns()
func (_Keymanager *KeymanagerTransactorSession) SetOwner(_newOwner common.Address) (*types.Transaction, error) {
	return _Keymanager.Contract.SetOwner(&_Keymanager.TransactOpts, _newOwner)
}

// UpdatePriKey is a paid mutator transaction binding the contract method 0x39468089.
//
// Solidity: function updatePriKey(string oldKey, string newKey) returns()
func (_Keymanager *KeymanagerTransactor) UpdatePriKey(opts *bind.TransactOpts, oldKey string, newKey string) (*types.Transaction, error) {
	return _Keymanager.contract.Transact(opts, "updatePriKey", oldKey, newKey)
}

// UpdatePriKey is a paid mutator transaction binding the contract method 0x39468089.
//
// Solidity: function updatePriKey(string oldKey, string newKey) returns()
func (_Keymanager *KeymanagerSession) UpdatePriKey(oldKey string, newKey string) (*types.Transaction, error) {
	return _Keymanager.Contract.UpdatePriKey(&_Keymanager.TransactOpts, oldKey, newKey)
}

// UpdatePriKey is a paid mutator transaction binding the contract method 0x39468089.
//
// Solidity: function updatePriKey(string oldKey, string newKey) returns()
func (_Keymanager *KeymanagerTransactorSession) UpdatePriKey(oldKey string, newKey string) (*types.Transaction, error) {
	return _Keymanager.Contract.UpdatePriKey(&_Keymanager.TransactOpts, oldKey, newKey)
}
