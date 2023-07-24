package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func ReturnTimePrefix() string {
	dt := time.Now()
	return dt.Format(settings.FileNameTimeStampFormat)
}

type SiteInfo struct {
	Name   string
	Footer string
}

var static_variables = SiteInfo{Name: "IrDA Box", Footer: `irDA Box 2023`}

type Battery struct {
	Name    string `json:"name"`
	Model   string `json:"model"`
	Percent int    `json:"percent"`
	Status  string `json:"status"`
}

func ReadFileLineStripped(file string) (string, error) {
	justFile, err := os.Open(file)
	// if we os.Open returns an error then handle it
	if err != nil {
		return "", err
	}
	// defer the closing of our justFile so that we can parse it later on
	defer justFile.Close()
	// read our opened justFile as a byte array.
	byteValue, err := ioutil.ReadAll(justFile)
	if err != nil {
		return "", err
	}
	read_line := strings.TrimSuffix(string(byteValue), "\n")
	return read_line, err
}

func batteryInfo() []Battery {
	//log.Printf("battery refresh\n")
	var rootSysfs = "/sys/class/power_supply"
	var batPrefix = "BAT"
	var Batteries = []Battery{}
	items, _ := ioutil.ReadDir(rootSysfs)
	for _, item := range items {
		if strings.Contains(item.Name(), batPrefix) {
			subitems, _ := ioutil.ReadDir(rootSysfs + "/" + item.Name())
			var model string
			var percent int
			var status string
			for _, subitem := range subitems {
				if !subitem.IsDir() {
					// handle file there
					batPath := fmt.Sprintf("%s/%s/%s", rootSysfs, item.Name(), subitem.Name())
					out, err := ReadFileLineStripped(batPath)
					if err == nil {
						switch subitem.Name() {
						case "model_name":
							model = out
						case "capacity":
							var perc int
							perc, _ = strconv.Atoi(out)
							percent = perc
						case "status":
							status = out
						}
					}
				}
			}
			bat := Battery{Name: item.Name(), Model: model, Percent: percent, Status: status}
			Batteries = append(Batteries, bat)
		}
	}
	return Batteries
}

func saveSettings() error {
	content, err := json.Marshal(settings)
	if err != nil {
		log.Printf("Unable to marshal settings data err: %s\n", err)
		return err
	}
	err = ioutil.WriteFile("settings.json", content, 0644)
	if err != nil {
		log.Printf("Unable to write to settings file: %s\n", err)
		return err
	}
	log.Printf("Settings have been successfully saved!\n")
	return nil
}

func StringJoinFix(data []string, sep string) string {
	string := ""
	for index, element := range data {
		if strings.Contains(element, " ") {
			if index >= len(data)-1 {
				string += fmt.Sprintf("\"%s\"", element)
			} else {
				string += fmt.Sprintf("\"%s\" ", element)
			}

		} else {
			if index >= len(data)-1 {
				string += element
			} else {
				string += element + " "
			}
		}
	}

	return string

	//return strings.Join(data, sep)
}

func Clone(a, b interface{}) {

	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	dec := gob.NewDecoder(buff)
	enc.Encode(a)
	dec.Decode(b)
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}
