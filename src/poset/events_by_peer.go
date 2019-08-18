package poset

import (
	"github.com/Fantom-foundation/go-lachesis/src/logger"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"strings"

	"github.com/Fantom-foundation/go-lachesis/src/hash"
)

// TODO: make EventsByPeer internal

type (
	// EventsByNode is a event hashes grouped by creator.
	// ( creator --> event hashes )
	EventsByPeer map[hash.Peer]hash.EventsSet
)

/*
 * eventsByNode's methods:
 */

// Add unions events into one.
func (ee EventsByPeer) Add(events EventsByPeer) (changed bool) {
	for creator, hashes := range events {
		if ee[creator] == nil {
			ee[creator] = hash.EventsSet{}
		}
		if ee[creator].Add(hashes.Slice()...) {
			changed = true
		}
	}
	return
}

// AddOne appends one event.
func (ee EventsByPeer) AddOne(event hash.Event, creator hash.Peer) (changed bool) {
	if ee[creator] == nil {
		ee[creator] = hash.EventsSet{}
	}
	if ee[creator].Add(event) {
		changed = true
	}
	return
}

// Contains returns true if event of node is in.
func (ee EventsByPeer) Contains(node hash.Peer, event hash.Event) bool {
	return ee[node] != nil && ee[node].Contains(event)
}

// Each returns range of all events.
func (ee EventsByPeer) Each() map[hash.Event]hash.Peer {
	res := make(map[hash.Event]hash.Peer)
	for creator, events := range ee {
		for e := range events {
			res[e] = creator
		}
	}
	return res
}

// String returns human readable string representation.
func (ee EventsByPeer) String() string {
	var ss []string
	for node, roots := range ee {
		ss = append(ss, node.String()+":"+roots.String())
	}
	return "byNode{" + strings.Join(ss, ", ") + "}"
}

type eventDescr struct {
	Creator hash.Peer
	Hash    hash.Event
}

func (ee EventsByPeer) EncodeRLP(w io.Writer) error {
	var arr []eventDescr
	for creator, hh := range ee {
		for hash_ := range hh {
			arr = append(arr, eventDescr{
				Creator: creator,
				Hash:    hash_,
			})
		}
	}
	return rlp.Encode(w, arr)
}

func (pp *EventsByPeer) DecodeRLP(s *rlp.Stream) error {
	if *pp == nil {
		*pp = EventsByPeer{}
	}
	ee := *pp

	var arr []eventDescr
	if err := s.Decode(&arr); err != nil {
		return err
	}

	for _, w := range arr {
		if ee[w.Creator] == nil {
			ee[w.Creator] = hash.EventsSet{}
		}
		if !ee[w.Creator].Add(w.Hash) {
			logger.Get().Fatal("double value is detected")
		}
	}

	return nil
}
