package schemas

type APIResponse struct {
	Status  int    `json:"status" binding:"required"`
	Message string `json:"message" binding:"required"`
	Data    any    `json:"data"`
}
