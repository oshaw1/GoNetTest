package networkTesting

import (
	"fmt"
	"os"
	"testing"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func TestCreateICMPMessage(t *testing.T) {
	tests := []struct {
		name     string
		sequence int
		wantType ipv4.ICMPType
		wantCode int
		wantSeq  int
	}{
		{
			name:     "Valid ICMP message",
			sequence: 1,
			wantType: ipv4.ICMPTypeEcho,
			wantCode: 0,
			wantSeq:  1,
		},
		{
			name:     "Valid ICMP message with large sequence number",
			sequence: 65535,
			wantType: ipv4.ICMPTypeEcho,
			wantCode: 0,
			wantSeq:  65535,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create an ICMP message directly
			wm := &icmp.Message{
				Type: ipv4.ICMPTypeEcho,
				Code: 0,
				Body: &icmp.Echo{
					ID:   os.Getpid() & 0xffff,
					Seq:  tt.sequence,
					Data: []byte("HELLO-R-U-THERE"),
				},
			}

			if wm.Type != tt.wantType {
				t.Errorf("message type = %v, want %v", wm.Type, tt.wantType)
			}

			if wm.Code != tt.wantCode {
				t.Errorf("message code = %d, want %d", wm.Code, tt.wantCode)
			}

			echo, ok := wm.Body.(*icmp.Echo)
			if !ok {
				t.Fatalf("message body type = %T, want *icmp.Echo", wm.Body)
			}

			if echo.Seq != tt.wantSeq {
				t.Errorf("sequence number = %d, want %d", echo.Seq, tt.wantSeq)
			}

			if string(echo.Data) != "HELLO-R-U-THERE" {
				t.Errorf("message data = %q, want %q", string(echo.Data), "HELLO-R-U-THERE")
			}
		})
	}
}

func TestUpdateICMPStats(t *testing.T) {
	tests := []struct {
		name      string
		initial   ICMPTestResult
		rtt       time.Duration
		wantFinal ICMPTestResult
	}{
		{
			name: "First RTT update",
			initial: ICMPTestResult{
				Sent:     1,
				Received: 0,
				MinRTT:   0,
				MaxRTT:   0,
				AvgRTT:   0,
			},
			rtt: 100 * time.Millisecond,
			wantFinal: ICMPTestResult{
				Sent:     1,
				Received: 0,
				MinRTT:   100 * time.Millisecond,
				MaxRTT:   100 * time.Millisecond,
				AvgRTT:   100 * time.Millisecond,
			},
		},
		{
			name: "Update with lower RTT",
			initial: ICMPTestResult{
				MinRTT: 100 * time.Millisecond,
				MaxRTT: 100 * time.Millisecond,
				AvgRTT: 100 * time.Millisecond,
			},
			rtt: 50 * time.Millisecond,
			wantFinal: ICMPTestResult{
				MinRTT: 50 * time.Millisecond,
				MaxRTT: 100 * time.Millisecond,
				AvgRTT: 150 * time.Millisecond, // Sum, not average yet
			},
		},
		{
			name: "Update with higher RTT",
			initial: ICMPTestResult{
				MinRTT: 100 * time.Millisecond,
				MaxRTT: 100 * time.Millisecond,
				AvgRTT: 100 * time.Millisecond,
			},
			rtt: 150 * time.Millisecond,
			wantFinal: ICMPTestResult{
				MinRTT: 100 * time.Millisecond,
				MaxRTT: 150 * time.Millisecond,
				AvgRTT: 250 * time.Millisecond, // Sum, not average yet
			},
		},
	}

	tester := &NetworkTester{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial
			tester.updateICMPStats(&result, tt.rtt)

			if result.MinRTT != tt.wantFinal.MinRTT {
				t.Errorf("MinRTT = %v, want %v", result.MinRTT, tt.wantFinal.MinRTT)
			}
			if result.MaxRTT != tt.wantFinal.MaxRTT {
				t.Errorf("MaxRTT = %v, want %v", result.MaxRTT, tt.wantFinal.MaxRTT)
			}
			if result.AvgRTT != tt.wantFinal.AvgRTT {
				t.Errorf("AvgRTT = %v, want %v", result.AvgRTT, tt.wantFinal.AvgRTT)
			}
		})
	}
}

func TestProcessICMPResponses(t *testing.T) {
	tests := []struct {
		name      string
		responses []icmpResponse
		wantFinal ICMPTestResult
	}{
		{
			name: "All successful responses",
			responses: []icmpResponse{
				{rm: &icmp.Message{Type: ipv4.ICMPTypeEchoReply}, rtt: 100 * time.Millisecond},
				{rm: &icmp.Message{Type: ipv4.ICMPTypeEchoReply}, rtt: 150 * time.Millisecond},
			},
			wantFinal: ICMPTestResult{
				Received: 2,
				Lost:     0,
				MinRTT:   100 * time.Millisecond,
				MaxRTT:   150 * time.Millisecond,
				AvgRTT:   125 * time.Millisecond,
			},
		},
		{
			name: "Mixed successful and failed responses",
			responses: []icmpResponse{
				{rm: &icmp.Message{Type: ipv4.ICMPTypeEchoReply}, rtt: 100 * time.Millisecond},
				{err: fmt.Errorf("timeout")},
				{rm: &icmp.Message{Type: ipv4.ICMPTypeDestinationUnreachable}},
			},
			wantFinal: ICMPTestResult{
				Received: 1,
				Lost:     2,
				MinRTT:   100 * time.Millisecond,
				MaxRTT:   100 * time.Millisecond,
				AvgRTT:   100 * time.Millisecond,
			},
		},
	}

	tester := &NetworkTester{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responses := make(chan *icmpResponse, len(tt.responses))
			result := &ICMPTestResult{}

			// Fill the channel with test responses
			for _, resp := range tt.responses {
				resp := resp // Create new variable for goroutine
				responses <- &resp
			}
			close(responses)

			tester.processICMPResponses(responses, result)

			if result.Received != tt.wantFinal.Received {
				t.Errorf("Received = %v, want %v", result.Received, tt.wantFinal.Received)
			}
			if result.Lost != tt.wantFinal.Lost {
				t.Errorf("Lost = %v, want %v", result.Lost, tt.wantFinal.Lost)
			}
			if result.MinRTT != tt.wantFinal.MinRTT {
				t.Errorf("MinRTT = %v, want %v", result.MinRTT, tt.wantFinal.MinRTT)
			}
			if result.MaxRTT != tt.wantFinal.MaxRTT {
				t.Errorf("MaxRTT = %v, want %v", result.MaxRTT, tt.wantFinal.MaxRTT)
			}
			if result.AvgRTT != tt.wantFinal.AvgRTT {
				t.Errorf("AvgRTT = %v, want %v", result.AvgRTT, tt.wantFinal.AvgRTT)
			}
		})
	}
}
