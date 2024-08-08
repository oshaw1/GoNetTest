package networkTesting

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/oshaw1/go-net-test/config"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type TestResult struct {
	Host     string
	Port     int
	Sent     int
	Received int
	Lost     int
	MinRTT   time.Duration
	MaxRTT   time.Duration
	AvgRTT   time.Duration
}

func TestNetwork(host string, port int) (*TestResult, error) {
	count, _, _ := getTestValues()
	dst, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve IP address: %w", err)
	}

	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, fmt.Errorf("failed to listen for ICMP packets: %w", err)
	}
	defer c.Close()

	result := &TestResult{
		Host: host,
		Port: port,
		Sent: count,
	}

	for i := 0; i < count; i++ {
		start := time.Now()

		if err := sendICMP(c, dst, i); err != nil {
			return nil, err
		}

		rm, err := receiveICMP(c)
		rtt := time.Since(start)

		if err != nil {
			result.Lost++
		} else {
			updateTestResult(result, rm.Type, rtt)
		}
	}

	calculateAverageRTT(result)
	return result, nil
}

func getTestValues() (count int, protocolIMCP int, timeoutSecond time.Duration) {
	config, err := config.NewConfig("config/config.json")
	if err != nil {
		fmt.Printf("failed to read config file defaulting to small values. caused by : %v", err)
		return 5, 1, 1
	}
	return config.Count, config.ProtocolIMCP, time.Duration(config.TimeoutSecond)
}

func sendICMP(c *icmp.PacketConn, dst *net.IPAddr, sequence int) error {
	wm := createICMPMessage(sequence)
	wb, err := wm.Marshal(nil)
	if err != nil {
		return fmt.Errorf("failed to marshal ICMP message: %w", err)
	}

	if _, err := c.WriteTo(wb, dst); err != nil {
		return fmt.Errorf("failed to send ICMP message: %w", err)
	}

	return nil
}

func receiveICMP(c *icmp.PacketConn) (*icmp.Message, error) {
	_, protocolICMP, timeoutSecond := getTestValues()
	err := c.SetReadDeadline(time.Now().Add(timeoutSecond * time.Second))
	if err != nil {
		return nil, fmt.Errorf("failed to set read deadline: %w", err)
	}

	rb := make([]byte, 1500)
	n, _, err := c.ReadFrom(rb)
	if err != nil {
		return nil, err
	}

	rm, err := icmp.ParseMessage(protocolICMP, rb[:n])
	if err != nil {
		return nil, fmt.Errorf("failed to parse ICMP message: %w", err)
	}

	return rm, nil
}

func createICMPMessage(sequence int) icmp.Message {
	return icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  sequence,
			Data: []byte("HELLO-R-U-THERE"),
		},
	}
}

func updateTestResult(result *TestResult, icmpType icmp.Type, rtt time.Duration) {
	if icmpType == ipv4.ICMPTypeEchoReply {
		result.Received++
		if rtt < result.MinRTT || result.MinRTT == 0 {
			result.MinRTT = rtt
		}
		if rtt > result.MaxRTT {
			result.MaxRTT = rtt
		}
		result.AvgRTT += rtt
	} else {
		result.Lost++
	}
}

func calculateAverageRTT(result *TestResult) {
	if result.Received > 0 {
		result.AvgRTT /= time.Duration(result.Received)
	}
}
