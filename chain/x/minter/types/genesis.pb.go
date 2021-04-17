// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: minter/v1/genesis.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
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

type Params struct {
	StartThreshold                uint64                                 `protobuf:"varint,1,opt,name=start_threshold,json=startThreshold,proto3" json:"start_threshold,omitempty"`
	MinterAddress                 string                                 `protobuf:"bytes,2,opt,name=minter_address,json=minterAddress,proto3" json:"minter_address,omitempty"`
	BridgeChainId                 uint64                                 `protobuf:"varint,3,opt,name=bridge_chain_id,json=bridgeChainId,proto3" json:"bridge_chain_id,omitempty"`
	SignedValsetsWindow           uint64                                 `protobuf:"varint,4,opt,name=signed_valsets_window,json=signedValsetsWindow,proto3" json:"signed_valsets_window,omitempty"`
	SignedBatchesWindow           uint64                                 `protobuf:"varint,5,opt,name=signed_batches_window,json=signedBatchesWindow,proto3" json:"signed_batches_window,omitempty"`
	SignedClaimsWindow            uint64                                 `protobuf:"varint,6,opt,name=signed_claims_window,json=signedClaimsWindow,proto3" json:"signed_claims_window,omitempty"`
	SlashFractionValset           github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,7,opt,name=slash_fraction_valset,json=slashFractionValset,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"slash_fraction_valset"`
	SlashFractionBatch            github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,8,opt,name=slash_fraction_batch,json=slashFractionBatch,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"slash_fraction_batch"`
	SlashFractionClaim            github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,9,opt,name=slash_fraction_claim,json=slashFractionClaim,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"slash_fraction_claim"`
	SlashFractionConflictingClaim github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,10,opt,name=slash_fraction_conflicting_claim,json=slashFractionConflictingClaim,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"slash_fraction_conflicting_claim"`
	Stopped                       bool                                   `protobuf:"varint,11,opt,name=stopped,proto3" json:"stopped,omitempty"`
	Commission                    github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,12,opt,name=commission,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"commission"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_43fc00fc33749c12, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetStartThreshold() uint64 {
	if m != nil {
		return m.StartThreshold
	}
	return 0
}

func (m *Params) GetMinterAddress() string {
	if m != nil {
		return m.MinterAddress
	}
	return ""
}

func (m *Params) GetBridgeChainId() uint64 {
	if m != nil {
		return m.BridgeChainId
	}
	return 0
}

func (m *Params) GetSignedValsetsWindow() uint64 {
	if m != nil {
		return m.SignedValsetsWindow
	}
	return 0
}

func (m *Params) GetSignedBatchesWindow() uint64 {
	if m != nil {
		return m.SignedBatchesWindow
	}
	return 0
}

func (m *Params) GetSignedClaimsWindow() uint64 {
	if m != nil {
		return m.SignedClaimsWindow
	}
	return 0
}

func (m *Params) GetStopped() bool {
	if m != nil {
		return m.Stopped
	}
	return false
}

// GenesisState struct
type GenesisState struct {
	Params           *Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params,omitempty"`
	StartMinterNonce uint64  `protobuf:"varint,2,opt,name=start_minter_nonce,json=startMinterNonce,proto3" json:"start_minter_nonce,omitempty"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_43fc00fc33749c12, []int{1}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetParams() *Params {
	if m != nil {
		return m.Params
	}
	return nil
}

func (m *GenesisState) GetStartMinterNonce() uint64 {
	if m != nil {
		return m.StartMinterNonce
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "minter.v1.Params")
	proto.RegisterType((*GenesisState)(nil), "minter.v1.GenesisState")
}

func init() { proto.RegisterFile("minter/v1/genesis.proto", fileDescriptor_43fc00fc33749c12) }

var fileDescriptor_43fc00fc33749c12 = []byte{
	// 518 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0x41, 0x6f, 0xd3, 0x30,
	0x18, 0x6d, 0xa0, 0x74, 0xab, 0xd7, 0x6d, 0xe0, 0x75, 0x22, 0x42, 0x22, 0x8b, 0x26, 0x31, 0x8a,
	0x04, 0x09, 0x1b, 0x37, 0x6e, 0x74, 0x08, 0xb4, 0x03, 0x13, 0x0a, 0x13, 0x48, 0x5c, 0x82, 0x63,
	0x7b, 0x89, 0x45, 0x63, 0x47, 0xb1, 0xd7, 0xc2, 0x8d, 0x9f, 0xc0, 0xff, 0xe0, 0x8f, 0xec, 0xb8,
	0x23, 0x42, 0x68, 0x42, 0xed, 0x1f, 0x41, 0xfd, 0xec, 0x95, 0x6e, 0xe2, 0x54, 0xed, 0x94, 0xe4,
	0xbd, 0xef, 0xbd, 0xe7, 0x2f, 0xd2, 0x33, 0xba, 0x5b, 0x0a, 0x69, 0x78, 0x1d, 0x0f, 0x77, 0xe3,
	0x9c, 0x4b, 0xae, 0x85, 0x8e, 0xaa, 0x5a, 0x19, 0x85, 0xdb, 0x96, 0x88, 0x86, 0xbb, 0xf7, 0xba,
	0xb9, 0xca, 0x15, 0xa0, 0xf1, 0xf4, 0xcd, 0x0e, 0x6c, 0xff, 0x68, 0xa1, 0xd6, 0x5b, 0x52, 0x93,
	0x52, 0xe3, 0x87, 0x68, 0x5d, 0x1b, 0x52, 0x9b, 0xd4, 0x14, 0x35, 0xd7, 0x85, 0x1a, 0x30, 0xdf,
	0x0b, 0xbd, 0x5e, 0x33, 0x59, 0x03, 0xf8, 0xe8, 0x02, 0xc5, 0x0f, 0xd0, 0x9a, 0xb5, 0x4d, 0x09,
	0x63, 0x35, 0xd7, 0xda, 0xbf, 0x11, 0x7a, 0xbd, 0x76, 0xb2, 0x6a, 0xd1, 0x17, 0x16, 0xc4, 0x3b,
	0x68, 0x3d, 0xab, 0x05, 0xcb, 0x79, 0x4a, 0x0b, 0x22, 0x64, 0x2a, 0x98, 0x7f, 0x13, 0xfc, 0x56,
	0x2d, 0xbc, 0x3f, 0x45, 0x0f, 0x18, 0xde, 0x43, 0x9b, 0x5a, 0xe4, 0x92, 0xb3, 0x74, 0x48, 0x06,
	0x9a, 0x1b, 0x9d, 0x8e, 0x84, 0x64, 0x6a, 0xe4, 0x37, 0x61, 0x7a, 0xc3, 0x92, 0xef, 0x2d, 0xf7,
	0x01, 0xa8, 0x39, 0x4d, 0x46, 0x0c, 0x2d, 0xf8, 0x4c, 0x73, 0x6b, 0x5e, 0xd3, 0xb7, 0x9c, 0xd3,
	0x3c, 0x45, 0x5d, 0xa7, 0xa1, 0x03, 0x22, 0xca, 0x99, 0xa4, 0x05, 0x12, 0x6c, 0xb9, 0x7d, 0xa0,
	0x9c, 0x22, 0x43, 0x9b, 0x7a, 0x40, 0x74, 0x91, 0x1e, 0xd7, 0x84, 0x1a, 0xa1, 0xa4, 0x3b, 0xa1,
	0xbf, 0x14, 0x7a, 0xbd, 0x4e, 0x3f, 0x3a, 0x3d, 0xdf, 0x6a, 0xfc, 0x3a, 0xdf, 0xda, 0xc9, 0x85,
	0x29, 0x4e, 0xb2, 0x88, 0xaa, 0x32, 0xa6, 0x4a, 0x97, 0x4a, 0xbb, 0xc7, 0x13, 0xcd, 0x3e, 0xc7,
	0xe6, 0x6b, 0xc5, 0x75, 0xf4, 0x92, 0xd3, 0x64, 0x03, 0xcc, 0x5e, 0x39, 0x2f, 0xbb, 0x10, 0xfe,
	0x84, 0xba, 0x57, 0x32, 0x60, 0x23, 0x7f, 0x79, 0xa1, 0x08, 0x7c, 0x29, 0x02, 0xf6, 0xff, 0x4f,
	0x02, 0xec, 0xef, 0xb7, 0xaf, 0x21, 0x01, 0x7e, 0x17, 0x1e, 0xa1, 0xf0, 0x6a, 0x82, 0x92, 0xc7,
	0x03, 0x41, 0x8d, 0x90, 0xb9, 0x4b, 0x43, 0x0b, 0xa5, 0xdd, 0xbf, 0x9c, 0xf6, 0xcf, 0xd5, 0x06,
	0xfb, 0x68, 0x49, 0x1b, 0x55, 0x55, 0x9c, 0xf9, 0x2b, 0xa1, 0xd7, 0x5b, 0x4e, 0x2e, 0x3e, 0xf1,
	0x21, 0x42, 0x54, 0x95, 0xa5, 0xd0, 0x5a, 0x28, 0xe9, 0x77, 0x16, 0x0a, 0x9f, 0x73, 0x78, 0xde,
	0xfc, 0xf6, 0x3b, 0x6c, 0x6c, 0xe7, 0xa8, 0xf3, 0xda, 0xf6, 0xeb, 0x9d, 0x21, 0x86, 0xe3, 0x47,
	0xa8, 0x55, 0x41, 0x79, 0xa0, 0x29, 0x2b, 0x7b, 0x77, 0xa2, 0x59, 0xdf, 0x22, 0xdb, 0xaa, 0xc4,
	0x0d, 0xe0, 0xc7, 0x08, 0xdb, 0x76, 0xb9, 0xea, 0x48, 0x25, 0x29, 0x87, 0xe2, 0x34, 0x93, 0xdb,
	0xc0, 0xbc, 0x01, 0xe2, 0x70, 0x8a, 0xf7, 0x0f, 0x4e, 0xc7, 0x81, 0x77, 0x36, 0x0e, 0xbc, 0x3f,
	0xe3, 0xc0, 0xfb, 0x3e, 0x09, 0x1a, 0x67, 0x93, 0xa0, 0xf1, 0x73, 0x12, 0x34, 0x3e, 0xc6, 0x73,
	0x87, 0xb7, 0x8a, 0x23, 0x4e, 0xca, 0xb8, 0x2c, 0x4e, 0xb2, 0x18, 0x7a, 0x16, 0x7f, 0x89, 0xdd,
	0x6d, 0x00, 0x9b, 0x64, 0x2d, 0x28, 0xfa, 0xb3, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xd8, 0x98,
	0x06, 0x6f, 0x24, 0x04, 0x00, 0x00,
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Commission.Size()
		i -= size
		if _, err := m.Commission.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x62
	if m.Stopped {
		i--
		if m.Stopped {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x58
	}
	{
		size := m.SlashFractionConflictingClaim.Size()
		i -= size
		if _, err := m.SlashFractionConflictingClaim.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x52
	{
		size := m.SlashFractionClaim.Size()
		i -= size
		if _, err := m.SlashFractionClaim.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x4a
	{
		size := m.SlashFractionBatch.Size()
		i -= size
		if _, err := m.SlashFractionBatch.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x42
	{
		size := m.SlashFractionValset.Size()
		i -= size
		if _, err := m.SlashFractionValset.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	if m.SignedClaimsWindow != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.SignedClaimsWindow))
		i--
		dAtA[i] = 0x30
	}
	if m.SignedBatchesWindow != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.SignedBatchesWindow))
		i--
		dAtA[i] = 0x28
	}
	if m.SignedValsetsWindow != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.SignedValsetsWindow))
		i--
		dAtA[i] = 0x20
	}
	if m.BridgeChainId != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.BridgeChainId))
		i--
		dAtA[i] = 0x18
	}
	if len(m.MinterAddress) > 0 {
		i -= len(m.MinterAddress)
		copy(dAtA[i:], m.MinterAddress)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.MinterAddress)))
		i--
		dAtA[i] = 0x12
	}
	if m.StartThreshold != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.StartThreshold))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.StartMinterNonce != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.StartMinterNonce))
		i--
		dAtA[i] = 0x10
	}
	if m.Params != nil {
		{
			size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.StartThreshold != 0 {
		n += 1 + sovGenesis(uint64(m.StartThreshold))
	}
	l = len(m.MinterAddress)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if m.BridgeChainId != 0 {
		n += 1 + sovGenesis(uint64(m.BridgeChainId))
	}
	if m.SignedValsetsWindow != 0 {
		n += 1 + sovGenesis(uint64(m.SignedValsetsWindow))
	}
	if m.SignedBatchesWindow != 0 {
		n += 1 + sovGenesis(uint64(m.SignedBatchesWindow))
	}
	if m.SignedClaimsWindow != 0 {
		n += 1 + sovGenesis(uint64(m.SignedClaimsWindow))
	}
	l = m.SlashFractionValset.Size()
	n += 1 + l + sovGenesis(uint64(l))
	l = m.SlashFractionBatch.Size()
	n += 1 + l + sovGenesis(uint64(l))
	l = m.SlashFractionClaim.Size()
	n += 1 + l + sovGenesis(uint64(l))
	l = m.SlashFractionConflictingClaim.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if m.Stopped {
		n += 2
	}
	l = m.Commission.Size()
	n += 1 + l + sovGenesis(uint64(l))
	return n
}

func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Params != nil {
		l = m.Params.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	if m.StartMinterNonce != 0 {
		n += 1 + sovGenesis(uint64(m.StartMinterNonce))
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StartThreshold", wireType)
			}
			m.StartThreshold = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.StartThreshold |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinterAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MinterAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BridgeChainId", wireType)
			}
			m.BridgeChainId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BridgeChainId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignedValsetsWindow", wireType)
			}
			m.SignedValsetsWindow = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SignedValsetsWindow |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignedBatchesWindow", wireType)
			}
			m.SignedBatchesWindow = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SignedBatchesWindow |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignedClaimsWindow", wireType)
			}
			m.SignedClaimsWindow = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SignedClaimsWindow |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SlashFractionValset", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.SlashFractionValset.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SlashFractionBatch", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.SlashFractionBatch.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SlashFractionClaim", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.SlashFractionClaim.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SlashFractionConflictingClaim", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.SlashFractionConflictingClaim.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 11:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Stopped", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Stopped = bool(v != 0)
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Commission", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Commission.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Params == nil {
				m.Params = &Params{}
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StartMinterNonce", wireType)
			}
			m.StartMinterNonce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.StartMinterNonce |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
