#!/usr/bin/env python3
"""
Script to delete all Firebase Auth users
WARNING: This script will permanently delete ALL users from Firebase Auth
"""

import firebase_admin
from firebase_admin import credentials, auth
import sys
import time

# Firebase service account key file
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

def get_all_users():
    """Get all Firebase Auth users"""
    try:
        users = []
        page = auth.list_users()
        
        while page:
            for user in page.users:
                users.append(user)
            
            # Get next page if available
            page = page.get_next_page() if page.has_next_page else None
        
        return users
    except Exception as e:
        print(f"✗ Error fetching users: {str(e)}")
        return []

def delete_user(uid):
    """Delete a single user by UID"""
    try:
        auth.delete_user(uid)
        return True
    except Exception as e:
        print(f"✗ Error deleting user {uid}: {str(e)}")
        return False

def main():
    """Main function to delete all Firebase Auth users"""
    print("🔥 Firebase Auth User Deletion Script")
    print("=" * 50)
    
    # Initialize Firebase
    if not initialize_firebase():
        sys.exit(1)
    
    # Get all users
    print("\n📋 Fetching all Firebase Auth users...")
    users = get_all_users()
    
    if not users:
        print("📭 No users found in Firebase Auth")
        return
    
    user_count = len(users)
    print(f"👥 Found {user_count} users")
    
    # Show warning and get confirmation
    print("\n⚠️  WARNING: This action will permanently delete ALL users!")
    print("This cannot be undone!")
    
    confirmation = input("\nType 'DELETE ALL USERS' to confirm: ")
    if confirmation != "DELETE ALL USERS":
        print("❌ Operation cancelled")
        return
    
    # Delete users
    print(f"\n🗑️  Deleting {user_count} users...")
    deleted_count = 0
    failed_count = 0
    
    for i, user in enumerate(users, 1):
        user_info = f"{user.uid}"
        if user.email:
            user_info += f" ({user.email})"
        
        print(f"[{i}/{user_count}] Deleting user: {user_info}")
        
        if user.email == "admin@kvmtcc.org":
            print(f"⏩ Skipping admin user: {user_info}")
            continue
        if user.email[0] in ['a','b','c','d','e','f','g','h','i','j','k','l','m','n','o','p','q','r','s','t','u','v','w','x','y','z']:
            print(f"⏩ Skipping user: {user_info}")
            continue
        if delete_user(user.uid):
            deleted_count += 1
            print(f"✓ Successfully deleted user: {user_info}")
        else:
            failed_count += 1
            print(f"✗ Failed to delete user: {user_info}")
        
        # Small delay to avoid rate limiting
        time.sleep(0.1)
    
    # Summary
    print("\n" + "=" * 50)
    print("📊 DELETION SUMMARY")
    print(f"✓ Successfully deleted: {deleted_count} users")
    print(f"✗ Failed to delete: {failed_count} users")
    print(f"📈 Total processed: {user_count} users")
    
    if failed_count > 0:
        print("\n⚠️  Some users could not be deleted. Check the error messages above.")
        sys.exit(1)
    else:
        print("\n🎉 All users successfully deleted!")

if __name__ == "__main__":
    main() 