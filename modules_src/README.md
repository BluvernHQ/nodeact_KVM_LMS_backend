# create batch
## /CreateBatch
Method: POST
Body:{
    Name: ""
}

# create Class
## /CreateClass
Method: POST
Body:{
    BatchId: ""
    Name: ""
}

# create session
## /CreateSession
Method: POST
Body:{
    Name: ""
    TimeFrom: ""
    TimeTo: ""
    SubjectId: ""
    ClassId: ""
    BatchId: ""
    Staff: ""
    day: ""
}

# create staff
## /CreateStaff
Method: POST
Body:{
    UID: ""
    Name: ""
    Role: ""
    TimeStamp: timestamp
    DOB: ""
    Qualification: ""
    Subjects: ""
    Experience: ""
    Phone: ""
    WorkingAt: ""
    OtherSpecialization: ""
    BatchId: ""
    ProfilePic: ""
}

# create student
## /CreateStudent
Method: POST
Body:{
    UID: ""
    Role: ""
    TimeStamp: timestamp
    DOB: ""
    Name: ""
    ClassId: ""
    Board: ""
    School: ""
    Address: ""
    Subjects: ""
    Mode: ""
    FatherName: ""
    FatherPhone: ""
    MotherName: ""
    MotherPhone: ""
    GuardianName: ""
    GuardianPhone: ""
    Status: ""
    BatchId: ""
    ProfilePic: ""
}

# create subject
## /CreateSubject
Method: POST
Body:{
    Name: ""
}

# delete document
## /DeleteDoc
Method: POST
Body:{
    Id: ""
    Collection: ""
}

# fetch document
## /FetchDoc
Method: GET
Body:{
    Query: {} //the db is populated by /CreateDoc
    Paging:{
        Page: 1
        Limit: 1
    }
    Projection: {} //what fields should it return and what not.like {_id:0,BatchId:1}
    Collection: ""
}

# get user
## /GetUser
Method: GET
Body:NULL

# mark attendance
## /MarkAttendance
Method: POST
Body:{
    SessionId: ""
    UIDs: []
    MarkedAt: timestamp
}

# update document
## /UpdateDoc
Method: POST
Body:{
    Id: ""
    Collection: ""
    Set: {} //used to change a field
}

# upload
## /Upload
Methd: POST
Body: multipartfile
Return: {
    URL:
}

# upload details
## /CreateUpload
Method: POST
Body:{
    URL: "" //this call is called as soon as /Upload
    BatchId: ""
    ClassId: ""
}


# NOTES

* This commands are all designed as mongoDB template
* "collections" are the name of the database collection
* Paging:{} is used to set Page and docs per page is set by Limit
* 