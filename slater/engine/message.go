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
	"fmt"
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

// Message : Data struct defination
type Message struct {
	Type          byte
	SerializeMode byte
	CompressMode  byte
	Uids          uint16
	UID           []int64
	PayloadLength uint32
	Payload       []byte
	buffer        *bytes.Buffer
	Stage         int
}

// NewMessage : Create a new message
/* {{{ [NewMessage] */
func NewMessage(buf *bytes.Buffer) (msg *Message) {
	return &Message{
		Type:          0,
		SerializeMode: MsgSerializeMsgPack,
		CompressMode:  MsgCompressNone,
		//Cmd:           0,
		buffer: buf,
		Stage:  MsgStageHeader,
	}
}

/* }}} */

// Parse : Try parse
/* {{{ [Parse] Try parse message */
func (msg *Message) Parse() bool {
	if msg.buffer == nil {
		return false
	}

	var remaining int
	enough := false
	for {
		remaining = msg.buffer.Len()
		enough = true
		switch msg.Stage {
		case MsgStageHeader:
			if remaining < 1 {
				enough = false
			} else {
				header, _ := msg.buffer.ReadByte()
				msg.Type = header >> 4
				msg.SerializeMode = (header >> 2) & 3
				msg.CompressMode = header & 3
				if MsgTypePing == msg.Type {
					msg.Stage = MsgStageComplete
					return true
				}

				msg.Stage++
			}
			break
		case MsgStageUids:
			if remaining < 2 {
				enough = false
			} else {
				b1, _ := msg.buffer.ReadByte()
				b2, _ := msg.buffer.ReadByte()
				msg.Uids = uint16(b1)<<8 | uint16(b2)
				msg.Stage++
			}
			break
		case MsgStageUID:
			if remaining < (8 * int(msg.Uids)) {
				enough = false
			} else {
				msg.UID = make([]int64, msg.Uids)
				var UID int64
				buf := make([]byte, 8)
				for idx := 0; idx < int(msg.Uids); idx++ {
					n, _ := msg.buffer.Read(buf)
					if 8 == n {
						r := bytes.NewReader(buf)
						binary.Read(r, binary.BigEndian, &UID)
						msg.UID[idx] = UID
					}
				}

				msg.Stage++
			}
			break
		case MsgStageLength:
			if remaining < 4 {
				enough = false
			} else {
				var length uint32
				buf := make([]byte, 4)
				n, _ := msg.buffer.Read(buf)
				if 4 == n {
					r := bytes.NewReader(buf)
					binary.Read(r, binary.BigEndian, &length)
					msg.PayloadLength = length
				}
				msg.Stage++
			}
			break
		case MsgStagePayload:
			if remaining < int(msg.PayloadLength) {
				enough = false
			} else {
				msg.Payload = make([]byte, msg.PayloadLength)
				msg.buffer.Read(msg.Payload)
				msg.Stage++
			}
			break
		case MsgStageComplete:
		default:
			// All done
			break
		}

		if MsgStageComplete == msg.Stage {
			return true
		}

		if false == enough {
			return false
		}
	}
}

/* }}} */

// Stream : Build bytes
/* {{{ [Stream] Serialize message to bytes */
func (msg *Message) Stream() ([]byte, error) {
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)

	// Ignore stage now
	// Write header
	header := ((msg.Type & 15) << 4) | ((msg.SerializeMode & 3) << 2) | (msg.CompressMode & 3)
	buf.WriteByte(byte(header))

	if MsgTypeDownward == msg.Type {
		// UIDs
		if msg.UID == nil {
			// No UID?
			fmt.Println("No UID")
			buf.WriteByte(0)
			buf.WriteByte(0)
		} else {
			fmt.Println("Has UID")
			uids := uint16(len(msg.UID))
			binary.Write(w, binary.BigEndian, uids)
			for idx := 0; idx < int(uids); idx++ {
				binary.Write(w, binary.BigEndian, msg.UID[idx])
			}

			w.Flush()
		}

		// Payload
		if msg.Payload == nil {
			fmt.Println("No payload")
			buf.Write([]byte{0, 0, 0, 0})
		} else {
			fmt.Println("Has payload")
			length := uint32(len(msg.Payload))
			binary.Write(w, binary.BigEndian, length)
			w.Flush()
			buf.Write(msg.Payload)
		}
	}

	return buf.Bytes(), nil
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
