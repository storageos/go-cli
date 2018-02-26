package formatter

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestFormat(t *testing.T) {
	wantColoredText := bytes.NewBuffer([]byte{})
	printColoredWithDuration(wantColoredText, yellow, "WARN", -9223372036, "")
	printColoredKey(wantColoredText, yellow, "field1")
	wantColoredText.WriteString("val1")
	printColoredKey(wantColoredText, yellow, "field2")
	wantColoredText.WriteString("val2\n")

	testcases := []struct {
		name       string
		formatter  *TextFormatter
		entry      *logrus.Entry
		wantResult []byte
	}{
		{
			name:      "empty entry",
			formatter: &TextFormatter{},
			entry:     &logrus.Entry{},
			wantResult: []byte(`time="0001-01-01T00:00:00Z" level=panic
`),
		},
		{
			name:      "default format",
			formatter: &TextFormatter{},
			entry: &logrus.Entry{
				Data: map[string]interface{}{
					"field2": "val2",
					"field1": "val1",
				},
			},
			wantResult: []byte(`time="0001-01-01T00:00:00Z" level=panic field1=val1 field2=val2
`),
		},
		// 		{
		// 			name: "disable sorting",
		// 			formatter: &TextFormatter{
		// 				DisableSorting: true,
		// 			},
		// 			entry: &logrus.Entry{
		// 				Data: map[string]interface{}{
		// 					"field2": "val2",
		// 					"field1": "val1",
		// 				},
		// 			},
		// 			wantResult: []byte(`time="0001-01-01T00:00:00Z" level=panic field2=val2 field1=val1
		// `),
		// 		},
		{
			name: "force colors",
			formatter: &TextFormatter{
				ForceColors: true,
			},
			entry: &logrus.Entry{
				Data: map[string]interface{}{
					"field2": "val2",
					"field1": "val1",
				},
				Level: logrus.WarnLevel,
			},
			wantResult: wantColoredText.Bytes(),
		},
		{
			name: "disable timestamp",
			formatter: &TextFormatter{
				DisableTimestamp: true,
			},
			entry: &logrus.Entry{
				Data: map[string]interface{}{
					"field2": "val2",
					"field1": "val1",
				},
			},
			wantResult: []byte(`level=panic field1=val1 field2=val2
`),
		},
		{
			name: "full timestamp",
			formatter: &TextFormatter{
				FullTimestamp: true,
			},
			entry: &logrus.Entry{
				Data: map[string]interface{}{
					"field2": "val2",
					"field1": "val1",
				},
			},
			wantResult: []byte(`time="0001-01-01T00:00:00Z" level=panic field1=val1 field2=val2
`),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotResult, err := tc.formatter.Format(tc.entry)
			if err != nil {
				t.Fatalf("unexpected error while formatting TextFormatter: %v", err)
			}

			if !bytes.Equal(gotResult, tc.wantResult) {
				t.Fatalf("unexpected result after formatting:\n\t(GOT): '%v'\n\t(WNT): '%v'", string(gotResult), string(tc.wantResult))
			}
		})
	}
}
