package domain

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginReponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
