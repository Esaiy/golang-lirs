package lirs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/esaiy/golang-lirs/simulator"
	"github.com/secnot/orderedmap"
)

type LIRS struct {
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
	cache        map[interface{}]bool
}

func NewLIRS(cacheSize, HIRSize int) *LIRS {
	if HIRSize > 100 || HIRSize < 0 {
		log.Fatal("HIRSize must be between 0 and 100")
	}
	LIRCapacity := (100 - HIRSize) * cacheSize / 100
	HIRCapacity := HIRSize * cacheSize / 100
	return &LIRS{
		cacheSize:    cacheSize,
		LIRSize:      LIRCapacity,
		HIRSize:      HIRCapacity,
		hit:          0,
		miss:         0,
		writeCount:   0,
		orderedStack: orderedmap.NewOrderedMap(),
		orderedList:  orderedmap.NewOrderedMap(),
		LIR:          make(map[interface{}]int, LIRCapacity),
		HIR:          make(map[interface{}]int, HIRCapacity),
		cache:        make(map[interface{}]bool, cacheSize),
	}
}

// func ssss(filePath string, totalCacheSize, hirPercentSize int, wg *sync.WaitGroup) error {
// 	LIRSObject := LIRS{
// 		cacheSize:    totalCacheSize,
// 		LIRSize:      totalCacheSize * 99 / 100,
// 		HIRSize:      totalCacheSize / 100,
// 		hit:          0,
// 		miss:         0,
// 		writeCount:   0,
// 		orderedStack: orderedmap.NewOrderedMap(),
// 		orderedList:  orderedmap.NewOrderedMap(),
// 		LIR:          make(map[interface{}]int, totalCacheSize),
// 		HIR:          make(map[interface{}]int, totalCacheSize),
// 	}

// 	start := time.Now()

// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)

// 	for scanner.Scan() {
// 		LIRSObject.get(scanner.Text())
// 	}

// 	if err := scanner.Err(); err != nil {
// 		log.Fatal(err)
// 	}

// 	duration := time.Since(start)

// 	hitRatio := 100 * float32(float32(LIRSObject.hit)/float32(LIRSObject.hit+LIRSObject.miss))

// 	fmt.Println("_______________________________________________________")
// 	fmt.Println("LIRS")
// 	fmt.Printf("cache size : %v\ncache hit : %v\ncache miss : %v\nhit ratio : %v\n", LIRSObject.cacheSize, LIRSObject.hit, LIRSObject.miss, hitRatio)
// 	fmt.Println("list size :", LIRSObject.orderedList.Len())
// 	fmt.Println("stack size :", LIRSObject.orderedStack.Len())
// 	fmt.Println("write count :", LIRSObject.writeCount)
// 	fmt.Printf("duration : %v\n", duration.Seconds())
// 	fmt.Printf("!LIRS|%v|%v|%v|\n", LIRSObject.cacheSize, LIRSObject.hit, LIRSObject.hit+LIRSObject.miss)

// 	wg.Done()
// 	return nil
// }

func (LIRSObject *LIRS) Get(trace simulator.Trace) (err error) {
	block := trace.Addr
	op := trace.Op
	if op == "W" {
		LIRSObject.writeCount++
	}

	if len(LIRSObject.LIR) < LIRSObject.LIRSize {
		// LIR is not full; there is space in cache
		LIRSObject.miss += 1
		if _, ok := LIRSObject.LIR[block]; ok {
			// block is in LIR, not a miss
			LIRSObject.miss -= 1
			LIRSObject.hit += 1
		}
		LIRSObject.addToStack(block)
		LIRSObject.makeLIR(block)
		return nil
	}

	if _, ok := LIRSObject.LIR[block]; ok {
		// hit, block is in LIR
		LIRSObject.handleLIRBlock(block)
	} else if _, ok := LIRSObject.orderedList.Get(strconv.Itoa(block)); ok {
		// hit, block is HIR resident
		LIRSObject.handleHIRResidentBlock(block)
	} else {
		// miss, blok is HIR non resident
		LIRSObject.HandleHIRNonResidentBlock(block)
	}
	return nil
}

func (LIRSObject *LIRS) PrintToFile(file *os.File, start time.Time) (err error) {
	duration := time.Since(start)
	hitRatio := 100 * float32(float32(LIRSObject.hit)/float32(LIRSObject.hit+LIRSObject.miss))
	result := fmt.Sprintf(`_______________________________________________________
LIRS
cache size : %v
cache hit : %v
cache miss : %v
hit ratio : %v
list size : %v
stack size : %v
lir capacity: %v
hir capacity: %v
write count : %v
duration : %v
!LIRS|%v|%v|%v
`, LIRSObject.cacheSize, LIRSObject.hit, LIRSObject.miss, hitRatio, LIRSObject.orderedList.Len(), LIRSObject.orderedStack.Len(), LIRSObject.LIRSize, LIRSObject.HIRSize, LIRSObject.writeCount, duration.Seconds(), LIRSObject.cacheSize, LIRSObject.hit, LIRSObject.hit+LIRSObject.miss)
	_, err = file.WriteString(result)
	return err
}

func (LIRSObject *LIRS) handleLIRBlock(block int) (err error) {
	LIRSObject.hit += 1
	key, _, ok := LIRSObject.orderedStack.GetFirst()
	if !ok {
		return errors.New("orderedStack is empty")
	}
	keyInInt, err := strconv.Atoi(key.(string))
	if err != nil {
		return err
	}
	if keyInInt == block {
		// block is in LIR and at the bottom of the stack
		// do stack pruning
		LIRSObject.stackPrunning(false)
	}
	LIRSObject.addToStack(block)
	return nil
}

func (LIRSObject *LIRS) handleHIRResidentBlock(block int) {
	LIRSObject.hit += 1
	if _, ok := LIRSObject.orderedStack.Get(strconv.Itoa(block)); ok {
		// block is in stack, move to LIR
		LIRSObject.makeLIR(block)
		LIRSObject.removeFromList(block)
		LIRSObject.stackPrunning(true)
	} else {
		// block is not in stack, move to end of list
		LIRSObject.orderedList.MoveLast(block)
	}
	LIRSObject.addToStack(block)
}

func (LIRSObject *LIRS) HandleHIRNonResidentBlock(block int) {
	LIRSObject.miss += 1
	LIRSObject.addToList(block)
	if _, ok := LIRSObject.orderedStack.Get(strconv.Itoa(block)); ok {
		// block is in stack, move to LIR
		LIRSObject.makeLIR(block)
		LIRSObject.removeFromList(block)
		LIRSObject.stackPrunning(true)
	} else {
		LIRSObject.makeHIR(block)
	}
	LIRSObject.addToStack(block)
}

func (LIRSObject *LIRS) addToStack(block int) {
	key := strconv.Itoa(block)
	if _, ok := LIRSObject.orderedStack.Get(key); ok {
		LIRSObject.orderedStack.MoveLast(key)
		return
	}
	LIRSObject.orderedStack.Set(key, 1)
}

// list
// front queue (paper) = last ordered list
// end queue (paper) = first ordered list
// func (LIRSObject *LIRS) addToListToLastIndex(block int) {
// 	key := strconv.Itoa(block)
// 	if LIRSObject.orderedList.Len() == LIRSObject.HIRSize {
// 		LIRSObject.orderedList.PopLast()
// 	}
// 	LIRSObject.orderedList.Set(key, 1)
// 	LIRSObject.orderedList.MoveFirst(key)
// }

func (LIRSObject *LIRS) addToList(block int) {
	key := strconv.Itoa(block)
	if LIRSObject.orderedList.Len() == LIRSObject.HIRSize {
		LIRSObject.orderedList.PopFirst()
	}
	LIRSObject.orderedList.Set(key, 1)
}

func (LIRSObject *LIRS) removeFromList(block int) {
	key := strconv.Itoa(block)
	LIRSObject.orderedList.Delete(key)
}

func (LIRSObject *LIRS) makeLIR(block int) {
	LIRSObject.LIR[block] = 1
	LIRSObject.removeFromList(block)
	delete(LIRSObject.HIR, block)
}

func (LIRSObject *LIRS) makeHIR(block int) {
	LIRSObject.HIR[block] = 1
	delete(LIRSObject.LIR, block)
}

func (LIRSObject *LIRS) stackPrunning(removeLIR bool) (err error) {
	key, _, ok := LIRSObject.orderedStack.PopFirst()
	if !ok {
		return errors.New("orderedStack is empty")
	}
	keyInInt, err := strconv.Atoi(key.(string))
	if err != nil {
		return err
	}
	if removeLIR {
		LIRSObject.makeHIR(keyInInt)
		LIRSObject.orderedList.Set(key, 1)
		LIRSObject.orderedList.MoveLast(key)
	}

	iter := LIRSObject.orderedStack.Iter()
	for k, _, ok := iter.Next(); ok; k, _, ok = iter.Next() {
		keyInInt, err := strconv.Atoi(k.(string))
		if err != nil {
			return err
		}
		if _, ok := LIRSObject.LIR[keyInInt]; ok {
			break
		}
		LIRSObject.orderedStack.PopFirst()
	}
	return nil
}
