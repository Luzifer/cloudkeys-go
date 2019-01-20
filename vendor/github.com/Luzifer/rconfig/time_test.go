package rconfig

import (
	"fmt"
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	type ts struct {
		Test    time.Time `flag:"time"`
		TestS   time.Time `flag:"other-time,o"`
		TestDef time.Time `default:"2006-01-02T15:04:05.999999999Z"`
		TestDE  time.Time `default:"18.09.2018 20:25:31"`
	}

	var (
		err  error
		args []string
		cfg  ts
	)

	for _, tf := range timeParserFormats {
		expect := time.Now().Format(tf)

		cfg = ts{}
		args = []string{
			fmt.Sprintf("--time=%s", expect),
			"-o", expect,
		}

		if err = parse(&cfg, args); err != nil {
			t.Fatalf("Time format %q did not parse: %s", tf, err)
		}

		for name, ti := range map[string]time.Time{
			"Long flag":    cfg.Test,
			"Short flag":   cfg.TestS,
			"Default flag": cfg.TestDef,
			"DE flag":      cfg.TestDE,
		} {
			if ti.IsZero() {
				t.Errorf("%s did parse to zero with format %q", name, tf)
			}
		}

		if e := cfg.Test.Format(tf); e != expect {
			t.Errorf("Parsed time %q did not match expectation %q", e, expect)
		}
	}
}
