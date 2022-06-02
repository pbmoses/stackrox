// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: storage/service_identity.proto

package storage

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ServiceType int32

const (
	ServiceType_UNKNOWN_SERVICE           ServiceType = 0
	ServiceType_SENSOR_SERVICE            ServiceType = 1
	ServiceType_CENTRAL_SERVICE           ServiceType = 2
	ServiceType_REMOTE_SERVICE            ServiceType = 3
	ServiceType_COLLECTOR_SERVICE         ServiceType = 4
	ServiceType_MONITORING_UI_SERVICE     ServiceType = 5
	ServiceType_MONITORING_DB_SERVICE     ServiceType = 6
	ServiceType_MONITORING_CLIENT_SERVICE ServiceType = 7
	ServiceType_BENCHMARK_SERVICE         ServiceType = 8
	ServiceType_SCANNER_SERVICE           ServiceType = 9
	ServiceType_SCANNER_DB_SERVICE        ServiceType = 10
	ServiceType_ADMISSION_CONTROL_SERVICE ServiceType = 11
)

var ServiceType_name = map[int32]string{
	0:  "UNKNOWN_SERVICE",
	1:  "SENSOR_SERVICE",
	2:  "CENTRAL_SERVICE",
	3:  "REMOTE_SERVICE",
	4:  "COLLECTOR_SERVICE",
	5:  "MONITORING_UI_SERVICE",
	6:  "MONITORING_DB_SERVICE",
	7:  "MONITORING_CLIENT_SERVICE",
	8:  "BENCHMARK_SERVICE",
	9:  "SCANNER_SERVICE",
	10: "SCANNER_DB_SERVICE",
	11: "ADMISSION_CONTROL_SERVICE",
}

var ServiceType_value = map[string]int32{
	"UNKNOWN_SERVICE":           0,
	"SENSOR_SERVICE":            1,
	"CENTRAL_SERVICE":           2,
	"REMOTE_SERVICE":            3,
	"COLLECTOR_SERVICE":         4,
	"MONITORING_UI_SERVICE":     5,
	"MONITORING_DB_SERVICE":     6,
	"MONITORING_CLIENT_SERVICE": 7,
	"BENCHMARK_SERVICE":         8,
	"SCANNER_SERVICE":           9,
	"SCANNER_DB_SERVICE":        10,
	"ADMISSION_CONTROL_SERVICE": 11,
}

func (x ServiceType) String() string {
	return proto.EnumName(ServiceType_name, int32(x))
}

func (ServiceType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_a988b93c2073ff63, []int{0}
}

type ServiceIdentity struct {
	SerialStr string `protobuf:"bytes,4,opt,name=serial_str,json=serialStr,proto3" json:"serial_str,omitempty" sql:"pk"`
	// Types that are valid to be assigned to Srl:
	//	*ServiceIdentity_Serial
	Srl                  isServiceIdentity_Srl `protobuf_oneof:"srl"`
	Id                   string                `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Type                 ServiceType           `protobuf:"varint,3,opt,name=type,proto3,enum=storage.ServiceType" json:"type,omitempty"`
	InitBundleId         string                `protobuf:"bytes,5,opt,name=init_bundle_id,json=initBundleId,proto3" json:"init_bundle_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *ServiceIdentity) Reset()         { *m = ServiceIdentity{} }
func (m *ServiceIdentity) String() string { return proto.CompactTextString(m) }
func (*ServiceIdentity) ProtoMessage()    {}
func (*ServiceIdentity) Descriptor() ([]byte, []int) {
	return fileDescriptor_a988b93c2073ff63, []int{0}
}
func (m *ServiceIdentity) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ServiceIdentity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ServiceIdentity.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ServiceIdentity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServiceIdentity.Merge(m, src)
}
func (m *ServiceIdentity) XXX_Size() int {
	return m.Size()
}
func (m *ServiceIdentity) XXX_DiscardUnknown() {
	xxx_messageInfo_ServiceIdentity.DiscardUnknown(m)
}

var xxx_messageInfo_ServiceIdentity proto.InternalMessageInfo

type isServiceIdentity_Srl interface {
	isServiceIdentity_Srl()
	MarshalTo([]byte) (int, error)
	Size() int
	Clone() isServiceIdentity_Srl
}

type ServiceIdentity_Serial struct {
	Serial int64 `protobuf:"varint,1,opt,name=serial,proto3,oneof" json:"serial,omitempty"`
}

func (*ServiceIdentity_Serial) isServiceIdentity_Srl() {}
func (m *ServiceIdentity_Serial) Clone() isServiceIdentity_Srl {
	if m == nil {
		return nil
	}
	cloned := new(ServiceIdentity_Serial)
	*cloned = *m

	return cloned
}

func (m *ServiceIdentity) GetSrl() isServiceIdentity_Srl {
	if m != nil {
		return m.Srl
	}
	return nil
}

func (m *ServiceIdentity) GetSerialStr() string {
	if m != nil {
		return m.SerialStr
	}
	return ""
}

// Deprecated: Do not use.
func (m *ServiceIdentity) GetSerial() int64 {
	if x, ok := m.GetSrl().(*ServiceIdentity_Serial); ok {
		return x.Serial
	}
	return 0
}

func (m *ServiceIdentity) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ServiceIdentity) GetType() ServiceType {
	if m != nil {
		return m.Type
	}
	return ServiceType_UNKNOWN_SERVICE
}

func (m *ServiceIdentity) GetInitBundleId() string {
	if m != nil {
		return m.InitBundleId
	}
	return ""
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*ServiceIdentity) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*ServiceIdentity_Serial)(nil),
	}
}

func (m *ServiceIdentity) MessageClone() proto.Message {
	return m.Clone()
}
func (m *ServiceIdentity) Clone() *ServiceIdentity {
	if m == nil {
		return nil
	}
	cloned := new(ServiceIdentity)
	*cloned = *m

	if m.Srl != nil {
		cloned.Srl = m.Srl.Clone()
	}
	return cloned
}

type ServiceCertificate struct {
	CertPem              []byte   `protobuf:"bytes,1,opt,name=cert_pem,json=certPem,proto3" json:"cert_pem,omitempty"`
	KeyPem               []byte   `protobuf:"bytes,2,opt,name=key_pem,json=keyPem,proto3" json:"key_pem,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ServiceCertificate) Reset()         { *m = ServiceCertificate{} }
func (m *ServiceCertificate) String() string { return proto.CompactTextString(m) }
func (*ServiceCertificate) ProtoMessage()    {}
func (*ServiceCertificate) Descriptor() ([]byte, []int) {
	return fileDescriptor_a988b93c2073ff63, []int{1}
}
func (m *ServiceCertificate) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ServiceCertificate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ServiceCertificate.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ServiceCertificate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServiceCertificate.Merge(m, src)
}
func (m *ServiceCertificate) XXX_Size() int {
	return m.Size()
}
func (m *ServiceCertificate) XXX_DiscardUnknown() {
	xxx_messageInfo_ServiceCertificate.DiscardUnknown(m)
}

var xxx_messageInfo_ServiceCertificate proto.InternalMessageInfo

func (m *ServiceCertificate) GetCertPem() []byte {
	if m != nil {
		return m.CertPem
	}
	return nil
}

func (m *ServiceCertificate) GetKeyPem() []byte {
	if m != nil {
		return m.KeyPem
	}
	return nil
}

func (m *ServiceCertificate) MessageClone() proto.Message {
	return m.Clone()
}
func (m *ServiceCertificate) Clone() *ServiceCertificate {
	if m == nil {
		return nil
	}
	cloned := new(ServiceCertificate)
	*cloned = *m

	if m.CertPem != nil {
		cloned.CertPem = make([]byte, len(m.CertPem))
		copy(cloned.CertPem, m.CertPem)
	}
	if m.KeyPem != nil {
		cloned.KeyPem = make([]byte, len(m.KeyPem))
		copy(cloned.KeyPem, m.KeyPem)
	}
	return cloned
}

type TypedServiceCertificate struct {
	ServiceType          ServiceType         `protobuf:"varint,1,opt,name=service_type,json=serviceType,proto3,enum=storage.ServiceType" json:"service_type,omitempty"`
	Cert                 *ServiceCertificate `protobuf:"bytes,2,opt,name=cert,proto3" json:"cert,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *TypedServiceCertificate) Reset()         { *m = TypedServiceCertificate{} }
func (m *TypedServiceCertificate) String() string { return proto.CompactTextString(m) }
func (*TypedServiceCertificate) ProtoMessage()    {}
func (*TypedServiceCertificate) Descriptor() ([]byte, []int) {
	return fileDescriptor_a988b93c2073ff63, []int{2}
}
func (m *TypedServiceCertificate) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TypedServiceCertificate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TypedServiceCertificate.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TypedServiceCertificate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TypedServiceCertificate.Merge(m, src)
}
func (m *TypedServiceCertificate) XXX_Size() int {
	return m.Size()
}
func (m *TypedServiceCertificate) XXX_DiscardUnknown() {
	xxx_messageInfo_TypedServiceCertificate.DiscardUnknown(m)
}

var xxx_messageInfo_TypedServiceCertificate proto.InternalMessageInfo

func (m *TypedServiceCertificate) GetServiceType() ServiceType {
	if m != nil {
		return m.ServiceType
	}
	return ServiceType_UNKNOWN_SERVICE
}

func (m *TypedServiceCertificate) GetCert() *ServiceCertificate {
	if m != nil {
		return m.Cert
	}
	return nil
}

func (m *TypedServiceCertificate) MessageClone() proto.Message {
	return m.Clone()
}
func (m *TypedServiceCertificate) Clone() *TypedServiceCertificate {
	if m == nil {
		return nil
	}
	cloned := new(TypedServiceCertificate)
	*cloned = *m

	cloned.Cert = m.Cert.Clone()
	return cloned
}

type TypedServiceCertificateSet struct {
	CaPem                []byte                     `protobuf:"bytes,1,opt,name=ca_pem,json=caPem,proto3" json:"ca_pem,omitempty"`
	ServiceCerts         []*TypedServiceCertificate `protobuf:"bytes,2,rep,name=service_certs,json=serviceCerts,proto3" json:"service_certs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *TypedServiceCertificateSet) Reset()         { *m = TypedServiceCertificateSet{} }
func (m *TypedServiceCertificateSet) String() string { return proto.CompactTextString(m) }
func (*TypedServiceCertificateSet) ProtoMessage()    {}
func (*TypedServiceCertificateSet) Descriptor() ([]byte, []int) {
	return fileDescriptor_a988b93c2073ff63, []int{3}
}
func (m *TypedServiceCertificateSet) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TypedServiceCertificateSet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TypedServiceCertificateSet.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TypedServiceCertificateSet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TypedServiceCertificateSet.Merge(m, src)
}
func (m *TypedServiceCertificateSet) XXX_Size() int {
	return m.Size()
}
func (m *TypedServiceCertificateSet) XXX_DiscardUnknown() {
	xxx_messageInfo_TypedServiceCertificateSet.DiscardUnknown(m)
}

var xxx_messageInfo_TypedServiceCertificateSet proto.InternalMessageInfo

func (m *TypedServiceCertificateSet) GetCaPem() []byte {
	if m != nil {
		return m.CaPem
	}
	return nil
}

func (m *TypedServiceCertificateSet) GetServiceCerts() []*TypedServiceCertificate {
	if m != nil {
		return m.ServiceCerts
	}
	return nil
}

func (m *TypedServiceCertificateSet) MessageClone() proto.Message {
	return m.Clone()
}
func (m *TypedServiceCertificateSet) Clone() *TypedServiceCertificateSet {
	if m == nil {
		return nil
	}
	cloned := new(TypedServiceCertificateSet)
	*cloned = *m

	if m.CaPem != nil {
		cloned.CaPem = make([]byte, len(m.CaPem))
		copy(cloned.CaPem, m.CaPem)
	}
	if m.ServiceCerts != nil {
		cloned.ServiceCerts = make([]*TypedServiceCertificate, len(m.ServiceCerts))
		for idx, v := range m.ServiceCerts {
			cloned.ServiceCerts[idx] = v.Clone()
		}
	}
	return cloned
}

func init() {
	proto.RegisterEnum("storage.ServiceType", ServiceType_name, ServiceType_value)
	proto.RegisterType((*ServiceIdentity)(nil), "storage.ServiceIdentity")
	proto.RegisterType((*ServiceCertificate)(nil), "storage.ServiceCertificate")
	proto.RegisterType((*TypedServiceCertificate)(nil), "storage.TypedServiceCertificate")
	proto.RegisterType((*TypedServiceCertificateSet)(nil), "storage.TypedServiceCertificateSet")
}

func init() { proto.RegisterFile("storage/service_identity.proto", fileDescriptor_a988b93c2073ff63) }

var fileDescriptor_a988b93c2073ff63 = []byte{
	// 562 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x53, 0xcb, 0x6e, 0xd3, 0x4c,
	0x14, 0xee, 0x38, 0xb7, 0xf6, 0x24, 0x7f, 0xfa, 0x33, 0xb4, 0xd4, 0x29, 0x10, 0x22, 0x8b, 0x45,
	0x04, 0x52, 0x22, 0x05, 0x24, 0x24, 0x76, 0xb1, 0x6b, 0x11, 0xab, 0xc9, 0xb8, 0x1a, 0xbb, 0x20,
	0xb1, 0xb1, 0x5c, 0x7b, 0x88, 0xac, 0x5c, 0x1c, 0xec, 0x01, 0x11, 0xb6, 0xbc, 0x04, 0x4f, 0xc1,
	0x4b, 0xb0, 0x61, 0xc9, 0x13, 0x20, 0x14, 0xde, 0x80, 0x27, 0x40, 0xe3, 0xd8, 0x8e, 0x55, 0x51,
	0x76, 0x39, 0xdf, 0x77, 0xce, 0x77, 0x89, 0x35, 0xd0, 0x8e, 0x79, 0x18, 0xb9, 0x53, 0xd6, 0x8f,
	0x59, 0xf4, 0x3e, 0xf0, 0x98, 0x13, 0xf8, 0x6c, 0xc9, 0x03, 0xbe, 0xee, 0xad, 0xa2, 0x90, 0x87,
	0xb8, 0x96, 0xf2, 0xa7, 0x47, 0xd3, 0x70, 0x1a, 0x26, 0x58, 0x5f, 0xfc, 0xda, 0xd2, 0xca, 0x57,
	0x04, 0x87, 0xd6, 0xf6, 0xd2, 0x48, 0x0f, 0xf1, 0x63, 0x80, 0x98, 0x45, 0x81, 0x3b, 0x77, 0x62,
	0x1e, 0xc9, 0xe5, 0x0e, 0xea, 0x1e, 0xa8, 0x8d, 0xdf, 0x3f, 0x1e, 0xec, 0xc7, 0x6f, 0xe7, 0xcf,
	0x95, 0xd5, 0x4c, 0xa1, 0x07, 0x5b, 0xde, 0xe2, 0x11, 0xbe, 0x07, 0xd5, 0xed, 0x20, 0xa3, 0x0e,
	0xea, 0x96, 0x54, 0x49, 0x46, 0xa3, 0x3d, 0x9a, 0x62, 0xb8, 0x09, 0x52, 0xe0, 0xcb, 0x92, 0x90,
	0xa0, 0x52, 0xe0, 0xe3, 0x2e, 0x94, 0xf9, 0x7a, 0xc5, 0xe4, 0x52, 0x07, 0x75, 0x9b, 0x83, 0xa3,
	0x5e, 0x1a, 0xae, 0x97, 0x46, 0xb0, 0xd7, 0x2b, 0x46, 0x93, 0x0d, 0xfc, 0x10, 0x9a, 0xc1, 0x32,
	0xe0, 0xce, 0xd5, 0xbb, 0xa5, 0x3f, 0x17, 0xad, 0xe4, 0x4a, 0xa2, 0xd2, 0x10, 0xa8, 0x9a, 0x80,
	0x86, 0xaf, 0x56, 0xa0, 0x14, 0x47, 0x73, 0x65, 0x04, 0x38, 0x55, 0xd0, 0x58, 0xc4, 0x83, 0x37,
	0x81, 0xe7, 0x72, 0x86, 0x5b, 0xb0, 0xef, 0xb1, 0x88, 0x3b, 0x2b, 0xb6, 0x48, 0xc2, 0x35, 0x68,
	0x4d, 0xcc, 0x17, 0x6c, 0x81, 0x4f, 0xa0, 0x36, 0x63, 0xeb, 0x84, 0x91, 0x12, 0xa6, 0x3a, 0x63,
	0xeb, 0x0b, 0xb6, 0x50, 0x3e, 0x21, 0x38, 0x11, 0x29, 0xfc, 0xbf, 0xe8, 0x3d, 0x83, 0x46, 0xf6,
	0x27, 0x27, 0x25, 0xd0, 0x3f, 0x4a, 0xd4, 0xe3, 0xdd, 0x80, 0xfb, 0x50, 0x16, 0xc6, 0x89, 0x55,
	0x7d, 0x70, 0xf7, 0xfa, 0x41, 0xc1, 0x83, 0x26, 0x8b, 0xca, 0x47, 0x38, 0xbd, 0x21, 0x84, 0xc5,
	0x38, 0x3e, 0x86, 0xaa, 0xe7, 0x16, 0x5a, 0x55, 0x3c, 0x57, 0x74, 0xd2, 0xe1, 0xbf, 0x2c, 0x9e,
	0x10, 0x89, 0x65, 0xa9, 0x53, 0xea, 0xd6, 0x07, 0x9d, 0xdc, 0xee, 0x06, 0x49, 0x9a, 0xb5, 0x12,
	0x58, 0xfc, 0xe8, 0x8b, 0x04, 0xf5, 0x42, 0x13, 0x7c, 0x1b, 0x0e, 0x2f, 0xc9, 0x39, 0x31, 0x5f,
	0x11, 0xc7, 0xd2, 0xe9, 0x4b, 0x43, 0xd3, 0xff, 0xdf, 0xc3, 0x18, 0x9a, 0x96, 0x4e, 0x2c, 0x93,
	0xe6, 0x18, 0x12, 0x8b, 0x9a, 0x4e, 0x6c, 0x3a, 0x1c, 0xe7, 0xa0, 0x24, 0x16, 0xa9, 0x3e, 0x31,
	0x6d, 0x3d, 0xc7, 0x4a, 0xf8, 0x18, 0x6e, 0x69, 0xe6, 0x78, 0xac, 0x6b, 0x76, 0xe1, 0xbe, 0x8c,
	0x5b, 0x70, 0x3c, 0x31, 0x89, 0x61, 0x9b, 0xd4, 0x20, 0x2f, 0x9c, 0x4b, 0x23, 0xa7, 0x2a, 0xd7,
	0xa8, 0x33, 0x35, 0xa7, 0xaa, 0xf8, 0x3e, 0xb4, 0x0a, 0x94, 0x36, 0x36, 0x74, 0x62, 0xe7, 0x74,
	0x4d, 0x78, 0xa9, 0x3a, 0xd1, 0x46, 0x93, 0x21, 0x3d, 0xcf, 0xe1, 0x7d, 0x91, 0xd5, 0xd2, 0x86,
	0x84, 0xe8, 0xbb, 0x00, 0x07, 0xf8, 0x0e, 0xe0, 0x0c, 0x2c, 0x58, 0x80, 0xb0, 0x18, 0x9e, 0x4d,
	0x0c, 0xcb, 0x32, 0x4c, 0xe2, 0x68, 0x26, 0xb1, 0xa9, 0xb9, 0xab, 0x58, 0x57, 0x9f, 0x7e, 0xdb,
	0xb4, 0xd1, 0xf7, 0x4d, 0x1b, 0xfd, 0xdc, 0xb4, 0xd1, 0xe7, 0x5f, 0xed, 0x3d, 0x68, 0x05, 0x61,
	0x2f, 0xe6, 0xae, 0x37, 0x8b, 0xc2, 0x0f, 0xdb, 0x77, 0x96, 0x7d, 0x83, 0xd7, 0xd9, 0x73, 0xbc,
	0xaa, 0x26, 0xf8, 0x93, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x0a, 0x66, 0x33, 0x15, 0xc0, 0x03,
	0x00, 0x00,
}

func (m *ServiceIdentity) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ServiceIdentity) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ServiceIdentity) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.InitBundleId) > 0 {
		i -= len(m.InitBundleId)
		copy(dAtA[i:], m.InitBundleId)
		i = encodeVarintServiceIdentity(dAtA, i, uint64(len(m.InitBundleId)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.SerialStr) > 0 {
		i -= len(m.SerialStr)
		copy(dAtA[i:], m.SerialStr)
		i = encodeVarintServiceIdentity(dAtA, i, uint64(len(m.SerialStr)))
		i--
		dAtA[i] = 0x22
	}
	if m.Type != 0 {
		i = encodeVarintServiceIdentity(dAtA, i, uint64(m.Type))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintServiceIdentity(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0x12
	}
	if m.Srl != nil {
		{
			size := m.Srl.Size()
			i -= size
			if _, err := m.Srl.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
		}
	}
	return len(dAtA) - i, nil
}

func (m *ServiceIdentity_Serial) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ServiceIdentity_Serial) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	i = encodeVarintServiceIdentity(dAtA, i, uint64(m.Serial))
	i--
	dAtA[i] = 0x8
	return len(dAtA) - i, nil
}
func (m *ServiceCertificate) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ServiceCertificate) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ServiceCertificate) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.KeyPem) > 0 {
		i -= len(m.KeyPem)
		copy(dAtA[i:], m.KeyPem)
		i = encodeVarintServiceIdentity(dAtA, i, uint64(len(m.KeyPem)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.CertPem) > 0 {
		i -= len(m.CertPem)
		copy(dAtA[i:], m.CertPem)
		i = encodeVarintServiceIdentity(dAtA, i, uint64(len(m.CertPem)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TypedServiceCertificate) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TypedServiceCertificate) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TypedServiceCertificate) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Cert != nil {
		{
			size, err := m.Cert.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintServiceIdentity(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.ServiceType != 0 {
		i = encodeVarintServiceIdentity(dAtA, i, uint64(m.ServiceType))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *TypedServiceCertificateSet) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TypedServiceCertificateSet) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TypedServiceCertificateSet) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.ServiceCerts) > 0 {
		for iNdEx := len(m.ServiceCerts) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ServiceCerts[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintServiceIdentity(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.CaPem) > 0 {
		i -= len(m.CaPem)
		copy(dAtA[i:], m.CaPem)
		i = encodeVarintServiceIdentity(dAtA, i, uint64(len(m.CaPem)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintServiceIdentity(dAtA []byte, offset int, v uint64) int {
	offset -= sovServiceIdentity(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ServiceIdentity) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Srl != nil {
		n += m.Srl.Size()
	}
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovServiceIdentity(uint64(l))
	}
	if m.Type != 0 {
		n += 1 + sovServiceIdentity(uint64(m.Type))
	}
	l = len(m.SerialStr)
	if l > 0 {
		n += 1 + l + sovServiceIdentity(uint64(l))
	}
	l = len(m.InitBundleId)
	if l > 0 {
		n += 1 + l + sovServiceIdentity(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *ServiceIdentity_Serial) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	n += 1 + sovServiceIdentity(uint64(m.Serial))
	return n
}
func (m *ServiceCertificate) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.CertPem)
	if l > 0 {
		n += 1 + l + sovServiceIdentity(uint64(l))
	}
	l = len(m.KeyPem)
	if l > 0 {
		n += 1 + l + sovServiceIdentity(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TypedServiceCertificate) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ServiceType != 0 {
		n += 1 + sovServiceIdentity(uint64(m.ServiceType))
	}
	if m.Cert != nil {
		l = m.Cert.Size()
		n += 1 + l + sovServiceIdentity(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TypedServiceCertificateSet) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.CaPem)
	if l > 0 {
		n += 1 + l + sovServiceIdentity(uint64(l))
	}
	if len(m.ServiceCerts) > 0 {
		for _, e := range m.ServiceCerts {
			l = e.Size()
			n += 1 + l + sovServiceIdentity(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovServiceIdentity(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozServiceIdentity(x uint64) (n int) {
	return sovServiceIdentity(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ServiceIdentity) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowServiceIdentity
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
			return fmt.Errorf("proto: ServiceIdentity: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ServiceIdentity: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Serial", wireType)
			}
			var v int64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServiceIdentity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Srl = &ServiceIdentity_Serial{v}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServiceIdentity
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
				return ErrInvalidLengthServiceIdentity
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthServiceIdentity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServiceIdentity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= ServiceType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SerialStr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServiceIdentity
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
				return ErrInvalidLengthServiceIdentity
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthServiceIdentity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SerialStr = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InitBundleId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServiceIdentity
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
				return ErrInvalidLengthServiceIdentity
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthServiceIdentity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.InitBundleId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipServiceIdentity(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthServiceIdentity
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ServiceCertificate) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowServiceIdentity
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
			return fmt.Errorf("proto: ServiceCertificate: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ServiceCertificate: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CertPem", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServiceIdentity
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
				return ErrInvalidLengthServiceIdentity
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthServiceIdentity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CertPem = append(m.CertPem[:0], dAtA[iNdEx:postIndex]...)
			if m.CertPem == nil {
				m.CertPem = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field KeyPem", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServiceIdentity
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
				return ErrInvalidLengthServiceIdentity
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthServiceIdentity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.KeyPem = append(m.KeyPem[:0], dAtA[iNdEx:postIndex]...)
			if m.KeyPem == nil {
				m.KeyPem = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipServiceIdentity(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthServiceIdentity
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TypedServiceCertificate) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowServiceIdentity
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
			return fmt.Errorf("proto: TypedServiceCertificate: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TypedServiceCertificate: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ServiceType", wireType)
			}
			m.ServiceType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServiceIdentity
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ServiceType |= ServiceType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Cert", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServiceIdentity
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
				return ErrInvalidLengthServiceIdentity
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthServiceIdentity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Cert == nil {
				m.Cert = &ServiceCertificate{}
			}
			if err := m.Cert.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipServiceIdentity(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthServiceIdentity
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TypedServiceCertificateSet) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowServiceIdentity
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
			return fmt.Errorf("proto: TypedServiceCertificateSet: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TypedServiceCertificateSet: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CaPem", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServiceIdentity
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
				return ErrInvalidLengthServiceIdentity
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthServiceIdentity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CaPem = append(m.CaPem[:0], dAtA[iNdEx:postIndex]...)
			if m.CaPem == nil {
				m.CaPem = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ServiceCerts", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServiceIdentity
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
				return ErrInvalidLengthServiceIdentity
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthServiceIdentity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ServiceCerts = append(m.ServiceCerts, &TypedServiceCertificate{})
			if err := m.ServiceCerts[len(m.ServiceCerts)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipServiceIdentity(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthServiceIdentity
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipServiceIdentity(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowServiceIdentity
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
					return 0, ErrIntOverflowServiceIdentity
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
					return 0, ErrIntOverflowServiceIdentity
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
				return 0, ErrInvalidLengthServiceIdentity
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupServiceIdentity
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthServiceIdentity
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthServiceIdentity        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowServiceIdentity          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupServiceIdentity = fmt.Errorf("proto: unexpected end of group")
)
