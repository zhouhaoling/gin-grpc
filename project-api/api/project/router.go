package project

import (
	"log"

	"test.com/project-api/api/middleware"

	"github.com/gin-gonic/gin"
	"test.com/project-api/router"
)

func init() {
	log.Println("init project router")
	router.Register(NewRouterProject())
}

type RouterProject struct {
}

func NewRouterProject() *RouterProject {
	return &RouterProject{}
}

// Route 项目路由模块
/*
 * @Description:存放用户相关路由
 */
func (ru *RouterProject) Route(r *gin.Engine) {
	InitGrpcProjectClient()
	p := NewHandlerProject()
	pg := r.Group("/project")
	pg.Use(middleware.NewJwtMiddlewareBuilder().Build())
	//pg.Use(Auth())
	pg.Use(ProjectsAuth())
	pg.POST("/index", p.index)
	pg.POST("/project/selfList", p.myProjectList)
	pg.POST("/project", p.myProjectList)
	pg.POST("/project_template", p.projectTemplate)
	pg.POST("/project/save", p.createProject)
	pg.POST("/project/read", p.readProject)
	pg.POST("/project/recycle", p.recycleProject)
	pg.POST("/project/recovery", p.recoveryProject)
	pg.POST("/project_collect/collect", p.collectProject)
	pg.POST("/project/edit", p.editProject)
	pg.POST("/project/getLogBySelfProject", p.getLogBySelfProject)
	pg.POST("/node", p.nodeList)

	t := NewHandlerTask()
	pg.POST("/task_stages", t.taskStages)
	pg.POST("/project_member/index", t.taskMemberList)
	pg.POST("/task_stages/tasks", t.taskList) //任务步骤详情
	pg.POST("/task/save", t.taskSave)
	pg.POST("/task/sort", t.taskSort)
	pg.POST("/task/selfList", t.myTaskList)
	pg.POST("/task/read", t.taskRead)
	pg.POST("/task_member", t.listTaskMember)
	pg.POST("/task/taskLog", t.taskLog)
	pg.POST("/task/_taskWorkTimeList", t.taskWorkTimeList)
	pg.POST("/task/saveTaskWorkTime", t.saveTaskWorkTime)
	pg.POST("/task/taskSources", t.taskSources)
	pg.POST("/task/createComment", t.taskCreateComment)
	pg.POST("/task_stages/save", t.taskStagesSave)
	pg.POST("/task_stages/delete", t.taskStagesDelete)
	pg.POST("/task_stages/edit", t.taskStagesEdit)

	f := NewHandlerFile()
	pg.POST("/file/uploadFiles", f.uploadFiles)

	a := NewHandlerAccount()
	pg.POST("/account", a.account)

	d := NewHandlerDepartment()
	pg.POST("/department", d.department)
	pg.POST("/department/save", d.saveDepartment)
	pg.POST("/department/read", d.readDepartment)

	auth := NewHandlerAuth()
	pg.POST("/auth", auth.authList)
	pg.POST("/auth/apply", auth.authApply)

	m := NewHandlerMenu()
	pg.POST("/menu/menu", m.menuList)
}
