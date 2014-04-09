package main

import (
	"os/exec"
	"strings"
)

type SslCheck struct {
	Url string
}

func (sslCheck *SslCheck) CheckSync() (bool, error) {
	isCached, err := cache.Exists("url:" + sslCheck.Url)
	if err != nil {
		return false, err
	}

	if isCached {
		cachedValue, err := cache.Get("url:" + sslCheck.Url)
		if err != nil {
			return false, err
		}

		return cachedValue == "1", nil
	}

	logger.Debugf("executing command: python /opt/local/heartbleed.py %s", sslCheck.Url)
	cmd := exec.Command("python", "./hb.py", sslCheck.Url)
	out, err := cmd.CombinedOutput()

	if err != nil {
		logger.Errorf("error running python script: %v, %s", err, string(out[:]))
		return false, err
	}

	outLines := strings.Split(string(out[:]), "\n")
	lastLine := outLines[len(outLines)-2]

	vulnerable := strings.Contains(lastLine, "server is vulnerable")

	if vulnerable {
		err = cache.Set("url:"+sslCheck.Url, "1")
	} else {
		err = cache.Set("url:"+sslCheck.Url, "0")
	}

	if err != nil {
		return false, err
	}

        cache.Expire("url:"+sslCheck.Url, 60*60)

	return vulnerable, nil
}
