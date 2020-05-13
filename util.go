package rng

// minBytes returns minimum number if bytes needed to store given number in
// binary form.
func minBytes(n uint64) (bytes uint) {
	switch {
	case n == 0:
		return 0
	case n <= 0xff:
		return 1
	case n <= 0xffff:
		return 2
	case n <= 0xffffff:
		return 3
	case n <= 0xffffffff:
		return 4
	case n <= 0xffffffffff:
		return 5
	case n <= 0xffffffffffff:
		return 6
	case n <= 0xffffffffffffff:
		return 7
	default:
		return 8
	}
}
