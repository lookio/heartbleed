package main

import (
	"os/exec"
	"strings"
        "bytes"
)

type SslCheck struct {
	Url string
}

func (sslCheck *SslCheck) CheckSync() (bool, error) {
        cmd := exec.Command("python", "./hb.py", sslCheck.Url)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		logger.Errorf("error running python script: %v, %s", err, out.String())
		return false, err
	}

	outLines := strings.Split(out.String(), "\n")
	lastLine := outLines[len(outLines)-2]

	return strings.Contains(lastLine, "server is vulnerable"), nil
}
