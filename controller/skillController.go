package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"my-personal-web/database"
	"my-personal-web/models"

	"github.com/gin-gonic/gin"
)

//  CREATE ======================

func CreateSkillsHandler(c *gin.Context) {
	cvIDStr := c.Param("cv_id")
	cvID, err := strconv.ParseUint(cvIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cv_id"})
		return
	}

	skillCatName := c.PostForm("skill_category_name")
	skillName := c.PostForm("skill_name")
	desc := c.PostForm("desc")
	nilaiStr := c.PostForm("nilai")

	var nilai int
	if nilaiStr != "" {
		nv, err := strconv.Atoi(nilaiStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "nilai must be an integer"})
			return
		}
		if nv < 0 || nv > 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "nilai must be between 0 and 100"})
			return
		}
		nilai = nv
	}

	imagePath := "public/image/default.jpg"
	file, err := c.FormFile("icon")
	if err == nil {
		ext := filepath.Ext(file.Filename)
		cleanName := strings.ReplaceAll(strings.ToLower(skillName), " ", "-")
		dateStr := time.Now().Format("20060102")
		uniqueName := fmt.Sprintf("%s-skill-icon-%s%s", cleanName, dateStr, ext)
		imagePath = "public/image/" + uniqueName

		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
	}

	skill := models.Skill{
		CvID:              uint64(cvID),
		SkillCategoryName: skillCatName,
		SkillName:         skillName,
		Desc:              desc,
		Icon:              imagePath,
		Nilai:             nilai,
	}

	if err := database.DB.Create(&skill).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create skill", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Skill created successfully",
		"skill":   skill,
	})
}

//  UPDATE ========================

func UpdateSkillsHandler(c *gin.Context) {
	skillIdStr := c.Param("skill_id")
	skillId, err := strconv.ParseUint(skillIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid skill_id"})
		return
	}

	var oldSkill models.Skill
	if err := database.DB.Where("skill_id = ?", skillId).First(&oldSkill).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}

	skillCatName := c.PostForm("skill_category_name")
	skillName := c.PostForm("skill_name")
	desc := c.PostForm("desc")
	nilaiStr := c.PostForm("nilai")

	if skillCatName != "" {
		oldSkill.SkillCategoryName = skillCatName
	}
	if skillName != "" {
		oldSkill.SkillName = skillName
	}
	if desc != "" {
		oldSkill.Desc = desc
	}
	if nilaiStr != "" {
		if n, err := strconv.Atoi(nilaiStr); err == nil {
			if n >= 0 && n <= 100 {
				oldSkill.Nilai = n
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "nilai must be between 0 and 100"})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "nilai must be a number"})
			return
		}
	}

	file, err := c.FormFile("icon")
	if err == nil {
		ext := filepath.Ext(file.Filename)

		baseName := oldSkill.SkillName
		if skillName != "" {
			baseName = skillName
		}

		cleanName := strings.ReplaceAll(strings.ToLower(baseName), " ", "-")
		dateStr := time.Now().Format("20060102")

		uniqueName := fmt.Sprintf("%s-skill-icon-%s%s", cleanName, dateStr, ext)
		newImagePath := "public/image/" + uniqueName

		if oldSkill.Icon != "public/image/default.jpg" {
			_ = os.Remove(oldSkill.Icon)
		}

		if err := c.SaveUploadedFile(file, newImagePath); err == nil {
			oldSkill.Icon = newImagePath
		}
	}

	if err := database.DB.Save(&oldSkill).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update skill"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Skill updated successfully",
		"skill":   oldSkill,
	})
}

//  GET ========================

// Get Hard Skills
func GetHardSkillsHandler(c *gin.Context) {
	getSkillsByCategory(c, "Hard Skill")
}

// Get Soft Skills
func GetSoftSkillsHandler(c *gin.Context) {
	getSkillsByCategory(c, "Soft Skill")
}

// Get Tools Skills
func GetToolsSkillsHandler(c *gin.Context) {
	getSkillsByCategory(c, "Tools Skill")
}

// Helper function untuk mengurangi duplikasi kode
func getSkillsByCategory(c *gin.Context, category string) {
	cvIDStr := c.Param("cv_id")
	cvID, err := strconv.ParseUint(cvIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cv_id"})
		return
	}

	var skills []models.Skill
	if err := database.DB.
		Where("cv_id = ? AND skill_category_name = ?", cvID, category).
		Find(&skills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch skills"})
		return
	}

	if len(skills) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("No %s found for this CV", category)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cv_id":    cvID,
		"category": category,
		"skills":   skills,
	})
}

// Delete ======================

func DeleteSkillHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid skill id"})
		return
	}

	var skill models.Skill
	if err := database.DB.First(&skill, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}

	// Hapus file image kalau bukan default
	if skill.Icon != "" && skill.Icon != "public/image/default.jpg" {
		if err := os.Remove(skill.Icon); err != nil && !os.IsNotExist(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete skill icon"})
			return
		}
	}

	// Hapus record dari DB
	if err := database.DB.Delete(&skill).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete skill"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Skill deleted successfully",
		"skill":   skill,
	})
}
