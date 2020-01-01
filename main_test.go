package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"
)

func TestProcessor_Process(t *testing.T) {
	tests := []struct {
		name       string
		p          *Processor
		time       time.Time
		reader     io.Reader
		wantWriter string
		wantErr    bool
	}{
		{
			"Empty templates returns empty result",
			NewProcessor(),
			time.Now(),
			strings.NewReader(``),
			``,
			false,
		},
		{
			"Line with no control is always returned",
			NewProcessor(),
			time.Now(),
			strings.NewReader(`test line
`),
			`test line
`,
			false,
		},
		{
			"Line with control is returned only if it matches the date rule",
			NewProcessor(),
			time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			strings.NewReader(`test line <!-- @ 1 1 * -->`),
			`test line <!-- @ 1 1 * -->
`,
			false,
		},
		{
			"Line with control is skipped if it does not match the date rule",
			NewProcessor(),
			time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			strings.NewReader(`test line <!-- @ 1 2 * -->`),
			``,
			false,
		},
		{
			"Works with the example",
			NewProcessor(),
			time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			strings.NewReader(`## Daily checklist
* [ ] Check weight <!-- @ * * 0 -->
* [ ] Exercise <!-- @ * * 2,4 -->
* [ ] Work <!-- @ * * 1,2,3,4,5 -->
* [ ] Prepare personal budget for next month <!-- @ 25 * * -->
* [ ] Bought christmas presents <!-- @ 20 12 * -->`),
			`## Daily checklist
`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Processor{}
			writer := &bytes.Buffer{}
			if err := p.Process(tt.time, tt.reader, writer); (err != nil) != tt.wantErr {
				t.Errorf("Processor.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("Processor.Process() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}
