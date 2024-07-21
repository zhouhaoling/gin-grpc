package model

type LoginResp struct {
	Member           Member             `json:"member"`
	TokenList        TokenList          `json:"tokenList"`
	OrganizationList []OrganizationList `json:"organizationList"`
}

type Member struct {
	//Id int64  `json:"id"`
	Name             string `json:"name"`
	Mobile           string `json:"mobile"`
	Status           int    `json:"status"`
	Code             string `json:"code"` //加密后的mid
	Email            string `json:"email"`
	CreateTime       string `json:"create_time"`
	LastLoginTime    string `json:"last_login_time"`
	OrganizationCode string `json:"organization_code"`
	Avatar           string `json:"avatar"`
}

type TokenList struct {
	AccessToken    string `json:"accessToken"`
	RefreshToken   string `json:"refreshToken"`
	TokenType      string `json:"tokenType"`
	AccessTokenExp int64  `json:"accessTokenExp"`
}

type OrganizationList struct {
	//Id          int64  `json:"id"`
	Code        string `json:"code"` //加密后的organization表的id
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	MemberId    int64  `json:"member_id"`
	OwnerCode   string `json:"owner_code"`
	//Mbid       string `json:"member_id"`
	CreateTime int64  `json:"create_time"`
	Personal   int32  `json:"personal"`
	Address    string `json:"address"`
	Province   int32  `json:"province"`
	City       int32  `json:"city"`
	Area       int32  `json:"area"`
}
