// gitops.go
package gitops

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/thywilljoshua/echo/generate"
)

func StartCommit(message *string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current working directory: %v", err)
	}

	diff, err := exec.Command("git", "diff", "--cached", dir).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running Git command: %v", err)
	}
	if string(diff) == "" {
		fmt.Println("No changes staged.")
		os.Exit(0)
	}

	fmt.Println("Generating message...")

	commitMsg, err := generate.GenerateWithOllama(string(diff))
	if err != nil {
		return fmt.Errorf("error generating diff message")
	}

	*message = commitMsg
	return nil
}
