package bananaphone

import (
	"bufio"
	"bytes"
	"crypto"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
)

func NewEncoder(wr io.Writer, encoder, spec, dict string) io.WriteCloser {
	switch encoder {
	case "random":
		return NewRandomEncoder(wr, spec, dict)
	case "markov":
		return NewMarkovEncoder(wr, spec, dict, 1)
	}

	panic("bad encoder")
}

type RandomEncoder struct {
	wr io.Writer

	dict  []byte
	model map[uint32][]string

	tokenizer bufio.SplitFunc
	hash      crypto.Hash
	bits      uint32
}

func NewRandomEncoder(wr io.Writer, spec string, dictfile string) *RandomEncoder {
	enc := &RandomEncoder{wr: wr}

	enc.tokenizer, enc.hash, enc.bits = parseencodingspec(spec)

	dicttext, err := ioutil.ReadFile(dictfile)
	if err != nil {
		return nil
	}

	enc.dict = dicttext

	enc.mkmodel()

	return enc
}

func (re *RandomEncoder) mkmodel() {
	m := make(map[uint32][]string)
	h := re.hash.New()

	scan := bufio.NewScanner(bytes.NewReader(re.dict))
	scan.Split(re.tokenizer)

	for scan.Scan() {
		h.Reset()
		h.Write(scan.Bytes())
		tokenhash := h.Sum(nil)
		tobyte := uint32(tokenhash[len(tokenhash)-1])
		m[tobyte] = append(m[tobyte], scan.Text())
	}

	//for i := []uint32(0); i < MaxUint32
	for i := uint32(0); i < 256; i++ {
		if _, ok := m[i]; !ok {
			panic("8 bit word space not full")
		}
	}

	re.model = m
}

func (re *RandomEncoder) Write(p []byte) (n int, err error) {
	for _, b := range p {
		words, ok := re.model[uint32(b)]
		if !ok {
			panic("no words for byte " + fmt.Sprintf("%X", b))
		}
		if _, err := io.WriteString(re.wr, words[rand.Intn(len(words))]); err != nil {
			return 0, err
		}
	}

	return len(p), nil
}

func (re *RandomEncoder) Close() error {
	if closer, ok := re.wr.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
