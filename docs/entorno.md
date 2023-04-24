---
hide:
  - footer
---

# Entorno de desarrollo

Se ha creado un entorno de desarrollo en docker en la carpeta devenv con el fin de poder levantar un pequeño laboratorio para mandar datos a Kinesis Data Streams y poder probar la infraestructura. Para ello, se ha creado un docker-compose con los siguientes servicios:

- Haproxy
- GoappX
- Fluentd

Para levantar los servicios `make dev-up` y para tirarlo `make dev-down`. En el fichero de fluentd.conf hay que configurar las variables para que pueda mandar los logs al stream de kinesis. Una vez rellenadas las variables y levantando el docker-compose, se puede realizar una sencilla prueba para ver el flow del log: `curl -v localhost:8100` y se verá en el log del contenedor fluentd que se ha mandado el log de la request a kinesis.

<p align="center">
  <img src="/images/dev-env.png">
</p>
