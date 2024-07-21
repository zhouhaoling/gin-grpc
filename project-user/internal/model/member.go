package model

// Member 数据库实体类
type Member struct {
	BaseModel
	Account       string   //账号
	Password      string   //密码
	Name          string   //用户昵称
	Mobile        string   //手机
	RealName      string   `gorm:"column:realname"` //真实姓名
	CreateTime    int64    //创建时间
	Status        int      //账号状态
	LastLoginTime int64    //最后登录时间
	Sex           int      //性别
	Avatar        string   //头像
	IdCard        string   `gorm:"column:idcard"` //身份证号
	Description   string   //备注
	Email         string   //邮箱
	DingTalk      DingTalk `gorm:"foreignkey:MId"`
	Address       Address  `gorm:"foreignkey:MId"`
}

func (m *Member) TableName() string {
	return "ms_member"
}

type DingTalk struct {
	BaseModel
	DingtalkOpenId  string `gorm:"column:dingtalk_openid"`  // 钉钉openid
	DingtalkUnionId string `gorm:"column:dingtalk_unionid"` // 钉钉unionid
	DingtalkUserid  string `gorm:"column:dingtalk_userid"`  // 钉钉userid
}

func (d *DingTalk) TableName() string {
	return "dingtalk"
}

type Address struct {
	BaseModel
	Province int    //省
	City     int    //市
	Area     int    //区
	Address  string //地址
}

func (a *Address) TableName() string {
	return "address"
}
