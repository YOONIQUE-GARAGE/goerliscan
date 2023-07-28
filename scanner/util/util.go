package util

import (
	"math/big"
)

func WeiToEth(wei *big.Int) *big.Float {
	if wei == nil {
		// Return a new *big.Float with value 0.0
		return new(big.Float).SetFloat64(0.0)
	}
	// Create a big.Float with 18 decimal places
	decimal := new(big.Float).SetInt64(1e18)
	weiToInt := new(big.Float).SetInt(wei)
	// Convert wei to ether by dividing by 10^18
	ethValue := new(big.Float).Quo(weiToInt, decimal)

	return ethValue
}
