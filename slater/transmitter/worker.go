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
)

type dataMsg struct {
	header uint8
	length uint32
	msg    []byte
}
type accessMsg struct {
	header byte
	user   uint
	length uint32
	body   dataMsg
}

type slaterMsg struct {
	header byte
	nUsers uint16
	users  []uint64
	length uint32
	body   dataMsg
}

// SlaterWorker : Client worker of network server
type SlaterWorker struct {
	clientAddr string
	recvBuffer *bytes.Buffer
	sendBuffer *bytes.Buffer
	recvChan   chan struct{}
	sendChan   chan struct{}
}

// Drive : Start worker
/* {{{ [Driver] */
func (worker *SlaterWorker) Drive() error {
	if worker == nil {
		return errors.New("Invalid worker object")
	}

	go func() {
		for {
			select {
			case <-worker.sendChan:
			}
		}
	}()

	return nil
}

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

//
/* }}} */

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
