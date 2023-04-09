# Usage of insecure download in Makefile

***
<font color="orange">Severity : 6.4 Medium</font>  
Date : 04/06/2023  
Reporter : Nicolas Viaud   
Weakness enumeration : [CWE-295: Improper Certificate Validation](https://cwe.mitre.org/data/definitions/295.html)  
CVSS v3.1 Vector : `AV:A/AC:H/PR:H/UI:N/S:U/C:H/I:H/A:H`
***

The makefile command below is dangerous:
```shell
@curl -L --insecure https://github.com/rossmcf/hey/releases/download/v1.0.0/installer.sh | bash
```
The `--insecure` option, during the download of the `installer.sh` file, force the http client to don't verify the TLS certificate of the `github.com` domain. This certificate, delivered by a trustable third party, is used to prove the legitimate identity of the `github` website.
Without the certificate verification, an attacker could intercept the request (MITM) and pretend to be the `github` website. The consequence of a such attack could be disastrous because the command execute directly the content of the `installer.sh` file in a shell.  
The impact could be, for an attacker:
- gaining access to the developer computer (steal user password, keylogger, access to shared document, access to source code, privilege escalation,..)
- gaining access to the CI workers (steal server password, access to applicative servers, access to source code, ..)

## Developer Fix

Remove the `--insecure` option. 