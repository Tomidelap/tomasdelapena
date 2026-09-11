# TP 1-4: de monolito a microservicios

**Alumno:** Tomás de la Peña

Separación del monolito de e-commerce en dos microservicios independientes, con caché para productos y publicación del evento `pedido.confirmado` en RabbitMQ.

```text
cliente -> clientes :8081
        -> pedidos  :8082

pedidos -- pedido.confirmado --> RabbitMQ (cola pedidos-confirmados) --> logística (conceptual)
```

## Estructura

```text
microservicios/
├── compose.yaml                 # RabbitMQ
├── clientes/
│   ├── controllers/
│   ├── services/
│   ├── repositories/
│   ├── models/
│   └── main.go                  # servidor :8081
└── pedidos/
    ├── controllers/
    ├── services/                # confirma el pedido y publica el evento
    ├── repositories/            # productos, caché de productos y pedidos
    ├── messaging/                # publisher de RabbitMQ
    ├── models/
    └── main.go                  # servidor :8082
```

## Microservicio `clientes` (puerto 8081)

| Método | Endpoint | Descripción |
| --- | --- | --- |
| `POST` | `/clientes` | Crea un cliente. Body: `{"nombre": "Ana Pérez"}`. Devuelve el cliente con su ID (`C-1`, `C-2`, ...). |
| `GET` | `/clientes/:id` | Devuelve el cliente o `404` si no existe. |

## Microservicio `pedidos` (puerto 8082)

| Método | Endpoint | Descripción |
| --- | --- | --- |
| `GET` | `/productos` | Lista los productos. Usa caché en memoria. |
| `POST` | `/pedidos` | Confirma un pedido. Body: `{"cliente_id": "C-1", "producto_id": "P-1"}`. |

**Caché:** `ProductosCache` envuelve al repositorio de productos. Si hay datos guardados y no pasó un minuto, los devuelve directo (log `CACHE HIT`). Si no, los pide al repositorio real, los guarda y reinicia el TTL (log `CACHE MISS`).

**Evento:** al confirmar un pedido se publica en la cola `pedidos-confirmados`:

```json
{
  "tipo": "pedido.confirmado",
  "pedido_id": "PED-1",
  "cliente_id": "C-1",
  "producto_id": "P-1"
}
```

Todo lo de RabbitMQ (conexión, cola, publicación) está en `pedidos/messaging/rabbitmq_publisher.go`. Logística no se implementa, es solo el destinatario conceptual del evento.

## Ejecución

1. Levantar RabbitMQ (desde esta carpeta):

```bash
docker compose up -d
```

2. Levantar cada microservicio en una terminal distinta:

```bash
cd clientes
go mod tidy
go run .
```

```bash
cd pedidos
go mod tidy
go run .
```

3. Probar (PowerShell):

```powershell
Invoke-RestMethod -Method Post -Uri http://localhost:8081/clientes -ContentType "application/json" -Body '{"nombre":"Ana Pérez"}'
Invoke-RestMethod http://localhost:8081/clientes/C-1
Invoke-RestMethod http://localhost:8082/productos
Invoke-RestMethod -Method Post -Uri http://localhost:8082/pedidos -ContentType "application/json" -Body '{"cliente_id":"C-1","producto_id":"P-1"}'
```

4. Verificar el evento en el panel de RabbitMQ: http://localhost:15672 (usuario `guest`, contraseña `guest`), cola `pedidos-confirmados`.

Más detalle de las decisiones de diseño en [`ARQUITECTURA.md`](../ARQUITECTURA.md).
