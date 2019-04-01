# Design
## General comments
I am alive messages are sent every X moment. We do not assume a PC is dead before we haven’t received messages from it in N seconds.

### What should happen if one of the nodes loses network?

If an elevator has people inside an has started to move, it should finish its orders alone, as we cannot have to errors at the same time. 

If there are no people in an elevator that is about to perform a job, and it loses network, stop the elevator and send another elevator. 

We assume that there are only people inside an elevator if someone clicked on a button inside.


### What should happen if one of the computers loses power for a brief moment?

If a computer loses power while there are people inside the elevator stops. The two other computers will then stop receiving I am alive messages, and will start sending it messages of what to do.

If a computer loses power while there are no one inside, one of the other elevators will take care of the order. (The computer that said it would complete the order will no longer complete it and it needs to be reassigned)

If an elevator has got an order to pick up someone at floor X, but then loses power or internet before  someone enters the elevator. We will wait for N seconds, before we send another elevator.

 
### What should happen if some unforeseen event causes the elevator to never reach its destination but communication remains intact?

If there are people inside we try a restart and hope for the best.
If there are no people inside we can either restart or send another elevator.

### Do all your nodes need to "agree" on an order for it to be accepted? In that case, how is a faulty node handled?

No, all the nodes do not have to agree. If a node has lost network or power (does not send IAA), we only need confirmation from one other node. The faulty node is then restarted.

How can you be sure that at least as many nodes as needs to agree on the order in fact agrees on the order?

All the nodes that sends IAA should agree on the order.

Do you share the entire state of the current orders , or just the changes as they occur? 
We share the entire state of the current orders with all the elevators that are alive. 
	What should happen when an elevator re-joins after having been offline?
	Re-assign all current orders


## Topology
### What kind of network topology do you want to implement? Peer to peer? Master slave? Circle? 

Peer to peer

## Technical implementation
### Will you be using blocking sockets & many threads, or nonblocking sockets and select()? 

### Do you want to build the necessary reliability into the module, or handle that at a higher level? 

No, handle it at a higher level

### How will you pack and unpack (serialize) data? 
#### Do you use structs, classes, tuples, lists, ...? 
	
Struct for each elevator woth information about state, floor and so on	

#### JSON, XML, or just plain strings?

Plain strings

#### Is serialization a part of the network module? 

Yes

#### Is detection (and handling) of things like lost messages or lost nodes a part of the network module? 

We take care of it outside the network module


## Main requirements
Once the light on an hall call button (buttons for calling an elevator to that floor; top 6 buttons on the control panel) is turned on, an elevator should arrive at that floor. 

Similarly for a cab call (for telling the elevator what floor you want to exit at; front 4 buttons on the control panel), but only the elevator at that specific workspace should take the order.
