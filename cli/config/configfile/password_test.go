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
