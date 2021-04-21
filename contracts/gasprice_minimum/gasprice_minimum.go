// Copyright 2017 The Celo Authors
// This file is part of the celo library.
//
// The celo library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The celo library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the celo library. If not, see <http://www.gnu.org/licenses/>.

package gasprice_minimum

import (
	"math/big"
	"strings"

	"github.com/celo-org/celo-blockchain/accounts/abi"
	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/contracts"
	"github.com/celo-org/celo-blockchain/contracts/blockchain_parameters"
	"github.com/celo-org/celo-blockchain/core/vm"
	"github.com/celo-org/celo-blockchain/params"
)

// TODO (jarmg 5/22/19): Store ABIs in a central location
const (
	gasPriceMinimumABIString = `[
    {
      "constant": true,
      "inputs": [
        {
          "name": "_tokenAddress",
          "type": "address"
        }
      ],
      "name": "getGasPriceMinimum",
      "outputs": [
        {
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "_blockGasTotal",
          "type": "uint256"
        },
        {
          "name": "_blockGasLimit",
          "type": "uint256"
        }
      ],
      "name": "updateGasPriceMinimum",
      "outputs": [
        {
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
      } 
    
  ]`
)

var (
	FallbackGasPriceMinimum *big.Int = big.NewInt(0) // gas price minimum to return if unable to fetch from contract
	suggestionMultiplier    *big.Int = big.NewInt(5) // The multiplier that we apply to the minimum when suggesting gas price
)

var (
	getGasPriceMinimumMethod    *contracts.BoundMethod
	updateGasPriceMinimumMethod *contracts.BoundMethod
)

func init() {
	gasPriceMinimumABI, err := abi.JSON(strings.NewReader(gasPriceMinimumABIString))
	if err != nil {
		panic(err)
	}

	getGasPriceMinimumMethod = contracts.NewRegisteredContractMethod(params.GasPriceMinimumRegistryId, &gasPriceMinimumABI, "getGasPriceMinimum", params.MaxGasForGetGasPriceMinimum)
	updateGasPriceMinimumMethod = contracts.NewRegisteredContractMethod(params.GasPriceMinimumRegistryId, &gasPriceMinimumABI, "updateGasPriceMinimum", params.MaxGasForUpdateGasPriceMinimum)
}

func GetGasPriceSuggestion(vmRunner vm.EVMRunner, currency *common.Address) (*big.Int, error) {
	gasPriceMinimum, err := GetGasPriceMinimum(vmRunner, currency)
	return new(big.Int).Mul(gasPriceMinimum, suggestionMultiplier), err
}

func GetGasPriceMinimum(vmRunner vm.EVMRunner, currency *common.Address) (*big.Int, error) {
	var currencyAddress common.Address
	var err error

	if currency == nil {
		currencyAddress, err = contracts.GetRegisteredAddress(vmRunner, params.GoldTokenRegistryId)

		if err == contracts.ErrSmartContractNotDeployed || err == contracts.ErrRegistryContractNotDeployed {
			return FallbackGasPriceMinimum, nil
		}
		if err == contracts.ErrNoInternalEvmHandlerSingleton {
			return FallbackGasPriceMinimum, nil
		}
		if err != nil {
			return FallbackGasPriceMinimum, err
		}
	} else {
		currencyAddress = *currency
	}

	var gasPriceMinimum *big.Int
	err = getGasPriceMinimumMethod.Query(vmRunner, &gasPriceMinimum, currencyAddress)

	if err == contracts.ErrSmartContractNotDeployed || err == contracts.ErrRegistryContractNotDeployed {
		return FallbackGasPriceMinimum, nil
	}
	if err == contracts.ErrNoInternalEvmHandlerSingleton {
		return FallbackGasPriceMinimum, nil
	}
	if err != nil {
		return FallbackGasPriceMinimum, err
	}

	return gasPriceMinimum, err
}

func UpdateGasPriceMinimum(lastUsedGas uint64, vmRunner vm.EVMRunner) (*big.Int, error) {
	var updatedGasPriceMinimum *big.Int

	// If an error occurs, the default block gas limit will be returned and a log statement will be produced by contract_comm
	gasLimit, _ := blockchain_parameters.GetBlockGasLimit(vmRunner)

	err := updateGasPriceMinimumMethod.Execute(vmRunner, &updatedGasPriceMinimum, common.Big0, big.NewInt(int64(lastUsedGas)), big.NewInt(int64(gasLimit)))

	if err != nil {
		return nil, err
	}
	return updatedGasPriceMinimum, err
}
