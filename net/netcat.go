// Chat client
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var (
	port = flag.Int("p", 3090, "port")
	host = flag.String("h", "localhost", "host")
)

func CopyContent(dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src)
	if err != nil {
		log.Fatal(err)
	}
}

// -> host:port
// Escribir -> host:port
// Leer -> host:port
// > Hola -> host:port -> [Hola]
func main() {
	flag.Parse()
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatal(err)
	}
	// Canal de control
	done := make(chan struct{})
	go func() {
		// Aca conn actua como io.Reader
		io.Copy(os.Stdout, conn)
		done <- struct{}{}
	}()
	// Aca conn actua como io.Writer
	CopyContent(conn, os.Stdin)
	conn.Close()
	// Blockeamos el programa hasta que el done channel termine
	<-done
}
