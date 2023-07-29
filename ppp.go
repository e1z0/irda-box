package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
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

func GenPPPConfigs() (bool, error) {
	ppp_peer := "/etc/ppp/peers/irda"
	ppp_connect := "/etc/ppp/connect.sh"
	ppp_disconnect := "/etc/ppp/disconnect.sh"

	TextReplacer := strings.NewReplacer(
		"{interface}", settings.WifiIface,
		"{ircomm}", settings.PPPSettings.IrComm,
		"{speed}", strconv.Itoa(settings.PPPSettings.Speed),
	)

    ppp_config := "connect /etc/ppp/connect.sh\ndisconnect /etc/ppp/disconnect.sh\n"+settings.PPPSettings.Options+"\n{speed}\n{ircomm}\n"

	// Create the directory if it does not exist
    path := "/etc/ppp/peers"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return false,err
		}
	}
	// WRITE PPP Profile
	// open file
	f, err := os.OpenFile(ppp_peer, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
    if err != nil {
        return false,err
    }
	// write to file
	pppconf := TextReplacer.Replace(ppp_config)
	_, err = f.WriteString(pppconf)
	if err != nil {
		return false,err
	}

	// close the file
    if err := f.Close(); err != nil {
        return false,err
    }
	// WRITE PPP Connect script
	f, err = os.OpenFile(ppp_connect, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
    if err != nil {
        return false,err
    }
	// write to file
	pppconn := TextReplacer.Replace(settings.PPPSettings.Connect)
	_, err = f.WriteString(pppconn+"\n")
	if err != nil {
		return false,err
	}

	// close the file
    if err := f.Close(); err != nil {
        return false,err
    }

	/// WRITE PPP Disconnect script
	f, err = os.OpenFile(ppp_disconnect, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
    if err != nil {
        return false,err
    }
	// write to file
	pppdisconn := TextReplacer.Replace(settings.PPPSettings.Disconnect)
	_, err = f.WriteString(pppdisconn+"\n")
	if err != nil {
		return false,err
	}

	// close the file
    if err := f.Close(); err != nil {
        return false,err
    }

	return true,nil
}

func PPPEvents(text string) {
   switch {
    // connected event
    case strings.Contains(text, "remote IP address"):
      log.Printf("Seems like peer connection have established")
      TurnLed("blue",true)
      break
    // peer disconnected
    case strings.Contains(text,"LCP terminated by peer"):
      log.Printf("Seems like peer connection have dropped")
      TurnLed("blue",false)
      break
    // dropped connection
    case strings.Contains(text,"Modem hangup"):
      log.Printf("Seems like peer connection have dropped")
      TurnLed("blue",false)
    default:
      break
    }
}

func StartPPP() {
	if !PPPDaemon.running && !IrUp.Running {

        userparams := strings.Split(settings.PPPSettings.Command," ")
        mergedParams := append(PPPDaemon.PPPArgs,userparams...)
		pppwriter([]byte(fmt.Sprintf("Generating config for ppp\n")))
		log.Printf("Generating config for pppp\n")
		ok, err := GenPPPConfigs()
		if err != nil || !ok {
			pppwriter([]byte(fmt.Sprintf("Config generation for ppp have failed: %s\n",err)))
			log.Printf("Config generation for ppp have failed: %s\n",err)
			return
		}
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
                          PPPEvents(string(s.Bytes()))
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
	if PPPDaemon.running {
		StopPPP()
	}
	time.Sleep(3 * time.Second)
	StartPPP()
}
