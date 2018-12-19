package storageos

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/storageos/go-api/types"
)

func TestLicence(t *testing.T) {
	body := `{
  "arrayUUID": "67f36f7d-a281-4070-8227-7e77d4b85dfa",
  "customerID": "abc",
  "customerName": "efg",
  "storage": 100,
  "validUntil": "0001-01-01T00:00:00Z",
  "licenceType": "basic",
  "features": {
    "HA": true
  },
  "unregistered": true
}`

	expected := &types.Licence{}
	if err := json.Unmarshal([]byte(body), expected); err != nil {
		t.Fatal(err)
	}

	client := newTestClient(&FakeRoundTripper{message: body, status: http.StatusOK})
	licence, err := client.Licence()
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(licence, expected) {
		t.Errorf("wrong return value. Want %#v. Got %#v.", expected, licence)
	}
}

func TestLicenceApply(t *testing.T) {
	licenceKey := `ABCDE`

	fakeRT := &FakeRoundTripper{status: http.StatusOK}
	client := newTestClient(fakeRT)
	err := client.LicenceApply(licenceKey)
	if err != nil {
		t.Fatal(err)
	}
	req := fakeRT.requests[0]
	expectedMethod := "POST"
	if req.Method != expectedMethod {
		t.Errorf("LicenceApply(): Wrong HTTP method. Want %s. Got %s.", expectedMethod, req.Method)
	}
	u, _ := url.Parse(client.getAPIPath(licenceAPIPrefix, url.Values{}, false))
	if req.URL.Path != u.Path {
		t.Errorf("LicenceApply(): Wrong request path. Want %q. Got %q.", u.Path, req.URL.Path)
	}
}

func TestLicenceDelete(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: "", status: http.StatusNoContent}
	client := newTestClient(fakeRT)
	err := client.LicenceDelete()
	if err != nil {
		t.Fatal(err)
	}
	req := fakeRT.requests[0]
	expectedMethod := "DELETE"
	if req.Method != expectedMethod {
		t.Errorf("Wrong HTTP method. Want %s. Got %s.", expectedMethod, req.Method)
	}
	path := licenceAPIPrefix
	u, _ := url.Parse(client.getAPIPath(path, url.Values{}, false))
	if req.URL.Path != u.Path {
		t.Errorf("Wrong request path. Want %q. Got %q.", u.Path, req.URL.Path)
	}
}
