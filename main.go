package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
)

func main() {
	log.Println("So it goes...")

	// Handle cancellation
	ctx := context.Background()
	// trap Ctrl+C and call cancel on the context
	c := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(ctx)
	signal.Notify(c, os.Interrupt)

	defer func() {
		signal.Stop(c)
	}()

	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	ms := exec.CommandContext(ctx, "monitor-sensor")
	stdout, err := ms.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := ms.Start(); err != nil {
		log.Fatal(err)
	}

	if out, err := ioutil.ReadAll(stdout); err != nil {
		log.Fatal(err)
	} else {
		log.Println(string(out))
	}

	if err := ms.Wait(); err != nil {
		log.Fatal(err)
	}
}

/* an example of some output from monitor-sensors
Has accelerometer (orientation: normal)
=== No ambient light sensor
    Accelerometer orientation changed: left-up
    Accelerometer orientation changed: normal
    Accelerometer orientation changed: right-up
    Accelerometer orientation changed: left-up
    Accelerometer orientation changed: bottom-up
    Accelerometer orientation changed: left-up
    Accelerometer orientation changed: normal
*/
