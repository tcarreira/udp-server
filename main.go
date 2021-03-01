package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	isClient      = flag.Bool("client", false, "whether it should be run as a client or server")
	port          = flag.Uint("port", 1337, "port to send to or receive from")
	host          = flag.String("host", "0.0.0.0", "address to send to or receive from")
	timeout       = flag.Duration("timeout", 15*time.Second, "read and write blocking deadlines")
	input         = flag.String("input", "", "file with contents to send over udp")
	maxBufferSize = flag.Int("buffer", 1024*8, "buffer size allocated for receiving UDP packets")
	hexDumpBytes  = flag.Int("hex-bytes", 32, "how many bytes do dump as hexadecimal")
)

// server wraps all the UDP echo server functionality.
// ps.: the server is capable of answering to a single
// client at a time.
func server(ctx context.Context, address string) (err error) {

	pc, err := net.ListenPacket("udp", address)
	if err != nil {
		fmt.Println("Fatal Error:", err)
		return
	}
	defer pc.Close()

	doneChan := make(chan error, 1)
	go func() {
		for {
			buffer := make([]byte, *maxBufferSize)
			n, addr, err := pc.ReadFrom(buffer)
			if err != nil {
				doneChan <- err
				return
			}

			hexDump := ""
			for i, v := range buffer {
				if i >= *hexDumpBytes {
					break
				}
				if i%8 == 0 {
					hexDump += fmt.Sprintf("\n%5d:     ", i)
				}
				hexDump += fmt.Sprintf("%02x ", v)
			}

			fmt.Printf("\n##### packet-received: bytes=%d from=%s\n", n, addr.String())
			if *hexDumpBytes > 0 {
				fmt.Printf("##### binary/hex dump (first %d bytes):", *hexDumpBytes)
				fmt.Printf("%s\n%s\n", hexDump, strings.Repeat("#", 40))
			}
			fmt.Println(string(buffer))

			// not sure if this is needed
			deadline := time.Now().Add(*timeout)
			err = pc.SetWriteDeadline(deadline)
			if err != nil {
				doneChan <- err
				return
			}
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("cancelled")
		err = ctx.Err()
	case err = <-doneChan:
	}

	return
}

func client(ctx context.Context, address string, reader io.Reader) (err error) {
	// defaults to 127.0.0.1 when using client, as 0.0.0.0 makes no sense in context
	address = strings.Replace(address, "0.0.0.0:", "127.0.0.1:", 1)

	fmt.Println("sending to " + address)

	raddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return
	}
	defer conn.Close()

	doneChan := make(chan error, 1)

	go func() {
		n, err := io.Copy(conn, reader)
		if err != nil {
			doneChan <- err
			return
		}

		fmt.Printf("packet-written: bytes=%d\n", n)
		doneChan <- nil
	}()

	select {
	case <-ctx.Done():
		fmt.Println("cancelled")
		err = ctx.Err()
	case err = <-doneChan:
	}

	return
}

func main() {
	flag.Parse()

	var (
		err     error
		address = fmt.Sprintf("%s:%d", *host, *port)
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Gracefully handle signals so that we can finalize any of our
	// blocking operations by cancelling their contexts.
	go func() {
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		cancel()
	}()

	if *isClient {
		var reader io.Reader
		if *input == "" {
			reader = strings.NewReader(strings.Join(flag.Args(), ", "))
		} else if *input == "-" {
			reader = os.Stdin
		} else {
			file, err := os.Open(*input)
			if err != nil {
				panic(err)
			}
			defer file.Close()
			reader = file
		}

		err = client(ctx, address, reader)
		if err != nil && err != context.Canceled {
			panic(err)
		}

		return
	}

	fmt.Println("running as a server on " + address)
	err = server(ctx, address)
	if err != nil && err != context.Canceled {
		panic(err)
	}
}
