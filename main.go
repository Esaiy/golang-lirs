package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	stack, list []int
	hit, miss       = 0, 0
	cacheSize   int = 500000
	LIRSize     int = cacheSize * 99 / 100
	HIRSize     int = cacheSize / 100
	LIR             = make(map[int]int)
	HIR             = make(map[int]int)
)

func main() {
	var (
		start    = time.Now()
		filePath = "./test/data/web1_input.txt"
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

	fmt.Printf("HIT : %v\nMISS : %v\nHIT RATE : %v\n", hit, miss, (hit+miss)/2)
	fmt.Printf("Duration : %v\n", duration.Seconds())
}

func lirs(line string) {
	var block = (strings.Split(line, ","))[0]

	var blockNum, err = strconv.Atoi(block)
	if err != nil {
		log.Fatal(err)
	}

	if len(stack) < LIRSize {
		miss += 1
		addToStack(blockNum)
		makeLIR(blockNum)
		return
	}

	if value, _ := LIR[blockNum]; value == 1 {
		hit += 1
		addToStack(blockNum)
		return
	}

	if value, _ := HIR[blockNum]; value == 0 {
		miss += 1
		removeListFront()

		if _, found := isBlockExists(blockNum, stack); found {
			makeLIR(blockNum)
			addToStack(blockNum)
			stackPrunning()
		} else {
			addToList(blockNum)
		}
		return
	} else {
		hit += 1
		if _, found := isBlockExists(blockNum, stack); found {
			makeLIR(blockNum)
			stackPrunning()
		} else {
			addToList(blockNum)
		}
		addToStack(blockNum)
	}

}

func addToStack(block int) {
	if pos, found := isBlockExists(block, stack); found {
		RemoveIndex(stack, pos)
	}

	stack = append([]int{block}, stack...)
}

func addToList(block int) {
	if pos, found := isBlockExists(block, list); found {
		RemoveIndex(list, pos)
	}

	list = append([]int{block}, list...)
}

func isBlockExists(block int, slice []int) (int, bool) {
	pos, found := Find(slice, block)
	if found {
		return pos, true
	}

	return -1, false
}

func Find(slice []int, val int) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func RemoveIndex(s []int, index int) []int {
	ret := make([]int, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

func makeLIR(block int) {
	LIR[block] = 1
	HIR[block] = 0
}

func makeHIR(block int) {
	HIR[block] = 1
	LIR[block] = 0
}

func removeListFront() {
	if len(list) == 0 {
		return
	}
	list = list[1:]
}

func stackPrunning() {
	var first bool = false
	for i := len(stack) - 1; i >= 0; i-- {
		if value, _ := LIR[stack[i]]; value == 1 {
			if !first {
				first = true
				makeHIR(stack[i])
			} else {
				stack = stack[:i-1]
			}
		}
	}
}
