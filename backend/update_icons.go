package main

import (
	"crypto/md5"
	"encoding/hex"
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
				
				database.DB.Model(&bookmark).Update("icon", iconPath)
				updated++
			}
		}
	}

	println("Updated", updated, "bookmarks with icon paths")
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
