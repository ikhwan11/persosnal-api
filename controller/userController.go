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
	"golang.org/x/crypto/bcrypt"
)

// =========== CREATE USER ===========

func CreateUserHandler(c *gin.Context) {
	name := c.PostForm("name")
	username := c.PostForm("username")
	password := c.PostForm("password")
	about := c.PostForm("about")

	slug := utils.GenerateSlug(name)

	imagePath := "public/image/default.jpg"

	// Upload image
	file, err := c.FormFile("image")
	if err == nil {
		ext := filepath.Ext(file.Filename)
		cleanName := strings.ReplaceAll(strings.ToLower(name), " ", "-")
		dateStr := time.Now().Format("20060102")

		uniqueName := fmt.Sprintf("%s-user-image-%s%s", cleanName, dateStr, ext)
		imagePath = "public/image/" + uniqueName

		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Name:     name,
		Slug:     slug,
		Username: username,
		Password: string(hashedPassword),
		Image:    imagePath,
		About:    about,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created",
		"user":    user,
	})
}

// =========== GET USERS ===========

func GetUsersHandler(c *gin.Context) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func GetUserBySlugHandler(c *gin.Context) {
	slug := c.Param("slug")
	var user models.User

	if err := database.DB.Where("slug = ?", slug).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// =========== UPDATE USER ===========

func UpdateUserHandler(c *gin.Context) {
	slug := c.Param("slug")
	var oldUser models.User

	// Cari user lama
	if err := database.DB.Where("slug = ?", slug).First(&oldUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Form data
	name := c.PostForm("name")
	username := c.PostForm("username")
	password := c.PostForm("password")
	about := c.PostForm("about")

	// Update slug kalau nama diubah
	if name != "" {
		oldUser.Name = name
		oldUser.Slug = utils.GenerateSlug(name)
	}

	if username != "" {
		oldUser.Username = username
	}

	if about != "" {
		oldUser.About = about
	}

	// Handle image upload
	file, err := c.FormFile("image")
	if err == nil {
		ext := filepath.Ext(file.Filename)
		cleanName := strings.ReplaceAll(strings.ToLower(oldUser.Name), " ", "-")
		dateStr := time.Now().Format("20060102")

		uniqueName := fmt.Sprintf("%s-user-image-%s%s", cleanName, dateStr, ext)
		newImagePath := "public/image/" + uniqueName

		// Hapus foto lama kalau bukan default
		if oldUser.Image != "public/image/default.jpg" {
			_ = os.Remove(oldUser.Image)
		}

		if err := c.SaveUploadedFile(file, newImagePath); err == nil {
			oldUser.Image = newImagePath
		}
	}

	// Update password
	if password != "" {
		newPass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		oldUser.Password = string(newPass)
	}

	// Simpan perubahan
	if err := database.DB.Save(&oldUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated",
		"user":    oldUser,
	})
}
