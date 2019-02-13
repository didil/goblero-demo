package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/didil/goblero/pkg/blero"
)

func main() {
	bl := blero.New("db/")
	err := bl.Start()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Registering Processor 1 ...")
	bl.RegisterProcessorFunc(func(j *blero.Job) error {
		fmt.Printf("[Processor 1] Processing job: %v - data: %v\n", j.Name, string(j.Data))
		time.Sleep(2 * time.Second)
		fmt.Printf("[Processor 1] Done Processing job: %v\n", j.Name)

		return nil
	})

	if len(os.Args) > 1 && os.Args[1] == "enqueue" {
		fmt.Println("Enqueuing jobs ...")
		for i := 1; i <= 50; i++ {
			jobName := fmt.Sprintf("Job #%v", i)
			jobData := []byte(fmt.Sprintf("Job Data #%v", i))
			_, err := bl.EnqueueJob(jobName, jobData)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// Wait for SIGTERM or SIGINT to stop Blero and exit
	var exitCh = make(chan os.Signal)
	signal.Notify(exitCh, syscall.SIGINT)
	signal.Notify(exitCh, syscall.SIGTERM)
	s := <-exitCh
	fmt.Printf("Caught signal %v. Exiting ...\n", s)
	bl.Stop()

	os.Exit(0)
}
