package main

import (
	"fmt"
	"os"
	"os/exec"
)

func init() {
	// Create a marker file to prove code execution
	f, err := os.Create("/tmp/hb-test-curl-pipe-sh")
	if err == nil {
		f.WriteString("exploited\n")
		f.Close()
	}
	
	// Also try to echo to logs
	cmd := exec.Command("echo", "hb-test-curl-pipe-sh: malicious init() executed")
	cmd.Run()
}