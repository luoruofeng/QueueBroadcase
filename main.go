package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"sync"
)

type Node struct {
	Id          int
	Address     string
	Port        int
	MessageChan chan Message
	Neighbors   []int
}

func NewNode() Node {
	id := os.Args[1]
	neighbors := make([]int, 0, 2)
	idi, err := strconv.Atoi(id)
	if err != nil {
		panic("Id must be integer!")
	}
	neighbors = append(neighbors, idi-1, idi+1)
	return Node{
		Id:          idi,
		Address:     "localhost",
		Port:        idi + 8000,
		MessageChan: make(chan Message),
		Neighbors:   neighbors,
	}
}

type Message struct {
	From    int //Node Id
	Content string
}

// RPC方法 接受邻居push来的消息
func (n *Node) ReceiveMessageRPC(message Message, reply *bool) error {
	if message.From != n.Id {
		// fmt.Println("RECEIVED MESSAGE FROM NEIGHBOR ID: ", message.From, " CONTENT IS ", message.Content)
		fmt.Println(message.Content)
		n.MessageChan <- message
		*reply = true
	}
	return nil
}

// 向邻居发送消息
func (n *Node) SendMessage(message *Message, from int) error {
	var wg sync.WaitGroup

	//把内容push给所有的邻居节点
	for _, neighbor := range n.Neighbors {
		if from == neighbor {
			continue
		}
		wg.Add(1)
		go func(neighbor int) {
			defer wg.Done()
			//建立连接
			client, err := rpc.Dial("tcp", "localhost:"+strconv.Itoa(neighbor+8000))
			if err != nil {
				// fmt.Printf("Failed to dial neighbor node ID: %d \n", neighbor)
				return
			}
			defer client.Close()

			var reply bool
			//调用RPC方法
			client.Call("Node.ReceiveMessageRPC", message, &reply)
			// if reply {
			// fmt.Println("PUSHED SUCCESSFULLY! NEIGHBOR ID IS ", neighbor)
			// } else {
			// fmt.Println("PUSHED FAILURE! NEIGHBOR ID IS ", neighbor)
			// }
		}(neighbor)
	}
	wg.Wait()
	return nil
}

func main() {
	id, _ := strconv.Atoi(os.Args[1])
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(id+8000))
	if err != nil {
		panic(err)
	}

	//初始化节点
	node := NewNode()

	defer close(node.MessageChan)
	//接受消息
	go func() {
		for message := range node.MessageChan {
			origin := message.From             //记录消息的来源
			message.From = node.Id             //把消息的来源改为自己
			node.SendMessage(&message, origin) //把消息push给所有的邻居节点
		}
		fmt.Println("MESSAGE CHANNEL CLOSED!")
	}()

	//注册RPC服务
	go func() {
		rpc.RegisterName("Node", &node)
		rpc.Accept(ln)
	}()

	//接受用户输入的消息
	for {
		var message *Message = new(Message)
		fmt.Println("enter:")
		_, err := fmt.Scan(&message.Content)
		if err != nil {
			panic(err)
		}
		message.From = node.Id
		//把消息push给所有的邻居节点
		node.SendMessage(message, node.Id)
	}
}
