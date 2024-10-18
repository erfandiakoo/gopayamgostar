package gopayamgostar

type JWT struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	//ExpiresAt    time.Time `json:"expiresAt"`
}
