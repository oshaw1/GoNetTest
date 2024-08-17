package networkTesting

import (
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
			msg := createICMPMessage(tt.sequence)

			if msg.Type != tt.wantType {
				t.Errorf("Expected message type %v, got %v", tt.wantType, msg.Type)
			}

			if msg.Code != tt.wantCode {
				t.Errorf("Expected message code %d, got %d", tt.wantCode, msg.Code)
			}

			echo, ok := msg.Body.(*icmp.Echo)
			if !ok {
				t.Fatalf("Expected message body to be *icmp.Echo, got %T", msg.Body)
			}

			if echo.Seq != tt.wantSeq {
				t.Errorf("Expected sequence number %d, got %d", tt.wantSeq, echo.Seq)
			}
		})
	}
}

func TestUpdateTestResult(t *testing.T) {
	tests := []struct {
		name      string
		initial   ICMBTestResult
		icmpType  ipv4.ICMPType
		rtt       time.Duration
		wantFinal ICMBTestResult
	}{
		{
			name:      "Successful ping",
			initial:   ICMBTestResult{Sent: 1},
			icmpType:  ipv4.ICMPTypeEchoReply,
			rtt:       100 * time.Millisecond,
			wantFinal: ICMBTestResult{Sent: 1, Received: 1, Lost: 0, MinRTT: 100 * time.Millisecond, MaxRTT: 100 * time.Millisecond},
		},
		{
			name:      "Lost ping",
			initial:   ICMBTestResult{Sent: 1, Received: 1},
			icmpType:  ipv4.ICMPTypeDestinationUnreachable,
			rtt:       0,
			wantFinal: ICMBTestResult{Sent: 1, Received: 1, Lost: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial
			updateTestResult(&result, tt.icmpType, tt.rtt)

			if result.Received != tt.wantFinal.Received {
				t.Errorf("Expected Received to be %d, got %d", tt.wantFinal.Received, result.Received)
			}
			if result.Lost != tt.wantFinal.Lost {
				t.Errorf("Expected Lost to be %d, got %d", tt.wantFinal.Lost, result.Lost)
			}
			if result.MinRTT != tt.wantFinal.MinRTT {
				t.Errorf("Expected MinRTT to be %v, got %v", tt.wantFinal.MinRTT, result.MinRTT)
			}
			if result.MaxRTT != tt.wantFinal.MaxRTT {
				t.Errorf("Expected MaxRTT to be %v, got %v", tt.wantFinal.MaxRTT, result.MaxRTT)
			}
		})
	}
}

func TestCalculateAverageRTT(t *testing.T) {
	tests := []struct {
		name      string
		initial   ICMBTestResult
		wantFinal ICMBTestResult
	}{
		{
			name:      "Calculate average with received packets",
			initial:   ICMBTestResult{Received: 3, AvgRTT: 300 * time.Millisecond},
			wantFinal: ICMBTestResult{Received: 3, AvgRTT: 100 * time.Millisecond},
		},
		{
			name:      "No received packets",
			initial:   ICMBTestResult{Received: 0, AvgRTT: 300 * time.Millisecond},
			wantFinal: ICMBTestResult{Received: 0, AvgRTT: 300 * time.Millisecond},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.initial
			calculateAverageRTT(&result)

			if result.AvgRTT != tt.wantFinal.AvgRTT {
				t.Errorf("Expected AvgRTT to be %v, got %v", tt.wantFinal.AvgRTT, result.AvgRTT)
			}
		})
	}
}
