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
	"bufio"
	"bytes"
	"encoding/binary"

	"github.com/ugorji/go/codec"
)

const (
	// MsgTypeReserved : Reserved
	MsgTypeReserved byte = iota
	// MsgTypeOnline : Online message
	MsgTypeOnline
	// MsgTypeOnlineAck : Online ACK
	MsgTypeOnlineAck
	// MsgTypeOffline : Offline message
	MsgTypeOffline
	// MsgTypeOfflineAck : Offline ACK
	MsgTypeOfflineAck
	// MsgTypeUpward : Upward message
	MsgTypeUpward
	// MsgTypeUpwardAck : Upward ACK
	MsgTypeUpwardAck
	// MsgTypeDownward : Downward message
	MsgTypeDownward
	// MsgTypeDownwardAck : Downward ACK
	MsgTypeDownwardAck
	// MsgTypePing : Heartbeat PING
	MsgTypePing
	// MsgTypePong : Heartbeat PONE
	MsgTypePong
)

/*
const (
	// MsgStageHeader : Header needed (Type + CompressMode + SerializeMode) (1 byte)
	MsgStageHeader int = iota
	// MsgStageUids : Uids needed (2 bytes)
	MsgStageUids
	// MsgStageUID : UID array needed ((8 * Uids) bytes)
	MsgStageUID
	// MsgStageLength : Payload length needed (4 bytes)
	MsgStageLength
	// MsgStagePayload : Payload data needed (length bytes)
	MsgStagePayload
	// MsgStageComplete : All done
	MsgStageComplete
)
*/
const (
	// MsgStageHeader : Header needed (Type + CompressMode + SerializeMode) + (BodyLength) (5 bytes)
	MsgStageHeader byte = iota
	// MsgStageBody : Body needed (n bytes)
	MsgStageBody
	// MsgStageComplete : All done
	MsgStageComplete
)

const (
	// MsgCompressNone : No compression
	MsgCompressNone byte = iota
	// MsgCompressDeflate : Deflate (GNU zip)
	MsgCompressDeflate
	// MsgCompressSnappy : Google snappy
	MsgCompressSnappy
	// MsgCompressLZ4 : LZ4
	MsgCompressLZ4
)

const (
	// MsgSerializeRaw : No serialization
	MsgSerializeRaw byte = iota
	// MsgSerializeJSON : JSON
	MsgSerializeJSON
	// MsgSerializeMsgPack : MessagePack
	MsgSerializeMsgPack
	// MsgSerializeAMF3 : AMF3
	MsgSerializeAMF3
)

// Body : Message body
type Body struct {
	App     string
	UID     []int64
	Payload []byte
}

// Message : Data struct defination
type Message struct {
	Type          byte
	SerializeMode byte
	CompressMode  byte
	BodyLength    uint32
	Body          Body
	buffer        *bytes.Buffer
	Stage         byte
}

// NewBody : Create a new body
/* {{{ [NewBody] */
func NewBody() (body *Body) {
	return &Body{}
}

/* }}} */

// NewMessage : Create a new message
/* {{{ [NewMessage] */
func NewMessage(buf *bytes.Buffer) (msg *Message) {
	return &Message{
		Type:          0,
		SerializeMode: MsgSerializeMsgPack,
		CompressMode:  MsgCompressNone,
		BodyLength:    0,
		buffer:        buf,
		Stage:         MsgStageHeader,
	}
}

/* }}} */

// Parse : Try parse
/* {{{ [Parse] Try parse message */
func (msg *Message) Parse() (bool, error) {
	if msg.buffer == nil {
		return false, nil
	}

	var remaining uint32
	var err error
	enough := false
	for {
		remaining = uint32(msg.buffer.Len())
		enough = true
		switch msg.Stage {
		case MsgStageHeader:
			if remaining < 5 {
				enough = false
			} else {
				header, _ := msg.buffer.ReadByte()
				msg.Type = header >> 4
				msg.SerializeMode = (header >> 2) & 3
				msg.CompressMode = header & 3
				var length uint32
				var n int
				buf := make([]byte, 4)
				n, err = msg.buffer.Read(buf)
				if 4 == n {
					r := bytes.NewReader(buf)
					binary.Read(r, binary.BigEndian, &length)
					msg.BodyLength = length
				}
				if MsgTypePing == msg.Type {
					// Force Zero
					msg.BodyLength = 0
				}

				msg.Stage++
			}
			break
		case MsgStageBody:
			if remaining < msg.BodyLength {
				enough = false
			} else {
				if msg.BodyLength > 0 {
					raw := make([]byte, msg.BodyLength)
					msg.buffer.Read(raw)
					var hdl codec.MsgpackHandle
					dec := codec.NewDecoderBytes(raw, &hdl)
					err = dec.Decode(&msg.Body)
				}

				msg.Stage++
			}
			break
		case MsgStageComplete:
		default:
			// All done
			break
		}

		if MsgStageComplete == msg.Stage {
			return true, err
		}

		if false == enough {
			return false, err
		}
	}
}

/* }}} */

// Stream : Build bytes
/* {{{ [Stream] Serialize message to bytes */
func (msg *Message) Stream() ([]byte, error) {
	var buf bytes.Buffer
	var err error
	w := bufio.NewWriter(&buf)

	// Ignore stage now
	// Write header
	header := ((msg.Type & 15) << 4) | ((msg.SerializeMode & 3) << 2) | (msg.CompressMode & 3)
	buf.WriteByte(byte(header))
	switch msg.Type {
	case MsgTypeDownward:
		// Pack body
		var raw []byte
		var hdl codec.MsgpackHandle
		hdl.EncodeOptions.StructToArray = true
		enc := codec.NewEncoderBytes(&raw, &hdl)
		err = enc.Encode(msg.Body)

		// Write length
		binary.Write(w, binary.BigEndian, uint32(len(raw)))
		w.Flush()
		buf.Write(raw)
		break
	case MsgTypePong:
		buf.Write([]byte{0, 0, 0, 0})
		break
	case MsgTypeUpwardAck:
		buf.Write([]byte{0, 0, 0, 0})
		break
	default:
		// Nothing to do
		break
	}

	return buf.Bytes(), err
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
