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
	"errors"
	"io"

	"github.com/drnp/slater/slater/engine"
	"github.com/drnp/slater/slater/runtime/utils"
)

// SlaterWorker : Client worker of network server
type SlaterWorker struct {
	Addr       string
	UID        uint64
	conn       io.ReadWriteCloser
	recvBuffer *bytes.Buffer
	sendBuffer *bytes.Buffer
	recvChan   chan struct{}
	sendChan   chan int
	closeChan  chan struct{}
	server     *SlaterServer
}

// Drive : Start worker
/* {{{ [Driver] worker handler */
func (worker *SlaterWorker) Drive() error {
	if worker == nil {
		return errors.New("Invalid worker object")
	}

	// Reader
	go func() {
		// Read data from socket
		var err error
		var n int
		var msg *engine.Message
		logger := utils.NewLogger("SLATER: ")

	loop:
		for {
			buf := make([]byte, 4096)
			n, err = worker.conn.Read(buf)
			if err != nil {
				if err != io.EOF {
					// Read error
					logger.Println("Socket read error")
				} else {
					// Worker closed
					if worker.server != nil && worker.server.OnClose != nil {
						worker.server.OnClose(worker)
					}

					close(worker.closeChan)
					break loop
				}
			} else {
				n, err = worker.recvBuffer.Write(buf[:n])
				if err != nil {
					// Write buffer error
					logger.Println("Buffer write error")
				} else {
					if worker.server != nil {
					TryMsg:
						for {
							if msg == nil {
								msg = engine.NewMessage(worker.recvBuffer)
							}

							ret, _ := msg.Parse()
							if ret {
								if engine.MsgTypePing == msg.Type {
									// Ping - Pong
									//logger.Println("Access ping")
									pong := engine.NewMessage(nil)
									pong.Type = engine.MsgTypePong
									worker.WriteMessage(pong)
								} else if worker.server.OnMessage != nil {
									worker.server.OnMessage(worker, msg)
								}
								msg = nil
							} else {
								break TryMsg
							}
						}
					} else {
						// Clear all
						worker.recvBuffer.Reset()
					}
				}
			}
		}
	}()

	// Writer
	go func() {
		var (
			err   error
			nData int
			nSent int
			n     int
		)
		//logger := utils.NewLogger("SLATER: ")
		buf := make([]byte, 4096)

	loop:
		for {
			<-worker.sendChan
			if worker.sendBuffer.Len() > 0 {
				nData, _ = worker.sendBuffer.Read(buf)
				nSent = 0
				// Write out
				for nSent < nData {
					n, err = worker.conn.Write(buf[:nData])
					if err != nil {
						if err == io.EOF {
							// Nothing to read
							break loop
						} else {
							// Error
						}
					} else {
						// Data sent
						nSent += n
					}
				}
			}
		}
	}()

loop:
	for {
		select {
		//case <-worker.sendChan:
		case <-worker.closeChan:
			break loop
		}
	}

	return nil
}

/* }}} */

// WriteRaw : Send data from worker
/* {{{ [Write] Send data */
func (worker *SlaterWorker) WriteRaw(data []byte) error {
	if worker == nil {
		return errors.New("Invalid worker object")
	}

	size, err := worker.sendBuffer.Write(data)
	if err != nil {
		return err
	}

	worker.sendChan <- size

	return nil
}

/* }}} */

// WriteMessage : Send engine.CommonCommand
/* {{{ [WriteMessage] Send command */
func (worker *SlaterWorker) WriteMessage(msg *engine.Message) error {
	if msg == nil {
		return errors.New("Invalid message object")
	}

	data, _ := msg.Stream()
	//fmt.Printf("Sending :\n%#v\n", data)
	//utils.DebugByteArray(data)
	size, err := worker.sendBuffer.Write(data)
	if err != nil {
		//fmt.Println(err.Error())
		return err
	}

	worker.sendChan <- size

	return nil
}

/* }}} */

// ReadAll : Read all data from worker
/* {{{ [ReadAll] */
func (worker *SlaterWorker) ReadAll() ([]byte, error) {
	if worker == nil {
		return nil, errors.New("Invalid worker object")
	}

	n := worker.recvBuffer.Len()
	ret := make([]byte, n)
	var err error
	n, err = worker.recvBuffer.Read(ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

/* }}} */

// AccessRequest : Default server handler
// Read request from client and return response
/* {{{ [AccessRequest] */
func AccessRequest(clientAddr string, request interface{}) (response interface{}) {
	return nil
}

/* }}} */

// clientHandler : Network client handler
/* {{{ [clientHandler] */
func clientHandler(server *SlaterServer, conn io.ReadWriteCloser, clientAddr string) {

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
