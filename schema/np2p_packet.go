package schema

import (
	"github.com/vmihailenco/msgpack/v5"
	"github.com/weaveworks/mesh"
)

type EncodableAndMergeable interface {
	Encode() [][]byte
	Merge(data mesh.GossipData) mesh.GossipData
}

// Np2pPacket is an implementation of GossipData
type Np2pPacket struct {
	Events []*Np2pEvent
	Reqs   []*Np2pReq
}

//// Np2pPacket implements GossipData.
//var _ mesh.GossipData = &Np2pPacket{}

// Np2pPacket implements EncodableAndMergeable.
var _ EncodableAndMergeable = &Np2pPacket{}

// Construct an empty Np2pPacket object, ready to receive updates.
// This is suitable to use at program start.
// Other peers will populate us with data.
func NewNp2pPacket(events *[]*Np2pEvent, req *[]*Np2pReq) *Np2pPacket {
	var events_ []*Np2pEvent = nil
	if events != nil {
		events_ = *events
	}
	var req_ []*Np2pReq = nil
	if req != nil {
		req_ = *req
	}

	return &Np2pPacket{
		Events: events_,
		Reqs:   req_,
	}
}

func NewNp2pPacketFromBytes(data []byte) (*Np2pPacket, error) {
	var bp Np2pPacket
	err := msgpack.Unmarshal(data, &bp)
	if err != nil {
		return nil, err
	}

	return &bp, nil
}

// Encode serializes Np2pPacket to a slice of byte-slices.
func (pkt *Np2pPacket) Encode() [][]byte {
	b, err := msgpack.Marshal(pkt)
	if err != nil {
		panic(err)
	}
	return [][]byte{b}
}

// Merge merges the other GossipData into this one,
// and returns our resulting, complete Np2pPacket.
func (st *Np2pPacket) Merge(other mesh.GossipData) (complete mesh.GossipData) {
	if st.Events != nil && other.(*Np2pPacket).Events != nil {
		st.Events = append(st.Events, other.(*Np2pPacket).Events...)
	} else if st.Events == nil && other.(*Np2pPacket).Events != nil {
		st.Events = other.(*Np2pPacket).Events
	}
	if st.Reqs != nil && other.(*Np2pPacket).Reqs != nil {
		st.Reqs = append(st.Reqs, other.(*Np2pPacket).Reqs...)
	} else if st.Reqs == nil && other.(*Np2pPacket).Reqs != nil {
		st.Reqs = other.(*Np2pPacket).Reqs
	}
	return st
}
