package pagination

import (
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type PageToken struct {
	LastId int `json:"last_id,omitempty"`
	Size   int `json:"size,omitempty"`
}

// String returns a string representation of the page token.
func (p *PageToken) String() string {
	var b strings.Builder
	base64Encoder := base64.NewEncoder(base64.URLEncoding, &b)
	gobEncoder := gob.NewEncoder(base64Encoder)
	jsonData, _ := json.Marshal(p)
	_ = gobEncoder.Encode(jsonData)
	_ = base64Encoder.Close()
	return b.String()
}

// DecodePageTokenStruct decodes an encoded page token into an arbitrary struct.
func (p *PageToken) DecodePageTokenStruct(s string) error {
	dec := gob.NewDecoder(base64.NewDecoder(base64.URLEncoding, strings.NewReader(s)))
	var b []byte
	if err := dec.Decode(&b); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("decode page token struct: %w", err)
	}
	_ = json.Unmarshal(b, p)
	return nil
}
