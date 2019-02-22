// +build linux

package main

import (
	"errors"
	"log"
	"os"

	ps "github.com/mitchellh/go-ps"
	"github.com/rjeczalik/notify"
)

var (
	directory string
)

func isDirectory(directoryPath string) (bool, error) {
	path, err := os.Stat(directoryPath)

	switch {
	case err != nil:
		log.Fatal(err)
	case path.IsDir():
		return true, nil
	}
	return false, err

}

func FindProcess(key string) (int, string, error) {
	pname := ""
	pid := 0
	err := errors.New("not found")
	ps, _ := ps.Processes()

	for i, _ := range ps {
		if ps[i].Executable() == key {
			pid = ps[i].Pid()
			pname = ps[i].Executable()
			err = nil
			break
		}
	}
	return pid, pname, err
}

func hupProcess(pid int, s string) error {
	log.Printf("HUPing process %s with PID of %d", s, pid)

	proc, err := os.FindProcess(pid)
	if err != nil {
		log.Fatal(err)
	}
	proc.Signal(os.Interrupt)
	return err
}

func main() {

	directory, ok := os.LookupEnv("WATCH_DIRECTORY")
	if !ok {
		log.Fatal("Could not find WATCH_DIRECTORY environment vairable, exiting.")
	}

	hupTarget, ok := os.LookupEnv("HUP_PROCESS_NAME")
	if !ok {
		log.Fatal("Could not find HUP_PROCESS_NAME environment vairable, exiting.")
	}

	log.Printf("Watching directory %s and will HUP %s", directory, hupTarget)

	hupTarget = hupTarget[0:15]

	_, err := isDirectory(directory)
	if err != nil {
		log.Fatal(err)
	}

	for {

		c := make(chan notify.EventInfo, 1)

		if err := notify.Watch(directory, c, notify.InCloseWrite, notify.InMovedTo); err != nil {
			log.Fatal(err)
		}
		defer notify.Stop(c)

		switch ei := <-c; ei.Event() {
		case notify.InCloseWrite, notify.InMovedTo:
			log.Printf("Found file %s", ei.Path())

			if pid, name, err := FindProcess(hupTarget); err == nil {
				hupProcess(pid, name)
			}

		}
	}
}
