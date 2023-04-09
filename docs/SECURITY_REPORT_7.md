#  Vulnerable and outdated third party components

***
<font color="orange">Severity : 7.8 High</font>  
Date : 04/04/2023  
Reporter : Nicolas Viaud   
Weakness enumeration : [CWE-1104: Use of Unmaintained Third Party Components](https://cwe.mitre.org/data/definitions/1104.html)  
CVSS v3.1 Vector : `AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H/E:U/RL:O/RC:U`  
***

Some third party components, used in the project, are outdated or not even maintained.  
Here the list of the third party components that were released more than 1 year ago:  

| Dependency                      | Version             | Release date | Usage      | Comment        |
|---------------------------------|---------------------|--------------|------------|----------------|
| github.com/dgrijalva/jwt-go     | v3.2.0+incompatible | Jun 8, 2021  | go.mod     | Not maintained |
| github.com/go-kit/kit           | v0.11.0             | Jul 4, 2021  | go.mod     |                |
| github.com/lib/pq               | v1.10.2             | May 17, 2021 | go.mod     |                |
| hub.docker.com/_/golang         | 1.18-rc-stretch     | Mar 4, 2022  | Dockerfile |                |
| github.com/kubernetes-sigs/kind | v0.13.0             | May 10, 2022 | Makefile   |                |
| github.com/rossmcf/hey          | v1.0.0              | Jul 6, 2021  | Makefile   |                |

These third party components can be subject to some common vulnerabilities and exposure (CVE). If an attacker identify some CVE, due to third party outdated library, he can try to exploit them to gain access to the system.

# Developer Fix

### Task 1
Change the `github.com/dgrijalva/jwt-go` library by the officially supported one `github.com/golang-jwt/jwt`.

### Task 2
Update to the latest version of all the third party dependencies in the [go.mod](../go.mod), [Dockerfile](../Dockerfile) and [Makefile](../Makefile).

### Task 3
Add a job in the CI to scan the CVE. The tool used for the scan could be for example [trivy](https://github.com/aquasecurity/trivy).  
If at least one CVE `Medium`, `High` or `Critical` is detected, the job should result with an error. The scan report should be sent to the security team to further analyses.

