package servestd

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
	Setup()
	srvr := httptest.NewServer(nil)
	defer srvr.Close()
	Pd.Addr = srvr.URL
	for _, route := range Pd.Rulebar {
		if len(route.Sinks) == 0 || len(route.Sinks[0].Name) == 0 {
			continue
		}
		t.Run(route.Name, func(t *testing.T) {
			url := srvr.URL + route.Base + "/"
			t.Logf("url %s", url)

			resp, err := http.Get(url)
			if err != nil {
				t.Fatal(err)
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			const findOut = "To find out more about "
			if !strings.Contains(string(body), findOut+route.Name) {
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
				for _, v := range Pd.Rulebar {
					rb = append(rb, v.String())
				}
				t.Logf("\n%s", strings.Join(rb, "\n"))
			}
			// for _, s := range route.Sinks {
			// 	t.Run(s.Name, func(t *testing.T) {
			//TODO check each sink
			// 	})
			// }
		})
	}
}
