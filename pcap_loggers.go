/*
 *    pcap_loggers.go - HoneyBadger core library for detecting TCP attacks
 *    such as handshake-hijack, segment veto and sloppy injection.
 *
 *    Copyright (C) 2014  David Stainton
 *
 *    This program is free software: you can redistribute it and/or modify
 *    it under the terms of the GNU General Public License as published by
 *    the Free Software Foundation, either version 3 of the License, or
 *    (at your option) any later version.
 *
 *    This program is distributed in the hope that it will be useful,
 *    but WITHOUT ANY WARRANTY; without even the implied warranty of
 *    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *    GNU General Public License for more details.
 *
 *    You should have received a copy of the GNU General Public License
 *    along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package HoneyBadger

import (
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/layers"
	"code.google.com/p/gopacket/pcapgo"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type PcapLogger struct {
	Dir    string
	Flow   TcpIpFlow
	writer *pcapgo.Writer
	file   *os.File
}

func NewPcapLogger(dir string, flow TcpIpFlow) *PcapLogger {
	var err error
	p := PcapLogger{
		Flow: flow,
		Dir:  dir,
	}
	p.file, err = os.OpenFile(filepath.Join(p.Dir, fmt.Sprintf("%s.pcap", p.Flow.String())), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}
	p.writer = pcapgo.NewWriter(p.file)
	err = p.writer.WriteFileHeader(65536, layers.LinkTypeEthernet) // XXX
	if err != nil {
		panic(err)
	}
	return &p
}

func (p *PcapLogger) WritePacket(rawPacket []byte, timestamp time.Time) {
	err := p.writer.WritePacket(gopacket.CaptureInfo{
		Timestamp:     timestamp,
		CaptureLength: len(rawPacket),
		Length:        len(rawPacket),
	}, rawPacket)
	if err != nil {
		panic(err)
	}
}

func (p *PcapLogger) Close() {
	p.file.Close()
}
