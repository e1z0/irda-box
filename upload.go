package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/e1z0/obexftp"
)

type IRUpload struct {
	Disabled bool `json:"disabled"`
    Running bool `json:"running"`
	LastError string `json:"last_error"`
	StartTime time.Time `json::"start_time"`
}

var (
	OBEX_OPEN_RETRY = 4
	OBEX_CONNECT_RETRY = 4
	OBEX_SEND_RETRY = 4
)

var IrUp = IRUpload{}


func Cleanup(filename string) {
	if _, err := os.Stat(filename); err == nil {
		err := os.Remove(filename) 
		if err != nil {
			log.Printf("Unable to remove file: %s err: %s\n",filename,err)
		}
	  }
}

func Upload(filename string) error {
	if !IrUp.Disabled && !IrUp.Running && !PPPDaemon.running {
		IrUp.Running = true
		if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
			IrUp.Running = false
			return err
		}
		filebase := filepath.Base(filename)
		var client obexftp.ObexFTPClient
		retry_count := 0
		// open IrDA device
		for {
			retry_count++
			cli, err := obexftp.Open()
    		if err != nil {
				if retry_count > OBEX_OPEN_RETRY {
				log.Printf("Finished trying to open IrDA device!\n")
				IrUp.Running = false
				return err	
				}
        		log.Printf("Error while trying to open IrDA device: %s\n",err)
    		} else {
				log.Printf("Connected to IrDA device!\n")
				client = cli
				break
			}
			log.Printf("Retrying to connect to IrDA device...\n")
            time.Sleep(2*time.Second)
		}
		retry_count = 0
		// connect to remote host over IrDA
        for {
			retry_count++
			err := obexftp.Connect(client)
    		if err != nil {
                if retry_count > OBEX_CONNECT_RETRY {
					log.Printf("Finished trying to connect to remote host!\n")
					IrUp.Running = false
					return err
				}
				log.Printf("Error while trying to connect to remote host!\n")
    		} else {
				break
			}
			log.Printf("Retrying to connect to remote device...\n")
			time.Sleep(2*time.Second)
		}
		retry_count = 0
		// push the file to the remote device
		for {
			retry_count++
			err := obexftp.Push(client,filename,filebase)
			if err != nil {
				if retry_count > OBEX_SEND_RETRY {
					log.Printf("Finished trying to send file to the remote host!\n")
					IrUp.Running = false
					return err
				}
				log.Printf("Trying to send file to the remote host!\n")
			} else {
				break
			}
			log.Printf("Retrying to send file to remote device...\n")
			time.Sleep(2*time.Second)
		}
		// disconnect and close the IrDA device handle
		err := obexftp.Disconnect(client)
		if err != nil {
			log.Printf("Unable to disconnect: %s\n",err)
		}
		obexftp.Close(client)
		IrUp.Running = false
	}

	Cleanup(filename)

	return nil
}