package random

import (
	"math/big"

	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/common/hexutil"
	"github.com/celo-org/celo-blockchain/contracts"
	"github.com/celo-org/celo-blockchain/contracts/abis"
	"github.com/celo-org/celo-blockchain/contracts/config"
	"github.com/celo-org/celo-blockchain/contracts/internal/n"
	"github.com/celo-org/celo-blockchain/core/vm"
	"github.com/celo-org/celo-blockchain/log"
)

const (
	maxGasForCommitments       uint64 = 2 * n.Million
	maxGasForComputeCommitment uint64 = 2 * n.Million
	maxGasForBlockRandomness   uint64 = 2 * n.Million
)

var (
	commitmentsMethod        = contracts.NewRegisteredContractMethod(config.RandomRegistryId, abis.Random, "commitments", maxGasForCommitments)
	computeCommitmentMethod  = contracts.NewRegisteredContractMethod(config.RandomRegistryId, abis.Random, "computeCommitment", maxGasForComputeCommitment)
	randomMethod             = contracts.NewRegisteredContractMethod(config.RandomRegistryId, abis.Random, "random", maxGasForBlockRandomness)
	getBlockRandomnessMethod = contracts.NewRegisteredContractMethod(config.RandomRegistryId, abis.Random, "getBlockRandomness", maxGasForBlockRandomness)
)

func IsRunning(vmRunner vm.EVMRunner) bool {
	randomAddress, err := contracts.GetRegisteredAddress(vmRunner, config.RandomRegistryId)

	if err == contracts.ErrSmartContractNotDeployed || err == contracts.ErrRegistryContractNotDeployed {
		log.Debug("Registry address lookup failed", "err", err, "contract", hexutil.Encode(config.RandomRegistryId[:]))
	} else if err != nil {
		log.Error(err.Error())
	}

	return err == nil && randomAddress != common.ZeroAddress
}

// GetLastCommitment returns up the last commitment in the smart contract
func GetLastCommitment(vmRunner vm.EVMRunner, validator common.Address) (common.Hash, error) {
	lastCommitment := common.Hash{}
	err := commitmentsMethod.Query(vmRunner, &lastCommitment, validator)
	if err != nil {
		log.Error("Failed to get last commitment", "err", err)
		return lastCommitment, err
	}

	if (lastCommitment == common.Hash{}) {
		log.Debug("Unable to find last randomness commitment in smart contract")
	}

	return lastCommitment, nil
}

// ComputeCommitment calulcates the commitment for a given randomness.
func ComputeCommitment(vmRunner vm.EVMRunner, randomness common.Hash) (common.Hash, error) {
	commitment := common.Hash{}
	// TODO(asa): Make an issue to not have to do this via StaticCall
	err := computeCommitmentMethod.Query(vmRunner, &commitment, randomness)
	if err != nil {
		log.Error("Failed to call computeCommitment()", "err", err)
		return common.Hash{}, err
	}

	return commitment, err
}

// Random performs an internal call to the EVM to retrieve the current randomness from the official Random contract.
func Random(vmRunner vm.EVMRunner) (common.Hash, error) {
	randomness := common.Hash{}
	err := randomMethod.Query(vmRunner, &randomness)
	return randomness, err
}

func BlockRandomness(vmRunner vm.EVMRunner, blockNumber uint64) (common.Hash, error) {
	randomness := common.Hash{}
	err := getBlockRandomnessMethod.Query(vmRunner, &randomness, big.NewInt(int64(blockNumber)))
	return randomness, err
}
