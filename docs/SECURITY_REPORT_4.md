# Insecure secret storage

***
<font color="orange">Severity : 6.4 Medium</font>  
Date : 04/04/2023  
Reporter : Nicolas Viaud   
Weakness enumeration : [CWE-522: Insufficiently Protected Credentials](https://cwe.mitre.org/data/definitions/522.html)  
CVSS v3.1 Vector : `AV:A/AC:H/PR:H/UI:N/S:U/C:H/I:H/A:H`
***

The secrets of the application are stored in a Kubernetes `ConfigMap`. This is not a safe place for sensitive data because they are stored in plain text and, according to the current configuration of the cluster, anyone with a cluster access can read or modify them.  

Without any cluster hardening, the Kubernetes `Secrets` are not a safe place neither because it also stores data in plain text by default.

> *Caution:*  
> *Kubernetes Secrets are, by default, stored unencrypted in the API server's underlying data store (etcd). Anyone with API access can retrieve or modify a Secret, and so can anyone with access to etcd. Additionally, anyone who is authorized to create a Pod in a namespace can use that access to read any Secret in that namespace; this includes indirect access such as the ability to create a Deployment*  
>
> **source** : https://kubernetes.io/docs/concepts/configuration/secret

The main requirements for managing sensitive data are:
* sensitive data need to be encrypted on the hard drive (at rest)
* sensitive data need to be encrypted when transmitting in the network (on transit)
* the access to the sensitive data need to be restricted with strong access control
* the access to the sensitive data has to be audited
* if the sensitive data is a secret, it has to be changed regularly

A secret manager solution meet all these requirements. A popular open source secret manager solution is [Vault from Hashicorp](https://www.hashicorp.com/products/vault) and could be used in the project.

# Developer Fix

All the application secrets and username (postgresql username and password, and the HS256 secret) has to be retrieved by API from the [Vault](https://www.hashicorp.com/products/vault). To do so, you can use the Vault API client library `github.com/hashicorp/vault/api`. The TLS communication version 1.3 has to be used in order to encrypted secret on transit.  
All the installation and configuration of the Vault server will be done by the security team.