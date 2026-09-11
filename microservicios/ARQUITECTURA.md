# Arquitectura: TP 1-4

**Alumna:** Sofía Peiretti

## 1. Punto de partida

El sistema original es un monolito de e-commerce en Go (gin) que corre en el puerto `8080`. Clientes, productos y pedidos conviven en el mismo proceso y en el mismo archivo de controladores, sin separación de responsabilidades.

## 2. Arquitectura resultante

```text
                 ┌─────────────────────┐
cliente ───────▶ │ clientes  :8081     │
   │             └─────────────────────┘
   │             ┌─────────────────────┐        ┌──────────────────────────┐
   └───────────▶ │ pedidos   :8082     │ ─────▶ │ RabbitMQ                 │ ─ ─ ▶ logística
                 │ (caché de productos)│ evento │ cola pedidos-confirmados │   (conceptual)
                 └─────────────────────┘        └──────────────────────────┘
```

```text
.
├── ARQUITECTURA.md
├── README.md
├── monolito/              # sistema original, sin cambios
└── microservicios/
    ├── README.md          # endpoints, ejecución y evidencia
    ├── compose.yaml       # RabbitMQ
    ├── clientes/
    ├── pedidos/
    └── eventos/
```

## 3. Descomposición en microservicios

La separación se hizo **por dominio de negocio**:

| Servicio | Responsabilidad | Endpoints |
| --- | --- | --- |
| `clientes` | Alta y consulta de clientes. | `POST /clientes`, `GET /clientes/:id` |
| `pedidos` | Catálogo de productos y confirmación de pedidos. | `GET /productos`, `POST /pedidos` |

Productos quedó dentro de `pedidos` porque confirmar un pedido requiere validar que el producto exista. Si productos fuera un servicio aparte, cada pedido dependería de una llamada de red a otro servicio.

Cada microservicio es un **módulo Go independiente**, con su propio `go.mod`, proceso y puerto, y sus propios datos en memoria. No comparten código ni datos, así que se pueden ejecutar, modificar y desplegar por separado.

`pedidos` no consulta a `clientes` para validar el `cliente_id`. Así se evita el acoplamiento sincrónico entre servicios: si `clientes` está caído, `pedidos` sigue funcionando.

## 4. Organización interna por capas

Los dos servicios siguen la misma estructura:

```text
controllers/   -> recibe la request HTTP, valida el body y arma la respuesta
services/      -> reglas de negocio (validaciones, confirmar pedido, publicar evento)
repositories/  -> acceso a datos (en memoria)
models/        -> estructuras de datos
main.go        -> cableado de dependencias y arranque del servidor
```

- El flujo siempre es **controlador → servicio → repositorio**. El controlador no conoce el repositorio, y el repositorio no conoce HTTP.
- El servicio depende de **interfaces** (`ClientesRepo`, `ProductosRepo`, `PedidosRepo`, `Publisher`), no de implementaciones concretas. Por eso los datos en memoria se pueden reemplazar por una base de datos, o RabbitMQ por otro broker, sin tocar la lógica de negocio.
- Las dependencias se crean y se conectan en `main.go`.
- Los repositorios en memoria usan `sync.Mutex`, porque gin atiende requests en paralelo.

## 5. Caché de productos

`GET /productos` usa una caché en memoria implementada como **decorador del repositorio**:

```text
PedidoService ──▶ ProductosCache ──(miss)──▶ ProductosMemoria
```

- `ProductosCache` implementa la misma interfaz `ProductosRepo` que el repositorio real, así que el servicio no sabe si está usando la caché.
- Si el listado está guardado y no venció, lo devuelve (**CACHE HIT**). Si no, lo pide al repositorio siguiente, lo guarda y le asigna un vencimiento (**CACHE MISS**).
- TTL de **1 minuto**: es el tiempo máximo que un cambio en el catálogo podría tardar en verse. Con datos fijos en memoria no hay riesgo, pero con una base real sería el compromiso entre rendimiento y datos actualizados.

## 6. Comunicación asincrónica con RabbitMQ

Cuando se confirma un pedido, `pedidos` publica el evento `pedido.confirmado` en la cola `pedidos-confirmados`:

```text
POST /pedidos
  -> PedidoService valida y guarda el pedido
  -> publica pedido.confirmado
  -> cola pedidos-confirmados
  -> logística (conceptual) prepararía el envío
```

```json
{
  "tipo": "pedido.confirmado",
  "pedido_id": "PED-1",
  "cliente_id": "C-1",
  "producto_id": "P-1"
}
```

**Por qué un evento y no una llamada HTTP a logística:**

- `pedidos` no necesita saber quién procesa el pedido confirmado, ni esperar a que lo haga.
- Si logística no está disponible, el mensaje queda en la cola hasta que lo consuma. La cola se declara *durable*.
- Se pueden agregar nuevos interesados en el evento (facturación, notificaciones) sin modificar `pedidos`.

La conexión, la declaración de la cola y la publicación están aisladas en `pedidos/messaging/rabbitmq_publisher.go`. El servicio solo conoce la interfaz `Publisher`.

## 7. Decisiones y alcance

- **Carpeta `eventos`:** el modelo del evento se definió en `pedidos/models`, porque por ahora es el único servicio que lo usa. Moverlo a un paquete compartido obligaría a crear un tercer módulo y acoplaría los servicios. Tendría sentido si se implementara logística como consumidor.
- **Fuera de alcance, según la consigna:** consumidor de logística, base de datos real, autenticación, pagos, stock real, frontend, DLQ y reintentos.
- **Limitación conocida:** el pedido se guarda antes de publicar el evento. Si RabbitMQ falla justo en ese momento, el pedido queda guardado pero la API responde con error y el evento no se envía. En un sistema real se resolvería con reintentos o con el patrón *outbox*.