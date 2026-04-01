package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main(){
	if len(os.Args) < 2 {
		fmt.Println("Enter a graph file path as an argument")
		os.Exit(1)
	}
	filePath := os.Args[1]

	// Parsing the file to make the graph
	adjacencyList, incomingEdges, inDegree, outDegree, allNodes, err := parseGraph(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	fmt.Printf("is_dag: %t\n", checkDAG(adjacencyList, inDegree, allNodes))
	fmt.Printf("max_in_degree: %d\n", getMaxInDegree(inDegree, allNodes))
	fmt.Printf("max_out_degree: %d\n", getMaxOutDegree(outDegree, allNodes))
	prMax, prMin := processPageRank(incomingEdges, outDegree, allNodes)
	fmt.Printf("pr_max: %f\n", prMax)
	fmt.Printf("pr_min: %f\n", prMin)
}

// parseGraph reads the CSV file and processes data
// It returns adjacency list, in-degrees, out-degrees, and all the set of nodes
func parseGraph(filePath string) (map[int][]int, map[int][]int, map[int]int, map[int]int, map[int]bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	defer file.Close()

	adjacencyList := make(map[int][]int)
	incomingEdges := make(map[int][]int) // incomingEdges[3] = [1, 2] means node 3 is pointed to by node 1 and node 2
	inDegree := make(map[int]int)
	outDegree := make(map[int]int)
	allNodes := make(map[int]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// skip empty lines
		if line == "" {
			continue
		}
		// Split it by the comma
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			return nil, nil, nil, nil, nil, fmt.Errorf("Error in file: %s", line)
		}

		// Converting the string numbers to integers
		src, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		dst, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err1 != nil || err2 != nil {
			return nil, nil, nil, nil, nil, fmt.Errorf("cannot read nodes on line: %s", line)
		}

		allNodes[src] = true
		allNodes[dst] = true

		adjacencyList[src] = append(adjacencyList[src], dst)
		incomingEdges[dst] = append(incomingEdges[dst], src) // track who points to dst
		outDegree[src]++
		inDegree[dst]++
	}

	return adjacencyList, incomingEdges, inDegree, outDegree, allNodes, scanner.Err()
}

func checkDAG(adjacencyList map[int][]int, inDegree map[int]int, allNodes map[int]bool) bool {
	// Copying the inDegree
	tempInDegree := make(map[int]int)
	for node := range allNodes {
		tempInDegree[node] = inDegree[node]
	}

	queue := []int{}
	for node := range allNodes {
		if tempInDegree[node] == 0 {
			queue = append(queue, node)
		}
	}
	removedCount := 0

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		removedCount++

		for _, neighbor := range adjacencyList[current] {
			tempInDegree[neighbor]--

			// If there are no incoming edges, add it to the queue
			if tempInDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	return removedCount == len(allNodes)
}

// Finding the node with the most incoming edges
func getMaxInDegree(inDegree map[int]int, allNodes map[int]bool) int {
	max := 0
	for node := range allNodes {
		if inDegree[node] > max {
			max = inDegree[node]
		}
	}
	return max
}

// Finding the node with most outgoing edges
func getMaxOutDegree(outDegree map[int]int, allNodes map[int]bool) int {
	max := 0
	for node := range allNodes {
		if outDegree[node] > max {
			max = outDegree[node]
		}
	}
	return max
}

// Calculating the PageRank algorithm
func processPageRank(incomingEdges map[int][]int, outDegree map[int]int, allNodes map[int]bool) (float64, float64) {

	N := len(allNodes)
	damping := 0.85

	// give all the nodes an equal starting value
	scores := make(map[int]float64)
	for node := range allNodes {
		scores[node] = 1.0 / float64(N)
	}

	// sort the list
	nodeList := []int{}
	for node := range allNodes {
		nodeList = append(nodeList, node)
	}

	// Run 20 iterations
	for i := 0; i < 20; i++ {

		// Collect scores from dangling nodes and spread to everyone
		danglingSum := 0.0
		for _, node := range nodeList {
			if outDegree[node] == 0 {
				danglingSum += scores[node]
			}
		}
		// Calculating new scores
		newScores := make(map[int]float64)
		for _, node := range nodeList {

			// Check if neighbor points to node
			newScores[node] = (1.0-damping)/float64(N) + damping*(danglingSum/float64(N))

			// sending equal share of the neighbour score to each node it points to
			for _, src := range incomingEdges[node] {
				newScores[node] += damping * (scores[src] / float64(outDegree[src]))
			}
		}

		scores = newScores
	}

	// Find max and min
	maxScore := 0.0
	minScore := 1.0
	for _, node := range nodeList {
		if scores[node] > maxScore {
			maxScore = scores[node]
		}
		if scores[node] < minScore {
			minScore = scores[node]
		}
	}

	return maxScore, minScore
}