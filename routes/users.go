package routes

import (
	"net/http"
	"strconv"
	"godoor/models"
	"godoor/serializers"
	"github.com/gin-gonic/gin"
)

// GetUsers returns list of users in current space
func GetUsers(c *gin.Context) {
	space := c.MustGet("space").(*models.Space)

	var userSpaces []models.UserSpace
	if err := models.DB.Where("space_id = ?", space.ID).Preload("User").Find(&userSpaces).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	users := make([]models.User, len(userSpaces))
	for i, us := range userSpaces {
		users[i] = us.User
	}

	c.JSON(http.StatusOK, serializers.SerializeUsers(users))
}

// GetUser returns a specific user
func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	space := c.MustGet("space").(*models.Space)

	var userSpace models.UserSpace
	if err := models.DB.Where("space_id = ? AND user_id = ?", space.ID, uint(id)).Preload("User").First(&userSpace).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, serializers.SerializeUser(&userSpace.User))
}

// UpdateUser updates a user (only first_name and last_name can be updated)
type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	space := c.MustGet("space").(*models.Space)
	currentUser := c.MustGet("user").(*models.User)

	// Check if user can update (only themselves or admin)
	if currentUser.ID != uint(id) && !currentUser.IsSuperuser {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	var user models.User
	if err := models.DB.First(&user, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if user is in current space
	var userSpace models.UserSpace
	if err := models.DB.Where("space_id = ? AND user_id = ?", space.ID, uint(id)).First(&userSpace).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not in current space"})
		return
	}

	// Update user
	user.FirstName = req.FirstName
	user.LastName = req.LastName

	if err := models.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, serializers.SerializeUser(&user))
}
