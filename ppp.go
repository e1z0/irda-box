package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"syscall"
	"time"
        "strings"
)

var PPPExecution *exec.Cmd
//var ppp_bin = ""
//var ppp_args = ""

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
//                var mergedParams []string
//                var userparams []string
                userparams := strings.Split(settings.PPPSettings.Command," ")
                mergedParams := append(PPPDaemon.PPPArgs,userparams...)
		pppwriter([]byte(fmt.Sprintf("Launching application: %s\n", PPPDaemon.PPPBinary)))
		pppwriter([]byte(fmt.Sprintf("We are now trying to run the command: %s with params: %v\n", PPPDaemon.PPPBinary, mergedParams)))
		log.Printf("Run command %s with parameters %v\n", PPPDaemon.PPPBinary, mergedParams)
		PPPExecution = exec.Command(PPPDaemon.PPPBinary, mergedParams...)
                PPPExecution.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
		stdout, err := PPPExecution.StdoutPipe()
		if err != nil {
			log.Printf("PPP Exec stdout capture error: %s\n",err)
			return
		}
		stderr, err := PPPExecution.StderrPipe()
		if err != nil {
			log.Printf("PPP Exec stderr capture error: %s\n",err)
			return
		}

		if err := PPPExecution.Start(); err != nil {
                        log.Printf("PPP Execute error: %s\n",err)
			return
		}
		PPPDaemon.running = true
		PPPDaemon.StartTime = time.Now()

		s := bufio.NewScanner(io.MultiReader(stdout, stderr))
		for s.Scan() {
			log.Println(string(s.Bytes()))
                        if len(string(s.Bytes())) > 0 {
			go pppwriter(s.Bytes())
                        }
		}

		exit_code := 0

		if err := PPPExecution.Wait(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				exit_code = exitError.ExitCode()
			}
                        log.Printf("PPP Execution finished: %s\n",err)
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
