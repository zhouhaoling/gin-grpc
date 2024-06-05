REM 该脚本用于生成proto文件,执行方法：cd到该目录下，/.user.bat执行该脚本 user_service.proto是可被更改的
protoc --proto_path=. --go_out=. --go-grpc_out=. user_service.proto