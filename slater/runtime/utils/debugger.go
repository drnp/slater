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

package utils

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/drnp/slater/slater/engine"
)

// DebugBuffer : Output content of buffer in hex mode
/* {{{ [DebugBuffer] Out put buffer */
func DebugBuffer(buf *bytes.Buffer) error {
	if buf == nil {
		return errors.New("Invalid buffer object")
	}

	raw := buf.Bytes()

	return DebugByteArray(raw)
}

/* }}} */

// DebugByteArray : Output content of byte array in hex mode
/* {{{ [DebugByteArray] Output byte array */
func DebugByteArray(ba []byte) error {
	if ba == nil {
		return errors.New("Invalid byte array")
	}

	fmt.Printf("=== Debug byte array [%d bytes] ===\n", len(ba))
	for idx := 0; idx < len(ba); idx++ {
		fmt.Printf("%02X ", ba[idx])
		if 7 == idx%8 {
			if 15 == idx%16 {
				fmt.Printf("\n")
			} else {
				fmt.Printf("  ")
			}
		}
	}

	fmt.Printf("\n\n")

	return nil
}

/* }}} */

// DebugMessage : Output detail of engine.Message
/* {{{ [DebugMessage] Output message */
func DebugMessage(msg *engine.Message) error {
	if msg == nil {
		return errors.New("Invalid message object")
	}

	fmt.Printf("=== Debug message ===\n")
	fmt.Printf("Message type : %d\n", msg.Type)
	fmt.Printf("Message compress mode : %d\n", msg.CompressMode)
	fmt.Printf("Message serialize mode : %d\n", msg.SerializeMode)
	fmt.Printf("Body : \n")

	DebugBody(&msg.Body)

	return nil
}

/* }}} */

// DebugBody : Output message body
/* {{{ [DebugBody] Output body */
func DebugBody(body *engine.Body) error {
	if body == nil {
		return errors.New("Invalid message body")
	}

	fmt.Printf("=== Debug message body ===\n")
	fmt.Printf("App : %s\n", body.App)
	fmt.Printf("Uid : %v\n", body.UID)
	fmt.Printf("Payload : \n")

	DebugByteArray(body.Payload)

	return nil
}

/* }}} */

// DebugCommand : Output engine.CommonCommand
/* {{{ [DebugCommand] Output command */
func DebugCommand(cmd *engine.CommonCommand) error {
	if cmd == nil {
		return errors.New("Invalid command object")
	}

	fmt.Printf("=== Debug CommonCommand ===\n")
	fmt.Printf("Command : %d\n", cmd.Command)
	fmt.Printf("Params : %v\n", cmd.Params)
	fmt.Printf("Additional : %v\n", cmd.Additional)

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
