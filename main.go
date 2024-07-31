package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

var (
	reboot     string
	journalctl string
	systemctl  string
	err        error
)

func main() {
	fmt.Println("Starting soukiller")
	journalctl, err = exec.LookPath("journalctl")
	checkError(err)

	reboot, err = exec.LookPath("reboot")
	checkError(err)

	systemctl, err = exec.LookPath("systemctl")
	checkError(err)

	fmt.Println("journalctl:", journalctl)
	fmt.Println("reboot:", reboot)
	fmt.Println("systemctl:", systemctl)

	monitor := exec.Command("journalctl", "-b", "-f")
	stdout, err := monitor.StdoutPipe()
	checkError(err)
	err = monitor.Start()
	checkError(err)

	defer monitor.Wait()
	scanner := bufio.NewScanner(stdout)
	readStuff(scanner)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func shellExec(command string, options ...string) string {

	cmd := exec.Command(command, options...)
	err := cmd.Run()
	checkError(err)
	result, err := cmd.Output()
	checkError(err)
	return string(result)
}

func readStuff(scanner *bufio.Scanner) {
	counter := 1000
	for scanner.Scan() {
		//fmt.Printf("New data:")
		logLine := scanner.Text()
		//fmt.Println(logLine)
		if grep("RIP: 0010:_nv012398rm+0xbd/0x130", logLine) {
			fmt.Println("NVIDIA kernel crash detected:", logLine)
			fmt.Println("Rebooting the system")
			shellExec(reboot, "-f")
		}

		if grep("nonono", logLine) {
			fmt.Println("Wayland crash detected:", logLine)
			fmt.Println("Restart SDDM")
			shellExec(systemctl, "restart", "--no-block", "sddm")
		}
		counter--
		if counter == 0 {
			fmt.Println("soulkiller stills alive")
			counter = 1000
		}
	}
	err := scanner.Err()
	checkError(err)
	fmt.Fprintln(os.Stderr, "reading standard input:", err)
}

func grep(pattern string, line string) bool {
	matched, err := regexp.MatchString(pattern, line)
	checkError(err)
	return matched

}
