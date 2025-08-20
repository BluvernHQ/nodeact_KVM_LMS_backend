import csv
import json
import argparse
import requests
from datetime import datetime

# Define the expected CSV headers
CSV_HEADERS = [
    "Role", "DOB", "Name", "Class", "Board", "School", "Address",
    "Subjects[0]", "Subjects[1]", "Subjects[2]", "Mode",
    "Father Name", "Father Phone", "Mother Name", "Mother Phone",
    "Guardian Name", "Guardian Phone", "Batch"
]

def row_to_payload(row):
    # Collect non-empty subjects
    subjects = [row.get(f"Subjects[{i}]") for i in range(3) if row.get(f"Subjects[{i}]")]
    return {
        "Name": row["Name"],
        "DOB": row["DOB"],
        "Board": row["Board"],
        "School": row["School"],
        "ClassId": row["Class"],
        "BatchId": row["Batch"],
        "Mode": row["Mode"],
        "Address": row["Address"],
        "FatherName": row["Father Name"],
        "FatherPhone": row["Father Phone"],
        "MotherName": row["Mother Name"],
        "MotherPhone": row["Mother Phone"],
        "GuardianName": row["Guardian Name"],
        "GuardianPhone": row["Guardian Phone"],
        "Subjects": subjects,
        "TimeStamp": str(datetime.now().isoformat()),
        "Sessions": []
    }
kl = 1
def main(csv_path, post_url, username):
    with open(csv_path, newline='', encoding='utf-8') as csvfile:
        reader = csv.DictReader(csvfile)
        for row in reader:
            payload = row_to_payload(row)
            # headers = {
            #     "Add-User-Name": f"{username}{reader.line_num - 2}@kvmtcc.org",
            #     "Add-User-Pwd": "123456"
            # }
            response = requests.post(post_url, json=payload)
            print(f"POST {post_url} => {response.status_code}")
            try:
                print(response.json())
            except Exception:
                print(response.text)
                print(kl)
                kl += 1
            # input("jk")

if __name__ == "__main__":
    main("11sb.csv", "http://localhost:5600/CreateStudent","12CB25")
