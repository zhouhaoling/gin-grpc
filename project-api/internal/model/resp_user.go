package model

type LoginResp struct {
	Member           Member             `json:"member"`
	TokenList        TokenList          `json:"token_list"`
	OrganizationList []OrganizationList `json:"organization_list"`
}

type Member struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
	Status int    `json:"status"`
}

type TokenList struct {
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
	TokenType      string `json:"token_type"`
	AccessTokenExp int64  `json:"access_token_exp"`
}

type OrganizationList struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	MemberId    int64  `json:"member_id"`
	CreateTime  int64  `json:"create_time"`
	Personal    int32  `json:"personal"`
	Address     string `json:"address"`
	Province    int32  `json:"province"`
	City        int32  `json:"city"`
	Area        int32  `json:"area"`
}
