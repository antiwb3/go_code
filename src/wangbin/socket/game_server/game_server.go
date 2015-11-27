package main

import (
	"fmt"
	"net"
	"os"
)

type ConnectedClient struct {
	clientIPPort string
	gsIPPort     string
	conn         net.Conn
}

func checkError(err error, msg string) (res bool) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s %s", msg, err.Error())
		return false
	}
	return true
}

func Handler(cclient *ConnectedClient, ch_cclients chan *ConnectedClient, messages chan string) {

	var conn net.Conn = cclient.conn

	fmt.Println("Handler connection is connected from ...", conn.RemoteAddr().String())

	buf := make([]byte, 1024)
	for {
		lenght, err := conn.Read(buf)
		if checkError(err, "Connection") == false {
			conn.Close()
			break
		}
		if lenght > 0 {
			buf[lenght] = 0
		}
		fmt.Println("Rec[",conn.RemoteAddr().String(),"] Say :" ,string(buf[0:lenght]))
		reciveStr := string(buf[0:lenght])
		messages <- reciveStr
		ch_cclients <- cclient

	}
}

////////////////////////////////////////////////////////
//
//服务器发送数据的线程
//
//参数
//      连接字典 conns
//      数据通道 messages
//
////////////////////////////////////////////////////////
func echoHandler(conns *map[string]*ConnectedClient, ch_cclients chan *ConnectedClient, messages chan string) {

	for {
		msg := <-messages
		cclient := <-ch_cclients
		//fmt.Println(msg)

		//for key,value := range *conns {

		//fmt.Println("connection is connected from ...", cclient.clientIPPort)
		_, err := cclient.conn.Write([]byte("game_server_send:" + msg))
		if err != nil {
			fmt.Println(err.Error())
			delete(*conns, cclient.clientIPPort)
		}

		// }
	}

}

func main() {
	port := "3724"
	service := ":" + port //strconv.Itoa(port);
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err, "ResolveTCPAddr")

	l, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err, "ListenTCP")

	conns := make(map[string]*ConnectedClient)
	messages := make(chan string, 10)
	ch_cclients := make(chan *ConnectedClient, 10)

	//启动服务器广播线程
	go echoHandler(&conns, ch_cclients, messages)

	for {
		fmt.Println("Listening ...")
		conn, err := l.Accept()
		checkError(err, "Accept")

		fmt.Println("Accepting ...")

		cclient := new(ConnectedClient)
        cclient.conn = conn
        cclient.clientIPPort = conn.RemoteAddr().String()

        if cclient.gsIPPort == "" { 
            fmt.Println("gsIPPort isnil") 
        }

		conns[conn.RemoteAddr().String()] = cclient
		//启动一个新线程
		go Handler(cclient, ch_cclients, messages)
	}
}
