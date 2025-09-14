package main

import "math/big"

func encodeDecimal(input string) string {
	data := []byte(input)

	num := new(big.Int).SetBytes(data)

	return num.String()
}
