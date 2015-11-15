package coordinate

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"raft/env"

	"fmt"
	"raft/proto"
	"time"
)

const (
	CONNECTING = 0
	CONNECTED  = 1
	DISCONNECT = 2
)

type Peer struct {
	Id     int
	Addr   string
	Client proto.RaftClient

	//used when current coordiator is leader
	NextIndex  uint64
	MatchIndex uint64

	dial_opts []grpc.DialOption
	call_opts []grpc.CallOption
}

func NewPeer(addr string, timeout time.Duration) (peer *Peer) {
	peer = &Peer{
		Addr: addr,
		Id:   -1,
	}

	peer.dial_opts = append(peer.dial_opts, grpc.WithInsecure())
	peer.dial_opts = append(peer.dial_opts, grpc.WithBlock())
	peer.dial_opts = append(peer.dial_opts, grpc.WithTimeout(time.Millisecond*timeout))
	return peer
}

func (peer *Peer) Connect() error {
	if conn, err := grpc.Dial(peer.Addr, peer.dial_opts...); err != nil {
		//env.Log.Info("connect peer [%s] error: [%v]", peer.Addr, err)
		return err
	} else {
		peer.Client = proto.NewRaftClient(conn)
		if resp, herr := peer.Client.Hello(context.Background(), &proto.Empty{}, peer.call_opts...); err != nil {
			//env.Log.Warn("send peer hello cmd error", herr)
			return herr
		} else {
			peer.Id = int(*resp.MyId)
		}
		if peer.Id == -1 {
			return fmt.Errorf("Invalid id")
		}
	}

	return nil
}

func (peer *Peer) RequestVote(req *proto.VoteRequest) (bool, error) {
	resp, err := peer.Client.RequestVote(context.Background(), req, peer.call_opts...)
	if err != nil {
		env.Log.Warn("RequestVote Error:%v", err)
		return false, err
	}

	//TODO check term
	return *resp.IsGranted, nil
}

func (peer *Peer) HeartBeat(req *proto.AppendEntriesRequest) (bool, error) {
	if _, err := peer.Client.AppendEntries(context.Background(), req, peer.call_opts...); err != nil {
		env.Log.Warn("Send HeartBeat Fails: %v", err)
		return false, err
	}

	return true, nil

}

func (peer *Peer) AppendEntries(req *proto.AppendEntriesRequest) (*proto.AppendEntriesResponse, error) {
	resp, err := peer.Client.AppendEntries(context.Background(), req, peer.call_opts...)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
