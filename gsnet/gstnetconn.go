package gsnet

import (
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type GSTNetConn struct {
	conn                                                        net.Conn
	rlent, wlent                                                atomic.Uint64
	ntTimeOut                                                   time.Duration
	setTimeOut, setTimeOutRead, setTimeOutWirite                float64 //seconds
	closed                                                      atomic.Bool
	localAddr, remoteAddr                                       string
	lock_setTimeOut, lock_setTimeOutRead, lock_setTimeOutWirite sync.Mutex
}

// blocking
func NewGSTNetConn(conn net.Conn) *GSTNetConn {
	return &GSTNetConn{
		conn:       conn,
		ntTimeOut:  0,
		localAddr:  conn.LocalAddr().String(),
		remoteAddr: conn.RemoteAddr().String(),
	}

}

// non-blocking
func NewGSTNetConnNonBlcoking(conn net.Conn, timeout time.Duration) *GSTNetConn {
	return &GSTNetConn{conn: conn, ntTimeOut: timeout}
}

func (c *GSTNetConn) Read(buf []byte) (n int, err error) {
	var err2 error
	if c.ntTimeOut != 0 {
		err2 = c.conn.SetReadDeadline(time.Now().Add(c.ntTimeOut))
	} else {
		//err2 = c.conn.SetReadDeadline(time.Time{})
	}

	if err != nil {
		return 0, err2
	}
	rn, err := c.conn.Read(buf)
	c.rlent.Add(uint64(rn))
	return rn, err
}

func (c *GSTNetConn) Write(buf []byte) (n int, err error) {
	var err2 error
	if c.ntTimeOut != 0 {
		err2 = c.conn.SetReadDeadline(time.Now().Add(c.ntTimeOut))
	} else {
		//err2 = c.conn.SetReadDeadline(time.Time{})
	}

	if err != nil {
		return 0, err2
	}
	rn, err := c.conn.Write(buf)
	c.wlent.Add(uint64(rn))
	return rn, err
}

func (c *GSTNetConn) Close() error {
	c.closed.Store(true)
	return c.conn.Close()
}
func (c *GSTNetConn) RemoteAddr() net.Addr { return c.conn.RemoteAddr() }
func (c *GSTNetConn) LocalAddr() net.Addr  { return c.conn.LocalAddr() }

func (c *GSTNetConn) SetDeadline(t time.Time) error {
	c.lock_setTimeOut.Lock()
	defer c.lock_setTimeOut.Unlock()
	c.setTimeOut = t.Sub(time.Now()).Seconds()
	return c.conn.SetDeadline(t)
}

func (c *GSTNetConn) SetReadDeadline(t time.Time) error {
	c.lock_setTimeOutRead.Lock()
	defer c.lock_setTimeOutRead.Unlock()
	c.setTimeOutRead = t.Sub(time.Now()).Seconds()
	return c.conn.SetReadDeadline(t)
}
func (c *GSTNetConn) SetWriteDeadline(t time.Time) error {
	c.lock_setTimeOutWirite.Lock()
	defer c.lock_setTimeOutWirite.Unlock()
	c.setTimeOutWirite = t.Sub(time.Now()).Seconds()
	return c.conn.SetWriteDeadline(t)
}
