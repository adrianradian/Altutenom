package main
import ("./elevio"
		"./network"
		"./fsm"
		"./orderHandler"
		"flag"
		."./typeDefinitions"
	)

func main(){
	var id string
	var port string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.StringVar(&port, "port", "", "port of this peer")
	flag.Parse()
	
	elevio.Init(port, N_FLOORS)
	elevio.SetMotorDirection(elevio.MD_Down)
	
	ordersFromFSMCh := make(chan ElevInfo)
	ordersToFSMCh := make(chan ElevInfo)
	ordersToOrderHandler := make(chan Orders)
	ordersFromOrderHandler := make(chan Orders)
	peersAliveCh := make(chan PeerUpdate)
	setLightCh := make(chan [3]int)

	go fsm.FSM(ordersToFSMCh, ordersFromFSMCh)
	go orderHandler.OrderHandler(ordersFromOrderHandler, ordersToOrderHandler, ordersToFSMCh, ordersFromFSMCh, peersAliveCh, setLightCh, id)
	go network.Network(ordersFromOrderHandler, ordersToOrderHandler, peersAliveCh, setLightCh, id)
	
	select{}
}
