# Arquitectura

**Alumno:** Tomás de la Peña

## Punto de partida

El monolito tenía todo junto: clientes, productos y pedidos en el mismo proceso, escuchando en el `8080`. La idea era separar eso en dos servicios sin perder funcionalidad.

```text
cliente -> clientes  :8081
        -> pedidos   :8082

pedidos -- pedido.confirmado --> RabbitMQ --> logística (no se implementa)
```

## Cómo separé los servicios

Quedaron dos microservicios, cada uno con su propio `go.mod` y su propio puerto:

- **clientes**: alta y consulta de clientes (`POST /clientes`, `GET /clientes/:id`).
- **pedidos**: productos y confirmación de pedidos (`GET /productos`, `POST /pedidos`).

Los productos los dejé dentro de `pedidos` en vez de armar un tercer servicio, porque para confirmar un pedido hay que chequear que el producto exista. Si eso fuera una llamada a otro servicio, cada pedido dependería de que ese servicio esté arriba, y no vale la pena esa complicación para este TP.

`pedidos` no le pregunta nada a `clientes` cuando confirma un pedido (no valida que el `cliente_id` exista de verdad). Así, si `clientes` se cae, `pedidos` sigue andando igual.

## Capas

Los dos servicios siguen más o menos la misma organización:

```text
controllers/   -> recibe el request, valida lo básico y arma la respuesta
services/      -> la lógica (confirmar pedido, publicar el evento, etc.)
repositories/  -> guarda y lee los datos (en memoria, con un mapa)
models/        -> los structs
main.go        -> arma todo y levanta el servidor
```

El controlador no toca el repositorio directamente, siempre pasa por el service. Y el repositorio no sabe nada de HTTP. Los repos usan un `sync.Mutex` porque gin puede atender varios requests al mismo tiempo y si no, se puede romper el mapa.

## Caché de productos

`GET /productos` tiene una caché en memoria bastante simple: cuando la piden, se fija si ya guardó una respuesta hace menos de un minuto. Si sí, la devuelve directo (log de `CACHE HIT`). Si no, la vuelve a pedir al repositorio real, la guarda y arranca de nuevo el contador (`CACHE MISS`).

La armé como un wrapper que cumple la misma interfaz que el repositorio de productos, entonces el service ni se entera si está usando la caché o el repo real.

Un minuto de TTL es medio arbitrario, pero como los datos son fijos en memoria no importa demasiado — la idea es mostrar que la caché funciona, no ajustar el tiempo perfecto para un caso real.

## RabbitMQ

Cuando `POST /pedidos` confirma un pedido, se guarda el pedido y después se publica el evento a la cola `pedidos-confirmados`:

```json
{
  "tipo": "pedido.confirmado",
  "pedido_id": "PED-1",
  "cliente_id": "C-1",
  "producto_id": "P-1"
}
```

Uso un evento en vez de llamar directamente a logística porque `pedidos` no tiene por qué saber quién consume eso ni esperar a que lo procesen. Si en algún momento se agrega logística de verdad (o facturación, notificaciones, lo que sea), se puede sumar como otro consumidor de la cola sin tocar `pedidos`.

Todo lo de la conexión a RabbitMQ, declarar la cola y publicar está en `pedidos/messaging/rabbitmq_publisher.go`. El service solo usa una interfaz `Publisher`, no depende de la librería de amqp directamente.

## Qué no hice (a propósito)

Según la consigna no hacía falta: autenticación, frontend, pagos, stock real, base de datos real, ni el consumidor de logística. Tampoco DLQ, reintentos ni nada de eso para la cola.

Un punto flojo que dejé así: el pedido se guarda antes de publicar el evento. Si en ese instante falla RabbitMQ, el pedido queda guardado pero la respuesta es un error y el evento nunca sale. Para este TP no lo resolví, pero en un caso real se podría manejar con reintentos o con un patrón tipo outbox.
