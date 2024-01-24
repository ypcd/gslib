package gsnet

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ypcd/gslib/gserror"
	"github.com/ypcd/gslib/gsrand"
)

func GetIP(conn net.Addr) string {
	s1 := conn.String()
	ix := strings.Index(s1, ":")
	if ix != -1 {
		return s1[:ix]
	}
	return s1
}

func NewHttpServer(listenAddr string) string {
	addchan := make(chan string)

	go func() {
		// 使用net.Listen随机选择一个可用的端口
		listener, err := net.Listen("tcp", listenAddr)
		gserror.CheckError_panic(err)
		defer listener.Close()

		// 获取实际监听的地址
		addchan <- listener.Addr().String()

		// 创建HTTP服务器
		server := &http.Server{
			Handler: nil,
		}

		// 启动HTTP服务器
		err = server.Serve(listener)
		gserror.CheckError_panic(err)
	}()
	return <-addchan
}

func GetNetLocalRDPort() string {
	return fmt.Sprintf("127.0.0.1:%d", gsrand.GetRD_netPortNumber())
}

func setTcpConnKeepAlive(conn net.Conn) {
	tcpconn1, ok := conn.(*net.TCPConn)
	if !ok {
		panic("!ok")
	}
	tcpconn1.SetKeepAlive(true)
	tcpconn1.SetKeepAlivePeriod(time.Second)
}

func NewNetConnRD(ip string) (net.Conn, net.Conn) {
	listenAddr := ip + ":"
	var conna, connc net.Conn

	serverAddr := ""

	wg := sync.WaitGroup{}
	listenDone := make(chan int)

	wg.Add(1)
	go func() {
		defer wg.Done()

		lst, err := net.Listen("tcp4", listenAddr)
		gserror.CheckError_panic(err)
		serverAddr = lst.Addr().String()

		listenDone <- 1
		//defer lst.Close()
		conna, err = lst.Accept()
		gserror.CheckError_exit(err)
	}()
	<-listenDone
	time.Sleep(time.Millisecond * 10)
	connc, err := net.Dial("tcp4", serverAddr)
	gserror.CheckError_exit(err)
	wg.Wait()

	setTcpConnKeepAlive(conna)
	setTcpConnKeepAlive(connc)

	return NewGSTNetConn(conna), NewGSTNetConn(connc)
}

func NewNetConnRDLocal() (net.Conn, net.Conn) {
	return NewNetConnRD("127.0.0.1")
}

func NewNetConnRDListen() (net.Conn, net.Conn, net.Listener) {
	var conna, connc net.Conn

	serverAddr := ""

	wg := sync.WaitGroup{}
	listenDone := make(chan int)
	var lst net.Listener
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()

		lst, err = net.Listen("tcp4", "127.0.0.1:")
		gserror.CheckError_panic(err)
		serverAddr = lst.Addr().String()

		listenDone <- 1
		//defer lst.Close()
		conna, err = lst.Accept()
		gserror.CheckError_exit(err)
	}()
	<-listenDone
	time.Sleep(time.Millisecond * 10)
	connc, err = net.Dial("tcp4", serverAddr)
	gserror.CheckError_exit(err)
	wg.Wait()

	setTcpConnKeepAlive(conna)
	setTcpConnKeepAlive(connc)

	return NewGSTNetConn(conna), NewGSTNetConn(connc), lst
}

func GetNetConnAddrString(str string, conn net.Conn) string {
	return fmt.Sprintf("%s: [localIp:%s  remoteIp:%s]\n",
		str, conn.LocalAddr().String(), conn.RemoteAddr().String())
}
