package controller

import (
	"net/http"
	"strconv"

	"my-personal-web/database"
	"my-personal-web/models"

	"github.com/gin-gonic/gin"
)

//  CREATE ======================

func CreateEducationsHandler(c *gin.Context) {
	cvIDStr := c.Param("cv_id")
	cvID, err := strconv.ParseUint(cvIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cv_id"})
		return
	}

	universitasName := c.PostForm("universitas_name")
	jurusan := c.PostForm("jurusan")
	tahun := c.PostForm("tahun")
	desc := c.PostForm("desc")
	ipkStr := c.PostForm("ipk")

	var ipk float64
	if ipkStr != "" {
		fv, err := strconv.ParseFloat(ipkStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ipk must be a number (e.g. 3.75)"})
			return
		}
		if fv < 0 || fv > 4 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ipk must be between 0.00 and 4.00"})
			return
		}
		ipk = fv
	}

	education := models.Education{
		CvID:            uint64(cvID),
		UniversitasName: universitasName,
		Jurusan:         jurusan,
		Tahun:           tahun,
		Desc:            desc,
		IPK:             ipk,
	}

	if err := database.DB.Create(&education).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to create education",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Education created successfully",
		"education": education,
	})
}

//  UPDATE ========================

func UpdateEducationsHandler(c *gin.Context) {
	eduIdStr := c.Param("edu_id")
	eduId, err := strconv.ParseUint(eduIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid edu_id"})
		return
	}

	var oldEdu models.Education
	if err := database.DB.Where("edu_id = ?", eduId).First(&oldEdu).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Education not found"})
		return
	}

	universitasName := c.PostForm("universitas_name")
	jurusan := c.PostForm("jurusan")
	tahun := c.PostForm("tahun")
	desc := c.PostForm("desc")
	ipkStr := c.PostForm("ipk")

	if universitasName != "" {
		oldEdu.UniversitasName = universitasName
	}
	if jurusan != "" {
		oldEdu.Jurusan = jurusan
	}
	if tahun != "" {
		oldEdu.Tahun = tahun
	}
	if desc != "" {
		oldEdu.Desc = desc
	}
	if ipkStr != "" {
		if fv, err := strconv.ParseFloat(ipkStr, 64); err == nil {
			if fv >= 0 && fv <= 4 {
				oldEdu.IPK = fv
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ipk must be between 0.00 and 4.00"})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ipk must be a number (e.g. 3.75)"})
			return
		}
	}

	if err := database.DB.Save(&oldEdu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update education"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Education updated successfully",
		"education": oldEdu,
	})
}

//  GET ========================

func GetEducationsHandler(c *gin.Context) {
	cvIDStr := c.Param("cv_id")
	cvID, err := strconv.ParseUint(cvIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cv_id"})
		return
	}

	var education []models.Education
	if err := database.DB.Where("cv_id = ?", cvID).Find(&education).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch education"})
		return
	}

	if len(education) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No education found for this CV"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cv_id":     cvID,
		"education": education,
	})
}

//  Delete ==================

func DeleteEducationHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid education id"})
		return
	}

	var edu models.Education
	if err := database.DB.First(&edu, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Education not found"})
		return
	}

	if err := database.DB.Delete(&edu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete education"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Education deleted successfully",
		"education": edu,
	})
}
