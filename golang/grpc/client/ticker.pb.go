// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: ticker.proto

package client

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// The request message containing the stock info
type StockRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbol string  `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Close  float64 `protobuf:"fixed64,2,opt,name=close,proto3" json:"close,omitempty"`
	High   float64 `protobuf:"fixed64,3,opt,name=high,proto3" json:"high,omitempty"`
	Low    float64 `protobuf:"fixed64,4,opt,name=low,proto3" json:"low,omitempty"`
	Date   string  `protobuf:"bytes,5,opt,name=date,proto3" json:"date,omitempty"`
}

func (x *StockRequest) Reset() {
	*x = StockRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ticker_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StockRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StockRequest) ProtoMessage() {}

func (x *StockRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ticker_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StockRequest.ProtoReflect.Descriptor instead.
func (*StockRequest) Descriptor() ([]byte, []int) {
	return file_ticker_proto_rawDescGZIP(), []int{0}
}

func (x *StockRequest) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *StockRequest) GetClose() float64 {
	if x != nil {
		return x.Close
	}
	return 0
}

func (x *StockRequest) GetHigh() float64 {
	if x != nil {
		return x.High
	}
	return 0
}

func (x *StockRequest) GetLow() float64 {
	if x != nil {
		return x.Low
	}
	return 0
}

func (x *StockRequest) GetDate() string {
	if x != nil {
		return x.Date
	}
	return ""
}

// The request message containing the stock info
type StockResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StockResponse) Reset() {
	*x = StockResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ticker_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StockResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StockResponse) ProtoMessage() {}

func (x *StockResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ticker_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StockResponse.ProtoReflect.Descriptor instead.
func (*StockResponse) Descriptor() ([]byte, []int) {
	return file_ticker_proto_rawDescGZIP(), []int{1}
}

// The request message containing the user's name.
type SummaryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sas string `protobuf:"bytes,1,opt,name=sas,proto3" json:"sas,omitempty"`
}

func (x *SummaryRequest) Reset() {
	*x = SummaryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ticker_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SummaryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SummaryRequest) ProtoMessage() {}

func (x *SummaryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ticker_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SummaryRequest.ProtoReflect.Descriptor instead.
func (*SummaryRequest) Descriptor() ([]byte, []int) {
	return file_ticker_proto_rawDescGZIP(), []int{2}
}

func (x *SummaryRequest) GetSas() string {
	if x != nil {
		return x.Sas
	}
	return ""
}

type SummaryResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Message:
	//	*SummaryResponse_Reply
	//	*SummaryResponse_Error
	Message isSummaryResponse_Message `protobuf_oneof:"message"`
}

func (x *SummaryResponse) Reset() {
	*x = SummaryResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ticker_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SummaryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SummaryResponse) ProtoMessage() {}

func (x *SummaryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ticker_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SummaryResponse.ProtoReflect.Descriptor instead.
func (*SummaryResponse) Descriptor() ([]byte, []int) {
	return file_ticker_proto_rawDescGZIP(), []int{3}
}

func (m *SummaryResponse) GetMessage() isSummaryResponse_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (x *SummaryResponse) GetReply() *SummaryReply {
	if x, ok := x.GetMessage().(*SummaryResponse_Reply); ok {
		return x.Reply
	}
	return nil
}

func (x *SummaryResponse) GetError() string {
	if x, ok := x.GetMessage().(*SummaryResponse_Error); ok {
		return x.Error
	}
	return ""
}

type isSummaryResponse_Message interface {
	isSummaryResponse_Message()
}

type SummaryResponse_Reply struct {
	// RateResponse returns the exchange rate when updated
	Reply *SummaryReply `protobuf:"bytes,1,opt,name=reply,proto3,oneof"`
}

type SummaryResponse_Error struct {
	// status returns rich error messages to the client
	Error string `protobuf:"bytes,2,opt,name=error,proto3,oneof"`
}

func (*SummaryResponse_Reply) isSummaryResponse_Message() {}

func (*SummaryResponse_Error) isSummaryResponse_Message() {}

// The response message containing the greetings
type SummaryReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Close       float64 `protobuf:"fixed64,1,opt,name=close,proto3" json:"close,omitempty"`
	High        float64 `protobuf:"fixed64,2,opt,name=high,proto3" json:"high,omitempty"`
	Low         float64 `protobuf:"fixed64,3,opt,name=low,proto3" json:"low,omitempty"`
	Average     float64 `protobuf:"fixed64,4,opt,name=average,proto3" json:"average,omitempty"`
	MinLP3      float64 `protobuf:"fixed64,5,opt,name=minLP3,proto3" json:"minLP3,omitempty"`
	Ema5        float64 `protobuf:"fixed64,6,opt,name=ema5,proto3" json:"ema5,omitempty"`
	Ema20       float64 `protobuf:"fixed64,7,opt,name=ema20,proto3" json:"ema20,omitempty"`
	Rsi         float64 `protobuf:"fixed64,8,opt,name=rsi,proto3" json:"rsi,omitempty"`
	Hl3         float64 `protobuf:"fixed64,9,opt,name=hl3,proto3" json:"hl3,omitempty"`
	Trend       float64 `protobuf:"fixed64,10,opt,name=trend,proto3" json:"trend,omitempty"`
	Buy         float64 `protobuf:"fixed64,11,opt,name=buy,proto3" json:"buy,omitempty"`
	Support     float64 `protobuf:"fixed64,12,opt,name=support,proto3" json:"support,omitempty"`
	Bullish     float64 `protobuf:"fixed64,13,opt,name=bullish,proto3" json:"bullish,omitempty"`
	Barish      float64 `protobuf:"fixed64,14,opt,name=barish,proto3" json:"barish,omitempty"`
	PreviousBuy float64 `protobuf:"fixed64,15,opt,name=previousBuy,proto3" json:"previousBuy,omitempty"`
}

func (x *SummaryReply) Reset() {
	*x = SummaryReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ticker_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SummaryReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SummaryReply) ProtoMessage() {}

func (x *SummaryReply) ProtoReflect() protoreflect.Message {
	mi := &file_ticker_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SummaryReply.ProtoReflect.Descriptor instead.
func (*SummaryReply) Descriptor() ([]byte, []int) {
	return file_ticker_proto_rawDescGZIP(), []int{4}
}

func (x *SummaryReply) GetClose() float64 {
	if x != nil {
		return x.Close
	}
	return 0
}

func (x *SummaryReply) GetHigh() float64 {
	if x != nil {
		return x.High
	}
	return 0
}

func (x *SummaryReply) GetLow() float64 {
	if x != nil {
		return x.Low
	}
	return 0
}

func (x *SummaryReply) GetAverage() float64 {
	if x != nil {
		return x.Average
	}
	return 0
}

func (x *SummaryReply) GetMinLP3() float64 {
	if x != nil {
		return x.MinLP3
	}
	return 0
}

func (x *SummaryReply) GetEma5() float64 {
	if x != nil {
		return x.Ema5
	}
	return 0
}

func (x *SummaryReply) GetEma20() float64 {
	if x != nil {
		return x.Ema20
	}
	return 0
}

func (x *SummaryReply) GetRsi() float64 {
	if x != nil {
		return x.Rsi
	}
	return 0
}

func (x *SummaryReply) GetHl3() float64 {
	if x != nil {
		return x.Hl3
	}
	return 0
}

func (x *SummaryReply) GetTrend() float64 {
	if x != nil {
		return x.Trend
	}
	return 0
}

func (x *SummaryReply) GetBuy() float64 {
	if x != nil {
		return x.Buy
	}
	return 0
}

func (x *SummaryReply) GetSupport() float64 {
	if x != nil {
		return x.Support
	}
	return 0
}

func (x *SummaryReply) GetBullish() float64 {
	if x != nil {
		return x.Bullish
	}
	return 0
}

func (x *SummaryReply) GetBarish() float64 {
	if x != nil {
		return x.Barish
	}
	return 0
}

func (x *SummaryReply) GetPreviousBuy() float64 {
	if x != nil {
		return x.PreviousBuy
	}
	return 0
}

var File_ticker_proto protoreflect.FileDescriptor

var file_ticker_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x76, 0x0a, 0x0c, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x12, 0x14, 0x0a,
	0x05, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x63, 0x6c,
	0x6f, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x69, 0x67, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x04, 0x68, 0x69, 0x67, 0x68, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x6f, 0x77, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x6c, 0x6f, 0x77, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x65, 0x22, 0x0f, 0x0a,
	0x0d, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x22,
	0x0a, 0x0e, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x10, 0x0a, 0x03, 0x73, 0x61, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x73,
	0x61, 0x73, 0x22, 0x61, 0x0a, 0x0f, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a, 0x05, 0x72, 0x65, 0x70, 0x6c, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x75, 0x6d,
	0x6d, 0x61, 0x72, 0x79, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x48, 0x00, 0x52, 0x05, 0x72, 0x65, 0x70,
	0x6c, 0x79, 0x12, 0x16, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x00, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x42, 0x09, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0xe0, 0x02, 0x0a, 0x0c, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72,
	0x79, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x68, 0x69, 0x67, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x68, 0x69, 0x67, 0x68,
	0x12, 0x10, 0x0a, 0x03, 0x6c, 0x6f, 0x77, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x6c,
	0x6f, 0x77, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x07, 0x61, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x6d, 0x69, 0x6e, 0x4c, 0x50, 0x33, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06, 0x6d, 0x69,
	0x6e, 0x4c, 0x50, 0x33, 0x12, 0x12, 0x0a, 0x04, 0x65, 0x6d, 0x61, 0x35, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x04, 0x65, 0x6d, 0x61, 0x35, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x32,
	0x30, 0x18, 0x07, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x32, 0x30, 0x12, 0x10,
	0x0a, 0x03, 0x72, 0x73, 0x69, 0x18, 0x08, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x72, 0x73, 0x69,
	0x12, 0x10, 0x0a, 0x03, 0x68, 0x6c, 0x33, 0x18, 0x09, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x68,
	0x6c, 0x33, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x72, 0x65, 0x6e, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x05, 0x74, 0x72, 0x65, 0x6e, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x62, 0x75, 0x79, 0x18,
	0x0b, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x62, 0x75, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75,
	0x70, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x01, 0x52, 0x07, 0x73, 0x75, 0x70,
	0x70, 0x6f, 0x72, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x62, 0x75, 0x6c, 0x6c, 0x69, 0x73, 0x68, 0x18,
	0x0d, 0x20, 0x01, 0x28, 0x01, 0x52, 0x07, 0x62, 0x75, 0x6c, 0x6c, 0x69, 0x73, 0x68, 0x12, 0x16,
	0x0a, 0x06, 0x62, 0x61, 0x72, 0x69, 0x73, 0x68, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06,
	0x62, 0x61, 0x72, 0x69, 0x73, 0x68, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f,
	0x75, 0x73, 0x42, 0x75, 0x79, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0b, 0x70, 0x72, 0x65,
	0x76, 0x69, 0x6f, 0x75, 0x73, 0x42, 0x75, 0x79, 0x32, 0x87, 0x01, 0x0a, 0x06, 0x54, 0x69, 0x63,
	0x6b, 0x65, 0x72, 0x12, 0x3a, 0x0a, 0x0b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x74, 0x6f,
	0x63, 0x6b, 0x12, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x6f, 0x63, 0x6b,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x53, 0x74, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x41, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x12, 0x15, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x75, 0x6d,
	0x6d, 0x61, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01,
	0x30, 0x01, 0x42, 0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x76, 0x73, 0x68, 0x65, 0x6f, 0x72, 0x61, 0x6e, 0x2f, 0x74, 0x72, 0x65, 0x6e, 0x64, 0x73,
	0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x3b, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ticker_proto_rawDescOnce sync.Once
	file_ticker_proto_rawDescData = file_ticker_proto_rawDesc
)

func file_ticker_proto_rawDescGZIP() []byte {
	file_ticker_proto_rawDescOnce.Do(func() {
		file_ticker_proto_rawDescData = protoimpl.X.CompressGZIP(file_ticker_proto_rawDescData)
	})
	return file_ticker_proto_rawDescData
}

var file_ticker_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_ticker_proto_goTypes = []interface{}{
	(*StockRequest)(nil),    // 0: proto.StockRequest
	(*StockResponse)(nil),   // 1: proto.StockResponse
	(*SummaryRequest)(nil),  // 2: proto.SummaryRequest
	(*SummaryResponse)(nil), // 3: proto.SummaryResponse
	(*SummaryReply)(nil),    // 4: proto.SummaryReply
}
var file_ticker_proto_depIdxs = []int32{
	4, // 0: proto.SummaryResponse.reply:type_name -> proto.SummaryReply
	0, // 1: proto.Ticker.UpdateStock:input_type -> proto.StockRequest
	2, // 2: proto.Ticker.GetSummary:input_type -> proto.SummaryRequest
	1, // 3: proto.Ticker.UpdateStock:output_type -> proto.StockResponse
	3, // 4: proto.Ticker.GetSummary:output_type -> proto.SummaryResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_ticker_proto_init() }
func file_ticker_proto_init() {
	if File_ticker_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_ticker_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StockRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ticker_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StockResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ticker_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SummaryRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ticker_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SummaryResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ticker_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SummaryReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_ticker_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*SummaryResponse_Reply)(nil),
		(*SummaryResponse_Error)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_ticker_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_ticker_proto_goTypes,
		DependencyIndexes: file_ticker_proto_depIdxs,
		MessageInfos:      file_ticker_proto_msgTypes,
	}.Build()
	File_ticker_proto = out.File
	file_ticker_proto_rawDesc = nil
	file_ticker_proto_goTypes = nil
	file_ticker_proto_depIdxs = nil
}
