package internal

import (
	"encoding/json"
	"fmt"
)

func Iferr(err error) {
	if err != nil {
		panic(err)
	}
}

func ifok(ok error) {
	if ok == nil {
		panic(ok)
	}
}

func contains[T comparable](list []T, element T) bool {
	for _, i := range list {
		if i == element {
			return true
		}
	}
	return false
}

func printJson(i any) {
	jsonb, err := json.Marshal(i)
	Iferr(err)
	fmt.Println(string(jsonb))
}
