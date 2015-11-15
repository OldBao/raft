package coordinate

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"math/rand"
	"net"
	"raft/env"
	"raft/proto"
	"sync"
	"time"
)

const (
	STATE_INITIAL   = 0
	STATE_FOLLOWER  = 1
	STATE_CANDIDATE = 2
	STATE_LEADER    = 3
)

type Coordinator struct {
	Id          int
	CurrentTerm uint64
	VoteFor     int

	//volatile state on all server
	CommitIndex uint64
	LastApplied uint64

	//volatile state on leader
	Conn        net.Listener
	RpcServer   *grpc.Server
	PeerMapLock sync.Mutex
	Peers       map[string]*Peer

	TickMin      time.Duration
	TickRange    int
	ElectTimeout time.Duration
	PeerTimeout  time.Duration

	HbReset     chan bool
	HbTicker    *time.Ticker
	ElectCancel chan bool
	ElectTicker *time.Ticker

	state     int
	stateLock sync.Mutex
}

func (coordinator *Coordinator) AppendEntries(ctx context.Context, req *proto.AppendEntriesRequest) (*proto.AppendEntriesResponse, error) {
	suc := true

	if *req.Term < coordinator.CurrentTerm {
		suc = false
	} else {
		coordinator.CurrentTerm = *req.Term
	}

	if suc == true {
		coordinator.transState(STATE_FOLLOWER)
	}

	resp := &proto.AppendEntriesResponse{
		Term:      &coordinator.CurrentTerm,
		IsSuccess: &suc,
	}

	return resp, nil
}

func (coordinator *Coordinator) RequestVote(ctx context.Context, req *proto.VoteRequest) (*proto.VoteResponse, error) {
	ct := coordinator.CurrentTerm
	granted := false

	resp := &proto.VoteResponse{
		Term:      &ct,
		IsGranted: &granted,
	}

	if *req.Term >= coordinator.CurrentTerm {
		if coordinator.VoteFor == -1 || coordinator.VoteFor == int(*req.CandidateId) {
			granted = true
			coordinator.VoteFor = int(*req.Term)
		}
	}

	return resp, nil
}

func (coordinator *Coordinator) Hello(ctx context.Context, req *proto.Empty) (*proto.HelloResponse, error) {
	env.Log.Info("[Hello]")

	resp := new(proto.HelloResponse)
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

	coordinator.VoteFor = -1
	coordinator.CurrentTerm = 0

	coordinator.TickRange = 150
	coordinator.TickMin = 150
	coordinator.ElectTimeout = 300
	coordinator.PeerTimeout = 150

	coordinator.state = STATE_INITIAL
	coordinator.HbReset = make(chan bool)
	return coordinator, nil
}

func (coordinator *Coordinator) ChattingOnePeer(peer *Peer) {
	for {
		if err := peer.Connect(); err != nil {
			time.Sleep(300 * time.Millisecond)
		} else {
			coordinator.PeerMapLock.Lock()
			if _, ok := coordinator.Peers[peer.Addr]; ok {
				//env.Log.Error("Conflict Id [%s:%d]", peer.Addr, peer.Id)
			} else {
				coordinator.Peers[peer.Addr] = peer
				coordinator.PeerMapLock.Unlock()
				break
			}
			coordinator.PeerMapLock.Unlock()
		}
	}
}

func (coordinator *Coordinator) ChattingToPeers() {
	coordinator.Peers = make(map[string]*Peer)

	for _, peerAddrConf := range env.Conf.Peer.Addr {
		peer := NewPeer(peerAddrConf, coordinator.PeerTimeout)
		go coordinator.ChattingOnePeer(peer)
	}
}

func (coordinator *Coordinator) AsCandidate() {
	coordinator.CurrentTerm += 1

	id := int32(coordinator.Id)
	ct := coordinator.CurrentTerm
	//lli := coordinator.LastLogIndex
	//llt := coordinator.LastLogTerm
	//TODO maintain log
	lli := uint64(1)
	llt := uint64(1)

	req := &proto.VoteRequest{
		Term:         &ct,
		CandidateId:  &id,
		LastLogIndex: &lli,
		LastLogTerm:  &llt,
	}

	env.Log.Info("AsCandidate. %d", coordinator.CurrentTerm)
	coordinator.PeerMapLock.Lock()
	total := len(coordinator.Peers)
	if total == 0 {
		env.Log.Info("No peer found, set my self leader")
		coordinator.PeerMapLock.Unlock()

		coordinator.transState(STATE_LEADER)
		return
	}

	granted := 1 //myself is granted
	rejected := 0
	c := make(chan bool)
	for _, peer := range coordinator.Peers {
		go func(r *proto.VoteRequest, ch chan bool) {
			ok, err := peer.RequestVote(r)
			if err != nil {
				coordinator.resetPeer(peer)
			}
			ch <- ok
		}(req, c)
	}
	coordinator.PeerMapLock.Unlock()

	coordinator.ElectTicker = time.NewTicker(coordinator.ElectTimeout * time.Millisecond)
	for {
		select {
		case <-coordinator.ElectTicker.C:
			//elect timeout
			//TODO abort older vote request
			env.Log.Info("proposer for %d timeout", coordinator.CurrentTerm)
			coordinator.transState(STATE_FOLLOWER)
			return
		case <-coordinator.ElectCancel:
			env.Log.Info("election cancelled")
			coordinator.transState(STATE_FOLLOWER)
			return
		case voteResult := <-c:
			if voteResult == true {
				granted += 1
				if granted > total/2 {
					//election success
					coordinator.transState(STATE_LEADER)
					return
				}
			} else {
				rejected += 1
				if rejected >= total/2 {
					//rejected too much
					coordinator.transState(STATE_FOLLOWER)
					return
				}
			}
		}
	}
}

func (coordinator *Coordinator) AsLeader() {
	env.Log.Info("AsLeader %d", coordinator.CurrentTerm)
	hbchan := make(chan bool)
	go coordinator.KeepLeaderShip(hbchan)
}

func (coordinator *Coordinator) AsFollower() {
	env.Log.Info("AsFollower")
	for {
		tickTimeout := coordinator.TickMin*2 + time.Duration(rand.Intn(coordinator.TickRange))
		coordinator.HbTicker =
			time.NewTicker(tickTimeout * time.Millisecond)

		select {
		case <-coordinator.HbTicker.C:
			//if timeout happens, try propose new election
			env.Log.Info("hb timeout")
			coordinator.transState(STATE_CANDIDATE)
			return
		case <-coordinator.HbReset:
			coordinator.HbTicker.Stop()
		}
	}
}

func (coordinator *Coordinator) Work() error {
	env.Log.Info("[CORD][WORK..]")
	go coordinator.ChattingToPeers()
	coordinator.transState(STATE_FOLLOWER)

	coordinator.RpcServer.Serve(coordinator.Conn)
	return nil
}

func (coordinator *Coordinator) KeepLeaderShip(stop chan bool) {
	for {
		req := coordinator.getAppendEntriesNullRequest()

		coordinator.PeerMapLock.Lock()
		for _, peer := range coordinator.Peers {
			go func(peer *Peer) {
				//env.Log.Info("Sending heartbeat to %s", peer.Id)
				_, err := peer.HeartBeat(req)
				if err != nil {
					coordinator.resetPeer(peer)
				}
			}(peer)
		}
		coordinator.PeerMapLock.Unlock()

		time.Sleep(coordinator.TickMin * time.Millisecond)
	}
}

func (coordinator *Coordinator) getAppendEntriesNullRequest() *proto.AppendEntriesRequest {
	t := coordinator.CurrentTerm
	li := int32(coordinator.Id)
	pli := uint64(1)
	plt := uint64(1)
	lc := uint64(1)
	req := &proto.AppendEntriesRequest{
		Term:         &t,
		LeaderId:     &li,
		PrevLogIndex: &pli,
		PrevLogTerm:  &plt,
		LeaderCommit: &lc,
	}
	return req
}

func (coordinator *Coordinator) transState(newState int) {
	coordinator.stateLock.Lock()
	defer coordinator.stateLock.Unlock()

	switch coordinator.state {
	case STATE_INITIAL:
		switch newState {
		case STATE_FOLLOWER:
			go coordinator.AsFollower()
		default:
			panic(fmt.Errorf("illegal state trans from %d=>%d", coordinator.state, newState))
		}
	case STATE_FOLLOWER:
		switch newState {
		case STATE_CANDIDATE:
			go coordinator.AsCandidate()
		case STATE_FOLLOWER:
			coordinator.HbReset <- true
		default:
			panic(fmt.Errorf("illegal state trans from %d=>%d", coordinator.state, newState))

		}

	case STATE_CANDIDATE:
		switch newState {
		case STATE_LEADER:
			go coordinator.AsLeader()
		case STATE_FOLLOWER:
			//cancel current selectiong if exists
			select {
			case coordinator.ElectCancel <- true:
			default:
			}

			go coordinator.AsFollower()
		default:
			panic(fmt.Errorf("illegal state trans from %d=>%d", coordinator.state, newState))
		}

	case STATE_LEADER:
		switch newState {
		case STATE_FOLLOWER:
			go coordinator.AsFollower()
		default:
			panic(fmt.Errorf("illegal state trans from %d=>%d", coordinator.state, newState))

		}
	default:
		panic(fmt.Errorf("illegal state trans from %d=>%d", coordinator.state, newState))
	}
	if coordinator.state != newState {
		env.Log.Info("Trans %d=>%d", coordinator.state, newState)
	}
	coordinator.state = newState
}

func (coordinator *Coordinator) resetPeer(peer *Peer) {
	coordinator.PeerMapLock.Lock()
	defer coordinator.PeerMapLock.Unlock()
	if _, ok := coordinator.Peers[peer.Addr]; ok {
		delete(coordinator.Peers, peer.Addr)
		go coordinator.ChattingOnePeer(peer)
	}

}
