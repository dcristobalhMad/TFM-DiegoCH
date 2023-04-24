---
hide:
  - footer
---

# Tecnologías

## **Pulumi**

![Logo-pulumi](images/pulumi.svg)

Pulumi es un framework de código abierto que permite crear, desplegar y gestionar infraestructura en la nube de forma programática. En este caso, he decidido usar Golang como lenguaje base para la creación de la infraestructura, pero también se puede usar Python, TypeScript o .NET.

## **Fluentd**

![Logo-fluentd](images/fluentd.jpeg)

Fluentd es un recolector de datos que permite unificar los logs de diferentes aplicaciones en un único punto. En este caso, se ha usado para recoger los logs de las peticiones al HAProxy de los backend de las aplicaciones y mandarlos a Kinesis Data Streams.

## **HAProxy**

![Logo-hap](images/Haproxy-logo.png)

HAProxy es un balanceador de carga de código abierto que permite distribuir las peticiones entre diferentes servidores.

## **GitHub Actions**

![Logo-actions](images/github_actions.png)

Github actions es un servicio de integración y despliegue continuo que permite automatizar tareas en un repositorio de GitHub. En este caso, se ha usado para crear un flujo de trabajo que permite desplegar la infraestructura en AWS cuando se hace un merge a master.
