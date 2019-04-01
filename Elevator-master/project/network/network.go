package network

import (
	"./bcast"
	"fmt"
	"time"
	"strconv"
	//"sort"
	."../typeDefinitions"
	"./peers"
)

func Network(
	orderFromOrderHandlerCh <-chan Orders,
	orderToOrderHandlerCh chan<- Orders,
	sendPeerCh chan<- PeerUpdate,
	setLightCh chan<- [3]int,
	id string) {

	peerUpdateCh := make(chan PeerUpdate)
	orderTxCh := make(chan Orders)
	orderRxCh := make(chan Orders)
	ackOrderTxCh := make(chan Acknowledgement)
	ackOrderRxCh := make(chan Acknowledgement)
	lightOrderTxCh := make(chan [3]int)
	lightOrderRxCh := make(chan [3]int)
	
	go peers.Transmitter(15000, id) 
	go peers.Receiver(15000, peerUpdateCh, id)
	go bcast.Transmitter(24000, orderTxCh, ackOrderTxCh, lightOrderTxCh)
	go bcast.Receiver(24000, orderRxCh, ackOrderRxCh, lightOrderRxCh)

	//disconnected := true
	myId, _ := strconv.Atoi(id) 
	NumberOfPrevPrevPeers := 0
	var peerInfo PeerUpdate
	var lastOrderReceived Orders 
	var unconfirmedOrdersMatrix [N_FLOORS][3] UnconfirmedOrder
	// unconfirmedOrdersMatrix is a matrix that keeps track of which orders that currently are being sendt
	// Each element in the matrix corresponds to a certain button press (floor, buttontype)

	for i := 0; i < N_FLOORS; i++ {
		for j := 0; j < N_BUTTONS; j++ {
			unconfirmedOrdersMatrix[i][j].Active = false
		}
	}

	ticker := time.NewTicker(time.Millisecond * 200)

	fmt.Println("Started")

	for {
		select {
		case orderToShare := <-orderFromOrderHandlerCh:
			if (orderToShare.ElevToExecute == myId) && (unconfirmedOrdersMatrix[orderToShare.Floor][orderToShare.Button].Active == false) { 
				unconfirmedOrdersMatrix[orderToShare.Floor][orderToShare.Button] = activateUnconfirmedOrder(orderToShare, len(peerInfo.Peers))
			}
			orderTxCh <- orderToShare
		case receivedOrder := <- orderRxCh:
			lastOrderReceived = receivedOrder
			lastOrderReceived.ElevToExecute = myId
			if receivedOrder.ElevToExecute != myId {
				orderToOrderHandlerCh <- receivedOrder
			} else if unconfirmedOrdersMatrix[receivedOrder.Floor][receivedOrder.Button].Active == false {
				unconfirmedOrdersMatrix[receivedOrder.Floor][receivedOrder.Button] = activateUnconfirmedOrder(receivedOrder, len(peerInfo.Peers))
			}
			ack := createAcknowledgement(receivedOrder, myId) 
			ackOrderTxCh <- ack
		case ack := <- ackOrderRxCh:
			if unconfirmedOrdersMatrix[ack.Floor][ack.Button].Active == true {
				index, indexError := correspondingIndex(peerInfo.Peers, ack.From)
				if (indexError == "nil") && (index < len(unconfirmedOrdersMatrix[ack.Floor][ack.Button].ReceivedAcks)) {
					unconfirmedOrdersMatrix[ack.Floor][ack.Button].ReceivedAcks[index] = 1
				} else {
					fmt.Println(indexError)
				}
				// We check if we have received acknowledgements from all the other elevators
				fmt.Println(unconfirmedOrdersMatrix[ack.Floor][ack.Button].ReceivedAcks)
				if sum(unconfirmedOrdersMatrix[ack.Floor][ack.Button].ReceivedAcks) == (len(peerInfo.Peers)) {
					unconfirmedOrdersMatrix[ack.Floor][ack.Button].Active = false
					orderToOrderHandlerCh <- unconfirmedOrdersMatrix[ack.Floor][ack.Button].OrderToShare
					lightOrder := [3]int{ack.Floor, ack.Button, ack.ElevToExecute}
					if (unconfirmedOrdersMatrix[ack.Floor][ack.Button].OrderToShare.NewOrder == true) {
						lightOrderTxCh <- lightOrder
						setLightCh <- lightOrder
					}
				}
			}
		case receiveLightOrder := <- lightOrderRxCh:
			setLightCh <- receiveLightOrder
		case latestPeerInfo := <-peerUpdateCh:

			// The old orders are shared with the newest elevators
			if (NumberOfPrevPrevPeers >= 3) {
				unconfirmedOrdersMatrix[lastOrderReceived.Floor][lastOrderReceived.Button] = activateUnconfirmedOrder(lastOrderReceived, len(latestPeerInfo.Peers))
			}

			NumberOfPrevPrevPeers = len(peerInfo.Peers)
			peerInfo = latestPeerInfo
			sendPeerCh <- peerInfo

			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", latestPeerInfo.Peers)
			fmt.Printf("  New:      %q\n", latestPeerInfo.New)
			fmt.Printf("  Lost:     %q\n", latestPeerInfo.Lost)

		case <-ticker.C: 
			for i := 0; i < N_FLOORS; i++ {
				for j := 0; j < N_BUTTONS; j++ {
					if unconfirmedOrdersMatrix[i][j].Active == true {
						unconfirmedOrdersMatrix[i][j].Counter++;
						if unconfirmedOrdersMatrix[i][j].Counter > 10 {
							unconfirmedOrdersMatrix[i][j].Active = false
						} else {
							orderTxCh <- unconfirmedOrdersMatrix[i][j].OrderToShare
						}
					}
				}
			}
		default	:
		}
	}
}
