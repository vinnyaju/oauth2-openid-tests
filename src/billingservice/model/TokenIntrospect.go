package model

type TokenIntrospect struct {
	Exp               int            `json:"exp"`
	Iat               int            `json:"iat"`
	AuthTime          int            `json:"auth_time"`
	Jti               string         `json:"jti"`
	Iss               string         `json:"iss"`
	Aud               string         `json:"aud"`
	Sub               string         `json:"sub"`
	Typ               string         `json:"typ"`
	Azp               string         `json:"azp"`
	SessionState      string         `json:"session_state"`
	Name              string         `json:"name"`
	GivenName         string         `json:"given_name"`
	FamilyName        string         `json:"family_name"`
	PreferredUsername string         `json:"preferred_username"`
	Email             string         `json:"email"`
	EmailVerified     bool           `json:"email_verified"`
	Acr               string         `json:"acr"`
	AllowedOrigins    []string       `json:"allowed-origins"`
	RealmAccess       RealmAccess    `json:"realm_access"`
	ResourceAccess    ResourceAccess `json:"resource_access"`
	Scope             string         `json:"scope"`
	ClientID          string         `json:"client_id"`
	Username          string         `json:"username"`
	Active            bool           `json:"active"`
}
