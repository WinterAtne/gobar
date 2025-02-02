package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Control
const tickrate time.Duration = time.Second * 30

// Seperators
const leftSeperator string = "["
const rightSeperator string = "]"

// Define a series of colors
const (
	defaultColor string = "FFFFFF" // White
	errorColor 			= "eb6f92" // red
	warnColor			= "f6c177" // yellow
	okColor				= "9ccfd8" // blue
	morningColor		= "ea9a97" // orange
	dayColor				= "e0def4" // white
	nightColor			= "a580d2" // dark purple
	dateColor			= "c4a7e7" // purple
)

// Makes a cell using color, saying text, surrounded by the left and right seperators
func makeCell(content string, color string) (string) {
	if color == "" { color = defaultColor }
	return fmt.Sprintf("^c#%s^ %s %s %s", color, leftSeperator, content, rightSeperator)
}

func removeRight(s string) string {
	return s[:len(s)-1]
}

func calcBattery() (string) {
	const (
		DyingPercent int 	= 25
		OkPercent			= 40
		FullPercent			= 70
	)
	const (
		Empty string	= "󰂎"
		Dying				= "󱊡"
		Ok					= "󱊢"
		Full				= "󱊣"
		EmptyCharge		= "󱊤"
		DyingCharge		= "󱊤"
		OkCharge			= "󱊥"
		FullCharge		= "󱊦"
	)
	
	const batteryPercentFile string = "/sys/class/power_supply/BAT0/capacity"
	const batteryStatusFile string = "/sys/class/power_supply/BAT0/status"

	var color string = okColor
	var percent int
	var symbol string
	var status bool = false // False is discharging
	var err error
	var data []byte // Temp variable for raw file data

	// Get the current battery percent
	data, err = os.ReadFile(batteryPercentFile)
	if err != nil {
		log.Println(err)
		return makeCell("?!", errorColor)
	}
	percent, err = strconv.Atoi(removeRight(string(data)))
	if err != nil {
		log.Println(err)
		return makeCell("?!", errorColor)
	}

	data, err = os.ReadFile(batteryStatusFile)
	if err != nil {
		log.Println(err)
		return makeCell("?!", errorColor)
	} else if data[0] == 'C' {
		status = true
	}

	if status {
		if percent >= FullPercent {
			color = warnColor
			symbol = FullCharge
		} else if percent >= OkPercent {
			symbol = OkCharge
		} else if percent >= DyingPercent {
			symbol = DyingCharge
		} else {
			symbol = EmptyCharge
		}
	} else {
		if percent >= FullPercent {
			symbol = Full
		} else if percent >= OkPercent {
			symbol = Ok
		} else if percent >= DyingPercent {
			symbol = Dying
			color = warnColor
		} else {
			symbol = Empty
			color = warnColor
		}
	}

	var content string = fmt.Sprintf("%s %d%%", symbol, percent)

	return makeCell(content, color)
}

func calcTime() (string) {
	const timeFormat string = "03:04"
	const hourFormat string = "15"

	var color string
	var outTime string = time.Now().Format(timeFormat)
	hour, _ := strconv.Atoi(time.Now().Format(hourFormat)) // Error should always be nil

	if hour >= 12 {
		outTime += "pm"
	} else {
		outTime += "am"
	}

	if 5 <= hour && 10 >= hour {
		color = morningColor
	} else if 11 <= hour && 18 >= hour {
		color = dayColor
	} else {
		color = nightColor
	}

	return makeCell(outTime, color)
}

// Calculates the date and returns it as a cell
func calcDate() (string) {
	const dateFormat string = "Mon Jan 02"
	var rtrn string = time.Now().Format(dateFormat)
	return makeCell(rtrn, dateColor)
}

// Sets the root name, setting the statusbar
func setRootName(name string) {
	cmd := exec.Command("xsetroot", "-name", name)
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
}

func main() {
	for {
		var status strings.Builder
		status.WriteString(calcBattery())
		status.WriteString(calcTime())
		status.WriteString(calcDate())

		setRootName(status.String())

		time.Sleep(tickrate)
	}
}

