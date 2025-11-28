package controllers

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/amirkhgraphic/go-arcaptcha-service/initializers"
	"github.com/amirkhgraphic/go-arcaptcha-service/models"
	"github.com/amirkhgraphic/go-arcaptcha-service/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type createUserRequest struct {
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Bio         string `json:"bio"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
	ChallengeID string `json:"challenge_id" binding:"required"`
}

type updateUserRequest struct {
	Username    *string `json:"username,omitempty"`
	Email       *string `json:"email,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	Nationality *string `json:"nationality,omitempty"`
	ChallengeID string  `json:"challenge_id" binding:"required"`
}

type userListResponse struct {
	Data []models.User `json:"data"`
	Meta pagination    `json:"meta"`
}

type pagination struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalItems int64  `json:"total_items"`
	TotalPages int    `json:"total_pages"`
	Sort       string `json:"sort"`
	Search     string `json:"search,omitempty"`
	Filters    gin.H  `json:"filters,omitempty"`
}

// CreateUser creates a new user after captcha validation.
// @Summary Create user
// @Accept json
// @Produce json
// @Param payload body createUserRequest true "User payload"
// @Success 201 {object} controllers.UserDoc
// @Failure 400 {object} controllers.ErrorResponse
// @Router /api/users [post]
func CreateUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "details": err.Error()})
		return
	}

	if err := services.Arcaptcha.ValidateChallenge(req.ChallengeID); err != nil {
		respondCaptchaError(c, err)
		return
	}

	user := models.User{
		Username:    req.Username,
		Email:       req.Email,
		Bio:         req.Bio,
		Gender:      strings.TrimSpace(req.Gender),
		Nationality: strings.TrimSpace(req.Nationality),
	}

	if err := initializers.DB.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "unique") {
			c.JSON(http.StatusConflict, gin.H{"error": "username or email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

// ListUsers returns paginated users with search/filter options.
// @Summary List users
// @Produce json
// @Param page query int false "page"
// @Param page_size query int false "page size"
// @Param sort query string false "sort (e.g. -created_at)"
// @Param search query string false "search in username/email"
// @Param username query string false "filter by username"
// @Param email query string false "filter by email"
// @Success 200 {object} controllers.UserListResponseDoc
// @Router /api/users [get]
func ListUsers(c *gin.Context) {
	page := parsePositiveInt(c.DefaultQuery("page", "1"), 1)
	pageSize := parsePositiveInt(c.DefaultQuery("page_size", "10"), 10)
	if pageSize > 100 {
		pageSize = 100
	}

	sort := sanitizeSort(c.DefaultQuery("sort", "-created_at"))
	search := strings.TrimSpace(c.Query("search"))
	usernameFilter := strings.TrimSpace(c.Query("username"))
	emailFilter := strings.TrimSpace(c.Query("email"))

	var users []models.User
	tx := initializers.DB.Model(&models.User{})

	if search != "" {
		searchValue := "%" + strings.ToLower(search) + "%"
		tx = tx.Where("LOWER(username) LIKE ? OR LOWER(email) LIKE ?", searchValue, searchValue)
	}
	if usernameFilter != "" {
		tx = tx.Where("username = ?", usernameFilter)
	}
	if emailFilter != "" {
		tx = tx.Where("email = ?", emailFilter)
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not count users"})
		return
	}

	offset := (page - 1) * pageSize
	if err := tx.Order(sort).Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch users"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages == 0 {
		totalPages = 1
	}

	c.JSON(http.StatusOK, userListResponse{
		Data: users,
		Meta: pagination{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: total,
			TotalPages: totalPages,
			Sort:       sort,
			Search:     search,
			Filters: gin.H{
				"username": usernameFilter,
				"email":    emailFilter,
			},
		},
	})
}

// GetUser fetches a user by id.
// @Summary Get user
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} controllers.UserDoc
// @Failure 400 {object} controllers.ErrorResponse
// @Router /api/users/{id} [get]
func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := initializers.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// UpdateUser updates a user after captcha validation.
// @Summary Update user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param payload body updateUserRequest true "Fields to update"
// @Success 200 {object} controllers.UserDoc
// @Failure 400 {object} controllers.ErrorResponse
// @Router /api/users/{id} [patch]
func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "details": err.Error()})
		return
	}

	if req.Username == nil && req.Email == nil && req.Bio == nil && req.Gender == nil && req.Nationality == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nothing to update"})
		return
	}

	if err := services.Arcaptcha.ValidateChallenge(req.ChallengeID); err != nil {
		respondCaptchaError(c, err)
		return
	}

	var user models.User
	if err := initializers.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch user"})
		return
	}

	updates := make(map[string]interface{})
	if req.Username != nil {
		updates["username"] = strings.TrimSpace(*req.Username)
	}
	if req.Email != nil {
		updates["email"] = strings.TrimSpace(*req.Email)
	}
	if req.Bio != nil {
		updates["bio"] = *req.Bio
	}
	if req.Gender != nil {
		updates["gender"] = strings.TrimSpace(*req.Gender)
	}
	if req.Nationality != nil {
		updates["nationality"] = strings.TrimSpace(*req.Nationality)
	}

	if err := initializers.DB.Model(&user).Updates(updates).Error; err != nil {
		if strings.Contains(err.Error(), "unique") {
			c.JSON(http.StatusConflict, gin.H{"error": "username or email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func respondCaptchaError(c *gin.Context, err error) {
	switch err {
	case services.ErrChallengeEmpty:
		c.JSON(http.StatusBadRequest, gin.H{"error": "challenge_id is required"})
	case services.ErrChallengeInvalid:
		c.JSON(http.StatusBadRequest, gin.H{"error": "captcha did not match the issued challenge"})
	case services.ErrChallengeNetwork:
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "captcha provider unavailable, try again"})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "captcha validation failed"})
	}
}

func sanitizeSort(raw string) string {
	allowed := map[string]bool{
		"username":   true,
		"email":      true,
		"created_at": true,
		"updated_at": true,
	}

	if raw == "" {
		return "created_at desc"
	}

	dir := "asc"
	field := raw
	if strings.HasPrefix(raw, "-") {
		dir = "desc"
		field = strings.TrimPrefix(raw, "-")
	} else if strings.HasPrefix(raw, "+") {
		field = strings.TrimPrefix(raw, "+")
	}

	if !allowed[field] {
		field = "created_at"
	}

	return field + " " + dir
}

func parsePositiveInt(value string, defaultVal int) int {
	n, err := strconv.Atoi(value)
	if err != nil || n <= 0 {
		return defaultVal
	}
	return n
}
