# Infraestructura del trabajo de fin de máster de Diego Cristóbal Herreros para la Universidad Alfonso X de Madrid

**\*Nota**: Este proyecto es la infraestructura, el flow completo del dato se explica en la literatura del trabajo\*

## Descripción del proyecto

La carpeta principal donde se encuentra el código de la infraestructura es la carpeta 'Infrastructure'. En ella se encuentran los archivos destinados a crear en Amazon Web Services los siguientes servicios:

- Kinesis data streams
- Kinesis data firehose
- Lambda
- S3

## Flow automático por github actions

- Creando una PR, se verá que infraestructura se crea o modifica en AWS
- Merge a master, se hará el despliegue/modificación de la infraestructura
- Tag con version de semver + tag 'destroy' (Example `0.0.1-destroy`). Se destruirá la infraestructura

## Ejecución en local

También se puede realizar el proceso en local, para ello se ha creado un Makefile con los siguientes comandos:

- Compilación de la lambda: `make build-lambda`
- Plan/preview de la infraestructura: `make check`
- Despliegue de la infraestructura: `make deploy`
- Destrucción de la infraestructura: `make destroy`

## Entorno de desarrollo

Se ha creado un entorno de desarrollo en docker en la carpeta devenv con el fin de poder levantar un pequeño laboratorio para mandar datos a Kinesis Data Streams y poder probar la infraestructura. Para ello, se ha creado un docker-compose con los siguientes servicios:

- Haproxy
- GoappX
- Fluentd

Para levantar los servicios `make dev-up` y para tirarlo `make dev-down`. En el fichero de fluentd.conf hay que configurar las variables para que pueda mandar los logs al stream de kinesis. Una vez rellenadas las variables y levantando el docker-compose, se puede realizar una sencilla prueba para ver el flow del log: `curl -v localhost:8100` y se verá en el log del contenedor fluentd que se ha mandado el log de la request a kinesis.
