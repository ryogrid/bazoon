package core

import (
	"fmt"
	"log"

	"github.com/ryogrid/nostrp2p/np2p_const"
	"github.com/ryogrid/nostrp2p/np2p_util"
	"github.com/ryogrid/nostrp2p/schema"
)

// Peer encapsulates state and implements mesh.Gossiper.
// It should be passed to mesh.Router.NewGossip,
// and the resulting Gossip registered in turn,
// before calling mesh.Router.Start.
type Np2pPeer struct {
	//send            *mesh.Gossip
	Actions    chan<- func()
	quit       chan struct{}
	logger     *log.Logger
	dataMan    *DataManager
	MessageMan *MessageManager
	SelfId     uint64 //mesh.PeerName
	SelfPubkey [np2p_const.PubkeySize]byte
	//Router          *mesh.Router // TODO: need to modify
	recvedEvtReqMap map[uint64]struct{}
}

// Construct a Np2pPeer with empty state.
// Be sure to Register a channel, later,
// so we can make outbound communication.
func NewPeer(self uint64, logger *log.Logger) *Np2pPeer {
	actions := make(chan func())
	dataMan := NewDataManager()
	msgMan := NewMessageManager(dataMan)

	p := &Np2pPeer{
		//send:            nil, // must .Register() later
		Actions:         actions,
		quit:            make(chan struct{}),
		logger:          logger,
		dataMan:         dataMan,
		MessageMan:      msgMan,
		SelfId:          self,
		recvedEvtReqMap: make(map[uint64]struct{}), // make(map[uint64]struct{}),
	}
	go p.loop(actions)

	// start event resender
	msgMan.evtReSender.Start()

	return p
}

func (p *Np2pPeer) loop(actions <-chan func()) {
	for {
		select {
		case f := <-actions:
			f()
		case <-p.quit:
			return
		}
	}
}

func (p *Np2pPeer) stop() {
	close(p.quit)
}

func (p *Np2pPeer) OnRecvBroadcast(src uint64, buf []byte) (received schema.EncodableAndMergeable, err error) {
	//var pkt schema.Np2pPacket
	///if err_ := gob.NewDecoder(bytes.NewReader(buf)).Decode(&pkt); err_ != nil {
	pkt, err_ := schema.NewNp2pPacketFromBytes(buf)
	if err_ != nil {
		return nil, err_
	}

	tmpEvts := make([]*schema.Np2pEvent, 0)
	tmpReqs := make([]*schema.Np2pReq, 0)
	retPkt := schema.NewNp2pPacket(&tmpEvts, &tmpReqs)
	if pkt.Events != nil {
		for _, evt := range pkt.Events {
			if _, ok := p.recvedEvtReqMap[np2p_util.ExtractUint64FromBytes(evt.Id[:])]; !ok {
				if evt.Verify() == false {
					// invalid signiture
					fmt.Println("invalid signiture")
					continue
				}

				err2 := p.MessageMan.handleRecvMsgBcastEvt(src, pkt, evt)
				if err2 != nil {
					panic(err2)
				}

				p.recvedEvtReqMap[np2p_util.ExtractUint64FromBytes(evt.Id[:])] = struct{}{}
				retPkt.Events = append(retPkt.Events, evt)
			} else {
				continue
			}
		}
	} else if pkt.Reqs != nil {
		for _, req := range pkt.Reqs {
			if _, ok := p.recvedEvtReqMap[req.Id]; !ok {
				err2 := p.MessageMan.handleRecvMsgBcastReq(src, pkt, req)
				if err2 != nil {
					panic(err2)
				}

				p.recvedEvtReqMap[req.Id] = struct{}{}
				retPkt.Reqs = append(retPkt.Reqs, req)
			} else {
				continue
			}
		}
	} else {
		return pkt, nil
	}

	if len(retPkt.Events) == 0 && len(retPkt.Reqs) == 0 {
		return nil, nil
	} else {
		return retPkt, nil
	}

	//return &pkt, nil
	//return &schema.Np2pPacket{}, nil
}

func (p *Np2pPeer) OnRecvUnicast(src uint64, buf []byte) (err error) {
	//var pkt schema.Np2pPacket
	//if err_ := gob.NewDecoder(bytes.NewReader(buf)).Decode(&pkt); err_ != nil {
	//	return err_
	//}
	pkt, err := schema.NewNp2pPacketFromBytes(buf)
	if err != nil {
		return err
	}

	err_ := p.MessageMan.handleRecvMsgUnicast(src, pkt)
	if err_ != nil {
		panic(err_)
	}

	return nil
}

/*
func (p *Np2pPeer) GetPeerList() []mesh.PeerName {
	tmpMap := p.Router.Routes.PeerNames()
	retArr := make([]mesh.PeerName, 0)
	for k, _ := range tmpMap {
		retArr = append(retArr, k)
	}
	return retArr
}
*/
