package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"nav-backend/database"
	"nav-backend/models"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func UpdateAllIcons(c *gin.Context) {
	var bookmarks []models.Bookmark
	if err := database.DB.Find(&bookmarks).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch bookmarks"})
		return
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
					continue
				}
				updated++
			}
		}
	}

	c.JSON(200, gin.H{
		"message": "Icons updated successfully",
		"total": len(bookmarks),
		"updated": updated,
	})
}

func extractDomain(urlStr string) string {
	if urlStr == "" {
		return ""
	}
	
	u, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}
	
	domain := u.Host
	if domain != "" {
		domain = strings.Split(domain, ":")[0]
	}
	
	return domain
}
