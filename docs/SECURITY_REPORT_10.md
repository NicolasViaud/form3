# No JWT token expiration 

***
<font color="orange">Severity : 4.2 Medium</font>  
Date : 04/06/2023  
Reporter : Nicolas Viaud   
Weakness enumeration : [CWE-613: Insufficient Session Expiration](https://cwe.mitre.org/data/definitions/613.html)  
CVSS v3.1 Vector : `AV:N/AC:H/PR:N/UI:R/S:U/C:L/I:L/A:N`
***

There is currently no expiration feature for JWT Token. When generated, the JWT token will be valid forever (or more precisely, until that the HS256 secrets changed). The risk is that if one token is stolen, it can be used to access protected resource by impersonating the identity of the token owner and will never expire.  
To mitigate the risk in case of token leak, the JWT token should have an expiration date. The duration of a token has to be defined by the company security policy, but it has to be the shortest possible.

If the application become more complex, we could plan to migrate to a more robust IAM solution that enforce all this token security concern, like for example the excellent open source [keycloak Solution](https://www.keycloak.org/) maintained by RedHat.

# Development Fix

During JWT token generation, in the code [cmd/token/main.go](../cmd/token/main.go), the fields creation date (`ias`) and expiration date (`exp`) should be added to the payload:
```
{
  "iat": 1516230000,
  "exp": 1516240000,
  "admin": true,
  "hotel": 456,
  "name": "Nico",
  "sub": "8790c514-73b6-400f-8f28-acc74d342a44"
}
```
The value of the creation date is set to the current timestamp.  
The value of the expiration date is set to the creation date + token lifetime.  
A flag with the token lifetime need to be added in the token generation CLI. By default, the value is 30 minutes.


During token validation, in the code [jwtauth/middleware.go](../jwtauth/middleware.go), the expiration date need to be verified. If the current timestamp is greater or equals than the expiration timestamp, the error `ErrTokenInvalid` should be thrown.

Some unit tests can be added:  

#### Nominal case for a non expiry JWT token
> **Given** a JWT token `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NzI1MzEyMDAsImV4cCI6MTY3MjUzMzAwMCwiYWRtaW4iOnRydWUsImhvdGVsIjo0NTYsIm5hbWUiOiJOaWNvIiwic3ViIjoiODc5MGM1MTQtNzNiNi00MDBmLThmMjgtYWNjNzRkMzQyYTQ0In0.nX2Evt6Mqp2tfNPIn_2a6uEChJ9fFci6ziop9-eS6Ek` used in the request  
> **AND** the verify token signature return always true (mock)  
> **AND** the current date is `Sun Jan 01 2023 00:01:00 GMT+0000`  
> **When** the middleware logic is called  
> **Then** the result should be success

#### Error case for a expiry JWT token
> **Given** a JWT token `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NzI1MzEyMDAsImV4cCI6MTY3MjUzMzAwMCwiYWRtaW4iOnRydWUsImhvdGVsIjo0NTYsIm5hbWUiOiJOaWNvIiwic3ViIjoiODc5MGM1MTQtNzNiNi00MDBmLThmMjgtYWNjNzRkMzQyYTQ0In0.nX2Evt6Mqp2tfNPIn_2a6uEChJ9fFci6ziop9-eS6Ek` used in the request  
> **AND** the verify token signature return always true (mock)  
> **AND** the current date is `Sun Jan 01 2023 00:31:00 GMT+0000`  
> **When** the middleware logic is called  
> **Then** the result should be the error `ErrTokenInvalid`  