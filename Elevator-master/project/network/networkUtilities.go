package network

import (
	"fmt"
	"strconv"
	."../typeDefinitions"
)

func arrayOfZeros(size int) []int{
	var array []int
	for i := 0; i < size; i++ {
		array = append(array, 0)
	}
	return array
}

func sum(array []int) int{
    sum := 0
    for i := 0; i < len(array); i++ {
        sum += array[i]
    }
    return sum
}

func correspondingIndex(array [] string, key int) (int, string){ 
	for i :=range array {
		intRep, error:=strconv.Atoi(array[i])
		if error != nil{
			fmt.Println(error)
			return 0, "Error: Stringconversion failed"
		}
		if intRep == key {
			return i, "nil"
		}
	}
	return 0, "Error: id not found"
}

func activateUnconfirmedOrder(order Orders, peersAlive int) UnconfirmedOrder {
	var newUnconfirmedOrder UnconfirmedOrder
	newUnconfirmedOrder.Active = true
	newUnconfirmedOrder.Counter = 0
	newUnconfirmedOrder.ReceivedAcks = arrayOfZeros(peersAlive)
	newUnconfirmedOrder.OrderToShare = order
	return newUnconfirmedOrder
}

func createAcknowledgement(receivedOrders Orders, myId int) Acknowledgement {
	var ack Acknowledgement
	ack.From = myId
	ack.Floor = receivedOrders.Floor
	ack.Button = receivedOrders.Button
	ack.ElevToExecute = receivedOrders.ElevToExecute
	return ack
}
