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
	"errors"
	"fmt"
	"io"
	"sync"
)

// HandlerFunc : Server handler function
type HandlerFunc func(clientAddr string, request interface{}) (response interface{})
type onConnectHandler func(worker *SlaterWorker) error
type onCloseHandler func(worker *SlaterWorker) error
type onDataHandler func(worker *SlaterWorker) (int, error)

// SlaterServer : Main server struct
type SlaterServer struct {
	// Addr : Address to listen to
	Addr string

	// Handler : Handler function for incoming request
	Handler HandlerFunc

	// nClients : Amount of clients
	nClients int

	// SendBufferSize : Size of send buffer
	// Default value is 0
	SendBufferSize int

	// RecvBufferSize : Size of recieve buffer
	// Default value is 0
	RecvBufferSize int

	// Listener : Socket listener of server
	Listener SlaterListener

	// stopChan : Send stop signal to server
	stopChan chan struct{}

	// Waiter : Symc waiter
	Waiter *sync.WaitGroup

	// Hooks
	OnConnect onConnectHandler
	OnClose   onCloseHandler
	OnData    onDataHandler
}

// serve : Go server
/* {{{ [serve] */
func (server *SlaterServer) serve() error {
	var err error
	if server == nil {
		return errors.New("Invalid server object")
	}

	if server.Handler == nil {
		return errors.New("Server handler cannot be nil")
	}

	if server.stopChan != nil {
		return errors.New("Server already running")
	}

	if server.Listener == nil {
		server.Listener = &defaultListener{}
	}

	err = server.Listener.Init(server.Addr)
	if err != nil {
		err = fmt.Errorf("Cannot listen to address [%s] : %s", server.Addr, err)
		return err
	}

	server.stopChan = make(chan struct{})
	var conn io.ReadWriteCloser
	var clientAddr string
	var worker *SlaterWorker

	// Let's go
	go func() {
		for {
			acceptChan := make(chan struct{})
			go func() {
				if conn, clientAddr, err = server.Listener.Accept(); err != nil {
					// Accept error
					fmt.Printf("Accept error\n")
					panic(err)
				}
				close(acceptChan)
			}()

			select {
			case <-server.stopChan:
				// Server close
				server.Listener.Close()
				server.Waiter.Done()
				<-acceptChan
			case <-acceptChan:
				server.nClients++
			}

			if err != nil {
				continue
			}

			worker = NewWorker(server, clientAddr, conn)
			go worker.Drive()
			if server.OnConnect != nil {
				server.OnConnect(worker)
			}
		}
	}()

	return nil
}

/* }}} */

// Start : Network server startup
/* {{{ [Start] Start server */
func (server *SlaterServer) Start() error {
	if server == nil {
		return errors.New("Invalid server object")
	}

	if server.Waiter == nil {
		return errors.New("Server has no sync waiter")
	}

	err := server.serve()
	if err != nil {
		return err
	}

	server.Waiter.Add(1)

	return nil
}

/* }}} */

// Stop : Stop network server
/* {{{ [Stop] Stop server */
func (server *SlaterServer) Stop() error {
	if server == nil {
		return errors.New("Invalid server object")
	}

	if server.stopChan == nil {
		return errors.New("Server not running")
	}

	close(server.stopChan)
	server.stopChan = nil
	if server.Waiter != nil {
		server.Waiter.Done()
	}

	return nil
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
