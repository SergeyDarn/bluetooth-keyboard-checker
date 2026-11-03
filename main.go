package main

import (
	"flag"
	"fmt"
	"os/exec"
	"regexp"
	"time"

	"github.com/charmbracelet/lipgloss"
)

const (
	errorColor   = lipgloss.Color("204")
	successColor = lipgloss.Color("#75d57f")
)

func color(output string, color lipgloss.Color) string {
	return lipgloss.NewStyle().Foreground(lipgloss.TerminalColor(color)).Render(output)
}

func logError(err string) {
	fmt.Println(color(err, errorColor))
}

func logSuccess(text string) {
	fmt.Println(color(text, successColor))
}

func wait(delay time.Duration) {
	time.Sleep(delay * time.Second)
}

func getKeyboardStatus() (bool, string) {
	bluetoothDevices, err := exec.Command("system_profiler", "SPBluetoothDataType").Output()
	if err != nil {
		return false, "Error while getting bluetooth devices: " + err.Error()
	}

	keyboardConnectedRegex := "(?i)connected(.|\n)*keyboard(.|\n)*not connected"
	isConnected, err := regexp.MatchString(keyboardConnectedRegex, string(bluetoothDevices))
	if err != nil {
		return false, "Error while getting keyboard status: " + err.Error()
	}

	return isConnected, ""
}

func main() {
	defaultDelayInSeconds := 5
	delayPtr := flag.Int("delay", defaultDelayInSeconds, "Delay between checks")
	flag.Parse()

	delay := time.Duration(*delayPtr)

	for {
		isConnected, err := getKeyboardStatus()
		if err != "" {
			logError(err)
			wait(delay)
			continue
		}

		if isConnected {
			logSuccess("Keyboard is connected")
		} else {
			logError("Keyboard is not connected")
		}

		wait(delay)
	}
}
