package main

import (
	"dt-extension/pkg/downloader"
	"dt-extension/pkg/logger"
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"time"
)

var (
	log           = logger.NewDTLogger()
	extDownloader = downloader.NewExtensionDownloader(log)
)

var callbacks = map[string]func() error{
	"extension":    checkForever,
	"health-check": healthCheck,
}

func healthCheck() error {
	downloaded := extDownloader.CheckIfDownloaded()
	if !downloaded {
		log.Info(fmt.Sprintf("The extension was not downloaded yet"))
	}
	log.Info(fmt.Sprintf("The extension was already downloaded successfully"))
	return nil

}

func checkForever() error {

	for {
		downloaded := extDownloader.CheckIfDownloaded()
		if !downloaded {
			err := extDownloader.Download()
			if err != nil {
				return err
			}
		}
		log.Info(fmt.Sprintf("The extension was already downloaded successfully"))
		time.Sleep(30 * time.Second)

	}
}

func main() {
	flags := pflag.NewFlagSet("health-check", pflag.ExitOnError)

	pflag.CommandLine.AddFlagSet(flags)
	pflag.Parse()

	command := "extension"
	if args := pflag.Args(); len(args) > 0 {
		command = args[0]
	}

	function := callbacks[command]
	if function == nil {
		os.Exit(1)
	}

	err := function()
	if err != nil {
		os.Exit(1)
	}

}
