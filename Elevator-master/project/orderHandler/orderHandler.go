package orderHandler
import (
	"../elevio"
	"strconv"
	."../typeDefinitions"
	"fmt"
)

func OrderHandler(orderToNetworkCh chan<- Orders, 
	orderFromNetworkCh <-chan Orders,
	ordersToFSMCh chan<- ElevInfo, 
	ordersFromFSMCh <-chan ElevInfo, 
	peersAliveCh <-chan PeerUpdate, 
	setLightCh <-chan [3] int, 
	ourId string){
	
	myId, _ := strconv.Atoi(ourId)
	var orderInfo Orders
	var peersAlive []int

	buttonEventCh := make(chan elevio.ButtonEvent)

	go elevio.PollButtons(buttonEventCh)

	for{
		select{
		case BT_Press := <-buttonEventCh:
			newOrderInfo := orderInfo
			send := false
			if((BT_Press.Button == elevio.BT_Cab) && !orderInfo.ElevInfos[myId].Orders[BT_Press.Floor][BT_Press.Button]){
				send = true
				newOrderInfo.ElevInfos[myId].Orders[BT_Press.Floor][BT_Press.Button] = true
				newOrderInfo.ElevToExecute = myId
			} else if (BT_Press.Button == elevio.BT_HallUp || BT_Press.Button == elevio.BT_HallDown){
				for _, k := range peersAlive{
					if !orderInfo.ElevInfos[k].Orders[BT_Press.Floor][BT_Press.Button]{
						send = true
					}
				}
				if send {
					bestElev := distributeTo(orderInfo, peersAlive, BT_Press.Floor)
					newOrderInfo.ElevInfos[bestElev].Orders[BT_Press.Floor][BT_Press.Button] = true
					newOrderInfo.ElevToExecute = bestElev
				}
			}
			if send{
				newOrderInfo.Floor = BT_Press.Floor
				newOrderInfo.Button = fromButtonTypeToInt(BT_Press.Button)
				newOrderInfo.NewOrder = true
				orderToNetworkCh <- newOrderInfo
			}

		case newOrderInfo := <-orderFromNetworkCh:
			orderInfo = updateOrders(newOrderInfo,orderInfo)
			for _, k := range peersAlive{
				if k != myId {
					orderInfo.ElevInfos[k].State = newOrderInfo.ElevInfos[k].State
					orderInfo.ElevInfos[k].Direction = newOrderInfo.ElevInfos[k].Direction
					orderInfo.ElevInfos[k].LastFloor = newOrderInfo.ElevInfos[k].LastFloor
				}
			}
			updateLightsOff(orderInfo,peersAlive,myId)
			ordersToFSMCh <- orderInfo.ElevInfos[myId]

		case orderComp := <-ordersFromFSMCh:
				orderInfo.ElevInfos[myId] = orderComp
				if orderInfo.ElevInfos[myId].State == Stuck {
					orderInfo = redistributeOrders(orderInfo, peersAlive, myId)
					orderInfo = removeHallCalls(orderInfo, myId)
					orderInfo.NewOrder = true
					orderInfo.Floor = 0
					orderInfo.Button = 1
				} else {
					updateLightsOff(orderInfo, peersAlive, myId)
					orderInfo.Floor = orderComp.LastFloor
					orderInfo.Button = directionToButton(orderComp.Direction)
					orderInfo.ElevToExecute = myId
					orderInfo.NewOrder = false					
				}
				orderToNetworkCh <- orderInfo
				
		case p := <-peersAliveCh:
			fmt.Println(p.Peers[0])
			newPeersAlive := make([]int, len(p.Peers))
			for i, k := range p.Peers {
				id, _ := strconv.Atoi(k)
				newPeersAlive[i] = id
			}
			if((len(peersAlive)-1) == len(newPeersAlive)) {
				redistribute := true
				if myId != newPeersAlive[0] {
					redistribute = false
				}
				if(redistribute){
					lostElevator, _ := strconv.Atoi(p.Lost[0])
					newOrders := redistributeOrders(orderInfo, newPeersAlive, lostElevator)
					newOrders = removeHallCalls(orderInfo, lostElevator)
					newOrders.NewOrder = false
					newOrders.Floor = 0
					newOrders.Button = 1
					newOrders.ElevToExecute = myId
					orderToNetworkCh <- newOrders
				}
			}
			if(len(newPeersAlive) == 1){
				orderInfo = removeHallCalls(orderInfo, myId)
				updateLightsOff(orderInfo, newPeersAlive, myId)
				ordersToFSMCh <- orderInfo.ElevInfos[myId]
			}
			peersAlive = newPeersAlive
			
		case newLight := <- setLightCh:
			if !((newLight[1] == fromButtonTypeToInt(elevio.BT_Cab)) && (newLight[2] != myId)){
				elevio.SetButtonLamp(fromIntToButtonType(newLight[1]), newLight[0], true)
			}

		}
	}
}
