package networkTesting

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type ICMPTestResult struct {
	Host      string
	Sent      int
	Received  int
	Lost      int
	MinRTT    time.Duration
	MaxRTT    time.Duration
	AvgRTT    time.Duration
	Timestamp time.Time
}

func TestNetwork(host string) (*ICMPTestResult, error) {
	count, err := getTestValue[int]("Count")
	if err != nil {
		return nil, err
	}
	dst, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve IP address: %w", err)
	}

	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, fmt.Errorf("failed to listen for ICMP packets: %w", err)
	}
	defer c.Close()

	result, err := performICMPTest(c, dst, count)
	if err != nil {
		return nil, fmt.Errorf("error performing icmp test")
	}

	return result, nil
}

func performICMPTest(c *icmp.PacketConn, dst *net.IPAddr, count int) (*ICMPTestResult, error) {
	result := &ICMPTestResult{
		Host:      dst.String(),
		Sent:      count,
		Timestamp: time.Now(),
	}

	responses := make(chan *icmpResponse, count)
	var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(seq int) {
			defer wg.Done()
			start := time.Now()
			if err := sendICMP(c, dst, seq); err != nil {
				responses <- &icmpResponse{err: err}
				return
			}
			rm, err := receiveICMP(c)
			rtt := time.Since(start)
			responses <- &icmpResponse{rm: rm, rtt: rtt, err: err}
		}(i)
	}

	go func() {
		wg.Wait()
		close(responses)
	}()

	for resp := range responses {
		if resp.err != nil {
			result.Lost++
		} else {
			updateTestResult(result, resp.rm.Type, resp.rtt)
		}
	}

	calculateAverageRTT(result)
	return result, nil
}

type icmpResponse struct {
	rm  *icmp.Message
	rtt time.Duration
	err error
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
	protocolICMP, err := getTestValue[int]("ProtocolICMP")
	if err != nil {
		return nil, err
	}

	timeoutSecond, err := getTestValue[int]("TimeoutSecond")
	if err != nil {
		return nil, err
	}

	timeout := time.Duration(timeoutSecond) * time.Second
	err = c.SetReadDeadline(time.Now().Add(timeout * time.Second))
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

func updateTestResult(result *ICMPTestResult, icmpType icmp.Type, rtt time.Duration) {
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

func calculateAverageRTT(result *ICMPTestResult) {
	if result.Received > 0 {
		result.AvgRTT /= time.Duration(result.Received)
	}
}
