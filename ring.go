package Ring
import (
	"github.com/vidarsolutions/Node"
)

type Ring struct{
	RingId				uint64
	Nodes				[]Node
	LastRing			uint64
}


var Rings = rings{
	allRings: 			make(map[uint64]Ring),        	//Known Rings
}



type rings struct {
	allRings map[uint64]Ring
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
