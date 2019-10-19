package db

import "fmt"

func handleError(err error) {
	fmt.Println("Error from db", err.Error())
}
