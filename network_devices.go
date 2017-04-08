package main

import (
	"fmt"
	"log"
	"strconv"
	str "strings"
	"time"

	gosnmp "github.com/soniah/gosnmp"
)

type IPAddr string

// Port

type Port struct {
	RemoteAddress IPAddr
	RemotePort    int
	LocalAddress  IPAddr
	LocalPort     int
}

func (p Port) String() string {
	return fmt.Sprintf("%s:%02d\t->\t%s\n", p.LocalAddress, p.LocalPort, p.RemoteAddress)
}

// interface NetworkDevice

type NetworkDevice interface {
	GetIP() IPAddr
	GetNeighbors() map[int]IPAddr
}

// NetworkSwitch

type NetworkSwitch struct {
	IP    IPAddr
	Ports map[int]IPAddr
}

func NewNetworkSwitch(IP IPAddr) *NetworkSwitch {
	return &NetworkSwitch{
		IP:    IP,
		Ports: make(map[int]IPAddr),
	}
}

func (n NetworkSwitch) GetIP() IPAddr {
	return n.IP
}

func (n NetworkSwitch) String() string {
	return fmt.Sprintf("\nSwitch: %s (%d connections)\n", n.IP, len(n.Ports))
}

func (n *NetworkSwitch) GetNeighbors() map[int]IPAddr {
	snmpConnection := &gosnmp.GoSNMP{
		Target:    string(n.IP),
		Port:      161,
		Community: "DS-public",
		Version:   gosnmp.Version2c,
		Timeout:   time.Duration(2) * time.Second,
	}

	fmt.Printf("Trying %s...", n.GetIP())
	err := snmpConnection.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	} else {
		defer snmpConnection.Conn.Close()
	}

	lldpOid := "1.0.8802.1.1.2.1.4.2.1"

	err2 := snmpConnection.BulkWalk(lldpOid, n.interpretValue)
	if err2 != nil {
		fmt.Println(" error, ignoring...")
	} else {
		fmt.Println("")
	}
	// fmt.Printf("%s", n.Ports)

	return n.Ports
}

func (n NetworkSwitch) ShowPorts() string {
	return fmt.Sprintf("%s", n.Ports)
}

func (n *NetworkSwitch) interpretValue(pdu gosnmp.SnmpPDU) error {

	splitOid := str.Split(pdu.Name, ".")

	portId, err := strconv.Atoi(splitOid[13])
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	remoteIP := IPAddr(str.Join(splitOid[17:], "."))

	if _, ok := n.Ports[portId]; !ok {
		n.Ports[portId] = remoteIP
	}

	return nil
}
