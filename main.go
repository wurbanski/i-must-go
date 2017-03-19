package main

import "fmt"

var ports = make(map[int]*Port)
var switches = make(map[string]*NetworkSwitch)

func main() {

	startIP := "10.1.1.0"
	startSwitch := NewNetworkSwitch(startIP)
	switches[startIP] = startSwitch

	startSwitch.FindRemotes()

	for index, nswitch := range switches {
		fmt.Printf("* Switch %s: %s\n", index, nswitch.ShowPorts())
	}
}
