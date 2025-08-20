# Attendance Marking System Documentation

## Overview
The attendance marking system is a feature that allows staff members to record attendance for students in specific sessions. The system is built with security in mind, requiring proper authentication and authorization.

## Authentication & Authorization
- Only users with `admin`, `staff`, or `DEVELOPER` roles can mark attendance
- Each request requires a valid Firebase authentication token
- The token is verified before any attendance operation is performed

## Data Structure

### Attendance Record
```json
{
    "Session": "string",      // Session ID for which attendance is being marked
    "UID": "string",         // User ID of the student/staff
    "Role": "string",        // Role (student/staff)
    "MarkedAt": "string"     // Timestamp when attendance was marked
}
```

### Session Structure
```json
{
    "Name": "string",        // Name of the session
    "TimeFrom": "string",    // Session start time
    "TimeTo": "string",      // Session end time
    "SubjectId": "string",   // Subject ID
    "ClassId": "string",     // Class ID
    "BatchId": "string",     // Batch ID
    "Staff": "string",       // Staff member assigned
    "Day": "string"         // Day of the session
}
```

## API Endpoint

### Mark Attendance
- **Endpoint**: `/MarkAttendance`
- **Method**: POST
- **Headers Required**:
  - `Authorization`: Firebase ID token
  - `Content-Type`: application/json

### Request Body
```json
{
    "Session": "session_id",
    "UIDs": ["student_id1", "student_id2", ...],
    "MarkedAt": "timestamp"
}
```

### Response
- Success (200 OK):
```json
{
    "message": "Attendances inserted successfully"
}
```
- Error Cases:
  - 401 Unauthorized: Invalid token or insufficient permissions
  - 400 Bad Request: Invalid JSON payload
  - 500 Internal Server Error: Database operation failure

## How It Works

1. **Authentication Check**:
   - The system verifies the Firebase ID token provided in the request header
   - Extracts the user's UID from the token

2. **Authorization Check**:
   - Queries the Users collection to check the role of the token holder
   - Only allows access if the user has admin, staff, or DEVELOPER role

3. **Attendance Recording**:
   - Creates an attendance record for the staff member marking attendance
   - Creates individual attendance records for each student in the UIDs array
   - All records are stored with the same session ID and timestamp
   - Each record includes the role of the attendee (staff/student)

4. **Database Storage**:
   - All attendance records are stored in the "Attendance" collection in MongoDB
   - Uses batch insert operation to store multiple records efficiently

## Important Notes

1. **Staff Attendance**:
   - The system automatically records attendance for the staff member marking attendance
   - Staff attendance is marked with role="staff"

2. **Student Attendance**:
   - Multiple students can be marked present in a single API call
   - Each student gets an individual attendance record
   - Student attendance is marked with role="student"

3. **Session Validation**:
   - The session ID must be valid and exist in the Sessions collection
   - The session is linked to specific batch, class, and subject

4. **Time Tracking**:
   - The `MarkedAt` timestamp helps track when the attendance was recorded
   - This can be useful for audit trails and attendance analytics

## Security Considerations

1. Only authenticated users can access the attendance marking system
2. Role-based access control ensures only authorized personnel can mark attendance
3. Firebase authentication ensures secure token-based authentication
4. MongoDB transactions ensure data consistency 