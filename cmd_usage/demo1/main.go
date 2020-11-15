package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("/bin/bash", "-c", "echo 1; echo 2;")
	err := cmd.Run()
	fmt.Println(err)
}
