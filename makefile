deploy:
	@echo "building modules"
	@make all
	@go build -o myapp main.go
	@echo "restating KVM_api in pm2"
	@pm2 restart KVM_api

all: CreateUpload CreateSession DeleteDoc Upload CreateClass CreateBatch CreateSubject CreateStaff CreateStudent MarkAttendance GetUser FetchDocs UpdateDoc

run:
	@pm2 start ./myapp --name KVM_api

CreateUpload:
	@echo "building CreateUpload"
	@go build -buildmode=plugin -o modules_bin/CreateUpload.so modules_src/CreateUpload.go

CreateSession:
	@echo "building CreateSession"
	@go build -buildmode=plugin -o modules_bin/CreateSession.so modules_src/CreateSession.go

DeleteDoc:
	@echo "building DeleteDoc"
	@go build -buildmode=plugin -o modules_bin/DeleteDoc.so modules_src/DeleteDoc.go

Upload:
	@echo "building Upload"
	@go build -buildmode=plugin -o modules_bin/Upload.so modules_src/Upload.go

CreateClass:
	@echo "building CreateClass"
	@go build -buildmode=plugin -o modules_bin/CreateClass.so modules_src/CreateClass.go

CreateBatch:
	@echo "building CreateSession"
	@go build -buildmode=plugin -o modules_bin/CreateBatch.so modules_src/CreateBatch.go

CreateSession:
	@echo "building CreateSession"
	@go build -buildmode=plugin -o modules_bin/CreateSession.so modules_src/CreateSession.go

CreateSubject:
	@echo "building CreateSubject"
	@go build -buildmode=plugin -o modules_bin/CreateSubject.so modules_src/CreateSubject.go

CreateStaff:
	@echo "building CreateStaff"
	@go build -buildmode=plugin -o modules_bin/CreateStaff.so modules_src/CreateStaff.go

CreateStudent:
	@echo "building CreateStudent"
	@go build -buildmode=plugin -o modules_bin/CreateStudent.so modules_src/CreateStudent.go

MarkAttendance:
	@echo "building MarkAttendance"
	@go build -buildmode=plugin -o modules_bin/MarkAttendance.so modules_src/MarkAttendance.go

GetUser:
	@echo "building GetUser"
	@go build -buildmode=plugin -o modules_bin/GetUser.so modules_src/GetUser.go

FetchDocs:
	@echo "building FetchDocs"
	@go build -buildmode=plugin -o modules_bin/FetchDocs.so modules_src/FetchDocs.go

UpdateDoc:
	@echo "building UpdateDoc"
	@go build -buildmode=plugin -o modules_bin/UpdateDoc.so modules_src/UpdateDoc.go