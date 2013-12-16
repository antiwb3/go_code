package netserver

func send_data(data *NetPacket) bool {

	return true
}

func (server *NServer) on_active() {
	l := server.data_list

	if l.Len() > 0 {
		server.send_netpacket()
	}
	return true
}

func (server *NServer) recive_netpacket(data *NetPacket) bool {
	if data != nil {
		server.data_list.PushBack(data)
	}

	return true
}

func (server *NServer) send_netpacket() bool {
	l := server.data_list
	for e := l.Front(); e != nil; e = e.Next() {
		send_data(data)
	}
	return true
}
