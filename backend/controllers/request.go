package controllers

type userCreateRequest struct {
	User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Image string `json:"image"`
	} `json:"user"`
}

type userLinkAccountRequest struct {
	Provider          string `json:"provider"`
	Type              string `json:"type"`
	ProviderAccountId string `json:"providerAccountId"`
	AccessToken       string `json:"access_token"`
	TokenType         string `json:"token_type"`
	Scope             string `json:"scope"`
	UserId            string `json:"id"`
}
