package request

// LoginByEmail is struct for login by email.
type LoginByEmail struct {
	Email    string `json:"email" binding:"email,required"`
	Password string `json:"password" binding:"required"`
}
