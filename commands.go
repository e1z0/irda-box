package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/lithammer/shortuuid"
)

var CurrentProgram = ProgramRunning{}
var Commands = []Command{}
var ProgramExecution *exec.Cmd
var StaticVariables = []StaticVariable{}

type StaticVariable struct {
	Hashname string
	Info     string
}

type ProgramRunning struct {
	Uid       string
	Exe       string
	running   bool
	info      string
	StartTime time.Time
}

type Command struct {
	Uid      string     `json:"uid"`
	Name     string     `json:"name"`
	Class    string     `json:"class"`
	Icon     string     `json:"icon"`
	Commands [][]string `json:"commands"`
	Info     string     `json:"info"`
}

func AddDummyStaticVariables() {
	StaticVariables = append(StaticVariables, StaticVariable{Hashname: "{datetimenow}", Info: "User defined datetime format prefix"})
	StaticVariables = append(StaticVariables, StaticVariable{Hashname: "{uuid}", Info: "Random generated file system safe uuid that shares across commands on application"})
}

func LoadCommands() {
	jsonFile, err := os.Open("commands.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Println(err)
	}
	log.Println("Successfully Opened commands.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	log.Printf("Loading commands to the memory...\n")
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &Commands)
	if err != nil {
		log.Printf("Unable to load commands from commands.json: %s\n", err)
		return
	}
	log.Printf("Commands loaded: \n")
	for _, item := range Commands {
		log.Printf("New command name: %s info: %s\n", item.Name, item.Info)
		for i := 0; i < len(item.Commands); i++ {
			log.Printf("Execution list%d: %v\n", i, item.Commands[i])
		}
	}
	AddDummyStaticVariables()
}

func SaveCommands() error {
	content, err := json.Marshal(Commands)
	if err != nil {
		log.Printf("Unable to marshal commands data err: %s\n", err)
		return err
	}
	err = ioutil.WriteFile("commands.json", content, 0644)
	if err != nil {
		log.Printf("Unable to write to commands file: %s\n", err)
		return err
	}
	log.Printf("Commands have been successfully saved!\n")
	return nil
}

func FindCommand(uid string) (Command, error) {
	for _, item := range Commands {
		if item.Uid == uid {
			return item, nil
		}
	}
	return Command{}, errors.New("Unable to find command item")
}

func PrepareRunCommand(uid string) error {
	item, err := FindCommand(uid)
	if err != nil {
		log.Printf("Unable to find command by uid: %s\n", err)
		return err
	}
	// create new dummy struct and use deep copy mechanism to pass the command
	var itm Command
	Clone(item, &itm)
	go runCmd(itm)
	return nil
}

// when passed with array returns application and it's arguments to the different variables
func ReturnProgramWithParams(params []string, uuid string) (string, []string) {
	appExe := params[0]
	otherParams := params[1:]
	return appExe, ReplaceCommandVariables(otherParams, uuid)
}

func ReplaceCommandVariables(params []string, uuid string) []string {
	//params := &params_orig
	for i := 0; i < len(params); i++ {
		// dinamic user defined variables
		for _, item := range settings.Variables {
			if strings.Contains(params[i], item.Variable) {
				//params[i] == item.Variable{
				TextReplacer := strings.NewReplacer(item.Variable, item.Setting)
				params[i] = TextReplacer.Replace(params[i])
			}
		}
		// static variables
		TextReplacer := strings.NewReplacer(
			"{datetimenow}", ReturnTimePrefix(),
			//"\u003e", ">>",
			"{uuid}", uuid,
		)
		params[i] = TextReplacer.Replace(params[i])
	}
	return params
}

func runKill() {
	err := ProgramExecution.Process.Signal(syscall.SIGTERM)
	if err != nil {
		log.Printf("Unable to kill the process!\n")
		log.Printf("Trying to terminate it by force...\n")
		ProgramExecution.Process.Kill()
	}
}

func cmdElapsedTime() string {
	elapsed := time.Since(CurrentProgram.StartTime)
	return fmt.Sprintf("Process %s is running for %s", CurrentProgram.Exe, elapsed.Round(time.Second))
}

func runCmd(command Command) {
	if !CurrentProgram.running {
		//command := &command_orig
		writer([]byte(fmt.Sprintf("Launching application: %s\n", command.Name)))
		// set the uid that share all commands in this application
		// file system safe uuid generation for filenames etc...
		alphabet := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxy"
		uuid := shortuuid.NewWithAlphabet(alphabet)

		for i := 0; i < len(command.Commands); i++ {
			app, params := ReturnProgramWithParams(command.Commands[i], uuid)
			writer([]byte(fmt.Sprintf("We are now trying to run the command: %s with params: %v\n", app, params)))
			log.Printf("Run command %s with parameters %v\n", app, params)
			ProgramExecution = exec.Command(app, params...)
			//cmd := exec.Command(app, params...)
			stdout, err := ProgramExecution.StdoutPipe()
			if err != nil {
				log.Println(err)
				return
			}
			stderr, err := ProgramExecution.StderrPipe()
			if err != nil {
				log.Println(err)
				return
			}

			if err := ProgramExecution.Start(); err != nil {
				log.Println(err)
				return
			}
			CurrentProgram.running = true
			CurrentProgram.Exe = command.Commands[i][0]
			CurrentProgram.Uid = command.Uid
			CurrentProgram.info = command.Info
			CurrentProgram.StartTime = time.Now()

			s := bufio.NewScanner(io.MultiReader(stdout, stderr))
			for s.Scan() {
				log.Println(string(s.Bytes()))
				go writer(s.Bytes())
			}

			exit_code := 0

			if err := ProgramExecution.Wait(); err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					exit_code = exitError.ExitCode()
				}
				log.Println(err)
			}
			CurrentProgram.running = false
			broadcast <- []byte(fmt.Sprintf("Process exited (%d)!", exit_code))
		}

	} else {
		broadcast <- []byte("Process is already running!")
	}
}
