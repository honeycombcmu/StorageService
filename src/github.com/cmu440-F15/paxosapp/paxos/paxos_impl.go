package paxos

import (
	"encoding/json"
	"errors"
	//"fmt"
	"github.com/cmu440-F15/paxosapp/rpc/paxosrpc"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
	"time"
)

const (
	TIMEOUT_LIMIT = 15
)

type paxosNode struct {
	numNodes              int
	id                    int
	hostMap               map[int]string
	clientMap             map[int]*rpc.Client
	proposalNum           chan int
	storage               map[string]interface{}
	acceptedVals          map[string]interface{}
	acceptedProps         map[string]int
	lastSeenProposeNumber map[string]int
	clientMutex           *sync.RWMutex
	storageMutex          *sync.RWMutex
	accpetedValsMutex     *sync.RWMutex
	proposeMutex          *sync.Mutex
}

// NewPaxosNode creates a new PaxosNode. This function should return only when
// all nodes have joined the ring, and should return a non-nil error if the node
// could not be started in spite of dialing the other nodes numRetries times.
//
// hostMap is a map from node IDs to their hostports, numNodes is the number
// of nodes in the ring, replace is a flag which indicates whether this node
// is a replacement for a node which failed.
func NewPaxosNode(myHostPort string, hostMap map[int]string, numNodes, srvId, numRetries int, replace bool) (PaxosNode, error) {
	pn := &paxosNode{
		hostMap:               hostMap,
		numNodes:              numNodes,
		id:                    srvId,
		clientMap:             make(map[int]*rpc.Client),
		acceptedVals:          make(map[string]interface{}),
		acceptedProps:         make(map[string]int),
		lastSeenProposeNumber: make(map[string]int),
		storage:               make(map[string]interface{}),
		clientMutex:           &sync.RWMutex{},
		storageMutex:          &sync.RWMutex{},
		accpetedValsMutex:     &sync.RWMutex{},
		proposeMutex:          &sync.Mutex{},
	}
	log.Printf("Total Number of Nodes:%d, Created srvId:%d \n", numNodes, srvId)
	listener, err := net.Listen("tcp", myHostPort)
	if err != nil {
		log.Println("Err")
		return nil, err
	}

	err = rpc.RegisterName("PaxosNode", paxosrpc.Wrap(pn))
	if err != nil {
		return nil, err
	}

	rpc.HandleHTTP()
	go http.Serve(listener, nil)

	// Connect other nodes in ring
	for i := 0; i < numNodes; i++ {
		numTries := 0
		port := hostMap[i]
		cli, err := rpc.DialHTTP("tcp", port)

		for err != nil {
			log.Println("Failed to dial http, retry")
			numTries++
			if numTries > numRetries {
				return nil, err
			}
			time.Sleep(time.Second)
			cli, err = rpc.DialHTTP("tcp", port)
		}
		pn.clientMap[i] = cli
	}

	if replace {
		//fmt.Println("coming into replace")
		catchupDone := make(chan *rpc.Call, numNodes)
		for i := 0; i < numNodes; i++ {
			if i != pn.id {
				arg := &paxosrpc.ReplaceCatchupArgs{}
				reply := &paxosrpc.ReplaceCatchupReply{}
				cli := pn.clientMap[i]
				cli.Go("PaxosNode.RecvReplaceCatchup", arg, reply, catchupDone)
				break
			}
		}

		//for cnt := 0; cnt < numNodes-1; cnt++ {
		call := <-catchupDone
		reply := call.Reply.(*paxosrpc.ReplaceCatchupReply)
		var m map[string]interface{}
		err = json.Unmarshal(reply.Data, &m)
		if err != nil {
			//fmt.Println(err)
			return nil, err
		}
		//fmt.Println(string(reply.Data))
		for k, v := range m {
			pn.storage[k] = v
		}
		//break
		//}

		replaceDone := make(chan *rpc.Call, numNodes)
		for i := 0; i < numNodes; i++ {
			if i != pn.id {
				arg := &paxosrpc.ReplaceServerArgs{
					SrvID:    pn.id, // Server being replaced
					Hostport: myHostPort,
				}
				reply := &paxosrpc.ReplaceServerReply{}
				cli := pn.clientMap[i]
				cli.Go("PaxosNode.RecvReplaceServer", arg, reply, replaceDone)
			}
		}

		for cnt := 0; cnt < numNodes-1; cnt++ {
			<-replaceDone
		}

		//pn.acceptedVals = temp_pn.acceptedVals
		//pn.acceptedProps = temp_pn.acceptedProps
		//pn.lastSeenProposeNumber = temp_pn.lastSeenProposeNumber
		//break
	} //else {
	//pn.acceptedVals = make(map[string]interface{})
	//pn.acceptedProps = make(map[string]int)
	//pn.lastSeenProposeNumber = make(map[string]int)
	//}

	return pn, nil
}

func (pn *paxosNode) GetNextProposalNumber(args *paxosrpc.ProposalNumberArgs, reply *paxosrpc.ProposalNumberReply) error {
	pn.proposeMutex.Lock()
	lastSeenProposeNum, ok := pn.lastSeenProposeNumber[args.Key]
	if !ok {
		//fmt.Printf("GetNextProposalNumber, id:%d lastSeenPPNum: none\n", pn.id)
		reply.N = pn.id
		//pn.lastSeenProposeNumber[args.Key] = pn.id
	} else {
		//fmt.Printf("GetNextProposalNumber, id:%d lastSeenPPNum:%d\n", pn.id, lastSeenProposeNum)
		lastSeenProposeNum = (lastSeenProposeNum/pn.numNodes+1)*pn.numNodes +
			pn.id
		pn.lastSeenProposeNumber[args.Key] = lastSeenProposeNum
		reply.N = lastSeenProposeNum
		//fmt.Printf("return GetNextProposalNumber:%d\n", lastSeenProposeNum)

	}
	pn.proposeMutex.Unlock()
	return nil
}

func (pn *paxosNode) Propose(args *paxosrpc.ProposeArgs, reply *paxosrpc.ProposeReply) error {
	timeout := make(chan bool, 1)
	go func() {
		select {
		case <-time.After(TIMEOUT_LIMIT * time.Second):
			timeout <- true
			break
		}
	}()
	N := args.N
	//var ok bool
	Key := args.Key
	V := args.V
	//fmt.Printf("Propose Number: %d; Key: %s \n", N, Key)
	prepare_count := 0
	prepareDone := make(chan *rpc.Call, pn.numNodes)
	for i := 0; i < pn.numNodes; i++ {
		pn.clientMutex.RLock()
		cli := pn.clientMap[i]
		pn.clientMutex.RUnlock()
		arg := &paxosrpc.PrepareArgs{Key: Key, N: N}
		var replys = &paxosrpc.PrepareReply{}
		// change to asynchronized version?
		cli.Go("PaxosNode.RecvPrepare", arg, replys, prepareDone)
	}
	largestNa := -1
	for cnt := 0; cnt < pn.numNodes; cnt++ {
		select {
		case call := <-prepareDone:
			//fmt.Println("return from rpc")
			reply := call.Reply.(*paxosrpc.PrepareReply)
			if reply.Status == paxosrpc.OK {
				prepare_count += 1
				if reply.V_a != nil && reply.N_a > largestNa {
					largestNa = reply.N_a
					V = reply.V_a
				}
			}
		case <-timeout:
			//fmt.Println("Propose: RecvPrepare Timeout")
			return errors.New("Timeout")
		}
	}
	//fmt.Printf("Prepare, Number of OK received: %d\n", prepare_count)
	if prepare_count <= pn.numNodes/2 {
		//fmt.Printf("Prepre reject, prepare_count:%d\n", prepare_count)
		reply.V = nil
		return errors.New("Rejected!")
	}
	//fmt.Println("Starting accept phase")
	acceptDone := make(chan *rpc.Call, pn.numNodes)
	accept_count := 0
	for i := 0; i < pn.numNodes; i++ {
		pn.clientMutex.RLock()
		//fmt.Printf("start getting %d th connection...", i)
		cli := pn.clientMap[i]
		//fmt.Println("done")
		pn.clientMutex.RUnlock()
		arg := &paxosrpc.AcceptArgs{Key: Key, N: N, V: V}
		replys := &paxosrpc.AcceptReply{}
		//fmt.Printf("start calling %d th connection...", i)
		cli.Go("PaxosNode.RecvAccept", arg, replys, acceptDone)
		//fmt.Println("done")
	}

	for cnt := 0; cnt < pn.numNodes; cnt++ {
		select {
		case call := <-acceptDone:
			//fmt.Println("received accept done channel")
			/*if call.Error != nil {
				fmt.Println(call.Error)
			}*/
			reply := call.Reply.(*paxosrpc.AcceptReply)
			if reply.Status == paxosrpc.OK {
				accept_count += 1
			}
		case <-timeout:
			//fmt.Println("Propose: RecvAccept Timeout")
			return errors.New("timeout")
		}
	}
	//fmt.Printf("Accept, Number of OK received: %d\n", accept_count)

	if accept_count <= pn.numNodes/2 {
		reply.V = nil
		//fmt.Printf("Accept reject, accept_count:%d\n", accept_count)
		return errors.New("Rejected!")
	}
	//fmt.Println("Start commit phase")
	commitDone := make(chan *rpc.Call, pn.numNodes)

	for i := 0; i < pn.numNodes; i++ {
		//pn.clientMutex.RLock()
		cli := pn.clientMap[i]
		//pn.clientMutex.RUnlock()
		arg := &paxosrpc.CommitArgs{Key: Key, V: V}
		replys := &paxosrpc.CommitReply{}
		cli.Go("PaxosNode.RecvCommit", arg, replys, commitDone)
	}
	for cnt := 0; cnt < pn.numNodes; cnt++ {
		select {
		case <-commitDone:
			continue
		case <-timeout:
			//fmt.Println("Propose: RecvCommit Timeout")
			return errors.New("timeout")
		}
	}

	reply.V = V

	return nil
}

func (pn *paxosNode) GetValue(args *paxosrpc.GetValueArgs, reply *paxosrpc.GetValueReply) error {
	val, ok := pn.storage[args.Key]
	if !ok {
		reply.Status = paxosrpc.KeyNotFound
		return nil
	}
	reply.V = val
	reply.Status = paxosrpc.KeyFound
	return nil
}

func (pn *paxosNode) RecvPrepare(args *paxosrpc.PrepareArgs, reply *paxosrpc.PrepareReply) error {
	//fmt.Printf("id:%d, RecvPrepare key:%s \n", pn.id, args.Key)
	pn.proposeMutex.Lock()
	lastSeenProposalNum, ok := pn.lastSeenProposeNumber[args.Key]
	//pn.proposeMutex.Unlock()
	reply.N_a = -1
	if !ok { // haven't seen any proposal yet
		//pn.proposeMutex.Lock()
		pn.lastSeenProposeNumber[args.Key] = args.N
		pn.proposeMutex.Unlock()
		//fmt.Printf("id:%d, LastSeenProposalNumber: %d; current_proposal: %d, return ok\n", pn.id, -1, args.N)
		reply.Status = paxosrpc.OK
	} else if lastSeenProposalNum <= args.N {
		if args.N > lastSeenProposalNum {
			//pn.proposeMutex.Lock()
			pn.lastSeenProposeNumber[args.Key] = args.N
			//pn.proposeMutex.Unlock()
		}
		pn.proposeMutex.Unlock()
		var accptedVal interface{}
		pn.accpetedValsMutex.Lock()
		accptedVal, ok = pn.acceptedVals[args.Key]
		if ok {
			//fmt.Printf("id:%d, Previous accpeted value: %d\n", pn.id, accptedVal.(uint32))
			reply.V_a = accptedVal
			reply.N_a = pn.acceptedProps[args.Key]
		}
		reply.Status = paxosrpc.OK
		pn.accpetedValsMutex.Unlock()
		//fmt.Printf("id:%d, LastSeenProposalNumber: %d; current_proposal: %d, return ok\n", pn.id, lastSeenProposalNum, args.N)
	} else {
		pn.proposeMutex.Unlock()
		//fmt.Printf("id:%d, LastSeenProposalNumber: %d; current_proposal: %d, return reject\n", pn.id, lastSeenProposalNum, args.N)
		reply.Status = paxosrpc.Reject
	}
	//pn.proposeMutex.Unlock()
	return nil
}

func (pn *paxosNode) RecvAccept(args *paxosrpc.AcceptArgs, reply *paxosrpc.AcceptReply) error {
	//*reply = paxosrpc.AcceptReply{}
	//fmt.Printf("id:%d, RecvAccept key:%s \n", pn.id, args.Key)
	pn.proposeMutex.Lock()
	lastSeenProposalNum, ok := pn.lastSeenProposeNumber[args.Key]
	if !ok {
		pn.lastSeenProposeNumber[args.Key] = args.N
		pn.proposeMutex.Unlock()
		//fmt.Printf("id:%d, LastSeenProposalNumber: %d; current_proposal: %d return ok\n", pn.id, -1, args.N)
		pn.accpetedValsMutex.Lock()
		pn.acceptedVals[args.Key] = args.V
		pn.acceptedProps[args.Key] = args.N
		reply.Status = paxosrpc.OK
		pn.accpetedValsMutex.Unlock()
		//return errors.New("Impossible to happen!")
	} else if lastSeenProposalNum > args.N {
		pn.proposeMutex.Unlock()
		//fmt.Printf("id:%d, LastSeenProposalNumber: %d; current_proposal: %d return reject\n", pn.id, lastSeenProposalNum, args.N)
		reply.Status = paxosrpc.Reject
	} else {
		pn.lastSeenProposeNumber[args.Key] = args.N
		pn.proposeMutex.Unlock()
		//fmt.Printf("id:%d, LastSeenProposalNumber: %d; current_proposal: %d return ok\n", pn.id, lastSeenProposalNum, args.N)
		pn.accpetedValsMutex.Lock()
		pn.acceptedVals[args.Key] = args.V
		pn.acceptedProps[args.Key] = args.N
		reply.Status = paxosrpc.OK
		pn.accpetedValsMutex.Unlock()
	}
	//pn.proposeMutex.Unlock()
	return nil
}

func (pn *paxosNode) RecvCommit(args *paxosrpc.CommitArgs, reply *paxosrpc.CommitReply) error {
	pn.storageMutex.Lock()
	*reply = paxosrpc.CommitReply{}
	//fmt.Printf("id:%d, RecvCommit key:%s \n", pn.id, args.Key)
	pn.storage[args.Key] = args.V
	pn.accpetedValsMutex.Lock()
	delete(pn.acceptedVals, args.Key)
	delete(pn.acceptedProps, args.Key)
	pn.accpetedValsMutex.Unlock()
	pn.storageMutex.Unlock()
	return nil
}

func (pn *paxosNode) RecvReplaceServer(args *paxosrpc.ReplaceServerArgs, reply *paxosrpc.ReplaceServerReply) error {
	replaceSrvId := args.SrvID
	replaceHostport := args.Hostport
	pn.clientMutex.Lock()
	cli, ok := pn.clientMap[replaceSrvId]
	if ok {
		cli.Close()
	}
	var err error
	cli, err = rpc.DialHTTP("tcp", replaceHostport)
	if err != nil {
		return err
	}
	pn.clientMap[replaceSrvId] = cli
	pn.hostMap[replaceSrvId] = replaceHostport
	pn.clientMutex.Unlock()
	return nil
}

func (pn *paxosNode) RecvReplaceCatchup(args *paxosrpc.ReplaceCatchupArgs, reply *paxosrpc.ReplaceCatchupReply) error {
	pn.proposeMutex.Lock()
	encoded, err := json.Marshal(pn.storage)
	//fmt.Printf("RecvReplaceCatchup id:%d", pn.id)
	if err != nil {
		return err
	}
	//fmt.Println(string(encoded))
	pn.proposeMutex.Unlock()
	reply.Data = encoded
	return nil
}
