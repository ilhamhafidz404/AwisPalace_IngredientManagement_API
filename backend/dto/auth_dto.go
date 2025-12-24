package dto

type GoogleAuthRequest struct {
	IDToken  string `json:"id_token" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name"`
	PhotoURL string `json:"photo_url"`
}

type AuthResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Data    AuthData `json:"data"`
}

type AuthData struct {
	Token string   `json:"token"`
	User  UserData `json:"user"`
}

type UserData struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	PhotoURL string `json:"photo_url"`
}
