package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
)

var site = flag.String("site", "scanme.nmap.org", "url to scan")

func main() {
	flag.Parse()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			// 1, 2, ..., 99
			// sitio:1, sitio:2, ..., sitio:99
			// 1 -> Open, 2 -> Closed, ...,
			conn, err := net.Dial(
				"tcp",
				fmt.Sprintf("%s:%d", *site, port))
			if err != nil {
				return
			}
			conn.Close()
			fmt.Printf("Port %d is open \n", port)
		}(i)
	}
	wg.Wait()
}
