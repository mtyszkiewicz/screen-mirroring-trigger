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

func main() {
	silence := flag.Int("silence", 30, "Seconds without packets before an active stream is considered ended")
	command := flag.String("command", "", "Command to execute when a screen mirroring stream starts (required)")
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

	log.Printf("Monitoring screen mirroring connections on awdl0 (silence: %ds)...\n", *silence)
	log.Printf("Will execute on stream start: %s\n", *command)

	// Edge-triggered state machine.
	//   IDLE   + packet  -> run command once, -> ACTIVE
	//   ACTIVE + packet  -> ignore (manual amp changes survive)
	//   ACTIVE + silence -> -> IDLE (no action)
	const (
		stateIdle   = "idle"
		stateActive = "active"
	)
	state := stateIdle
	var lastPacket time.Time

	silenceTimeout := time.Duration(*silence) * time.Second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	packets := packetSource.Packets()
	for {
		select {
		case packet, ok := <-packets:
			if !ok {
				log.Printf("Packet source closed; exiting.\n")
				return
			}
			_ = packet
			lastPacket = time.Now()

			if state == stateIdle {
				timestamp := time.Now().Format("2006-01-02 15:04:05")
				fmt.Printf("[%s] Screen mirroring stream started, executing command...\n", timestamp)
				runCommand(*command, timestamp)
				state = stateActive
			}
		case <-ticker.C:
			if state == stateActive && time.Since(lastPacket) >= silenceTimeout {
				timestamp := time.Now().Format("2006-01-02 15:04:05")
				fmt.Printf("[%s] No packets for %ds, stream considered ended.\n", timestamp, *silence)
				state = stateIdle
			}
		}
	}
}

func runCommand(command, timestamp string) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return
	}
	cmd := exec.Command(parts[0], parts[1:]...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("[%s] Command failed: %v\n", timestamp, err)
	} else {
		fmt.Printf("[%s] Command executed successfully\n", timestamp)
	}
	if len(output) > 0 {
		fmt.Printf("[%s] Output: %s\n", timestamp, string(output))
	}
}
