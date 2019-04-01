package peers

import (
	"../conn"
	"fmt"
	"net"
	"sort"
	"time"
	."../../typeDefinitions"
)

const interval = 15 * time.Millisecond
const timeout = 100 * time.Millisecond

func Transmitter(port int, id string) {

	conn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))

	enable := true
	for {
		select {
		case <-time.After(interval):
		}
		if enable {
			conn.WriteTo([]byte(id), addr)
		}
	}
}

func Receiver(port int, peerUpdateCh chan<- PeerUpdate, MyId string) {

	var buf [1024]byte
	var p PeerUpdate
	lastSeen := make(map[string]time.Time)

	conn := conn.DialBroadcastUDP(port)

	for {
		updated := false

		conn.SetReadDeadline(time.Now().Add(interval))
		n, _, _ := conn.ReadFrom(buf[0:])

		id := string(buf[:n])

		// Adding new connection
		p.New = ""
		if id != "" {
			if _, idExists := lastSeen[id]; !idExists {
				p.New = id
				updated = true
			}

			lastSeen[id] = time.Now()
		}

		// Removing dead connection
		p.Lost = make([]string, 0)
		for k, v := range lastSeen {
			if time.Now().Sub(v) > timeout {
				updated = true
				p.Lost = append(p.Lost, k)
				delete(lastSeen, k)
			}
		}

		// Sending update
		if updated {
			p.Peers = make([]string, 0, len(lastSeen))

			for k, _ := range lastSeen {
				p.Peers = append(p.Peers, k)
			}
			AddYourself := true
			for _, k := range p.Peers{
				if(k == MyId){
					AddYourself = false
				}
			}
			if(AddYourself){
				p.Peers = append(p.Peers,MyId)
			}
			sort.Strings(p.Peers)
			sort.Strings(p.Lost)
			if p.New != MyId {
				peerUpdateCh <- p				
			}

		}
	}
}
