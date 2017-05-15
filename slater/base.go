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

package slater

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"runtime"
	"sync"

	"github.com/drnp/slater/slater/engine"
	"github.com/drnp/slater/slater/runtime/config"
	"github.com/drnp/slater/slater/runtime/utils"
	"github.com/drnp/slater/slater/transmitter"
)

// Global waiter
var globalWaiter sync.WaitGroup

// Logger to stdout
var logger *log.Logger

// Conf : Iniial configuration
type Conf struct {
	CustomConf map[string]interface{}
	Game       string
	Standalone bool
}

func onConnectFunc(worker *transmitter.SlaterWorker) error {
	if worker == nil {
		return errors.New("Invalid worker object")
	}

	fmt.Printf("Client %s connected\n", worker.Addr)

	return nil
}

func onCloseFunc(worker *transmitter.SlaterWorker) error {
	if worker == nil {
		return errors.New("Invalid worker object")
	}

	fmt.Printf("Client %s disconnected\n", worker.Addr)

	return nil
}

func onDataFunc(worker *transmitter.SlaterWorker) (int, error) {
	if worker == nil {
		return 0, errors.New("Invalid worker object")
	}

	data, err := worker.ReadAll()
	len := binary.Size(data)
	fmt.Printf("Read %d bytes from client %s\n", len, worker.Addr)

	return len, err
}

func onMessageFunc(worker *transmitter.SlaterWorker, msg *engine.Message) error {
	if msg == nil {
		return errors.New("Invalid message object")
	}

	utils.DebugMessage(msg)
	downmsg := engine.NewMessage(nil)
	if msg.Payload != nil {
		cmd, _ := engine.CmdDecode(msg.Payload, msg.SerializeMode)
		if cmd != nil {
			utils.DebugCommand(cmd)
		}

		downmsg.Type = engine.MsgTypeDownward
		downmsg.SerializeMode = msg.SerializeMode
		downmsg.UID = []int64{100, 200, 300, 400}
		downmsg.Payload, _ = engine.CmdEncode(cmd, downmsg.SerializeMode)
		worker.WriteMessage(downmsg)
	} else {
		if engine.MsgTypePing == msg.Type {
			fmt.Println("Heartbeat !")
			downmsg.Type = engine.MsgTypePong
			worker.WriteMessage(downmsg)
		}
	}

	return nil
}

// Start : Slater startup
/* {{{ [slater.Start] Startup */
func Start(c *Conf) (err error) {
	if nil == c {
		return fmt.Errorf("No valid configuration")
	}

	// Runtime
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Default values
	if nil != c.CustomConf {
		for key, value := range c.CustomConf {
			config.SetDefault(key, value)
		}
	}

	// Override config
	config.Load()

	// Start
	//fmt.Printf("Start slater engine with game <%s> ...\n", c.Game)

	// Engine
	engine.Start(logger)

	// TCPServer
	if !c.Standalone {
		s := transmitter.NewTCPServer(config.Get("server_addr").(string), transmitter.AccessRequest)
		s.Waiter = &globalWaiter
		s.OnConnect = onConnectFunc
		s.OnClose = onCloseFunc
		s.OnData = onDataFunc
		s.OnMessage = onMessageFunc
		err = s.Start()
		if err != nil {
			//return err
			panic(err)
		}
	}

	globalWaiter.Wait()

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
