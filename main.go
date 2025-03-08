package main

import "fmt"

// func main() {
// 	fmt.Printf("Hello, world!")
// }

type Node struct {
	value int
	next  *Node
}

type LinkedList struct {
	head *Node
}

func (ll *LinkedList) Insert(value int) {
	newNode := &Node{value: value}
	if ll.head == nil {
		ll.head = newNode
	} else {
		current := ll.head
		for current.next != nil {
			current = current.next
		}
		current.next = newNode
	}
}

func (ll *LinkedList) Display() {
	current := ll.head
	for current != nil {
		fmt.Printf("%d -> ", current.value)
		if current.next == nil {
			fmt.Println("Omg it is a NEO!")
		}
		current = current.next
	}
	fmt.Println("nil")
}

func main() {
	ll := LinkedList{}
	ll.Insert(1)
	ll.Insert(2)
	ll.Insert(3)
	ll.Display()
}
