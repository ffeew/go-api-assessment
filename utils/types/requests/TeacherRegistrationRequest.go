package requests

type TeacherRegistrationRequest struct {
	Email    string `json:"email" validate:"required,email,min=6,max=40"`
	Password string `json:"password" validate:"required,min=8,max=20"`
	Name     string `json:"name" validate:"required,min=3,max=40"`
	Age      uint8  `json:"age" validate:"required,min=18,max=100"`
}
