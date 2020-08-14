package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"
	"sync"
	"time"
)

type PortScan struct {
	timeOut       time.Duration
	protocol      string
	cores         int
	threads       int
	addr          string
	availablePort []int
	wasteTime     int64
	rwLocker      sync.RWMutex
}

var protocol = flag.String("p", "tcp", "Use the specified protocol")
var cores = flag.Int("c", 1, "Use the specified core number")
var threads = flag.Int("t", 0, "Use the specified thread number")
var timeout = flag.Int("d", 5, "Use the specified duration")

func main() {
	var portScanner *PortScan
	portScanner = new(PortScan)
	if err := portScanner.Scan(); err != nil {
		fmt.Println(err)
	}
}

func (ps *PortScan) Scan() error {
	start := time.Now().Unix()
	if err := ps.parse(); err != nil {
		return err
	}
	fmt.Println("开始端口扫描...")
	concurrent := make(chan bool, ps.threads)
	ports := 65535
	for port := 1; port <= ports; port++ {
		concurrent <- true
		go ps.check(port, concurrent)
	}
	end := time.Now().Unix()
	ps.wasteTime = end - start
	fmt.Printf("扫描结束,耗时%d秒,%s:%s共扫描%d个可用端口\r\n", ps.wasteTime, ps.protocol, ps.addr, len(ps.availablePort))
	return nil
}

func (ps *PortScan) appendAvailablePort(port int) {
	ps.rwLocker.Lock()
	defer ps.rwLocker.Unlock()
	ps.availablePort = append(ps.availablePort, port)
}

func (ps *PortScan) check(port int, concurrent <-chan bool) {
	ps.connect(port)
	<-concurrent
}

func (ps *PortScan) connect(port int) {
	cn := make(chan int, 1)
	go func() {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ps.addr, port))
		if err == nil {
			conn.Close()
			cn <- port
		} else {
			cn <- 0
		}
	}()
	select {
	case result := <-cn:
		if result > 0 {
			ps.appendAvailablePort(result)
			ps.stdout(result)
		}
	case <-time.After(ps.timeOut):
	}
	return
}

func (ps *PortScan) parse() error {
	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	ps.protocol = *protocol
	ps.cores = *cores
	ps.threads = *threads
	ps.timeOut = time.Duration(*timeout) * time.Second
	ipAddr, err := net.ResolveIPAddr("ip", flag.Arg(0))
	if err != nil {
		return err
	}
	ps.addr = ipAddr.String()
	if ps.threads > 0 {
		if ps.cores > 0 {
			if ps.cores <= runtime.NumCPU() {
				runtime.GOMAXPROCS(ps.cores)
			} else {
				runtime.GOMAXPROCS(runtime.NumCPU())
			}
		} else {
			runtime.GOMAXPROCS(1)
		}
	} else {
		ps.threads = 1
	}
	return nil
}

func (ps *PortScan) stdout(port int) {
	fmt.Println(port)
}
