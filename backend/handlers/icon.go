package handlers

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"nav-backend/config"

	"github.com/gin-gonic/gin"
)

func init() {
	if err := os.MkdirAll(config.IconsDir, 0755); err != nil {
		panic(err)
	}
}

func GetFavicon(c *gin.Context) {
	// 记录请求开始时间
	startTime := time.Now()
	
	// 记录请求基本信息
	log.Printf("[Icon Fetch] Request started at: %s", startTime.Format("2006-01-02 15:04:05"))
	
	websiteURL := c.Query("url")
	log.Printf("[Icon Fetch] Request parameter - URL: %s", websiteURL)
	
	if websiteURL == "" {
		log.Printf("[Icon Fetch] Error: URL parameter is required")
		c.JSON(400, gin.H{"error": "URL parameter is required"})
		return
	}

	parsedURL, err := url.Parse(websiteURL)
	if err != nil {
		log.Printf("[Icon Fetch] Error parsing URL: %v", err)
		c.JSON(400, gin.H{"error": "Invalid URL"})
		return
	}

	domain := parsedURL.Host
	if domain == "" {
		log.Printf("[Icon Fetch] Error: Invalid domain")
		c.JSON(400, gin.H{"error": "Invalid URL"})
		return
	}

	// 使用域名和尺寸信息生成文件名，确保不同尺寸的图标可以共存
	hash := md5.Sum([]byte(domain + "_large"))
	filename := hex.EncodeToString(hash[:]) + ".png"
	filePath := filepath.Join(config.IconsDir, filename)
	
	log.Printf("[Icon Fetch] Icon info - Domain: %s, Type: PNG, Size: Large, Filename: %s", domain, filename)

	// 从 favicon.im 获取图标（使用 larger=true 参数获取大尺寸图标）
	faviconURL := "https://wsrv.nl/?url=https://favicon.im/zh/" + domain + "?larger=true"
	log.Printf("[Icon Fetch] Fetching from URL: %s", faviconURL)
	
	downloadStartTime := time.Now()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	resp, err := client.Get(faviconURL)
	if err != nil {
		log.Printf("[Icon Fetch] Error fetching icon - Type: Network Error, Code: 0, Description: %v", err)
		c.JSON(200, gin.H{
			"icon":    "/api/icons/default.svg",
			"message": "Failed to fetch icon, using default",
			"error":   err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[Icon Fetch] Error fetching icon - Type: HTTP Error, Code: %d, Description: Server returned non-200 status", resp.StatusCode)
		c.JSON(200, gin.H{
			"icon":    "/api/icons/default.svg",
			"message": "Failed to fetch icon, using default",
			"status":  resp.StatusCode,
		})
		return
	}
	
	downloadDuration := time.Since(downloadStartTime)
	log.Printf("[Icon Fetch] Download completed - Duration: %v", downloadDuration)

	// 读取图片数据以获取详细信息
	imageData := new(bytes.Buffer)
	imageSize, err := imageData.ReadFrom(resp.Body)
	if err != nil {
		log.Printf("[Icon Fetch] Error reading image data: %v", err)
		c.JSON(200, gin.H{
			"icon":    "/api/icons/default.svg",
			"message": "Failed to read icon data, using default",
			"error":   err.Error(),
		})
		return
	}
	
	// 解析图片以获取尺寸信息
	imgConfig, format, err := image.DecodeConfig(bytes.NewReader(imageData.Bytes()))
	if err != nil {
		log.Printf("[Icon Fetch] Warning: Could not decode image config: %v", err)
	} else {
		log.Printf("[Icon Fetch] Image details - Format: %s, Size: %d bytes, Dimensions: %dx%d, Hash: %s", 
			format, imageSize, imgConfig.Width, imgConfig.Height, hex.EncodeToString(hash[:]))
	}

	// 保存图标文件
	saveStartTime := time.Now()
	out, err := os.Create(filePath)
	if err != nil {
		log.Printf("[Icon Fetch] Error saving icon - Type: File System Error, Code: 0, Description: %v", err)
		c.JSON(200, gin.H{
			"icon":    "/api/icons/default.svg",
			"message": "Failed to save icon, using default",
			"error":   err.Error(),
		})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, bytes.NewReader(imageData.Bytes())); err != nil {
		log.Printf("[Icon Fetch] Error saving icon - Type: File System Error, Code: 0, Description: %v", err)
		c.JSON(200, gin.H{
			"icon":    "/api/icons/default.svg",
			"message": "Failed to save icon, using default",
			"error":   err.Error(),
		})
		return
	}
	
	// 获取文件权限
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Printf("[Icon Fetch] Warning: Could not get file info: %v", err)
	} else {
		log.Printf("[Icon Fetch] Save info - Path: %s, Filename: %s, Permissions: %v, Status: Success", 
			filePath, filename, fileInfo.Mode())
	}
	
	saveDuration := time.Since(saveStartTime)
	totalDuration := time.Since(startTime)
	
	log.Printf("[Icon Fetch] Processing time - Download: %v, Save: %v, Total: %v", 
		downloadDuration, saveDuration, totalDuration)
	log.Printf("[Icon Fetch] Request completed successfully at: %s", time.Now().Format("2006-01-02 15:04:05"))

	c.JSON(200, gin.H{
		"icon":    "/api/icons/" + filename,
		"message": "Icon fetched successfully",
	})
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

	// 使用域名和尺寸信息生成文件名，确保不同尺寸的图标可以共存
	hash := md5.Sum([]byte(domain + "_large"))
	filename := hex.EncodeToString(hash[:]) + ".png"
	filePath := filepath.Join(config.IconsDir, filename)

	if _, err := os.Stat(filePath); err == nil {
		return "/api/icons/" + filename, nil
	}

	// 从 favicon.im 获取图标（使用 larger=true 参数获取大尺寸图标）
	faviconURL := "https://wsrv.nl/?url=https://favicon.im/zh/" + domain + "?larger=true"
	
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

	// 添加缓存头，避免浏览器不断请求图标
	c.Header("Cache-Control", "public, max-age=86400") // 缓存1天
	c.Header("Expires", time.Now().Add(24*time.Hour).Format(time.RFC1123))

	if filename == "default.png" {
		svgPath := filepath.Join(config.IconsDir, "default.svg")
		if _, err := os.Stat(svgPath); err == nil {
			c.Header("Content-Type", "image/svg+xml")
			c.File(svgPath)
			return
		}
	}

	filePath := filepath.Join(config.IconsDir, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		svgPath := filepath.Join(config.IconsDir, "default.svg")
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
	entries, err := os.ReadDir(config.IconsDir)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read icons directory"})
		return
	}

	cleaned := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		if entry.Name() == "default.png" || entry.Name() == "default.svg" {
			continue
		}

		filepath := filepath.Join(config.IconsDir, entry.Name())
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

func DeleteIcon(c *gin.Context) {
	iconURL := c.Query("url")
	if iconURL == "" {
		c.JSON(400, gin.H{"error": "URL parameter is required"})
		return
	}

	// 从图标路径中提取文件名
	if strings.HasPrefix(iconURL, "/api/icons/") {
		filename := strings.TrimPrefix(iconURL, "/api/icons/")
		filePath := filepath.Join(config.IconsDir, filename)
		
		// 不允许删除默认图标
		if filename == "default.svg" || filename == "default.png" {
			c.JSON(400, gin.H{"error": "Cannot delete default icon"})
			return
		}
		
		// 检查文件是否存在
		if _, err := os.Stat(filePath); err != nil {
			// 文件不存在，直接返回成功响应
			log.Printf("[Icon Delete] File not found, returning success: %s", filePath)
			c.JSON(200, gin.H{
				"message": "Icon deleted successfully",
			})
			return
		}
		
		// 删除文件
		if err := os.Remove(filePath); err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete icon"})
			return
		}
		
		// 文件存在并成功删除
		log.Printf("[Icon Delete] File deleted successfully: %s", filePath)
		c.JSON(200, gin.H{
			"message": "Icon deleted successfully",
		})
		return
	}
	
	c.JSON(400, gin.H{"error": "Invalid icon URL"})
}

func UploadIcon(c *gin.Context) {
	file, err := c.FormFile("icon")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file uploaded"})
		return
	}

	// 检查文件类型
	ext := filepath.Ext(file.Filename)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".svg" && ext != ".ico" {
		c.JSON(400, gin.H{"error": "Unsupported file type. Only PNG, JPG, JPEG, SVG, and ICO are allowed."})
		return
	}

	// 生成唯一文件名
	hash := md5.Sum([]byte(file.Filename + time.Now().String()))
	filename := hex.EncodeToString(hash[:]) + ext
	filePath := filepath.Join(config.IconsDir, filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(500, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Icon uploaded successfully",
		"icon":    "/api/icons/" + filename,
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
	filepath := filepath.Join(config.IconsDir, filename)

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
