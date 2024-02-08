# Golang Demo

### Note: 
Install all the packages in the project and then run project by: 
```go run main.go```

This api is created in golang using gin framework, this api supports postgresql and used sqlc to generate code for the sql files provided. This project also uses the pgx/v5 as sql_package. six Endpoints will be exposed from this project.

This repo contains 3 different assignments

* Users Module
* Re-Arrange String Module
* Sql to Swap two consecutive students.

### Users Module
A RESTful API using Golang with the Gin framework and PostgreSQL as the database and pgx as a database driver and sqlc for queries. The API should provide functionality to create a user, generate OTP for the user, and verify the OTP.

* **POST /api/users** >> Create Users
	* Accept JSON payload containing "name" & "phone_number"
	* Phone Number should be unique
	* Validation for phone number is important, it accepts phone number with (+).
	* Stores the user in the database.
* **POST /api/users/generateotp** >> Generates the OTP for user
	* Accepts JSON payload with `phone_number`.
	* If the `phone_number` does not exist, return a 404 error.
	* Generate a random 4-digit OTP and set its expiration time to 1 minute from the current time.
*  	**POST /api/users/verifyotp** >> Verify the OTP
	* Accepts JSON payload with `phone_number` and `otp`.
	* Check if the OTP is correct and not expired (compare with `otp_expiration_time`).
	* If the OTP is correct and not expired, return a success message.
	* If the OTP is incorrect, return an error message.
	* If the OTP is expired, return an error message indicating that the OTP has expired.

### Re-Arrangement of String
Given a string s, rearrange the characters of s so that any two adjacent characters are not the same.
Return any possible rearrangement of s or return "" if not possible.

**Example:** 
`Input: s = "aab"`
`Output: "aba"`

**Endpoint**

`POST "/api/rearrange-string"`
**Body**`{"s": "aab"}`

### Swap Students Module
Solution to re-arrange the seat_id of every two consecutive students. If no of students is odd, the id of last student is not swapped.

GET "/api/swap-students"

**Query**

``` 
WITH RankedStudents AS (
    SELECT id, student,
           ROW_NUMBER() OVER (ORDER BY id) AS RowNum
    FROM Seat
),
SwappedStudents AS (
    SELECT id,
           CASE
               WHEN RowNum % 2 = 0 THEN (
                   SELECT id FROM RankedStudents
                   WHERE RowNum = s.RowNum - 1
               )
               ELSE (
                   SELECT id FROM RankedStudents
                   WHERE RowNum = s.RowNum + 1
               )
           END AS student
    FROM RankedStudents s
)
SELECT id, student
FROM SwappedStudents
ORDER BY id;
```






