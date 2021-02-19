package allnodes

import (
	"fmt"
	"sync"
)

//constants are declared for scalibility of the datastructure should it be employed in other projects albeit certain numbers are still hardcoded.
const (
	order       = 4
	minChildren = order / 2
	maxChildren = order
	minKeys     = 1
	maxKeys     = order/2 + 1
)

//key is an abritrary value to prevent duplicate UIDs it is initialized at 0.
var (
	key   int
	mutex sync.Mutex
	wg    sync.WaitGroup
)

//BPlusTree stores all the information in a B tree
type BPlusTree struct {
	root  *Node
	mutex sync.Mutex
}

//Node stores the Keys and Pointers to the Next Node
type Node struct {
	parent      **Node
	Pointers    []interface{}
	Keys        []int
	NumKeys     int
	NumPointers int
	isLeaf      bool
	Next        **Node
}

//NewBPlusTree Initializes an empty instance of the tree
func NewBPlusTree() *BPlusTree {
	return &BPlusTree{}
}

//SeedData since there is no db
func (b *BPlusTree) SeedData() {
	b.InsertTree(CreateLocation("Woodlands", "Inaccessible Place in the North", 5))
	b.InsertTree(CreateLocation("Pasir Ris", "Inaccessible Place in the east", 2))
	b.InsertTree(CreateLocation("Chinatown", "Chinese Cultural Place", 3))
	b.InsertTree(CreateLocation("Orchard", "Singapore's Shopping District", 6))
	b.InsertTree(CreateLocation("Serangoon", "Neighborhood", 8))
	b.InsertTree(CreateLocation("NTU", "University in the west", 10))
	b.InsertTree(CreateLocation("NUS", "Singapore's main university", 1))
	b.InsertTree(CreateLocation("SMU", "Business School", 0))
	b.InsertTree(CreateLocation("CBD", "Central Business District", 15))
	b.InsertTree(CreateLocation("Toa Payoh", "Another Neighborhood", 16))
	b.InsertTree(CreateLocation("Zoo", "Zoo", 12))
	b.InsertTree(CreateLocation("Mount Faber", "Sight seeing", 13))
	b.InsertTree(CreateLocation("MBS", "Gambling spot", 13))
}

//Search searches through the tree based on the UID and returns the pointer to the data if it finds it
func (b *BPlusTree) Search(dataKey int) (*Location, error) {
	leafNode := b.findLeafNode(dataKey, &b.root)
	for i := 0; i < (*leafNode).NumKeys; i++ {
		if (*leafNode).Keys[i] == dataKey {
			return (*leafNode).Pointers[i].(*Location), nil
		}
	}
	return nil, fmt.Errorf("Location ID not found")
}

//InsertTree finds the respective position in the tree and inserts the node
func (b *BPlusTree) InsertTree(location *Location) {
	if b.root == nil {
		b.firstNode(location)
		return
	}
	key++
	leafNode := *b.findLeafNode(key, &b.root)
	leafNode.Keys[leafNode.NumKeys] = key
	leafNode.Pointers[leafNode.NumPointers] = location
	leafNode.NumKeys++
	leafNode.NumPointers++
	return
}

//findLeafNode always checks if the node is full before splitting. This will always ensure that there is room for insertion at a cost of additional space.
func (b *BPlusTree) findLeafNode(insertKey int, node **Node) **Node {
	if (*node).NumKeys == maxKeys {
		return b.splitNode(insertKey, &(*node))
	}

	if (*node).isLeaf {
		return node
	}
	for i := 0; i < (*node).NumKeys; i++ {
		if insertKey < (*node).Keys[i] {
			nodeV := (*node).Pointers[i].(*Node)
			return b.findLeafNode(insertKey, &nodeV)
		}
	}
	searchNode := (*node).Pointers[(*node).NumKeys].(*Node)
	return b.findLeafNode(insertKey, &searchNode)
}

//splitNode pushes the key into the parent and recursively calling back findLeafNode.
func (b *BPlusTree) splitNode(insertKey int, node **Node) **Node {
	middleKey := (*node).Keys[maxKeys/2]
	leftKeys := make([]int, maxKeys)
	rightKeys := make([]int, maxKeys)
	leftPointers := make([]interface{}, maxChildren)
	rightPointers := make([]interface{}, maxChildren)
	rightNumPointers := maxKeys/2 + 1
	leftNumPointers := maxKeys/2 + 1
	leftNumKeys := maxKeys / 2
	rightNumKeys := maxKeys / 2
	parent := &(*node)
	if (*node) != b.root {
		parent = &(*(*node).parent)
	}
	if (*node).isLeaf {
		for i := 0; i < maxKeys/2; i++ {
			leftPointers[i] = (*node).Pointers[i]

			leftKeys[i] = (*node).Keys[i]
			(*node).NumPointers--
			(*node).NumKeys--
		}
		var j int
		for i := maxKeys / 2; i < maxKeys; i++ {
			rightPointers[j] = (*node).Pointers[i]
			rightKeys[j] = (*node).Keys[i]
			(*node).NumPointers--
			(*node).NumKeys--
			j++
		}
		(*node).NumPointers++
	} else {
		for i := 0; i < maxChildren/2; i++ {
			leftPointers[i] = (*node).Pointers[i]
			if i < maxKeys/2 {
				leftKeys[i] = (*node).Keys[i]
				(*node).NumKeys--
			}
			(*node).NumPointers--

		}
		var j int
		for i := maxKeys/2 + 1; i < maxChildren; i++ {
			rightPointers[j] = (*node).Pointers[i]
			if i < maxKeys {
				rightKeys[j] = (*node).Keys[i]
				(*node).NumKeys--
			}
			(*node).NumPointers--
			j++
		}
		(*node).NumPointers++
	}

	var rightNextNode **Node
	var leftNextNode **Node
	if (*node).isLeaf {
		leftNumPointers = maxChildren/2 - 1
		rightNumPointers = maxChildren / 2
		leftNumKeys = 1
		rightNumKeys = 2
		if (*node).Next == nil {
			rightNextNode = nil
		} else {
			rightNextNode = &(*(*node).Next)
		}

	}

	rightPointer := &Node{
		parent:      parent,
		Pointers:    rightPointers,
		Keys:        rightKeys,
		NumKeys:     rightNumKeys,
		NumPointers: rightNumPointers,
		isLeaf:      (*node).isLeaf,
		Next:        rightNextNode,
	}
	if (*node).isLeaf {
		leftNextNode = &rightPointer

	}
	var leftPointer *Node
	if (*node) == b.root {
		leftPointer = &Node{
			parent:      parent,
			Pointers:    leftPointers,
			Keys:        leftKeys,
			NumKeys:     leftNumKeys,
			NumPointers: leftNumPointers,
			isLeaf:      (*node).isLeaf,
			Next:        leftNextNode,
		}
	} else {
		(*node).parent = parent
		(*node).Pointers = leftPointers
		(*node).Keys = leftKeys
		(*node).NumKeys = leftNumKeys
		(*node).NumPointers = leftNumPointers
		(*node).Next = leftNextNode
		leftPointer = *node
	}

	if !(*node).isLeaf {
		for i := 0; i < (leftPointer).NumPointers; i++ {
			leftPointer.Pointers[i].(*Node).parent = &leftPointer
		}
		for i := 0; i < rightPointer.NumPointers; i++ {
			rightPointer.Pointers[i].(*Node).parent = &rightPointer
		}
	}
	if (*node) == b.root {
		(*node).Keys = make([]int, maxKeys)
		(*node).Pointers = make([]interface{}, maxChildren)
		(*node).Keys[0] = middleKey
		(*node).Pointers[0] = leftPointer
		(*node).Pointers[1] = rightPointer
		(*node).isLeaf = false
		(*node).NumKeys = 1
		(*node).NumPointers = 2
		return b.findLeafNode(insertKey, &(*node))
	}

	b.insertIntoParent(middleKey, *(*node).parent, &rightPointer, &leftPointer)
	return b.findLeafNode(insertKey, (*node).parent)

}

//insertIntoParent inserts the left nad right pointers into the respective positions in the parent node and adjusts the values
func (b *BPlusTree) insertIntoParent(insertKey int, node *Node, right **Node, left **Node) {
	(*node).Keys[(*node).NumKeys] = insertKey
	(*node).NumKeys++
	(*node).Pointers[(*node).NumPointers-1] = (*left)
	(*node).Pointers[(*node).NumPointers] = (*right)
	(*node).NumPointers++
}

//firstNode created to initialize the database.
func (b *BPlusTree) firstNode(location *Location) {
	b.root = &Node{
		parent:      nil,
		Keys:        make([]int, order-1),
		Pointers:    make([]interface{}, order),
		NumKeys:     1,
		NumPointers: 1,
		isLeaf:      true,
		Next:        nil,
	}
	key++
	b.root.Keys[0] = key
	b.root.Pointers[0] = location
}

//Delete finds a leaf node. Should the leaf node contain more than 1 key, deletion is carried out with no issues. If the leafnode is left with only 1 key, it shall attempt to borrow nodes before merging nodes
func (b *BPlusTree) Delete(deletedKey int) (bool, error) {
	node := b.findLeafNode(deletedKey, &b.root)
	var i int
	for i = 0; i < (*node).NumKeys; i++ {
		if (*node).Keys[i] == deletedKey {
			break
		}
	}
	if (*node).Keys[i] != deletedKey {
		return false, fmt.Errorf("Invalid Key")
	}
	indexInParent, err := b.getIndexInParent(node)
	if err != nil {
		return false, err
	}
	if (*node).NumKeys > 1 {
		for j := i; j < (*node).NumKeys-1; j++ {
			(*node).Keys[j] = (*node).Keys[j+1]
			(*node).Pointers[j] = (*node).Pointers[j+1]
		}
		(*node).Keys[(*node).NumKeys-1] = 0
		(*node).Pointers[(*node).NumPointers-1] = 0
		(*node).NumKeys--
		(*node).NumPointers--
		var keyIndexInParent int
		if indexInParent > 0 {
			keyIndexInParent = indexInParent - 1
		}
		if (*(*node).parent).Keys[keyIndexInParent] != deletedKey {
			nodeContainingKey, index, ok := b.searchParents((*node).Keys[keyIndexInParent], (*node).parent)
			if ok {
				deRefNode := (*nodeContainingKey).Pointers[index+1].(*Node)
				(*nodeContainingKey).Keys[index] = b.smallestInSubTree(&deRefNode)
			}
		} else if indexInParent > 0 {
			(*(*node).parent).Keys[keyIndexInParent] = (*node).Keys[0]
		}

		return true, nil
	}
	(*node).Keys[0] = 0
	(*node).Pointers[0] = nil
	(*node).NumKeys--
	(*node).NumPointers--

	var keyIndexInParent int
	if indexInParent > 0 {
		keyIndexInParent = indexInParent - 1
	}
	ok := b.borrowFromLeft(keyIndexInParent, (*node).parent, deletedKey)
	if !ok {
		ok = b.borrowFromRight(indexInParent, (*node).parent, deletedKey)
	}
	if !ok {
		ok = b.mergeRightChild(indexInParent, (*node).parent, deletedKey)
		if !ok {
			ok = b.mergeLeftChild(indexInParent, (*node).parent, deletedKey)
		}
	}
	return true, nil

}

//Returns the index of the node in the parent for quick access
func (b *BPlusTree) getIndexInParent(node **Node) (int, error) {
	var parentNode **Node
	if node == &b.root {
		parentNode = node
	} else {
		parentNode = (*node).parent
	}
	for i := 0; i < (*parentNode).NumPointers; i++ {
		if (*parentNode).Pointers[i].(*Node) == *node {
			return i, nil
		}
	}
	return -1, fmt.Errorf("Pointer not found")
}

//mergeRightChild will check if the child nodes are leaves, if they are, they will be merged. Upon merging should the parent node not satisfy the tree conditions, it will recursively call on mergeInternalNode until the tree is balanced
func (b *BPlusTree) mergeRightChild(indexInParent int, node **Node, deletedKey int) bool {
	if indexInParent == (*node).NumPointers-1 {
		return false
	}
	var keyIndexInParent int
	if indexInParent > 0 {
		keyIndexInParent = indexInParent - 1
	}
	removedKey := (*node).Keys[keyIndexInParent]
	leftChild := &(*node).Pointers[indexInParent]
	if (*leftChild).(*Node).isLeaf {
		(*node).Pointers[indexInParent] = *(*leftChild).(*Node).Next
		(*node).NumPointers--
		for i := indexInParent; i < (*node).NumPointers-1; i++ {
			(*node).Pointers[i] = (*node).Pointers[i+1]
		}
		if (*node).NumKeys > 1 {
			for i := keyIndexInParent; i < (*node).NumKeys; i++ {
				(*node).Keys[i] = (*node).Keys[i+1]
			}
			(*node).NumKeys--
			(*node).Keys[(*node).NumKeys] = 0
		}
		(*node).Pointers[(*node).NumPointers] = nil
	}
	if removedKey != deletedKey {
		nodeContainingKey, index, ok := b.searchParents(deletedKey, &(*node))
		if ok {
			deRefNode := (*nodeContainingKey).Pointers[index+1].(*Node)
			(*nodeContainingKey).Keys[index] = b.smallestInSubTree(&deRefNode)
		}
	}
	if (*node).NumPointers == 1 {
		indexInParent, err := b.getIndexInParent(node)
		if err != nil {
			fmt.Println(err)
			return false
		}
		ok := b.borrowFromLeft(indexInParent, (*node).parent, deletedKey)
		if !ok {
			ok = b.borrowFromRight(indexInParent, (*node).parent, deletedKey)
		}
		if !ok {
			ok = b.mergeLeftInternalNode(indexInParent, (*node).parent)
			if !ok {
				ok = b.mergeRightInternalNode(indexInParent, (*node).parent)
			}
		}
	}
	return true

}

//mergeLeftChild will check if the child nodes are leaves, if they are, they will be merg. Upon merging should the parent node not satisfy the tree conditions, it will recursively call on mergeInternalNode until the tree is balanced
func (b *BPlusTree) mergeLeftChild(indexInParent int, node **Node, deletedKey int) bool {
	if indexInParent == 0 {
		return false
	}
	var keyIndexInParent int
	if indexInParent > 0 {
		keyIndexInParent = indexInParent - 1
	}
	leftChild := &(*node).Pointers[indexInParent-1]
	rightChild := (*node).Pointers[indexInParent].(*Node).Next
	removedKey := (*node).Keys[keyIndexInParent]
	if (*leftChild).(*Node).isLeaf {
		fmt.Println((*leftChild))
		(*leftChild).(*Node).Next = rightChild
		(*node).NumPointers--

		for i := indexInParent; i < (*node).NumPointers-1; i++ {
			(*node).Pointers[i] = (*node).Pointers[i+1]
		}
		for i := keyIndexInParent; i < (*node).NumKeys; i++ {
			(*node).Keys[i] = (*node).Keys[i+1]
		}
	}
	if removedKey != deletedKey {
		nodeContainingKey, index, ok := b.searchParents((*node).Keys[keyIndexInParent], &(*node))
		if ok {
			deRefNode := (*nodeContainingKey).Pointers[index+1].(*Node)
			(*nodeContainingKey).Keys[index] = b.smallestInSubTree(&deRefNode)
		}
	}
	if (*node).NumPointers == 1 {
		indexInParent, err := b.getIndexInParent(node)
		if err != nil {
			fmt.Println(err)
			return false
		}
		ok := b.borrowFromLeft(indexInParent, (*node).parent, deletedKey)
		if !ok {
			ok = b.borrowFromRight(indexInParent, (*node).parent, deletedKey)
		}
		if !ok {
			ok = b.mergeLeftInternalNode(indexInParent, (*node).parent)
			if !ok {
				ok = b.mergeRightInternalNode(indexInParent, (*node).parent)
			}
		}
	}
	return true

}

//mergeLeftInternalNode checks if the parent node has more than 1 key, if it does, it brings the key down and combines with the children node and apnds to the new parent node which has 1 key and pointer less.
//In the even the parent node only has 1 key, the children Pointers are brought  and combined with the node to maintain the tree's constant height
func (b *BPlusTree) mergeLeftInternalNode(indexInParent int, node **Node) bool {
	if indexInParent == 0 {
		return b.mergeRightInternalNode(indexInParent, node)
	}
	var keyIndexInParent int
	leftPointerIndex := indexInParent - 1
	if indexInParent > 0 {
		keyIndexInParent = indexInParent - 1
	}
	if (*node).NumKeys == 1 {
		//attach the children key and pointer to the parennode and call on the merge internalnode again
		NextIndexInParent, err := b.getIndexInParent(node)
		if err != nil {
			fmt.Println(err)
			return false
		}
		var leftKey int
		var rightKey int
		if (*node).isLeaf {
			leftKey = (*node).Keys[0]
			rightKey = (*node).Keys[0]
		} else {
			leftKey = (*node).Pointers[0].(*Node).Keys[0]
			rightKey = (*node).Pointers[1].(*Node).Keys[0]
		}

		if rightKey == 0 {
			rightKey = (*node).Keys[0]
		}
		firstPointer := (*node).Pointers[0].(*Node).Pointers[0]
		secondPointer := (*node).Pointers[0].(*Node).Pointers[1]
		thirdPointer := (*node).Pointers[1].(*Node).Pointers[0]
		(*(thirdPointer.(*Node))).parent = (*(firstPointer.(*Node))).parent
		(*node).Keys[0] = leftKey
		(*node).Keys[1] = rightKey
		(*node).Pointers[0] = firstPointer
		(*node).Pointers[1] = secondPointer
		(*node).Pointers[2] = thirdPointer
		(*node).NumKeys++
		(*node).NumPointers++
		return b.mergeLeftInternalNode(NextIndexInParent, (*node).parent)

	}
	//bring the parent key down and attach to the rging node
	//update parentKey and pointer registersfter
	newKey := (*node).Keys[keyIndexInParent]
	newPointer := (*node).Pointers[indexInParent].(*Node)
	for i := keyIndexInParent; i < (*node).NumKeys; i++ {
		(*node).Keys[i] = (*node).Keys[i+1]
	}
	for i := indexInParent; i < (*node).NumPointers; i++ {
		(*node).Pointers[i] = (*node).Pointers[i+1]
	}
	(*node).NumKeys--
	(*node).NumPointers--
	leftNode := (*node).Pointers[leftPointerIndex].(*Node)
	leftNode.Keys[leftNode.NumKeys] = newKey
	newPointer.parent = &leftNode
	leftNode.Pointers[leftNode.NumPointers] = newPointer
	leftNode.NumKeys++
	leftNode.NumPointers++
	if (*node).NumKeys == 0 && (*node) == b.root {
		(*node).Pointers[0].(*Node).parent = nil
		b.root = (*node).Pointers[0].(*Node)
	}
	return true

}

//mergeRightInternalNode checks if the parent node has more than 1 key, if it does, it brings the key down and combines with the children node and appends to the new parent node which has 1 key and pointer less.
//In the even the parent node only has 1 key, the children Pointers are brought uand combined with the node to maintain the tree's constant height
func (b *BPlusTree) mergeRightInternalNode(indexInParent int, node **Node) bool {
	if indexInParent == (*node).NumPointers-1 {
		return b.mergeLeftInternalNode(indexInParent, node)
	}
	var keyIndexInParent int
	rightPointerIndex := indexInParent + 1
	if indexInParent > 0 {
		keyIndexInParent = indexInParent
	}
	if (*node).NumKeys == 1 && (*node) != b.root {
		//attach the children key and pointer to the parentnode and call on the merge internalnode again
		NextIndexInParent, err := b.getIndexInParent(node)
		if err != nil {
			fmt.Println(err)
		}
		var leftKey int
		var rightKey int
		if (*node).isLeaf {
			leftKey = (*node).Keys[0]
			rightKey = (*node).Keys[1]
		} else {
			leftKey = (*node).Keys[0]
			rightKey = (*node).Pointers[1].(*Node).Keys[0]
		}
		if leftKey == 0 {
			leftKey = (*node).Pointers[1].(*Node).Pointers[0].(*Node).Keys[0]
		}

		firstPointer := (*node).Pointers[indexInParent].(*Node).Pointers[0]
		secondPointer := (*node).Pointers[rightPointerIndex].(*Node).Pointers[0]
		thirdPointer := (*node).Pointers[rightPointerIndex].(*Node).Pointers[1]
		(*(firstPointer.(*Node))).parent = (*(thirdPointer.(*Node))).parent
		(*node).Keys[0] = leftKey
		(*node).Keys[1] = rightKey
		(*node).Pointers[0] = firstPointer
		(*node).Pointers[1] = secondPointer
		(*node).Pointers[2] = thirdPointer
		(*node).NumKeys++
		(*node).NumPointers++
		if NextIndexInParent < 0 {
			NextIndexInParent = 1
		}
		return b.mergeRightInternalNode(NextIndexInParent, (*node).parent)

	}
	//bring the parent key down and attach to the merging node
	//update parentKey and pointer registers after

	newKey := (*node).Keys[keyIndexInParent]
	newPointer := (*node).Pointers[indexInParent].(*Node)
	rightNode := (*node).Pointers[rightPointerIndex].(*Node)
	for i := keyIndexInParent; i < (*node).NumKeys; i++ {
		(*node).Keys[i] = (*node).Keys[i+1]
	}
	for i := indexInParent; i < (*node).NumPointers; i++ {
		(*node).Pointers[i] = (*node).Pointers[i+1]
	}
	(*node).NumKeys--
	(*node).NumPointers--
	rightNode.Keys[rightNode.NumKeys] = rightNode.Keys[0]
	rightNode.Keys[0] = newKey
	newPointer.parent = &rightNode
	rightNode.Pointers[2] = rightNode.Pointers[1]
	rightNode.Pointers[1] = rightNode.Pointers[0]
	rightNode.Pointers[0] = newPointer
	rightNode.NumKeys++
	rightNode.NumPointers++
	if (*node).NumKeys == 0 && (*node) == b.root {
		(*node).Pointers[0].(*Node).parent = nil
		b.root = (*node).Pointers[0].(*Node)
	}
	return true
}

//search for the node that contains the key to be deleted if it isn't directly abovehe deleted node and returns the address of said node
func (b *BPlusTree) searchParents(deletedKey int, node **Node) (**Node, int, bool) {
	for i := 0; i < (*node).NumKeys; i++ {
		if (*node).Keys[i] == deletedKey {
			return node, i, true
		}
	}
	if (*node) == b.root {
		return nil, -1, false
	}
	return b.searchParents(deletedKey, (*node).parent)

}

//
func (b *BPlusTree) searchNode(searchKey int, node **Node) **Node {
	var childNode *Node
	if (*node).Pointers[0].(*Node).NumKeys > 2 {
		childNode = (*node).Pointers[0].(*Node)
	} else {
		childNode = (*node).Pointers[1].(*Node)
	}
	keyVal, nodePointer := b.largestNode(&childNode)
	(*node).Keys[1] = (*node).Keys[0]
	(*node).Keys[0] = keyVal
	(*node).Pointers[2] = (*node).Pointers[1]
	(*node).Pointers[1] = nodePointer
	(*node).NumKeys++
	(*node).NumPointers++
	return b.findLeafNode(searchKey, node)
}

//returns the largest node in the subtree
func (b *BPlusTree) largestNode(node **Node) (int, *Node) {
	keyVal := (*node).Keys[(*node).NumKeys-1]
	(*node).Keys[(*node).NumKeys-1] = 0
	var nodePointer *Node

	if (*node).isLeaf {
		KeysArray := make([]int, maxKeys)
		KeysArray[0] = keyVal
		PointersArray := make([]interface{}, maxChildren)
		PointersArray[0] = (*node).Pointers[(*node).NumPointers-1].(*Location)
		nodePointer = &Node{
			parent:      (*node).parent,
			Pointers:    PointersArray,
			Keys:        KeysArray,
			NumKeys:     1,
			NumPointers: 1,
			isLeaf:      (*node).isLeaf,
			Next:        (*(*node).parent).Pointers[1].(**Node),
		}
	} else {
		nodePointer = (*node).Pointers[(*node).NumPointers-1].(*Node)
	}

	(*node).Pointers[(*node).NumPointers-1] = nil
	(*node).NumKeys--
	(*node).NumPointers--
	return keyVal, nodePointer

}

//borrows the Keys and Pointers from the left node given that the left no has more than 1 key/pointer
func (b *BPlusTree) borrowFromLeft(indexInParent int, node **Node, deletedKey int) bool {
	if indexInParent == 0 {
		return false
	}
	if (*node).Pointers[indexInParent-1].(*Node).NumKeys == 1 {
		return false
	}
	leftNode := (*node).Pointers[indexInParent-1].(*Node)
	borrowedKey := leftNode.Keys[leftNode.NumKeys-1]
	borrowingNode := (*node).Pointers[indexInParent].(*Node)
	var borrowedPointer *interface{}
	borrowedPointer = &leftNode.Pointers[leftNode.NumPointers-1]
	keyIndexInParent := indexInParent - 1
	if (*node).isLeaf {
		borrowingNode.Pointers[0] = borrowedPointer
		borrowingNode.Keys[0] = borrowedKey
		(*node).Keys[keyIndexInParent] = borrowedKey
	} else {
		borrowingNode.Pointers[1] = borrowingNode.Pointers[0]
		borrowingNode.Pointers[0] = borrowedPointer
		borrowingNode.Keys[0] = borrowingNode.Pointers[1].(*Node).Keys[0]
		(*borrowingNode.parent).Keys[keyIndexInParent] = borrowedKey
	}
	borrowingNode.NumKeys++
	borrowingNode.NumPointers++
	b.updateKeysInLeftSibling(&leftNode, keyIndexInParent, deletedKey)
	return true

}

//borrows the Keys and Pointers from the right node given that the right ne has more than 1 key/pointer
func (b *BPlusTree) borrowFromRight(indexInParent int, node **Node, deletedKey int) bool {
	if indexInParent == (*node).NumPointers-1 {
		return false
	}
	if (*node).Pointers[indexInParent+1].(*Node).NumKeys == 1 {
		return false
	}
	var borrowedPointer interface{}
	var keyIndexInParent int
	rightNode := (*node).Pointers[indexInParent+1].(*Node)
	borrowingNode := (*node).Pointers[indexInParent].(*Node)
	borrowedKey := rightNode.Keys[0]
	borrowedPointer = rightNode.Pointers[0]

	if indexInParent > 0 {
		keyIndexInParent = indexInParent - 1
	}

	if borrowingNode.isLeaf {
		borrowingNode.Pointers[0] = borrowedPointer
		borrowingNode.Keys[0] = borrowedKey

	} else {
		borrowingNode.Pointers[1] = borrowedPointer
		borrowingNode.Keys[0] = (*borrowingNode.parent).Keys[keyIndexInParent]
		(*borrowingNode.parent).Keys[keyIndexInParent] = borrowedKey
	}
	borrowingNode.NumKeys++
	borrowingNode.NumPointers++
	b.updateKeysInRightSibling(&rightNode, keyIndexInParent, deletedKey)
	return true

}

//update the sibiling node after borrowing the key and data
func (b *BPlusTree) updateKeysInRightSibling(node **Node, indexInParent int, deletedKey int) {

	for i := 0; i < (*node).NumKeys-1; i++ {
		(*node).Keys[i] = (*node).Keys[i+1]
	}
	for j := 0; j < (*node).NumPointers-1; j++ {
		(*node).Pointers[j] = (*node).Pointers[j+1]
	}
	(*node).Keys[(*node).NumKeys-1] = 0
	(*node).Pointers[(*node).NumPointers-1] = nil
	(*node).NumKeys--
	(*node).NumPointers--
	newKey := (*node).Keys[0]
	b.updateKeyInParent((*node).parent, indexInParent, newKey, deletedKey)
}

//update the sibiling node after borrowing the key and data
func (b *BPlusTree) updateKeysInLeftSibling(node **Node, indexInParent int, deletedKey int) {
	newKey := (*node).Keys[(*node).NumKeys-1]
	(*node).Keys[(*node).NumKeys-1] = 0
	(*node).Pointers[(*node).NumPointers-1] = nil
	(*node).NumKeys--
	(*node).NumPointers--
	b.updateKeyInParent((*node).parent, indexInParent, newKey, deletedKey)

}

// Updates the key in the parent node after borrowing the key from the sibling no
func (b *BPlusTree) updateKeyInParent(node **Node, indexInParent int, key int, deletedKey int) {
	for i := indexInParent; i < (*node).NumKeys-1; i++ {
		(*node).Keys[i] = (*node).Keys[i+1]
	}
	removedKey := (*node).Keys[(*node).NumKeys-1]
	(*node).Keys[(*node).NumKeys-1] = key
	if removedKey != deletedKey {
		nodeContainingKey, index, ok := b.searchParents(deletedKey, &(*node))
		if ok {
			deRefNode := (*nodeContainingKey).Pointers[index+1].(*Node)
			(*nodeContainingKey).Keys[index] = b.smallestInSubTree(&deRefNode)
		}
	}
}

//returns the smallest in the subtree to update the nodef the deleted key doesn't correspond to the deleted value
func (b *BPlusTree) smallestInSubTree(node **Node) int {
	if (*node).isLeaf {
		return (*node).Keys[0]
	}
	NextNode := (*node).Pointers[0].(*Node)
	return b.smallestInSubTree(&NextNode)

}

//FirstNode recursively calls on smallest until the left most node is found and returned for ease of traversing the data
func (b *BPlusTree) FirstNode() **Node {
	return b.smallest(&b.root)
}
func (b *BPlusTree) smallest(node **Node) **Node {
	if (*node).isLeaf {
		return node
	}
	NextNode := (*node).Pointers[0].(*Node)
	return b.smallest(&NextNode)

}

//PrintLeaves prints all nodes in t tree helper function for debugging and visualizaing the tree
func (b *BPlusTree) PrintLeaves() {
	fmt.Println("START PRINTING LEAVES ========================================================================")
	fmt.Println(&b.root, "ROOT ADDRESS")
	fmt.Println(b.root.Pointers[0])
	fmt.Println(b.root.Pointers[1])
	if key >= 8 {
		// fmt.Println(*(b.root.Pointers[0].(*Node).Pointers[0].(*Node).Pointers[0.(*Node).Next))
		// fmt.Println(b.root.Pointers[0].(*Node).Pointers[0].(*Node).Pointers[1].*Node))
		// fmt.Println(&*(b.root.Pointers[0].(*Node).Pointers[0].(*Node).Pointers[].(*Node)))
		// fmt.Println(b.root.Pointers[0].(*Node).Pointers[0])
		// fmt.Println(b.root.Pointers[0].(*Node).Pointers[1])
		// fmt.Println(b.root.Pointers[0].(*Node).Pointers[0].(*Node).Pointers[0])
		// fmt.Println(b.root.Pointers[0].(*Node).Pointers[0].(*Node).Pointers[1])
		// fmt.Println(b.root.Pointers[0].(*Node).Pointers[1].(*Node).Pointers[0])
		// fmt.Println(b.root.Pointers[0].(*Node).Pointers[1].(*Node).Pointers[1])
		// fmt.Println(b.root.Pointers[0].(*Node).Pointers[2]
		// fmt.Println(b.root.Pointers[1].(*Node).Pointers[0].(*Node).Pointers[0])
		// fmt.Println(b.root.Pointers[1].(*Node).Pointers[0].(*Node).Pointers[1])
		// fmt.Println(b.root.Pointers[1].(*Node).Pointers[1].(*Node).Pointers[1])
		// fmt.Println(b.root.Pointers[1].(*Node).Pointers[1].(*Node).Pointers[1])
		// fmt.Println(b.root.Pointers[1].(*Node).Pointers[2].(*Node).Pointers[0])
		// fmt.Println(b.root.Pointers[1].(*Node).Pointers[2].(*Node).Pointers[1])
		// fmt.Println(b.root.Pointers[1].(*Node).Pointers[2].(*Node).Pointers[0])
		// fmt.Println(b.root.Pointers[1].(*Node).Pointers[2].(*Node).Pointers[1])
		// fmt.Println(b.root.Pointers[1].(*Node).Pointers[2].(*Node).Pointers[2])
		// fmt.Println(b.root.Pointers[2].(*Node).Pointers[1].(*Node).Pointers[1])
		// fmt.Println(b.root.Pointers[2].(*Node).Pointers[1].(*Node).Pointers[2])
	}
	// fmt.Println(b.root.Pointers[0].(*Node).Pointers[1])
	// fmt.Println(b.root.Pointers[1].(*Node).Pointers[0])
	// fmt.Println(b.root.Pointers[1].(*Node).Pointers[1])
	fmt.Println("END PRINTING LEAVES ===========================================================================")
}

//PrintAll prints the entire tree from the root
func (b *BPlusTree) PrintAll() {
	b.printTree(&b.root)
}

func (b *BPlusTree) printTree(node **Node) {
	fmt.Println(*node)
	if !(*node).isLeaf {
		for i := 0; i < (*node).NumPointers; i++ {
			d := (*node).Pointers[i].(*Node)
			b.printTree(&d)
		}
	}

}
