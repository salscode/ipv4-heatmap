// Code generated by protoc-gen-go.
// source: locdata.proto
// DO NOT EDIT!

/*
Package locdata is a generated protocol buffer package.

It is generated from these files:
	locdata.proto

It has these top-level messages:
	LocDisplayData
*/
package locdata

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type LocDisplayData struct {
	Lldata           []*LocDisplayData_LatLon `protobuf:"bytes,1,rep,name=lldata" json:"lldata,omitempty"`
	Count            *int32                   `protobuf:"varint,2,req,name=count" json:"count,omitempty"`
	Coordformat      *string                  `protobuf:"bytes,3,req,name=coordformat" json:"coordformat,omitempty"`
	XXX_unrecognized []byte                   `json:"-"`
}

func (m *LocDisplayData) Reset()         { *m = LocDisplayData{} }
func (m *LocDisplayData) String() string { return proto.CompactTextString(m) }
func (*LocDisplayData) ProtoMessage()    {}

func (m *LocDisplayData) GetLldata() []*LocDisplayData_LatLon {
	if m != nil {
		return m.Lldata
	}
	return nil
}

func (m *LocDisplayData) GetCount() int32 {
	if m != nil && m.Count != nil {
		return *m.Count
	}
	return 0
}

func (m *LocDisplayData) GetCoordformat() string {
	if m != nil && m.Coordformat != nil {
		return *m.Coordformat
	}
	return ""
}

type LocDisplayData_LatLon struct {
	Latitude         *float32 `protobuf:"fixed32,1,req,name=latitude" json:"latitude,omitempty"`
	Longitude        *float32 `protobuf:"fixed32,2,req,name=longitude" json:"longitude,omitempty"`
	Intensity        *float32 `protobuf:"fixed32,3,req,name=intensity" json:"intensity,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *LocDisplayData_LatLon) Reset()         { *m = LocDisplayData_LatLon{} }
func (m *LocDisplayData_LatLon) String() string { return proto.CompactTextString(m) }
func (*LocDisplayData_LatLon) ProtoMessage()    {}

func (m *LocDisplayData_LatLon) GetLatitude() float32 {
	if m != nil && m.Latitude != nil {
		return *m.Latitude
	}
	return 0
}

func (m *LocDisplayData_LatLon) GetLongitude() float32 {
	if m != nil && m.Longitude != nil {
		return *m.Longitude
	}
	return 0
}

func (m *LocDisplayData_LatLon) GetIntensity() float32 {
	if m != nil && m.Intensity != nil {
		return *m.Intensity
	}
	return 0
}

func init() {
	proto.RegisterType((*LocDisplayData)(nil), "locdata.LocDisplayData")
	proto.RegisterType((*LocDisplayData_LatLon)(nil), "locdata.LocDisplayData.LatLon")
}
