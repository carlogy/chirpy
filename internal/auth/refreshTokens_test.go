package auth

import "testing"

func TestMakeRefreshToken(t *testing.T) {
	got, err := MakeRefreshToken()
	if err != nil {
		t.Errorf("Expecting nil err, Got:\t%v", err)
	}

	if len(got) != 64 {
		t.Errorf("Hexadecimal string not of the correct length\nGot:\t%v\tLen:\t%v", got, len(got))
	}

}
