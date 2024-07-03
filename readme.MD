# Prometheus Config API Using Gin framework (Golang)

## Table of Contents

- [Features](#features)
- [API Endpoints](#api-endpoints)

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

  - Success (200):

    ```json
    {
      "status": "success",
      "message": "Jobs retrieved successfully",
      "data": [
        {
          "job_name": "job1",
          "static_configs": [
            {
              "targets": ["11.111.111.111:26", "11.111.111.111:27"]
            }
          ]
        },
        {
          "job_name": "job2",
          "static_configs": [
            {
              "targets": ["11.111.111.112:26", "11.111.111.112:27"]
            }
          ]
        }
      ]
    }
    ```

  - Fail (500):

    ```json
    {
      "status": "error",
      "message": "Failed to load configuration",
      "data": null
    }
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

  - Success (200):
    ```json
    {
      "status": "success",
      "message": "Job added successfully",
      "data": {
        "job_name": "newjob",
        "static_configs": [
          {
            "targets": ["11.111.111.111:26", "11.111.111.111:27"]
          }
        ]
      }
    }

    ```
  - Fail (400):
    ```json
    {
      "status": "error",
      "message": "Invalid request body",
      "data": null
    }
    ```

  - Fail (409):
    ```json
    {
      "status": "error",
      "message": "Job name already exists",
      "data": null
    }
    ```
  - Fail (409):
    ```json
    {
      "status": "error",
      "message": "Job with this IP address already exists",
      "data": null
    }
    ```
  - Fail (500):
    ```json
    {
      "status": "error",
      "message": "Failed to load configuration",
      "data": null
    }
    ```

  - Fail (500):
    ```json
    {
      "status": "error",
      "message": "Failed to save configuration",
      "data": null
    }
    ```
  - Fail (500):
    ```json
    {
      "status": "error",
      "message": "docker die roy",
      "data": null
    }
    ```

### Remove a Job

**URL**: `/jobs/:job_name`  
**Method**: `DELETE`  
**Description**: Removes a job specified by job_name.

**Request:**
none

**Response:**
  - Success (200):
  ```json
  {
    "status": "success",
    "message": "Job removed successfully",
    "data": "job1"
  }
  ```

  - Fail (404):
  ```json
  {
    "status": "error",
    "message": "Job not found",
    "data": "null"
  }
  ```

  - Fail (500):
    ```json
    {
      "status": "error",
      "message": "Failed to load configuration",
      "data": null
    }
    ```

  - Fail (500):
    ```json
    {
      "status": "error",
      "message": "Failed to save configuration",
      "data": null
    }
    ```
  - Fail (500):
    ```json
    {
      "status": "error",
      "message": "docker die roy",
      "data": null
    }
    ```

### Search Job by IP Address

**URL**: `/jobs/search`  
**Method**: `GET`  
**Description**: Searches for a job by IP address.

**Request:**
Query parameter: `ip` (string) - The IP address to search for. An example below:

```bash
/jobs/search?ip=11.111.111.111
```


**Response:**

  - Success (200):
    ```json
    {
      "status": "success",
      "message": "Job(s) retrieved successfully",
      "data": [
        {
          "job_name": "job1",
          "static_configs": [
            {
              "targets": ["11.111.111.111:26", "11.111.111.111:27"]
            }
          ]
        }
      ]
    }
    ```

  - Fail (400):
    ```json
    {
      "status": "error",
      "message": "IP address is required",
      "data": null
    }
    ```

  - Fail (404):
    ```json
    {
      "status": "error",
      "message": "No job found with the specified IP address",
      "data": null
    }
    ```

  - Fail (500):
    ```json
    {
      "status": "error",
      "message": "Failed to load configuration",
      "data": null
    }
    ```
