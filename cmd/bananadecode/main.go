package main

import (
	"github.com/mischief/bananaphone"
	"io"
	"os"
)

func main() {
	dec := bananaphone.NewDecoder(os.Stdout, "")

	io.Copy(dec, os.Stdin)
}
