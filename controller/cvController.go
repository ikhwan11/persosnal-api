package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"my-personal-web/database"
	"my-personal-web/models"
	"my-personal-web/utils"

	"github.com/gin-gonic/gin"
)

// CREATE =============================

func CreateCVHandler(c *gin.Context) {
	name := c.PostForm("name")
	tagline := c.PostForm("tagline")
	about := c.PostForm("about")

	slug := utils.GenerateSlug(name)

	imagePath := "public/image/default.jpg"

	file, err := c.FormFile("image")
	if err == nil {
		ext := filepath.Ext(file.Filename)
		cleanName := strings.ReplaceAll(strings.ToLower(name), " ", "-")
		dateStr := time.Now().Format("20060102")

		uniqueName := fmt.Sprintf("%s-cv-image-%s%s", cleanName, dateStr, ext)
		imagePath = "public/image/" + uniqueName

		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
	}

	cv := models.CV{
		Name:    name,
		Slug:    slug,
		Image:   imagePath,
		Tagline: tagline,
		About:   about,
	}

	if err := database.DB.Create(&cv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create CV"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "CV created successfully",
		"cv":      cv,
	})
}

// GET =============================
func GetCVHandler(c *gin.Context) {
	var cv []models.CV
	if err := database.DB.Find(&cv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": cv})
}

// UPDATE =============================
func UpdateCVHandler(c *gin.Context) {
	slug := c.Param("slug")
	var oldcv models.CV

	if err := database.DB.Where("slug = ?", slug).First(&oldcv).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	name := c.PostForm("name")
	tagline := c.PostForm("tagline")
	about := c.PostForm("about")

	if name != "" {
		oldcv.Name = name
		oldcv.Slug = utils.GenerateSlug(name)
	}

	if tagline != "" {
		oldcv.Tagline = tagline
	}

	if about != "" {
		oldcv.About = about
	}

	file, err := c.FormFile("image")
	if err == nil {
		ext := filepath.Ext(file.Filename)
		cleanName := strings.ReplaceAll(strings.ToLower(name), " ", "-")
		dateStr := time.Now().Format("20060102")

		uniqueName := fmt.Sprintf("%s-cv-image-%s%s", cleanName, dateStr, ext)
		newImagePath := "public/image/" + uniqueName

		if oldcv.Image != "public/image/default.jpg" {
			_ = os.Remove(oldcv.Image)
		}

		if err := c.SaveUploadedFile(file, newImagePath); err == nil {
			oldcv.Image = newImagePath
		}
	}

	if err := database.DB.Save(&oldcv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cv"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "cv updated",
		"user":    oldcv,
	})
}
