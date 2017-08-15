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

				penstylus := "Wacom HID 4822 Pen stylus"
				peneraser := "Wacom HID 4822 Pen eraser"
				fingertouch := "Wacom HID 4822 Finger touch"

				// Now parse the line and do something with it
				switch parsed {
				case "normal":
					setXrandr("normal")
					setXwacom(penstylus, "none")
					setXwacom(peneraser, "none")
					setXwacom(fingertouch, "none")
				case "bottom-up":
					setXrandr("inverted")
					setXwacom(penstylus, "half")
					setXwacom(peneraser, "half")
					setXwacom(fingertouch, "half")
				case "left-up":
					setXrandr("left")
					setXwacom(penstylus, "ccw")
					setXwacom(fingertouch, "ccw")
					setXwacom(fingertouch, "ccw")
				case "right-up":
					setXrandr("right")
					setXwacom(penstylus, "cw")
					setXwacom(peneraser, "cw")
					setXwacom(fingertouch, "cw")
				}
			}
		}
	}()

	if err := ms.Wait(); err != nil {
		log.Fatal(err)
	}
}

func setXrandr(dir string) {
	tmp := exec.Command("xrandr", "-o", dir)
	if err:=tmp.Run(); err !=nil{
		log.Println(err)
	}
}
func setXwacom(dev, dir string) {
	tmp := exec.Command("xsetwacom", "set", dev, "rotate", dir)
	if err:=tmp.Run(); err !=nil{
		log.Println(err)
	}
}

/* Example of xsetwacom output
xsetwacom --list devices
Wacom HID 4822 Pen stylus               id: 11  type: STYLUS
Wacom HID 4822 Finger touch             id: 12  type: TOUCH
Wacom HID 4822 Pen eraser               id: 18  type: ERASER
*/

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
