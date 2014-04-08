package main

import (
    "os/exec"
    "strings"
)

type SslCheck struct {
    Url string
}

func (sslCheck *SslCheck) CheckSync() (bool, error) {
    app := "python"
    arg0 := "./hb.py"
    arg1 := sslCheck.Url

    cmd := exec.Command(app, arg0, arg1)
    output, err := cmd.Output()

    if err != nil {
        logger.Errorf("error running python script: %v", err)
        return false, err
    }

    outLines := strings.Split(string(output), "\n")
    lastLine := outLines[len(outLines)-2]

    return strings.Contains(lastLine, "server is vulnerable"), nil
}
