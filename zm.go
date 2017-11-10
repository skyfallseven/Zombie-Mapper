/*
Zombie Mapper Version 0.001 Alpha
Command & Control System for IoT Devices
Keep your devices secure w/ a central management system
--------------------------------
Features:
libpcap sniffing of packets (dev)
	-Track unique IP addresses that contact device
Analysis of packets and real-time alerts (dev)
Reporting to AbuseIPDB (todo)
-----------
Auto generation of keys after breaches (todo)
Cron jobs running daily (todo)
PGP Authentication with server (todo)
Web interface for server (todo)
MySQL database of all devices (todo)
*/

package main

import (
	"fmt"
	"log"
	"time"
	"net"
	"bufio"
	"strings"
	"os"
	"github.com/google/go-cmp/cmp"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/layers"
)

var isSniffing bool = false	//This is checked by the sniffing loop

var IPlist []net.IP //Keep a list of IPs that have contacted us

var UniqueIPs int //Counter

var OurIP net.IP //localhost IP

//TODO: Get the interface from the command line
var (
	device string = "wlan0"
	snapshot_len int32 = 1024
	promiscuous bool = false
	err error
	timeout time.Duration = 30* time.Second
	handle *pcap.Handle
)


type Hit struct {
	//What else are we caring about?
	sourceIP net.IP
	destPort layers.TCPPort
	//numHits int
}

var HitList []Hit //Track all 'hits' during sniff
//TODO: Check memory usage and find a way to optimize
// 		Or offload this data. Don't want to fill the device

//Get local IP from device
//Used to filter out packets not meant for monitoring
//TODO: Have a configuration system on the server
func getLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

//This is the mechanism the C2 uses to stop sniffing
func setSniffer(x bool) {
	if x {
		isSniffing = true
	} else {
		isSniffing = false
	}
}

//Filters out packets not directed at us
func pktForDevice(packet gopacket.Packet) bool {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		if cmp.Equal(ip.DstIP, OurIP)  {
			return true
		}
	}
	return false
}

//Returns the source IP address of the device sending data to us
func ipTracker(packet gopacket.Packet) net.IP {
	var newIP net.IP
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		newIP = ip.SrcIP
	}
	return newIP
}

//Our own compare function to see if IP is already tracked
//TODO, if yes, add to counter for that IP (if structs are added)
func ipExist(s []net.IP, x net.IP) bool {
	for i, _ := range s {
		if cmp.Equal(s[i], x) {
			return true
		}
	}
	return false
}

//Detecting application layers if there are any
func portDetect(packet gopacket.Packet) layers.TCPPort {
	var dstPort layers.TCPPort

	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		dstPort = tcp.DstPort
	}
	return dstPort
}

//Run the main sniffing program
func sniff() []Hit {
	fmt.Println("Sniffing started!")
	//Open Device
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {log.Fatal(err)}
	defer handle.Close()

	HitList = make([]Hit, 0, 1) //initialize the list

	//Use handle as a packet source to process all packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		if pktForDevice(packet) { //Only packets that are for us
			newHit := new(Hit) //init new object
			sIP := ipTracker(packet) //Get source IP of this pkt

			if !(ipExist(IPlist, sIP)) {
				IPlist = append(IPlist, sIP) //add it
				UniqueIPs++ //increment counter

				newHit.sourceIP = sIP // set source
				//newHit.numHits++ //inc # of encounters
				newHit.destPort = portDetect(packet)

			}
			HitList = append(HitList, *newHit)
		}

		if !isSniffing {
			break
		}

	}
	return HitList
}

func main() {
	fmt.Println("Zombie Mapper v0.001-DEV")
	OurIP = getLocalIP()
	fmt.Println("Starting on", OurIP)
	UniqueIPs = 0
	IPlist = make([]net.IP, 0, 1)

	ln, _ := net.Listen("tcp", "129.21.146.90:8081") //Listen on all interfaces
	fmt.Println("Listening...")
	conn, _ := ln.Accept() //Accept connections
	fmt.Println("Connected!")

	for {
		msg, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println("Message Received:", msg)
		cmd := strings.ToLower(msg)

		switch cmd {
			case "start":
				setSniffer(true)
				sniff()
			case "stop":
				setSniffer(false)
			case "exit":
				os.Exit(1)
				
		}
	}

}