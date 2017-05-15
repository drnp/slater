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

package engine

import (
	"errors"

	"github.com/ugorji/go/codec"
)

// CommonCommand : Common command
type CommonCommand struct {
	Additional map[string]string      `cmd:"Additional"`
	Command    int                    `cmd:"Command"`
	Params     map[string]interface{} `cmd:"Params"`
}

// CmdEncode : Encode command struct into bytes
/* {{{ [CmdEncode] Encode command */
func CmdEncode(cmd interface{}, t byte) ([]byte, error) {
	if cmd == nil {
		return nil, errors.New("Invalid command object")
	}

	var (
		ret []byte
		err error
	)

	switch t {
	case MsgSerializeMsgPack:
		var hdl codec.MsgpackHandle
		hdl.EncodeOptions.StructToArray = true
		enc := codec.NewEncoderBytes(&ret, &hdl)
		err = enc.Encode(cmd)
		break
	case MsgSerializeJSON:
		var hdl codec.JsonHandle
		enc := codec.NewEncoderBytes(&ret, &hdl)
		err = enc.Encode(cmd)
		break
	case MsgSerializeAMF3:
		break
	case MsgSerializeRaw:
	default:
		break
	}

	return ret, err
}

/* }}} */

// CmdDecode : Decode bytes into struct
/* {{{ [CmdDecode] */
func CmdDecode(raw []byte, t byte) (*CommonCommand, error) {
	if raw == nil {
		return nil, errors.New("Invalid stream")
	}

	var (
		ret CommonCommand
		err error
	)

	switch t {
	case MsgSerializeMsgPack:
		var hdl codec.MsgpackHandle
		dec := codec.NewDecoderBytes(raw, &hdl)
		err = dec.Decode(&ret)
		break
	case MsgSerializeJSON:
		var hdl codec.JsonHandle
		dec := codec.NewDecoderBytes(raw, &hdl)
		err = dec.Decode(&ret)
		break
	case MsgSerializeAMF3:
		break
	case MsgSerializeRaw:
	default:
		break
	}
	return &ret, err
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
