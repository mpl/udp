package main

import (
	"log"
	"net"
)

/*
type Client struct {
	L *net.UDPConn
}

var clients = make(map[string]Client)

var clientLock sync.RWMutex

func handle(client net.Addr, data []byte, listener *net.UDPConn) {
	var clientConn Client
	fmt.Println(client)
	clientLock.Lock()
	if _, ok := clients[client.String()]; !ok {
		fmt.Println("NOT FOUND")
		addr2 := "127.0.0.1:8090"
		udpAddr2, err := net.ResolveUDPAddr("udp", addr2)
		if err != nil {
			log.Fatal("Unable to start UDP Server: %s", err)
		}
		listener2, err := net.DialUDP("udp", nil, udpAddr2)

		// fmt.Println("Start Listener")
		if err != nil {
			log.Fatal("Unable to start UDP Server: %s", err)
		}
		// fmt.Println("Listener", listener2)
		clientConn = Client{
			L: listener2,
		}
		fmt.Println("Listener", clientConn.L)

		fmt.Printf("Write %s\n", client)
		clients[client.String()] = clientConn
		fmt.Println(clients)
		go func() {
			fmt.Println("Waiting response")
			bytes := make([]byte, 1024*1024)
			for {
				n, _, _ := listener2.ReadFrom(bytes)

				listener.WriteTo(bytes[:n], client)
			}
		}()
	} else {
		fmt.Println("Already found")
	}
	clientLock.Unlock()
	// fmt.Println("Listener", clientConn.L)

	listenerT := clients[client.String()].L
	listenerT.Write(data)
}

func startUdpEchoServer(addr string) {
	udpConn, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	bytes := make([]byte, 1024*1024)

	for {
		// fmt.Println("Waiting")
		n, client, err := udpConn.ReadFrom(bytes)
		// fmt.Println("receive bytes")
		if err != nil {
			log.Fatal("Unable to start UDP Server: %s", err)
		}
		// udpConn.WriteTo(bytes[:n], client)
		handle(client, bytes[:n], udpConn)
	}
}
*/

type Handler interface {
	ServeUDP(*dataFlow)
}

type EchoServer struct {
	listenAddr string
}

type dataFlow struct {
	data chan []byte // TODO: async?
	conn net.PacketConn
	peer net.Addr
}

func (s *EchoServer) Start() {
	listener, err := net.ListenPacket("udp", s.listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	data := make([]byte, 1024*1024)
	conns := make(map[net.Addr]*dataFlow)
	for {
		n, client, err := listener.ReadFrom(data)
		if err != nil {
			log.Printf("Error while reading: %v", err)
			continue
		}
		println("BYTES READ: ", n)

		var df *dataFlow
		df, ok := conns[client]
		if !ok {
			df = &dataFlow{
				conn: listener,
				peer: client,
				data: make(chan []byte),
			}
			conns[client] = df
			go func() {
				defer delete(conns, client)
				s.ServeUDP(df)
			}()
		}
		df.data <- data[:n]

	}
}

func (s *EchoServer) ServeUDP(d *dataFlow) {
	for {
		println("WAITING")
		data := <-d.data
		n, err := d.conn.WriteTo(data, d.peer)
		if err != nil {
			log.Printf("WATUP: %v", err)
		}
		println("BYTES SENT: ", n)
	}
}

type Proxy struct {
	listenAddr  string
	forwardAddr string
}

func main() {
	s := EchoServer{listenAddr: ":8081"}
	s.Start()
}
