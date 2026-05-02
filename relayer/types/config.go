package types

import "math/big"

type ContractConfig struct {
	SafeFactory   string `json:"safe_factory"`
	SafeMultisend string `json:"safe_multisend"`
}

func GetContractConfig(chainId *big.Int) *ContractConfig {
	switch chainId.Int64() {
	case 137:
		return &ContractConfig{
			SafeFactory:   "0xaacFeEa03eb1561C4e67d661e40682Bd20E3541b",
			SafeMultisend: "0xA238CBeb142c10Ef7Ad8442C6D1f9E89e07e7761",
		}
	default:
		panic("invalid chain id")
	}
}

type DepositWalletContractConfig struct {
	DepositWalletFactory        string
	DepositWalletImplementation string
}

func GetDepositWalletContractConfig(chainId *big.Int) *DepositWalletContractConfig {
	switch chainId.Int64() {
	case 137:
		return &DepositWalletContractConfig{
			DepositWalletFactory:        "0x00000000000Fb5C9ADea0298D729A0CB3823Cc07",
			DepositWalletImplementation: "0x58CA52ebe0DadfdF531Cde7062e76746de4Db1eB",
		}
	default:
		panic("invalid chain id")
	}
}
