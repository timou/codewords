package codewords

import (
	"fmt"
	"math/rand"
	"time"
)

// Factory builds codewords.
type Factory struct {
	random     *rand.Rand
	adjectives []string
	nouns      []string
}

// NewFactory creates a new codeword Factory.
func NewFactory() *Factory {
	return &Factory{
		random:     rand.New(rand.NewSource(time.Now().UTC().UnixNano())),
		adjectives: adjectives,
		nouns:      nouns,
	}
}

// Generate creates a new codeword of the form <adjective>-<noun>.
// They are not guaranteed to be unique. They may be offensive or
// inappropriate for the target audience. Words are randomly
// selected from the Princeton WordNet database. Use at own peril.
func (f *Factory) Generate() string {
	adjIdx := f.random.Intn(len(f.adjectives))
	nounIdx := f.random.Intn(len(f.nouns))

	return fmt.Sprintf("%s-%s", f.adjectives[adjIdx], f.nouns[nounIdx])
}
