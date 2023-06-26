// Copyright 2021 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package misc

import (
	"math/big"

	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/common/math"
	gpm "github.com/celo-org/celo-blockchain/contracts/gasprice_minimum"
	"github.com/celo-org/celo-blockchain/core/types"
	"github.com/celo-org/celo-blockchain/core/vm"
	"github.com/celo-org/celo-blockchain/log"
	"github.com/celo-org/celo-blockchain/params"
)

// // VerifyEip1559Header verifies some header attributes which were changed in EIP-1559,
// // - gas limit check
// // - basefee check
// func VerifyEip1559Header(config *params.ChainConfig, parent, header *types.Header) error {
// 	// Verify that the gas limit remains within allowed bounds
// 	parentGasLimit := parent.GasLimit
// 	if !config.IsLondon(parent.Number) {
// 		parentGasLimit = parent.GasLimit * params.ElasticityMultiplier
// 	}
// 	if err := VerifyGaslimit(parentGasLimit, header.GasLimit); err != nil {
// 		return err
// 	}
// 	// Verify the header is not malformed
// 	if header.BaseFee == nil {
// 		return fmt.Errorf("header is missing baseFee")
// 	}
// 	// Verify the baseFee is correct based on the parent header.
// 	expectedBaseFee := CalcBaseFee(config, parent)
// 	if header.BaseFee.Cmp(expectedBaseFee) != 0 {
// 		return fmt.Errorf("invalid baseFee: have %s, want %s, parentBaseFee %s, parentGasUsed %d",
// 			expectedBaseFee, header.BaseFee, parent.BaseFee, parent.GasUsed)
// 	}
// 	return nil
// }

// CalcBaseFee calculates the basefee for the header.
// If the gasPriceMinimum contract fails to retrieve the new gas price minimum, it uses the
// ethereum's default baseFee calculation
func CalcBaseFee(config *params.ChainConfig, parent *types.Header, vmRunner vm.EVMRunner) *big.Int {
	var newBaseFee *big.Int
	var err error
	// If the current block is the first GFork block, return the gasPriceMinimum on the state.
	if !config.IsGFork(parent.Number) {
		newBaseFee, err = gpm.GetRealGasPriceMinimum(vmRunner, nil)
		if err != nil {
			log.Warn("CalcBaseFee error, contract call getPriceMinimumMethod", "error", err, "header.Number", parent.Number.Uint64()+1)
			// Will return the initialBaseFee param, but this way, the parameters will be isolated to the
			// ethereum functions
			newBaseFee = calcBaseFeeEthereum(config, parent)
		}
		return newBaseFee
	}

	// // If an error occurs, the default block gas limit will be returned and a log statement will be produced by GetBlockGasLimitOrDefault
	// gasLimit, err := blockchain_parameters.GetBlockGasLimit(vmRunner)
	// if err != nil {
	// 	return nil, err
	// }

	newBaseFee, err = gpm.GetUpdatedGasPriceMinimum(vmRunner, parent.GasUsed, parent.GasLimit)
	if err != nil {
		log.Warn("CalcBaseFee error, contract call UpdatedGasPriceMinimum", "error", err, "header.Number", parent.Number.Uint64()+1)
		newBaseFee = calcBaseFeeEthereum(config, parent)
	}
	return newBaseFee
}

// CalcBaseFee calculates the basefee of the header.
func calcBaseFeeEthereum(config *params.ChainConfig, parent *types.Header) *big.Int {
	// If the current block is the first EIP-1559 block, return the InitialBaseFee.
	if !config.IsLondon(parent.Number) {
		return new(big.Int).SetUint64(params.InitialBaseFee)
	}

	var (
		parentGasTarget          = parent.GasLimit / params.ElasticityMultiplier
		parentGasTargetBig       = new(big.Int).SetUint64(parentGasTarget)
		baseFeeChangeDenominator = new(big.Int).SetUint64(params.BaseFeeChangeDenominator)
	)
	// If the parent gasUsed is the same as the target, the baseFee remains unchanged.
	if parent.GasUsed == parentGasTarget {
		return new(big.Int).Set(parent.BaseFee)
	}
	if parent.GasUsed > parentGasTarget {
		// If the parent block used more gas than its target, the baseFee should increase.
		gasUsedDelta := new(big.Int).SetUint64(parent.GasUsed - parentGasTarget)
		x := new(big.Int).Mul(parent.BaseFee, gasUsedDelta)
		y := x.Div(x, parentGasTargetBig)
		baseFeeDelta := math.BigMax(
			x.Div(y, baseFeeChangeDenominator),
			common.Big1,
		)

		return x.Add(parent.BaseFee, baseFeeDelta)
	} else {
		// Otherwise if the parent block used less gas than its target, the baseFee should decrease.
		gasUsedDelta := new(big.Int).SetUint64(parentGasTarget - parent.GasUsed)
		x := new(big.Int).Mul(parent.BaseFee, gasUsedDelta)
		y := x.Div(x, parentGasTargetBig)
		baseFeeDelta := x.Div(y, baseFeeChangeDenominator)

		return math.BigMax(
			x.Sub(parent.BaseFee, baseFeeDelta),
			common.Big0,
		)
	}
}
