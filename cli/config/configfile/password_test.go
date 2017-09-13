package configfile

import (
	"encoding/json"
	"testing"
)

func TestPasswordEncoding(t *testing.T) {
	enc := encodedPassword("foo")

	b, err := json.Marshal(&enc)
	if err != nil {
		t.Fatal(err)
	}

	if string(b) != `"Zm9v"` {
		t.Log(string(b) + ` != "Zm9v"`)
		t.FailNow()
	}

	var recovered encodedPassword

	err = json.Unmarshal(b, &recovered)
	if err != nil {
		t.Fatal(err)
	}

	if recovered != "foo" {
		t.Log(recovered + ` != foo`)
		t.FailNow()
	}
}

func TestCredentialMarshal(t *testing.T) {
	pass := encodedPassword("bar")

	buf, err := json.Marshal(&credentials{
		Username:    "foo",
		Password:    &pass,
		UseKeychain: false,
	})

	if err != nil {
		t.Fatal(err)
	}

	// Expect password to be marshaled as the keychain is not used
	expect := `{"username":"foo","password":"YmFy"}`
	if string(buf) != expect {
		t.Log(string(buf) + " != " + expect)
		t.Fail()
	}
}

func TestCredentialMarshalKeychain(t *testing.T) {
	pass := encodedPassword("bar")

	buf, err := json.Marshal(&credentials{
		Username:    "foo",
		Password:    &pass,
		UseKeychain: true,
	})

	if err != nil {
		t.Fatal(err)
	}

	// Expect password not to be marshaled as the keychain is used
	expect := `{"username":"foo","useKeychain":true}`
	if string(buf) != expect {
		t.Log(string(buf) + " != " + expect)
		t.Fail()
	}
}
