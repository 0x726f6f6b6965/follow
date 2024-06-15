package pagination

import "testing"

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
