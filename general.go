package pgrst

import (
	"bytes"
	"encoding/json"
	"io"
)

// NewReader - create new instance of reader from a given structure
func NewReader(i any) (r io.Reader, err error) {

	iBytes, err := json.Marshal(i)
	if nil != err {
		return
	}

	return bytes.NewReader(iBytes), nil
}
