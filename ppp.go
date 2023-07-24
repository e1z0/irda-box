package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"syscall"
	"time"
)

var PPPExecution *exec.Cmd
var ppp_bin = ""
var ppp_args = ""

type PPP struct {
	PPPBinary string
	PPPArgs   []string
	running   bool
	info      string
	StartTime time.Time
}

var PPPDaemon = PPP{}

func PPPRuntime() string {
	elapsed := time.Since(PPPDaemon.StartTime)
	return fmt.Sprintf("PPP is running for %s", elapsed.Round(time.Second))
}

func StartPPP() {
	if !PPPDaemon.running {
		pppwriter([]byte(fmt.Sprintf("Launching application: %s\n", PPPDaemon.PPPBinary)))
		pppwriter([]byte(fmt.Sprintf("We are now trying to run the command: %s with params: %v\n", PPPDaemon.PPPBinary, PPPDaemon.PPPArgs)))
		log.Printf("Run command %s with parameters %v\n", PPPDaemon.PPPBinary, PPPDaemon.PPPArgs)
		PPPExecution = exec.Command(PPPDaemon.PPPBinary, PPPDaemon.PPPArgs...)
		stdout, err := PPPExecution.StdoutPipe()
		if err != nil {
			log.Println(err)
			return
		}
		stderr, err := PPPExecution.StderrPipe()
		if err != nil {
			log.Println(err)
			return
		}

		if err := PPPExecution.Start(); err != nil {
			log.Println(err)
			return
		}
		PPPDaemon.running = true
		PPPDaemon.StartTime = time.Now()

		s := bufio.NewScanner(io.MultiReader(stdout, stderr))
		for s.Scan() {
			log.Println(string(s.Bytes()))
			go pppwriter(s.Bytes())
		}

		exit_code := 0

		if err := PPPExecution.Wait(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				exit_code = exitError.ExitCode()
			}
			log.Println(err)
		}
		PPPDaemon.running = false
		ppp_broadcast <- []byte(fmt.Sprintf("PPP Process exited (%d)!", exit_code))

	} else {
		ppp_broadcast <- []byte("PPP Process is already running!")
	}
}

func StopPPP() {
	err := PPPExecution.Process.Signal(syscall.SIGTERM)
	if err != nil {
		log.Printf("Unable to kill the ppp process!\n")
		log.Printf("Trying to terminate ppp process by force...\n")
		PPPExecution.Process.Kill()
	}
}

func RestartPPP() {
	StopPPP()
	time.Sleep(3 * time.Second)
	StartPPP()
}
