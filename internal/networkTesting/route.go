package networkTesting

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type RouteHop struct {
	Number  int           `json:"hop_number"`
	Address string        `json:"address"`
	RTT     time.Duration `json:"rtt"`
	Lost    bool          `json:"packet_lost"`
}

type RouteTestResult struct {
	Timestamp time.Time  `json:"timestamp"`
	Target    string     `json:"target"`
	Hops      []RouteHop `json:"hops"`
	Status    string     `json:"status"`
	Error     error      `json:"error,omitempty"`
}

func (t *NetworkTester) RunRouteTest() (*RouteTestResult, error) {
	dst, err := net.ResolveIPAddr("ip4", t.config.Tests.RouteTest.Target)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve target IP: %w", err)
	}

	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, fmt.Errorf("failed to create ICMP connection: %w", err)
	}
	defer conn.Close()

	result := &RouteTestResult{
		Timestamp: time.Now(),
		Target:    t.config.Tests.RouteTest.Target,
		Hops:      make([]RouteHop, 0),
	}

	maxHops := t.config.Tests.RouteTest.MaxHops
	for ttl := 1; ttl <= maxHops; ttl++ {
		hop := probeRouteHop(conn, dst, ttl, t.config.Tests.ICMP.TimeoutSeconds)
		result.Hops = append(result.Hops, hop)

		if !hop.Lost && hop.Address == dst.String() {
			result.Status = "SUCCESS"
			return result, nil
		}
	}

	result.Status = "INCOMPLETE"
	return result, nil
}

func probeRouteHop(conn *icmp.PacketConn, dst *net.IPAddr, ttl int, timeout int) RouteHop {
	hop := RouteHop{
		Number: ttl,
		Lost:   true,
	}

	// Debug logging
	fmt.Printf("Probing hop %d\n", ttl)

	if err := conn.IPv4PacketConn().SetTTL(ttl); err != nil {
		fmt.Printf("SetTTL error: %v\n", err)
		return hop
	}

	start := time.Now()

	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  ttl,
			Data: []byte("TRACEROUTE"),
		},
	}

	wb, err := wm.Marshal(nil)
	if err != nil {
		return hop
	}

	if _, err := conn.WriteTo(wb, dst); err != nil {
		return hop
	}

	rb := make([]byte, 1500)
	if err := conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second)); err != nil {
		return hop
	}

	n, peer, err := conn.ReadFrom(rb)
	if err != nil {
		return hop
	}

	rm, err := icmp.ParseMessage(1, rb[:n])
	if err != nil {
		return hop
	}

	switch rm.Type {
	case ipv4.ICMPTypeTimeExceeded:
		hop.Lost = false
		hop.Address = peer.String()
		hop.RTT = time.Since(start)
	case ipv4.ICMPTypeEchoReply:
		hop.Lost = false
		hop.Address = dst.String()
		hop.RTT = time.Since(start)
	}

	return hop
}
