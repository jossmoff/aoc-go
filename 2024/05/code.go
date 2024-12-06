package main

import (
	"bufio"
	"io"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/dominikbraun/graph"
	"github.com/jpillora/puzzler/harness/aoc"
)

func main() {
	aoc.Harness(run)
}

func parseInput(r io.Reader) (graph.Graph[int, int], [][]int) {
	g := graph.New(graph.IntHash, graph.Directed())
	var updateSequences [][]int
	scanner := bufio.NewScanner(r)
	parsingDependencies := true

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			parsingDependencies = false
			continue
		}

		if parsingDependencies {
			// Parse dependencies
			parts := strings.Split(line, "|")
			if len(parts) != 2 {
				continue
			}
			left, err1 := strconv.Atoi(parts[0])
			right, err2 := strconv.Atoi(parts[1])
			if err1 != nil || err2 != nil {
				continue
			}
			_ = g.AddVertex(left)
			_ = g.AddVertex(right)
			_ = g.AddEdge(left, right)
		} else {
			// Parse update sequences
			parts := strings.Split(line, ",")
			var sequence []int
			for _, part := range parts {
				page, err := strconv.Atoi(part)
				if err != nil {
					continue
				}
				sequence = append(sequence, page)
			}
			updateSequences = append(updateSequences, sequence)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return g, updateSequences
}

func hasNodeTraversal(g graph.Graph[int, int], nodeSequence []int) bool {
	for i := 0; i < len(nodeSequence)-1; i++ {
		if _, err := g.Edge(nodeSequence[i], nodeSequence[i+1]); err != nil {
			return false
		}
	}
	return true
}

func sortRelativeTo(reference, toSort []int) []int {
	orderMap := make(map[int]int)
	for index, value := range reference {
		orderMap[value] = index
	}

	sorted := make([]int, len(toSort))
	copy(sorted, toSort)

	sort.Slice(sorted, func(i, j int) bool {
		return orderMap[sorted[i]] < orderMap[sorted[j]]
	})

	return sorted
}

func getSubgraph(g graph.Graph[int, int], nodes []int) graph.Graph[int, int] {
	subgraph := graph.New(graph.IntHash, graph.Directed())

	// Add nodes to the subgraph
	for _, node := range nodes {
		_ = subgraph.AddVertex(node)
	}

	adjacencyMap, _ := g.AdjacencyMap()

	// Add edges to the subgraph
	for _, fromNode := range nodes {
		if edgesMap, exists := adjacencyMap[fromNode]; exists {
			for toNode, _ := range edgesMap {
				if slices.Contains(nodes, toNode) {
					_ = subgraph.AddEdge(fromNode, toNode)
				}
			}
		}
	}

	return subgraph
}

func part1(input string) int {
	g, updateSequences := parseInput(strings.NewReader(input))
	var validMidpoints []int

	for _, sequence := range updateSequences {
		if hasNodeTraversal(g, sequence) {
			validMidpoints = append(validMidpoints, sequence[len(sequence)/2])
		}
	}

	return sum(validMidpoints)
}

func part2(input string) int {
	dag, updateSequences := parseInput(strings.NewReader(input))
	var reorderedMidpoints []int

	for _, sequence := range updateSequences {
		if !hasNodeTraversal(dag, sequence) {
			subGraph := getSubgraph(dag, sequence)
			sortedSequence, _ := graph.TopologicalSort(subGraph)
			reorderedMidpoints = append(reorderedMidpoints, sortedSequence[len(sortedSequence)/2])
		}
	}

	return sum(reorderedMidpoints)
}

func run(runPart2 bool, input string) any {
	if runPart2 {
		return part2(input)
	}
	return part1(input)

}

func sum(slice []int) int {
	total := 0
	for _, value := range slice {
		total += value
	}
	return total
}
