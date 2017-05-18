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
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/drnp/slater/slater/engine"
	"github.com/drnp/slater/slater/runtime/utils"
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
func NewWorker(server *SlaterServer, addr string, conn io.ReadWriteCloser) *SlaterWorker {
	return &SlaterWorker{
		Addr:       addr,
		UID:        0,
		conn:       conn,
		recvBuffer: bytes.NewBuffer(nil),
		sendBuffer: bytes.NewBuffer(nil),
		recvChan:   make(chan struct{}),
		sendChan:   make(chan int),
		closeChan:  make(chan struct{}),
		server:     server,
	}
}

/* }}} */

// SendMessage : Send command (message) to remote
/* {{{ [SendCommand] Send message */
func SendMessage(msg *engine.Message) {

}

/* }}} */

// DefaultOnConnect : Default behavior
func DefaultOnConnect(worker *SlaterWorker) error {
	if worker == nil {
		return errors.New("Invalid worker object")
	}

	fmt.Printf("Client %s connected\n", worker.Addr)

	return nil
}

// DefaultOnClose : Default behavior
func DefaultOnClose(worker *SlaterWorker) error {
	if worker == nil {
		return errors.New("Invalid worker object")
	}

	fmt.Printf("Client %s disconnected\n", worker.Addr)

	return nil
}

// DefaultOnData : Default behavior
func DefaultOnData(worker *SlaterWorker) (int, error) {
	if worker == nil {
		return 0, errors.New("Invalid worker object")
	}

	data, err := worker.ReadAll()
	len := binary.Size(data)
	fmt.Printf("Read %d bytes from client %s\n", len, worker.Addr)

	return len, err
}

// DefaultOnnMessage : Default behavior
func DefaultOnnMessage(worker *SlaterWorker, msg *engine.Message) error {
	if msg == nil {
		return errors.New("Invalid message object")
	}

	utils.DebugMessage(msg)
	if msg.Body.Payload != nil {
		cmd, _ := engine.CmdDecode(msg.Body.Payload, msg.SerializeMode)
		if cmd != nil {
			utils.DebugCommand(cmd)
		}
	} else {
		if engine.MsgTypePing == msg.Type {
			downmsg := engine.NewMessage(nil)
			downmsg.Type = engine.MsgTypePong
			worker.WriteMessage(downmsg)
		}
	}

	return nil
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
