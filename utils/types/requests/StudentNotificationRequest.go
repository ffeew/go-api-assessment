package requests

type StudentNotificationRequest struct {
	TeacherEmail string `json:"teacher" validate:"required,email"`
	Notification string `json:"notification" validate:"required"`
}
