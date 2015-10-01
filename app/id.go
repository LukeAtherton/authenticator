// Copyright (c) Luke Atherton 2015

package authenticator

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"code.google.com/p/go-uuid/uuid"
	//gouuid "github.com/nu7hatch/gouuid"
)

var (
	seq  uint32
	node = getNodeUint32()
)

type ID []byte

//implements TextMarshaler for json encoding
func (id ID) MarshalText() (text []byte, err error) {
	return []byte(id.String()), nil
}

//implements TextUnmarshaler for json encoding
func (id ID) UnmarshalText(text []byte) error {
	id = DecodeIdString(string(text))
	return nil
}

//implements JSONUnmarshaler for json encoding
func (id *ID) UnmarshalJSON(text []byte) (err error) {
	decoded := DecodeIdString(string(text))
	*id = decoded
	return
}

func NewUUID() ID {
	// Get the unique ID
	/*u4, err := gouuid.NewV4()
	if err != nil {
		panic(err)
	}*/

	id, _ := parse(uuid.NewRandom())

	return id
}

func getNodeUint32() uint32 {
	n := uuid.NodeID()
	return binary.BigEndian.Uint32(n)
}

// 8 bytes of UNIXNANO
// 4 bytes of counter
// 4 bytes of hardware address
//type UUID []byte

func parse(b []byte) (u ID, err error) {
	if len(b) != 16 {
		err = errors.New("Given slice is not valid UUID sequence")
		return
	}
	u = make([]byte, 16)
	copy(u[:], b)
	return
}

// Returns unparsed version of the generated UUID sequence.
func (u ID) String() string {
	uBytes := u.Bytes()
	if len(uBytes) != 16 {
		return ""
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", uBytes[0:4], uBytes[4:6], uBytes[6:8], uBytes[8:10], uBytes[10:])
}

func (u ID) Equals(u2 ID) bool {
	return bytes.Equal(u, u2)
}

func DecodeIdString(idString string) ID {
	g := strings.Replace(string(idString), "-", "", -1)
	g = strings.Replace(g, "\"", "", -1)
	b, err := hex.DecodeString(g)
	if err != nil {
		fmt.Printf("decode: error while decoding uuid: %v", err)
	}
	return b
}

func NewSequentialUUID() ID {
	uuid := make([]byte, 16)

	nano := time.Now().UnixNano()
	incr := atomic.AddUint32(&seq, 1)

	binary.BigEndian.PutUint64(uuid[0:], uint64(nano))
	binary.BigEndian.PutUint32(uuid[8:], node)
	binary.BigEndian.PutUint32(uuid[12:], incr)

	return uuid
}

func (u ID) Bytes() []byte {
	return []byte(u)
}

/*func (u seqUUID) UnixNano() int64 {
	return int64(binary.BigEndian.Uint64([]byte(u)))
}

func NewUUIDPrefix(nsec int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(nsec))
	return b
}

func (u seqUUID) Time() time.Time {
	nsec := binary.BigEndian.Uint64([]byte(u))
	return time.Unix(0, int64(nsec))
}
func (u seqUUID) Node() uint32 {
	return binary.BigEndian.Uint32([]byte(u)[12:])
}
func (u seqUUID) Sequence() uint32 {
	return binary.BigEndian.Uint32([]byte(u)[8:])
}

func (u seqUUID) After(another seqUUID) bool {
	if u.Node() != another.Node() {
		panic("Can't match UUIDs from different nodes")
	}
	t1 := u.Time()
	t2 := another.Time()
	if t1 == t2 {
		const halfway uint32 = 0xFFFFFFFF / 2
		// clocks match, let us compare sequences with wrap
		s1 := u.Sequence()
		s2 := u.Sequence()

		return s1-s2 < halfway

	} else {
		return t1.After(t2)
	}
}*/
