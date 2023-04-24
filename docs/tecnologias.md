---
hide:
  - footer
---

# Tecnologías

## **Pulumi**

<p align="center">
<img src="/images/pulumi.svg">
</p>

Pulumi es un framework de código abierto que permite crear, desplegar y gestionar infraestructura en la nube de forma programática. En este caso, he decidido usar Golang como lenguaje base para la creación de la infraestructura, pero también se puede usar Python, TypeScript o .NET.

## **Fluentd**

<p align="center">
  <img src="/images/fluentd.jpeg">
</p>

Fluentd es un recolector de datos que permite unificar los logs de diferentes aplicaciones en un único punto. En este caso, se ha usado para recoger los logs de las peticiones al HAProxy de los backend de las aplicaciones y mandarlos a Kinesis Data Streams.

## **HAProxy**

<p align="center">
  <img src="/images/Haproxy-logo.png">
</p>

HAProxy es un balanceador de carga de código abierto que permite distribuir las peticiones entre diferentes servidores.

## **GitHub Actions**

<p align="center">
  <img src="/images/github_actions.png">
</p>

Github actions es un servicio de integración y despliegue continuo que permite automatizar tareas en un repositorio de GitHub. En este caso, se ha usado para crear un flujo de trabajo que permite desplegar la infraestructura en AWS cuando se hace un merge a master.
