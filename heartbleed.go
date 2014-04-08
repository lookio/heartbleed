package main

import (
	"code.google.com/p/gcfg"
	"flag"
	"fmt"
	"github.com/howbazaar/loggo"
	"os"
)

type Config struct {
	Server struct {
		ListenAddress string
		ListenPort    string
	}
}

var (
	printVersion bool
	verbose      bool
	configFile   string
	serverConfig Config

	logger = loggo.GetLogger("")
)

const releaseVersion = "0.0.1"

func init() {
	flag.BoolVar(&printVersion, "version", false, "print the version and exit")
	flag.BoolVar(&verbose, "v", false, "verbose logging")
	flag.StringVar(&configFile, "config", "/etc/heartbleed.conf", "location of the config file")
}

func main() {
	flag.Parse()

	logFilename := "/var/log/heartbleed"
	var logFile *os.File
	if _, err := os.Stat(logFilename); err == nil {
		logFile, err = os.OpenFile(logFilename, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			logger.Errorf("unable to open %s for appending: %v", logFilename, err)
		}
	} else {
		logFile, err = os.Create(logFilename)
		if err != nil {
			logger.Errorf("unable to open %s for creating: %v", logFilename, err)
		}
	}

	logWriter := loggo.NewSimpleWriter(logFile, &loggo.DefaultFormatter{})
	if verbose {
		logger.SetLogLevel(loggo.DEBUG)
		loggo.RegisterWriter("file", logWriter, loggo.DEBUG)
	} else {
		loggo.RegisterWriter("file", logWriter, loggo.WARNING)
	}

	if printVersion {
		fmt.Println(releaseVersion)
		os.Exit(0)
	}

	logger.Debugf("starting heartbleed")

	// set some defaults
	serverConfig.Server.ListenAddress = "0.0.0.0"
	serverConfig.Server.ListenPort = "8000"

	if _, err := os.Stat(configFile); err == nil {
		err := gcfg.ReadFileInto(&serverConfig, configFile)
		if err != nil {
			logger.Errorf("unable to read config file: %v", err)
			os.Exit(1)
		}
	} else {
		logger.Warningf("config file not found, using defaults")
	}

	logger.Debugf("config: %v", serverConfig)

	webServer := WebServer{}
	webServer.Listen()

}
