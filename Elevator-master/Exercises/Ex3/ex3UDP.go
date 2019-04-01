package main
import "net"
import "fmt"
import "time"
import "runtime"
//ServerIP = 10.100.23.242

func receiver(ipAndPort string){
	addr, _ := net.ResolveUDPAddr("udp", ipAndPort)
	conn, err := net.ListenUDP("udp", addr)
	if err!=nil{
		fmt.Println(err)
	}
	buf := make([]byte, 1024)
	fmt.Println("Starting")
	for {
		n, err := conn.Read(buf)
		if err!=nil{
			fmt.Println(err)
			return
		}
		fmt.Println(string(buf[:n]))	
	}
}


func sender(ipAndPort string){
	// dialUDP
	addr, err2 := net.ResolveUDPAddr("udp", ipAndPort)
	conn,_ :=net.DialUDP("udp",nil,addr)
	message := []byte("Hello, from 20001")
	if err2!=nil{
		fmt.Println(err2)
		return
	}

	for{
		// Skriv noe
		_, err := conn.Write(message)
		
		if err!=nil{
			fmt.Println(err)
			return
		}

		time.Sleep(1*time.Second)
	}
}



func main(){
	runtime.GOMAXPROCS(2)
	
	
	//receiver(":30000")
	
	//sender("10.100.23.242:20001")
	
	

	
	go sender("10.100.23.242:20001")
	go receiver(":20001")

	
	
	select{}	
//simple write
//pc.WriteTo([]byte("Hello from client"), addr)
}



