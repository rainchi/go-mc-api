package navigate

import (
	"github.com/rain931215/go-mc-api/api"
	"math"
)

type pathFinder struct {
	c                                     *api.Client
	startPointX, startPointY, startPointZ float64
	endPointX, endPointY, endPointZ       float64
	startPos, endPos                      pos
	openNodeList, closeNodeList           map[pos]*node
	count                                 uint8
}

func setNewPath(x, y, z float64, c *api.Client) *pathFinder {
	f := new(pathFinder)
	f.c = c
	f.startPointX = 0
	f.startPointY = 0
	f.startPointZ = 0
	f.startPos = pos{x: int(math.Floor(f.c.GetX())), y: int(math.Floor(f.c.GetY())), z: int(math.Floor(f.c.GetZ()))}
	f.endPointX = x - c.GetX()
	f.endPointY = y - c.GetY()
	f.endPointZ = z - c.GetZ()
	f.endPos = pos{x: int(math.Floor(f.endPointX)), y: int(math.Floor(f.endPointY)), z: int(math.Floor(f.endPointZ))}
	f.openNodeList = make(map[pos]*node)
	f.closeNodeList = make(map[pos]*node)
	pos := pos{x: 0, y: 0, z: 0}
	f.openNodeList[pos] = newNode(pos, new(node))
	f.count++
	return f
}

func (f *pathFinder) getNodes() []*node {
	var nodes = make([]*node, 1)
	for {
		if f.count < 1 {
			println("wrong")
			return nodes
		}
		var (
			FList   []uint8
			getNode = make(map[uint8]*node)
		)
		for _, node := range f.openNodeList {
			F := node.cost + node.getGuessCost(f.endPos)
			FList = append(FList, F)
			getNode[F] = node
		}
		thisNode := getNode[min(FList)]
		nodePos := thisNode.pos
		if thisNode.pos == f.endPos {
			println("finish")
			nodes = append(nodes, thisNode)
			nodes = thisNode.returnNodes(nodes)
			return nodes
		}
		delete(f.openNodeList, nodePos)
		f.count--
		f.closeNodeList[nodePos] = thisNode
		for x := -1; x < 2; x += 2 {
			pos := pos{x: x + nodePos.x, y: nodePos.y, z: nodePos.z}
			if f.nodeRule(thisNode, pos) {
				f.openNodeList[pos] = newNode(pos, thisNode)
				f.count++
			}
		}
		for y := -1; y < 2; y += 2 {
			pos := pos{x: nodePos.x, y: y + nodePos.y, z: nodePos.z}
			if f.nodeRule(thisNode, pos) {
				f.openNodeList[pos] = newNode(pos, thisNode)
				f.count++
			}
		}
		for z := -1; z < 2; z += 2 {
			pos := pos{x: nodePos.x, y: nodePos.y, z: z + nodePos.z}
			if f.nodeRule(thisNode, pos) {
				f.openNodeList[pos] = newNode(pos, thisNode)
				f.count++
			}
		}
	}
}

func (f *pathFinder) nodeRule(node *node, p pos) bool {
	var pass bool
	if _, ok := f.closeNodeList[p]; ok == true {
		return false
	}
	x := f.startPos.x + p.x
	y := f.startPos.y + p.y
	z := f.startPos.z + p.z
	//println(x, y, z, f.c.World.GetBlockStatus(x, y, z))
	if (y < 0 || f.c.World.GetBlockStatus(x, y, z) == 0) && (y+1 < 0 || f.c.World.GetBlockStatus(x, y+1, z) == 0) {
		pass = true
	}
	return pass
}

func min(l []uint8) (min uint8) {
	min = l[0]
	for _, v := range l {
		if v < min {
			min = v
		}
	}
	return
}