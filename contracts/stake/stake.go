package stake

//go:generate abigen --sol contract/HasStake.sol --pkg contract --out contract/stake.go

import (
  "github.com/ethereum/go-ethereum/accounts/abi/bind"
  "github.com/ethereum/go-ethereum/common"
  "github.com/ethereum/go-ethereum/contracts/stake/contract"
)

var (
  TestNetAddress = common.HexToAddress("0x42e4fb08443f93c736f3bb6544ea08628eb1c30f")
)

type Stake struct {
  *contract.HasStakeCallerSession
  contractBackend bind.ContractBackend
}

func HasStake(contractBackend bind.ContractBackend, address common.Address) (bool, error) {
  caller, err := contract.NewHasStakeCaller(TestNetAddress, contractBackend)
  if err != nil {
    return false, err
  }

  return caller.HasStake(nil, address)
}
