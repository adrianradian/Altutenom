package orderHandler

import (
	"../elevio"
	"../fsm"
	"math"
	."../typeDefinitions"
)

// The function checks for changes in previous set of orders and the current set of orders
// If there is a difference the previous set is updated to the current set
func updateOrders(newOrders Orders, ordersInfo Orders) Orders{
	for i := 0; i < N_ELEVATORS; i++{
		for j := 0; j < N_FLOORS; j++{
			for k := 0; k < N_BUTTONS; k++{
				if ordersInfo.ElevInfos[i].Orders[j][k] != newOrders.ElevInfos[i].Orders[j][k] {
					ordersInfo.ElevInfos[i].Orders[j][k] = !ordersInfo.ElevInfos[i].Orders[j][k]
				}
			}
		}
	}
	return ordersInfo
}


func updateLightsOff(ordersInfo Orders, peersAlive []int, myId int){
	for k := elevio.BT_HallUp; k < elevio.BT_Cab; k++{
		for j := 0; j < N_FLOORS; j++{
			turnOffLight := false
			for _, i := range peersAlive{
				if(ordersInfo.ElevInfos[i].Orders[j][k]){
					turnOffLight = true
				}
			}
			if(!turnOffLight){
				elevio.SetButtonLamp(k, j, false)
			}
		}
	}
	for j := 0; j < N_FLOORS; j++{
		if !ordersInfo.ElevInfos[myId].Orders[j][elevio.BT_Cab]{
				elevio.SetButtonLamp(elevio.BT_Cab, j, false)
			}
	}
}

func timeToIdle(elevInfo ElevInfo) int{
    duration := 0
    switch(elevInfo.State){
    case Idle:
        elevInfo.Direction = fsm.OrdersChooseDirection(elevInfo)
        if elevInfo.Direction == 0{
            return duration
        }
    case Moving:
        duration += TRAVEL_TIME/2
        elevInfo.LastFloor += elevInfo.Direction
    case OpenDoor:
        duration += DOOR_OPEN_TIME/2
    }
    for{
        if(fsm.OrdersShouldStop(elevInfo)){
            elevInfo = fsm.OrdersClearAtCurrentFloor(elevInfo)
            duration += DOOR_OPEN_TIME
            elevInfo.Direction = fsm.OrdersChooseDirection(elevInfo)
            if(elevInfo.Direction == 0){
                return duration
            }
        }
        elevInfo.LastFloor += elevInfo.Direction
        duration += TRAVEL_TIME
    }
}

// The function calculates which elevator that is going to execute the order
func distributeTo(orderInfo Orders, peersAlive []int, floor int) int{
	best := timeToIdle(orderInfo.ElevInfos[peersAlive[0]])
	bestElev := peersAlive[0]
	bestElevDist := math.Abs(float64(floor - orderInfo.ElevInfos[peersAlive[0]].LastFloor))
	for _,k := range peersAlive{
		duration := timeToIdle(orderInfo.ElevInfos[k])
		dist := math.Abs(float64(floor - orderInfo.ElevInfos[k].LastFloor))
		if(duration < best && orderInfo.ElevInfos[k].State != Stuck){
			best = duration
			bestElev = k
			bestElevDist = math.Abs(float64(floor - orderInfo.ElevInfos[k].LastFloor))
		} else if (duration == best) && (bestElevDist > dist){
				bestElev = k
				bestElevDist = dist
		}
	}
	return bestElev
}

// The function redistributes the orders from the lost elevator to the other active elevators
func redistributeOrders(ordersInfo Orders, peersAlive [] int, elevLost int) Orders{
	for j := 0; j < N_FLOORS; j++{
		for k := elevio.BT_HallUp; k < elevio.BT_Cab; k++{
			if ordersInfo.ElevInfos[elevLost].Orders[j][k]{
				var best int
				bestElev := -1
				for _, x := range peersAlive{
					if x != elevLost{
						duration := timeToIdle(ordersInfo.ElevInfos[x])
						if (bestElev == -1) || (duration < best) && (ordersInfo.ElevInfos[x].State != Stuck){
							best = duration
							bestElev = x
						}
					}
				}
                ordersInfo.ElevInfos[bestElev].Orders[j][k] = true
			}
		}
	}
	return ordersInfo
}

// The function remove all the hall calls for the given elevator
func removeHallCalls(ordersInfo Orders, elevLost int) Orders{
	for i := 0; i < N_FLOORS; i++{
		for j := elevio.BT_HallUp; j < elevio.BT_Cab; j++{
			ordersInfo.ElevInfos[elevLost].Orders[i][j] = false
		}
	}
	return ordersInfo
}

func fromIntToButtonType(b int) elevio.ButtonType{
	if b == 0 {
		return elevio.BT_HallUp
	}
	if b == 1 {
		return elevio.BT_HallDown
	}
	return elevio.BT_Cab
}

func fromButtonTypeToInt(b elevio.ButtonType) int{
	if b == elevio.BT_HallUp{
		return 0
	}
	if b == elevio.BT_HallDown{
		return 1
	}
	return 2
}

func directionToButton(direction int)int{
	if(direction == 1){
		return 1
	}else if(direction == -1){
		return 0
	} else {
		return 2
	}	
}
