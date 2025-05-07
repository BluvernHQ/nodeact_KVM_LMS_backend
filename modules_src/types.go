package main

import ()

type Student struct {
	UID           string   `json:"Id" bson:"Id"`
	Role          string   `json:"Role" bson:"Role"`
	TimeStamp     string   `json:"TimeStamp" bson:"TimeStamp"`
	DOB           string   `json:"DOB" bson:"DOB"`
	Name          string   `json:"Name" bson:"Name"`
	Class         string   `json:"Class" bson:"Class"`
	Board         string   `json:"Board" bson:"Board"`
	School        string   `json:"School" bson:"School"`
	Address       string   `json:"Address" bson:"Address"`
	Subjects      []string `json:"Subjects" bson:"Subjects"`
	Mode          string   `json:"Mode" bson:"Mode"`
	FatherName    string   `json:"FatherName" bson:"FatherName"`
	FatherPhone   string   `json:"FatherPhone" bson:"FatherPhone"`
	MotherName    string   `json:"MotherName" bson:"MotherName"`
	MotherPhone   string   `json:"MotherPhone" bson:"MotherPhone"`
	GuardianName  string   `json:"GuardianName" bson:"GuardianName"`
	GuardianPhone string   `json:"GuardianPhone" bson:"GuardianPhone"`
	Status        string   `json:"Status" bson:"Status"`
	Batch         string   `json:"Batch" bson:"Batch"`
	ProfilePic    string   `json:"ProfilePic" bson:"ProfilePic"`
}

type Staff struct {
	UID                 string   `json:"Id" bson:"Id"`
	Role                string   `json:"Role" bson:"Role"`
	TimeStamp           string   `json:"TimeStamp" bson:"TimeStamp"`
	DOB                 string   `json:"DOB" bson:"DOB"`
	Qualification       string   `json:"Qualification" bson:"Qualification"`
	Subjects            []string `json:"Subjects" bson:"Subjects"`
	Experience          string   `json:"Experience" bson:"Experience"`
	Phone               string   `json:"Phone" bson:"Phone"`
	WorkingAt           string   `json:"WorkingAt" bson:"WorkingAt"`
	OtherSpecialization string   `json:"OtherSpecialization" bson:"OtherSpecialization"`
	Batch               []string `json:"Batch" bson:"Batch"`
	ProfilePic          string   `json:"ProfilePic" bson:"ProfilePic"`
}

type RequestPayload struct {
	Query  map[string]interface{} `json:"query"`
	Paging struct {
		Page  int64 `json:"page"`
		Limit int64 `json:"limit"`
	} `json:"paging"`
}

type Fees struct {
	Amount           string `json:"Amount"`
	DueDate          string `json:"DueDate"`
	PaymentTimeStamp string `json:"PaymentDate"`
	Status           string `json:"Status"`
	StudentId        string `json:"StudentId"`
	TimeStamp        string `json:"TimeStamp"`
}

type Leave struct {
	FromDate     string `json:"FromDate"`
	ToDate       string `json:"ToDate"`
	Reason       string `json:"Reason"`
	TimeStamp    string `json:"TimeStamp"`
	CreatedBy    string `json:"CreatedBy"`
	Status       string `json:"Status"`
	AuthorisedBy string `json:"AuthorisedBy"`
}

type Attendance struct {
	AttendanceBy string `json:"AttendanceBy"`
	Session      string `json:"Session"`
	AttendanceTo string `json:"AttendanceTo"`
	MarkedAt   string `json:"MarkedAt"`
}
