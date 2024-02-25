package schema

import (
	"github.com/holiman/uint256"
	"github.com/ryogrid/buzzoon/buzz_const"
	"github.com/vmihailenco/msgpack/v5"
	"math/big"
)

type BuzzEvent struct {
	Id         uint64                      // random 64bit uint value to identify the event (Not sha256 32bytes hash)
	Pubkey     [buzz_const.PubkeySize]byte // encoded 256bit uint (holiman/uint256)
	Created_at int64                       // unix timestamp in seconds
	Kind       uint16                      // integer between 0 and 65535
	Tags       map[string][]interface{}    // Key: tag string, Value: any type
	Content    string
	Sig        *[buzz_const.SignatureSize]byte // 64-bytes integr of the signature of the sha256 hash of the serialized event data
}

func (e *BuzzEvent) GetPubkey() *big.Int {
	fixed256 := uint256.NewInt(0)
	return fixed256.SetBytes(e.Pubkey[:]).ToBig()
}

func (e *BuzzEvent) SetPubkey(pubkey *big.Int) {
	fixed256, isOverflow := uint256.FromBig(pubkey)
	if isOverflow {
		panic("overflow")
	}
	fixed256.WriteToArray32(&e.Pubkey)
}

func (e *BuzzEvent) Encode() []byte {
	b, err := msgpack.Marshal(e)
	if err != nil {
		panic(err)
	}
	return b
}

func NewBuzzEventFromBytes(b []byte) (*BuzzEvent, error) {
	var e BuzzEvent
	if err := msgpack.Unmarshal(b, &e); err != nil {
		return nil, err
	}
	return &e, nil
}
