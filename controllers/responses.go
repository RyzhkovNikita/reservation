package controllers

type ProfileResponse struct {
	Id         uint64 `json:"id"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Role       uint   `json:"role"`
}

type BarInfoResponse struct {
	Id          uint64 `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Address     string `json:"address"`
	LogoUrl     string `json:"logo_url,omitempty"`
}

type AuthorizationPayload struct {
	Role         uint   `json:"role"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
