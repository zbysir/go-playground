package merge_input

import (
	"sync"
)

// 题目说明:
//   启动两个线程, 一个输出 1,3,5,7…99, 另一个输出 2,4,6,8…100 最后 STDOUT 中按序输出 1,2,3,4,5…100

// 扩展思维: 将多个接收到的多个无序UDP包组成有序的包

func MergeInput(a chan int, b chan int, d chan int) {
	var buff = NewOrderBuff(100, d)

	for {
		if a == nil && b == nil {
			break
		}

		select {
		case i, ok := <-a:
			if ok {
				buff.Put(i)
			} else {
				a = nil
			}
		case j, ok := <-b:
			if ok {
				buff.Put(j)
			} else {
				b = nil
			}
		}
	}

	buff.Stop()

	return
}

type OrderBuff struct {
	buff  []int
	out   chan int
	index int // 当前处理到的消息序号
	lock  sync.Mutex
}

// 注意 长度 一定需要比数据序号大
// 由于Put不是阻塞的,
// 如何处理这个问题:
//
func NewOrderBuff(l int, out chan int) OrderBuff {
	return OrderBuff{buff: make([]int, l), out: out, index: 1}
}

func (b *OrderBuff) Stop() {
	b.Refresh()
}

func (b *OrderBuff) Refresh() {
	// 检查缓冲中是否有已经排序好的完整的数据
	// 将缓冲区中的完整数据输出
	var j int
	for j = 0; j < len(b.buff); j++ {
		if b.buff[j] == 0 {
			break
		}
	}

	if j != 0 {
		for jj := 0; jj < j; jj++ {
			b.out <- b.buff[jj]
			b.index++
		}
		copy(b.buff, b.buff[j:])
	}

}

func (b *OrderBuff) Put(i int) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if i == b.index {
		b.out <- i
		b.index++
		copy(b.buff, b.buff[1:])
	} else {
		// 放在缓冲中
		// 判断是否溢出
		// 如果溢出了如何处理? 暂时不知道如何处理
		b.buff[i-b.index] = i

		b.Refresh()
	}
}
