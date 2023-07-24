package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var settings = Settings{}

type Variables struct {
	Variable string `json:"variable"`
	Setting  string `json:"setting"`
}

type Settings struct {
	RootFolder              string      `json:"rootfolder"`
	CaptureOutput           string      `json:"captureoutput"`
	FileNameTimeStampFormat string      `json:"filenametimestampformat"`
	Variables               []Variables `json:"variables"`
}

func LoadSettings() {
	jsonFile, err := os.Open("settings.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Println(err)
	}
	log.Println("Successfully Opened settings.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	log.Printf("Loading settings to the memory...\n")
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &settings)
	if err != nil {
		log.Printf("Unable to load settings from settings.json: %s\n", err)
		return
	}
	for _, item := range settings.Variables {
		log.Printf("New variable: %s as setting: %s have been registered\n", item.Variable, item.Setting)
	}
	log.Printf("Settings loaded! \n")
}

func LoadProgramData() {
	LoadCommands()
	LoadSettings()
}
