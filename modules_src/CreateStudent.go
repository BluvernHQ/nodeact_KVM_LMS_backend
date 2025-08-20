package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	// "strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"firebase.google.com/go/v4/auth"
)

type Student struct {
	UID           string   `json:"UID" bson:"UID"`
	Role          string   `json:"Role" bson:"Role"`
	TimeStamp     string   `json:"TimeStamp" bson:"TimeStamp"`
	DOB           string   `json:"DOB" bson:"DOB"`
	Name          string   `json:"Name" bson:"Name"`
	ClassId       string   `json:"ClassId" bson:"ClassId"`
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
	BatchId       string   `json:"BatchId" bson:"BatchId"`
	ProfilePic    string   `json:"ProfilePic" bson:"ProfilePic"`
	Sessions      []string `json:"Sessions" bson:"Sessions"`
	RegNo         string   `json:"RegNo" bson:"RegNo"`
	Email         string   `json:"Email" bson:"Email"`
}

type Class struct {
	Name    string `json:"Name" bson:"Name"`
	BatchId string `json:"BatchId" bson:"BatchId"`
	BoardId string `json:"BoardId" bson:"BoardId"`
}
type Board struct {
	Name	string   `json:"Name" bson:"Name"`
}
type Batch struct {
	Name	string   `json:"Name" bson:"Name"`
}

func CreateStudent(w http.ResponseWriter, r *http.Request, db *mongo.Client, authClient *auth.Client) {
	fmt.Println("HTTP Method:", r.Method)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// token, err := authClient.VerifyIDToken(ctx, r.Header.Get("Authorization"))
	// if err != nil {
	// 	http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
	// 	return
	// }
	// fmt.Println("Token verified:", token.UID)

	// collection := db.Database("KVM").Collection("Users")
	// var Role bson.M
	// err = collection.FindOne(ctx, bson.M{"UID": token.UID}, options.FindOne().SetProjection(bson.M{"Role": 1, "_id": 0})).Decode(&Role)
	// if err != nil {
	// 	http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// if Role["Role"] == nil {
	// 	http.Error(w, "Role not found", http.StatusNotFound)
	// 	return
	// }
	// if Role["Role"] != "admin" && Role["Role"] != "DEVELOPER" {
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	return
	// }

	var data Student
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	//get class
	collection := db.Database("KVM").Collection("Classes")
	var class Class
	objectId, err := primitive.ObjectIDFromHex(data.ClassId)
	if err != nil {
		fmt.Println("Invalid ClassId format")
		http.Error(w, "Invalid ClassId format: "+err.Error(), http.StatusBadRequest)
		return
	}
	err = collection.FindOne(ctx, bson.M{"_id": objectId}, options.FindOne().SetProjection(bson.M{"Name": 1, "_id": 0})).Decode(&class)
	if err != nil {
		fmt.Println("Class not found")
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//get batch
	collection = db.Database("KVM").Collection("Batches")
	var batch Batch
	objectId, err = primitive.ObjectIDFromHex(data.BatchId)
	if err != nil {
		fmt.Println("Invalid BatchId format ", data.BatchId)
		http.Error(w, "Invalid BatchId format: "+err.Error(), http.StatusBadRequest)
		return
	}
	err = collection.FindOne(ctx, bson.M{"_id": objectId}, options.FindOne().SetProjection(bson.M{"Name": 1, "_id": 0})).Decode(&batch)
	if err != nil {
		fmt.Println("Batch not found")
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//get board
	collection = db.Database("KVM").Collection("Boards")
	var board Board
	objectId, err = primitive.ObjectIDFromHex(data.Board)
	if err != nil {
		fmt.Println("Invalid Board format")
		http.Error(w, "Invalid Board format: "+err.Error(), http.StatusBadRequest)
		return
	}
	err = collection.FindOne(ctx, bson.M{"_id": objectId}, options.FindOne().SetProjection(bson.M{"Name": 1, "_id": 0})).Decode(&board)
	if err != nil {
		fmt.Println("Board not found")
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//total users
	collection = db.Database("KVM").Collection("Users")
	totalUsers, err := collection.CountDocuments(ctx, bson.M{"Role": "student"})
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	totalUsersInt := int(totalUsers)

	classStr, _ := json.Marshal(class)
	batchStr, _ := json.Marshal(batch)
	boardStr, _ := json.Marshal(board)
	fmt.Println("class ,",string(classStr))
	fmt.Println("batch ,",string(batchStr))
	fmt.Println("board ,",string(boardStr))
	fmt.Println("totalUsers ,",totalUsersInt)


	className := class.Name
	batchName := batch.Name
	boardName := board.Name

	if batchName == "2025-2026"{
		batchName = "25"
	}
	if boardName == "CBSE" {
		boardName = "CB"
	}
	className = strings.Split(className, " ")[0]

	

	ROLLNUMBER := className + boardName + batchName + fmt.Sprintf("%03d", totalUsersInt+1)
	data.RegNo = ROLLNUMBER
	data.Email = ROLLNUMBER+"@kvmtcc.org"

	fmt.Println("Email:", data.Email)
	fmt.Println("Password:", data.RegNo)

	params := (&auth.UserToCreate{}).
		Email(data.Email).
		Password(data.RegNo)
	userRecord, err := authClient.CreateUser(context.Background(), params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating UID: %v", err), http.StatusInternalServerError)
		return
	}
	data.UID = userRecord.UID
	data.Role = "student"

	collection = db.Database("KVM").Collection("Users")
	result, err := collection.InsertOne(ctx, data)
	if err != nil {
		http.Error(w, "Insert error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "User created successfully",
		"insertedID": result.InsertedID,
		"UID":        data.UID,
	})
}
