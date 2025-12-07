deploy:
	@echo "building modules"
	@make all
	@go build -o myapp main.go
	@echo "restating KVM-API in pm2"
	@pm2 restart KVM-BCK

all: ExportDocs MarkForTest CreateSection CreateMessage CreateUpload CreateSession DeleteDoc Upload CreateClass CreateBatch CreateBoard CreateSubject CreateStaff CreateStudent MarkAttendance GetUser FetchDocs UpdateDoc

run_backend:
	@pm2 start ./myapp --name KVM_api

run_frontend:
	cd /root/kvm_web
	@pm2 start npm --name KVM_admin -- run start

ExportDocs:
	@echo "building ExportDocs"
	@go build -buildmode=plugin -o modules_bin/ExportDocs.so modules_src/ExportDocs.go

MarkForTest:
	@echo "building MarkForTest"
	@go build -buildmode=plugin -o modules_bin/MarkForTest.so modules_src/MarkForTest.go

CreateUpload:
	@echo "building CreateUpload"
	@go build -buildmode=plugin -o modules_bin/CreateUpload.so modules_src/CreateUpload.go

CreateBoard:
	@echo "building CreateBoard"
	@go build -buildmode=plugin -o modules_bin/CreateBoard.so modules_src/CreateBoard.go

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

CreateMessage:
	@echo "building CreateMessage"
	@go build -buildmode=plugin -o modules_bin/CreateMessage.so modules_src/CreateMessage.go

CreateSection:
	@echo "building CreateSection"
	@go build -buildmode=plugin -o modules_bin/CreateSection.so modules_src/CreateSection.go

CreateFees:
	@echo "building CreateFees"
	@go build -buildmode=plugin -o modules_bin/CreateFees.so modules_src/CreateFees.go