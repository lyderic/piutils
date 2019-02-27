package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	if !onpi() {
		fmt.Println("This program report a Raspberry Pi CPU temperature.")
		fmt.Println("It works only on a Raspberry Pi!")
		return
	}
	content, err := ioutil.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		panic(err)
	}
	cpuTemperature, err := strconv.ParseFloat(strings.TrimSpace(string(content)), 64)
	if err != nil {
		panic(err)
	}
	cpuTemperature = cpuTemperature / 1000
	fmt.Printf("%.1f\u00b0\n", cpuTemperature)
}

func onpi() bool {
	modelFile := "/proc/device-tree/model"
	if _, err := os.Stat(modelFile); os.IsNotExist(err) {
		return false
	}
	content, err := ioutil.ReadFile(modelFile)
	if err != nil {
		panic(err)
	}
	return strings.HasPrefix(string(content), "Raspberry Pi")
}
