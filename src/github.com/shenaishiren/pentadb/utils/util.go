package utils

import (
	"errors"
	"math/rand"
)

func RandomChoice(list []interface{}, k int) ([]interface{}, error) {
	if k <= 0 {
		return nil, errors.New("invalid k: k must be > 0")
	}
	pool := list
	n := len(pool)
	result := make([]interface{}, k)
	for i := 0; i < k; i++ {
		j := rand.Intn(n - i)
		result[i] = pool[j]
		pool[j] = result[n - i - 1]
	}

	return result, nil
}

func Filter(lambda func(interface{}) bool, list []interface{}) []interface{} {
	result := make([]interface{}, len(list))
	for _, v := range list {
		if lambda(v) {
			result = append(result, v)
		}
	}

	return result
}