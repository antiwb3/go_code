package netserver

import "container/list"

type NetPacket struct {
	dst_ip string

	data_size uint
	data      []byte
}

type NServer struct {
	data_list list.List
}
