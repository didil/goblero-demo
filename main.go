package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/didil/goblero/pkg/blero"
)

func main() {
	// Parse flags
	n := flag.Int("n", 1, "number of processors")
	flag.Parse()

	// Create a new Blero backend
	bl := blero.New("db/")
	// Start Blero
	err := bl.Start()
	if err != nil {
		log.Fatal(err)
	}

	// Register processor(s)
	for i := 1; i <= *n; i++ {
		pI := i
		fmt.Printf("Registering Processor %v ...\n", pI)
		bl.RegisterProcessorFunc(func(j *blero.Job) error {
			fmt.Printf("[Processor %v] Processing job: %v - data: %v\n", pI, j.Name, string(j.Data))
			// Simulate processing
			time.Sleep(2 * time.Second)
			fmt.Printf("[Processor %v] Done Processing job: %v\n", pI, j.Name)

			return nil
		})
	}

	// Enqueue jobs
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

	// Stop Blero
	bl.Stop()
	os.Exit(0)
}
