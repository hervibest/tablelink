package converter

import (
	"tablelink/internal/domain"
	"tablelink/internal/model"
)

func UserToResponse(user *domain.User) *model.UserResponse {
	return &model.UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		RoleID:     user.RoleID,
		LastAccess: user.LastAccess,
	}
}

func UsersToResponses(users []*domain.User) []*model.UserResponse {
	responses := make([]*model.UserResponse, 0)
	for _, user := range users {
		response := UserToResponse(user)
		responses = append(responses, response)
	}
	return responses
}
