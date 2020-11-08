package codewords

import (
	"fmt"
	"testing"
)

func TestNewFactory(t *testing.T) {
	cwf := NewFactory()
	cw := cwf.Generate()
	if len(cw) == 0 {
		t.Error("generated empty string")
	}
}

func TestSamples(t *testing.T) {
	cwf := NewFactory()
	for i := 0; i < 16; i++ {
		fmt.Println(cwf.Generate())
	}
}
