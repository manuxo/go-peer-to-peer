package main

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
)

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
	src := fmt.Sprintf("%d-%s-%s", block.Index, block.Timestamp.String(), block.Data)
	return base64.StdEncoding.EncodeToString([]byte(src))
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
	fmt.Println(1)
	return block
}

func (blockChain *BlockChain) GetLatesBlock() Block {
	n := len(blockChain.Chain)
	return blockChain.Chain[n-1]
}

func (blockChain *BlockChain) AddBlock(block Block) {
	block.Timestamp = time.Now()
	block.Index = blockChain.GetLatesBlock().Index + 1
	block.PreviousHash = blockChain.GetLatesBlock().Hash
	block.Hash = block.CalculateHash()
	blockChain.Chain = append(blockChain.Chain, block)
}

func (blockChain *BlockChain) IsChainValid() bool {
	n := len(blockChain.Chain)
	for i := 1; i < n; i++ {
		currentBlock := blockChain.Chain[i]
		previousBlock := blockChain.Chain[i-1]
		if currentBlock.Hash != currentBlock.CalculateHash() {
			return false
		}
		if currentBlock.PreviousHash != previousBlock.Hash {
			return false
		}
	}
	return true
}

func CreateBlockChain() BlockChain {
	bc := BlockChain{}
	genesisBlock := bc.CreateGenesisBlock()
	bc.Chain = append(bc.Chain, genesisBlock)
	return bc
}

func main() {
	bc := CreateBlockChain()
	mr1 := MedicalRecord{
		Name:       "MR1",
		Year:       "2019",
		Hospital:   "H1",
		Doctor:     "D1",
		Diagnostic: "DG1",
		Medication: "M1",
		Procedure:  "P1",
	}
	block1 := Block{
		Data: mr1,
	}
	bc.AddBlock(block1)
	mr2 := MedicalRecord{
		Name:       "MR2",
		Year:       "2022",
		Hospital:   "H2",
		Doctor:     "D2",
		Diagnostic: "DG2",
		Medication: "M2",
		Procedure:  "P2",
	}
	block2 := Block{
		Data: mr2,
	}
	bc.AddBlock(block2)
	//src, _ := json.Marshal(bc)
	//fmt.Println(string(src))
	fmt.Printf("Is blockchain valid? %s\n", strconv.FormatBool(bc.IsChainValid()))
	bc.Chain[1].Data.Doctor = "Manuel"
	fmt.Printf("Is blockchain valid? %s\n", strconv.FormatBool(bc.IsChainValid()))
}
