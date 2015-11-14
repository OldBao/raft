package coordinate

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"raft/env"
	"sync"
	"time"

	"raft/proto"
)

type Coordinator struct {
	Id          int
	Conn        net.Listener
	RpcServer   *grpc.Server
	PeerMapLock sync.Mutex
	Peers       map[int]*Peer
}

func (coordinator *Coordinator) AppendEntries(ctx context.Context, req *proto.AppendEntriesRequest) (*proto.AppendEntriesResponse, error) {
	env.Log.Info("[APPEND ENTRIES] %+V %+V", ctx, req)
	resp := &proto.AppendEntriesResponse{}
	return resp, nil
}

func (coordinator *Coordinator) Vote(ctx context.Context, req *proto.VoteRequest) (*proto.VoteResponse, error) {
	env.Log.Info("[APPEND ENTRIES] %+V %+V", ctx, req)
	resp := &proto.VoteResponse{}
	return resp, nil
}

func (coordinator *Coordinator) Hello(ctx context.Context, req *proto.Empty) (*proto.HelloResponse, error) {
	env.Log.Info("[Hello] %+V", ctx)

	resp := &proto.HelloResponse{}
	id := int32(coordinator.Id)
	resp.MyId = &id
	return resp, nil
}

func NewCoordinator(id int, addr string) (coordinator *Coordinator, err error) {
	coordinator = &Coordinator{
		Id: id,
	}
	if coordinator.Conn, err = net.Listen("tcp", addr); err != nil {
		env.Log.Warn("Listen [%s] Error", addr)
		return
	}
	env.Log.Info("[CORD][LISTEN][%s]", addr)

	coordinator.RpcServer = grpc.NewServer()
	proto.RegisterRaftServer(coordinator.RpcServer, coordinator)
	return coordinator, nil
}

func (coordinator *Coordinator) ChattingOnePeer(peer *Peer) {
	for {
		if err := peer.Connect(); err != nil {
			time.Sleep(300 * time.Millisecond)
		} else {
			coordinator.PeerMapLock.Lock()
			if _, ok := coordinator.Peers[peer.Id]; ok {
				env.Log.Error("Conflict Id [%s:%d]", peer.Addr, peer.Id)
			} else {
				coordinator.Peers[peer.Id] = peer
			}
			coordinator.PeerMapLock.Unlock()

			for {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

func (coordinator *Coordinator) ChattingToPeers() {
	coordinator.Peers = make(map[int]*Peer)

	for _, peerAddrConf := range env.Conf.Peer.Addr {
		peer := NewPeer(peerAddrConf)
		go coordinator.ChattingOnePeer(peer)
	}

	for {
		time.Sleep(1 * time.Second)
		coordinator.PeerMapLock.Lock()
		env.Log.Debug("Peers: %+V", coordinator.Peers)
		coordinator.PeerMapLock.Unlock()
	}
}

func (coordinator *Coordinator) Work() error {
	env.Log.Info("[CORD][WORK..]")
	go coordinator.ChattingToPeers()

	coordinator.RpcServer.Serve(coordinator.Conn)
	return nil
}
