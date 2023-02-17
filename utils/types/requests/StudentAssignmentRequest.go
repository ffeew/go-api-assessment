package requests

type StudentAssignmentRequest struct {
	TeacherEmail string   `json:"teacher" validate:"required,email"`
	StudentEmail []string `json:"students" validate:"required"`
}
