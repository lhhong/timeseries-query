// Code generated by protoc-gen-go. DO NOT EDIT.
// source: sectionInfo.proto

package repository

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type SectionInfo struct {
	Groupname            string   `protobuf:"bytes,1,opt,name=Groupname,proto3" json:"Groupname,omitempty"`
	Series               string   `protobuf:"bytes,2,opt,name=Series,proto3" json:"Series,omitempty"`
	Smooth               int32    `protobuf:"varint,3,opt,name=Smooth,proto3" json:"Smooth,omitempty"`
	StartSeq             int64    `protobuf:"varint,4,opt,name=StartSeq,proto3" json:"StartSeq,omitempty"`
	Sign                 int32    `protobuf:"varint,5,opt,name=Sign,proto3" json:"Sign,omitempty"`
	Height               float64  `protobuf:"fixed64,6,opt,name=Height,proto3" json:"Height,omitempty"`
	Width                int64    `protobuf:"varint,7,opt,name=Width,proto3" json:"Width,omitempty"`
	NextSeq              int64    `protobuf:"varint,8,opt,name=NextSeq,proto3" json:"NextSeq,omitempty"`
	PrevSeq              int64    `protobuf:"varint,9,opt,name=PrevSeq,proto3" json:"PrevSeq,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SectionInfo) Reset()         { *m = SectionInfo{} }
func (m *SectionInfo) String() string { return proto.CompactTextString(m) }
func (*SectionInfo) ProtoMessage()    {}
func (*SectionInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_5f43b93927de6136, []int{0}
}

func (m *SectionInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SectionInfo.Unmarshal(m, b)
}
func (m *SectionInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SectionInfo.Marshal(b, m, deterministic)
}
func (m *SectionInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SectionInfo.Merge(m, src)
}
func (m *SectionInfo) XXX_Size() int {
	return xxx_messageInfo_SectionInfo.Size(m)
}
func (m *SectionInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_SectionInfo.DiscardUnknown(m)
}

var xxx_messageInfo_SectionInfo proto.InternalMessageInfo

func (m *SectionInfo) GetGroupname() string {
	if m != nil {
		return m.Groupname
	}
	return ""
}

func (m *SectionInfo) GetSeries() string {
	if m != nil {
		return m.Series
	}
	return ""
}

func (m *SectionInfo) GetSmooth() int32 {
	if m != nil {
		return m.Smooth
	}
	return 0
}

func (m *SectionInfo) GetStartSeq() int64 {
	if m != nil {
		return m.StartSeq
	}
	return 0
}

func (m *SectionInfo) GetSign() int32 {
	if m != nil {
		return m.Sign
	}
	return 0
}

func (m *SectionInfo) GetHeight() float64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *SectionInfo) GetWidth() int64 {
	if m != nil {
		return m.Width
	}
	return 0
}

func (m *SectionInfo) GetNextSeq() int64 {
	if m != nil {
		return m.NextSeq
	}
	return 0
}

func (m *SectionInfo) GetPrevSeq() int64 {
	if m != nil {
		return m.PrevSeq
	}
	return 0
}

func init() {
	proto.RegisterType((*SectionInfo)(nil), "repository.SectionInfo")
}

func init() { proto.RegisterFile("sectionInfo.proto", fileDescriptor_5f43b93927de6136) }

var fileDescriptor_5f43b93927de6136 = []byte{
	// 207 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x90, 0xb1, 0x4e, 0xc3, 0x40,
	0x0c, 0x86, 0x75, 0xb4, 0x49, 0x1b, 0x33, 0x61, 0x21, 0x64, 0x21, 0x86, 0x88, 0xe9, 0x26, 0x16,
	0x1e, 0x02, 0x58, 0x10, 0xba, 0x1b, 0x98, 0x0b, 0x98, 0xe6, 0x86, 0x9e, 0xc3, 0xc5, 0x20, 0x78,
	0x68, 0xde, 0x01, 0xc5, 0x3d, 0xda, 0xcd, 0xdf, 0xff, 0xd9, 0xff, 0x60, 0x38, 0x9b, 0xf8, 0x55,
	0x93, 0xe4, 0x87, 0xfc, 0x2e, 0x37, 0x63, 0x11, 0x15, 0x84, 0xc2, 0xa3, 0x4c, 0x49, 0xa5, 0xfc,
	0x5c, 0xff, 0x3a, 0x38, 0x8d, 0xc7, 0x0d, 0xbc, 0x82, 0xee, 0xae, 0xc8, 0xe7, 0x98, 0x37, 0x3b,
	0x26, 0xd7, 0x3b, 0xdf, 0x85, 0x63, 0x80, 0x17, 0xd0, 0x46, 0x2e, 0x89, 0x27, 0x3a, 0x31, 0x55,
	0xc9, 0xf2, 0x9d, 0x88, 0x0e, 0xb4, 0xe8, 0x9d, 0x6f, 0x42, 0x25, 0xbc, 0x84, 0x75, 0xd4, 0x4d,
	0xd1, 0xc8, 0x1f, 0xb4, 0xec, 0x9d, 0x5f, 0x84, 0x03, 0x23, 0xc2, 0x32, 0xa6, 0x6d, 0xa6, 0xc6,
	0x2e, 0x6c, 0x9e, 0x7b, 0xee, 0x39, 0x6d, 0x07, 0xa5, 0xb6, 0x77, 0xde, 0x85, 0x4a, 0x78, 0x0e,
	0xcd, 0x73, 0x7a, 0xd3, 0x81, 0x56, 0x56, 0xb2, 0x07, 0x24, 0x58, 0x3d, 0xf2, 0xb7, 0x95, 0xaf,
	0x2d, 0xff, 0xc7, 0xd9, 0x3c, 0x15, 0xfe, 0x9a, 0x4d, 0xb7, 0x37, 0x15, 0x5f, 0x5a, 0x7b, 0xc1,
	0xed, 0x5f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x64, 0x24, 0x92, 0x1b, 0x17, 0x01, 0x00, 0x00,
}
