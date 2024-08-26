package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"encoding/binary"
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

	writeEncodedFile("encoded.huff", tree, encoded_text)

	//step 5: decode the "encoded.huff" file

	encodedFile, err := os.ReadFile("output.huff")
	if err != nil {
		log.Fatal(err)
	}

	tree, err := decodeHuffmanTree(encodedFile[:treeSize])
	if err != nil {
		log.Fatal(err)
	}

	encodedText := string(encodedFile[treeSize:])
	decodedText := decodeText(encodedText, tree)
	fmt.Println("Decoded Text:", decodedText)
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
Function to generate prefix table, where rune - binary output of that rune (or character, pick your vice)
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

// Function to write and encode tree data structure along with encoded text recieved to a file
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

	// Layman's term: I’m taking everything from the box (buf.Bytes()) and putting it into another box called encoded_tree
	encoded_tree := buf.Bytes()

	tree_size := int32(len(encoded_tree))

	// append bytes converted `encoded_text` (original text mapped to binary using prefix_code_table) to `encoded_tree` (or the huffman tree we created)
	encoded_tree = append(encoded_tree, []byte(encoded_text)...)

	// Create a buffer for the header and encoded data, or empty container for blocks
	var out_buf bytes.Buffer

	/*
		Write the size of the Huffman Tree as a 4-byte integer

		&out_buf: You’re telling the tool where to write the data, which is into our container (out_buf).

		binary.BigEndian: This specifies how to format the data. It’s like deciding how you want to write numbers so that everyone reads them the same way.

		tree_size: This is the size of our Huffman Tree (think of it as how big your tree is). You’re writing this size as a 4-byte number into your container.
	*/
	if err := binary.Write(&out_buf, binary.BigEndian, tree_size); err != nil {
		return err
	}

	// Write the a packed-up version of Huffman Tree
	out_buf.Write(encoded_tree)

	// Write the encoded text
	out_buf.WriteString(string(encoded_text))

	/*
		file is opened with the os.O_WRONLY (write-only) flag,
		the os.O_CREATE flag (which creates the file if it doesn’t exist),
		and the os.O_TRUNC flag (which truncates the file if it already exists).

		0644 is permissions -  (readable and writable by the owner, readable by others).
	*/
	return os.WriteFile(filename, out_buf.Bytes(), 0644)
}
