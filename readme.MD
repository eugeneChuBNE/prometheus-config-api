# Prometheus Config API Using Gin framework (Golang)

## Table of Contents

- [Features](#features)
- [API Endppoints](#api-endpoints)
- [Example curls](#example-curls)

## Features

This project provides APIs to list all jobs in JSON format, as well as add or remove a job

## API Endpoints:

### List All Jobs

**URL**: `/jobs`  
**Method**: `GET`

### Add a Job

**URL**: `/jobs`  
**Method**: `POST`

### Remove a Job

**URL**: `/jobs/:job_name`  
**Method**: `DELETE`

## Example curls

### list all jobs 

```sh
    curl -X GET http://localhost:8080/jobs
```

### Add a job 

```sh
    curl -X POST http://localhost:8080/jobs -H "Content-Type: application/json" -d '{"job_name": "newjob", "ip_address": "11.111.111.111"}'
```

Run this twice if you want to see what happen if job already exists

### remove a job 

```sh
    curl -X DELETE http://localhost:8080/jobs/newjob
```
