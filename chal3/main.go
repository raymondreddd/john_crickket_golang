package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

/*
COmperession tool

1. Create hashmap of frequencies of each char/rune
2. Use heap package of std lib to create min priority queue of a huffman tree
*/

// Huffman Node
type HuffmanNode struct {
	Char      rune
	Frequency int
	left      *HuffmanNode
	right     *HuffmanNode
}

// Priority queue for Huffman Nodes
type PriorityQueue []*HuffmanNode

func (pq PriorityQueue) Len() int { return len(pq) }

// since we want the order to be in ascending.
// Pop first 2 join them, and continue
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Frequency < pq[j].Frequency
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
	defer file.Close()

	// Step 1 - frequency map for each rune or char
	freq_map, err := genFrequencyMap(file)
	check(err)

	// Step 2 - huffman tree (or minimum priority queue based on freq) out of the map
	tree := buildHuffmanTree(freq_map)
	printTree(tree, 0)

	// Step 3 -  Prefix code table for each rune the respective prefixed binary code
	prefix_code_table := make(map[rune]string)
	generatePrefixTableMap(tree, "", prefix_code_table)

	fmt.Println("Generated Huffman Codes:")
	for char, code := range prefix_code_table {
		fmt.Printf("%q: %s\n", char, code)
	}

	// Step 4a - get all text from file
	text := getTextFromFile(file)

	// step 4b - main encode
	encoded_text := encodeText(text, prefix_code_table)

	if err := writeEncodedFile("encoded.huff", tree, encoded_text); err != nil {
		log.Fatal(err)
	}

}

func getTextFromFile(file *os.File) string {

	var text string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		text += line
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file in getTextFromFile:", err)
	}

	return text
}

// using dfs to print the tree
func printTree(node *HuffmanNode, depth int) {
	if node == nil {
		return
	}

	if node.left == nil && node.right == nil {
		fmt.Printf("Node character: %c at depth: %d with frequency: %d \n", node.Char, depth, node.Frequency)
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

	// ch is character and frq is frequency
	for ch, frq := range freq_map {
		fmt.Printf("Character: %c %d \n", ch, frq)
	}
	return freq_map, nil
}

func buildHuffmanTree(frequencies map[rune]int) *HuffmanNode {

	// here I'm initializing the priority queue data structure as an empty slice.
	pq := make(PriorityQueue, 0)

	// Initialize the heap with the priority queue.
	heap.Init(&pq)

	for char, freq := range frequencies {
		heap.Push(&pq, &HuffmanNode{Char: char, Frequency: freq})
	}

	// As long as there are 2 nodes left in heap, keep doing below
	for pq.Len() > 1 {
		// left and right are first two (with minimum frequencies) nodes in pq
		// remove them,
		left := heap.Pop(&pq).(*HuffmanNode)
		right := heap.Pop(&pq).(*HuffmanNode)

		// merge them and them again to the heap, where it will reorganized in it
		merged := &HuffmanNode{
			Frequency: left.Frequency + right.Frequency,
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
func generatePrefixTableMap(node *HuffmanNode, prefix string, code_table map[rune]string) {
	if node == nil {
		return
	}

	// If its leaf node, then its a rune, assign the rune to the prefix in code_table
	if node.left == nil && node.right == nil {
		code_table[node.Char] = prefix
	}

	// left edge will keep adding 0 and vice versa
	generatePrefixTableMap(node.left, prefix+"0", code_table)
	generatePrefixTableMap(node.right, prefix+"1", code_table)
}

func encodeText(text string, prefix_code_table map[rune]string) string {
	var encoded_text string

	//extract chars from string and use the map to get the ecoded bit value
	for _, char := range text {
		encoded_text += prefix_code_table[char]
	}
	return encoded_text
}

func writeEncodedFile(filename string, tree *HuffmanNode, encoded_text string) error {
	// To temporarily hold the binary data of the Huffman tree before it is written to a file.
	// laymans term: packing all the tiny blocks (bytes) created by your gob machine.
	var buf bytes.Buffer

	// creating machine(complex data to simple bits) to encode a large text into stream of bytes
	encoder := gob.NewEncoder(&buf)

	// using the machine
	if err := encoder.Encode(tree); err != nil {
		return err
	}

	// Layman's term: I’m taking everything from the box (buf.Bytes()) and putting it into another box called encoded_data
	encoded_data := buf.Bytes()

	// append bytes converted `encoded_text` (original text mapped to binary using prefix_code_table) to `encoded_data` (or the huffman tree we created)
	encoded_data = append(encoded_data, []byte(encoded_text)...)

	// Open the file for writing
	/*
		file is opened with the os.O_WRONLY (write-only) flag,
		the os.O_CREATE flag (which creates the file if it doesn’t exist),
		and the os.O_TRUNC flag (which truncates the file if it already exists).

		0644 is permissions -  (readable and writable by the owner, readable by others).
	*/
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	check(err)
	defer file.Close()

	// Write the encoded data to the file
	_, err = file.Write(encoded_data)
	return err
}
