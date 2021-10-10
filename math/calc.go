package math

import (
	"log"
	"math"
	"math/rand"
	"time"
)

func calcReainder(dividend, divisor int) {
	if divisor == 0 {
		log.Fatal("divisor is not zero")
	}

	// divisor 除数, dividend 被除数  dividend/divisor; dividend%divisor
	log.Printf("[除法]运算结果:%d, [求余]运算结果:%d", dividend/divisor, dividend%divisor)

	// 被除数 % 除数：如果被除数小于除数时，其结果商等于0，且余数是被除数本身

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := rand.Intn(math.MaxInt32 - 1)
	log.Println("init value:", index)

	len := 7 // 可能变化，但 result 必须 [0, len-1]
	for i := 0; i < 10; i++ {
		result := index % len
		index = (index + 1) % len
		log.Printf("result:%d, index:%d", result, index)
	}
}
