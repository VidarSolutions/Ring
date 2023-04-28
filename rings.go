package Ring

import(
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"github.com/vidarsolutions/Node"
	"github.com/vidarsolutions/Transfer"
)

type Bytes32 [32]byte
type Bytes64 [64]byte

var Rings = rings{
	AllRings: 			make(map[uint64]Ring),        	//Known Rings
	RingMasters: 		make(map[uint64]Node.VidarNode),			//Map of Nodes allowed to generate new ring and node ids.
	Update:				time.Now(),
	NodeIDs:			0,
	LastRing:			0,
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
				if ringIndex+1 < uint64(len(knownRings)) {
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
	_, found := r.RingMasters[node.NodeID]		//only ring masters may call this function
	if found{
		if r.isRingMaster(node, sig, msg){
			r.LastRing += 1
			ring := Ring{}
			ring.RingId =r.LastRing
			nodes := []Node.VidarNode{*node}
			
			ring.Nodes = nodes
			r.AddRing(ring)
		}
	}
    return r.LastRing
}
func (r *rings)isRingMaster(node *Node.VidarNode, sig Bytes32, msg string) bool{
	var nodeId = node.NodeID
	var pubKey = node.PubKey
	var validMsg = r.LastRing+nodeId
	var rm bool  = false 
	m, _ := strconv.ParseUint(msg, 10, 64)
	if m==validMsg{
		//Check if Node signature is valid
		today := time.Now()
		tooLong := today.Add(-15 * time.Minute)
		if r.Update.Before(tooLong){
			//Run RingMasterUpdate
			r.RingMasterUpdate()
		}
		//add code to verify signature
		
		rm = ed25519.Verify(pubKey,[]byte(msg) , []byte(sig[:]))
	}
	return rm
}


func (r *rings) RingMasterUpdate(){
	for _, rm := range r.RingMasters {
		//dial out over tor to sync rings with ringmasters
		var t =Transfer.Dialer("127.0.0.1:9050")
		// Encode the struct as JSON
		jsonData, err := json.Marshal(r.AllRings)
		if err != nil {
			fmt.Println("Could not convert Rings to Json")
		}else{			
			 var resp *http.Response 
	         x := 0
		     for {
				x++
				if x > 10 {
					//Add code to Report Down RingMaster
					break
				}
					resp = t.Request("Post", rm.Tor, jsonData)
					if resp.StatusCode == http.StatusOK {
						break
					}else {
						time.Sleep(1 * time.Second) // Wait for 1 second before trying again
					}
				}
				
			}
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
       return r.GetRings()
    }

    // Decode the JSON string into a Rings struct
    var savedRings = r.GetRings()
    err = json.Unmarshal(fileBytes, &savedRings)
    if err != nil {
        return r.GetRings()
    }
	
	fmt.Println("Rings loaded from file:")
	
	
		fmt.Printf("%+v\n", savedRings)
	
	return savedRings
}
