/*
 * Copyright (c) 2016, 2017
 *     PC-Game of Qihu.360. All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 * 3. Neither the name of the University nor the names of its contributors
 *    may be used to endorse or promote products derived from this software
 *    without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE REGENTS AND CONTRIBUTORS ``AS IS'' AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED.  IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
 * OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
 * LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
 * OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
 * SUCH DAMAGE.
 */

package transmitter

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

// SlaterListener : Listener interface of slater engine
/* {{{ [SlaterListener] Listener interface */
type SlaterListener interface {
	// Init : Called on server start
	Init(addr string) error

	// Accept : Return incoming connection from client.
	// We support TCP / UNIX / TLS here
	Accept() (conn io.ReadWriteCloser, clientAddr string, err error)

	// Close : Close the listener
	// Server shutdown
	Close() error

	// ListenAddr : Listener's network address
	ListenAddr() net.Addr
}

/* }}} */

// defaultListener : TCPServer listener
/* {{{ [defaultListener] */
type defaultListener struct {
	L net.Listener
}

func (dl *defaultListener) Init(addr string) (err error) {
	fmt.Printf("Listen on addr : %s\n", addr)
	dl.L, err = net.Listen("tcp", addr)

	return
}

func (dl *defaultListener) Accept() (conn io.ReadWriteCloser, clientAddr string, err error) {
	c, err := dl.L.Accept()
	if err != nil {
		return nil, "", err
	}

	return c, c.RemoteAddr().String(), nil
}

func (dl *defaultListener) Close() error {
	return dl.L.Close()
}

func (dl *defaultListener) ListenAddr() net.Addr {
	if dl.L != nil {
		return dl.L.Addr()
	}

	return nil
}

/* }}} */

// netListener : General listener
/* {{{ [netListener] */
type netListener struct {
	F func(addr string) (net.Listener, error)
	L net.Listener
}

/* }}} */

// NewTCPServer : Create a new TCP server
/* {{{ [NewTCPServer] Create TCP server */
func NewTCPServer(addr string, handler HandlerFunc) *SlaterServer {
	return &SlaterServer{
		Addr:     addr,
		Handler:  handler,
		Listener: &defaultListener{},
	}
}

/* }}} */

// NewWorker : Create a new worker
/* {{{ [NewWorker] Create worker */
func NewWorker(addr string, conn io.ReadWriteCloser) *SlaterWorker {
	return &SlaterWorker{
		Addr:       addr,
		conn:       conn,
		recvBuffer: bytes.NewBuffer(nil),
		sendBuffer: bytes.NewBuffer(nil),
		recvChan:   make(chan struct{}),
		sendChan:   make(chan int),
		closeChan:  make(chan struct{}),
	}
}

/* }}} */

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
