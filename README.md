# TP0: Docker + Comunicaciones + Concurrencia

En el presente repositorio se provee un ejemplo de cliente-servidor el cual corre en containers con la ayuda de [docker-compose](https://docs.docker.com/compose/). El mismo es un ejemplo práctico brindado por la cátedra para que los alumnos tengan un esqueleto básico de cómo armar un proyecto de cero en donde todas las dependencias del mismo se encuentren encapsuladas en containers. El cliente (Golang) y el servidor (Python) fueron desarrollados en diferentes lenguajes simplemente para mostrar cómo dos lenguajes de programación pueden convivir en el mismo proyecto con la ayuda de containers.

Por otro lado, se presenta una guía de ejercicios que los alumnos deberán resolver teniendo en cuenta las consideraciones generales descriptas al pie de este archivo.

## Instrucciones de uso
El repositorio cuenta con un **Makefile** que posee encapsulado diferentes comandos utilizados recurrentemente en el proyecto en forma de targets. Los targets se ejecutan mediante la invocación de:

* **make \<target\>**:
Los target imprescindibles para iniciar y detener el sistema son **docker-compose-up** y **docker-compose-down**, siendo los restantes targets de utilidad para el proceso de _debugging_ y _troubleshooting_.

Los targets disponibles son:
* **docker-compose-up**: Inicializa el ambiente de desarrollo (buildear docker images del servidor y cliente, inicializar la red a utilizar por docker, etc.) y arranca los containers de las aplicaciones que componen el proyecto.
* **docker-compose-down**: Realiza un `docker-compose stop` para detener los containers asociados al compose y luego realiza un `docker-compose down` para destruir todos los recursos asociados al proyecto que fueron inicializados. Se recomienda ejecutar este comando al finalizar cada ejecución para evitar que el disco de la máquina host se llene.
* **docker-compose-logs**: Permite ver los logs actuales del proyecto. Acompañar con `grep` para lograr ver mensajes de una aplicación específica dentro del compose.
* **docker-image**: Buildea las imágenes a ser utilizadas tanto en el servidor como en el cliente. Este target es utilizado por **docker-compose-up**, por lo cual se lo puede utilizar para testear nuevos cambios en las imágenes antes de arrancar el proyecto.
* **build**: Compila la aplicación cliente para ejecución en el _host_ en lugar de en docker. La compilación de esta forma es mucho más rápida pero requiere tener el entorno de Golang instalado en la máquina _host_.

### Servidor
El servidor del presente ejemplo es un EchoServer: los mensajes recibidos por el cliente son devueltos inmediatamente. El servidor actual funciona de la siguiente forma:
1. Servidor acepta una nueva conexión.
2. Servidor recibe mensaje del cliente y procede a responder el mismo.
3. Servidor desconecta al cliente.
4. Servidor procede a recibir una conexión nuevamente.

### Cliente
El cliente del presente ejemplo se conecta reiteradas veces al servidor y envía mensajes de la siguiente forma.
1. Cliente se conecta al servidor.
2. Cliente genera mensaje incremental.
recibe mensaje del cliente y procede a responder el mismo.
3. Cliente envía mensaje al servidor y espera mensaje de respuesta.
Servidor desconecta al cliente.
4. Cliente verifica si aún debe enviar un mensaje y si es así, vuelve al paso 2.

Al ejecutar el comando `make docker-compose-up` para comenzar la ejecución del ejemplo y luego el comando `make docker-compose-logs`, se observan los siguientes logs:

```
client1  | 2024-08-21 22:11:15 INFO     action: config | result: success | client_id: 1 | server_address: server:12345 | loop_amount: 5 | loop_period: 5s | log_level: DEBUG
client1  | 2024-08-21 22:11:15 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°1
server   | 2024-08-21 22:11:14 DEBUG    action: config | result: success | port: 12345 | listen_backlog: 5 | logging_level: DEBUG
server   | 2024-08-21 22:11:14 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:15 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:15 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°1
server   | 2024-08-21 22:11:15 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:20 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:20 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°2
server   | 2024-08-21 22:11:20 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:20 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°2
server   | 2024-08-21 22:11:25 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:25 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°3
client1  | 2024-08-21 22:11:25 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°3
server   | 2024-08-21 22:11:25 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:30 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:30 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°4
server   | 2024-08-21 22:11:30 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:30 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°4
server   | 2024-08-21 22:11:35 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:35 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°5
client1  | 2024-08-21 22:11:35 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°5
server   | 2024-08-21 22:11:35 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:40 INFO     action: loop_finished | result: success | client_id: 1
client1 exited with code 0
```

## Parte 1: Introducción a Docker
En esta primera parte del trabajo práctico se plantean una serie de ejercicios que sirven para introducir las herramientas básicas de Docker que se utilizarán a lo largo de la materia. El entendimiento de las mismas será crucial para el desarrollo de los próximos TPs.

### Ejercicio N°1:
Además, definir un script de bash `generar-compose.sh` que permita crear una definición de DockerCompose con una cantidad configurable de clientes.  El nombre de los containers deberá seguir el formato propuesto: client1, client2, client3, etc. 

El script deberá ubicarse en la raíz del proyecto y recibirá por parámetro el nombre del archivo de salida y la cantidad de clientes esperados:

`./generar-compose.sh docker-compose-dev.yaml 5`

Considerar que en el contenido del script pueden invocar un subscript de Go o Python:

```
#!/bin/bash
echo "Nombre del archivo de salida: $1"
echo "Cantidad de clientes: $2"
python3 mi-generador.py $1 $2
```



### Ejercicio N°2:
Modificar el cliente y el servidor para lograr que realizar cambios en el archivo de configuración no requiera un nuevo build de las imágenes de Docker para que los mismos sean efectivos. La configuración a través del archivo correspondiente (`config.ini` y `config.yaml`, dependiendo de la aplicación) debe ser inyectada en el container y persistida afuera de la imagen (hint: `docker volumes`).



### Ejercicio N°3:
Crear un script de bash `validar-echo-server.sh` que permita verificar el correcto funcionamiento del servidor utilizando el comando `netcat` para interactuar con el mismo. Dado que el servidor es un EchoServer, se debe enviar un mensaje al servidor y esperar recibir el mismo mensaje enviado.

En caso de que la validación sea exitosa imprimir: `action: test_echo_server | result: success`, de lo contrario imprimir:`action: test_echo_server | result: fail`.

El script deberá ubicarse en la raíz del proyecto. Netcat no debe ser instalado en la máquina _host_ y no se puede exponer puertos del servidor para realizar la comunicación (hint: `docker network`). `



### Ejercicio N°4:
Modificar servidor y cliente para que ambos sistemas terminen de forma _graceful_ al recibir la signal SIGTERM. Terminar la aplicación de forma _graceful_ implica que todos los _file descriptors_ (entre los que se encuentran archivos, sockets, threads y procesos) deben cerrarse correctamente antes que el thread de la aplicación principal muera. Loguear mensajes en el cierre de cada recurso (hint: Verificar que hace el flag `-t` utilizado en el comando `docker compose down`).



## Parte 2: Repaso de Comunicaciones

Las secciones de repaso del trabajo práctico plantean un caso de uso denominado **Lotería Nacional**. Para la resolución de las mismas deberá utilizarse como base al código fuente provisto en la primera parte, con las modificaciones agregadas en el ejercicio 4.



### Ejercicio N°5:
Modificar la lógica de negocio tanto de los clientes como del servidor para nuestro nuevo caso de uso.

#### Cliente
Emulará a una _agencia de quiniela_ que participa del proyecto. Existen 5 agencias. Deberán recibir como variables de entorno los campos que representan la apuesta de una persona: nombre, apellido, DNI, nacimiento, numero apostado (en adelante 'número'). Ej.: `NOMBRE=Santiago Lionel`, `APELLIDO=Lorca`, `DOCUMENTO=30904465`, `NACIMIENTO=1999-03-17` y `NUMERO=7574` respectivamente.

Los campos deben enviarse al servidor para dejar registro de la apuesta. Al recibir la confirmación del servidor se debe imprimir por log: `action: apuesta_enviada | result: success | dni: ${DNI} | numero: ${NUMERO}`.



#### Servidor
Emulará a la _central de Lotería Nacional_. Deberá recibir los campos de la cada apuesta desde los clientes y almacenar la información mediante la función `store_bet(...)` para control futuro de ganadores. La función `store_bet(...)` es provista por la cátedra y no podrá ser modificada por el alumno.
Al persistir se debe imprimir por log: `action: apuesta_almacenada | result: success | dni: ${DNI} | numero: ${NUMERO}`.

#### Comunicación:
Se deberá implementar un módulo de comunicación entre el cliente y el servidor donde se maneje el envío y la recepción de los paquetes, el cual se espera que contemple:
* Definición de un protocolo para el envío de los mensajes.
* Serialización de los datos.
* Correcta separación de responsabilidades entre modelo de dominio y capa de comunicación.
* Correcto empleo de sockets, incluyendo manejo de errores y evitando los fenómenos conocidos como [_short read y short write_](https://cs61.seas.harvard.edu/site/2018/FileDescriptors/).



### Ejercicio N°6:
Modificar los clientes para que envíen varias apuestas a la vez (modalidad conocida como procesamiento por _chunks_ o _batchs_). La información de cada agencia será simulada por la ingesta de su archivo numerado correspondiente, provisto por la cátedra dentro de `.data/datasets.zip`.
Los _batchs_ permiten que el cliente registre varias apuestas en una misma consulta, acortando tiempos de transmisión y procesamiento.

En el servidor, si todas las apuestas del *batch* fueron procesadas correctamente, imprimir por log: `action: apuesta_recibida | result: success | cantidad: ${CANTIDAD_DE_APUESTAS}`. En caso de detectar un error con alguna de las apuestas, debe responder con un código de error a elección e imprimir: `action: apuesta_recibida | result: fail | cantidad: ${CANTIDAD_DE_APUESTAS}`.

La cantidad máxima de apuestas dentro de cada _batch_ debe ser configurable desde config.yaml. Respetar la clave `batch: maxAmount`, pero modificar el valor por defecto de modo tal que los paquetes no excedan los 8kB. 

El servidor, por otro lado, deberá responder con éxito solamente si todas las apuestas del _batch_ fueron procesadas correctamente.



### Ejercicio N°7:
Modificar los clientes para que notifiquen al servidor al finalizar con el envío de todas las apuestas y así proceder con el sorteo.
Inmediatamente después de la notificacion, los clientes consultarán la lista de ganadores del sorteo correspondientes a su agencia.
Una vez el cliente obtenga los resultados, deberá imprimir por log: `action: consulta_ganadores | result: success | cant_ganadores: ${CANT}`.

El servidor deberá esperar la notificación de las 5 agencias para considerar que se realizó el sorteo e imprimir por log: `action: sorteo | result: success`.
Luego de este evento, podrá verificar cada apuesta con las funciones `load_bets(...)` y `has_won(...)` y retornar los DNI de los ganadores de la agencia en cuestión. Antes del sorteo, no podrá responder consultas por la lista de ganadores.
Las funciones `load_bets(...)` y `has_won(...)` son provistas por la cátedra y no podrán ser modificadas por el alumno.



## Parte 3: Repaso de Concurrencia

### Ejercicio N°8:
Modificar el servidor para que permita aceptar conexiones y procesar mensajes en paralelo.
En este ejercicio es importante considerar los mecanismos de sincronización a utilizar para el correcto funcionamiento de la persistencia.

En caso de que el alumno implemente el servidor Python utilizando _multithreading_,  deberán tenerse en cuenta las [limitaciones propias del lenguaje](https://wiki.python.org/moin/GlobalInterpreterLock).



## Consideraciones Generales
Se espera que los alumnos realicen un _fork_ del presente repositorio para el desarrollo de los ejercicios.El _fork_ deberá contar con una sección de README que indique como ejecutar cada ejercicio.

La Parte 2 requiere una sección donde se explique el protocolo de comunicación implementado.
La Parte 3 requiere una sección que expliquen los mecanismos de sincronización utilizados.

Cada ejercicio deberá resolverse en una rama independiente con nombres siguiendo el formato `ej${Nro de ejercicio}`. Se permite agregar commits en cualquier órden, así como crear una rama a partir de otra, pero al momento de la entrega deben existir 8 ramas llamadas: ej1, ej2, ..., ej7, ej8.

(hint: verificar listado de ramas y últimos commits con `git ls-remote`)

Puden obtener un listado del último commit de cada rama ejecutando `git ls-remote`.

Finalmente, se pide a los alumnos leer atentamente y **tener en cuenta** los criterios de corrección provistos [en el campus](https://campusgrado.fi.uba.ar/mod/page/view.php?id=73393).


# Soluciones

## Parte 1

### Ejercicio 1

Para la resolución de este ejercicio se creo un programa en python para generar un archivo de docker compose. Primero se creo el script principal generar-compose.sh, que recibe dos parametros: el nombre del archivo de salida y la cantidad de clientes, y luego ejecuta el programa de python (mi-generador.py) pasandole esos parametros.

El el programa mi-generador.py genera un archivo de docker-compose con la configuración de un servidor y la cantidad de clientes que se le pase por parametro.

Para correr la solución del ejercicio, parado sobre la rama ej1, se ejecuta:

./generar-compose.sh docker-compose-dev.yaml x

siendo x la cantidad de clientes a configurar.
Luego se puede seguir con los comandos de make docker-compose-up y make docker-compose-logs para correr el programa usando la configuración generada.

### Ejercicio 2
Se modificó el programa mi-generador.py añadiendo un volumen en el servidor que mapea el config.ini a una ruta dentro del contenedor, y lo mismo para el cliente con el archivo config.yaml.

Por ultimo elimine una linea en el dockerfile del cliente que copiaba el archivo config.yaml dentro del la imagen del contenedor cuando se construía, la cual ya no era necesaria.

Igual que en el ejercicio anterior, para correr la solución del ejercicio, parado sobre la rama ej2, se ejecuta:

./generar-compose.sh docker-compose-dev.yaml x
siendo x la cantidad de clientes a configurar.
Luego se puede seguir con los comandos de make docker-compose-up y make docker-compose-logs para correr el programa usando la configuración generada.


### Ejercicio 3

Se creó el script validar-echo-server.sh, el cual verifica el funcionamiento del echo server usando netcat.

En este script el servidor y un contenedor de prueba se contectan a una red, mediante la cual el contenedor temporal le mandará un mensaje de prueba al servidor. Luego el servidor le debería responder con el mismo mensaje si esta funcionando correctamente. Si el mensaje recibido es igual al enviado, el script imprimirá "action: test_echo_server | result: success", de lo contrario imprimirá "action: test_echo_server | result: fail".

Para correrlo: Primero verificamos que tenga permisos de ejecución y luego lo corremos:

chmod +x validar-echo-server.sh
./validar-echo-server.sh


### Ejercicio 4

Cliente
Se agregó un canal sigChan para señales del sistema, para poder capturar el SIGTERTM, si se recibe la señal se llama a un metodo del cliente que cierra los recursos correctamente.

Servidor
Para el servidor se realizo algo parecido, se utiliza el modulo signal de python para captural la señal, y si se recibe la señal, se detiene el bucle de aceptacion de conexiones y se llama a una funcion handle sigterm que maneja el cierre de los recursos correctamente.

## Parte 2

### Ejercicio 5
Tipos de mensajes en el protocolo:

- BET = 1: Mensaje enviado por el cliente al servidor con los datos de la apuesta.
- OK = 2: Mensaje enviado por el servidor al cliente para confirmar la recepción de la apuesta.
- ERROR = 3: Mensaje enviado por el servidor al cliente para indicar que hubo un error al recibir la apuesta.
  
Cliente:
El cliente ahora recibe los datos de una persona por variable de entorno, y se los envia al servidor mediante un protocolo el cual serializa los datos a enviar y luego espera la respuesta de confirmacion del servidor. El protocolo maneja la comunicación entre el cliente y un servidor usando sockets y utiliza funciones como htonl y ntohl para convertir enteros a bytes en formato big-endian, y readFully y writeFully para garantizar que todos los datos se envíen y reciban completamente, evitando short reads y writes. Se utiliza el método sendBet para serializar y envía una apuesta, incluyendo todos sus componentes, mientras que receiveMessage lee y procesa el mensaje de respuesta del servidor.

Servidor:
El servidor recibe los datos de la apuesta de los clientes mediante el protocolo, almacena la información y envia nuevamente mediante el protocolo el mensaje de éxito al cliente. Este protocolo tiene el método de recibir la apuesta en donde recibe los datos aplicando la conversión de endianness y usando una función como receive_exact para evitar los short reads y usando el metodo de sendall de socket para evitar los short writes.

### Ejercicio 6
Tipos de mensajes en el protocolo:

- BET = 1: Mensaje enviado por el cliente al servidor con los datos de la apuesta.
- OK = 2: Mensaje enviado por el servidor al cliente para confirmar la recepción de la apuesta.
- ERROR = 3: Mensaje enviado por el servidor al cliente para indicar que hubo un error al recibir la apuesta.
- FINISH = 4: Mensaje enviado por el cliente al servidor para indicar que terminó de enviar todas las apuestas.

Cliente

El cliente lee el archivo de apuestas de la agencia y envía los datos al servidor en chunks de tamaño configurable. Cada batch es enviado en un mensaje de tipo BET. Luego, el cliente espera la respuesta del servidor. Si el servidor responde con un mensaje de tipo OK, el cliente envía el siguiente batch. Si el servidor responde con un mensaje de tipo ERROR, el cliente imprime un mensaje de error y sigue con el siguiente batch. Cuando termina de enviar todas las apuestas, el cliente envía un mensaje de tipo FINISH al servidor.

El protocolo de comunicación entre el cliente y el servidor sigue siendo el mismo del ejercicio anterior, con la adición de los mensajes FINISH (que el cliente envia y el servidor recibe) y y que la funcion de sendBet en el cliente ahora recibe un array de apuestas en lugar de una sola apuesta.

Servidor
El servidor recibe los mensajes de los clientes y procesa las apuesta. Si todas las apuestas del batch fueron procesadas correctamente, el servidor responde con un mensaje de tipo OK. Si hubo un error con alguna de las apuestas, el servidor responde con un mensaje de tipo ERROR. Cuando recibe un mensaje de tipo FINISH, el servidor cierra la conexión con el cliente.


### Ejercicio 7
Tipos de mensajes en el protocolo:
- BET = 1: Mensaje enviado por el cliente al servidor con los datos de la apuesta.
- OK = 2: Mensaje enviado por el servidor al cliente para confirmar la recepción de la apuesta.
- ERROR = 3: Mensaje enviado por el servidor al cliente para indicar que hubo un error al recibir la apuesta.
- FINISH = 4: Mensaje enviado por el cliente al servidor para indicar que terminó de enviar todas las apuestas.


Cliente
Se modificó que ahora despues de que el cliente envía un mensaje FINISH al servidor cuando termina de enviar todas las apuestas, espera un mensaje del servidor con la lista de ganadores de la agencia. 

Servidor
Se modificó que ahora el servidor espera la notificación de FINISH las N agencias para realizar el sorteo. Luego de este evento, verifica cada apuesta con las funciones load_bets(...) y has_won(...) y retorna los DNI de los ganadores de la agencia en cuestión a cada cliente.

El protocolo de comunicación entre el cliente y el servidor sigue siendo el mismo del ejercicio anterior, con la adición de que el cliente ahora tiene un metodo para recibir la lista de ganadores del servidor y el servidor ahora tiene un metodo para enviar la lista de ganadores al cliente.

## Parte 3 

### Ejercicio 8

Se modificó el servidor para que acepte conexiones y procese mensajes en paralelo. Para ello se utilizó la librería ThreadPoolExecutor de Python. Se creó un pool de threads para los clientes (agencias) que se conectan al servidor, y se ejecuta la función handle_client en un thread distinto para cada cliente, así cada sigue procesando mensajes en paralelo mientras el servidor acepta nuevas conexiones.

Se utilizó un Lock para sincronizar el acceso al archivo bets.csv, el cual se utiliza para persistir las apuestas y asi evitar race conditions.

Al finalizar el server, se cierran los threads y termina la ejecución.

El protocolo de comunicación entre el cliente y el servidor sigue siendo el mismo del ejercicio anterior que todos los ejercicios anteriores y no se modificó para este ejercicio.

Lo vuelvo a comentar por las dudas: 

Tipos de mensajes en el protocolo:
- BET = 1: Mensaje enviado por el cliente al servidor con los datos de la apuesta.
- OK = 2: Mensaje enviado por el servidor al cliente para confirmar la recepción de la apuesta.
- ERROR = 3: Mensaje enviado por el servidor al cliente para indicar que hubo un error al recibir la apuesta.
- FINISH = 4: Mensaje enviado por el cliente al servidor para indicar que terminó de enviar todas las apuestas.

Algunas funciones del protocolo del cliente: 
- sendBets envía un chunk de apuestas al servidor.
- receiveMessage: recibe un mensaje del servidor.
- sendFinish: envía un mensaje de FINISH al servidor.
- receiveWinners: recibe la lista de ganadores del servidor.

Algunas funciones del protocolo del servidor:
- receiveBets: recibe un chunk de las apuesta del cliente.
- sendSuccess: envía un mensaje de OK al cliente.
- sendError: envía un mensaje de ERROR al cliente.

