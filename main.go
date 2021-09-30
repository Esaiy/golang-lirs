package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/secnot/orderedmap"
)

var (
	// HIRSize      int = 1
	hit, miss        = 0, 0
	cacheSize    int = 500000
	LIRSize      int = cacheSize * 99 / 100
	HIRSize      int = cacheSize / 100
	orderedStack     = orderedmap.NewOrderedMap()
	orderedList      = orderedmap.NewOrderedMap()
	LIR              = make(map[int]int)
	HIR              = make(map[int]int)
)

func main() {
	var (
		start    = time.Now()
		filePath = "./data/web1_input.txt"
		// filePath = "./data/test_input.txt"
	)

	fmt.Println(LIRSize, HIRSize)

	var file, err = os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lirs(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var duration = time.Since(start)

	fmt.Println("List size :", orderedList.Len())
	fmt.Println("Stack size :", orderedStack.Len())
	fmt.Printf("HIT : %v\nMISS : %v\nHIT RATE : %v\n", hit, miss, float32(hit)/float32((hit+miss)))
	fmt.Printf("Duration : %v\n", duration.Seconds())
}

func lirs(line string) {
	var block = (strings.Split(line, ","))[0]
	// fmt.Printf("-------\nCurrent key : %v\n", block)
	// printStack()
	// printList()
	// printCache()

	var blockNum, err = strconv.Atoi(block)
	if err != nil {
		log.Fatal(err)
	}

	if len(LIR) < LIRSize {
		miss += 1
		// fmt.Println("---miss---")
		addToStack(blockNum)
		makeLIR(blockNum)
		return
	}

	if _, ok := LIR[blockNum]; ok {
		hit += 1
		// fmt.Println("---HIT---")
		key, _, _ := orderedStack.GetFirst()
		keyInInt, _ := strconv.Atoi(key.(string))
		if keyInInt == blockNum {
			stackPrunning(false)
		}
		addToStack(blockNum)
		return
	}

	if _, ok := orderedList.Get(strconv.Itoa(blockNum)); ok {
		hit += 1
		// fmt.Println("---hit---")
		if _, ok := orderedStack.Get(strconv.Itoa(blockNum)); ok {
			makeLIR(blockNum)
			removeFromList(blockNum)
			stackPrunning(true)
		} else {
			addToList(blockNum)
		}
		addToStack(blockNum)
		return
	} else {
		miss += 1
		// fmt.Println("---MISS---")
		if orderedList.Len() == HIRSize {
			key, _, _ := orderedList.PopLast()
			keyInInt, _ := strconv.Atoi(key.(string))
			removeFromList(keyInInt)
		}

		if _, ok := orderedStack.Get(strconv.Itoa(blockNum)); ok {
			makeLIR(blockNum)
			stackPrunning(true)
		} else {
			makeHIR(blockNum)
			addToList(blockNum)
		}
		addToStack(blockNum)
	}

}

func addToStack(block int) {
	key := strconv.Itoa(block)
	if _, ok := orderedStack.Get(key); ok {
		orderedStack.MoveLast(key)
		return
	}
	orderedStack.Set(key, 1)
}

func addToList(block int) {
	key := strconv.Itoa(block)
	if orderedList.Len() == HIRSize {
		orderedList.PopLast()
	}
	orderedList.Set(key, 1)
	orderedList.MoveFirst(key)
}

func removeFromList(block int) {
	key := strconv.Itoa(block)
	orderedList.Delete(key)
}

func makeLIR(block int) {
	LIR[block] = 1
	removeFromList(block)
	delete(HIR, block)
}

func makeHIR(block int) {
	HIR[block] = 1
	delete(LIR, block)
}

func stackPrunning(removeLIR bool) {
	key, _, _ := orderedStack.PopFirst()
	keyInInt, _ := strconv.Atoi(key.(string))
	if removeLIR {
		makeHIR(keyInInt)
		if orderedList.Len() == HIRSize {
			orderedList.PopLast()
		}
		orderedList.Set(key, 1)
		orderedList.MoveFirst(key)
	}

	iter := orderedStack.Iter()
	for k, _, ok := iter.Next(); ok; k, _, ok = iter.Next() {
		keyInInt, _ := strconv.Atoi(k.(string))
		if _, ok := LIR[keyInInt]; ok {
			break
		}
		orderedStack.PopFirst()
	}
}

func printStack() {
	iter := orderedStack.Iter()
	fmt.Printf("Stack : ")
	for k, _, ok := iter.Next(); ok; k, _, ok = iter.Next() {
		fmt.Printf("%v ", k)
	}
	fmt.Println()
}

func printList() {
	iter := orderedList.Iter()
	fmt.Printf("List : ")
	for k, _, ok := iter.Next(); ok; k, _, ok = iter.Next() {
		fmt.Printf("%v ", k)
	}
	fmt.Println()
}

func printCache() {
	fmt.Printf("Cache : ")
	for key, _ := range cache {
		fmt.Printf("%v ", key)
	}
	fmt.Println()
}
