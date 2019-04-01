package typeDefinitions
 const TRAVEL_TIME = 5
const DOOR_OPEN_TIME = 5

const N_FLOORS = 4
const N_BUTTONS_HALL = 2
const N_ELEVATORS = 3
const N_BUTTONS = 3
type State int
const (
	Idle = 0
	Moving = 1
	OpenDoor = 2
	Stuck = 3
)

type Orders struct{
	ElevInfos [N_ELEVATORS]ElevInfo
	Floor int
	Button int
	ElevToExecute int
	NewOrder bool
}

type ButtonPress struct{
	Floor int 
	Button int
	ElevToExecute int
	NewOrder bool
}

type ElevInfo struct{
		Orders [N_FLOORS][N_BUTTONS]bool
		LastFloor int
		Direction int
		State State		
}
type UnconfirmedOrder struct{
	Counter int
	Active bool
	ReceivedAcks []int
	OrderToShare Orders
}

type Acknowledgement struct {
	From int
	Floor int
	Button int
	ElevToExecute int
}

type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}
