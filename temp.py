import firebase_admin
from firebase_admin import credentials, auth
import sys
import time
import pymongo

client = pymongo.MongoClient("mongodb://admin:machine2003@localhost:27017/admin")
db = client["KVM"]
collection = db["Users"]

SERVICE_ACCOUNT_KEY_FILE = "nodeact-kvm-firebase-adminsdk-fbsvc-c9fe0118fb.json"

def initialize_firebase():
    """Initialize Firebase Admin SDK"""
    try:
        cred = credentials.Certificate(SERVICE_ACCOUNT_KEY_FILE)
        firebase_admin.initialize_app(cred)
        print("✓ Firebase Admin SDK initialized successfully")
        return True
    except Exception as e:
        print(f"✗ Error initializing Firebase: {str(e)}")
        return False

def create_user(uid, email, password):
    user = auth.create_user(
    uid=uid,
    email=email,
    password=password
)

def delete_user(uid):
    """Delete a single user by UID"""
    try:
        print(f"Deleting user {uid}")
        auth.delete_user(uid)
        return True
    except Exception as e:
        print(f"✗ Error deleting user {uid}: {str(e)}")
        return False

users = collection.find()

initialize_firebase()

reg =  0

for user in users:
    if(user["Role"] != "student"):
        continue
    reg += 1
    prev = user["RegNo"]
    print(prev[:-3])
    newReg = prev[:-3] + str(reg).zfill(3)
    newMail = newReg + "@kvmtcc.org"
    print(prev , " -> " , newReg)
    collection.update_one({"UID": user["UID"]}, {"$set": {"RegNo": newReg, "Email": newMail}})
    create_user(user["UID"], newMail, "123456*#")
    

# import requests
# for user in users:
#     if user["Role"] == "student":
#         continue
#     # print(user)
# create_user("4GFRySJMttelFdMyuWzGChhgdOi2", "11CB25300@kvmtcc.org", "123456")