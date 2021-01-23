package main

import (
	"context"
	"log"
	"net/http"
	"os/exec"
	"sync"

	flag "github.com/spf13/pflag"
)

type State int

const (
	Running State = iota
	Ready
)

type Worker struct {
	status State
	mx     *sync.Mutex
}

var (
	command = flag.String("command", "", "command to execute")
	args    = flag.StringSlice("args", []string{}, "arguments to pass to command")
)

func (w *Worker) Start(ctx context.Context) ([]byte, error) {
	w.mx.Lock()
	defer w.mx.Unlock()
	w.status = Running
	defer w.SetTerminated()
	log.Println("call: ", *command, (*args))
	o, err := exec.CommandContext(ctx, *command, *args...).Output()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return o, nil
}

func (w *Worker) SetTerminated() {
	w.status = Ready
}

func (w *Worker) Status() string {
	switch w.status {
	case Running:
		return "running"
	default:
		return "ready"
	}
}

func main() {
	flag.Parse()
	if *command == "" {
		log.Fatalln("Please provide command")
	}
	worker := &Worker{
		status: Ready,
		mx:     &sync.Mutex{},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		if worker.status == Running {
			http.Error(w, "target is busy", http.StatusServiceUnavailable)
			return
		}
		o, err := worker.Start(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(o)
		return
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(worker.Status()))
		return
	})
	log.Fatal(http.ListenAndServe(":8080", mux))
}
