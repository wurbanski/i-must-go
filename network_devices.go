package main

import (
	"fmt"
	"log"
	"strconv"
	str "strings"
	"time"

	gosnmp "github.com/soniah/gosnmp"
)

type Port struct {
	RemoteSwitch *NetworkSwitch
}

type NetworkSwitch struct {
	IP    string
	Ports map[int]*Port
}

func NewNetworkSwitch(IP string) *NetworkSwitch {
	return &NetworkSwitch{
		IP:    IP,
		Ports: make(map[int]*Port),
	}
}

func (n NetworkSwitch) FindRemotes() {
	snmpConnection := &gosnmp.GoSNMP{
		Target:    n.IP,
		Port:      161,
		Community: "DS-public",
		Version:   gosnmp.Version2c,
		Timeout:   time.Duration(2) * time.Second,
	}

	err := snmpConnection.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer snmpConnection.Conn.Close()

	lldpOid := "1.0.8802.1.1.2.1.4.2.1"

	err2 := snmpConnection.BulkWalk(lldpOid, n.interpretValue)
	if err2 != nil {
		log.Fatalf("Get() err: %v", err2)
	}

}

func (n NetworkSwitch) String() string {
	return fmt.Sprintf("Switch: %s\nPorts: %s", n.IP, n.Ports)
}

func (n NetworkSwitch) ShowPorts() string {
	return fmt.Sprintf("%s", n.Ports)
}

func (p Port) String() string {
	return fmt.Sprintf("Remote IP:%s\n", p.RemoteSwitch.IP)
}

func (n NetworkSwitch) unpackPduName(s string) (int, string) {
	splitOid := str.Split(s, ".")

	portId, err := strconv.Atoi(splitOid[13])
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	remoteIP := str.Join(splitOid[17:], ".")

	return portId, remoteIP
}

func (n *NetworkSwitch) interpretValue(pdu gosnmp.SnmpPDU) error {

	portId, remoteIP := n.unpackPduName(pdu.Name)

	nswitch, ok := switches[remoteIP]
	if !ok {
		nswitch = NewNetworkSwitch(remoteIP)
		switches[remoteIP] = nswitch
	}

	_, ok = n.Ports[portId]
	if !ok {
		n.Ports[portId] = &Port{nswitch}
	}

	return nil
}
