package responses

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
	Id          uint64      `json:"id"`
	OwnerId     uint64      `json:"owner_id"`
	Email       string      `json:"email"`
	Name        string      `json:"name"`
	Phone       string      `json:"phone"`
	Description string      `json:"description"`
	Address     string      `json:"address"`
	LogoUrl     string      `json:"logo_url,omitempty"`
	IsVisible   bool        `json:"is_visible_for_clients"`
	Admins      []uint64    `json:"admin_ids"`
	WorkHours   []WorkHours `json:"work_hours"`
}

type WorkHours struct {
	Weekday uint   `json:"weekday"`
	From    string `json:"from"`
	To      string `json:"to"`
}

type AuthorizationPayload struct {
	Role         uint   `json:"role"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
