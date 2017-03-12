package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strconv"
	str "strings"

	gosnmp "github.com/soniah/gosnmp"
)

const (
	ChassisIdSubtype = 4 + iota // 4
	ChassisId                   // 5
	PortIdSubtype               // 6
	PortId                      // 7
	PortDesc                    // 8
	SysName                     // 9
	SysDesc                     // 10
)

type Port struct {
	ChassisIdSubtype *big.Int // 4
	ChassisId        string   // 5
	PortIdSubtype    *big.Int // 6
	PortId           *big.Int // 7
	PortDesc         string   // 8
	SysName          string   // 9
	SysDesc          string   // 10
}

func (p Port) String() string {
	return fmt.Sprintf("Remote:\n Chassis: %s\n Port: %d (%s)\n System: %s\n", p.ChassisId, p.PortId, p.PortDesc, p.SysName)
}

func unpackPduName(s string) (int, int) {
	splitOid := str.Split(s, ".")
	functId, err := strconv.Atoi(splitOid[11])
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	portId, err := strconv.Atoi(splitOid[13])
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	return functId, portId
}

func interpretValue(pdu gosnmp.SnmpPDU) error {

	funcId, portId := unpackPduName(pdu.Name)

	_, ok := ports[portId]
	if !ok {
		ports[portId] = &Port{}
	}

	switch funcId {
	case ChassisIdSubtype:
		ports[portId].ChassisIdSubtype = gosnmp.ToBigInt(pdu.Value)
	case ChassisId:
		hexEncoded := hex.EncodeToString(pdu.Value.([]byte))
		ports[portId].ChassisId = hexEncoded
	case PortIdSubtype:
		ports[portId].PortIdSubtype = gosnmp.ToBigInt(pdu.Value)
	case PortId:
		ports[portId].PortId = gosnmp.ToBigInt(pdu.Value)
	case PortDesc:
		ports[portId].PortDesc = string(pdu.Value.([]byte))
	case SysName:
		ports[portId].SysName = string(pdu.Value.([]byte))
	case SysDesc:
		ports[portId].SysDesc = string(pdu.Value.([]byte))
	}

	return nil
}

var ports = make(map[int]*Port)

func main() {
	gosnmp.Default.Target = "10.1.1.1"
	gosnmp.Default.Community = "DS-public"
	err := gosnmp.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer gosnmp.Default.Conn.Close()

	lldpOid := "1.0.8802.1.1.2.1.4.1.1"

	err2 := gosnmp.Default.BulkWalk(lldpOid, interpretValue)
	if err2 != nil {
		log.Fatalf("Get() err: %v", err2)
	}

	for index, port := range ports {
		fmt.Printf("Port %d:\n%s", index, port)
	}
}
