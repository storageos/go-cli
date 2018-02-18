package formatter

import (
	"errors"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestMarshalEntry(t *testing.T) {
	time1, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05+07:00")

	testcases := []struct {
		name      string
		message   []byte
		wantEntry *logrus.Entry
		wantError error
	}{
		{
			name:    "message containing all the known fields",
			message: []byte(`{"time": "2006-01-02T15:04:05+07:00", "level": "warning", "msg": "foo"}`),
			wantEntry: &logrus.Entry{
				Time:    time1,
				Level:   logrus.WarnLevel,
				Message: "foo",
			},
		},
		{
			name:      "empty message",
			message:   []byte(`{}`),
			wantEntry: &logrus.Entry{},
		},
		{
			name:    "message containing unknown fields",
			message: []byte(`{"level": "info", "field1": "val1", "field2": "val2"}`),
			wantEntry: &logrus.Entry{
				Level: logrus.InfoLevel,
				Data: map[string]interface{}{
					"field1": "val1",
					"field2": "val2",
				},
			},
		},
		{
			name:      "invalid byte message ",
			message:   []byte(`{`),
			wantError: errors.New("unexpected end of JSON input"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotEntry, err := marshalEntry(tc.message)
			if err == nil {
				if tc.wantError != nil {
					t.Fatalf("unexpected error while marshalling message:\n\t(GOT): %v\n\t(WNT): %v", err, tc.wantError)
				}

				if !logrusEntryEqual(gotEntry, tc.wantEntry) {
					t.Errorf("unexpected result after marshalling:\n\t(GOT): %v\n\t(WNT): %v", gotEntry, tc.wantEntry)
				}

			} else if err.Error() != tc.wantError.Error() {
				t.Fatalf("unexpected error while marshalling message:\n\t(GOT): %v\n\t(WNT): %v", err, tc.wantError)
			}
		})
	}
}

// logrusEntryEqual compares only a set of fields we are interested in.
func logrusEntryEqual(a, b *logrus.Entry) bool {
	if !a.Time.Equal(b.Time) || a.Level != b.Level || a.Message != b.Message {
		return false
	}

	// Since Data is of type logrus.Fields (map[string]interface{}), we can't
	// compare the interfaces. Checking just the existence of the same keys.
	for k := range a.Data {
		if _, ok := b.Data[k]; !ok {
			return false
		}
	}

	return true
}
