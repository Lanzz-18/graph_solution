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
		fmt.Println()
		os.Exit(1)
	}
	filePath := os.Args[1]

	// Parsing the file to make the graph
	adjacentList, inDegree, outDegree, allNodes, err := parseGraph(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	fmt.Printf("is_dag: %t\n", checkIsDAG(adjacentList, inDegree, allNodes))
	fmt.Printf("max_in_degree: %d\n", getMaxInDegree(inDegree, allNodes))
	fmt.Printf("max_out_degree: %d\n", getMaxOutDegree(outDegree, allNodes))
}

// parseGraph reads the CSV file and processes data
// It returns adjacency list, in-degrees, out-degrees, and all the set of nodes
func parseGraph(filePath string) (map[int]int, map[int]int, map[int]int, map[int]bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer file.Close() 

	// adjacentList
	adjacencyList := make(map[int][]int)

	// inDegree
	inDegree := make(map[int]int)

	// outDegree
	outDegree := make(map[int]int)

	// allNodes
	allNodes := make(map[int]bool)

	// Reading the file line by line
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
			return nil, nil, nil, nil, fmt.Errorf("Error in file: %s", line)
		}

		// Converting the string numbers to integers
		src, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		dst, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err1 != nil || err2 != nil {
			return nil, nil, nil, nil, fmt.Errorf("cannot read nodes on line: %s", line)
		}

		// Adding nodes to set
		allNodes[src] = true
		allNodes[dst] = true

		adjacencyList[src] = append(adjacencyList[src], dst)
		outDegree[src]++
		inDegree[dst]++
	}
	return adjacencyList, inDegree, outDegree, allNodes, scanner.Err()
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
		queue + queue[1:]
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
