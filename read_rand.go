package rng

import (
	"io"
)

// ReadUint64Bits reads a random uint64 value in range [0, 2^n) from a random
// source. In other words returned uint64 will have n least significant bits set
// to random values, other bits will be set to 0.
//
// It will panic if random source returns read error.
func ReadUint64Bits(src io.Reader, n uint) (r uint64) {
	if n > 64 {
		panic("abr.Uint64Bits can not return more than 64 random bits")
	}

	bytes := (n + 7) / 8 // number of bytes to read from random source
	b := make([]byte, 8) // initial value will be all zeros
	if _, e := io.ReadFull(src, b[:bytes]); e != nil {
		panic(e)
	}
	// treat b as little endian value
	r = uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
	// mask extra bits
	return r & ((1 << n) - 1)
}

// ReadIntn returns a non negative int in [0, n) reading randomness from a given
// source. It will panic if n <= 0.
func ReadIntn(src io.Reader, n int) int {
	if n <= 0 {
		panic("invalid argument to Intn")
	}

	N := uint64(n)
	// minimum number of random bits that will be read be read from entropy
	// source, it is always a multiple of 8 because reads have byte
	// granularity
	bits := minBytes(N-1) * 8
	// if N is a power of two single read is always sufficient
	if N&(N-1) == 0 {
		return int(ReadUint64Bits(src, bits) & (N - 1))
	}

	// call Uint64Bits(bits) will always return values in range [0; M)
	// bits=8, M=0x100
	// bits=16, M=0x10000
	// ...
	// bits=56, M=0x100000000000000
	// bits=64, M=0x0, because of unint64 overflow, its OK
	M := uint64(1 << bits)
	// limit is upper bound (non inclusive) of values that can be used to
	// generate random number in range [0; N)
	// we could have limit set to (M - M % N) but it would produce
	// invalid results with M=0, thus we manually reduce single modulus step
	// by using (M - N) instead of M.
	limit := M - (M-N)%N
	// Example: N=65, bits=8, M=256
	// we can expand 256 values to three full intervals of 65 elements
	//
	// 0        65       130      195    256
	// |        |        |        |      |
	// |<0---64>|<0---64>|<0---64>|unused|
	//
	// M / N = 3, number of full intervals of 65 elements
	// M % N = 61, number of unused values
	//
	// in this case limit will be 195, if we draw a number in range [0; 195)
	// we can use it to reduce it to [0; 65) if we get number >= 195 we try
	// drawing again.
	for {
		r := ReadUint64Bits(src, bits)
		if r < limit {
			return int(r % N)
		}
	}
}

// ReadFloat64 returns a random number in [0.0,1.0) reading randomness from a
// given source.
func ReadFloat64(src io.Reader) float64 {
	return float64(ReadUint64Bits(src, 53)) / float64(1<<53)
}

// ReadPerm returns, as a slice of n ints, a random permutation of the integers
// [0,n) reading randomness from a given source.
func ReadPerm(src io.Reader, n int) []int {
	m := make([]int, n)
	for i := 0; i < n; i++ {
		j := ReadIntn(src, i+1)
		m[i] = m[j]
		m[j] = i
	}
	return m
}

// ReadSample returns random k integers from a range [0 n). If k > n then only n
// integers are returned.
//
// This function consumes entropy from a given entroy source src.
func ReadSample(src io.Reader, n int, k int) []int {
	if k > n {
		k = n
	}

	if k > n/2 {
		return ReadPerm(src, n)[0:k]
	}

	sample := make([]int, 0, k)
	cache := make(map[int]struct{})
	for len(sample) < k {
		r := ReadIntn(src, n)
		if _, ok := cache[r]; ok {
			continue
		}
		cache[r] = struct{}{}
		sample = append(sample, r)
	}
	return sample
}
