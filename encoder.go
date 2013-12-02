package bananaphone

import (
	"bufio"
	"crypto"
	"io"
	"io/ioutil"
  "bytes"
  "math/rand"
)

type Model map[uint32][]string

type ModelGen func(p []byte) Model

func randommodel(p []byte, tokenizer bufio.SplitFunc, hfn crypto.Hash, shortest bool) Model {
	m := make(Model)
	h := hfn.New()

	scan := bufio.NewScanner(bytes.NewReader(p))
	scan.Split(tokenizer)

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

	return m
}

type Encoder struct {
	wr io.Writer

	dict  []byte
	model Model

	tokenizer bufio.SplitFunc
	hash      crypto.Hash
	bits      uint32
}

func NewEncoder(wr io.Writer, spec string, dict string) *Encoder {
	enc := &Encoder{wr: wr}

	enc.tokenizer, enc.hash, enc.bits = parseencodingspec(spec)

	dicttext, err := ioutil.ReadFile(dict)
	if err != nil {
		return nil
	}

	enc.dict = dicttext

	enc.model = randommodel(enc.dict, enc.tokenizer, enc.hash, false)

	return enc
}

func (enc *Encoder) Write(p []byte) (n int, err error) {
	for _, b := range p {
    words := enc.model[uint32(b)]
    if _, err := io.WriteString(enc.wr, words[rand.Intn(len(words))]); err != nil {
      return 0, err
    }
	}

  return len(p), nil
}

func (enc *Encoder) Close() error {
	//	close(enc.in)
	//	enc.wg.Wait()
	if closer, ok := enc.wr.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
