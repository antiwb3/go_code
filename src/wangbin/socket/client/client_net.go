package main

import (
    "fmt"
    "net"
    "os"
)

func chatSend(conn net.Conn){
    
    var input string
    username := conn.LocalAddr().String()
    for {
        
        fmt.Scanln(&input)
        if input == "/quit"{
            fmt.Println("ByeBye..")
            conn.Close()
            os.Exit(0);
        }
        
        
        lens,err :=conn.Write([]byte(username + " Say :::" + input))
        fmt.Println(lens)
        if(err != nil){
            fmt.Println(err.Error())
            conn.Close()
            break
        }
        
    }
    
}

func main() {
    //if len(os.Args) != 2 {
     //   fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
      //  os.Exit(1)
   // }
    service := "127.0.0.1:9090"//os.Args[1]

    tcpAddr, err := net.ResolveTCPAddr("tcp4", service)

    checkError(err, "ResolveTCPAddr")
    fmt.Printf("tcp %s\n", tcpAddr)

    conn, err := net.DialTCP("tcp", nil, tcpAddr)
    checkError(err, "DialTCP")

    fmt.Printf("finish dialtcp\n")

   //启动客户端发送线程
    go chatSend(conn)   
    
    //开始客户端轮训
    buf := make([]byte,1024)
    for{
        
        lenght, err := conn.Read(buf)
        if(checkError(err,"Connection")==false){
            conn.Close()
            fmt.Println("Server is dead ...ByeBye")
           // os.Exit(0)
        }
        fmt.Println(string(buf[0:lenght]))
        
    }
}

func checkError(err error, msg string) (res bool) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s %s ", msg, err.Error())
       // os.Exit(1)
        return false
    }
    return true
}