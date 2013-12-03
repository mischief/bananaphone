package bananaphone

import (
	"bufio"
	"crypto"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
)

// Prefix is a Markov chain prefix of one or more words.
type Prefix []string

// String returns the Prefix as a string (for use as a map key).
func (p Prefix) String() string {
	return strings.Join(p, " ")
}

// Shift removes the first word from the Prefix and appends the given word.
func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

// Chain contains a map ("chain") of prefixes to a list of suffixes.
// A prefix is a string of prefixLen words joined with spaces.
// A suffix is a single word. A prefix can have multiple suffixes.
type Chain struct {
	Chain     map[string][]string
	prefixLen int
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func NewChain(prefixLen int) *Chain {
	return &Chain{make(map[string][]string), prefixLen}
}

// Build reads text from the provided Reader and
// parses it into prefixes and suffixes that are stored in Chain.
func (c *Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(Prefix, c.prefixLen)
	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}
		key := p.String()
		c.Chain[key] = append(c.Chain[key], s)
		p.Shift(s)
	}
}

// Generate returns a string of at most n words generated from Chain.
func (c *Chain) Generate(n int) string {
	p := make(Prefix, c.prefixLen)
	nprefix := len(c.Chain)
	rnd := rand.Intn(nprefix)
	i := 0
	for pref := range c.Chain {
		if i == rnd {
			p = strings.SplitN(pref, " ", c.prefixLen)
		}
		i++
	}

	var words []string
	for i := 0; i < n; i++ {
		choices := c.Chain[p.String()]
		if len(choices) == 0 {
			break
		}
		next := choices[rand.Intn(len(choices))]
		words = append(words, next)
		p.Shift(next)
	}
	return strings.Join(words, " ")
}

type MarkovEncoder struct {
	wr        io.Writer
	chain     *Chain
	tokenizer bufio.SplitFunc
	hash      crypto.Hash
	bits      uint32
}

func NewMarkovEncoder(wr io.Writer, spec string, dictfile string, order int) *MarkovEncoder {
	me := &MarkovEncoder{wr: wr}

	me.tokenizer, me.hash, me.bits = parseencodingspec(spec)

	dict, err := os.Open(dictfile)
	if err != nil {
		return nil
	}
	defer dict.Close()

	me.chain = NewChain(order)
	me.chain.Build(dict)

	return me
}

// Writes p bytes as reverse hash encoded data through markov chains.
func (me *MarkovEncoder) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("fix me")
	/*
		for _, b := range p {
			words := re.model[uint32(b)]
			if _, err := io.WriteString(re.wr, words[rand.Intn(len(words))]); err != nil {
				return 0, err
			}
		}
	*/
	return len(p), nil
}

func (me *MarkovEncoder) Close() error {
	if closer, ok := me.wr.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
