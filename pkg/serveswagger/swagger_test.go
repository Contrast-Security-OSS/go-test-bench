package serveswagger

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func TestCheckTimestamps(t *testing.T) {
	var stderr bytes.Buffer
	cmdGenerate := exec.Command("go", "generate", "../../cmd/go-swagger/restapi")
	cmdGenerate.Stderr = &stderr
	if err := cmdGenerate.Run(); err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		t.Fatal(err.Error())
	}
	if stderr.Len() == 0 {
		t.Fatal("There is no output from generate command.")
	}

	cmd := exec.Command("git", "diff", "../../pkg/serveswagger")

	stdout, err := cmd.Output()

	if err != nil {
		t.Fatal(err.Error())
	}

	lines := strings.Split(string(stdout), "\n")
	var changes []string
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		if (line[0] == '-' && line[1] != '-') || (line[0] == '+' && line[1] != '+') {
			changes = append(changes, line)
		}
	}

	// changes now contains every changed/added/removed line, and no
	// context lines. Every one of those lines should match the regexp.

	generatedLines := regexp.MustCompile(`// Generated at [0-9]{4}-[0-9]{2}`)
	for _, line := range changes {
		if !generatedLines.Match([]byte(line)) {
			t.Errorf("changed line %s does not match regexp", line)
		}
	}
}
