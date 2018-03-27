package stake

//go:generate abigen --sol contract/StakeInterface.sol --pkg contract --out contract/stake.go

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/ens"
	ensInternal "github.com/ethereum/go-ethereum/contracts/ens/contract"
	"github.com/ethereum/go-ethereum/contracts/stake/contract"
)

var (
	MainNetENS = "stake.mainframehq.eth"
	TestNetENS = "stake.mainframe.test"
)

type Stake struct {
	*contract.StakeInterfaceCallerSession
	contractBackend bind.ContractBackend
}

func HasStakeENS(contractBackend bind.ContractBackend, ensContract common.Address, staker common.Address) (bool, error) {
  stakingContractENS := MainNetENS
  if ensContract != ens.MainNetAddress {
    stakingContractENS = TestNetENS
  }

	stakingContract, err := resolveStakingContractAddress(contractBackend, ensContract, stakingContractENS)
	if err != nil {
		return false, err
	}

	return HasStake(contractBackend, stakingContract, staker)
}

func HasStake(contractBackend bind.ContractBackend, stakingContract common.Address, staker common.Address) (bool, error) {
	caller, err := contract.NewStakeInterfaceCaller(stakingContract, contractBackend)
	if err != nil {
		return false, err
	}

	return caller.HasStake(nil, staker)
}

func resolveStakingContractAddress(contractBackend bind.ContractBackend, ensContract common.Address, stakingContractENS string) (common.Address, error) {
	ensResolver, err := ens.NewENS(&bind.TransactOpts{}, ensContract, contractBackend)
	if err != nil {
		return common.Address{}, err
	}

	node := ens.EnsNode(stakingContractENS)
	resolverContract, err := ensResolver.Resolver(node)
	if err != nil {
		return common.Address{}, err
	}

	caller, err := ensInternal.NewPublicResolverCaller(resolverContract, contractBackend)
	if err != nil {
		return common.Address{}, err
	}

	return caller.Addr(nil, node)
}
