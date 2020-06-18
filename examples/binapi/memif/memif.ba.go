// Code generated by GoVPP's binapi-generator. DO NOT EDIT.
// versions:
//  binapi-generator: v0.4.0-alpha-1-g435c3f4-dirty
//  VPP:              20.01-45~g7a071e370~b63
// source: /usr/share/vpp/api/plugins/memif.api.json

/*
Package memif contains generated code for VPP binary API defined by memif.api (version 3.0.0).

It consists of:
	  2 aliases
	  8 enums
	 10 messages
*/
package memif

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"math"
	"strconv"

	api "git.fd.io/govpp.git/api"
	codec "git.fd.io/govpp.git/codec"
	struc "github.com/lunixbochs/struc"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the GoVPP api package it is being compiled against.
// A compilation error at this line likely means your copy of the
// GoVPP api package needs to be updated.
const _ = api.GoVppAPIPackageIsVersion2 // please upgrade the GoVPP api package

const (
	// ModuleName is the name of this module.
	ModuleName = "memif"
	// APIVersion is the API version of this module.
	APIVersion = "3.0.0"
	// VersionCrc is the CRC of this module.
	VersionCrc = 0x88dc56c9
)

// IfStatusFlags represents VPP binary API enum 'if_status_flags'.
type IfStatusFlags uint32

const (
	IF_STATUS_API_FLAG_ADMIN_UP IfStatusFlags = 1
	IF_STATUS_API_FLAG_LINK_UP  IfStatusFlags = 2
)

var (
	IfStatusFlags_name = map[uint32]string{
		1: "IF_STATUS_API_FLAG_ADMIN_UP",
		2: "IF_STATUS_API_FLAG_LINK_UP",
	}
	IfStatusFlags_value = map[string]uint32{
		"IF_STATUS_API_FLAG_ADMIN_UP": 1,
		"IF_STATUS_API_FLAG_LINK_UP":  2,
	}
)

func (x IfStatusFlags) String() string {
	s, ok := IfStatusFlags_name[uint32(x)]
	if ok {
		return s
	}
	return "IfStatusFlags(" + strconv.Itoa(int(x)) + ")"
}

// IfType represents VPP binary API enum 'if_type'.
type IfType uint32

const (
	IF_API_TYPE_HARDWARE IfType = 1
	IF_API_TYPE_SUB      IfType = 2
	IF_API_TYPE_P2P      IfType = 3
	IF_API_TYPE_PIPE     IfType = 4
)

var (
	IfType_name = map[uint32]string{
		1: "IF_API_TYPE_HARDWARE",
		2: "IF_API_TYPE_SUB",
		3: "IF_API_TYPE_P2P",
		4: "IF_API_TYPE_PIPE",
	}
	IfType_value = map[string]uint32{
		"IF_API_TYPE_HARDWARE": 1,
		"IF_API_TYPE_SUB":      2,
		"IF_API_TYPE_P2P":      3,
		"IF_API_TYPE_PIPE":     4,
	}
)

func (x IfType) String() string {
	s, ok := IfType_name[uint32(x)]
	if ok {
		return s
	}
	return "IfType(" + strconv.Itoa(int(x)) + ")"
}

// LinkDuplex represents VPP binary API enum 'link_duplex'.
type LinkDuplex uint32

const (
	LINK_DUPLEX_API_UNKNOWN LinkDuplex = 0
	LINK_DUPLEX_API_HALF    LinkDuplex = 1
	LINK_DUPLEX_API_FULL    LinkDuplex = 2
)

var (
	LinkDuplex_name = map[uint32]string{
		0: "LINK_DUPLEX_API_UNKNOWN",
		1: "LINK_DUPLEX_API_HALF",
		2: "LINK_DUPLEX_API_FULL",
	}
	LinkDuplex_value = map[string]uint32{
		"LINK_DUPLEX_API_UNKNOWN": 0,
		"LINK_DUPLEX_API_HALF":    1,
		"LINK_DUPLEX_API_FULL":    2,
	}
)

func (x LinkDuplex) String() string {
	s, ok := LinkDuplex_name[uint32(x)]
	if ok {
		return s
	}
	return "LinkDuplex(" + strconv.Itoa(int(x)) + ")"
}

// MemifMode represents VPP binary API enum 'memif_mode'.
type MemifMode uint32

const (
	MEMIF_MODE_API_ETHERNET    MemifMode = 0
	MEMIF_MODE_API_IP          MemifMode = 1
	MEMIF_MODE_API_PUNT_INJECT MemifMode = 2
)

var (
	MemifMode_name = map[uint32]string{
		0: "MEMIF_MODE_API_ETHERNET",
		1: "MEMIF_MODE_API_IP",
		2: "MEMIF_MODE_API_PUNT_INJECT",
	}
	MemifMode_value = map[string]uint32{
		"MEMIF_MODE_API_ETHERNET":    0,
		"MEMIF_MODE_API_IP":          1,
		"MEMIF_MODE_API_PUNT_INJECT": 2,
	}
)

func (x MemifMode) String() string {
	s, ok := MemifMode_name[uint32(x)]
	if ok {
		return s
	}
	return "MemifMode(" + strconv.Itoa(int(x)) + ")"
}

// MemifRole represents VPP binary API enum 'memif_role'.
type MemifRole uint32

const (
	MEMIF_ROLE_API_MASTER MemifRole = 0
	MEMIF_ROLE_API_SLAVE  MemifRole = 1
)

var (
	MemifRole_name = map[uint32]string{
		0: "MEMIF_ROLE_API_MASTER",
		1: "MEMIF_ROLE_API_SLAVE",
	}
	MemifRole_value = map[string]uint32{
		"MEMIF_ROLE_API_MASTER": 0,
		"MEMIF_ROLE_API_SLAVE":  1,
	}
)

func (x MemifRole) String() string {
	s, ok := MemifRole_name[uint32(x)]
	if ok {
		return s
	}
	return "MemifRole(" + strconv.Itoa(int(x)) + ")"
}

// MtuProto represents VPP binary API enum 'mtu_proto'.
type MtuProto uint32

const (
	MTU_PROTO_API_L3   MtuProto = 1
	MTU_PROTO_API_IP4  MtuProto = 2
	MTU_PROTO_API_IP6  MtuProto = 3
	MTU_PROTO_API_MPLS MtuProto = 4
	MTU_PROTO_API_N    MtuProto = 5
)

var (
	MtuProto_name = map[uint32]string{
		1: "MTU_PROTO_API_L3",
		2: "MTU_PROTO_API_IP4",
		3: "MTU_PROTO_API_IP6",
		4: "MTU_PROTO_API_MPLS",
		5: "MTU_PROTO_API_N",
	}
	MtuProto_value = map[string]uint32{
		"MTU_PROTO_API_L3":   1,
		"MTU_PROTO_API_IP4":  2,
		"MTU_PROTO_API_IP6":  3,
		"MTU_PROTO_API_MPLS": 4,
		"MTU_PROTO_API_N":    5,
	}
)

func (x MtuProto) String() string {
	s, ok := MtuProto_name[uint32(x)]
	if ok {
		return s
	}
	return "MtuProto(" + strconv.Itoa(int(x)) + ")"
}

// RxMode represents VPP binary API enum 'rx_mode'.
type RxMode uint32

const (
	RX_MODE_API_UNKNOWN   RxMode = 0
	RX_MODE_API_POLLING   RxMode = 1
	RX_MODE_API_INTERRUPT RxMode = 2
	RX_MODE_API_ADAPTIVE  RxMode = 3
	RX_MODE_API_DEFAULT   RxMode = 4
)

var (
	RxMode_name = map[uint32]string{
		0: "RX_MODE_API_UNKNOWN",
		1: "RX_MODE_API_POLLING",
		2: "RX_MODE_API_INTERRUPT",
		3: "RX_MODE_API_ADAPTIVE",
		4: "RX_MODE_API_DEFAULT",
	}
	RxMode_value = map[string]uint32{
		"RX_MODE_API_UNKNOWN":   0,
		"RX_MODE_API_POLLING":   1,
		"RX_MODE_API_INTERRUPT": 2,
		"RX_MODE_API_ADAPTIVE":  3,
		"RX_MODE_API_DEFAULT":   4,
	}
)

func (x RxMode) String() string {
	s, ok := RxMode_name[uint32(x)]
	if ok {
		return s
	}
	return "RxMode(" + strconv.Itoa(int(x)) + ")"
}

// SubIfFlags represents VPP binary API enum 'sub_if_flags'.
type SubIfFlags uint32

const (
	SUB_IF_API_FLAG_NO_TAGS           SubIfFlags = 1
	SUB_IF_API_FLAG_ONE_TAG           SubIfFlags = 2
	SUB_IF_API_FLAG_TWO_TAGS          SubIfFlags = 4
	SUB_IF_API_FLAG_DOT1AD            SubIfFlags = 8
	SUB_IF_API_FLAG_EXACT_MATCH       SubIfFlags = 16
	SUB_IF_API_FLAG_DEFAULT           SubIfFlags = 32
	SUB_IF_API_FLAG_OUTER_VLAN_ID_ANY SubIfFlags = 64
	SUB_IF_API_FLAG_INNER_VLAN_ID_ANY SubIfFlags = 128
	SUB_IF_API_FLAG_MASK_VNET         SubIfFlags = 254
	SUB_IF_API_FLAG_DOT1AH            SubIfFlags = 256
)

var (
	SubIfFlags_name = map[uint32]string{
		1:   "SUB_IF_API_FLAG_NO_TAGS",
		2:   "SUB_IF_API_FLAG_ONE_TAG",
		4:   "SUB_IF_API_FLAG_TWO_TAGS",
		8:   "SUB_IF_API_FLAG_DOT1AD",
		16:  "SUB_IF_API_FLAG_EXACT_MATCH",
		32:  "SUB_IF_API_FLAG_DEFAULT",
		64:  "SUB_IF_API_FLAG_OUTER_VLAN_ID_ANY",
		128: "SUB_IF_API_FLAG_INNER_VLAN_ID_ANY",
		254: "SUB_IF_API_FLAG_MASK_VNET",
		256: "SUB_IF_API_FLAG_DOT1AH",
	}
	SubIfFlags_value = map[string]uint32{
		"SUB_IF_API_FLAG_NO_TAGS":           1,
		"SUB_IF_API_FLAG_ONE_TAG":           2,
		"SUB_IF_API_FLAG_TWO_TAGS":          4,
		"SUB_IF_API_FLAG_DOT1AD":            8,
		"SUB_IF_API_FLAG_EXACT_MATCH":       16,
		"SUB_IF_API_FLAG_DEFAULT":           32,
		"SUB_IF_API_FLAG_OUTER_VLAN_ID_ANY": 64,
		"SUB_IF_API_FLAG_INNER_VLAN_ID_ANY": 128,
		"SUB_IF_API_FLAG_MASK_VNET":         254,
		"SUB_IF_API_FLAG_DOT1AH":            256,
	}
)

func (x SubIfFlags) String() string {
	s, ok := SubIfFlags_name[uint32(x)]
	if ok {
		return s
	}
	return "SubIfFlags(" + strconv.Itoa(int(x)) + ")"
}

// InterfaceIndex represents VPP binary API alias 'interface_index'.
type InterfaceIndex uint32

// MacAddress represents VPP binary API alias 'mac_address'.
type MacAddress [6]uint8

// MemifCreate represents VPP binary API message 'memif_create'.
type MemifCreate struct {
	Role       MemifRole  `binapi:"memif_role,name=role" json:"role,omitempty"`
	Mode       MemifMode  `binapi:"memif_mode,name=mode" json:"mode,omitempty"`
	RxQueues   uint8      `binapi:"u8,name=rx_queues" json:"rx_queues,omitempty"`
	TxQueues   uint8      `binapi:"u8,name=tx_queues" json:"tx_queues,omitempty"`
	ID         uint32     `binapi:"u32,name=id" json:"id,omitempty"`
	SocketID   uint32     `binapi:"u32,name=socket_id" json:"socket_id,omitempty"`
	RingSize   uint32     `binapi:"u32,name=ring_size" json:"ring_size,omitempty"`
	BufferSize uint16     `binapi:"u16,name=buffer_size" json:"buffer_size,omitempty"`
	NoZeroCopy bool       `binapi:"bool,name=no_zero_copy" json:"no_zero_copy,omitempty"`
	HwAddr     MacAddress `binapi:"mac_address,name=hw_addr" json:"hw_addr,omitempty"`
	Secret     string     `binapi:"string[24],name=secret" json:"secret,omitempty" struc:"[24]byte"`
}

func (m *MemifCreate) Reset()                        { *m = MemifCreate{} }
func (*MemifCreate) GetMessageName() string          { return "memif_create" }
func (*MemifCreate) GetCrcString() string            { return "b1b25061" }
func (*MemifCreate) GetMessageType() api.MessageType { return api.RequestMessage }

func (m *MemifCreate) Size() int {
	if m == nil {
		return 0
	}
	var size int
	// field[1] m.Role
	size += 4
	// field[1] m.Mode
	size += 4
	// field[1] m.RxQueues
	size += 1
	// field[1] m.TxQueues
	size += 1
	// field[1] m.ID
	size += 4
	// field[1] m.SocketID
	size += 4
	// field[1] m.RingSize
	size += 4
	// field[1] m.BufferSize
	size += 2
	// field[1] m.NoZeroCopy
	size += 1
	// field[1] m.HwAddr
	size += 6
	// field[1] m.Secret
	size += 24
	return size
}
func (m *MemifCreate) Marshal(b []byte) ([]byte, error) {
	o := binary.BigEndian
	_ = o
	pos := 0
	_ = pos
	var buf []byte
	if b == nil {
		buf = make([]byte, m.Size())
	} else {
		buf = b
	}
	// field[1] m.Role
	o.PutUint32(buf[pos:pos+4], uint32(m.Role))
	pos += 4
	// field[1] m.Mode
	o.PutUint32(buf[pos:pos+4], uint32(m.Mode))
	pos += 4
	// field[1] m.RxQueues
	buf[pos] = uint8(m.RxQueues)
	pos += 1
	// field[1] m.TxQueues
	buf[pos] = uint8(m.TxQueues)
	pos += 1
	// field[1] m.ID
	o.PutUint32(buf[pos:pos+4], uint32(m.ID))
	pos += 4
	// field[1] m.SocketID
	o.PutUint32(buf[pos:pos+4], uint32(m.SocketID))
	pos += 4
	// field[1] m.RingSize
	o.PutUint32(buf[pos:pos+4], uint32(m.RingSize))
	pos += 4
	// field[1] m.BufferSize
	o.PutUint16(buf[pos:pos+2], uint16(m.BufferSize))
	pos += 2
	// field[1] m.NoZeroCopy
	if m.NoZeroCopy {
		buf[pos] = 1
	}
	pos += 1
	// field[1] m.HwAddr
	for i := 0; i < 6; i++ {
		var x uint8
		if i < len(m.HwAddr) {
			x = uint8(m.HwAddr[i])
		}
		buf[pos] = uint8(x)
		pos += 1
	}
	// field[1] m.Secret
	copy(buf[pos:pos+24], m.Secret)
	pos += 24
	return buf, nil
}
func (m *MemifCreate) Unmarshal(tmp []byte) error {
	o := binary.BigEndian
	_ = o
	pos := 0
	_ = pos
	// field[1] m.Role
	m.Role = MemifRole(o.Uint32(tmp[pos : pos+4]))
	pos += 4
	// field[1] m.Mode
	m.Mode = MemifMode(o.Uint32(tmp[pos : pos+4]))
	pos += 4
	// field[1] m.RxQueues
	m.RxQueues = uint8(tmp[pos])
	pos += 1
	// field[1] m.TxQueues
	m.TxQueues = uint8(tmp[pos])
	pos += 1
	// field[1] m.ID
	m.ID = uint32(o.Uint32(tmp[pos : pos+4]))
	pos += 4
	// field[1] m.SocketID
	m.SocketID = uint32(o.Uint32(tmp[pos : pos+4]))
	pos += 4
	// field[1] m.RingSize
	m.RingSize = uint32(o.Uint32(tmp[pos : pos+4]))
	pos += 4
	// field[1] m.BufferSize
	m.BufferSize = uint16(o.Uint16(tmp[pos : pos+2]))
	pos += 2
	// field[1] m.NoZeroCopy
	m.NoZeroCopy = tmp[pos] != 0
	pos += 1
	// field[1] m.HwAddr
	for i := 0; i < len(m.HwAddr); i++ {
		m.HwAddr[i] = uint8(tmp[pos])
		pos += 1
	}
	// field[1] m.Secret
	{
		nul := bytes.Index(tmp[pos:pos+24], []byte{0x00})
		m.Secret = codec.DecodeString(tmp[pos : pos+nul])
		pos += 24
	}
	return nil
}

// MemifCreateReply represents VPP binary API message 'memif_create_reply'.
type MemifCreateReply struct {
	Retval    int32          `binapi:"i32,name=retval" json:"retval,omitempty"`
	SwIfIndex InterfaceIndex `binapi:"interface_index,name=sw_if_index" json:"sw_if_index,omitempty"`
}

func (m *MemifCreateReply) Reset()                        { *m = MemifCreateReply{} }
func (*MemifCreateReply) GetMessageName() string          { return "memif_create_reply" }
func (*MemifCreateReply) GetCrcString() string            { return "5383d31f" }
func (*MemifCreateReply) GetMessageType() api.MessageType { return api.ReplyMessage }

func (m *MemifCreateReply) Size() int {
	if m == nil {
		return 0
	}
	var size int
	// field[1] m.Retval
	size += 4
	// field[1] m.SwIfIndex
	size += 4
	return size
}
func (m *MemifCreateReply) Marshal(b []byte) ([]byte, error) {
	o := binary.BigEndian
	_ = o
	pos := 0
	_ = pos
	var buf []byte
	if b == nil {
		buf = make([]byte, m.Size())
	} else {
		buf = b
	}
	// field[1] m.Retval
	o.PutUint32(buf[pos:pos+4], uint32(m.Retval))
	pos += 4
	// field[1] m.SwIfIndex
	o.PutUint32(buf[pos:pos+4], uint32(m.SwIfIndex))
	pos += 4
	return buf, nil
}
func (m *MemifCreateReply) Unmarshal(tmp []byte) error {
	o := binary.BigEndian
	_ = o
	pos := 0
	_ = pos
	// field[1] m.Retval
	m.Retval = int32(o.Uint32(tmp[pos : pos+4]))
	pos += 4
	// field[1] m.SwIfIndex
	m.SwIfIndex = InterfaceIndex(o.Uint32(tmp[pos : pos+4]))
	pos += 4
	return nil
}

// MemifDelete represents VPP binary API message 'memif_delete'.
type MemifDelete struct {
	SwIfIndex InterfaceIndex `binapi:"interface_index,name=sw_if_index" json:"sw_if_index,omitempty"`
}

func (m *MemifDelete) Reset()                        { *m = MemifDelete{} }
func (*MemifDelete) GetMessageName() string          { return "memif_delete" }
func (*MemifDelete) GetCrcString() string            { return "f9e6675e" }
func (*MemifDelete) GetMessageType() api.MessageType { return api.RequestMessage }

func (m *MemifDelete) Size() int {
	if m == nil {
		return 0
	}
	var size int
	// field[1] m.SwIfIndex
	size += 4
	return size
}
func (m *MemifDelete) Marshal(b []byte) ([]byte, error) {
	o := binary.BigEndian
	_ = o
	pos := 0
	_ = pos
	var buf []byte
	if b == nil {
		buf = make([]byte, m.Size())
	} else {
		buf = b
	}
	// field[1] m.SwIfIndex
	o.PutUint32(buf[pos:pos+4], uint32(m.SwIfIndex))
	pos += 4
	return buf, nil
}
func (m *MemifDelete) Unmarshal(tmp []byte) error {
	o := binary.BigEndian
	_ = o
	pos := 0
	_ = pos
	// field[1] m.SwIfIndex
	m.SwIfIndex = InterfaceIndex(o.Uint32(tmp[pos : pos+4]))
	pos += 4
	return nil
}

// MemifDeleteReply represents VPP binary API message 'memif_delete_reply'.
type MemifDeleteReply struct {
	Retval int32 `binapi:"i32,name=retval" json:"retval,omitempty"`
}

func (m *MemifDeleteReply) Reset()                        { *m = MemifDeleteReply{} }
func (*MemifDeleteReply) GetMessageName() string          { return "memif_delete_reply" }
func (*MemifDeleteReply) GetCrcString() string            { return "e8d4e804" }
func (*MemifDeleteReply) GetMessageType() api.MessageType { return api.ReplyMessage }

func (m *MemifDeleteReply) Size() int {
	if m == nil {
		return 0
	}
	var size int
	// field[1] m.Retval
	size += 4
	return size
}
func (m *MemifDeleteReply) Marshal(b []byte) ([]byte, error) {
	o := binary.BigEndian
	_ = o
	pos := 0
	_ = pos
	var buf []byte
	if b == nil {
		buf = make([]byte, m.Size())
	} else {
		buf = b
	}
	// field[1] m.Retval
	o.PutUint32(buf[pos:pos+4], uint32(m.Retval))
	pos += 4
	return buf, nil
}
func (m *MemifDeleteReply) Unmarshal(tmp []byte) error {
	o := binary.BigEndian
	_ = o
	pos := 0
	_ = pos
	// field[1] m.Retval
	m.Retval = int32(o.Uint32(tmp[pos : pos+4]))
	pos += 4
	return nil
}

// MemifDetails represents VPP binary API message 'memif_details'.
type MemifDetails struct {
	SwIfIndex  InterfaceIndex `binapi:"interface_index,name=sw_if_index" json:"sw_if_index,omitempty"`
	HwAddr     MacAddress     `binapi:"mac_address,name=hw_addr" json:"hw_addr,omitempty"`
	ID         uint32         `binapi:"u32,name=id" json:"id,omitempty"`
	Role       MemifRole      `binapi:"memif_role,name=role" json:"role,omitempty"`
	Mode       MemifMode      `binapi:"memif_mode,name=mode" json:"mode,omitempty"`
	ZeroCopy   bool           `binapi:"bool,name=zero_copy" json:"zero_copy,omitempty"`
	SocketID   uint32         `binapi:"u32,name=socket_id" json:"socket_id,omitempty"`
	RingSize   uint32         `binapi:"u32,name=ring_size" json:"ring_size,omitempty"`
	BufferSize uint16         `binapi:"u16,name=buffer_size" json:"buffer_size,omitempty"`
	Flags      IfStatusFlags  `binapi:"if_status_flags,name=flags" json:"flags,omitempty"`
	IfName     string         `binapi:"string[64],name=if_name" json:"if_name,omitempty" struc:"[64]byte"`
}

func (m *MemifDetails) Reset()                        { *m = MemifDetails{} }
func (*MemifDetails) GetMessageName() string          { return "memif_details" }
func (*MemifDetails) GetCrcString() string            { return "d0382c4c" }
func (*MemifDetails) GetMessageType() api.MessageType { return api.ReplyMessage }

func (m *MemifDetails) Size() int {
	if m == nil {
		return 0
	}
	var size int
	// field[1] m.SwIfIndex
	size += 4
	// field[1] m.HwAddr
	size += 6
	// field[1] m.ID
	size += 4
	// field[1] m.Role
	size += 4
	// field[1] m.Mode
	size += 4
	// field[1] m.ZeroCopy
	size += 1
	// field[1] m.SocketID
	size += 4
	// field[1] m.RingSize
	size += 4
	// field[1] m.BufferSize
	size += 2
	// field[1] m.Flags
	size += 4
	// field[1] m.IfName
	size += 64
	return size
}
func (m *MemifDetails) Marshal(b []byte) ([]byte, error) {
	o := binary.BigEndian
	_ = o
	pos := 0
	_ = pos
	var buf []byte
	if b == nil {
		buf = make([]byte, m.Size())
	} else {
		buf = b
	}
	// field[1] m.SwIfIndex
	o.PutUint32(buf[pos:pos+4], uint32(m.SwIfIndex))
	pos += 4
	// field[1] m.HwAddr
	for i := 0; i < 6; i++ {
		var x uint8
		if i < len(m.HwAddr) {
			x = uint8(m.HwAddr[i])
		}
		buf[pos] = uint8(x)
		pos += 1
	}
	// field[1] m.ID
	o.PutUint32(buf[pos:pos+4], uint32(m.ID))
	pos += 4
	// field[1] m.Role
	o.PutUint32(buf[pos:pos+4], uint32(m.Role))
	pos += 4
	// field[1] m.Mode
	o.PutUint32(buf[pos:pos+4], uint32(m.Mode))
	pos += 4
	// field[1] m.ZeroCopy
	if m.ZeroCopy {
		buf[pos] = 1
	}
	pos += 1
	// field[1] m.SocketID
	o.PutUint32(buf[pos:pos+4], uint32(m.SocketID))
	pos += 4
	// field[1] m.RingSize
	o.PutUint32(buf[pos:pos+4], uint32(m.RingSize))
	pos += 4
	// field[1] m.BufferSize
	o.PutUint16(buf[pos:pos+2], uint16(m.BufferSize))
	pos += 2
	// field[1] m.Flags
	o.PutUint32(buf[pos:pos+4], uint32(m.Flags))
	pos += 4
	// field[1] m.IfName
	copy(buf[pos:pos+64], m.IfName)
	pos += 64
	return buf, nil
}
func (m *MemifDetails) Unmarshal(tmp []byte) error {
	o := binary.BigEndian
	_ = o
	pos := 0
	_ = pos
	// field[1] m.SwIfIndex
	m.SwIfIndex = InterfaceIndex(o.Uint32(tmp[pos : pos+4]))
	pos += 4
	// field[1] m.HwAddr
	for i := 0; i < len(m.HwAddr); i++ {
		m.HwAddr[i] = uint8(tmp[pos])
		pos += 1
	}
	// field[1] m.ID
	m.ID = uint32(o.Uint32(tmp[pos : pos+4]))
	pos += 4
	// field[1] m.Role
	m.Role = MemifRole(o.Uint32(tmp[pos : pos+4]))
	pos += 4
	// field[1] m.Mode
	m.Mode = MemifMode(o.Uint32(tmp[pos : pos+4]))
	pos += 4
	// field[1] m.ZeroCopy
	m.ZeroCopy = tmp[pos] != 0
	pos += 1
	// field[1] m.SocketID
	m.SocketID = uint32(o.Uint32(tmp[pos : pos+4]))
	pos += 4
	// field[1] m.RingSize
	m.RingSize = uint32(o.Uint32(tmp[pos : pos+4]))
	pos += 4
	// field[1] m.BufferSize
	m.BufferSize = uint16(o.Uint16(tmp[pos : pos+2]))
	pos += 2
	// field[1] m.Flags
	m.Flags = IfStatusFlags(o.Uint32(tmp[pos : pos+4]))
	pos += 4
	// field[1] m.IfName
	{
		nul := bytes.Index(tmp[pos:pos+64], []byte{0x00})
		m.IfName = codec.DecodeString(tmp[pos : pos+nul])
		pos += 64
	}
	return nil
}

// MemifDump represents VPP binary API message 'memif_dump'.
type MemifDump struct{}

func (m *MemifDump) Reset()                        { *m = MemifDump{} }
func (*MemifDump) GetMessageName() string          { return "memif_dump" }
func (*MemifDump) GetCrcString() string            { return "51077d14" }
func (*MemifDump) GetMessageType() api.MessageType { return api.RequestMessage }

func (m *MemifDump) Size() int {
	if m == nil {
		return 0
	}
	var size int
	return size
}
func (m *MemifDump) Marshal(b []byte) ([]byte, error) {
	o := binary.BigEndian
	_ = o
	pos := 0
	_ = pos
	var buf []byte
	if b == nil {
		buf = make([]byte, m.Size())
	} else {
		buf = b
	}
	return buf, nil
}
func (m *MemifDump) Unmarshal(tmp []byte) error {
	o := binary.BigEndian
	_ = o
	pos := 0
	_ = pos
	return nil
}

// MemifSocketFilenameAddDel represents VPP binary API message 'memif_socket_filename_add_del'.
type MemifSocketFilenameAddDel struct {
	IsAdd          bool   `binapi:"bool,name=is_add" json:"is_add,omitempty"`
	SocketID       uint32 `binapi:"u32,name=socket_id" json:"socket_id,omitempty"`
	SocketFilename string `binapi:"string[108],name=socket_filename" json:"socket_filename,omitempty" struc:"[108]byte"`
}

func (m *MemifSocketFilenameAddDel) Reset()                        { *m = MemifSocketFilenameAddDel{} }
func (*MemifSocketFilenameAddDel) GetMessageName() string          { return "memif_socket_filename_add_del" }
func (*MemifSocketFilenameAddDel) GetCrcString() string            { return "a2ce1a10" }
func (*MemifSocketFilenameAddDel) GetMessageType() api.MessageType { return api.RequestMessage }

func (m *MemifSocketFilenameAddDel) Size() int {
	if m == nil {
		return 0
	}
	var size int
	// field[1] m.IsAdd
	size += 1
	// field[1] m.SocketID
	size += 4
	// field[1] m.SocketFilename
	size += 108
	return size
}
func (m *MemifSocketFilenameAddDel) Marshal(b []byte) ([]byte, error) {
	o := binary.BigEndian
	_ = o
	pos := 0
	_ = pos
	var buf []byte
	if b == nil {
		buf = make([]byte, m.Size())
	} else {
		buf = b
	}
	// field[1] m.IsAdd
	if m.IsAdd {
		buf[pos] = 1
	}
	pos += 1
	// field[1] m.SocketID
	o.PutUint32(buf[pos:pos+4], uint32(m.SocketID))
	pos += 4
	// field[1] m.SocketFilename
	copy(buf[pos:pos+108], m.SocketFilename)
	pos += 108
	return buf, nil
}
func (m *MemifSocketFilenameAddDel) Unmarshal(tmp []byte) error {
	o := binary.BigEndian
	_ = o
	pos := 0
	_ = pos
	// field[1] m.IsAdd
	m.IsAdd = tmp[pos] != 0
	pos += 1
	// field[1] m.SocketID
	m.SocketID = uint32(o.Uint32(tmp[pos : pos+4]))
	pos += 4
	// field[1] m.SocketFilename
	{
		nul := bytes.Index(tmp[pos:pos+108], []byte{0x00})
		m.SocketFilename = codec.DecodeString(tmp[pos : pos+nul])
		pos += 108
	}
	return nil
}

// MemifSocketFilenameAddDelReply represents VPP binary API message 'memif_socket_filename_add_del_reply'.
type MemifSocketFilenameAddDelReply struct {
	Retval int32 `binapi:"i32,name=retval" json:"retval,omitempty"`
}

func (m *MemifSocketFilenameAddDelReply) Reset() { *m = MemifSocketFilenameAddDelReply{} }
func (*MemifSocketFilenameAddDelReply) GetMessageName() string {
	return "memif_socket_filename_add_del_reply"
}
func (*MemifSocketFilenameAddDelReply) GetCrcString() string            { return "e8d4e804" }
func (*MemifSocketFilenameAddDelReply) GetMessageType() api.MessageType { return api.ReplyMessage }

func (m *MemifSocketFilenameAddDelReply) Size() int {
	if m == nil {
		return 0
	}
	var size int
	// field[1] m.Retval
	size += 4
	return size
}
func (m *MemifSocketFilenameAddDelReply) Marshal(b []byte) ([]byte, error) {
	o := binary.BigEndian
	_ = o
	pos := 0
	_ = pos
	var buf []byte
	if b == nil {
		buf = make([]byte, m.Size())
	} else {
		buf = b
	}
	// field[1] m.Retval
	o.PutUint32(buf[pos:pos+4], uint32(m.Retval))
	pos += 4
	return buf, nil
}
func (m *MemifSocketFilenameAddDelReply) Unmarshal(tmp []byte) error {
	o := binary.BigEndian
	_ = o
	pos := 0
	_ = pos
	// field[1] m.Retval
	m.Retval = int32(o.Uint32(tmp[pos : pos+4]))
	pos += 4
	return nil
}

// MemifSocketFilenameDetails represents VPP binary API message 'memif_socket_filename_details'.
type MemifSocketFilenameDetails struct {
	SocketID       uint32 `binapi:"u32,name=socket_id" json:"socket_id,omitempty"`
	SocketFilename string `binapi:"string[108],name=socket_filename" json:"socket_filename,omitempty" struc:"[108]byte"`
}

func (m *MemifSocketFilenameDetails) Reset()                        { *m = MemifSocketFilenameDetails{} }
func (*MemifSocketFilenameDetails) GetMessageName() string          { return "memif_socket_filename_details" }
func (*MemifSocketFilenameDetails) GetCrcString() string            { return "7ff326f7" }
func (*MemifSocketFilenameDetails) GetMessageType() api.MessageType { return api.ReplyMessage }

func (m *MemifSocketFilenameDetails) Size() int {
	if m == nil {
		return 0
	}
	var size int
	// field[1] m.SocketID
	size += 4
	// field[1] m.SocketFilename
	size += 108
	return size
}
func (m *MemifSocketFilenameDetails) Marshal(b []byte) ([]byte, error) {
	o := binary.BigEndian
	_ = o
	pos := 0
	_ = pos
	var buf []byte
	if b == nil {
		buf = make([]byte, m.Size())
	} else {
		buf = b
	}
	// field[1] m.SocketID
	o.PutUint32(buf[pos:pos+4], uint32(m.SocketID))
	pos += 4
	// field[1] m.SocketFilename
	copy(buf[pos:pos+108], m.SocketFilename)
	pos += 108
	return buf, nil
}
func (m *MemifSocketFilenameDetails) Unmarshal(tmp []byte) error {
	o := binary.BigEndian
	_ = o
	pos := 0
	_ = pos
	// field[1] m.SocketID
	m.SocketID = uint32(o.Uint32(tmp[pos : pos+4]))
	pos += 4
	// field[1] m.SocketFilename
	{
		nul := bytes.Index(tmp[pos:pos+108], []byte{0x00})
		m.SocketFilename = codec.DecodeString(tmp[pos : pos+nul])
		pos += 108
	}
	return nil
}

// MemifSocketFilenameDump represents VPP binary API message 'memif_socket_filename_dump'.
type MemifSocketFilenameDump struct{}

func (m *MemifSocketFilenameDump) Reset()                        { *m = MemifSocketFilenameDump{} }
func (*MemifSocketFilenameDump) GetMessageName() string          { return "memif_socket_filename_dump" }
func (*MemifSocketFilenameDump) GetCrcString() string            { return "51077d14" }
func (*MemifSocketFilenameDump) GetMessageType() api.MessageType { return api.RequestMessage }

func (m *MemifSocketFilenameDump) Size() int {
	if m == nil {
		return 0
	}
	var size int
	return size
}
func (m *MemifSocketFilenameDump) Marshal(b []byte) ([]byte, error) {
	o := binary.BigEndian
	_ = o
	pos := 0
	_ = pos
	var buf []byte
	if b == nil {
		buf = make([]byte, m.Size())
	} else {
		buf = b
	}
	return buf, nil
}
func (m *MemifSocketFilenameDump) Unmarshal(tmp []byte) error {
	o := binary.BigEndian
	_ = o
	pos := 0
	_ = pos
	return nil
}

func init() { file_memif_binapi_init() }
func file_memif_binapi_init() {
	api.RegisterMessage((*MemifCreate)(nil), "memif.MemifCreate")
	api.RegisterMessage((*MemifCreateReply)(nil), "memif.MemifCreateReply")
	api.RegisterMessage((*MemifDelete)(nil), "memif.MemifDelete")
	api.RegisterMessage((*MemifDeleteReply)(nil), "memif.MemifDeleteReply")
	api.RegisterMessage((*MemifDetails)(nil), "memif.MemifDetails")
	api.RegisterMessage((*MemifDump)(nil), "memif.MemifDump")
	api.RegisterMessage((*MemifSocketFilenameAddDel)(nil), "memif.MemifSocketFilenameAddDel")
	api.RegisterMessage((*MemifSocketFilenameAddDelReply)(nil), "memif.MemifSocketFilenameAddDelReply")
	api.RegisterMessage((*MemifSocketFilenameDetails)(nil), "memif.MemifSocketFilenameDetails")
	api.RegisterMessage((*MemifSocketFilenameDump)(nil), "memif.MemifSocketFilenameDump")
}

// Messages returns list of all messages in this module.
func AllMessages() []api.Message {
	return []api.Message{
		(*MemifCreate)(nil),
		(*MemifCreateReply)(nil),
		(*MemifDelete)(nil),
		(*MemifDeleteReply)(nil),
		(*MemifDetails)(nil),
		(*MemifDump)(nil),
		(*MemifSocketFilenameAddDel)(nil),
		(*MemifSocketFilenameAddDelReply)(nil),
		(*MemifSocketFilenameDetails)(nil),
		(*MemifSocketFilenameDump)(nil),
	}
}

// Reference imports to suppress errors if they are not otherwise used.
var _ = api.RegisterMessage
var _ = codec.DecodeString
var _ = bytes.NewBuffer
var _ = context.Background
var _ = io.Copy
var _ = strconv.Itoa
var _ = struc.Pack
var _ = binary.BigEndian
var _ = math.Float32bits
