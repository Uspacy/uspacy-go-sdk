package auth

const VersionUrl = "auth/v1/"

const (
	RefreshTokenUrl = "auth/refresh_token"
)

type RefreshOutput struct {
	Jwt             string `json:"jwt"`
	RefreshToken    string `json:"refreshToken"`
	ExpireInSeconds int64  `json:"expireInSeconds"`
}

type JwtClaims struct {
	Iss         string   `json:"iss"`
	Iat         int64    `json:"iat"`
	Exp         int64    `json:"exp"`
	Nbf         int64    `json:"nbf"`
	Jti         string   `json:"jti"`
	Id          int64    `json:"id"`
	Email       string   `json:"email"`
	Sub         string   `json:"sub"`
	AuthUserId  int      `json:"authUserId"`
	FirstName   string   `json:"firstName"`
	LastName    string   `json:"lastName"`
	WorkspaceId string   `json:"workspaceId"`
	Domain      string   `json:"domain"`
	Roles       []string `json:"roles"`
	Departments []int    `json:"departments"`
	Permissions struct {
		Create []string `json:"create"`
		View   []string `json:"view"`
		Edit   []string `json:"edit"`
		Delete []string `json:"delete"`
	} `json:"permissions"`
	TariffId int `json:"tariffId"`
}
