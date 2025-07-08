package utils

type Node[T comparable] struct {
	data T
	next *Node[T]
}

func NewNode[T comparable](data T, next *Node[T]) *Node[T] {
	return &Node[T]{data, next}
}

type NodeList[T comparable] struct {
	head *Node[T]
	tail *Node[T]
	size int
}

func NewNodeList[T comparable]() *NodeList[T] {
	return &NodeList[T]{
		head: nil,
		tail: nil,
		size: 0,
	}
}

func (n *NodeList[T]) Add(data T) {
	newNode := NewNode(data, nil)
	if n.head == nil {
		n.head = newNode
		n.tail = newNode
	} else {
		n.tail.next = newNode
		n.tail = newNode
	}
	n.size++
}

func (n *NodeList[T]) Pop() *Node[T] {
	if n.head == nil {
		return nil
	}

	if n.head.next == nil {
		n.head = nil
		n.size--
		return nil
	}

	var curr *Node[T]
	var prev *Node[T]
	curr, prev = n.head, nil

	for {
		if curr == nil {
			break
		}

		if curr.next != nil {
			prev = curr
			curr = curr.next
		} else {
			prev.next = nil
			n.size--
			break
		}
	}
	return curr
}

func (n *NodeList[T]) Peek() *Node[T] {
	return n.head
}

func (n *NodeList[T]) Size() int {
	return n.size
}
