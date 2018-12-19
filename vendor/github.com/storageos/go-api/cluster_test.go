package storageos

import (
	"net/http"
	"testing"
)

func TestMaintenanceGet(t *testing.T) {
	msg := `
    {
        "enabled": false,
        "updatedBy": "storageos",
        "updatedAt": "2018-11-14T12:03:07.165496683Z"
    }
	`
	client := newTestClient(&FakeRoundTripper{message: msg, status: http.StatusOK})

	res, err := client.Maintenance()
	if err != nil {
		t.Fatal(err)
	}

	if res.UpdatedBy != "storageos" {
		t.Fatalf("expect %s\ngot%v", msg, res)
	}
}

func TestMaintenanceEnable(t *testing.T) {
	msg := `
    {
        "enabled": true,
        "updatedBy": "storageos",
        "updatedAt": "2018-11-14T12:03:07.165496683Z"
    }
	`
	client := newTestClient(&FakeRoundTripper{message: msg, status: http.StatusOK})

	if err := client.EnableMaintenance(); err != nil {
		t.Fatal(err)
	}
}

func TestMaintenanceDisable(t *testing.T) {
	msg := `
    {
        "enabled": false,
        "updatedBy": "storageos",
        "updatedAt": "2018-11-14T12:03:07.165496683Z"
    }
	`
	client := newTestClient(&FakeRoundTripper{message: msg, status: http.StatusOK})

	if err := client.DisableMaintenance(); err != nil {
		t.Fatal(err)
	}
}
