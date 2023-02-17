package requests

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,min=6,max=40"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}
