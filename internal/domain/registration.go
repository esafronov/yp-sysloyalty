package domain

type RegistrationRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegistrationReponse struct {
	AccessToken string `json:"accessToken"`
	//	RefreshToken string `json:"refreshToken"`
}
