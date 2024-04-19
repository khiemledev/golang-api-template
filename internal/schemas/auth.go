package schemas

type AuthLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthLoginUserResponse struct {
	ID       uint   `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

type AuthLoginResponse struct {
	Status  int                   `json:"status" binding:"required"`
	Message string                `json:"message" binding:"required"`
	User    AuthLoginUserResponse `json:"user"`
}

type AuthRegisterRequest struct {
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Name            string `json:"name" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type AuthRegisterResponse struct {
	Status        int    `json:"status" binding:"required"`
	Message       string `json:"message" binding:"required"`
	CreatedUserId uint   `json:"created_user_id" binding:"required"`
}
