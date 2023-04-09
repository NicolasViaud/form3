# Usage of hostPath in kubernetes PersistentVolume

***
<font color="#ffcb0d">Severity : 3.4 Low</font>  
Date : 04/04/2023  
Reporter : Nicolas Viaud   
Weakness enumeration : [CWE-653: Improper Isolation or Compartmentalization](https://cwe.mitre.org/data/definitions/653.html)  
CVSS v3.1 Vector : `AV:L/AC:H/PR:H/UI:N/S:U/C:L/I:L/A:L/E:U/RL:O/RC:C`
***

Resource sharing between nodes and pods should be avoided because the main security concern about containerization and virtualization is to isolate the container as much as possible from the host.

The current configuration of the Kubernetes create a `PersistentVolume` named `postgres-pv-volume`. This `PersistentVolume` mount the directory `/mnt/data` of the node file system, with a read and write access, via the [hostPath](https://kubernetes.io/docs/concepts/storage/volumes/#hostpath) directive.
```yaml
  accessModes:
    - ReadWriteMany
  hostPath:
    path: "/mnt/data"
```
Even if the `PersistentVolume` appear to be not used by the postgres deployment due to a misconfiguration, if an attacker has the right to create a pod, he can take the opportunity to mount it into his pod and access to the `/mnt/data` directory on the node file system as root. This could give an opportunity for the attacker to gain access to the node server.


### Development Fix

Because this `PersistentVolume` is meant to be used by the posgresql deployment to persist data from the database, it needs to be persistent (can't use an `emptyDir`). Depending on where the application is hosted, I would recommend to mount a file system from a bucket from the cloud provider directly.

The security team is enforcing `Pod Security Admission` rules and will deprecate the usage of the `hostpath` for Q3.
