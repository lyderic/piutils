package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/lyderic/tools"
)

type Information struct {
	Hostname   string
	Model      string
	Celsius    string
	Farenheit  string
	Networking []NIC
}

type NIC struct {
	Name      string
	IpAddress string
	State     string
}

const (
	VERSION       = "0.0.2"
	MODEL_FILE    = "/proc/device-tree/model"
	CPU_TEMP_FILE = "/sys/class/thermal/thermal_zone0/temp"
)

var (
	showHostname    bool
	showModel       bool
	showTemperature bool
	showCelsius     bool
	showFarenheit   bool
	showNetworking  bool
	showAll         bool
)

func main() {
	flag.BoolVar(&showHostname, "h", false, "show hostname")
	flag.BoolVar(&showModel, "m", false, "show Rasperry Pi model")
	flag.BoolVar(&showTemperature, "t", false, "show temperature")
	flag.BoolVar(&showCelsius, "c", false, "show temperature (celsius only)")
	flag.BoolVar(&showFarenheit, "f", false, "show temperature (farenheit only)")
	flag.BoolVar(&showNetworking, "n", false, "show networking")
	flag.BoolVar(&showAll, "a", false, "show all information")
	flag.Usage = usage
	flag.Parse()
	information := getInformation()
	if len(os.Args) == 1 || showAll {
		fmt.Println(information)
		return
	}
	if showHostname {
		fmt.Println(information.Hostname)
	}
	if showModel {
		fmt.Println(information.Model)
	}
	if showTemperature {
		fmt.Println(information.Celsius, information.Farenheit)
	}
	if showCelsius {
		fmt.Println(information.Celsius)
	}
	if showFarenheit {
		fmt.Println(information.Farenheit)
	}
	if showNetworking {
		fmt.Println(displayNetworking(information.Networking))
	}
}

func getInformation() (i Information) {
	i.Hostname = getHostname()
	i.Model = getModel()
	i.Celsius = getCelsius()
	i.Farenheit = getFarenheit()
	i.Networking = getNetworking()
	return
}

func getHostname() (hostname string) {
	hostname, err := os.Hostname()
	if err != nil {
		tools.PrintColorf(tools.RED, "Cannot get hostname! %s\n", err)
	}
	return
}

func getModel() (model string) {
	model, err := getFileString(MODEL_FILE)
	if err != nil {
		tools.PrintColorf(tools.RED, "Cannot get model. Are you sure you are running this on a Raspberry Pi? %s\n", err)
	}
	return
}

func getCelsius() (celsius string) {
	return fmt.Sprintf("%.1f\u00b0C", getCelsiusTemperature())
}

func getFarenheit() (farenheit string) {
	return fmt.Sprintf("%.1f\u00b0F", (getCelsiusTemperature()*1.8)+32)
}

func getCelsiusTemperature() (celsius float64) {
	rawtemperature, err := getFileString(CPU_TEMP_FILE)
	if err != nil {
		tools.PrintColorf(tools.RED, "Cannot get CPU temperature! Are you sure you are running this on a Raspberry Pi? %s\n", err)
		return
	}
	kcelsius, err := strconv.ParseFloat(rawtemperature, 64)
	if err != nil {
		tools.PrintColorf(tools.RED, "Cannot parse raw temperature from %q: %s\n", CPU_TEMP_FILE, rawtemperature)
		return
	}
	celsius = kcelsius / 1000
	return
}

func getNetworking() (nics []NIC) {
	nics, err := getNICs()
	if err != nil {
		panic(err)
	}
	return nics
}

func getNICs() (nics []NIC, err error) {
	cmd := exec.Command("ip", "-brief", "address")
	output, err := cmd.CombinedOutput()
	if err != nil {
		tools.PrintColorf(tools.RED, "ip command failed! %v\nOutput: %s\n", cmd.Args, string(output))
		return
	}
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		var nic NIC
		fields := strings.Fields(scanner.Text())
		nic.Name = fields[0]
		nic.State = fields[1]
		if len(fields) > 2 {
			nic.IpAddress = fields[2]
		}
		if nic.Name != "lo" {
			nics = append(nics, nic)
		}
	}
	return
}

func getFileString(f string) (s string, err error) {
	if _, err = os.Stat(f); os.IsNotExist(err) {
		return
	}
	content, err := ioutil.ReadFile(f)
	if err != nil {
		return
	}
	s = strings.TrimSpace(string(content))
	return
}

func usage() {
	fmt.Printf("rpi-info v.%s (c) Lyderic Landry, London 2019\n", VERSION)
	fmt.Println("Usage: rpi-info <flags>")
	flag.PrintDefaults()
}

func (information Information) String() string {
	var buffer strings.Builder
	buffer.WriteString(fmt.Sprintf("%-12.12s: %s\n", "Hostname", information.Hostname))
	buffer.WriteString(fmt.Sprintf("%-12.12s: %s\n", "Model", information.Model))
	buffer.WriteString(fmt.Sprintf("%-12.12s: %s %s\n", "Temperature", information.Celsius, information.Farenheit))
	buffer.WriteString(fmt.Sprintf("%-12.12s:\n", "Networking"))
	buffer.WriteString(displayNetworking(information.Networking))
	return buffer.String()
}

func displayNetworking(networking []NIC) (output string) {
	var buffer strings.Builder
	buffer.WriteString(fmt.Sprintf(" %-8.8s %-8.8s %s\n", "Name", "State", "IP Address"))
	buffer.WriteString(fmt.Sprintf(" %-8.8s %-8.8s %s\n", "----", "-----", "----------"))
	for idx, nic := range networking {
		buffer.WriteString(nic.String())
		if idx != len(networking)-1 {
			buffer.WriteString("\n")
		}
	}
	return buffer.String()
}

func (nic NIC) String() string {
	return fmt.Sprintf(" %-8.8s %-8.8s %s", nic.Name, nic.State, nic.IpAddress)
}
