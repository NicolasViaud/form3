# Docker image run as root

***
<font color="#ffcb0d">Severity : 3.4 Low</font>  
Date : 04/04/2023  
Reporter : Nicolas Viaud   
Weakness enumeration : [CWE-250: Execution with Unnecessary Privileges](https://cwe.mitre.org/data/definitions/250.html)  
CVSS v3.1 Vector : `AV:A/AC:H/PR:H/UI:N/S:U/C:L/I:L/A:L/E:U/RL:O/RC:X`
***

The Dockerfile of the innsecure application is run as root.
```shell
root@innsecure-645cdb9bbd-n9rxp:/src# id
uid=0(root) gid=0(root) groups=0(root)
```
The risk is that if an attacker gain access to a shell inside the docker container by exploiting a web application vulnerability, he will have the same right as the application. 
Because the web application is run as root, an attacker could gain root access inside the container (after a successful vulnerability exploitation). This high privilege inside the container allow for an attacker to:
* access sensitive data that only root can access, like sensitive configuration file
* keep trying finding some vulnerability inside the container with root privilege

Another security concern is that the container and the host share the same kernel. If an attacker exploit another vulnerability in the container to gain access to the host, he will be root on the host as well.

### Development Fix

In the Dockerfile, the user is created but not used
```dockerfile
RUN useradd -u 1200 builder
```
To use it, you need to add the command
```dockerfile
USER builder
```

It's important as well to add a group to the user `builder` because it would belong to the root group by default:
```dockerfile
RUN groupadd -g 1200 builder && useradd -r -u 1200 -g builder builder
```

To test the configuration, you can execute the command on your local kind cluster:
```shell
kubectl exec --context kind-kind -it -n innsecure $(kubectl get pods -n innsecure -l app=innsecure -o jsonpath='{.items[*].metadata.name}') -- id
```
The ouput should be 
```shell
uid=1200(builder) gid=1200(builder) groups=1200(builder)
```