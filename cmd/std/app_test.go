package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseTemplates(t *testing.T) {
	err := parseTemplates()
	if err != nil {
		t.Error(err)
	}
}

//ensure template is populated
func TestServedTemplate(t *testing.T) {
	setup()
	srvr := httptest.NewServer(nil)
	defer srvr.Close()
	pd.Addr = srvr.URL
	resp, err := http.Get(srvr.URL + "/cmdInjection/")
	if err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	const findOut = "To find out more about"
	if !strings.Contains(string(body), findOut+" Command Injection") {
		lines := strings.Split(string(body), "\n")
		var i int
		var found bool
		for i = range lines {
			if strings.Contains(lines[i], findOut) {
				found = true
				break
			}
		}
		if found {
			t.Errorf("missing command name at %d:\n--  %s\n==> %s\n-- %s\n", i, lines[i-1], lines[i], lines[i+2])
		} else {
			t.Errorf("missing command name in\n%s", body)
		}
		var rb []string
		for _, v := range pd.Rulebar {
			rb = append(rb, v.String())
		}
		t.Logf("\n%s", strings.Join(rb, "\n"))
	}
}
