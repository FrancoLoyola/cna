package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// getOutboundIP :
//
// Get preferred outbound ip of this machine
func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println("Coudln't determine the default interface", err.Error())
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

// probeTCPPort
//
// Tries to connect to a given port for 1 second only
func probeTCPPort(ip, port string) {
	timeout := time.Duration(1 * time.Second)
	_, err := net.DialTimeout("tcp", ip+":"+port, timeout)
	if err != nil {
		return
	}
	fmt.Println("### Bingo: ", ip, port, "is listening! ###")
}

// Logic here
func main() {
	// Program flags
	mask := flag.String("netmask", "/24", "Mask for the IP, use the /8, /16 style notation")
	initPort := flag.Int("init-port", 1, "First port start the scan")
	endPort := flag.Int("end-port", 65535, "Last port to scan")
	flag.Parse()

	// Sanity checks
	if *initPort < 0 || *endPort < 0 {
		fmt.Println("Ports can't be negative...")
		os.Exit(1)
	}
	if *initPort > *endPort {
		fmt.Println("Initial port cannot be higher than the end...")
		os.Exit(1)
	}
	if *endPort > 65535 {
		fmt.Println("End port cannot be higher than 65535...")
		os.Exit(1)
	}

	// Find main interface IP and it's whole network based on the mask
	ip := getOutboundIP()
	ip, n, err := net.ParseCIDR(ip.String() + *mask)
	if err != nil {
		fmt.Println("Couldn't determine the net this host belongs to...", err.Error())
		os.Exit(1)
	}

	// Use as a base the network
	sIP := strings.Split(n.IP.String(), ".")
	// Not check for errors, should always be within 0-255 from the net package
	octet1, _ := strconv.Atoi(sIP[0])
	octet2, _ := strconv.Atoi(sIP[1])
	octet3, _ := strconv.Atoi(sIP[2])
	octet4, _ := strconv.Atoi(sIP[3])
	octet4++

	// Try to probe all ports on all hosts in the network instead of pinging, maybe is better to invert the loops host{port}
	// instead of port{host}
	fmt.Println("Going to probe all the IPs within this network", n.IP, *mask)
	fmt.Println("From port", *initPort, "to", *endPort)
	for i := *initPort; i <= *endPort; i++ {
		// Since the length of the network might vary (/16 vs /24) loop until there are no more IPs in the network
		for {
			// Craft the IP based on the octet and check that we are still within the network
			tmp := strconv.Itoa(octet1) + "." + strconv.Itoa(octet2) + "." + strconv.Itoa(octet3) + "." + strconv.Itoa(octet4)
			tmpIP := net.ParseIP(tmp)
			if !n.Contains(tmpIP) {
				// Reset over to the start
				octet1, _ = strconv.Atoi(sIP[0])
				octet2, _ = strconv.Atoi(sIP[1])
				octet3, _ = strconv.Atoi(sIP[2])
				octet4, _ = strconv.Atoi(sIP[3])
				octet4++
				break
			}
			go probeTCPPort(tmp, strconv.Itoa(i))
			// Increase / reset... Ugly, maybe bitwise operations instead of int/str?
			if octet4 == 255 {
				octet4 = 0
				octet3++
				continue
			}
			if octet3 == 255 {
				octet3 = 0
				octet2++
				continue
			}
			if octet2 == 255 {
				octet2 = 0
				octet1++
				continue
			}
			if octet1 > 255 {
				break
			}
			// Last thing to do is increase to the next address
			octet4++
		}
		if i%3000 == 0 {
			fmt.Println("Scanned up to port:", i)
		}
	}
	// Just to give time to possible messages coming back
	time.Sleep(1 * time.Second)
	fmt.Println("Done")
}
