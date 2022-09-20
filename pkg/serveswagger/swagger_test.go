package serveswagger

import (
	"os/exec"
	"regexp"
	"testing"
)

func TestCheckTimestamps(t *testing.T) {
	cmd := exec.Command("git", "diff", "../../pkg/serverswagger")

	stdout, err := cmd.Output()

	if err != nil {
		t.Fatal(err.Error())
	}

	updates, _ := regexp.Compile(`@@.*@@`)
	generatedLines, _ := regexp.Compile(`// Generated at [0-9]{4}-[0-9]{2}`)

	matchesUpdates := updates.FindAllStringSubmatch(string(stdout), -1)
	matchesGeneratedLines := generatedLines.FindAllStringSubmatch(string(stdout), -1)

	if len(matchesUpdates) != len(matchesGeneratedLines) {
		t.Fatal("Not only generated files are changed!")
	}
}
