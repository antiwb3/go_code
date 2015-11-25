package main

import (
    "fmt"
    "net"
    "os"
)


func checkError(err error, msg string) (res bool) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s %s", msg, err.Error())
        return false
    }
    return true
}

func Handler(conn net.Conn,messages chan string){
    
    fmt.Println("Handler connection is connected from ...",conn.RemoteAddr().String())
    
    buf := make([]byte,1024)
    for{
        lenght, err := conn.Read(buf)
        if(checkError(err,"Connection")==false){
            conn.Close()
            break
        }
        if lenght > 0{
            buf[lenght]=0
        }
        //fmt.Println("Rec[",conn.RemoteAddr().String(),"] Say :" ,string(buf[0:lenght]))
        reciveStr :=string(buf[0:lenght])
        messages <- reciveStr
            
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
func echoHandler(conns *map[string]net.Conn,messages chan string){
    
    
    for{
        msg:= <- messages
        fmt.Println(msg)
        
        for key,value := range *conns {
            
            fmt.Println("connection is connected from ...",key)
            _,err :=value.Write([]byte(msg))
            if(err != nil){
                fmt.Println(err.Error())
                delete(*conns,key)
            }
            
        }
    }
    
}


func main() {
    port := "9090"
    service:=":"+port //strconv.Itoa(port);
    tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
    checkError(err,"ResolveTCPAddr")

    l,err := net.ListenTCP("tcp",tcpAddr)
    checkError(err,"ListenTCP")

    conns:=make(map[string]net.Conn)
    messages := make(chan string,10)

//启动服务器广播线程
    go echoHandler(&conns,messages)

    for {
            fmt.Println("Listening ...")
            conn,err := l.Accept()
            checkError(err,"Accept")

            fmt.Println("Accepting ...")
            conns[conn.RemoteAddr().String()] = conn
            //启动一个新线程
            go Handler(conn, messages)       
    }
}
