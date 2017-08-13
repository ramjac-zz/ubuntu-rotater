package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
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

	go func() {
		buf := bufio.NewReader(stdout)
		for {
			if out, _, err := buf.ReadLine(); err != nil {
				if err == io.EOF {
					log.Println("Done reading")
					break
				}
				log.Fatal(err)
			} else {
				log.Println(string(out))
				parsed := string(out)

				// If we can't find a colon, this isn't a line we're looking for
				if len(strings.Split(parsed, ":")) <= 1 {
					continue
				}

				// clean up the orientation
				parsed = strings.Split(parsed, ":")[1]
				parsed = strings.Replace(parsed, ")", "", -1)
				parsed = strings.TrimSpace(parsed)
				// Now parse the line and do something with it
				switch parsed {
				case "normal":
					//xrandr -o normal
					tmp := exec.Command("xrandr", "-o", "normal")
					tmp.Run()
				case "bottom-up":
					//xrandr -o inverted
					tmp := exec.Command("xrandr", "-o", "inverted")
					tmp.Run()
				case "left-up":
					//xrandr -o left
					tmp := exec.Command("xrandr", "-o", "left")
					tmp.Run()
				case "right-up":
					//xrandr -o right
					tmp := exec.Command("xrandr", "-o", "right")
					tmp.Run()
				}
			}
		}
	}()

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
