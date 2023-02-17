package requests

type SuspendStudentRequest struct {
	StudentEmail string `json:"student" validate:"required,email"`
}
