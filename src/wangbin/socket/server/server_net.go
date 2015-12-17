package main

import (
	"fmt"
	"net"
	"os"
    "strings"
)

type ConnectedClient struct {
	clientIPPort string
	gsIPPort     string
    gsConn       net.Conn
	conn         net.Conn
}

type CommonConnect struct {
    id      int
    conn    net.Conn
}

m_aclient_ids := make(map[string]int)
m_gclient_ids := make(map[string]int)
m_aclient := make(map[int]*CommonConnect)
m_gclient := make(map[int]*ConnectedClient)

func checkError(err error, msg string) (res bool) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s %s", msg, err.Error())
		return false
	}
	return true
}

func parseIP(address string) (res string) {
    res = strings.Split(address, ":")
    if res == nil {
        return nil
    }

    return res[0]
}

func getAClientId(address string) (id int) {
    sip :=parseIP(address)
    if m_aclient_ids[sip] == nil {
        m_aclient_ids[sip] = 0
    }
    else {
        m_aclient_ids[sip]++;
    }
    return m_aclient_ids[sip]
}

func getGClientId(address string) (id int) {
    sip :=parseIP(address)
    if m_gclient_ids[sip] == nil {
        m_gclient_ids[sip] = 0
    }
    else {
        m_gclient_ids[sip]++;
    }
    return m_gclient_ids[sip]
}

func ProccessGSData(cclient *ConnectedClient) {
    var conn net.Conn = cclient.gsConn

    fmt.Println("Handler gs c is connected from ...", conn.RemoteAddr().String())

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
        _, err1 := cclient.conn.Write(buf[0:lenght])
        if err1 != nil {
            fmt.Println(err1.Error())
        }
    }
}

func ConnectGS(service string) (net.Conn) {
    tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
    checkError(err, "ResolveTCPAddr")
    conn, err := net.DialTCP("tcp", nil, tcpAddr)
    checkError(err, "DialTCP")

    fmt.Println("finish connectGS:" + service)

    return conn
}

func ProcessClientData(cclient *ConnectedClient, ch_cclients chan *ConnectedClient, messages chan string) {

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
// gs_ipport:127.0.0.1:3724
////////////////////////////////////////////////////////
func ProcessGC2GSData(conns *map[string]*ConnectedClient, ch_cclients chan *ConnectedClient, messages chan string) {

	for {
		msg := <-messages // data come from game client
		cclient := <-ch_cclients

        if cclient.gsConn == nil {
            fmt.Println("msg:" + msg)
            if strings.HasPrefix(msg, "gs_ipport:") {
                cclient.gsIPPort = string(msg[len("gs_ipport:"):])
                cclient.gsConn = ConnectGS(cclient.gsIPPort)
                go ProccessGSData(cclient)
            } 
            continue
        }
      
		fmt.Println("cclient.gsConn:" + msg)

		//for key,value := range *conns {

		//fmt.Println("connection is connected from ...", cclient.clientIPPort)
		_, err := cclient.gsConn.Write([]byte(msg))
		if err != nil {
			fmt.Println(err.Error())
			delete(*conns, cclient.clientIPPort)
		}
	}

}

func ReceiveAClientData(cconn *CommonConnect, ch_cconn chan *CommonConnect, messages chan string) {
    var conn net.Conn = cconn.conn

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
        ch_cconn <- cconn

    }
}

func ProcessAClientData(ch_cconn chan *CommonConnect, messages chan string) {
    for {
        msg := <-messages
        cconn := <-ch_cconn
        fmt.Println("msg:" + msg)
        if strings.HasPrefix(msg, "gs_ipport:") {
            var cclient = m_gclient[cconn.id]
            cclient.gsIPPort = string(msg[len("gs_ipport:"):])
            cclient.gsConn = connectGS(cclient.gsIPPort)
            if cclient.gsConn != nil {
                go ProccessGSData(cclient)    
            }
        } 
      
        fmt.Println("cclient.gsConn:" + msg)

        //for key,value := range *conns {

        //fmt.Println("connection is connected from ...", cclient.clientIPPort)
        _, err := cclient.gsConn.Write([]byte(msg))
        if err != nil {
            fmt.Println(err.Error())
            delete(*conns, cclient.clientIPPort)
        }

        // }
    }
}

func startListenAssistTCP() {
    port := "9091"
    service := ":" + port //strconv.Itoa(port);
    tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
    checkError(err, "ResolveTCPAddr")

    l, err := net.ListenTCP("tcp", tcpAddr)
    checkError(err, "ListenTCP")

    assistConns := make(map[string]*ConnectedClient)
    messages := make(chan string, 10)
    ch_cclients := make(chan *ConnectedClient, 10)

    //启动服务器广播线程
    go ProcessAClientData(&conns, ch_cclients, messages)

    for {
        fmt.Println("Listening ...")
        conn, err := l.Accept()
        checkError(err, "Accept")

        fmt.Println("Accepting ...")

        address := conn.RemoteAddr().String()
        id := getAClientId(address)
        cconn := new(CommonConnect)
        cconn.id = id
        cconn.conn = conn
        m_aclient[id] = cconn

        if cclient.gsIPPort == "" { 
            fmt.Println("gsIPPort isnil") 
        }

        assistConns[conn.RemoteAddr().String()] = cclient
        //启动一个新线程
        go ReceiveAClientData(cclient, ch_cclients, messages)
    }
}

func main() {
    goto startListenAssistTCP()

	port := "9090"
	service := ":" + port //strconv.Itoa(port);
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err, "ResolveTCPAddr")

	l, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err, "ListenTCP")

	conns := make(map[string]*ConnectedClient)
	messages := make(chan string, 10)
	ch_cclients := make(chan *ConnectedClient, 10)

	//启动服务器广播线程
	go ProcessGC2GSData(&conns, ch_cclients, messages)

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
		go ProcessClientData(cclient, ch_cclients, messages)
	}
}
