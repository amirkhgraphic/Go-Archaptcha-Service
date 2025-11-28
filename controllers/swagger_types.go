package controllers

// UserDoc is a Swagger-only representation of User (avoids gorm.Model embedding).
type UserDoc struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Bio         string `json:"bio"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

type PaginationDoc struct {
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalItems int64             `json:"total_items"`
	TotalPages int               `json:"total_pages"`
	Sort       string            `json:"sort"`
	Search     string            `json:"search,omitempty"`
	Filters    map[string]string `json:"filters,omitempty"`
}

type UserListResponseDoc struct {
	Data []UserDoc     `json:"data"`
	Meta PaginationDoc `json:"meta"`
}

type GroupUsersResponseDoc struct {
	GroupBy []string                 `json:"group_by"`
	Data    []map[string]interface{} `json:"data"`
}

type ChallengeVerifyRequest struct {
    ChallengeID string `json:"challenge_id"`
}
type ChallengeVerifyResponse struct {
    Valid bool   `json:"valid"`
    Error string `json:"error,omitempty"`
}

type ChallengeResponse struct {
    ChallengeID string `json:"challenge_id"`
    Note        string `json:"note,omitempty"`
}

type ErrorResponse struct {
    Error   string `json:"error"`
    Details string `json:"details,omitempty"`
}
