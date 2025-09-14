package controller

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"my-personal-web/api/config"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Image    string `json:"image"`
	About    string `json:"about"`
}

func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}

func CreateUserHandler(c *gin.Context) {
	name := c.PostForm("name")
	username := c.PostForm("username")
	password := c.PostForm("password")
	about := c.PostForm("about")

	slug := generateSlug(name)

	imagePath := "public/image/default.jpg"

	file, err := c.FormFile("image")
	if err == nil {
		ext := filepath.Ext(file.Filename)
		uniqueName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		imagePath = "public/image/" + uniqueName

		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	_, err = config.DB.Exec(`
		INSERT INTO users (name, slug, username, password, image, about)
		VALUES (?, ?, ?, ?, ?, ?)`,
		name, slug, username, string(hashedPassword), imagePath, about)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created",
		"image":   imagePath,
	})
}
