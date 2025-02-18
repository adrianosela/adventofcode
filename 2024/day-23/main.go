package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"os"
	"sort"
	"strings"

	"github.com/adrianosela/adventofcode/utils/set"
)

type node struct {
	id    string
	peers map[string]*node
}

type network struct {
	nodes map[string]*node
}

func main() {
	debug := false

	sampleInputNetwork, err := loadInput("sample-input.txt")
	if err != nil {
		log.Fatalf("failed to load sample input data: %v", err)
	}
	if debug {
		fmt.Println(sampleInputNetwork.String())
	}
	log.Printf(
		"[Answer to Sample in Part 1] The number of triplets is: %d (should be 12), and %d have a node that starts with 't' (should be 7)",
		len(sampleInputNetwork.getAllTriplets()),
		part1(sampleInputNetwork),
	)

	inputNetwork, err := loadInput("input.txt")
	if err != nil {
		log.Fatalf("failed to load sample input data: %v", err)
	}
	log.Printf("[Answer to Part 1] The result is: %d", part1(inputNetwork))

	log.Printf("[Answer to Sample in Part 2] Password is %s (should be 'co,de,ka,ta')", part2(sampleInputNetwork))
	log.Printf("[Answer to Part 2] Password is %s", part2(inputNetwork))
}

func loadInput(filename string) (*network, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	nodes := make(map[string]*node)
	scanner := bufio.NewScanner(file)
	for lineNo := 0; scanner.Scan(); lineNo++ {
		line := scanner.Text()

		nodeIDA, nodeIDB, ok := strings.Cut(line, "-")
		if !ok {
			return nil, fmt.Errorf("invalid input in line %d: \"%s\" not in the form ${peer_a}-${peer_b}", lineNo, line)
		}

		nodeA, ok := nodes[nodeIDA]
		if !ok {
			nodes[nodeIDA] = &node{
				id:    nodeIDA,
				peers: make(map[string]*node),
			}
			nodeA = nodes[nodeIDA]
		}

		nodeB, ok := nodes[nodeIDB]
		if !ok {
			nodes[nodeIDB] = &node{
				id:    nodeIDB,
				peers: make(map[string]*node),
			}
			nodeB = nodes[nodeIDB]
		}

		nodeA.peers[nodeIDB] = nodeB
		nodeB.peers[nodeIDA] = nodeA
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan input file: %v", err)
	}

	return &network{nodes: nodes}, nil
}

func (n *network) String() string {
	asMap := map[string][]string{}
	for id, node := range n.nodes {
		peerIDs := []string{}
		for id := range maps.Keys(node.peers) {
			peerIDs = append(peerIDs, id)
		}
		sort.Strings(peerIDs)
		asMap[id] = peerIDs
	}
	byt, err := json.Marshal(&asMap)
	if err != nil {
		// cannot ever happen
		panic(fmt.Errorf("failed to json encode network: %v", err))
	}
	return string(byt)
}

func (n *network) getAllTriplets() [][]string {
	triplets := [][]string{}
	visited := set.New[string]()

	for _, node := range n.nodes {
		// we can skip nodes with less than two peers
		// because we know they will never form a triplet.
		if len(node.peers) < 2 {
			continue
		}

		for firstPeerID, peerNode := range node.peers {
			for secondPeerID := range peerNode.peers {
				if secondPeerID != node.id && node.peers[secondPeerID] != nil {
					triplet := []string{node.id, firstPeerID, secondPeerID}
					sort.Strings(triplet)
					tripletKey := fmt.Sprintf("%v", triplet)
					if !visited.Has(tripletKey) {
						visited.Put(tripletKey)
						triplets = append(triplets, triplet)
					}
				}
			}
		}
	}
	return triplets
}

func (n *network) isClique(ids ...string) bool {
	for i, idA := range ids {
		for j, idB := range ids {
			if i == j {
				continue
			}
			if _, ok := n.nodes[idA].peers[idB]; !ok {
				return false
			}
		}
	}
	return true
}

func (n *network) findLargestClique() []string {
	largestClique := []string{}

	var findCliques func(currentClique []string, nodesLeft []string, nodesExcluded []string)
	findCliques = func(currentClique []string, nodesLeft []string, nodesExcluded []string) {
		if len(nodesLeft) == 0 && len(nodesExcluded) == 0 {
			if len(currentClique) > len(largestClique) {
				largestClique = make([]string, len(currentClique))
				copy(largestClique, currentClique)
			}
		} else if len(nodesLeft) != 0 {
			considering := nodesLeft[0]
			newClique := append(currentClique, considering)
			if n.isClique(newClique...) {
				newNodesLeft := []string{}
				for _, nodeID := range nodesLeft[1:] {
					if _, ok := n.nodes[considering].peers[nodeID]; ok {
						newNodesLeft = append(newNodesLeft, nodeID)
					}
				}
				newNodesExcluded := []string{}
				for _, nodeID := range nodesExcluded {
					if _, ok := n.nodes[considering].peers[nodeID]; ok {
						newNodesExcluded = append(newNodesExcluded, nodeID)
					}
				}
				findCliques(newClique, newNodesLeft, newNodesExcluded)
			}
			findCliques(currentClique, nodesLeft[1:], append(nodesExcluded, considering))
		}
	}

	nodeIds := make([]string, 0, len(n.nodes))
	for id := range n.nodes {
		nodeIds = append(nodeIds, id)
	}

	findCliques([]string{}, nodeIds, []string{})

	return largestClique
}

func part1(n *network) int {
	startsWithT := 0
	for _, triplet := range n.getAllTriplets() {
		for _, nodeID := range triplet {
			if strings.HasPrefix(nodeID, "t") {
				startsWithT++
				break
			}
		}
	}
	return startsWithT
}

func part2(n *network) string {
	ids := n.findLargestClique()
	sort.Strings(ids)
	return strings.Join(ids, ",")
}
