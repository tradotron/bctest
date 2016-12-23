// Code generated by protoc-gen-go.
// source: common/common.proto
// DO NOT EDIT!

/*
Package common is a generated protocol buffer package.

It is generated from these files:
	common/common.proto
	common/configuration.proto

It has these top-level messages:
	Header
	ChainHeader
	SignatureHeader
	Payload
	Envelope
	Block
	BlockHeader
	BlockData
	BlockMetadata
	ConfigurationEnvelope
	SignedConfigurationItem
	ConfigurationItem
	ConfigurationSignature
	Policy
	SignaturePolicyEnvelope
	SignaturePolicy
*/
package common

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// These status codes are intended to resemble selected HTTP status codes
type Status int32

const (
	Status_UNKNOWN                  Status = 0
	Status_SUCCESS                  Status = 200
	Status_BAD_REQUEST              Status = 400
	Status_FORBIDDEN                Status = 403
	Status_NOT_FOUND                Status = 404
	Status_REQUEST_ENTITY_TOO_LARGE Status = 413
	Status_INTERNAL_SERVER_ERROR    Status = 500
	Status_SERVICE_UNAVAILABLE      Status = 503
)

var Status_name = map[int32]string{
	0:   "UNKNOWN",
	200: "SUCCESS",
	400: "BAD_REQUEST",
	403: "FORBIDDEN",
	404: "NOT_FOUND",
	413: "REQUEST_ENTITY_TOO_LARGE",
	500: "INTERNAL_SERVER_ERROR",
	503: "SERVICE_UNAVAILABLE",
}
var Status_value = map[string]int32{
	"UNKNOWN":                  0,
	"SUCCESS":                  200,
	"BAD_REQUEST":              400,
	"FORBIDDEN":                403,
	"NOT_FOUND":                404,
	"REQUEST_ENTITY_TOO_LARGE": 413,
	"INTERNAL_SERVER_ERROR":    500,
	"SERVICE_UNAVAILABLE":      503,
}

func (x Status) String() string {
	return proto.EnumName(Status_name, int32(x))
}
func (Status) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type HeaderType int32

const (
	HeaderType_MESSAGE                   HeaderType = 0
	HeaderType_CONFIGURATION_TRANSACTION HeaderType = 1
	HeaderType_CONFIGURATION_ITEM        HeaderType = 2
	HeaderType_ENDORSER_TRANSACTION      HeaderType = 3
	HeaderType_ORDERER_TRANSACTION       HeaderType = 4
)

var HeaderType_name = map[int32]string{
	0: "MESSAGE",
	1: "CONFIGURATION_TRANSACTION",
	2: "CONFIGURATION_ITEM",
	3: "ENDORSER_TRANSACTION",
	4: "ORDERER_TRANSACTION",
}
var HeaderType_value = map[string]int32{
	"MESSAGE":                   0,
	"CONFIGURATION_TRANSACTION": 1,
	"CONFIGURATION_ITEM":        2,
	"ENDORSER_TRANSACTION":      3,
	"ORDERER_TRANSACTION":       4,
}

func (x HeaderType) String() string {
	return proto.EnumName(HeaderType_name, int32(x))
}
func (HeaderType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

// This enum enlist indexes of the block metadata array
type BlockMetadataIndex int32

const (
	BlockMetadataIndex_SIGNATURES          BlockMetadataIndex = 0
	BlockMetadataIndex_LAST_CONFIGURATION  BlockMetadataIndex = 1
	BlockMetadataIndex_TRANSACTIONS_FILTER BlockMetadataIndex = 2
)

var BlockMetadataIndex_name = map[int32]string{
	0: "SIGNATURES",
	1: "LAST_CONFIGURATION",
	2: "TRANSACTIONS_FILTER",
}
var BlockMetadataIndex_value = map[string]int32{
	"SIGNATURES":          0,
	"LAST_CONFIGURATION":  1,
	"TRANSACTIONS_FILTER": 2,
}

func (x BlockMetadataIndex) String() string {
	return proto.EnumName(BlockMetadataIndex_name, int32(x))
}
func (BlockMetadataIndex) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type Header struct {
	ChainHeader     *ChainHeader     `protobuf:"bytes,1,opt,name=chainHeader" json:"chainHeader,omitempty"`
	SignatureHeader *SignatureHeader `protobuf:"bytes,2,opt,name=signatureHeader" json:"signatureHeader,omitempty"`
}

func (m *Header) Reset()                    { *m = Header{} }
func (m *Header) String() string            { return proto.CompactTextString(m) }
func (*Header) ProtoMessage()               {}
func (*Header) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Header) GetChainHeader() *ChainHeader {
	if m != nil {
		return m.ChainHeader
	}
	return nil
}

func (m *Header) GetSignatureHeader() *SignatureHeader {
	if m != nil {
		return m.SignatureHeader
	}
	return nil
}

// Header is a generic replay prevention and identity message to include in a signed payload
type ChainHeader struct {
	Type int32 `protobuf:"varint,1,opt,name=type" json:"type,omitempty"`
	// Version indicates message protocol version
	Version int32 `protobuf:"varint,2,opt,name=version" json:"version,omitempty"`
	// Timestamp is the local time when the message was created
	// by the sender
	Timestamp *google_protobuf.Timestamp `protobuf:"bytes,3,opt,name=timestamp" json:"timestamp,omitempty"`
	// Identifier of the chain this message is bound for
	ChainID string `protobuf:"bytes,4,opt,name=chainID" json:"chainID,omitempty"`
	// An unique identifier that is used end-to-end.
	//  -  set by higher layers such as end user or SDK
	//  -  passed to the endorser (which will check for uniqueness)
	//  -  as the header is passed along unchanged, it will be
	//     be retrieved by the committer (uniqueness check here as well)
	//  -  to be stored in the ledger
	TxID string `protobuf:"bytes,5,opt,name=txID" json:"txID,omitempty"`
	// The epoch in which this header was generated, where epoch is defined based on block height
	// Epoch in which the response has been generated. This field identifies a
	// logical window of time. A proposal response is accepted by a peer only if
	// two conditions hold:
	// 1. the epoch specified in the message is the current epoch
	// 2. this message has been only seen once during this epoch (i.e. it hasn't
	//    been replayed)
	Epoch uint64 `protobuf:"varint,6,opt,name=epoch" json:"epoch,omitempty"`
	// Extension that may be attached based on the header type
	Extension []byte `protobuf:"bytes,7,opt,name=extension,proto3" json:"extension,omitempty"`
}

func (m *ChainHeader) Reset()                    { *m = ChainHeader{} }
func (m *ChainHeader) String() string            { return proto.CompactTextString(m) }
func (*ChainHeader) ProtoMessage()               {}
func (*ChainHeader) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ChainHeader) GetTimestamp() *google_protobuf.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

type SignatureHeader struct {
	// Creator of the message, specified as a certificate chain
	Creator []byte `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	// Arbitrary number that may only be used once. Can be used to detect replay attacks.
	Nonce []byte `protobuf:"bytes,2,opt,name=nonce,proto3" json:"nonce,omitempty"`
}

func (m *SignatureHeader) Reset()                    { *m = SignatureHeader{} }
func (m *SignatureHeader) String() string            { return proto.CompactTextString(m) }
func (*SignatureHeader) ProtoMessage()               {}
func (*SignatureHeader) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

// Payload is the message contents (and header to allow for signing)
type Payload struct {
	// Header is included to provide identity and prevent replay
	Header *Header `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	// Data, the encoding of which is defined by the type in the header
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *Payload) Reset()                    { *m = Payload{} }
func (m *Payload) String() string            { return proto.CompactTextString(m) }
func (*Payload) ProtoMessage()               {}
func (*Payload) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Payload) GetHeader() *Header {
	if m != nil {
		return m.Header
	}
	return nil
}

// Envelope wraps a Payload with a signature so that the message may be authenticated
type Envelope struct {
	// A marshaled Payload
	Payload []byte `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	// A signature by the creator specified in the Payload header
	Signature []byte `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (m *Envelope) Reset()                    { *m = Envelope{} }
func (m *Envelope) String() string            { return proto.CompactTextString(m) }
func (*Envelope) ProtoMessage()               {}
func (*Envelope) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

// This is finalized block structure to be shared among the orderer and peer
// Note that the BlockHeader chains to the previous BlockHeader, and the BlockData hash is embedded
// in the BlockHeader.  This makes it natural and obvious that the Data is included in the hash, but
// the Metadata is not.
type Block struct {
	Header   *BlockHeader   `protobuf:"bytes,1,opt,name=Header" json:"Header,omitempty"`
	Data     *BlockData     `protobuf:"bytes,2,opt,name=Data" json:"Data,omitempty"`
	Metadata *BlockMetadata `protobuf:"bytes,3,opt,name=Metadata" json:"Metadata,omitempty"`
}

func (m *Block) Reset()                    { *m = Block{} }
func (m *Block) String() string            { return proto.CompactTextString(m) }
func (*Block) ProtoMessage()               {}
func (*Block) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Block) GetHeader() *BlockHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Block) GetData() *BlockData {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Block) GetMetadata() *BlockMetadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type BlockHeader struct {
	Number       uint64 `protobuf:"varint,1,opt,name=Number" json:"Number,omitempty"`
	PreviousHash []byte `protobuf:"bytes,2,opt,name=PreviousHash,proto3" json:"PreviousHash,omitempty"`
	DataHash     []byte `protobuf:"bytes,3,opt,name=DataHash,proto3" json:"DataHash,omitempty"`
}

func (m *BlockHeader) Reset()                    { *m = BlockHeader{} }
func (m *BlockHeader) String() string            { return proto.CompactTextString(m) }
func (*BlockHeader) ProtoMessage()               {}
func (*BlockHeader) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type BlockData struct {
	Data [][]byte `protobuf:"bytes,1,rep,name=Data,proto3" json:"Data,omitempty"`
}

func (m *BlockData) Reset()                    { *m = BlockData{} }
func (m *BlockData) String() string            { return proto.CompactTextString(m) }
func (*BlockData) ProtoMessage()               {}
func (*BlockData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

type BlockMetadata struct {
	Metadata [][]byte `protobuf:"bytes,1,rep,name=Metadata,proto3" json:"Metadata,omitempty"`
}

func (m *BlockMetadata) Reset()                    { *m = BlockMetadata{} }
func (m *BlockMetadata) String() string            { return proto.CompactTextString(m) }
func (*BlockMetadata) ProtoMessage()               {}
func (*BlockMetadata) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func init() {
	proto.RegisterType((*Header)(nil), "common.Header")
	proto.RegisterType((*ChainHeader)(nil), "common.ChainHeader")
	proto.RegisterType((*SignatureHeader)(nil), "common.SignatureHeader")
	proto.RegisterType((*Payload)(nil), "common.Payload")
	proto.RegisterType((*Envelope)(nil), "common.Envelope")
	proto.RegisterType((*Block)(nil), "common.Block")
	proto.RegisterType((*BlockHeader)(nil), "common.BlockHeader")
	proto.RegisterType((*BlockData)(nil), "common.BlockData")
	proto.RegisterType((*BlockMetadata)(nil), "common.BlockMetadata")
	proto.RegisterEnum("common.Status", Status_name, Status_value)
	proto.RegisterEnum("common.HeaderType", HeaderType_name, HeaderType_value)
	proto.RegisterEnum("common.BlockMetadataIndex", BlockMetadataIndex_name, BlockMetadataIndex_value)
}

func init() { proto.RegisterFile("common/common.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 771 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x5c, 0x54, 0xcd, 0x6e, 0xf2, 0x46,
	0x14, 0x8d, 0x63, 0x7e, 0xc2, 0x35, 0xcd, 0xe7, 0x0e, 0xdf, 0x97, 0xb8, 0xa8, 0x51, 0x90, 0xa5,
	0x56, 0x28, 0x51, 0x41, 0x4d, 0x55, 0xa9, 0x5b, 0x83, 0x07, 0x62, 0x95, 0x8c, 0xd3, 0xb1, 0x9d,
	0xaa, 0xdd, 0x58, 0x06, 0x26, 0x80, 0x0a, 0x36, 0x32, 0x26, 0x4a, 0xb6, 0xed, 0xbe, 0xaa, 0xd4,
	0x6e, 0xfb, 0x02, 0x7d, 0x92, 0xbe, 0x41, 0x5f, 0xa4, 0x52, 0xb7, 0xd5, 0x8c, 0x6d, 0xc0, 0xac,
	0x98, 0x73, 0xcf, 0xfd, 0x39, 0xe7, 0xce, 0x60, 0x68, 0x4c, 0xa2, 0xd5, 0x2a, 0x0a, 0xbb, 0xe9,
	0x4f, 0x67, 0x1d, 0x47, 0x49, 0x84, 0x2a, 0x29, 0x6a, 0x5e, 0xcf, 0xa2, 0x68, 0xb6, 0x64, 0x5d,
	0x11, 0x1d, 0x6f, 0x9f, 0xbb, 0xc9, 0x62, 0xc5, 0x36, 0x49, 0xb0, 0x5a, 0xa7, 0x89, 0xfa, 0xcf,
	0x12, 0x54, 0xee, 0x59, 0x30, 0x65, 0x31, 0xfa, 0x1a, 0x94, 0xc9, 0x3c, 0x58, 0x84, 0x29, 0xd4,
	0xa4, 0x96, 0xd4, 0x56, 0xee, 0x1a, 0x9d, 0xac, 0x6f, 0x7f, 0x4f, 0xd1, 0xc3, 0x3c, 0x64, 0xc0,
	0xbb, 0xcd, 0x62, 0x16, 0x06, 0xc9, 0x36, 0x66, 0x59, 0xe9, 0xa9, 0x28, 0xbd, 0xcc, 0x4b, 0x9d,
	0x22, 0x4d, 0x8f, 0xf3, 0xf5, 0x7f, 0x24, 0x50, 0x0e, 0xfa, 0x23, 0x04, 0xa5, 0xe4, 0x6d, 0xcd,
	0x84, 0x84, 0x32, 0x15, 0x67, 0xa4, 0x41, 0xf5, 0x85, 0xc5, 0x9b, 0x45, 0x14, 0x8a, 0xf6, 0x65,
	0x9a, 0x43, 0xf4, 0x0d, 0xd4, 0x76, 0xae, 0x34, 0x59, 0x8c, 0x6e, 0x76, 0x52, 0xdf, 0x9d, 0xdc,
	0x77, 0xc7, 0xcd, 0x33, 0xe8, 0x3e, 0x99, 0xf7, 0x14, 0x4e, 0x2c, 0x53, 0x2b, 0xb5, 0xa4, 0x76,
	0x8d, 0xe6, 0x50, 0x28, 0x78, 0xb5, 0x4c, 0xad, 0x2c, 0xc2, 0xe2, 0x8c, 0xde, 0x43, 0x99, 0xad,
	0xa3, 0xc9, 0x5c, 0xab, 0xb4, 0xa4, 0x76, 0x89, 0xa6, 0x00, 0x7d, 0x0a, 0x35, 0xf6, 0x9a, 0xb0,
	0x50, 0x28, 0xab, 0xb6, 0xa4, 0x76, 0x9d, 0xee, 0x03, 0xba, 0x01, 0xef, 0x8e, 0xdc, 0x8b, 0xa1,
	0x31, 0x0b, 0x92, 0x28, 0x5d, 0x71, 0x9d, 0xe6, 0x90, 0x0f, 0x08, 0xa3, 0x70, 0xc2, 0x84, 0xc1,
	0x3a, 0x4d, 0x81, 0x8e, 0xa1, 0xfa, 0x18, 0xbc, 0x2d, 0xa3, 0x60, 0x8a, 0x3e, 0x87, 0xca, 0xfc,
	0xf0, 0x72, 0xce, 0xf3, 0x0d, 0x67, 0x8b, 0xcd, 0x58, 0xae, 0x7e, 0x1a, 0x24, 0x41, 0xd6, 0x47,
	0x9c, 0xf5, 0x1e, 0x9c, 0xe1, 0xf0, 0x85, 0x2d, 0xa3, 0x74, 0x97, 0xeb, 0xb4, 0x65, 0x2e, 0x21,
	0x83, 0xdc, 0xcd, 0xee, 0x72, 0xb2, 0xf2, 0x7d, 0x40, 0xff, 0x55, 0x82, 0x72, 0x6f, 0x19, 0x4d,
	0x7e, 0x42, 0xb7, 0xf9, 0xab, 0x39, 0x7e, 0x26, 0x82, 0xce, 0xe5, 0x64, 0x8e, 0x3f, 0x83, 0x92,
	0x99, 0xcb, 0x51, 0xee, 0x3e, 0x2e, 0xa4, 0x72, 0x82, 0x0a, 0x1a, 0x7d, 0x09, 0x67, 0x0f, 0x2c,
	0x09, 0x84, 0xf2, 0xf4, 0x1a, 0x3f, 0x14, 0x52, 0x73, 0x92, 0xee, 0xd2, 0x74, 0x06, 0xca, 0xc1,
	0x40, 0x74, 0x01, 0x15, 0xb2, 0x5d, 0x8d, 0x33, 0x55, 0x25, 0x9a, 0x21, 0xa4, 0x43, 0xfd, 0x31,
	0x66, 0x2f, 0x8b, 0x68, 0xbb, 0xb9, 0x0f, 0x36, 0xf3, 0xcc, 0x58, 0x21, 0x86, 0x9a, 0x70, 0xc6,
	0x55, 0x08, 0x5e, 0x16, 0xfc, 0x0e, 0xeb, 0xd7, 0x50, 0xdb, 0x89, 0xe5, 0xcb, 0x15, 0x6e, 0xa4,
	0x96, 0xcc, 0x97, 0xcb, 0xcf, 0xfa, 0x2d, 0x7c, 0x54, 0x90, 0xc8, 0xbb, 0xed, 0xbc, 0xa4, 0x89,
	0x3b, 0x7c, 0xf3, 0x97, 0x04, 0x15, 0x27, 0x09, 0x92, 0xed, 0x06, 0x29, 0x50, 0xf5, 0xc8, 0xb7,
	0xc4, 0xfe, 0x9e, 0xa8, 0x27, 0xa8, 0x0e, 0x55, 0xc7, 0xeb, 0xf7, 0xb1, 0xe3, 0xa8, 0x7f, 0x4b,
	0x48, 0x05, 0xa5, 0x67, 0x98, 0x3e, 0xc5, 0xdf, 0x79, 0xd8, 0x71, 0xd5, 0xdf, 0x64, 0x74, 0x0e,
	0xb5, 0x81, 0x4d, 0x7b, 0x96, 0x69, 0x62, 0xa2, 0xfe, 0x2e, 0x30, 0xb1, 0x5d, 0x7f, 0x60, 0x7b,
	0xc4, 0x54, 0xff, 0x90, 0xd1, 0x15, 0x68, 0x59, 0xb6, 0x8f, 0x89, 0x6b, 0xb9, 0x3f, 0xf8, 0xae,
	0x6d, 0xfb, 0x23, 0x83, 0x0e, 0xb1, 0xfa, 0xa7, 0x8c, 0x9a, 0xf0, 0xc1, 0x22, 0x2e, 0xa6, 0xc4,
	0x18, 0xf9, 0x0e, 0xa6, 0x4f, 0x98, 0xfa, 0x98, 0x52, 0x9b, 0xaa, 0xff, 0xca, 0x48, 0x83, 0x06,
	0x0f, 0x59, 0x7d, 0xec, 0x7b, 0xc4, 0x78, 0x32, 0xac, 0x91, 0xd1, 0x1b, 0x61, 0xf5, 0x3f, 0xf9,
	0xe6, 0x17, 0x09, 0x20, 0xdd, 0xae, 0xcb, 0xff, 0x85, 0x0a, 0x54, 0x1f, 0xb0, 0xe3, 0x18, 0x43,
	0xac, 0x9e, 0xa0, 0x2b, 0xf8, 0xa4, 0x6f, 0x93, 0x81, 0x35, 0xf4, 0xa8, 0xe1, 0x5a, 0x36, 0xf1,
	0x5d, 0x6a, 0x10, 0xc7, 0xe8, 0xf3, 0xb3, 0x2a, 0xa1, 0x0b, 0x40, 0x45, 0xda, 0x72, 0xf1, 0x83,
	0x7a, 0x8a, 0x34, 0x78, 0x8f, 0x89, 0x69, 0x53, 0x07, 0xd3, 0x42, 0x85, 0x8c, 0x2e, 0xa1, 0x61,
	0x53, 0x13, 0xd3, 0x23, 0xa2, 0x74, 0xe3, 0x01, 0x2a, 0xec, 0xd7, 0x0a, 0xa7, 0xec, 0x15, 0x9d,
	0x03, 0x38, 0xd6, 0x90, 0x18, 0xae, 0x47, 0xb1, 0xa3, 0x9e, 0xf0, 0x81, 0x23, 0xc3, 0x71, 0xfd,
	0xc2, 0x54, 0x55, 0xe2, 0x6d, 0x0f, 0xda, 0x39, 0xfe, 0xc0, 0x1a, 0xb9, 0x98, 0xaa, 0xa7, 0xbd,
	0x2f, 0x7e, 0xbc, 0x9d, 0x2d, 0x92, 0xf9, 0x76, 0xcc, 0xdf, 0x59, 0x77, 0xfe, 0xb6, 0x66, 0xf1,
	0x92, 0x4d, 0x67, 0x2c, 0xee, 0x3e, 0x07, 0xe3, 0x78, 0x31, 0x49, 0x3f, 0x9b, 0x9b, 0xec, 0xd3,
	0x3a, 0xae, 0x08, 0xf8, 0xd5, 0xff, 0x01, 0x00, 0x00, 0xff, 0xff, 0x8a, 0x69, 0x37, 0x92, 0x72,
	0x05, 0x00, 0x00,
}
