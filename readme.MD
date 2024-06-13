# Prometheus Config API Using Gin framework (Golang)

## Table of Contents

- [Features](#features)
- [API Endpoints](#api-endpoints)
- [Example curls](#example-curls)

## Features

This project provides APIs to list all jobs in JSON format, as well as add or remove a job.

## API Endpoints

### List All Jobs

**URL**: `/jobs`  
**Method**: `GET`  
**Description**: Retrieves a list of all jobs in JSON format.

**Request:**
None

**Response:**
```json
[
  {
    "job_name": "job1",
    "ip_address": "11.111.111.111"
  },
  {
    "job_name": "job2",
    "ip_address": "22.222.222.222"
  }
]
```

### Add a Job

**URL**: `/jobs`  
**Method**: `POST`  
**Description**: Adds a new job. Requires job_name and ip_address in the request body.

**Request:**
```json

{
  "job_name": "newjob",
  "ip_address": "11.111.111.111"
}
```

**Response:**
```json
{
  "status": "ok",
  "message": "Job added successfully"
}
```

### Remove a Job

**URL**: `/jobs/:job_name`  
**Method**: `DELETE`  
**Description**: Removes a job specified by job_name.

**Request:**
none

**Response:**
```json
{
  "status": "ok",
  "message": "Job removed successfully"
}
```