## Backend for A Restaurant Management System 
Check the routes folder for incoming routes. Project using mongoDB so make sure the mongo server is running.
Check the models folder for data points to be sent and controllers folder to check modify functionality if required 

### Usage 
  *  Start the Server  ``` go run main.go ```
  *  Now you can request the endpoint on localhost:8000 using curl or POSTMAN
  *  Example using curl -  ``` curl -X POST http://localhost:8000/users/signup  -d "@data.json" ```

Signup data.json -

    {
       'First_name':'John',
       'Last_name': 'Doe',
       'Password': 'password',
       'Email': 'test@test.com',
       'Phone': '1234567890',
       'User_type': 'ADMIN' 
    }


  #### After signup or login a token is returned that needs to be passed for all further calls for that user in the header as 
  ``` {'token': 'value_ returned_after_login/signup_request'}```


