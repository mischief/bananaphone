package bananaphone

import (
	"bufio"
  "crypto"
  _ "crypto/sha1"
)

func words(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanWords(data, atEOF)
	if err == nil && token != nil {
		token = append(token, ' ')
	}
	return
}

func parseencodingspec(spec string) (bufio.SplitFunc, crypto.Hash, uint32) {
	return words, crypto.Hash(crypto.SHA1), 8
}
