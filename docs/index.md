---
hide:
  - footer
---

# **Infraestructura del trabajo de fin de máster de Diego Cristóbal Herreros para la Universidad Alfonso X de Madrid**

![Logo-UAX](images/Logo-UAX.png)

**\*Nota**: Este proyecto es la infraestructura, el flow completo del dato se explica en la literatura del trabajo\*

## Descripción general del proyecto

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
