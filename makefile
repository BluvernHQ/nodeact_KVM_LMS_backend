deploy:
	@echo "building modules"
	@make all
	@go build -o myapp main.go
	@echo "restating KVM_api in pm2"
	@pm2 restart KVM_api

all: CreateStaff CreateStudent MarkAttendance GetAttendance GetUser FetchUsers

CreateStaff:
	@echo "building CreateStaff"
	@go build -buildmode=plugin -o modules_bin/CreateStaff.so modules_src/CreateStaff.go

CreateStudent:
	@echo "building CreateStudent"
	@go build -buildmode=plugin -o modules_bin/CreateStudent.so modules_src/CreateStudent.go

MarkAttendance:
	@echo "building MarkAttendance"
	@go build -buildmode=plugin -o modules_bin/MarkAttendance.so modules_src/MarkAttendance.go

GetAttendance:
	@echo "building GetAttendance"
	@go build -buildmode=plugin -o modules_bin/GetAttendance.so modules_src/GetAttendance.go

GetUser:
	@echo "building GetUser"
	@go build -buildmode=plugin -o modules_bin/GetUser.so modules_src/GetUser.go

FetchUsers:
	@echo "building FetchUsers"
	@go build -buildmode=plugin -o modules_bin/FetchUsers.so modules_src/FetchUsers.go
