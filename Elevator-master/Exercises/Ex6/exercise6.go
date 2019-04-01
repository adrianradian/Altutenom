package main
import "net"
import "fmt"
import "time"
import "strconv"
//import "runtime"
import "os/exec"

//exec.Command("gnome-terminal -x","go run","main.go")


func timer(){
	
}

//ServerIP = 10.100.23.242
func sender(ipAndPort string, i int){
	cmd := exec.Command("gnome-terminal", "-x","go", "run", "/home/student/Desktop/Gruppe54/main.go")
	err := cmd.Run()
	time.Sleep(2*time.Second)
	if(err != nil){
		fmt.Println(err)
	}
	fmt.Println("heihei")
	addr, err2 := net.ResolveUDPAddr("udp", "localhost"+ipAndPort)
	fmt.Println(addr)
	conn, err2 := net.DialUDP("udp",nil,addr)

		if err2!=nil{
		fmt.Println(err2)
		return
	}
	
	for{
		
		message := []byte(fmt.Sprint(i))
		i+=1
		fmt.Println("sending:", i)
		_, err := conn.Write(message)
		fmt.Println(i)
		//fmt.Println("%d\n",i)
		if err!=nil{
			fmt.Println(err)
			return
		}

		time.Sleep(1*time.Second)
	}
}/*
func timeOut(chan number ->string, chan spawn<- bool){
	while(1){


		for{
			select{
				case 
				//hvis vi får inn et tall, send false på spawn
				//hvis ikke send true
			}
		}

	}

	
	
}
*/

func receiver(ipAndPort string,i int)(int, bool){
	addr, _ := net.ResolveUDPAddr("udp", ipAndPort)
	conn, err := net.ListenUDP("udp", addr)
	defer conn.Close()
	var istring = "error"
	if err!=nil{
		fmt.Println(err)
	}
	buf := make([]byte, 1024)
	fmt.Println("Starting")
	
	for {
		conn.SetReadDeadline(time.Now().Local().Add(time.Duration(3)*1000000000))
		fmt.Println("reading number")
		n, err := conn.Read(buf)
		if err!=nil{
			fmt.Println("UDPerr:", err)
			return i, false
		}
	
		istring=string(buf[:n])
		i, err = strconv.Atoi(istring)
		if err != nil{
			fmt.Println(err)
		}
		fmt.Println(i)
	}
}



func main(){
	var i int = 0
	var primaryAlive bool = true
	
	for primaryAlive{
		i, primaryAlive = receiver(":12341",i)
		
	}
	sender(":12341", i)
	 //CAstes

	
	
	select{}	
//simple write
//pc.WriteTo([]byte("Hello from client"), addr)
}


// A new process starts and it assumes that it is the backup
	// The process sets up a UDP listen a certain port and waits X-seconds before it decides that it is the primary
	// Depending on whether or not we recieved a IAA from primary, we set the state to primary or backup
	// An if or switch statement is used to place us in the correct state

	//Primary state
	// A backup is created
	// The Primary loop begins

		//Primary Loop
			// We

	// Backup state
	// We move into a while loop until we stop recieving IAA messages within a certain time)

		// Backup Loop
			// We listen for IAA messages and resets the timer when we recieve one
	// When we stop recieving IAA we exit the loop and jump into the primary state 
