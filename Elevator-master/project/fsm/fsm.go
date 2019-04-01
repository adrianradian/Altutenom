package fsm

import(
	"../elevio"
	"time"
	. "../typeDefinitions"
)

func FSM(ordersInCh <-chan ElevInfo,ordersOutCh chan<- ElevInfo){
	
	var thisElev ElevInfo
	var currentFloor int
	thisElev.State = OpenDoor
	counter := 0
	sumCounter := 0
	prevSum := 0
	
	lastFloorCh := make(chan int)
	closeDoorCh := make(chan bool)

	go elevio.PollFloorSensor(lastFloorCh)
	elevio.SetDoorOpenLamp(false)
	ticker:=time.NewTicker(time.Millisecond * 200)
	
	for{
		select{
		case newOrders := <-ordersInCh:
			thisElev.Orders = newOrders.Orders
			if (thisElev.State == Idle) && OrdersShouldStop(thisElev) && checkOrderAtFloor(thisElev){
				elevio.SetMotorDirection(elevio.MD_Stop)
				thisElev.State = OpenDoor
				elevio.SetDoorOpenLamp(true)
				go doorOpen(closeDoorCh)
			}
			if thisElev.State != OpenDoor{
				thisElev.Direction = OrdersChooseDirection(thisElev)
				elevio.SetMotorDirection(intToMotorDirection(thisElev.Direction))
				thisElev.State = Idle
				if (thisElev.Direction == -1) || (thisElev.Direction == 1) {
					thisElev.State = Moving
				}
			}
		case thisElev.LastFloor = <-lastFloorCh:
			if OrdersShouldStop(thisElev) {
				elevio.SetMotorDirection(elevio.MD_Stop)
				thisElev.State = OpenDoor
				elevio.SetDoorOpenLamp(true)
				go doorOpen(closeDoorCh)
			}
		case <-closeDoorCh:
			elevio.SetDoorOpenLamp(false)
			thisElev.State = Idle
			thisElev = OrdersClearAtCurrentFloor(thisElev)
			ordersOutCh <- thisElev
			if thisElev.State != OpenDoor{
				thisElev.Direction = OrdersChooseDirection(thisElev)
				elevio.SetMotorDirection(intToMotorDirection(thisElev.Direction))
				if (thisElev.Direction == -1) || (thisElev.Direction == 1){
						thisElev.State = Moving
				}
			}
		case <-ticker.C:
			sum := numberOfOrders(thisElev)
			if sum != prevSum{
				prevSum = sum
				sumCounter = 0
			}else if sum != 0 {
				sumCounter++
			}
			if sumCounter > 100{
				thisElev.State = Stuck
				ordersOutCh <- thisElev
				sumCounter = 0
			}
			if counter == 0{
				currentFloor = thisElev.LastFloor
				counter ++
			}
			if currentFloor != thisElev.LastFloor{
				counter = 0
			} else if thisElev.State == Moving{
				counter++
				if counter > 20 {
					thisElev.State = Stuck
					ordersOutCh <- thisElev
					counter = 0
				}
			}
		}
	}
}
