package lirs

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/secnot/orderedmap"
)

type LIRSStruct struct {
	cacheSize    int
	LIRSize      int
	HIRSize      int
	hit          int
	miss         int
	writeCount   int
	orderedStack *orderedmap.OrderedMap
	orderedList  *orderedmap.OrderedMap
	LIR          map[interface{}]int
	HIR          map[interface{}]int
}

func LIRS(filePath string, totalCacheSize, hirPercentSize int, wg *sync.WaitGroup) error {
	LIRSObject := LIRSStruct{
		cacheSize:    totalCacheSize,
		LIRSize:      totalCacheSize * 99 / 100,
		HIRSize:      totalCacheSize / 100,
		hit:          0,
		miss:         0,
		writeCount:   0,
		orderedStack: orderedmap.NewOrderedMap(),
		orderedList:  orderedmap.NewOrderedMap(),
		LIR:          make(map[interface{}]int, totalCacheSize),
		HIR:          make(map[interface{}]int, totalCacheSize),
	}

	start := time.Now()

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		LIRSObject.get(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	duration := time.Since(start)

	hitRatio := 100 * float32(float32(LIRSObject.hit)/float32(LIRSObject.hit+LIRSObject.miss))

	fmt.Println("_______________________________________________________")
	fmt.Println("LIRS")
	fmt.Printf("cache size : %v\ncache hit : %v\ncache miss : %v\nhit ratio : %v\n", LIRSObject.cacheSize, LIRSObject.hit, LIRSObject.miss, hitRatio)
	fmt.Println("list size :", LIRSObject.orderedList.Len())
	fmt.Println("stack size :", LIRSObject.orderedStack.Len())
	fmt.Println("write count :", LIRSObject.writeCount)
	fmt.Printf("duration : %v\n", duration.Seconds())
	fmt.Printf("!LIRS|%v|%v|%v|", LIRSObject.cacheSize, LIRSObject.hit, LIRSObject.hit+LIRSObject.miss)

	wg.Done()
	return nil
}

func (LIRSObject *LIRSStruct) get(line string) {
	block := (strings.Split(line, ","))[0]
	op := (strings.Split(line, ","))[1]
	if op == "W" {
		(*LIRSObject).writeCount++
	}

	blockNum, err := strconv.Atoi(block)
	if err != nil {
		log.Fatal(err)
	}

	if len((*LIRSObject).LIR) < (*LIRSObject).LIRSize {
		(*LIRSObject).miss += 1
		if _, ok := (*LIRSObject).LIR[blockNum]; ok {
			(*LIRSObject).miss -= 1
			(*LIRSObject).hit += 1
			key, _, _ := (*LIRSObject).orderedStack.GetFirst()
			keyInInt, _ := strconv.Atoi(key.(string))
			if keyInInt == blockNum {
				(*LIRSObject).stackPrunning(false)
			}
		}
		(*LIRSObject).addToStack(blockNum)
		(*LIRSObject).makeLIR(blockNum)
		return
	}

	if _, ok := (*LIRSObject).LIR[blockNum]; ok {
		(*LIRSObject).hit += 1
		key, _, _ := (*LIRSObject).orderedStack.GetFirst()
		keyInInt, _ := strconv.Atoi(key.(string))
		if keyInInt == blockNum {
			(*LIRSObject).stackPrunning(false)
		}
		(*LIRSObject).addToStack(blockNum)
		return
	}

	if _, ok := (*LIRSObject).orderedList.Get(strconv.Itoa(blockNum)); ok {
		(*LIRSObject).hit += 1
		if _, ok := (*LIRSObject).orderedStack.Get(strconv.Itoa(blockNum)); ok {
			(*LIRSObject).makeLIR(blockNum)
			(*LIRSObject).removeFromList(blockNum)
			(*LIRSObject).stackPrunning(true)
		} else {
			(*LIRSObject).addToList(blockNum)
		}
		(*LIRSObject).addToStack(blockNum)
		return
	} else {
		(*LIRSObject).miss += 1
		if (*LIRSObject).orderedList.Len() == (*LIRSObject).HIRSize {
			key, _, _ := (*LIRSObject).orderedList.PopLast()
			keyInInt, _ := strconv.Atoi(key.(string))
			(*LIRSObject).removeFromList(keyInInt)
		}

		(*LIRSObject).addToList(blockNum)
		if _, ok := (*LIRSObject).orderedStack.Get(strconv.Itoa(blockNum)); ok {
			(*LIRSObject).makeLIR(blockNum)
			(*LIRSObject).removeFromList(blockNum)
			(*LIRSObject).stackPrunning(true)
		} else {
			(*LIRSObject).makeHIR(blockNum)
		}
		(*LIRSObject).addToStack(blockNum)
	}

}

func (LIRSObject *LIRSStruct) addToStack(block int) {
	key := strconv.Itoa(block)
	if _, ok := (*LIRSObject).orderedStack.Get(key); ok {
		(*LIRSObject).orderedStack.MoveLast(key)
		return
	}
	(*LIRSObject).orderedStack.Set(key, 1)
}

func (LIRSObject *LIRSStruct) addToList(block int) {
	key := strconv.Itoa(block)
	if (*LIRSObject).orderedList.Len() == (*LIRSObject).HIRSize {
		(*LIRSObject).orderedList.PopLast()
	}
	(*LIRSObject).orderedList.Set(key, 1)
	(*LIRSObject).orderedList.MoveFirst(key)
}

func (LIRSObject *LIRSStruct) removeFromList(block int) {
	key := strconv.Itoa(block)
	(*LIRSObject).orderedList.Delete(key)
}

func (LIRSObject *LIRSStruct) makeLIR(block int) {
	(*LIRSObject).LIR[block] = 1
	(*LIRSObject).removeFromList(block)
	delete((*LIRSObject).HIR, block)
}

func (LIRSObject *LIRSStruct) makeHIR(block int) {
	(*LIRSObject).HIR[block] = 1
	delete((*LIRSObject).LIR, block)
}

func (LIRSObject *LIRSStruct) stackPrunning(removeLIR bool) {
	key, _, _ := (*LIRSObject).orderedStack.PopFirst()
	keyInInt, _ := strconv.Atoi(key.(string))
	if removeLIR {
		(*LIRSObject).makeHIR(keyInInt)
		if (*LIRSObject).orderedList.Len() == (*LIRSObject).HIRSize {
			(*LIRSObject).orderedList.PopLast()
		}
		(*LIRSObject).orderedList.Set(key, 1)
		(*LIRSObject).orderedList.MoveFirst(key)
	}

	iter := (*LIRSObject).orderedStack.Iter()
	for k, _, ok := iter.Next(); ok; k, _, ok = iter.Next() {
		keyInInt, _ := strconv.Atoi(k.(string))
		if _, ok := (*LIRSObject).LIR[keyInInt]; ok {
			break
		}
		(*LIRSObject).orderedStack.PopFirst()
	}
}

// func printStack() {
// 	iter := orderedStack.Iter()
// 	fmt.Printf("Stack : ")
// 	for k, _, ok := iter.Next(); ok; k, _, ok = iter.Next() {
// 		fmt.Printf("%v ", k)
// 	}
// 	fmt.Println()
// }

// func printList() {
// 	iter := orderedList.Iter()
// 	fmt.Printf("List : ")
// 	for k, _, ok := iter.Next(); ok; k, _, ok = iter.Next() {
// 		fmt.Printf("%v ", k)
// 	}
// 	fmt.Println()
// }
