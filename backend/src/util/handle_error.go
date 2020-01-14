package util

import "fmt"

func HandleError(err error) {
	fmt.Println(err.Error())
}
