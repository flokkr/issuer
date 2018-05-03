# Kerberos ticket issuer

This is a simple REST application to generate (and cache) kerberos tickets and java ssl keys and truststores.

It should run together with the krb5 server and it will create the required kerberos keytab files and publish it over the rest interface.

*FOR DEVELOPMENT ONLY*: As it can give you any keytab this is strictly for development only and not for production use..

## Rest endpoints

### Create/get keytabs

Generates a keytab for the service SERVICENAME/HOST@EXAMPLE.COM

*URL*: http://localhost:80801/keytab/HOST/SERVICENAME

*Source*: https://github.com/flokkr/issuer/blob/master/bash/keytab.sh

### Create/get java keystore and truststore


*URL for keystore*: http://localhost:8081/keystore/NAME

*URL for truststore: http://localhost:8081/truststore

*Source*: https://github.com/flokkr/issuer/blob/master/bash/root.sh

*Source*: https://github.com/flokkr/issuer/blob/master/bash/issue.sh


## Development

The bash scripts are included with the help of [go-bindata](https://github.com/jteeuwen/go-bindata). You need it on the path.

The easiest way to build it is using [goreleaser](https://github.com/goreleaser/goreleaser) utility,
