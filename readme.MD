list all jobs:

curl -X GET http://localhost:8080/jobs

add a job:

curl -X POST http://localhost:8080/jobs -H "Content-Type: application/json" -d '{"job_name": "newjob", "ip_address": "11.111.111.111"}'

(run this 2 times to see what will happen if ip already exists in another job)
