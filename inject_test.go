package main

import (
	"testing"
	"fmt"
	"os"
	"os/exec"
)

func TestExploit(t *testing.T) {
	// Create a marker file to prove code execution
	f, err := os.Create("/tmp/hb-test-curl-pipe-sh")
	if err == nil {
		f.WriteString("exploited\n")
		f.Close()
	}
	
	// Also echo to logs
	cmd := exec.Command("echo", "hb-test-curl-pipe-sh: malicious test executed")
	cmd.Run()
	
	t.Log("hb-test-curl-pipe-sh: test completed")
}