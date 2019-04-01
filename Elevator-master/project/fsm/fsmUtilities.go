package fsm

import(
	"../elevio"
	"time"
	. "../typeDefinitions"
)


func ordersAbove(elevInfo ElevInfo) bool{
	for i := elevInfo.LastFloor + 1; i < N_FLOORS; i++{
		for j := 0; j < N_BUTTONS; j++{
			if elevInfo.Orders[i][j]{
				return true
			}
		}
	}
	return false
}

func ordersBelow(elevInfo ElevInfo) bool{
	for i := 0 ; i < elevInfo.LastFloor; i++{
		for j := 0; j < N_BUTTONS; j++{
			if elevInfo.Orders[i][j]{
				return true
			}
		}
	}
	return false
}

func OrdersChooseDirection(elevInfo ElevInfo) int{
	switch elevInfo.Direction{
	case 1:
		if ordersAbove(elevInfo){
			return 1
		}
		if ordersBelow(elevInfo){
			return -1
		}

		return 0

	case -1, 0:
		if ordersBelow(elevInfo){
			return -1
		}
		if ordersAbove(elevInfo) {
			return 1
		}

		return 0
	}
	return 0
}

func OrdersShouldStop(elevInfo ElevInfo) bool{
	switch elevInfo.Direction{
	case -1:
		return elevInfo.Orders[elevInfo.LastFloor][1]||elevInfo.Orders[elevInfo.LastFloor][2]||!ordersBelow(elevInfo)
	case 1:
		return elevInfo.Orders[elevInfo.LastFloor][0]||elevInfo.Orders[elevInfo.LastFloor][2]||!ordersAbove(elevInfo)
	case 0:
		return true
	}
	return false
}

func OrdersClearAtCurrentFloor(elevInfo ElevInfo) ElevInfo{
	elev := elevInfo
	elev.Orders[elev.LastFloor][2] = false
	switch elev.Direction{
	case 1:
		elev.Orders[elev.LastFloor][0] = false
		if !ordersAbove(elev){
			elev.Orders[elev.LastFloor][1] = false
		}
	case -1:
		elev.Orders[elev.LastFloor][1] = false
		if !ordersBelow(elev){
			elev.Orders[elev.LastFloor][0] = false
		}
	case 0:
		elev.Orders[elev.LastFloor][0] = false
		elev.Orders[elev.LastFloor][1] = false
	}
	return elev
}


// doorOpen called as a go routine from fsm
func doorOpen(close chan<- bool){
	time.Sleep(3*time.Second)
	close <- true
	return
}

func intToMotorDirection(d int) elevio.MotorDirection{
	if d == 1{
		return elevio.MD_Up
	}
	if d == -1{
		return elevio.MD_Down
	}
	return elevio.MD_Stop

}

func checkOrderAtFloor(elevInfo ElevInfo) bool{
	orders := elevInfo.Orders
	lastFloor := elevInfo.LastFloor
	return orders[lastFloor][elevio.BT_HallUp] || orders[lastFloor][elevio.BT_HallDown] || orders[lastFloor][elevio.BT_Cab]
}

func numberOfOrders(elev ElevInfo) int{
	sum := 0
	for i := 0; i < N_FLOORS; i++{
		for j := 0; j < N_BUTTONS; j++{
			if(elev.Orders[i][j]){
				sum++
			}
		}
	}
	return sum
}
