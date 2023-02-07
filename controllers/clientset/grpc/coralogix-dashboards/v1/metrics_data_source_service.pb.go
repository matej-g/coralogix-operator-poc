// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.8
// source: com/coralogixapis/dashboards/v1/services/metrics_data_source_service.proto

package v1

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SearchMetricsTimeSeriesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TimeFrame   *TimeFrame              `protobuf:"bytes,1,opt,name=time_frame,json=timeFrame,proto3" json:"time_frame,omitempty"`
	Interval    *durationpb.Duration    `protobuf:"bytes,2,opt,name=interval,proto3" json:"interval,omitempty"`
	PromqlQuery *wrapperspb.StringValue `protobuf:"bytes,3,opt,name=promql_query,json=promqlQuery,proto3" json:"promql_query,omitempty"`
	Limit       *wrapperspb.Int32Value  `protobuf:"bytes,4,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *SearchMetricsTimeSeriesRequest) Reset() {
	*x = SearchMetricsTimeSeriesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchMetricsTimeSeriesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchMetricsTimeSeriesRequest) ProtoMessage() {}

func (x *SearchMetricsTimeSeriesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchMetricsTimeSeriesRequest.ProtoReflect.Descriptor instead.
func (*SearchMetricsTimeSeriesRequest) Descriptor() ([]byte, []int) {
	return file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_rawDescGZIP(), []int{0}
}

func (x *SearchMetricsTimeSeriesRequest) GetTimeFrame() *TimeFrame {
	if x != nil {
		return x.TimeFrame
	}
	return nil
}

func (x *SearchMetricsTimeSeriesRequest) GetInterval() *durationpb.Duration {
	if x != nil {
		return x.Interval
	}
	return nil
}

func (x *SearchMetricsTimeSeriesRequest) GetPromqlQuery() *wrapperspb.StringValue {
	if x != nil {
		return x.PromqlQuery
	}
	return nil
}

func (x *SearchMetricsTimeSeriesRequest) GetLimit() *wrapperspb.Int32Value {
	if x != nil {
		return x.Limit
	}
	return nil
}

type SearchMetricsTimeSeriesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TimeSeries      []*TimeSeries `protobuf:"bytes,1,rep,name=time_series,json=timeSeries,proto3" json:"time_series,omitempty"`
	IsLimitExceeded bool          `protobuf:"varint,2,opt,name=is_limit_exceeded,json=isLimitExceeded,proto3" json:"is_limit_exceeded,omitempty"`
}

func (x *SearchMetricsTimeSeriesResponse) Reset() {
	*x = SearchMetricsTimeSeriesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchMetricsTimeSeriesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchMetricsTimeSeriesResponse) ProtoMessage() {}

func (x *SearchMetricsTimeSeriesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchMetricsTimeSeriesResponse.ProtoReflect.Descriptor instead.
func (*SearchMetricsTimeSeriesResponse) Descriptor() ([]byte, []int) {
	return file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_rawDescGZIP(), []int{1}
}

func (x *SearchMetricsTimeSeriesResponse) GetTimeSeries() []*TimeSeries {
	if x != nil {
		return x.TimeSeries
	}
	return nil
}

func (x *SearchMetricsTimeSeriesResponse) GetIsLimitExceeded() bool {
	if x != nil {
		return x.IsLimitExceeded
	}
	return false
}

var File_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto protoreflect.FileDescriptor

var file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_rawDesc = []byte{
	0x0a, 0x4a, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x72, 0x61, 0x6c, 0x6f, 0x67, 0x69, 0x78, 0x61,
	0x70, 0x69, 0x73, 0x2f, 0x64, 0x61, 0x73, 0x68, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x73, 0x2f, 0x76,
	0x31, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x28, 0x63, 0x6f,
	0x6d, 0x2e, 0x63, 0x6f, 0x72, 0x61, 0x6c, 0x6f, 0x67, 0x69, 0x78, 0x61, 0x70, 0x69, 0x73, 0x2e,
	0x64, 0x61, 0x73, 0x68, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x1a, 0x2f, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x72, 0x61,
	0x6c, 0x6f, 0x67, 0x69, 0x78, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x64, 0x61, 0x73, 0x68, 0x62, 0x6f,
	0x61, 0x72, 0x64, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x75, 0x64, 0x69, 0x74, 0x5f, 0x6c, 0x6f,
	0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x37, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x72,
	0x61, 0x6c, 0x6f, 0x67, 0x69, 0x78, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x64, 0x61, 0x73, 0x68, 0x62,
	0x6f, 0x61, 0x72, 0x64, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x5f, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x38, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x72, 0x61, 0x6c, 0x6f, 0x67, 0x69, 0x78, 0x61,
	0x70, 0x69, 0x73, 0x2f, 0x64, 0x61, 0x73, 0x68, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x73, 0x2f, 0x76,
	0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x73, 0x65,
	0x72, 0x69, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70,
	0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9d, 0x02, 0x0a, 0x1e, 0x53,
	0x65, 0x61, 0x72, 0x63, 0x68, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x54, 0x69, 0x6d, 0x65,
	0x53, 0x65, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x50, 0x0a,
	0x0a, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x31, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x6f, 0x72, 0x61, 0x6c, 0x6f, 0x67, 0x69,
	0x78, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x64, 0x61, 0x73, 0x68, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x46,
	0x72, 0x61, 0x6d, 0x65, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x12,
	0x35, 0x0a, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x12, 0x3f, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x6d, 0x71, 0x6c,
	0x5f, 0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53,
	0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x6d,
	0x71, 0x6c, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x31, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x22, 0xa2, 0x01, 0x0a, 0x1f, 0x53,
	0x65, 0x61, 0x72, 0x63, 0x68, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x54, 0x69, 0x6d, 0x65,
	0x53, 0x65, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x53,
	0x0a, 0x0b, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x32, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x6f, 0x72, 0x61, 0x6c, 0x6f,
	0x67, 0x69, 0x78, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x64, 0x61, 0x73, 0x68, 0x62, 0x6f, 0x61, 0x72,
	0x64, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x53, 0x65, 0x72, 0x69, 0x65, 0x73, 0x52, 0x0a, 0x74, 0x69, 0x6d, 0x65, 0x53, 0x65, 0x72,
	0x69, 0x65, 0x73, 0x12, 0x2a, 0x0a, 0x11, 0x69, 0x73, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x5f,
	0x65, 0x78, 0x63, 0x65, 0x65, 0x64, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0f,
	0x69, 0x73, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x45, 0x78, 0x63, 0x65, 0x65, 0x64, 0x65, 0x64, 0x32,
	0xfd, 0x01, 0x0a, 0x18, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x44, 0x61, 0x74, 0x61, 0x53,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0xe0, 0x01, 0x0a,
	0x17, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x54, 0x69,
	0x6d, 0x65, 0x53, 0x65, 0x72, 0x69, 0x65, 0x73, 0x12, 0x48, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x63,
	0x6f, 0x72, 0x61, 0x6c, 0x6f, 0x67, 0x69, 0x78, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x64, 0x61, 0x73,
	0x68, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x73, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x73, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x65, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x49, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x6f, 0x72, 0x61, 0x6c, 0x6f, 0x67,
	0x69, 0x78, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x64, 0x61, 0x73, 0x68, 0x62, 0x6f, 0x61, 0x72, 0x64,
	0x73, 0x2e, 0x76, 0x31, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x53, 0x65,
	0x61, 0x72, 0x63, 0x68, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x54, 0x69, 0x6d, 0x65, 0x53,
	0x65, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x30, 0xba,
	0xb8, 0x02, 0x2c, 0x0a, 0x2a, 0x66, 0x65, 0x74, 0x63, 0x68, 0x20, 0x74, 0x69, 0x6d, 0x65, 0x20,
	0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x20, 0x66, 0x72, 0x6f, 0x6d, 0x20, 0x6d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x73, 0x20, 0x64, 0x61, 0x74, 0x61, 0x20, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x42,
	0x04, 0x5a, 0x02, 0x2e, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_rawDescOnce sync.Once
	file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_rawDescData = file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_rawDesc
)

func file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_rawDescGZIP() []byte {
	file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_rawDescOnce.Do(func() {
		file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_rawDescData)
	})
	return file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_rawDescData
}

var file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_goTypes = []interface{}{
	(*SearchMetricsTimeSeriesRequest)(nil),  // 0: com.coralogixapis.dashboards.v1.services.SearchMetricsTimeSeriesRequest
	(*SearchMetricsTimeSeriesResponse)(nil), // 1: com.coralogixapis.dashboards.v1.services.SearchMetricsTimeSeriesResponse
	(*TimeFrame)(nil),                       // 2: com.coralogixapis.dashboards.v1.common.TimeFrame
	(*durationpb.Duration)(nil),             // 3: google.protobuf.Duration
	(*wrapperspb.StringValue)(nil),          // 4: google.protobuf.StringValue
	(*wrapperspb.Int32Value)(nil),           // 5: google.protobuf.Int32Value
	(*TimeSeries)(nil),                      // 6: com.coralogixapis.dashboards.v1.common.TimeSeries
}
var file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_depIdxs = []int32{
	2, // 0: com.coralogixapis.dashboards.v1.services.SearchMetricsTimeSeriesRequest.time_frame:type_name -> com.coralogixapis.dashboards.v1.common.TimeFrame
	3, // 1: com.coralogixapis.dashboards.v1.services.SearchMetricsTimeSeriesRequest.interval:type_name -> google.protobuf.Duration
	4, // 2: com.coralogixapis.dashboards.v1.services.SearchMetricsTimeSeriesRequest.promql_query:type_name -> google.protobuf.StringValue
	5, // 3: com.coralogixapis.dashboards.v1.services.SearchMetricsTimeSeriesRequest.limit:type_name -> google.protobuf.Int32Value
	6, // 4: com.coralogixapis.dashboards.v1.services.SearchMetricsTimeSeriesResponse.time_series:type_name -> com.coralogixapis.dashboards.v1.common.TimeSeries
	0, // 5: com.coralogixapis.dashboards.v1.services.MetricsDataSourceService.SearchMetricsTimeSeries:input_type -> com.coralogixapis.dashboards.v1.services.SearchMetricsTimeSeriesRequest
	1, // 6: com.coralogixapis.dashboards.v1.services.MetricsDataSourceService.SearchMetricsTimeSeries:output_type -> com.coralogixapis.dashboards.v1.services.SearchMetricsTimeSeriesResponse
	6, // [6:7] is the sub-list for method output_type
	5, // [5:6] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_init() }
func file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_init() {
	if File_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto != nil {
		return
	}
	file_com_coralogixapis_dashboards_v1_audit_log_proto_init()
	file_com_coralogixapis_dashboards_v1_common_time_frame_proto_init()
	file_com_coralogixapis_dashboards_v1_common_time_series_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchMetricsTimeSeriesRequest); i {
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
		file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchMetricsTimeSeriesResponse); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_goTypes,
		DependencyIndexes: file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_depIdxs,
		MessageInfos:      file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_msgTypes,
	}.Build()
	File_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto = out.File
	file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_rawDesc = nil
	file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_goTypes = nil
	file_com_coralogixapis_dashboards_v1_services_metrics_data_source_service_proto_depIdxs = nil
}
