CreateStaff:
	@echo "building CreateStaff.go"
	go build -buildmode=plugin -o modules_bin/CreateStaff.so modules_src/CreateStaff.go

CreateStudent:
	@echo "building CreateStudent.go"
	go build -buildmode=plugin -o modules_bin/CreateStudent.so modules_src/CreateStudent.go

MarkAttendance:
	@echo "building MarkAttendance"
	go build -buildmode=plugin -o modules_bin/MarkAttendance.so modules_src/MarkAttendance.go

GetAttendance:
	@echo "building GetAttendance"
	go build -buildmode=plugin -o modules_bin/GetAttendance.so modules_src/GetAttendance.go

GetUser:
	@echo "building GetUser"
	go build -buildmode=plugin -o modules_bin/GetUser.so modules_src/GetUser.go

FetchUsers:
	@echo "building FetchUsers"
	go build -buildmode=plugin -o modules_bin/FetchUsers.so modules_src/FetchUsers.go
