package serveswagger

import (
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func TestCheckTimestamps(t *testing.T) {
	cmdGenerate := exec.Command("go", "generate", "./cmd/go-swagger/restapi")
	if err := cmdGenerate.Run(); err != nil{
		t.Fatal(err.Error())
	}

	cmd := exec.Command("git", "diff", "../../pkg/serveswagger")

	stdout, err := cmd.Output()

	if err != nil {
		t.Fatal(err.Error())
	}

	lines := strings.Split(string(stdout), "\n")
	var changes []string
	for _ ,line := range lines {
		if len(line) == 0 {
			continue
		}
		if line[0] == '-' || line[0] == '+' {
			changes = append(changes, line)
		}
	}

	// changes now contains every changed/added/removed line, and no
	// context lines. Every one of those lines should match the regexp.

	generatedLines := regexp.MustCompile(`// Generated at [0-9]{4}-[0-9]{2}`)
	for _, line := range changes {
		if !generatedLines.Match([]byte(line)) {
			t.Errorf("changed line %s does not match regexp",line)
		}
	}
}
