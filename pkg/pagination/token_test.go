package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetToken(t *testing.T) {
	token := &PageToken{
		LastId: 10,
		Size:   100,
	}
	tokenString := token.String()

	newToken := &PageToken{}
	err := newToken.DecodePageTokenStruct(tokenString)
	if err != nil {
		t.Fatal(err)
	}
	if newToken.LastId != token.LastId || newToken.Size != token.Size {
		t.Fatal("decode error")
	}
}

func TestInvaildToken(t *testing.T) {
	token := &PageToken{
		LastId: 10,
		Size:   100,
	}

	err := token.DecodePageTokenStruct("abc")
	assert.NotNil(t, err)
	if token.LastId != 10 || token.Size != 100 {
		t.Fatal("decode error")
	}
}
