package main

import (
	"fmt"
	"os/exec"
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

const timeFormat string = "03:04PM"
func calcTime() (string) {
	var time string = time.Now().Format(timeFormat)
	return makeCell(time, nightColor)
}

// Calculates the date and returns it as a cell
const dateFormat string = "Mon Jan 02"
func calcDate() (string) {
	var date string = time.Now().Format(dateFormat)
	return makeCell(date, dateColor)
}

// Sets the root name, setting the statusbar
func setRootName(name string) (error) {
	cmd := exec.Command("xsetroot", "-name", name)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var err error = nil
	var index int = 0
	var status strings.Builder
	for err == nil {
		status.Write([]byte(calcTime()))
		status.Write([]byte(calcDate()))
		err = setRootName(status.String())
		index++

		status.Reset()
		time.Sleep(tickrate)
	}

	panic(err)
}
