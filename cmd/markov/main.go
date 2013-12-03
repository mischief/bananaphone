// Helper command for debugging the markov chain algorithm, and having fun.
package main

import (
	"flag"
	"fmt"
	"github.com/mischief/bananaphone"
	"math/rand"
	"os"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	dump := flag.Bool("dump", false, "dump markov model and exit")
	numWords := flag.Int("words", 100, "maximum number of words to print")
	prefixLen := flag.Int("prefix", 2, "prefix length in words")

	flag.Parse()

	c := bananaphone.NewChain(*prefixLen)
	c.Build(os.Stdin)

	if *dump {
		for p, s := range c.Chain {
			fmt.Printf("%-20s : %s\n", p, strings.Join(s, " / "))
		}
	} else {
		fmt.Println(c.Generate(*numWords))
	}
}
