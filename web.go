package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []byte)
var ppp_broadcast = make(chan []byte)
var socketLock bool

func httpPool() {
	router := mux.NewRouter()
	router.HandleFunc("/", rootHandler).Methods("GET")
	router.HandleFunc("/term", termHandler).Methods("GET") // for testing purposes only
	router.HandleFunc("/cmd", cmdHandler).Methods("POST")
	router.HandleFunc("/upload", UploadHandler).Methods(("POST"))
	router.HandleFunc("/kill", killHandler).Methods(("POST"))
	router.HandleFunc("/ppp-start", pppstartHandler).Methods(("POST"))
	router.HandleFunc("/ppp-stop", pppstopHandler).Methods(("POST"))
	router.HandleFunc("/ppp-restart", ppprestartHandler).Methods(("POST"))
	router.HandleFunc("/settings", settingsHandler).Methods("GET", "POST")
	router.HandleFunc("/commands", commandsHandler).Methods("GET", "POST")
	router.HandleFunc("/ws", wsHandler)

	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	go broadcaster()
	go ppp_broadcaster()
	log.Fatal(http.ListenAndServe(":8000", router))
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	type HtmlStruct struct {
		Name     string
		Title    string
		Footer   string
		Commands []Command
	}
	t, err := template.New("home.html").Funcs(template.FuncMap{
		"StringsJoin": strings.Join,
	}).ParseFiles("tpl/home.html")
	t.New("header").ParseFiles("tpl/header.html")
	t.New("footer").ParseFiles("tpl/footer.html")

	if err != nil {
		log.Printf("Unable to parse template: %s\n", err)
	}
	err = t.ExecuteTemplate(w, "home.html", HtmlStruct{Name: static_variables.Name, Title: "Home", Footer: static_variables.Footer, Commands: Commands})
	if err != nil {
		log.Printf("Error when parsing html template: %s\n", err)
		http.Error(w, "Internal Server Error", 500)
	}
}

func termHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/term" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "term.html")
}

func writer(msg []byte) {
	broadcast <- msg
}

func statusUpdate() {

	broadcast <- []byte("just_a_status_update")
}

func pppwriter(msg []byte) {
	ppp_broadcast <- msg
}

func killHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if CurrentProgram.running {
		runKill()
		w.WriteHeader(200)
		json.NewEncoder(w).Encode("ok")
		return
	}
	log.Printf("You want to kill something but nothing left alive :(\n")
	w.WriteHeader(500)
	json.NewEncoder(w).Encode("error")
}

func pppstopHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if PPPDaemon.running {
		StopPPP()
		w.WriteHeader(200)
		json.NewEncoder(w).Encode("ok")
		return
	}
	log.Printf("You want to kill something but nothing left alive :(\n")
	w.WriteHeader(500)
	json.NewEncoder(w).Encode("error")
}

func pppstartHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Printf("ppp start requested!\n")
	defer r.Body.Close()
	StartPPP()
	//		if err == nil {
	w.WriteHeader(200)
	json.NewEncoder(w).Encode("ok")
	return
	//		}
	//	}

	// w.WriteHeader(500)
	// json.NewEncoder(w).Encode(err)
}

func ppprestartHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Printf("ppp restart requested!\n")
	defer r.Body.Close()
	RestartPPP()
	w.WriteHeader(200)
	json.NewEncoder(w).Encode("ok")
	return
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//w.Header().Set("Content-Type", "application/json")
	//
	// Parse our multipart form, 10 << 20 specifies a maximum
    // upload of 10 MB files.
    err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return
	}
	files := r.MultipartForm.File["file"]
	for i, _ := range files {
			log.Printf("Uploading file: %s\n",files[i].Filename)
			file, err := files[i].Open()
 			if err != nil {
				w.WriteHeader(500)
				fmt.Fprint(w, err)
 				return
 			}

			 fp := filepath.Join("/tmp/", files[i].Filename)
			 out, err := os.Create(fp)
			 defer out.Close()
			 if err != nil {
				w.WriteHeader(500)
				fmt.Fprint(w, err)
				return
			 }

			 _, err = io.Copy(out, file)

			 if err != nil {
				 log.Printf("Error storing file %s data, err: %s\n",files[i].Filename,err)
				 w.WriteHeader(500)
				 fmt.Fprint(w, err)
				 return
			 }
			 err = Upload(fp)
			 if err != nil {
				log.Printf("Unable to upload: %s\n",err)
				w.WriteHeader(500)
				fmt.Fprint(w, err)
				return
			 }
		}
		log.Printf("Upload finished!")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode("ok")

}

func cmdHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-Type", "application/json")
	log.Printf("Cmd handler requested!\n")
	r.ParseForm()
	uid := r.Form.Get("uid")
	defer r.Body.Close()
	if uid != "" {
		log.Printf("Got request to run command with uid: %s\n", uid)
		err = PrepareRunCommand(uid)
		if err == nil {
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(uid)
			return
		}
	}

	w.WriteHeader(500)
	json.NewEncoder(w).Encode(err)
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		log.Printf("got some post")
		sets := Settings{}
		err := json.NewDecoder(r.Body).Decode(&sets)
		if err != nil {
			log.Printf("Can't decode Json received in post %s\n", err)
			//w.Header().Set("Content-Type", "application/json")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
log.Printf("new settings: %#v\n",sets)
		settings = sets
		err = saveSettings()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"data": "ok"})
		return
	} else {

		type HtmlStruct struct {
			Name     string
			Title    string
			Footer   string
			Settings Settings
                        Ifaces   []Iface
                        IrcommFaces []string
		}
                ifaces,_ := ListInterfaces()
                comfaces, _ := ReturnIrcommIfaces()
		t, err := template.New("settings.html").Funcs(template.FuncMap{
			"StringsJoin": strings.Join,
		}).ParseFiles("tpl/settings.html")
		t.New("header").ParseFiles("tpl/header.html")
		t.New("footer").ParseFiles("tpl/footer.html")

		if err != nil {
			log.Printf("Unable to parse template: %s\n", err)
		}
		err = t.ExecuteTemplate(w, "settings.html", HtmlStruct{Name: static_variables.Name, Title: "Settings", Footer: static_variables.Footer, Settings: settings, Ifaces: ifaces, IrcommFaces: comfaces})
		if err != nil {
			log.Printf("Error when parsing html template: %s\n", err)
			http.Error(w, "Internal Server Error", 500)
		}
	}
}

func commandsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		log.Printf("Got commands save request\n")
		coms := []Command{}
		err := json.NewDecoder(r.Body).Decode(&coms)
		if err != nil {
			log.Printf("Can't decode Json received in post %s\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		Commands = coms
		err = SaveCommands()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"data": "ok"})
		return
	} else {
		type HtmlStruct struct {
			Name      string
			Title     string
			Footer    string
			Commands  []Command
			Variables []Variables
		}
		t, err := template.New("commands.html").Funcs(template.FuncMap{
			"StringsJoin": StringJoinFix,
		}).ParseFiles("tpl/commands.html")
		t.New("header").ParseFiles("tpl/header.html")
		t.New("footer").ParseFiles("tpl/footer.html")

		if err != nil {
			log.Printf("Unable to parse template: %s\n", err)
		}
		var vars []Variables
		vars = append(vars,Variables{Variable: "{interface}", Setting: settings.WifiIface})
		vars = append(vars,Variables{Variable: "{ircomm}", Setting: settings.PPPSettings.IrComm})
		vars = append(vars,Variables{Variable: "{speed}", Setting: strconv.Itoa(settings.PPPSettings.Speed)})
		vars = append(vars,settings.Variables...)
		err = t.ExecuteTemplate(w, "commands.html", HtmlStruct{Name: static_variables.Name, Title: "Commands", Footer: static_variables.Footer, Commands: Commands, Variables: vars})
		if err != nil {
			log.Printf("Error when parsing html template: %s\n", err)
			http.Error(w, "Internal Server Error", 500)
		}
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// register client
	clients[ws] = true
}

// 3

type WebSocketMessage struct {
	Messages         []byte `json:"messages"`
	Running          bool   `json:"running"`
	Eta              string `json:"eta"`
	StatusUpdate     bool   `json:"statusupdate"`
	Batteries        []Battery
	PPPStatusRunning bool   `json:"ppp_running"`
	PPPStatus        []byte `json:"ppp_status"`
	PPPRuntime       string `json:"ppp_runtime"`
	Upload           IRUpload `json:"ir_upload"`
}

func StatusLoop() {
	for {
		statusUpdate()
		time.Sleep(1 * time.Second)
	}
}

func ppp_broadcaster() {
	for {
		val := <-ppp_broadcast
		msg := WebSocketMessage{}

		//if string(val) != "just_a_status_update" {
		msg = WebSocketMessage{Messages: []byte(""), Running: CurrentProgram.running, Eta: cmdElapsedTime(), StatusUpdate: true, Batteries: batteryInfo(), PPPStatus: val, PPPStatusRunning: PPPDaemon.running, PPPRuntime: PPPRuntime(), Upload: IrUp}
		//}
		for client := range clients {
			if !socketLock {
				socketLock = true
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("PPP Broadcaster Websocket error: %s", err)
					client.Close()
					delete(clients, client)
				}
				socketLock = false

			}
		}
	}
}

func broadcaster() {
	for {
		val := <-broadcast
		msg := WebSocketMessage{}

		if string(val) == "just_a_status_update" {
			msg = WebSocketMessage{Messages: []byte(""), Running: CurrentProgram.running, Eta: cmdElapsedTime(), StatusUpdate: true, Batteries: batteryInfo(), PPPStatusRunning: PPPDaemon.running, PPPRuntime: PPPRuntime(), Upload: IrUp}
		} else {
			msg = WebSocketMessage{Messages: val, Running: CurrentProgram.running, Eta: cmdElapsedTime(), Batteries: batteryInfo(), PPPStatusRunning: PPPDaemon.running, PPPRuntime: PPPRuntime(),Upload: IrUp}
		}

		//val2 := <-ppp_broadcast
		//msg = WebSocketMessage{Messages: []byte(""), Running: CurrentProgram.running, Eta: cmdElapsedTime(), StatusUpdate: true, Batteries: batteryInfo(), PPPStatus: val2, PPPStatusRunning: PPPDaemon.running}

		for client := range clients {
			if !socketLock {
				socketLock = true
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("Websocket error: %s", err)
					client.Close()
					delete(clients, client)
				}
				socketLock = false

			}
		}
	}
}
