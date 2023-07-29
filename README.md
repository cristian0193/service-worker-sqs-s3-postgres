# Service Worker SQS - S3 - Postgres

![Diagrama3 1](https://github.com/cristian0193/service-worker-sqs-s3-postgres/assets/11803196/c49fdc92-cbb2-4c99-acaf-2219610a67c7)

## Tabla de contenido
1. [Contexto](#contexto)
2. [Tecnologías](#tecnologías)
3. [Arquitectura](#arquitectura)
    * [Estructura del proyecto](#estructura-del-proyecto)
5. [Despliegues](#despliegues)
    * [Local](#local)
6. [Endpoints](#endpoints)
7. [Queues](#queues)
8. [Buckets](#buckets)


<a name="contexto"></a>
# Contexto 📋

El objetivo de este template es proporcionar una estructura básica para crear un servicio de trabajador (worker) que establezca conexión con servicios de AWS, como SQS y RDS. Esta plantilla tiene como finalidad simplificar la construcción de servicios similares, al proporcionar una configuración sencilla y reducir el tiempo necesario para su desarrollo.

- **Estructurar un servicio worker**: Crear una estructura clara y organizada para desarrollar un servicio de trabajador que cumpla con los requisitos del proyecto.
- **Establecer conexión con servicios AWS**: Configurar la conexión con servicios de AWS, como Amazon SQS (Simple Queue Service) y RDS (Relational Database Service).
- **Simplificar la configuración**: Proporcionar una configuración sencilla que permita a los desarrolladores personalizar y adaptar el servicio según sus necesidades específicas.
- **Reducir tiempos de construcción**: Agilizar el proceso de desarrollo al brindar una plantilla predefinida y una estructura base que sirva como punto de partida para la construcción del servicio.
- **Mejorar la reutilización de código**: Fomentar la reutilización de código al proporcionar una estructura común y componentes predefinidos que puedan ser utilizados en diferentes proyectos de servicios de trabajador.


<a name="tecnologías"></a>
# Tecnologias 💻

**Dependencies** 🤝
Las siguientes dependencias se utilizan en el desarrollo para llevar a cabo depliegue de servidor http, conexiones SQS y RDS, entre otros.

* [github.com/aws/aws-sdk-go](https://github.com/aws/aws-sdk-go): SDK oficial de AWS para el lenguaje de programación Go.

* [github.com/labstack/echo/v4](https://github.com/labstack/echo): Echo es un framework de aplicaciones web Go de código abierto, extensible y centrado en el rendimiento.

* [gorm.io/gorm](https://gorm.io/): Es una libreria que permite el mapeo objeto-relacional (ORM), es una técnica que permite consultar y manipular datos de una base de datos utilizando un paradigma orientado a objetos.

* [gorm.io/driver/postgres](https://github.com/go-gorm/postgres): Controlador que permite establecer la conexion entre una base de datos postgres y un cliente.

* [go.uber.org/zap](https://pkg.go.dev/go.uber.org/zap): Permite la configuracion de logs.

**Framework**

* [Echo](https://echo.labstack.com/)
* [Gorm](https://gorm.io/)

**Servicios AWS**

* [SQS](https://aws.amazon.com/es/sqs/)
* [S3](https://aws.amazon.com/es/s3/)

**Bases de datos**

* [Postgresql](https://www.postgresql.org/)

<a name="arquitectura"></a>
# Arquitectura 🏢

Para del proyecto se toma como base los principios de las arquitecturas limpias, utilizando en este caso gran parte del concepto de **arquitectura multicapas**, lo cual permite la independencia de frameworks, entidades externas y UI, por medio de capas con responsabilidad únicas que permite ser testeables mediante el uso de sus interfaces. Como parte de las buenas prácticas la solución cuenta en su gran mayoría con la aplicación de los principios SOLID, garantizando un código limpio, mantenible, reutilizable y escalable.

![Diagrama3](https://github.com/cristian0193/service-worker-sqs-s3-postgres/assets/11803196/a2020517-4724-410b-83ea-ee4e96815ec2)

<a name="estructura-del-proyecto"></a>
### * **Estructura del proyecto** 🧱

- [x] `config/`: establece las configuraciones iniciales a los servicios
    - [ ] `cmd/`: administra los recursos de llamados al api
        - [ ] `builder/`: construye cada una de las instancias transversales
- [x] `core/`: establece la logica de negocio
    - [ ] `domain/`: administracion de los datos de manera transversal
    - [ ] `usecases/`: define los casos de uso utilizados por el handler
- [x] `dataproviders/`: contiene la implementacion de los clients externos
    - [ ] `awss3/`: define el cliente para aws s3
       - [ ] `downloader/`: genera la creacion del archivo temporal descargado desde s3
    - [ ] `awssqs/`: define el cliente para aws sqs
    - [ ] `consumer/`: define la logica para obtener los mensajes desde el consumidor
    - [ ] `postgres/`: define el cliente que permite la conexion a base de dato
      - [ ] `repository/`: define las consultas, actualizacion o inserciones a la base de datos
    - [ ] `processor/`: define el inicio del proceso para la lectura de mensajes desde SQS 
    - [ ] `server/`: define la configuracion para correr el server http
    - [ ] `utils/`: define las funciones transversales
- [x] `entrypoints/`: administra los recursos de llamados al api
    - [ ] `controllers/`: define los handler

<a name="despliegues"></a>
# Despliegues 🚀

Para la fase de despliegue a nivel local, se utilizaron algunas herramientas que nos permite agilizar este proceso. Como fase se mostrara el paso a paso:

> **Nota:** Para el proceso se deben definir las variables de ambiente que nos permite establecer conexion a los diferentes servicios.

```
APPLICATION_ID=
SERVER_PORT=
LOG_LEVEL=INFO

AWS_ACCESS_KEY=
AWS_SECRET_KEY=
AWS_REGION=

AWS_SQS_URL=
AWS_SQS_MAX_MESSAGES=
AWS_SQS_VISIBILITY_TIMEOUT=

AWS_S3_BUCKET=

DB_PORT=
DB_HOST=
DB_NAME=
DB_USERNAME=
DB_PASSWORD=
```

<a name="local"></a>
### * **Local** 

En el proceso local podemos utilizar despliegues de contenedores con postgres RDS o local.

    1. Instalacion postgres local o aws (rds)
        - docker pull postgres
        - https://aws.amazon.com/es/rds/

    2. Creacion de bucket S3 en AWS
        - https://aws.amazon.com/es/s3/        

    3. Creacion de SQS en AWS
        - https://aws.amazon.com/es/sqs/

    4. Editar politica de acceso en SQS para reportar eventos desde S3
        - https://docs.aws.amazon.com/es_es/AmazonS3/latest/userguide/ways-to-add-notification-config-to-bucket.html

    5. Automigracion en gorm activa

    6. Definir variables de entorno

    7. Start 'go run main.go'

<a name="endpoints"></a>
# Endpoints 🤖

**FileData**

- **GET**    http://localhost:8080/s3/filedata/:id
```
curl --location --request GET 'http://localhost:8080/s3/filedata/:id'
```

- **Response**
```
  {
    "id": "7a312c5a-e69e-4935-9b33-5dc33919a76f",
    "message": "Hola Mundo!!",
    "owner": "charodriguez",   
    "date": "2023-06-13T17:48:05-05:00"
  }
```

**MetaData**

- **GET**    http://localhost:8080/s3/metadata/:trackid
```
curl --location --request GET 'http://localhost:8080/s3/metadata/:trackid'
```

- **Response**
```
  {
    "trackid": "7a312c5a-e69e-4935-9b33-5dc33919a76f",
    "bucket": "s3-service-worker",
    "filename": "file-test.csv",
    "key": "/files/file-test.csv",
    "size": 350
  }
```

<a name="queues"></a>
# Queues 📨

- **URL**    https://sqs.us-east-1.amazonaws.com/XXXXXXXX/service-worker-sqs-s3-postgres

- **Message**
```
    {
     "Records": [
       {
         "eventSource": {
           "aws:s3": {
             "s3": {
               "bucket": {
                 "name": "s3-service-worker"
               },
               "object": {
                 "key": "/files/file-test.csv",
                 "size": 350
               }
             }
           }
         }
       }
     ]
   }
```

<a name="buckets"></a>
# Buckets 📂

- **Name**    s3-service-worker
- **Folder**  /files

# Author 🧑‍💻
```
- Christian Alexis Rodriguez Castillo
- Sr Software Engineer - Mercadolibre
- cristian010193@gmail.com
```
