# Weak secret for HS256 JWT signature algorithm

***
<font color="#df3d03">Severity : 8.8 High</font>  
Date : 04/04/2023  
Reporter : Nicolas Viaud   
Weakness enumeration : [CWE-1391: Use of Weak Credentials](https://cwe.mitre.org/data/definitions/1391.html)  
CVSS v3.1 Vector : `AV:N/AC:L/PR:L/UI:N/S:U/C:H/I:H/A:H`  
***

The secrets used by the algorithm HS256 to sign and verify JWT tokens is not secure because:
* it doesn't have enough characters
* it use human readable words
 
An attacker with a valid JWT token, could easily recover this secret by brut force or dictionary attack. Furthermore, this secret is a really sensitive data because it is used for authentication mechanism. With this secret, an attacker could create an admin JWT token and access to all endpoints as an admin.

### Developer Fix

The recommendation is to use a random secret with at least 64 characters, generated the secure random number generator tool provide by the security team.  
Secrets have to be different for each environment (local, production, ..).

