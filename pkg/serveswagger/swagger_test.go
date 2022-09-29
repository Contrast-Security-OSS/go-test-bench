package serveswagger

import (
	"bytes"
	"os/exec"
	"regexp"
	"strings"
	"testing"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
)

func TestCheckTimestamps(t *testing.T) {
	rapi, err := common.LocateDir("cmd/go-swagger/restapi", 5)
	if err != nil {
		t.Fatal(err)
	}

	var stderr bytes.Buffer
	cmdGenerate := exec.Command("go", "generate", rapi)
	cmdGenerate.Stderr = &stderr
	if err := cmdGenerate.Run(); err != nil {
		t.Fatalf("Running %v: error %s\noutput:\n%s", cmdGenerate.Args, err, stderr.String())
	}
	if stderr.Len() == 0 {
		t.Fatalf("There is no output from generate command %v.", cmdGenerate.Args)
	}

	sswag, err := common.LocateDir("pkg/serveswagger", 5)
	if err != nil {
		t.Fatal(err)
	}

	generatedLines := regexp.MustCompile(`// Generated at [0-9]{4}-[0-9]{2}`)

	const maxMismatches = 10
	var mismatchCount int
	var mismatchLines = make([]string, 0, maxMismatches)

	// check for differences in pkg/serveswagger and cmd/go-swagger/restapi
	for _, dir := range []string{rapi, sswag} {
		diff := exec.Command("git", "diff", dir)

		stdout, err := diff.Output()

		if err != nil {
			t.Fatal(err)
		}

		lines := strings.Split(string(stdout), "\n")
		var changes []string
		for _, line := range lines {
			if len(line) < 2 {
				continue
			}
			if line[0] != '-' && line[0] != '+' {
				continue
			}
			if line[1] != '-' && line[1] != '+' {
				changes = append(changes, line)
			}
		}

		// changes now contains every changed/added/removed line, and no
		// context lines. Every one of those lines should match the regexp.

		for _, line := range changes {
			if !generatedLines.Match([]byte(line)) {
				mismatchCount++
				if mismatchCount <= maxMismatches {
					mismatchLines = append(mismatchLines, line)
				}
			}
		}
	}
	t.Logf("%d mismatches", mismatchCount)
	if mismatchCount > maxMismatches {
		t.Logf("first %d mismatches follow:", maxMismatches)
	}
	for _, line := range mismatchLines {
		t.Errorf("changed line %q does not match regexp", line)
	}
}
