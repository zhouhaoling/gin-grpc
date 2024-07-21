package model

import (
	"strings"

	"github.com/jinzhu/copier"
	"test.com/common/encrypts"
	"test.com/common/tms"
)

type ProjectAuth struct {
	Id               int64  `json:"id"`
	OrganizationCode int64  `json:"organization_code"`
	Title            string `json:"title"`
	CreateAt         int64  `json:"create_at"`
	Sort             int    `json:"sort"`
	Status           int    `json:"status"`
	Desc             string `json:"desc"`
	CreateBy         int64  `json:"create_by"`
	IsDefault        int    `json:"is_default"`
	Type             string `json:"type"`
}

func (*ProjectAuth) TableName() string {
	return "ms_project_auth"
}

func (a *ProjectAuth) ToDisplay() *ProjectAuthDisplay {
	p := &ProjectAuthDisplay{}
	copier.Copy(p, a)
	p.OrganizationCode = encrypts.EncryptInt64NoErr(a.OrganizationCode)
	p.CreateAt = tms.FormatByMill(a.CreateAt)
	if a.Type == "admin" || a.Type == "member" {
		p.CanDelete = 0
	} else {
		p.CanDelete = 1
	}
	return p
}

type ProjectAuthDisplay struct {
	Id               int64  `json:"id"`
	OrganizationCode string `json:"organization_code"`
	Title            string `json:"title"`
	CreateAt         string `json:"create_at"`
	Sort             int    `json:"sort"`
	Status           int    `json:"status"`
	Desc             string `json:"desc"`
	CreateBy         int64  `json:"create_by"`
	IsDefault        int    `json:"is_default"`
	Type             string `json:"type"`
	CanDelete        int    `json:"canDelete"`
}

type ProjectNodeAuthTree struct {
	Id       int64
	Node     string
	Title    string
	IsMenu   int
	IsLogin  int
	IsAuth   int
	Pnode    string
	Key      string
	Checked  bool
	Children []*ProjectNodeAuthTree
}

func ToAuthNodeTreeList(list []*ProjectNode, checkedList []string) []*ProjectNodeAuthTree {
	checkedMap := make(map[string]struct{})
	for _, v := range checkedList {
		checkedMap[v] = struct{}{}
	}
	var roots []*ProjectNodeAuthTree
	for _, v := range list {
		paths := strings.Split(v.Node, "/")
		if len(paths) == 1 {
			checked := false
			if _, ok := checkedMap[v.Node]; ok {
				checked = true
			}
			//根节点
			root := &ProjectNodeAuthTree{
				Id:       v.Id,
				Node:     v.Node,
				Pnode:    "",
				IsLogin:  v.IsLogin,
				IsMenu:   v.IsMenu,
				IsAuth:   v.IsAuth,
				Title:    v.Title,
				Children: []*ProjectNodeAuthTree{},
				Checked:  checked,
				Key:      v.Node,
			}
			roots = append(roots, root)
		}
	}
	for _, v := range roots {
		addAuthNodeChild(list, v, 2, checkedMap)
	}
	return roots
}

func addAuthNodeChild(list []*ProjectNode, root *ProjectNodeAuthTree, level int, checkedMap map[string]struct{}) {
	for _, v := range list {
		if strings.HasPrefix(v.Node, root.Node+"/") && len(strings.Split(v.Node, "/")) == level {
			//此根节点子节点
			checked := false
			if _, ok := checkedMap[v.Node]; ok {
				checked = true
			}

			child := &ProjectNodeAuthTree{
				Id:       v.Id,
				Node:     v.Node,
				Pnode:    "",
				IsLogin:  v.IsLogin,
				IsMenu:   v.IsMenu,
				IsAuth:   v.IsAuth,
				Title:    v.Title,
				Children: []*ProjectNodeAuthTree{},
				Checked:  checked,
				Key:      v.Node,
			}
			root.Children = append(root.Children, child)
		}
	}
	for _, v := range root.Children {
		addAuthNodeChild(list, v, level+1, checkedMap)
	}
}

type ProjectAuthNode struct {
	Id   int64
	Auth int64
	Node string
}

func (p *ProjectAuthNode) TableName() string {
	return "ms_project_auth_node"
}

func NewProjectAuthNode(list []*ProjectAuthNode, auth int64) []*ProjectAuthNode {
	res := make([]*ProjectAuthNode, len(list))
	for i, v := range list {
		res[i] = &ProjectAuthNode{
			Auth: auth,
			Node: v.Node,
		}
	}
	return res
}
