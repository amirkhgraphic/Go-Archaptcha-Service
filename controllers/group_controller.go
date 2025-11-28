package controllers

import (
	"net/http"
	"strings"

	"github.com/amirkhgraphic/go-arcaptcha-service/initializers"
	"github.com/amirkhgraphic/go-arcaptcha-service/models"
	"github.com/gin-gonic/gin"
)

// GroupUsers aggregates users by the specified fields (gender, nationality) directly from the DB.
// @Summary Group users
// @Description Aggregate users by gender/nationality
// @Param group_by query string false "Comma separated fields (gender,nationality)"
// @Success 200 {object} controllers.GroupUsersResponseDoc
// @Router /api/users/group [get]
func GroupUsers(c *gin.Context) {
	groupBy := parseGroupBy(c.DefaultQuery("group_by", "gender,nationality"))
	if len(groupBy) == 0 {
		groupBy = []string{"gender", "nationality"}
	}

	selectParts := append([]string{}, groupBy...)
	selectParts = append(selectParts, "count(*) as count")
	groupClause := strings.Join(groupBy, ", ")

	var rows []map[string]interface{}
	if err := initializers.DB.Model(&models.User{}).
		Select(strings.Join(selectParts, ", ")).
		Group(groupClause).
		Find(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not group users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"group_by": groupBy,
		"data":     rows,
	})
}

func parseGroupBy(raw string) []string {
	allowed := map[string]bool{
		"gender":      true,
		"nationality": true,
	}
	var out []string
	for _, part := range strings.Split(raw, ",") {
		p := strings.TrimSpace(part)
		if allowed[p] {
			out = append(out, p)
		}
	}
	// remove duplicates
	seen := map[string]bool{}
	var dedup []string
	for _, p := range out {
		if !seen[p] {
			seen[p] = true
			dedup = append(dedup, p)
		}
	}
	return dedup
}
