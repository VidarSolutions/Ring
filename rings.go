package Ring

import(
	"time"
	"github.com/vidarsolutions/Node"
)

type Bytes32 [32]byte
type Bytes64 [64]byte

var Rings = rings{
	AllRings: 			make(map[uint64]Ring),        	//Known Rings
	RingMasters: 		make(map[uint64]Node.VidarNode),			//Map of Nodes allowed to generate new ring and node ids.
	Update:				time.Now(),
	NodeIDs:			0,
	LastRing			0,
}



type rings struct {
	AllRings map[uint64]Ring
	RingMasters map[uint64]Node.VidarNode
	Update		time.Time
	NodeIDs		uint64
	LastRing	uint64
}


	
func (r *rings) GetRing(ringId uint64) Ring	{
	return r.AllRings[ringId]

}


func (r *rings) GetRings() map[uint64]Ring {
	return r.AllRings
}

func (r *rings) AddRing(newRing Ring) {
	r.AllRings[newRing.RingId] = newRing
}


func (r *rings)getNewNodeId() uint64 {
    r.NodeIDs += 1
	
    return r.NodeIDs
}
func (r *rings)getRingPeers(node *Node.VidarNode) ([]uint64, Node.VidarNode) {
		knownRings := r.GetRings()
		const ringSize = 7
		node.NodeID = r.getNewNodeId()
		// Determine the ring number of the node based on its ID
		ringNum := node.NodeID / ringSize
		node.RingId = uint64(ringNum)
		nextRing := ringNum + 1

		// Create a slice to store the rings that this node will learn about and backup
		var peersRing []uint64
		

		if len(knownRings) < 50 {
			for _, ring := range knownRings {
				peersRing = append(peersRing, ring.RingId)
			}
		} else {
			peersRing = append(peersRing, knownRings[nextRing].RingId)
			ringIndex := node.NodeID - 1
			// Loop through the next 7 rings to assign peers to
			for i := 0; i < ringSize; i++ {
				// Calculate the index of the ring in the knownRings slice
				ringIndex = ringIndex + 7
				// Add the ring ID to the peersRing slice
				if ringIndex+1 < len(knownRings) {
					peersRing = append(peersRing, knownRings[ringIndex].RingId)
				} else {
					ringIndex = 0
					peersRing = append(peersRing, knownRings[ringIndex].RingId)
				}
			}
		}

		// Return the peersRing slice
		return peersRing, *node
}

func (r *rings)newRing(node *Node.VidarNode, sig Bytes32, msg string) uint64 {
	rm, found := r.ringMaster[node.NodeID]		//only ring masters may call this function
	if found{
		if r.isRingMaster(node, sig, msg){
			r.LastRing += 1
			ring := Ring{}
			ring.RingId =r.LastRing
			ring.LastRing =r.LastRing
			nodes := []Node.VidarNode{node}
			
			ring.Nodes = nodes
			r.AddRing(ring)
		}
	}
    return r.LastRing
}
func (r *rings)isRingMaster(node *Node.VidarNode, sig Bytes32, msg string) bool{
	var nodeId = node.NodeID
	pubKey = node.PubKey
	validMsg = r.lastRing+nodeId
	var rm bool
	rm := false 
	m, err := strconv.ParseInt(msg, 10, 64)
	if m==validMsg{
		//Check if Node signature is valid
		today := time.Now()
		tooLong = today.Add(-15 * time.Minute)
		if r.Update > tooLong{
			//Run RingMasterUpdate
			r.RingMasterUpdate()
		}
		//add code to verify signature
		rm = true
	}
	return rm
}

func (r *rings) RingMasterUpdate(){
	for k, rm := range r.Rings.ringMaster {
		//dial out over tor to sync rings with ringmasters
		
	}
	r.Update = time.Now()
}

func (r *rings)saveRings(){
// Encode the rings map as JSON
    ringsJSON, err := json.Marshal(r)
    if err != nil {
        panic(err)
    }

    // Write the JSON string to a file
    err = ioutil.WriteFile("rings.json", ringsJSON, 0644)
    if err != nil {
        panic(err)
    }

    fmt.Println("Rings have been saved to rings.json")

    

}

func (r *rings)loadRings() map[uint64]Ring{
// Read the rings from the file
    fileBytes, err := ioutil.ReadFile("rings.json")
    if err != nil {
       return GetRings()
    }

    // Decode the JSON string into a Rings struct
    var savedRings = Ring.Rings.GetRings()
    err = json.Unmarshal(fileBytes, &savedRings)
    if err != nil {
        return GetRings()
    }
	
	fmt.Println("Rings loaded from file:")
	
	if(Logging){
		fmt.Printf("%+v\n", savedRings)
	}
	return savedRings
}
