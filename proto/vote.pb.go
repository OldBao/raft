// Code generated by protoc-gen-go.
// source: vote.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	vote.proto

It has these top-level messages:
	Empty
	LogEntry
	AppendEntriesRequest
	AppendEntriesResponse
	VoteRequest
	VoteResponse
	HelloResponse
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type LogAction int32

const (
	LogAction_NOP LogAction = 0
	LogAction_SET LogAction = 1
	LogAction_INC LogAction = 2
	LogAction_DEC LogAction = 3
	LogAction_DEL LogAction = 4
)

var LogAction_name = map[int32]string{
	0: "NOP",
	1: "SET",
	2: "INC",
	3: "DEC",
	4: "DEL",
}
var LogAction_value = map[string]int32{
	"NOP": 0,
	"SET": 1,
	"INC": 2,
	"DEC": 3,
	"DEL": 4,
}

func (x LogAction) Enum() *LogAction {
	p := new(LogAction)
	*p = x
	return p
}
func (x LogAction) String() string {
	return proto1.EnumName(LogAction_name, int32(x))
}
func (x *LogAction) UnmarshalJSON(data []byte) error {
	value, err := proto1.UnmarshalJSONEnum(LogAction_value, data, "LogAction")
	if err != nil {
		return err
	}
	*x = LogAction(value)
	return nil
}

type Empty struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto1.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}

type LogEntry struct {
	Action           *LogAction `protobuf:"varint,1,req,name=Action,enum=proto.LogAction" json:"Action,omitempty"`
	Value            *int32     `protobuf:"varint,2,opt,name=Value" json:"Value,omitempty"`
	XXX_unrecognized []byte     `json:"-"`
}

func (m *LogEntry) Reset()         { *m = LogEntry{} }
func (m *LogEntry) String() string { return proto1.CompactTextString(m) }
func (*LogEntry) ProtoMessage()    {}

func (m *LogEntry) GetAction() LogAction {
	if m != nil && m.Action != nil {
		return *m.Action
	}
	return LogAction_NOP
}

func (m *LogEntry) GetValue() int32 {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return 0
}

type AppendEntriesRequest struct {
	Term             *uint64     `protobuf:"varint,1,req,name=Term" json:"Term,omitempty"`
	LeaderId         *int32      `protobuf:"varint,2,req,name=LeaderId" json:"LeaderId,omitempty"`
	PrevLogIndex     *uint64     `protobuf:"varint,3,req,name=PrevLogIndex" json:"PrevLogIndex,omitempty"`
	PrevLogTerm      *uint64     `protobuf:"varint,4,req,name=PrevLogTerm" json:"PrevLogTerm,omitempty"`
	Entries          []*LogEntry `protobuf:"bytes,5,rep,name=Entries" json:"Entries,omitempty"`
	LeaderCommit     *uint64     `protobuf:"varint,6,req,name=LeaderCommit" json:"LeaderCommit,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *AppendEntriesRequest) Reset()         { *m = AppendEntriesRequest{} }
func (m *AppendEntriesRequest) String() string { return proto1.CompactTextString(m) }
func (*AppendEntriesRequest) ProtoMessage()    {}

func (m *AppendEntriesRequest) GetTerm() uint64 {
	if m != nil && m.Term != nil {
		return *m.Term
	}
	return 0
}

func (m *AppendEntriesRequest) GetLeaderId() int32 {
	if m != nil && m.LeaderId != nil {
		return *m.LeaderId
	}
	return 0
}

func (m *AppendEntriesRequest) GetPrevLogIndex() uint64 {
	if m != nil && m.PrevLogIndex != nil {
		return *m.PrevLogIndex
	}
	return 0
}

func (m *AppendEntriesRequest) GetPrevLogTerm() uint64 {
	if m != nil && m.PrevLogTerm != nil {
		return *m.PrevLogTerm
	}
	return 0
}

func (m *AppendEntriesRequest) GetEntries() []*LogEntry {
	if m != nil {
		return m.Entries
	}
	return nil
}

func (m *AppendEntriesRequest) GetLeaderCommit() uint64 {
	if m != nil && m.LeaderCommit != nil {
		return *m.LeaderCommit
	}
	return 0
}

type AppendEntriesResponse struct {
	Term             *uint64 `protobuf:"varint,1,req,name=Term" json:"Term,omitempty"`
	IsSuccess        *bool   `protobuf:"varint,2,req,name=IsSuccess" json:"IsSuccess,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *AppendEntriesResponse) Reset()         { *m = AppendEntriesResponse{} }
func (m *AppendEntriesResponse) String() string { return proto1.CompactTextString(m) }
func (*AppendEntriesResponse) ProtoMessage()    {}

func (m *AppendEntriesResponse) GetTerm() uint64 {
	if m != nil && m.Term != nil {
		return *m.Term
	}
	return 0
}

func (m *AppendEntriesResponse) GetIsSuccess() bool {
	if m != nil && m.IsSuccess != nil {
		return *m.IsSuccess
	}
	return false
}

type VoteRequest struct {
	Term             *uint64 `protobuf:"varint,1,req,name=Term" json:"Term,omitempty"`
	CandidateId      *int32  `protobuf:"varint,2,req,name=CandidateId" json:"CandidateId,omitempty"`
	LastLogIndex     *uint64 `protobuf:"varint,3,req,name=LastLogIndex" json:"LastLogIndex,omitempty"`
	LastLogTerm      *uint64 `protobuf:"varint,4,req,name=LastLogTerm" json:"LastLogTerm,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *VoteRequest) Reset()         { *m = VoteRequest{} }
func (m *VoteRequest) String() string { return proto1.CompactTextString(m) }
func (*VoteRequest) ProtoMessage()    {}

func (m *VoteRequest) GetTerm() uint64 {
	if m != nil && m.Term != nil {
		return *m.Term
	}
	return 0
}

func (m *VoteRequest) GetCandidateId() int32 {
	if m != nil && m.CandidateId != nil {
		return *m.CandidateId
	}
	return 0
}

func (m *VoteRequest) GetLastLogIndex() uint64 {
	if m != nil && m.LastLogIndex != nil {
		return *m.LastLogIndex
	}
	return 0
}

func (m *VoteRequest) GetLastLogTerm() uint64 {
	if m != nil && m.LastLogTerm != nil {
		return *m.LastLogTerm
	}
	return 0
}

type VoteResponse struct {
	Term             *uint64 `protobuf:"varint,1,req,name=Term" json:"Term,omitempty"`
	IsGranted        *bool   `protobuf:"varint,2,req,name=IsGranted" json:"IsGranted,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *VoteResponse) Reset()         { *m = VoteResponse{} }
func (m *VoteResponse) String() string { return proto1.CompactTextString(m) }
func (*VoteResponse) ProtoMessage()    {}

func (m *VoteResponse) GetTerm() uint64 {
	if m != nil && m.Term != nil {
		return *m.Term
	}
	return 0
}

func (m *VoteResponse) GetIsGranted() bool {
	if m != nil && m.IsGranted != nil {
		return *m.IsGranted
	}
	return false
}

type HelloResponse struct {
	MyId             *int32 `protobuf:"varint,1,req,name=MyId,def=-1" json:"MyId,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *HelloResponse) Reset()         { *m = HelloResponse{} }
func (m *HelloResponse) String() string { return proto1.CompactTextString(m) }
func (*HelloResponse) ProtoMessage()    {}

const Default_HelloResponse_MyId int32 = -1

func (m *HelloResponse) GetMyId() int32 {
	if m != nil && m.MyId != nil {
		return *m.MyId
	}
	return Default_HelloResponse_MyId
}

func init() {
	proto1.RegisterType((*Empty)(nil), "proto.Empty")
	proto1.RegisterType((*LogEntry)(nil), "proto.LogEntry")
	proto1.RegisterType((*AppendEntriesRequest)(nil), "proto.AppendEntriesRequest")
	proto1.RegisterType((*AppendEntriesResponse)(nil), "proto.AppendEntriesResponse")
	proto1.RegisterType((*VoteRequest)(nil), "proto.VoteRequest")
	proto1.RegisterType((*VoteResponse)(nil), "proto.VoteResponse")
	proto1.RegisterType((*HelloResponse)(nil), "proto.HelloResponse")
	proto1.RegisterEnum("proto.LogAction", LogAction_name, LogAction_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Client API for Raft service

type RaftClient interface {
	AppendEntries(ctx context.Context, in *AppendEntriesRequest, opts ...grpc.CallOption) (*AppendEntriesResponse, error)
	RequestVote(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteResponse, error)
	Hello(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*HelloResponse, error)
}

type raftClient struct {
	cc *grpc.ClientConn
}

func NewRaftClient(cc *grpc.ClientConn) RaftClient {
	return &raftClient{cc}
}

func (c *raftClient) AppendEntries(ctx context.Context, in *AppendEntriesRequest, opts ...grpc.CallOption) (*AppendEntriesResponse, error) {
	out := new(AppendEntriesResponse)
	err := grpc.Invoke(ctx, "/proto.raft/AppendEntries", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftClient) RequestVote(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteResponse, error) {
	out := new(VoteResponse)
	err := grpc.Invoke(ctx, "/proto.raft/RequestVote", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftClient) Hello(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*HelloResponse, error) {
	out := new(HelloResponse)
	err := grpc.Invoke(ctx, "/proto.raft/Hello", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Raft service

type RaftServer interface {
	AppendEntries(context.Context, *AppendEntriesRequest) (*AppendEntriesResponse, error)
	RequestVote(context.Context, *VoteRequest) (*VoteResponse, error)
	Hello(context.Context, *Empty) (*HelloResponse, error)
}

func RegisterRaftServer(s *grpc.Server, srv RaftServer) {
	s.RegisterService(&_Raft_serviceDesc, srv)
}

func _Raft_AppendEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(AppendEntriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(RaftServer).AppendEntries(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Raft_RequestVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(VoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(RaftServer).RequestVote(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Raft_Hello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(RaftServer).Hello(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _Raft_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.raft",
	HandlerType: (*RaftServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AppendEntries",
			Handler:    _Raft_AppendEntries_Handler,
		},
		{
			MethodName: "RequestVote",
			Handler:    _Raft_RequestVote_Handler,
		},
		{
			MethodName: "Hello",
			Handler:    _Raft_Hello_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}
