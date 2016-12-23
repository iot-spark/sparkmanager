package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var (
	ManagerPort *string
	DbFile      *string
)

func init() {
	ManagerPort = flag.String("port", "8080", "MGMT manager port")
	DbFile = flag.String("dbname", "iotspark.db", "The name of an SQLite database")

	flag.Parse()
}

func main() {
	var (
		// terminate
		exit = make(chan os.Signal, 1)
	)

	signal.Notify(exit,
		os.Interrupt,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGINT,
	)

	go Manager()
	fmt.Printf("MGMT manager has been started on port %s\n", *ManagerPort)

	for {
		select {
		case <-exit:
			go func() {
				os.Exit(0)
			}()
		}
	}
}
