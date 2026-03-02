package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const iconsDir = "icons"

func init() {
	if err := os.MkdirAll(iconsDir, 0755); err != nil {
		panic(err)
	}
}

func GetFavicon(c *gin.Context) {
	websiteURL := c.Query("url")
	if websiteURL == "" {
		c.JSON(400, gin.H{"error": "URL parameter is required"})
		return
	}

	parsedURL, err := url.Parse(websiteURL)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid URL"})
		return
	}

	domain := parsedURL.Host
	if domain == "" {
		c.JSON(400, gin.H{"error": "Invalid URL"})
		return
	}

	hash := md5.Sum([]byte(domain))
	filename := hex.EncodeToString(hash[:]) + ".png"
	filePath := filepath.Join(iconsDir, filename)

	if _, err := os.Stat(filePath); err == nil {
		c.File(filePath)
		return
	}

	faviconURL := "https://favicon.yandex.net/favicon/" + domain
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	resp, err := client.Get(faviconURL)
	if err != nil {
		svgPath := filepath.Join(iconsDir, "default.svg")
		c.Header("Content-Type", "image/svg+xml")
		c.File(svgPath)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		svgPath := filepath.Join(iconsDir, "default.svg")
		c.Header("Content-Type", "image/svg+xml")
		c.File(svgPath)
		return
	}

	out, err := os.Create(filePath)
	if err != nil {
		svgPath := filepath.Join(iconsDir, "default.svg")
		c.Header("Content-Type", "image/svg+xml")
		c.File(svgPath)
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		svgPath := filepath.Join(iconsDir, "default.svg")
		c.Header("Content-Type", "image/svg+xml")
		c.File(svgPath)
		return
	}

	c.File(filePath)
}

func FetchFavicon(websiteURL string) (string, error) {
	parsedURL, err := url.Parse(websiteURL)
	if err != nil {
		return "", err
	}

	domain := parsedURL.Host
	if domain == "" {
		return "", err
	}

	hash := md5.Sum([]byte(domain))
	filename := hex.EncodeToString(hash[:]) + ".png"
	filePath := filepath.Join(iconsDir, filename)

	if _, err := os.Stat(filePath); err == nil {
		return "/api/icons/" + filename, nil
	}

	faviconURL := "https://favicon.yandex.net/favicon/" + domain
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	resp, err := client.Get(faviconURL)
	if err != nil {
		return "/api/icons/default.svg", nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "/api/icons/default.svg", nil
	}

	out, err := os.Create(filePath)
	if err != nil {
		return "/api/icons/default.svg", nil
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "/api/icons/default.svg", nil
	}

	return "/api/icons/" + filename, nil
}

func ServeIcon(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(400, gin.H{"error": "Filename is required"})
		return
	}

	if filename == "default.png" {
		svgPath := filepath.Join(iconsDir, "default.svg")
		if _, err := os.Stat(svgPath); err == nil {
			c.Header("Content-Type", "image/svg+xml")
			c.File(svgPath)
			return
		}
	}

	filePath := filepath.Join(iconsDir, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		svgPath := filepath.Join(iconsDir, "default.svg")
		if _, err := os.Stat(svgPath); err == nil {
			c.Header("Content-Type", "image/svg+xml")
			c.File(svgPath)
			return
		}
		c.JSON(404, gin.H{"error": "Icon not found"})
		return
	}

	c.File(filePath)
}

func CleanInvalidIcons(c *gin.Context) {
	entries, err := os.ReadDir(iconsDir)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read icons directory"})
		return
	}

	cleaned := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		if entry.Name() == "default.png" {
			continue
		}

		filepath := filepath.Join(iconsDir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}

		if time.Since(info.ModTime()) > 30*24*time.Hour {
			if err := os.Remove(filepath); err == nil {
				cleaned++
			}
		}
	}

	c.JSON(200, gin.H{
		"message": "Cleaned successfully",
		"cleaned": cleaned,
	})
}

func GetIconFromURL(iconURL string) string {
	if iconURL == "" {
		return "/api/icons/default.png"
	}

	if strings.HasPrefix(iconURL, "/api/icons/") {
		return iconURL
	}

	parsedURL, err := url.Parse(iconURL)
	if err != nil {
		return "/api/icons/default.png"
	}

	if parsedURL.Host == "" {
		return "/api/icons/default.png"
	}

	hash := md5.Sum([]byte(parsedURL.Host))
	filename := hex.EncodeToString(hash[:]) + ".png"
	filepath := filepath.Join(iconsDir, filename)

	if _, err := os.Stat(filepath); err == nil {
		return "/api/icons/" + filename
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	resp, err := client.Get(iconURL)
	if err != nil {
		return "/api/icons/default.png"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "/api/icons/default.png"
	}

	out, err := os.Create(filepath)
	if err != nil {
		return "/api/icons/default.png"
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "/api/icons/default.png"
	}

	return "/api/icons/" + filename
}
