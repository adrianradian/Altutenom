package main
import "net"
import "fmt"
import "time"


func listener(){
	addr,err:= net.ResolveTCPAddr("tcp", ":20001")
	if err!=nil{
		fmt.Println(err)
		return
	}
	ln, err := net.ListenTCP("tcp", addr)
	if err!=nil{
		fmt.Println(err)
		return
	}
	conn, err := ln.Accept()
	defer conn.Close()
	if err!=nil{
		fmt.Println(err)
		return
	}
	fmt.Println("Server connected to us")
}

func reciver(conn *net.TCPConn){
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


func sender(conn *net.TCPConn){
	message := []byte("Hei\u0000")

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
	addr,err2:= net.ResolveTCPAddr("tcp", "10.100.23.242:33546")
	if err2!=nil{
		fmt.Println(err2)
		return
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err!=nil{
		fmt.Println(err2)
		return
	}
	conn.Write([]byte("Connect to: 10.100.23.180:20001\u0000"))


	//go sender(conn)
	//go reciver(conn)
	go listener()
	select{}	
}

