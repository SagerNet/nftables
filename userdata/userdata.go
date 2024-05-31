// Copyright 2018 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package userdata implements a TLV parser/serializer for libnftables-compatible comments
package userdata

import (
	"encoding/binary"
)

type Type byte

// TLV type values are defined in:
// https://git.netfilter.org/iptables/tree/iptables/nft.c?id=73611d5582e72367a698faf1b5301c836e981465#n1659
const (
	TypeComment Type = iota
	TypeEbtablesPolicy

	TypesCount
)

func Append(udata []byte, typ Type, data []byte) []byte {
	udata = append(udata, byte(typ), byte(len(data)))
	udata = append(udata, data...)

	return udata
}

func Get(udata []byte, styp Type) []byte {
	for {
		if len(udata) < 2 {
			break
		}

		typ := Type(udata[0])
		length := int(udata[1])
		data := udata[2 : 2+length]

		if styp == typ {
			return data
		}

		if len(udata) < 2+length {
			break
		} else {
			udata = udata[2+length:]
		}
	}

	return nil
}

func AppendUint32(udata []byte, typ Type, num uint32) []byte {
	var data [4]byte
	binary.LittleEndian.PutUint32(data[:], num)

	return Append(udata, typ, data[:])
}

func GetUint32(udata []byte, typ Type) (uint32, bool) {
	data := Get(udata, typ)
	if data == nil {
		return 0, false
	}

	return binary.LittleEndian.Uint32(data), true
}

func AppendString(udata []byte, typ Type, str string) []byte {
	data := append([]byte(str), 0)
	return Append(udata, typ, data)
}

func GetString(udata []byte, typ Type) (string, bool) {
	data := Get(udata, typ)
	if data == nil {
		return "", false
	}

	if data[len(data)-1] == 0 {
		data = data[:len(data)-1]
	}

	return string(data), true
}
