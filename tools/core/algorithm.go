package core

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

func ShuffleSlice(slice any) {
	rand.Seed(time.Now().UnixNano())

	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		fmt.Println("Input is not a slice.")
		return
	}

	// Convert the interface{} to a []interface{}
	length := rv.Len()
	shuffled := make([]any, length)
	for i := 0; i < length; i++ {
		shuffled[i] = rv.Index(i).Interface()
	}

	// Shuffle the []interface{}
	rand.Shuffle(length, func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	// Copy the shuffled values back to the original slice
	for i := 0; i < length; i++ {
		rv.Index(i).Set(reflect.ValueOf(shuffled[i]))
	}
}

func FisherYatesShuffle(slice []string) {
	rand.Seed(time.Now().UnixNano())

	n := len(slice)
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}
