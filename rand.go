package rng

import "crypto/rand"

// Uint64Bits generates a random uint64 value in range [0, 2^n). In other words
// returned uint64 will have n least significant bits set to random values,
// other bits will be set to 0.
//
// It will panic if there is error reading from crypto/rand source.
func Uint64Bits(n uint) (r uint64) {
	return ReadUint64Bits(rand.Reader, n)
}

// Intn returns a non negative int in [0, n).
// It will panic if n <= 0.
func Intn(n int) int {
	return ReadIntn(rand.Reader, n)
}

// Float64 returns a random number in [0.0,1.0)
func Float64() float64 {
	return ReadFloat64(rand.Reader)
}

// Perm returns, as a slice of n ints, a random permutation of the integers
// [0,n).
func Perm(n int) []int {
	return ReadPerm(rand.Reader, n)
}
