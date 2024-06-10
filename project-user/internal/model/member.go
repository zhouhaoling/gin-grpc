package model

// Member 数据库实体类
type Member struct {
	BaseModel
	Account       string   `json:"account"`                         //账号
	Password      string   `json:"password"`                        //密码
	Name          string   `json:"name"`                            //用户昵称
	Mobile        string   `json:"mobile"`                          //手机
	RealName      string   `json:"realname" gorm:"column:realname"` //真实姓名
	CreateTime    int64    `json:"create_time"`                     //创建时间
	Status        int      `json:"status"`                          //账号状态
	LastLoginTime int64    `json:"last_login_time"`                 //最后登录时间
	Sex           int      `json:"sex"`                             //性别
	Avatar        string   `json:"avatar"`                          //头像
	IdCard        string   `json:"id_card" gorm:"column:idcard"`    //身份证号
	Description   string   `json:"description"`                     //备注
	Email         string   `json:"email"`                           //邮箱
	DingTalk      DingTalk `gorm:"foreignkey:MId"`
	Address       Address  `gorm:"foreignkey:MId"`
}

func (m *Member) TableName() string {
	return "ms_member"
}

type DingTalk struct {
	BaseModel
	DingtalkOpenId  string `json:"dingtalk_open_id" gorm:"column:dingtalk_openid"`   // 钉钉openid
	DingtalkUnionId string `json:"dingtalk_union_id" gorm:"column:dingtalk_unionid"` // 钉钉unionid
	DingtalkUserid  string `json:"dingtalk_user_id" gorm:"column:dingtalk_userid"`   // 钉钉userid
}

func (d *DingTalk) TableName() string {
	return "dingtalk"
}

type Address struct {
	BaseModel
	Province int    `json:"province"` //省
	City     int    `json:"city"`     //市
	Area     int    `json:"area"`     //区
	Address  string `json:"address"`  //地址
}

func (a *Address) TableName() string {
	return "address"
}
