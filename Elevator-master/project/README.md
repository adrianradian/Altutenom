# TTK4145 - Elevator project

## Overview
A peer to peer system is implemented consisting of a network module, a module to handle orders and a FSM module. 

## Compiling
When compiling main.go, two flags need to be passed as argument `go run main.go -id=# -port=localhost:######`, where the # represents an int.

## Type definitions for the project
We have declared several types to fit the project. All of them are declared in [typeDefinitions.go](https://github.com/IngeborgE/Elevator/blob/master/project/typeDefinitions/typeDefinitions.go)

## Disclaimer
The following packages in [network](./network) where handed out by TTK4145, and we are not the author:
- [bcast](./network/bcast)
- [conn](./network/conn)
- [localip](./network/localip)
- [peers](./network/peers) - Minor changes where done in the file [peers.go](./network/peers/peers.go)

The package [elevio](./elevio) was also given by TTK4145, and works as a driver to connect the software with the elevator hardware,

We have also used following standard go libraries throughout the project:
- `"fmt"`
- `"strconv"`
- `"os"`
- `"flag"`
- `"time"`
- `"sort"`
- `"math"`
