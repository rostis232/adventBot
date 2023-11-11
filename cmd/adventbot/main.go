package main

import (
	"fmt"

	"github.com/rostis232/adventBot/internal/repository"
)

func main(){
	db, err := repository.NewSQLiteDB("sqlite3.db")
	if err != nil {
		fmt.Println("Помилка!", err)
	} else {
		fmt.Println("Ok!")
	}
	db.Close()
}