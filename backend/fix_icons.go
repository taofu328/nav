package main

import (
	"crypto/md5"
	"encoding/hex"
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

	updated := 0
	for _, bookmark := range bookmarks {
		if bookmark.Icon == "" || bookmark.Icon == "/api/icons/default.svg" {
			domain := extractDomain(bookmark.URL)
			if domain != "" {
				hash := md5.Sum([]byte(domain))
				filename := hex.EncodeToString(hash[:]) + ".png"
				iconPath := "/api/icons/" + filename
				
				if err := database.DB.Model(&bookmark).Update("icon", iconPath).Error; err != nil {
					fmt.Printf("Failed to update bookmark %d: %v\n", bookmark.ID, err)
				} else {
					updated++
					fmt.Printf("Updated bookmark %d (%s): %s\n", bookmark.ID, bookmark.Title, iconPath)
				}
			}
		} else {
			fmt.Printf("Bookmark %d (%s) already has icon: %s\n", bookmark.ID, bookmark.Title, bookmark.Icon)
		}
	}

	fmt.Printf("\nTotal updated: %d/%d\n", updated, len(bookmarks))
}

func extractDomain(url string) string {
	if len(url) == 0 {
		return ""
	}
	
	domain := url
	if len(domain) > 7 && domain[:7] == "http://" {
		domain = domain[7:]
	} else {
		if len(domain) > 8 && domain[:8] == "https://" {
			domain = domain[8:]
		}
	}
	
	for i, c := range domain {
		if c == '/' || c == ':' {
			return domain[:i]
		}
	}
	
	return domain
}
