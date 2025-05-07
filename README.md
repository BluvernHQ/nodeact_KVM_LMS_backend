# Base architechture

![alt text](image.png)

Main runtime runt he go program and uses plugin module to dynamically import and use the functions.

`/modules_src/` contains all functions's source code

`/modules_bin/` contains all functions's compiled output file


# Writing functions
Each functions in this architechture can be edited via a dashboard and have its own twmplate.

## function parameter
All functions get a `req` object that contains all the details of request **URL** and request **header**

## funtion return
The return type of the function is a **http response** in json object.

# Security
The overlying HTTPS server is Nginx.