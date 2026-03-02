package main

import (
	"fmt"
	"nav-backend/database"
	"nav-backend/models"
)

func main() {
	database.InitDB()

	var bookmarks []models.Bookmark
	if err := database.DB.Find(&bookmarks).Error; err != nil {
		panic(err)
	}

	fmt.Println("Total bookmarks:", len(bookmarks))
	fmt.Println("\nFirst 5 bookmarks:")
	for i := 0; i < 5 && i < len(bookmarks); i++ {
		b := bookmarks[i]
		fmt.Printf("%d. ID: %d, Title: %s, Icon: '%s'\n", i+1, b.ID, b.Title, b.Icon)
	}
}
