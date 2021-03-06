package HoneyBadger

import (
	"github.com/david415/HoneyBadger/types"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
	"net"
	"testing"
	"time"
)

func TestOrderedCoalesceUsedPages(t *testing.T) {
	maxBufferedPagesTotal := 1024
	maxBufferedPagesPerFlow := 1024
	streamRing := types.NewRing(40)
	pager := NewPager()

	ipFlow, _ := gopacket.FlowFromEndpoints(layers.NewIPEndpoint(net.IPv4(1, 2, 3, 4)), layers.NewIPEndpoint(net.IPv4(2, 3, 4, 5)))
	tcpFlow, _ := gopacket.FlowFromEndpoints(layers.NewTCPPortEndpoint(layers.TCPPort(1)), layers.NewTCPPortEndpoint(layers.TCPPort(2)))
	flow := types.NewTcpIpFlowFromFlows(ipFlow, tcpFlow)

	var nextSeq types.Sequence = types.Sequence(1)

	dummyClose := func() {
		log.Print("dummyClose\n")
	}
	coalesce := NewOrderedCoalesce(dummyClose, nil, flow, pager, streamRing, maxBufferedPagesTotal, maxBufferedPagesPerFlow, false)

	ip := layers.IPv4{
		SrcIP:    net.IP{1, 2, 3, 4},
		DstIP:    net.IP{2, 3, 4, 5},
		Version:  4,
		TTL:      64,
		Protocol: layers.IPProtocolTCP,
	}
	tcp := layers.TCP{
		Seq:     3,
		SYN:     false,
		SrcPort: 1,
		DstPort: 2,
	}
	p := PacketManifest{
		Timestamp: time.Now(),
		Flow:      flow,
		IP:        ip,
		TCP:       tcp,
		Payload:   []byte{1, 2, 3, 4, 5, 6, 7},
	}

	coalesce.pager.Start()
	coalesce.insert(p, nextSeq)

	if coalesce.pager.Used() != 1 {
		t.Errorf("coalesce.pager.Used() not equal to 1\n")
		t.Fail()
	}

	coalesce.pager.Stop()
}
