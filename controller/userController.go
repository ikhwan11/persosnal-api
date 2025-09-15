package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"my-personal-web/api/config"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Username string `json:"username"`
	Password string `json:"password"`
	Image    string `json:"image"`
	About    string `json:"about"`
}

// =========== CREATE

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

		// generate nama file: nama-user-image-YYYYMMDD.ext
		cleanName := strings.ToLower(name)
		cleanName = strings.ReplaceAll(cleanName, " ", "-")
		dateStr := time.Now().Format("20060102") // YYYYMMDD

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

	// Simpan ke DB
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

// ======= GET

// GET all users
func GetUsersHandler(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT id, name, slug, username, image, about 
		FROM users
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Slug, &u.Username, &u.Image, &u.About); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user"})
			return
		}
		users = append(users, u)
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

// GET user by slug
func GetUserBySlugHandler(c *gin.Context) {
	slug := c.Param("slug")

	var u User
	err := config.DB.QueryRow(`
		SELECT id, name, slug, username, image, about 
		FROM users 
		WHERE slug = ?`, slug).
		Scan(&u.ID, &u.Name, &u.Slug, &u.Username, &u.Image, &u.About)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		}
		return
	}

	c.JSON(http.StatusOK, u)
}

// =========== PUT

func UpdateUserHandler(c *gin.Context) {
	slug := c.Param("slug")

	// Ambil user lama dulu
	var oldUser User
	err := config.DB.QueryRow(`
		SELECT id, name, slug, username, password, image, about 
		FROM users WHERE slug = ?`, slug).
		Scan(&oldUser.ID, &oldUser.Name, &oldUser.Slug, &oldUser.Username, &oldUser.Password, &oldUser.Image, &oldUser.About)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		}
		return
	}

	// Ambil form data
	name := c.PostForm("name")
	username := c.PostForm("username")
	password := c.PostForm("password")
	about := c.PostForm("about")

	// Update slug kalau nama diubah
	newSlug := oldUser.Slug
	if name != "" {
		newSlug = generateSlug(name)
	} else {
		name = oldUser.Name
	}

	// Handle image upload
	imagePath := oldUser.Image
	file, err := c.FormFile("image")
	if err == nil {
		ext := filepath.Ext(file.Filename)

		cleanName := strings.ToLower(name)
		cleanName = strings.ReplaceAll(cleanName, " ", "-")
		dateStr := time.Now().Format("20060102")

		uniqueName := fmt.Sprintf("%s-user-image-%s%s", cleanName, dateStr, ext)
		newImagePath := "public/image/" + uniqueName

		// Hapus foto lama kalau bukan default
		if oldUser.Image != "public/image/default.jpg" {
			if removeErr := os.Remove(oldUser.Image); removeErr != nil && !os.IsNotExist(removeErr) {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove old image"})
				return
			}
		}

		// Simpan foto baru
		if err := c.SaveUploadedFile(file, newImagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}

		imagePath = newImagePath
	}

	// Handle password update (hanya kalau ada input baru)
	hashedPassword := oldUser.Password
	if password != "" {
		newPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		hashedPassword = string(newPass)
	}

	// Update ke DB
	_, err = config.DB.Exec(`
		UPDATE users
		SET name = ?, slug = ?, username = ?, password = ?, image = ?, about = ?
		WHERE slug = ?`,
		name, newSlug, username, hashedPassword, imagePath, about, slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated",
		"slug":    newSlug,
		"image":   imagePath,
	})
}
