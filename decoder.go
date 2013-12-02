package bananaphone

import (
	"bufio"
	"crypto"
	"io"
)

type Decoder struct {
	out io.Writer

	tokenizer bufio.SplitFunc
	hash      crypto.Hash
	bits      uint32

	inpipew *io.PipeWriter
	inpiper *io.PipeReader
}

func NewDecoder(wr io.Writer, spec string) *Decoder {
	dec := &Decoder{
		out: wr,
	}

	dec.tokenizer, dec.hash, dec.bits = parseencodingspec(spec)

	dec.inpiper, dec.inpipew = io.Pipe()

	go func() {
		h := dec.hash.New()
		scan := bufio.NewScanner(dec.inpiper)
		scan.Split(dec.tokenizer)
		for scan.Scan() {
			h.Reset()
			h.Write(scan.Bytes())
			tokenhash := h.Sum(nil)
			hashbyte := tokenhash[len(tokenhash)-1:]
			dec.out.Write(hashbyte)
		}

		if closer, ok := dec.out.(io.Closer); ok {
			closer.Close()
		}
	}()

	return dec
}

func (dec *Decoder) Write(p []byte) (n int, err error) {
	return dec.inpipew.Write(p)
}

func (dec *Decoder) Close() error {
	dec.inpipew.Close()
	return nil
}
