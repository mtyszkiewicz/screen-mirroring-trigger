package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var lastTriggerTime time.Time

func main() {
	cooldown := flag.Int("cooldown", 30, "Cooldown period in seconds between triggers")
	command := flag.String("command", "", "Command to execute when screen mirroring connection detected (required)")
	flag.Parse()

	if *command == "" {
		fmt.Println("Error: -command is required")
		flag.Usage()
		os.Exit(1)
	}

	handle, err := pcap.OpenLive("awdl0", 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	err = handle.SetBPFFilter("dst port 7000")
	if err != nil {
		log.Fatal(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	
	log.Printf("Monitoring screen mirroring connections on awdl0 (cooldown: %ds)...\n", *cooldown)
	log.Printf("Will execute: %s\n", *command)

	for range packetSource.Packets() {
		if time.Since(lastTriggerTime) >= time.Duration(*cooldown)*time.Second {
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fmt.Printf("[%s] Screen mirroring connection detected, executing command...\n", timestamp)

			// Parse command and args
			parts := strings.Fields(*command)
			cmd := exec.Command(parts[0], parts[1:]...)
			
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("[%s] Command failed: %v\n", timestamp, err)
				if len(output) > 0 {
					fmt.Printf("[%s] Output: %s\n", timestamp, string(output))
				}
			} else {
				fmt.Printf("[%s] Command executed successfully\n", timestamp)
				if len(output) > 0 {
					fmt.Printf("[%s] Output: %s\n", timestamp, string(output))
				}
			}

			lastTriggerTime = time.Now()
		}
	}
}
