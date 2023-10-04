package follows

import (
	"fmt"
	"peargram/database"
)

func GetFollowerAmount(username string) int {
	var amount int

	DB := database.ConnectDB()
	err := DB.QueryRow("SELECT COUNT(*) FROM follows WHERE target=?", username).Scan(&amount)
	if err != nil {
		fmt.Println(err)
	}

	return amount
}

func GetFollowingAmount(username string) int {
	var amount int

	DB := database.ConnectDB()
	err := DB.QueryRow("SELECT COUNT(*) FROM follows WHERE actor=?", username).Scan(&amount)
	if err != nil {
		fmt.Println(err)
	}

	return amount
}

func IsFollowing(username string, target string) bool {
	var amount int

	DB := database.ConnectDB()
	err := DB.QueryRow("SELECT COUNT(*) FROM follows WHERE actor=? AND target=?", username, target).Scan(&amount)
	if err != nil {
		fmt.Println(err)
	}

	return amount == 1
}

func GetFollowers(username string) []string {
	var users []string

	DB := database.ConnectDB()
	err := DB.Select(&users, "SELECT actor FROM follows WHERE target=?", username)
	if err != nil {
		fmt.Println(err)
	}

	return users
}

func GetFollowings(username string) []string {
	var users []string

	DB := database.ConnectDB()
	err := DB.Select(&users, "SELECT target FROM follows WHERE actor=?", username)
	if err != nil {
		fmt.Println(err)
	}

	return users
}
