package rng

import (
	"flag"
	"os"
	"testing"
)

var cfg struct {
	long bool
}

func TestMain(m *testing.M) {
	flag.BoolVar(&cfg.long, "long", false, "Enable long RNG tests")
	flag.Parse()
	os.Exit(m.Run())
}
