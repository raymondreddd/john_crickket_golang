package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

/*
COmperession tool

1. Create hashmap of frequencies of each char/rune
2. Use heap package of std lib to create min priority queue of a huffman tree
*/

// Huffman Node
type HuffmanNode struct {
	char      rune
	frequency int
	left      *HuffmanNode
	right     *HuffmanNode
}

// Priority queue for Huffman Nodes
type PriorityQueue []*HuffmanNode

func (pq PriorityQueue) Len() int { return len(pq) }

// since we want the order to be in ascending.
// Pop first 2 join them, and continue
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].frequency < pq[j].frequency
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	node := x.(*HuffmanNode)
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() interface{} {
	// creating temporary pq, and replacing the pointer with temporary nodes except last
	old := *pq
	n := len(old)
	node := old[n-1]
	*pq = old[0 : n-1]
	return node
}

func main() {
	file := readFile("test.txt")

	// Step 1 - frequency map for each rune or char
	freq_map, err := genFrequencyMap(file)
	check(err)

	// Step 2 - huffman tree (or minimum priority queue based on freq) out of the map
	tree := buildHuffmanTree(freq_map)
	printTree(tree, 0)

	// Step 3 -  Prefix code table for each rune the respective prefixed binary code
	code_table := make(map[rune]string)
	generateCodes(tree, "", code_table)

	fmt.Println("Generated Huffman Codes:")
	for char, code := range code_table {
		fmt.Printf("%q: %s\n", char, code)
	}

}

// using dfs to print the tree
func printTree(node *HuffmanNode, depth int) {
	if node == nil {
		return
	}

	if node.left == nil && node.right == nil {
		fmt.Printf("Node character: %c at depth: %d with frequency: %d \n", node.char, depth, node.frequency)
		return
	}

	printTree(node.left, depth+1)
	printTree(node.right, depth+1)
}

func check(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func readFile(file_name string) *os.File {
	file, err := os.Open(file_name)
	check(err)
	return file
}

func genFrequencyMap(file *os.File) (map[rune]int, error) {
	reader := bufio.NewReader(file)
	defer file.Close()

	freq_map := make(map[rune]int)

	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Println("Error reading file:", err)
			return nil, err
		}
		freq_map[char]++
	}

	keys := make([]rune, 0, len(freq_map))

	// ch is character and frq is frequency
	for ch, frq := range freq_map {
		fmt.Printf("Character: %c %d \n", ch, frq)
		keys = append(keys, ch)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return freq_map[keys[i]] < freq_map[keys[j]]
	})

	fmt.Println(keys)
	return freq_map, nil
}

func buildHuffmanTree(frequencies map[rune]int) *HuffmanNode {

	// here I'm initializing the priority queue data structure as an empty slice.
	pq := make(PriorityQueue, 0)

	// Initialize the heap with the priority queue.
	heap.Init(&pq)

	for char, freq := range frequencies {
		heap.Push(&pq, &HuffmanNode{char: char, frequency: freq})
	}

	// As long as there are 2 nodes left in heap, keep doing below
	for pq.Len() > 1 {
		// left and right are first two (with minimum frequencies) nodes in pq
		// remove them,
		left := heap.Pop(&pq).(*HuffmanNode)
		right := heap.Pop(&pq).(*HuffmanNode)

		// merge them and them again to the heap, where it will reorganized in it
		merged := &HuffmanNode{
			frequency: left.frequency + right.frequency,
			left:      left,
			right:     right,
		}

		// we merge smaller frequencies one because they will keep getting
		// pushed down as tree develops
		heap.Push(&pq, merged)
	}

	return heap.Pop(&pq).(*HuffmanNode)
}

/*
Function to genretate prefix table, where rune - binary output of that rune (or character, pick your vice)
*/
func generateCodes(node *HuffmanNode, prefix string, code_table map[rune]string) {
	if node == nil {
		return
	}

	// If its leaf node, then its a rune, assign the rune to the prefix in code_table
	if node.left == nil && node.right == nil {
		code_table[node.char] = prefix
	}

	// left edge will keep adding 0 and vice versa
	generateCodes(node.left, prefix+"0", code_table)
	generateCodes(node.right, prefix+"1", code_table)
}
