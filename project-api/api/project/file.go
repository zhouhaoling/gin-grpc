package project

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"test.com/common/encrypts"

	"test.com/common/errs"
	"test.com/project-api/config"
	pg "test.com/project-grpc/project_grpc"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"test.com/common"
	"test.com/common/fs"
	"test.com/common/tms"
	"test.com/project-api/internal/model"
)

type HandlerFile struct {
}

func NewHandlerFile() *HandlerFile {
	return &HandlerFile{}
}

// uploadFiles 上传文件
func (f *HandlerFile) uploadFiles(c *gin.Context) {
	res := common.NewResponseData()
	req := model.UploadFileReq{}
	c.ShouldBind(&req)

	fmt.Println("req:", req)
	fmt.Println("req.TaskCode:", encrypts.DecryptInt64(req.TaskCode), "req.Name", req.Filename, "req.ProjectCode", encrypts.DecryptInt64(req.ProjectCode))
	//处理文件
	multipartForm, _ := c.MultipartForm()
	file := multipartForm.File
	key := ""
	//只上传一个文件
	uploadFile := file["file"][0]
	//第一种不分片上传文件,上传到服务器主机
	if req.TotalChunks == 1 {
		path := "upload/" + req.ProjectCode + "/" + req.TaskCode + "/" + tms.FormatYMD(time.Now())
		if !fs.IsExist(path) {
			os.MkdirAll(path, os.ModePerm)
		}
		dst := path + "/" + req.Filename
		key = dst
		err := c.SaveUploadedFile(uploadFile, dst)
		if err != nil {
			zap.L().Error("保存文件失败", zap.Error(err))
			res.ResponseError(c, common.CodeServerBusy)
			return
		}
	}
	if req.TotalChunks > 1 {
		//分片上传 合起来即可
		path := "upload/" + req.ProjectCode + "/" + req.TaskCode + "/" + tms.FormatYMD(time.Now())
		if !fs.IsExist(path) {
			os.MkdirAll(path, os.ModePerm)
		}
		fileName := path + "/" + req.Filename
		openFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			zap.L().Error("打开文件失败", zap.Error(err))
			res.ResponseError(c, common.CodeServerBusy)
			return
		}
		open, err := uploadFile.Open()
		if err != nil {
			zap.L().Error("打开文件失败", zap.Error(err))
			res.ResponseError(c, common.CodeServerBusy)
			return
		}
		defer open.Close()
		buf := make([]byte, req.CurrentChunkSize)
		open.Read(buf)
		openFile.Write(buf)
		openFile.Close()

		if req.TotalChunks == req.ChunkNumber {
			//最后一个分片,对文件改名
			newPath := path + "/" + req.Filename
			key = newPath
			os.Rename(fileName, newPath)
		}
	}

	//调用grpc服务，存储文件信息和资源路径
	//调用服务 存入file表
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	fileUrl := "http://localhost/" + key
	msg := &pg.TaskFileReqMessage{
		TaskCode:         req.TaskCode,
		ProjectCode:      req.ProjectCode,
		OrganizationCode: c.GetString(config.CtxOrganizationIDKey),
		PathName:         key,
		FileName:         req.Filename,
		Size:             int64(req.TotalSize),
		Extension:        path.Ext(key),
		FileUrl:          fileUrl,
		FileType:         file["file"][0].Header.Get("Content-Type"),
		MemberId:         c.GetInt64(config.CtxMemberIDKey),
	}
	if req.TotalChunks == req.ChunkNumber {
		_, err := taskServiceClient.SaveTaskFile(ctx, msg)
		if err != nil {
			code, msg := errs.ParseGrpcError(err)
			res.ResponseErrorWithMsg(c, code, msg)
			return
		}
	}
	res.ResponseSuccess(c, gin.H{
		"file":        key,
		"hash":        "",
		"key":         key,
		"url":         "http://localhost/" + key,
		"projectName": req.ProjectName,
	})
}
