package main

import (
	"bytes"
	"log"
	"math/rand"
	"net"
	"time"
)

var loremIpsum = []byte(`
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed luctus pellentesque mauris, non laoreet sem hendrerit eu. Nulla mi ex, interdum nec libero quis, hendrerit feugiat erat. Aenean blandit rutrum mollis. Nullam consectetur odio ac turpis imperdiet accumsan. Nam lobortis fermentum lectus, in auctor nisl condimentum in. Mauris cursus leo sed dolor facilisis, eget rutrum ligula pretium. Proin justo odio, aliquam quis vulputate vel, viverra molestie ligula. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas eget ultrices urna. Pellentesque vehicula, magna et sodales vulputate, nulla nulla malesuada ex, eu elementum mauris nulla vel arcu. Fusce vel dui risus. Donec pharetra mattis nisl id scelerisque. Vivamus ut luctus nibh, vitae cursus tortor.

In eu tellus quam. Ut vitae urna id arcu tincidunt scelerisque ac ac odio. Quisque gravida nisl eget lacus faucibus, at tempus dolor feugiat. In porttitor eget urna sit amet tincidunt. Phasellus in interdum nisi, ac convallis sem. Fusce et felis non nunc pharetra commodo. Nullam at sagittis nunc. Vivamus fringilla quis mauris nec blandit. Fusce laoreet volutpat urna a aliquet. Suspendisse eget orci feugiat, luctus nunc id, ultrices augue. Pellentesque vitae auctor lacus, consectetur consequat quam. Cras vel leo sapien.

Praesent consequat orci non cursus fermentum. Morbi et sapien blandit, tempor erat at, ultrices nisi. Fusce vitae finibus enim. Suspendisse tincidunt odio ultrices facilisis efficitur. Morbi eleifend sem sit amet ipsum feugiat semper. Integer nisl neque, ultrices non leo eu, porta tristique velit. Sed sollicitudin, diam at rutrum egestas, lorem dolor tempus est, a bibendum quam risus in velit. Etiam sed turpis sit amet mauris pretium egestas vel sed nisl. Nunc tortor lectus, gravida et lobortis ornare, faucibus ut tellus. Curabitur tincidunt imperdiet iaculis. Etiam iaculis sit amet arcu quis rutrum. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum sed velit efficitur, fringilla justo id, lobortis nunc. Mauris consequat nibh quis magna sagittis efficitur. Phasellus gravida sollicitudin leo, id tristique mauris iaculis id.

Ut a elit nisl. Aliquam ultricies lorem a pellentesque fermentum. Sed placerat enim lacus, ut volutpat massa auctor in. Donec id faucibus libero. Donec in mollis elit. Nulla vitae imperdiet ipsum. Aenean non diam purus. Proin euismod ullamcorper magna et euismod. Donec eu nunc congue, sollicitudin augue id, iaculis mi. Sed a neque nisi. Proin odio justo, tempus sit amet vulputate ac, sodales nec neque. Mauris malesuada eros at porta bibendum.

Fusce vel porttitor arcu. Etiam ultrices vitae neque at lacinia. Curabitur leo est, vestibulum vel faucibus sed, euismod nec mauris. Proin egestas ultrices eros vitae mattis. Ut justo massa, bibendum ac ultricies sed, sagittis quis nunc. Pellentesque nec massa commodo, finibus tortor volutpat, bibendum tortor. Nam et gravida sem.
`)

func main() {
	conn, err := net.Dial("udp", "localhost:8081")
	if err != nil {
		log.Fatal(err)
	}
	stop := time.Now().Add(10 * time.Second)
	idx := 0
	remoteAddr := conn.RemoteAddr()
	println("SENDING TO: ", remoteAddr.String())
	udpConn, ok := conn.(*net.UDPConn)
	if !ok {
		log.Fatalf("NOPE: %T", conn)
	}
	//	bytesRead := make([]byte, 0, 1024)
	bytesRead := make([]byte, 1024)
	var response []byte
	waitforit := make(chan bool)
	go func() {
		for {
			if time.Now().After(stop) {
				break
			}
			//			n, _, err := udpConn.ReadFrom(bytesRead)
			n, err := udpConn.Read(bytesRead)
			if err != nil {
				log.Fatal(err)
			}
			response = append(response, bytesRead[:n]...)
			if len(response) >= len(loremIpsum) {
				break
			}
		}
		if !bytes.Equal(response, loremIpsum) {
			log.Print("response mismatch")
			if len(response) != len(loremIpsum) {
				log.Fatalf("length mismatch: got %d, wanted %d", len(response), len(loremIpsum))
			}
			for k, v := range loremIpsum {
				if val := response[k]; val != v {
					log.Fatalf("mismatch at %d, got %v, wanted %v", k, val, v)
				}
			}
		}
		waitforit <- true
	}()
	for {
		if time.Now().After(stop) {
			break
		}
		inc := rand.Intn(30)
		upper := idx + inc
		if upper > len(loremIpsum) {
			upper = len(loremIpsum)
			inc = upper - idx
		}
		//		n, err := udpConn.WriteTo(loremIpsum[idx:upper], remoteAddr) // interesting error
		n, err := conn.Write(loremIpsum[idx:upper])
		if err != nil {
			log.Fatal(err)
		}
		if n != inc {
			log.Fatalf("wanted: %d, got: %d", inc, n)
		}
		idx += n
		if idx == len(loremIpsum) {
			break
		}

		time.Sleep(10 * time.Millisecond)
	}

	<-waitforit

}
