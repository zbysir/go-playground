package merge_input

import (
	"testing"
)

func Test_mergeInput(t *testing.T) {
	a := make(chan int, 50)
	b := make(chan int, 50)
	go func() {
		for i := 0; i < 50; i++ {
			a <- i*2 + 1
		}
		close(a)
	}()
	go func() {
		for i := 0; i < 30; i++ {
			b <- i*2 + 2
		}
		close(b)
	}()

	d := make(chan int, 10)

	go func() {
		MergeInput(a, b, d)
		close(d)
	}()

	for i := range d {
		t.Log(i)
	}

	t.Log("finish")
}

func TestOrderBuff(t *testing.T) {
	var d = make(chan int, 10)
	x := NewOrderBuff(10, d)
	x.Put(1)
	x.Put(3)
	x.Put(2)

	x.Stop()

	for i := range d {
		println(i)
	}
}
