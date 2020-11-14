package utils

import (
	"bytes"
	"compress/flate"
	"io/ioutil"
)

func FlateEncode(input string) (result []byte, err error) {
	var buf bytes.Buffer
	w, err := flate.NewWriter(&buf, -1)
	w.Write([]byte(input))
	w.Close()
	result = buf.Bytes()
	return
}

func FlateDecode(input []byte) (result []byte, err error) {
	result, err = ioutil.ReadAll(flate.NewReader(bytes.NewReader(input)))
	return
}
