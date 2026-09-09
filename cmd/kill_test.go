package cmd

import "testing"

func TestParsePort(t *testing.T) {
	for _, ok := range []string{"80", "8080", "65535", "1"} {
		if _, err := parsePort(ok); err != nil {
			t.Fatalf("esperava porta válida %q: %v", ok, err)
		}
	}
	for _, bad := range []string{"", "abc", "99999999", "0", "65536", "-1", "80a", " 8080"} {
		if _, err := parsePort(bad); err == nil {
			t.Fatalf("esperava erro para porta %q", bad)
		}
	}
}
