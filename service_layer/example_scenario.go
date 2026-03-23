package servicelayer

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// fake data
var fakeData = &DummyModel{
	ID:    1,
	Stock: 1000000,
}

var saveLogAfterHit []int64 // write all result stock after hitted

func ExampleScenario() {
	start := time.Now()
	log := log.Default()

	log.SetPrefix(strconv.Itoa(time.Now().Nanosecond()) + "-")
	fakeRequesterWrite := make(chan func(fakeAmount int64) int64)
	fakeRequesterRead := make(chan func() int64)
	cancel := make(chan struct{})
	var wg sync.WaitGroup
	worker := func(writer <-chan func(int64) int64, reader chan func() int64, cancel <-chan struct{}) {
		defer wg.Done()
		for {
			select {
			case w := <-writer:
				res := w(1)                                    // exec/hit the source data <<<<============[
				saveLogAfterHit = append(saveLogAfterHit, res) // save history hitted
			case r := <-reader:
				read := r()
				log.Println("read Stock: ", read)
			case <-cancel:
				return
			}
		}
	}

	// add 100000 worker for hit at the same time
	wg.Add(100000)
	for range 100000 {
		go worker(fakeRequesterWrite, fakeRequesterRead, cancel)
	}

	writeTotal := 0
	go func() {
		for {
			rand := rand.Intn(2) // we want rendom probability scenario for Read/Write
			if writeTotal < 1000000 {
				switch rand {
				case 0:
					// put all Workers at the same time
					writeTotal += 100000
					for range 100000 {
						fakeRequesterWrite <- fakeData.Writer
					}
				case 1:
					fakeRequesterRead <- fakeData.Reader
				}
			} else {
				close(cancel)
				break
			}
		}
	}()
	wg.Wait()

	duplicates := checkDuplicate(saveLogAfterHit)
	fmt.Println("=== PROOF RACE CONDITION ===")
	fmt.Println("see on function: checkDuplicate() Line 87")
	count := 0
	for _, appear := range duplicates {
		if appear > 1 {
			count++
		}
	}
	fmt.Println(count)
	fmt.Println("done : ", time.Since(start))
}

func checkDuplicate(l []int64) map[int64]int {
	counts := make(map[int64]int)
	for _, val := range l {
		counts[val]++
	}
	return counts
}
