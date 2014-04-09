package main

import (
	"os/exec"
	"strings"
)

type SslCheck struct {
	Url string
}

func (sslCheck *SslCheck) CheckSync() (bool, error) {
        logger.Debugf("executing command: python ./hb.py %s", sslCheck.Url)
        cmd := exec.Command("python", "./hb.py", sslCheck.Url)
        out, err := cmd.CombinedOutput()

	if err != nil {
            logger.Errorf("error running python script: %v, %s", err, string(out[:]))
		return false, err
	}

        outLines := strings.Split(string(out[:]), "\n")
	lastLine := outLines[len(outLines)-2]

	return strings.Contains(lastLine, "server is vulnerable"), nil
}
