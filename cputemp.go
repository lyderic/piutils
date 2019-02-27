package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
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
	cpuTemperature, err := strconv.Atoi(string(content))
	if err != nil {
		panic(err)
	}
	fmt.Println(cpuTemperature)
}
