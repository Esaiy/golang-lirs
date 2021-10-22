package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/esaiy/golang-lirs/lirs"
)

func main() {
	var filePath string
	var cacheList []int
	var wg sync.WaitGroup

	if len(os.Args) < 3 {
		fmt.Println("program [file] [cachesize]...")
		os.Exit(1)
	}

	filePath = os.Args[1]
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("%v does not exists\n", filePath)
		os.Exit(1)
	}

	cacheList = checkCacheSize(os.Args[2:])

	for i := 2; i < len(os.Args); i++ {
		cacheSize, err := strconv.Atoi(os.Args[i])
		if err != nil {
			fmt.Printf("%v not an int\n", os.Args[i])
			os.Exit(1)
		}
		wg.Add(1)
		go lirs.LIRS(filePath, cacheSize, 1, &wg)
		fmt.Println(filePath, cacheSize, cacheList)
	}
	wg.Wait()
}

func checkCacheSize(cacheSize []string) []int {
	var cacheList []int
	for _, size := range cacheSize {
		cache, err := strconv.Atoi(size)
		if err != nil {
			fmt.Printf("%v not an int\n", size)
			os.Exit(1)
		}
		cacheList = append(cacheList, cache)

	}
	return cacheList
}

func readFile(i int, filePath string, wg *sync.WaitGroup) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fWrite, err := os.Create(strconv.Itoa(i) + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fWrite.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fWrite.WriteString(scanner.Text() + "\n")
	}
	wg.Done()
}
