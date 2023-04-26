package Ring
import (
	"crypto/ed25519"

)

type Ring struct{
	RingId				uint64
	Nodes				[]Node
	LastRing			uint64
}


var{
	Rings 			map[uint64]ring        	//Known Rings
}

func getRing(RingID uint64)(Ring)	{
	return Rings[RingID]

}

type rings struct {
	allRings map[uint64]Ring
}


var Rings = rings{
	allRings 			map[uint64]ring,        	//Known Rings
}


	
func (r *rings) GetRing(ringId uint64) Ring	{
	return r.allRings[ringId]

}


func (r *rings) GetRings() map[uint64]Ring {
	return r.allRings
}

func (r *rings) AddNode(newRing Ring) {
	r.allRings[newRing.RingId] = newRing
}
