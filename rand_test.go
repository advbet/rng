package rng

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt64Overflow(t *testing.T) {
	// Test M overflow handling
	N := uint64(0x7fffffffffffffff) // max positive int
	bits := minBytes(N-1) * 8
	M := uint64(1 << bits)
	assert.Equal(t, uint64(0), M) // assert M overflows to 0

	assert.Equal(t, uint64(0x8000000000000001), M-N)
	limit := M - (M-N)%N
	assert.Equal(t, 2*N, limit)
	assert.Equal(t, M-2, limit)
}

func TestUint64BitsSourceError(t *testing.T) {
	origRand := rand.Reader
	defer func() {
		rand.Reader = origRand
	}()

	// This test replaces random source with source that would return error
	// on read. We test if error is converted to panic.

	rand.Reader = bytes.NewBuffer([]byte{})
	assert.Panics(t, func() {
		Uint64Bits(8)
	})
}

func TestUint64BitsRead(t *testing.T) {
	origRand := rand.Reader
	defer func() {
		rand.Reader = origRand
	}()

	// Test if Uint64Bits read from random source and read as little bytes
	// as possible. In this test we replace random source with fixed length
	// buffer containing number of bytes thats is enough for generating
	// expected length random value. If Uint64Bits reads more panic will
	// occur.

	tests := []struct {
		source []byte
		bits   uint
		value  uint64
	}{
		{[]byte{}, 0, 0},
		{[]byte{0x12}, 8, 0x12},
		{[]byte{0x12, 0x34}, 16, 0x3412},
		{[]byte{0x12, 0x34, 0x56}, 24, 0x563412},
		{[]byte{0x12, 0x34, 0x56, 0x78}, 32, 0x78563412},
		{[]byte{0x12, 0x34, 0x56, 0x78, 0x9a}, 40, 0x9a78563412},
		{[]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc}, 48, 0xbc9a78563412},
		{[]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde}, 56, 0xdebc9a78563412},
		{[]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0}, 64, 0xf0debc9a78563412},
	}

	for _, test := range tests {
		rand.Reader = bytes.NewBuffer(test.source)
		n := Uint64Bits(test.bits)
		assert.Equal(t, test.value, n)
	}
}

func TestUint64BitsMask(t *testing.T) {
	origRand := rand.Reader
	defer func() {
		rand.Reader = origRand
	}()

	// This test checks if Uint64Bits masks extra bits if nuber of bits does
	// not make last byte full.
	source := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	tests := []struct {
		bits  uint
		value uint64
	}{
		{0, 0x0},
		{1, 0x1},
		{2, 0x3},
		{3, 0x7},
		{4, 0xf},
		{5, 0x1f},
		{6, 0x3f},
		{7, 0x7f},
		{8, 0xff},
		{9, 0x1ff},
		{60, 0x0fffffffffffffff},
		{61, 0x1fffffffffffffff},
		{62, 0x3fffffffffffffff},
		{63, 0x7fffffffffffffff},
		{64, 0xffffffffffffffff},
	}

	for _, test := range tests {
		rand.Reader = bytes.NewBuffer(source)
		n := Uint64Bits(test.bits)
		assert.Equal(t, test.value, n)
	}
}

func TestMinBytes(t *testing.T) {
	tests := []struct {
		arg uint64
		val uint
	}{
		{0, 0},
		{0x1, 1},
		{0xff, 1},
		{0x100, 2},
		{0xffff, 2},
		{0x10000, 3},
		{0xffffff, 3},
		{0x1000000, 4},
		{0xffffffff, 4},
		{0x100000000, 5},
		{0xffffffffff, 5},
		{0x10000000000, 6},
		{0xffffffffffff, 6},
		{0x1000000000000, 7},
		{0xffffffffffffff, 7},
		{0x100000000000000, 8},
		{0x7fffffffffffffff, 8},
		{0xffffffffffffffff, 8},
	}

	for _, test := range tests {
		n := minBytes(test.arg)
		assert.Equal(t, test.val, n)
	}
}

func TestIntnPanics(t *testing.T) {
	assert.Panics(t, func() {
		Intn(0)
	})
	assert.Panics(t, func() {
		Intn(-1)
	})
}

func TestIntn(t *testing.T) {
	assert.Equal(t, 0, Intn(1))

	for i := uint(1); i <= 62; i++ {
		max := (1 << i) + 1
		n := Intn(max)
		assert.True(t, n < max, fmt.Sprintf("Intn(%d) = %d, must be < %d", max, n, max))
		assert.True(t, n >= 0, fmt.Sprintf("Intn(%d) = %d, must be >= 0", max, n))
	}
}

func TestFloat64(t *testing.T) {
	for i := 0; i < 10; i++ {
		r := Float64()
		assert.True(t, r < 1.0, fmt.Sprintf("Float64() = %g, must be < 1.0", r))
		assert.True(t, r >= 0.0, fmt.Sprintf("Float64() = %g, must be >= 0.0", r))
	}
}

func TestPerm(t *testing.T) {
	N := 20

	p := Perm(N)
	unique := make(map[int]bool)

	for _, val := range p {
		unique[val] = true

		// All values must be within range [0, N)
		assert.True(t, val >= 0 && val < N)
	}

	// All numbers from 0 to N-1 must occur only once
	assert.Len(t, unique, N)
}

func TestSample(t *testing.T) {
	tests := []struct {
		n int
		k int
	}{
		{10, 0},
		{10, 1},
		{10, 5},
		{10, 8},
		{10, 10},
	}

	for _, test := range tests {
		s := Sample(test.n, test.k)

		unique := make(map[int]struct{})
		for _, val := range s {
			unique[val] = struct{}{}
			assert.True(t, val >= 0 && val < test.n)
		}

		assert.Len(t, s, test.k)
		assert.Len(t, unique, test.k)
	}

	assert.Len(t, Sample(2, 10), 2)
}

func TestIntnFrequencyMonobit(t *testing.T) {
	if !cfg.long {
		t.Skip("skipping, run with --long to enable long RNG tests")
	}
	lim := 2
	sum := 0
	N := 1000
	Pmin := 0.01

	for i := 0; i < N; i++ {
		x := Intn(lim)
		if x < 0 || x >= lim {
			t.Errorf("Expected 0 <= Intn(%d) <= %d, got %d", lim, lim-1, x)
		}
		sum += x*2 - 1
	}
	Sobs := math.Abs(float64(sum)) / math.Sqrt(float64(N))
	P := math.Erfc(Sobs / math.Sqrt(2.0))
	t.Log("P = ", P)

	if P < Pmin {
		t.Errorf("sequence appears to be non-random, P = %f (< %f)", P, Pmin)
	}
}

func TestFloat64Sum(t *testing.T) {
	// This test computes sum of N random float64 variables in range
	// [-0.5; 0.5) and checks the result against expected sum distribution.
	if !cfg.long {
		t.Skip("skipping, run with --long to enable long RNG tests")
	}
	N := 1000
	Pmin := 0.01

	sum := 0.0
	for i := 0; i < N; i++ {
		sum += Float64() - 0.5
	}
	// V = (max-min)^2/12 = 1/12
	// V_sum = N*V
	// sigma_sum = sqrt(V_sum) = sqrt(N*V) = sqrt(N/12)
	Sobs := math.Abs(sum) / math.Sqrt(float64(N)/12.0)
	P := math.Erfc(Sobs / math.Sqrt(2.0))
	t.Log("P = ", P)

	if P < Pmin {
		t.Errorf("sequence appears to be non-random, P = %f (< %f)", P, Pmin)
	}
}

func TestIntnSum(t *testing.T) {
	// This test computes sum of N random integer variables in range [-5; 5]
	// and checks the result against expected sum distribution.
	if !cfg.long {
		t.Skip("skipping, run with --long to enable long RNG tests")
	}
	N := 1000
	Pmin := 0.01

	sum := 0
	for i := 0; i < N; i++ {
		// Naive RNG approach
		//sum += (int(Uint64Bits(8)) % 11) - 5
		sum += Intn(11) - 5
	}
	// V = ((max - min + 1)^2 - 1) / 12 = 10
	// V_sum = N*V
	// sigma_sum = sqrt(V_sum) = sqrt(N*V) = sqrt(N*10)
	Sobs := math.Abs(float64(sum)) / math.Sqrt(float64(N)*10.0)
	P := math.Erfc(Sobs / math.Sqrt(2.0))
	t.Log("P = ", P)

	if P < Pmin {
		t.Errorf("sequence appears to be non-random, P = %f (< %f)", P, Pmin)
	}
}

func TestIntnHistogram(t *testing.T) {
	if !cfg.long {
		t.Skip("skipping, run with --long to enable long RNG tests")
	}

	eps := 0.01
	N := 1000 * 1000
	max := 10
	hist := make([]int, max)

	for i := 0; i < N; i++ {
		// Naive RNG approach
		//hist[int(Uint64Bits(8))%max]++
		hist[Intn(max)]++
	}

	expected := 1.0 / float64(max)
	for i, count := range hist {
		actual := float64(count) / float64(N)
		assert.InEpsilon(t, expected, actual, eps, fmt.Sprintf("P(%d)_expected = %f, P(%d)_actual = %f", i, expected, i, actual))
	}
}

func TestFloat64Histogram(t *testing.T) {
	if !cfg.long {
		t.Skip("skipping, run with --long to enable long RNG tests")
	}

	eps := 0.01
	N := 1000 * 1000
	max := 10
	limits := make([]float64, max)
	hist := make([]int, max)

	for i := range limits {
		limits[i] = float64(i+1) / float64(max)
	}

	for i := 0; i < N; i++ {
		f := Float64()
		for n, limit := range limits {
			if f < limit {
				hist[n]++
				break
			}
		}
	}

	expected := 1.0 / float64(max)
	n := 0
	for i, count := range hist {
		n += count
		actual := float64(count) / float64(N)
		assert.InEpsilon(t, expected, actual, eps, fmt.Sprintf("P(%d)_expected = %f, P(%d)_actual = %f", i, expected, i, actual))
	}
	assert.Equal(t, N, n)
}
