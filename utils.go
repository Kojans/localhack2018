package main

import (
	"encoding/binary"
	"math"
)

func float64ToByte(num float64) [8]byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], math.Float64bits(num))
	return buf
}
