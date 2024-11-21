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
	Timestamp time.Time
	TestType  string
	Sent      int
	Received  int
	Lost      int
	MinRTT    time.Duration
	MaxRTT    time.Duration
	AvgRTT    time.Duration
}

type icmpResponse struct {
	rm  *icmp.Message
	rtt time.Duration
	err error
}

func (t *NetworkTester) runICMPTest() (*ICMPTestResult, error) {
	dst, err := net.ResolveIPAddr("ip4", "8.8.8.8")
	if err != nil {
		return nil, fmt.Errorf("failed to resolve IP address: %w", err)
	}

	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, fmt.Errorf("failed to listen for ICMP packets: %w", err)
	}
	defer c.Close()

	return t.performICMPTest(c, dst)
}

func (t *NetworkTester) performICMPTest(c *icmp.PacketConn, dst *net.IPAddr) (*ICMPTestResult, error) {
	count := t.config.Tests.ICMP.PacketCount
	result := &ICMPTestResult{
		Host:      dst.String(),
		Timestamp: time.Now(),
		TestType:  "icmp",
		Sent:      count,
	}

	responses := make(chan *icmpResponse, count)
	var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(seq int) {
			defer wg.Done()
			t.sendAndReceiveICMP(c, dst, seq, responses)
		}(i)
	}

	go func() {
		wg.Wait()
		close(responses)
	}()

	t.processICMPResponses(responses, result)
	return result, nil
}

func (t *NetworkTester) sendAndReceiveICMP(c *icmp.PacketConn, dst *net.IPAddr, seq int, responses chan<- *icmpResponse) {
	start := time.Now()
	if err := t.sendICMPPacket(c, dst, seq); err != nil {
		responses <- &icmpResponse{err: err}
		return
	}

	rm, err := t.receiveICMPPacket(c)
	responses <- &icmpResponse{
		rm:  rm,
		rtt: time.Since(start),
		err: err,
	}
}

func (t *NetworkTester) sendICMPPacket(c *icmp.PacketConn, dst *net.IPAddr, sequence int) error {
	wm := &icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  sequence,
			Data: []byte("HELLO-R-U-THERE"),
		},
	}

	wb, err := wm.Marshal(nil)
	if err != nil {
		return fmt.Errorf("failed to marshal ICMP message: %w", err)
	}

	if _, err := c.WriteTo(wb, dst); err != nil {
		return fmt.Errorf("failed to send ICMP message: %w", err)
	}

	return nil
}

func (t *NetworkTester) receiveICMPPacket(c *icmp.PacketConn) (*icmp.Message, error) {
	timeout := time.Duration(t.config.Tests.ICMP.TimeoutSeconds) * time.Second
	if err := c.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return nil, fmt.Errorf("failed to set read deadline: %w", err)
	}

	rb := make([]byte, 150)
	n, _, err := c.ReadFrom(rb)
	if err != nil {
		return nil, err
	}

	rm, err := icmp.ParseMessage(1, rb[:n])
	if err != nil {
		return nil, fmt.Errorf("failed to parse ICMP message: %w", err)
	}

	return rm, nil
}

func (t *NetworkTester) processICMPResponses(responses <-chan *icmpResponse, result *ICMPTestResult) {
	for resp := range responses {
		if resp.err != nil {
			result.Lost++
			continue
		}

		if resp.rm.Type == ipv4.ICMPTypeEchoReply {
			result.Received++
			t.updateICMPStats(result, resp.rtt)
		} else {
			result.Lost++
		}
	}

	if result.Received > 0 {
		result.AvgRTT /= time.Duration(result.Received)
	}
}

func (t *NetworkTester) updateICMPStats(result *ICMPTestResult, rtt time.Duration) {
	if rtt < result.MinRTT || result.MinRTT == 0 {
		result.MinRTT = rtt
	}
	if rtt > result.MaxRTT {
		result.MaxRTT = rtt
	}
	result.AvgRTT += rtt
}
