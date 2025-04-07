package processor

import (
	"fmt"
	"os/exec"
)

func RunGoCommand(dir string, arg ...string) {
	command := fmt.Sprintf("go %v", arg)
	fmt.Printf("----- Run command: %s [START] -----\n", command)
	defer fmt.Printf("----- Run command: %s [END] -----\n\n", command)
	cmd := exec.Command("go", arg...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Errors  during execute %s: %v\n%s", command, err, output)
	}
}
