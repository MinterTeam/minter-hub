// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: peggy/v1/types.proto

package types

import (
	fmt "fmt"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// BridgeValidator represents a validator's ETH address and its power
type BridgeValidator struct {
	Power           uint64 `protobuf:"varint,1,opt,name=power,proto3" json:"power,omitempty"`
	EthereumAddress string `protobuf:"bytes,2,opt,name=ethereum_address,json=ethereumAddress,proto3" json:"ethereum_address,omitempty"`
}

func (m *BridgeValidator) Reset()         { *m = BridgeValidator{} }
func (m *BridgeValidator) String() string { return proto.CompactTextString(m) }
func (*BridgeValidator) ProtoMessage()    {}
func (*BridgeValidator) Descriptor() ([]byte, []int) {
	return fileDescriptor_1488ca6080c6185d, []int{0}
}
func (m *BridgeValidator) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BridgeValidator) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BridgeValidator.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BridgeValidator) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BridgeValidator.Merge(m, src)
}
func (m *BridgeValidator) XXX_Size() int {
	return m.Size()
}
func (m *BridgeValidator) XXX_DiscardUnknown() {
	xxx_messageInfo_BridgeValidator.DiscardUnknown(m)
}

var xxx_messageInfo_BridgeValidator proto.InternalMessageInfo

func (m *BridgeValidator) GetPower() uint64 {
	if m != nil {
		return m.Power
	}
	return 0
}

func (m *BridgeValidator) GetEthereumAddress() string {
	if m != nil {
		return m.EthereumAddress
	}
	return ""
}

// Valset is the Ethereum Bridge Multsig Set, each peggy validator also
// maintains an ETH key to sign messages, these are used to check signatures on
// ETH because of the significant gas savings
type Valset struct {
	Nonce   uint64             `protobuf:"varint,1,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Members []*BridgeValidator `protobuf:"bytes,2,rep,name=members,proto3" json:"members,omitempty"`
	Height  uint64             `protobuf:"varint,3,opt,name=height,proto3" json:"height,omitempty"`
}

func (m *Valset) Reset()         { *m = Valset{} }
func (m *Valset) String() string { return proto.CompactTextString(m) }
func (*Valset) ProtoMessage()    {}
func (*Valset) Descriptor() ([]byte, []int) {
	return fileDescriptor_1488ca6080c6185d, []int{1}
}
func (m *Valset) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Valset) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Valset.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Valset) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Valset.Merge(m, src)
}
func (m *Valset) XXX_Size() int {
	return m.Size()
}
func (m *Valset) XXX_DiscardUnknown() {
	xxx_messageInfo_Valset.DiscardUnknown(m)
}

var xxx_messageInfo_Valset proto.InternalMessageInfo

func (m *Valset) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *Valset) GetMembers() []*BridgeValidator {
	if m != nil {
		return m.Members
	}
	return nil
}

func (m *Valset) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

// It's difficult to serialize and deserialize
// interfaces, instead we can make this struct
// that stores all the data the interface requires
// and use it to store and then re-create a interface
// object with all the same properties.
type GenericClaim struct {
	EventNonce   uint64 `protobuf:"varint,1,opt,name=event_nonce,json=eventNonce,proto3" json:"event_nonce,omitempty"`
	ClaimType    int32  `protobuf:"varint,2,opt,name=claim_type,json=claimType,proto3" json:"claim_type,omitempty"`
	Hash         []byte `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`
	EventClaimer string `protobuf:"bytes,4,opt,name=event_claimer,json=eventClaimer,proto3" json:"event_claimer,omitempty"`
}

func (m *GenericClaim) Reset()         { *m = GenericClaim{} }
func (m *GenericClaim) String() string { return proto.CompactTextString(m) }
func (*GenericClaim) ProtoMessage()    {}
func (*GenericClaim) Descriptor() ([]byte, []int) {
	return fileDescriptor_1488ca6080c6185d, []int{2}
}
func (m *GenericClaim) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenericClaim) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenericClaim.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenericClaim) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenericClaim.Merge(m, src)
}
func (m *GenericClaim) XXX_Size() int {
	return m.Size()
}
func (m *GenericClaim) XXX_DiscardUnknown() {
	xxx_messageInfo_GenericClaim.DiscardUnknown(m)
}

var xxx_messageInfo_GenericClaim proto.InternalMessageInfo

func (m *GenericClaim) GetEventNonce() uint64 {
	if m != nil {
		return m.EventNonce
	}
	return 0
}

func (m *GenericClaim) GetClaimType() int32 {
	if m != nil {
		return m.ClaimType
	}
	return 0
}

func (m *GenericClaim) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *GenericClaim) GetEventClaimer() string {
	if m != nil {
		return m.EventClaimer
	}
	return ""
}

type LockedCoins struct {
	List []types.Coin `protobuf:"bytes,1,rep,name=list,proto3" json:"list"`
}

func (m *LockedCoins) Reset()         { *m = LockedCoins{} }
func (m *LockedCoins) String() string { return proto.CompactTextString(m) }
func (*LockedCoins) ProtoMessage()    {}
func (*LockedCoins) Descriptor() ([]byte, []int) {
	return fileDescriptor_1488ca6080c6185d, []int{3}
}
func (m *LockedCoins) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LockedCoins) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LockedCoins.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LockedCoins) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LockedCoins.Merge(m, src)
}
func (m *LockedCoins) XXX_Size() int {
	return m.Size()
}
func (m *LockedCoins) XXX_DiscardUnknown() {
	xxx_messageInfo_LockedCoins.DiscardUnknown(m)
}

var xxx_messageInfo_LockedCoins proto.InternalMessageInfo

func (m *LockedCoins) GetList() []types.Coin {
	if m != nil {
		return m.List
	}
	return nil
}

func init() {
	proto.RegisterType((*BridgeValidator)(nil), "peggy.v1.BridgeValidator")
	proto.RegisterType((*Valset)(nil), "peggy.v1.Valset")
	proto.RegisterType((*GenericClaim)(nil), "peggy.v1.GenericClaim")
	proto.RegisterType((*LockedCoins)(nil), "peggy.v1.LockedCoins")
}

func init() { proto.RegisterFile("peggy/v1/types.proto", fileDescriptor_1488ca6080c6185d) }

var fileDescriptor_1488ca6080c6185d = []byte{
	// 408 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x92, 0xb1, 0x8e, 0xd3, 0x40,
	0x10, 0x86, 0x63, 0xce, 0x17, 0xb8, 0x49, 0xd0, 0xa1, 0x55, 0x84, 0xcc, 0x49, 0xf8, 0x4e, 0xa1,
	0x09, 0xcd, 0xae, 0x42, 0x9e, 0x80, 0xa4, 0x80, 0x02, 0x28, 0xac, 0xd3, 0x15, 0x34, 0xd1, 0xda,
	0x1e, 0xd9, 0xab, 0x8b, 0x77, 0xad, 0xdd, 0x4d, 0x20, 0x4f, 0x40, 0xcb, 0x63, 0x5d, 0x79, 0x25,
	0x15, 0x42, 0xc9, 0x8b, 0x20, 0xcf, 0xc6, 0x12, 0xa2, 0x9b, 0xf9, 0xf4, 0xeb, 0x1f, 0xff, 0xfe,
	0x17, 0x26, 0x2d, 0x56, 0xd5, 0x5e, 0xec, 0xe6, 0xc2, 0xef, 0x5b, 0x74, 0xbc, 0xb5, 0xc6, 0x1b,
	0xf6, 0x8c, 0x28, 0xdf, 0xcd, 0xaf, 0xd2, 0xc2, 0xb8, 0xc6, 0x38, 0x91, 0x4b, 0x87, 0x62, 0x37,
	0xcf, 0xd1, 0xcb, 0xb9, 0x28, 0x8c, 0xd2, 0x41, 0x79, 0x35, 0xa9, 0x4c, 0x65, 0x68, 0x14, 0xdd,
	0x14, 0xe8, 0x34, 0x83, 0xcb, 0xa5, 0x55, 0x65, 0x85, 0x77, 0x72, 0xa3, 0x4a, 0xe9, 0x8d, 0x65,
	0x13, 0x38, 0x6f, 0xcd, 0x37, 0xb4, 0x49, 0x74, 0x13, 0xcd, 0xe2, 0x2c, 0x2c, 0xec, 0x2d, 0xbc,
	0x40, 0x5f, 0xa3, 0xc5, 0x6d, 0xb3, 0x96, 0x65, 0x69, 0xd1, 0xb9, 0xe4, 0xc9, 0x4d, 0x34, 0xbb,
	0xc8, 0x2e, 0x7b, 0xfe, 0x3e, 0xe0, 0xe9, 0x3d, 0x0c, 0xef, 0xe4, 0xc6, 0xa1, 0xef, 0xac, 0xb4,
	0xd1, 0x05, 0xf6, 0x56, 0xb4, 0xb0, 0x05, 0x3c, 0x6d, 0xb0, 0xc9, 0xd1, 0x76, 0x0e, 0x67, 0xb3,
	0xd1, 0xbb, 0x57, 0xbc, 0x4f, 0xc1, 0xff, 0xfb, 0x98, 0xac, 0x57, 0xb2, 0x97, 0x30, 0xac, 0x51,
	0x55, 0xb5, 0x4f, 0xce, 0xc8, 0xeb, 0xb4, 0x4d, 0x7f, 0x44, 0x30, 0xfe, 0x80, 0x1a, 0xad, 0x2a,
	0x56, 0x1b, 0xa9, 0x1a, 0x76, 0x0d, 0x23, 0xdc, 0xa1, 0xf6, 0xeb, 0x7f, 0x2f, 0x03, 0xa1, 0x2f,
	0x74, 0xfe, 0x35, 0x40, 0xd1, 0x29, 0xd7, 0xdd, 0x7f, 0xa4, 0x0c, 0xe7, 0xd9, 0x05, 0x91, 0xdb,
	0x7d, 0x8b, 0x8c, 0x41, 0x5c, 0x4b, 0x57, 0xd3, 0x99, 0x71, 0x46, 0x33, 0x7b, 0x03, 0xcf, 0x83,
	0x27, 0xc9, 0xd0, 0x26, 0x31, 0x25, 0x1f, 0x13, 0x5c, 0x05, 0x36, 0x5d, 0xc2, 0xe8, 0x93, 0x29,
	0xee, 0xb1, 0x5c, 0x19, 0xa5, 0x1d, 0x5b, 0x40, 0xbc, 0x51, 0xce, 0x27, 0xd1, 0x29, 0x62, 0xa8,
	0x87, 0x77, 0xf5, 0xf0, 0x53, 0x3d, 0xbc, 0x53, 0x2e, 0xe3, 0x87, 0xdf, 0xd7, 0x83, 0x8c, 0xc4,
	0xcb, 0x8f, 0x0f, 0x87, 0x34, 0x7a, 0x3c, 0xa4, 0xd1, 0x9f, 0x43, 0x1a, 0xfd, 0x3c, 0xa6, 0x83,
	0xc7, 0x63, 0x3a, 0xf8, 0x75, 0x4c, 0x07, 0x5f, 0x79, 0xa5, 0x7c, 0xbd, 0xcd, 0x79, 0x61, 0x1a,
	0xf1, 0x59, 0x69, 0x8f, 0xf6, 0x16, 0x65, 0x23, 0x9a, 0x7a, 0x9b, 0x8b, 0xa2, 0x96, 0x4a, 0x8b,
	0xef, 0x22, 0xbc, 0x10, 0x7a, 0x1e, 0xf9, 0x90, 0xfa, 0x5d, 0xfc, 0x0d, 0x00, 0x00, 0xff, 0xff,
	0x2e, 0xd2, 0x67, 0xb0, 0x37, 0x02, 0x00, 0x00,
}

func (m *BridgeValidator) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BridgeValidator) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BridgeValidator) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.EthereumAddress) > 0 {
		i -= len(m.EthereumAddress)
		copy(dAtA[i:], m.EthereumAddress)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.EthereumAddress)))
		i--
		dAtA[i] = 0x12
	}
	if m.Power != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Power))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Valset) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Valset) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Valset) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Height != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Members) > 0 {
		for iNdEx := len(m.Members) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Members[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTypes(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Nonce != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Nonce))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *GenericClaim) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenericClaim) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenericClaim) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.EventClaimer) > 0 {
		i -= len(m.EventClaimer)
		copy(dAtA[i:], m.EventClaimer)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.EventClaimer)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Hash) > 0 {
		i -= len(m.Hash)
		copy(dAtA[i:], m.Hash)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Hash)))
		i--
		dAtA[i] = 0x1a
	}
	if m.ClaimType != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.ClaimType))
		i--
		dAtA[i] = 0x10
	}
	if m.EventNonce != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.EventNonce))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *LockedCoins) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LockedCoins) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LockedCoins) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.List) > 0 {
		for iNdEx := len(m.List) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.List[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTypes(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintTypes(dAtA []byte, offset int, v uint64) int {
	offset -= sovTypes(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *BridgeValidator) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Power != 0 {
		n += 1 + sovTypes(uint64(m.Power))
	}
	l = len(m.EthereumAddress)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}

func (m *Valset) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Nonce != 0 {
		n += 1 + sovTypes(uint64(m.Nonce))
	}
	if len(m.Members) > 0 {
		for _, e := range m.Members {
			l = e.Size()
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	if m.Height != 0 {
		n += 1 + sovTypes(uint64(m.Height))
	}
	return n
}

func (m *GenericClaim) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.EventNonce != 0 {
		n += 1 + sovTypes(uint64(m.EventNonce))
	}
	if m.ClaimType != 0 {
		n += 1 + sovTypes(uint64(m.ClaimType))
	}
	l = len(m.Hash)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.EventClaimer)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}

func (m *LockedCoins) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.List) > 0 {
		for _, e := range m.List {
			l = e.Size()
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *BridgeValidator) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: BridgeValidator: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BridgeValidator: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Power", wireType)
			}
			m.Power = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Power |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EthereumAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EthereumAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Valset) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Valset: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Valset: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nonce", wireType)
			}
			m.Nonce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Nonce |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Members", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Members = append(m.Members, &BridgeValidator{})
			if err := m.Members[len(m.Members)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *GenericClaim) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GenericClaim: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenericClaim: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EventNonce", wireType)
			}
			m.EventNonce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EventNonce |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClaimType", wireType)
			}
			m.ClaimType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ClaimType |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash = append(m.Hash[:0], dAtA[iNdEx:postIndex]...)
			if m.Hash == nil {
				m.Hash = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EventClaimer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EventClaimer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *LockedCoins) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: LockedCoins: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LockedCoins: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field List", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.List = append(m.List, types.Coin{})
			if err := m.List[len(m.List)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTypes(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTypes
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTypes
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTypes
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTypes        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTypes          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTypes = fmt.Errorf("proto: unexpected end of group")
)
