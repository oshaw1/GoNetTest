package networkTesting

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type JitterTestResult struct {
	Timestamp   time.Time       `json:"timestamp"`
	Target      string          `json:"target"`
	PacketCount int             `json:"packet_count"`
	AvgJitter   time.Duration   `json:"avg_jitter"`
	MaxJitter   time.Duration   `json:"max_jitter"`
	MinJitter   time.Duration   `json:"min_jitter"`
	PacketLoss  float64         `json:"packet_loss"`
	RTTs        []time.Duration `json:"rtts"`
	Status      string          `json:"status"`
	Error       error           `json:"error,omitempty"`
}

func (t *NetworkTester) RunJitterTest() (*JitterTestResult, error) {
	dst, err := net.ResolveIPAddr("ip4", t.config.Tests.JitterTest.Target)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve target IP: %w", err)
	}

	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, fmt.Errorf("failed to create ICMP connection: %w", err)
	}
	defer conn.Close()

	result := &JitterTestResult{
		Timestamp:   time.Now(),
		Target:      t.config.Tests.JitterTest.Target,
		PacketCount: t.config.Tests.JitterTest.PacketCount,
		RTTs:        make([]time.Duration, 0),
	}

	var lastRTT time.Duration
	var totalJitter time.Duration
	var lostPackets int

	for i := 0; i < result.PacketCount; i++ {
		rtt, err := sendJitterPing(conn, dst, i, t.config.Tests.ICMP.TimeoutSeconds)
		if err != nil {
			lostPackets++
			continue
		}

		result.RTTs = append(result.RTTs, rtt)

		if i > 0 {
			jitter := abs(rtt - lastRTT)
			totalJitter += jitter

			if jitter > result.MaxJitter {
				result.MaxJitter = jitter
			}
			if jitter < result.MinJitter || result.MinJitter == 0 {
				result.MinJitter = jitter
			}
		}

		lastRTT = rtt
		time.Sleep(50 * time.Millisecond) // Space out packets
	}

	if len(result.RTTs) == 0 {
		result.Status = "FAILED"
		result.Error = fmt.Errorf("all packets lost")
		return result, result.Error
	}

	result.PacketLoss = float64(lostPackets) / float64(result.PacketCount) * 100
	if len(result.RTTs) > 1 {
		result.AvgJitter = totalJitter / time.Duration(len(result.RTTs)-1)
	}
	result.Status = "SUCCESS"

	return result, nil
}

func sendJitterPing(conn *icmp.PacketConn, dst *net.IPAddr, seq int, timeout int) (time.Duration, error) {
	start := time.Now()

	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  seq,
			Data: []byte("JITTER"),
		},
	}

	wb, err := wm.Marshal(nil)
	if err != nil {
		return 0, err
	}

	if _, err := conn.WriteTo(wb, dst); err != nil {
		return 0, err
	}

	rb := make([]byte, 1500)
	if err := conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second)); err != nil {
		return 0, err
	}

	n, _, err := conn.ReadFrom(rb)
	if err != nil {
		return 0, err
	}

	rm, err := icmp.ParseMessage(1, rb[:n])
	if err != nil {
		return 0, err
	}

	if rm.Type != ipv4.ICMPTypeEchoReply {
		return 0, fmt.Errorf("non-echo reply received")
	}

	return time.Since(start), nil
}

func abs(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}
