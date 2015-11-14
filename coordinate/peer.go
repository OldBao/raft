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
	Addr   string
	Client proto.RaftClient
	Id     int

	opts []grpc.DialOption
}

func NewPeer(addr string) (peer *Peer) {
	peer = &Peer{
		Addr: addr,
		Id:   -1,
	}

	peer.opts = append(peer.opts, grpc.WithInsecure())
	peer.opts = append(peer.opts, grpc.WithBlock())
	peer.opts = append(peer.opts, grpc.WithTimeout(3*time.Second))
	return peer
}

func (peer *Peer) Connect() error {
	if conn, err := grpc.Dial(peer.Addr, peer.opts...); err != nil {
		env.Log.Info("connect peer [%s] error: [%v]", peer.Addr, err)
		return err
	} else {
		peer.Client = proto.NewRaftClient(conn)
		if resp, herr := peer.Client.Hello(context.Background(), &proto.Empty{}); err != nil {
			env.Log.Warn("send peer hello cmd error", herr)
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
