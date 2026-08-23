package captor

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var (
	callback func([]byte)
	ticker *time.Ticker
	stopCh chan bool
)

func Init(cb func([]byte)) {
	callback = cb
	stopCh = make(chan bool)
}

func HandleStart(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Interval int `json:"interval"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if req.Interval <= 0 { req.Interval = 5 }

	if ticker!= nil { ticker.Stop() }
	ticker = time.NewTicker(time.Duration(req.Interval) * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				doCapture()
			case <-stopCh:
				return
			}
		}
	}()
	log.Printf("Monitor iniciado: %d min\n", req.Interval)
	w.Write([]byte("started"))
}

func HandleStop(w http.ResponseWriter, r *http.Request) {
	if ticker!= nil { ticker.Stop() }
	select { case stopCh <- true: default: }
	w.Write([]byte("stopped"))
}

func HandleCapture(w http.ResponseWriter, r *http.Request) {
	go doCapture()
	w.Write([]byte("captured"))
}

func doCapture() {
	img, err := CaptureScreen()
	if err!= nil {
		log.Println("erro captura:", err)
		return
	}
	if callback!= nil {
		callback(img)
	}
}
