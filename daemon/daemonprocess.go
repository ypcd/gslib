package daemon

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"sync/atomic"
	"time"

	"github.com/ypcd/gslib/gserror"
)

func checkErr(err error) {
	gserror.CheckErrorEx_panic(err)
}

// 阻塞调用。
// 创建一个守护进程和一个子进程。
// 守护进程：	守护进程监视子进程的运行，当子进程退出后，守护进程创建一个新的子进程。
// 子进程：		子进程运行参数f函数。
func Daemon(f func()) {
	if len(os.Args) > 1 && os.Args[1] == "child" {
		// 子进程执行的代码
		childProcess_func(f, os.Args[2])
		return
	}
	serverAddr, msg1 := rpc_server()
	//defer g_exit.Store(true)
	defer msg1.setExit(true)

	for {
		msg1.setExit(false)
		// Fork出子进程
		cmd := exec.Command(os.Args[0], "child", serverAddr)

		pwd, err := os.Getwd()
		checkErr(err)
		fmt.Println("pwd:", pwd)
		outf, err := os.Create("child.log")
		checkErr(err)
		defer outf.Close()
		errf, err := os.Create("child.err.log")
		checkErr(err)
		defer errf.Close()
		cmd.Stdout = outf
		cmd.Stderr = errf

		// 父进程执行的代码

		// 启动子进程
		err = cmd.Start()
		if err != nil {
			fmt.Println("启动子进程失败:", err)
			return
		}

		fmt.Println("父进程启动子进程，子进程PID:", cmd.Process.Pid)

		// 等待子进程退出
		err = cmd.Wait()
		//checkErr(err)
		fmt.Println("cmd.Wait():err:", err)
		msg1.setExit(true)
	}
}

func childProcess_wait_close(serverAddr string) {
	fmt.Println("子进程启动")

	fmt.Println("Server addr:", serverAddr)
	client, err := rpc.DialHTTP("tcp", serverAddr)
	checkErr(err)
	var reply bool

	for {
		// 子进程每秒输出一条日志
		err := client.Call("Msg.Get", 1, &reply)
		checkErr(err)
		if reply {
			os.Exit(1)
		}
		time.Sleep(time.Second)
	}
}

func childProcess_func(f func(), serverAddr string) {
	fmt.Println("子进程启动")
	go childProcess_wait_close(serverAddr)
	f()
}

type Msg struct {
	exit atomic.Bool
}

func (m *Msg) Get(v1 int, v2 *bool) error {
	*v2 = m.exit.Load()
	return nil
}

func (m *Msg) setExit(v bool) {
	m.exit.Store(v)
}

func rpc_server() (listenAddr string, msg *Msg) {
	msg1 := &Msg{}

	rpc.Register(msg1)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", "127.0.0.1:")
	checkErr(err)
	go http.Serve(l, nil)
	fmt.Println("Listen addr:", l.Addr().String())

	return l.Addr().String(), msg1
}
