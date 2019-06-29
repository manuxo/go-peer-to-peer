package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
)

var HOSTS []string
var LOCALHOST string = "25.16.224.61:5000"
var PUBLICHOST string

type MessageType int32

const (
	NEWHOST  MessageType = 0
	ADDHOST  MessageType = 1
	PROTOCOL             = "tcp"
)

type RequestBody struct {
	Message     string
	MessageType MessageType
}

type MedicalRecord struct {
	Name       string
	Year       string
	Hospital   string
	Doctor     string
	Diagnostic string
	Medication string
	Procedure  string
}

type Block struct {
	Index        int
	Timestamp    time.Time
	Data         MedicalRecord
	PreviousHash string
	Hash         string
}

func (block *Block) CalculateHash() string {
	src, _ := json.Marshal(block)
	return base64.StdEncoding.EncodeToString(src)
}

type BlockChain struct {
	Chain []Block
}

func (blockChain *BlockChain) CreateGenesisBlock() Block {
	block := Block{
		Index:        0,
		Timestamp:    time.Now(),
		Data:         MedicalRecord{},
		PreviousHash: "0",
	}
	block.Hash = block.CalculateHash()
	return block
}

func (blockChain *BlockChain) GetLatesBlock() Block {
	n := len(blockChain.Chain)
	return blockChain.Chain[n-1]
}

func (blockChain *BlockChain) AddBlock(block Block) {
	block.PreviousHash = blockChain.GetLatesBlock().Hash
	block.Hash = block.CalculateHash()
	blockChain.Chain = append(blockChain.Chain, block)
}

func GetMessage(conn net.Conn) string {
	reader := bufio.NewReader(conn)
	data, _ := reader.ReadString('\n')
	return strings.TrimSpace(data)
}

func SendMessage(toHost string, message string) string {
	conn, _ := net.Dial(PROTOCOL, toHost)
	defer conn.Close()
	fmt.Fprintln(conn, message)
	return GetMessage(conn)
}

func SendMessageNoReply(toHost string, message string) {
	conn, _ := net.Dial(PROTOCOL, toHost)
	defer conn.Close()
	fmt.Fprintln(conn, message)
}

func RemoveHost(index int, hosts []string) []string {
	n := len(hosts)
	hosts[index] = hosts[n-1]
	hosts[n-1] = ""
	return hosts[:n-1]
}

func RemoveHostByValue(ip string, hosts []string) []string {
	for index, host := range hosts {
		if host == ip {
			return RemoveHost(index, hosts)
		}
	}
	return hosts
}

func Broadcast(newHost string) {
	for _, host := range HOSTS {
		data := append(HOSTS, newHost, PUBLICHOST)
		data = RemoveHostByValue(host, data)
		requestBroadcast := RequestBody{
			Message:     strings.Join(data, ","),
			MessageType: ADDHOST,
		}
		broadcastMessage, _ := json.Marshal(requestBroadcast)
		SendMessageNoReply(host, string(broadcastMessage))
	}
}

func IpServer(end chan<- int) {
	ln, _ := net.Listen(PROTOCOL, LOCALHOST)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		defer conn.Close()
		request := RequestBody{}
		data := GetMessage(conn)
		_ = json.Unmarshal([]byte(data), &request)
		if request.MessageType == NEWHOST {
			fmt.Printf("NEWHOST: %s\n", request.Message)
			message := strings.Join(append(HOSTS, PUBLICHOST), ",")
			requestClient := RequestBody{
				Message:     message,
				MessageType: ADDHOST,
			}
			clientMessage, _ := json.Marshal(requestClient)
			SendMessageNoReply(request.Message, string(clientMessage))
			Broadcast(request.Message)
			HOSTS = append(HOSTS, request.Message)
		} else if request.MessageType == ADDHOST {
			fmt.Printf("ADDHOST: %s\n", request.Message)
			HOSTS = strings.Split(request.Message, ",")
		}
		fmt.Println(HOSTS)
	}
	end <- 0
}

func main() {
	var dest string
	end := make(chan int)
	fmt.Print("Ingresa tu host: ")
	fmt.Scanf("%s\n", &PUBLICHOST)
	fmt.Print("Ingresa host destino(Vacío para ser el primer nodo): ")
	fmt.Scanf("%s\n", &dest)
	go IpServer(end)
	if dest != "" {
		requestBody := &RequestBody{
			Message:     PUBLICHOST,
			MessageType: NEWHOST,
		}
		requestMessage, _ := json.Marshal(requestBody)
		SendMessageNoReply(dest, string(requestMessage))
	}
	<-end
}