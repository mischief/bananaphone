package main

import (
	"github.com/mischief/bananaphone"
	"io"
	"os"
)

func main() {
	enc := bananaphone.NewEncoder(os.Stdout, "random", "", "/usr/share/dict/words")

	io.Copy(enc, os.Stdin)
}
