package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	// Port to listen on
	Port uint16
	// Directory for saving uploaded files
	UploadDirectory string
	// File containing user info with passwords
	PostersInfoFile string
	// True if video upload is enabled
	Videos bool
	// Alternative video upload url
	VideoUploadURL string
	// Conference Homepage URL
	ConferencePageURL string
	// Support email address displayed on the page
	SupportEmail string
	// Number of file versions to keep
	KeepVersions int
	// Date as YYYY-MM-DD string when the poster submission is closed
	SubmissionClosedDate string
	// Text when the poster submission is closed
	SubmissionClosedText string
	// Text when the video upload is closed
	SubmissionClosedVideoText string
	// File whitelisted email addresses can be uploaded to
	WhitelistFile string
	// Password for whitelist email address upload
	WhitelistPW string
}

func defaultConfig() *Config {
	return &Config{
		Port:                      3000,
		UploadDirectory:           "uploads",
		PostersInfoFile:           "posters.json",
		Videos:                    false,
		VideoUploadURL:            "",
		ConferencePageURL:         "https://www.bernstein-network.de/en/bernstein-conference/",
		SupportEmail:              "bernstein.conference@fz-juelich.de",
		KeepVersions:              5,
		SubmissionClosedDate:      "2100-12-31",
		SubmissionClosedText:      "Sunday, Sep 19, 2021, 8 pm CEST",
		SubmissionClosedVideoText: "Friday, Sep 17, 1 pm CEST",
		WhitelistFile:             "whitelist.txt",
		WhitelistPW:               fmt.Sprint(time.Now().UnixNano()),
	}
}

func readConfig(configFileName string) *Config {
	configFile, err := os.Open(configFileName)
	if err != nil {
		log.Printf("[os.Open] Error reading config file %q: %s", configFileName, err.Error())
		os.Exit(1)
	}
	defer configFile.Close()

	data, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Printf("[ioutil.ReadAll] Error reading config file %q: %s", configFileName, err.Error())
		os.Exit(1)
	}

	config := defaultConfig() // set defaults first
	if err := yaml.Unmarshal(data, config); err != nil {
		log.Printf("[yaml.Unmarshall] Error reading config file (%q): %s", configFileName, err.Error())
		os.Exit(1)
	}
	// create upload directory (if it doesn't exist)
	err = os.MkdirAll(config.UploadDirectory, 0777)
	if err != nil {
		log.Printf("[os.MkdirAll] Error creating upload directory (%s): %s", config.UploadDirectory, err.Error())
		os.Exit(1)
	}
	return config
}

// writeConfig writes the default configuration values to the specified file.
func writeConfig(cfgFileName string) {
	// using fmt.Print for error messages here since it's run interactively and
	// the log-style formatting with timestamps makes it noisy.
	cfgYml, err := yaml.Marshal(defaultConfig())
	if err != nil {
		fmt.Printf("Error marshalling default config: %s\n", err.Error())
		os.Exit(1)
	}

	cfgFile, err := os.Create(cfgFileName)
	if err != nil {
		fmt.Printf("Error creating config file: %s\n", err.Error())
		os.Exit(1)
	}
	defer cfgFile.Close()

	if _, err := cfgFile.Write(cfgYml); err != nil {
		fmt.Printf("Error writing default config: %s\n", err.Error())
		os.Exit(1)
	}
}

// prompt reads input from the command line and returns
// the read result as string. Returns an empty string in
// case of an error.
func prompt(msg string) string {
	var response string
	fmt.Printf("%s: ", msg)
	_, err := fmt.Scanln(&response)
	if err != nil {
		fmt.Printf("Error reading input: %s", err.Error())
		return ""
	}
	return response
}
