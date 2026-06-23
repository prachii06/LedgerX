# Deployment Diagram

## Local Development

```text
+--------------------------------------+
| Docker Compose                       |
|                                      |
|  +-------------------------------+   |
|  | API Service                   |   |
|  +-------------------------------+   |
|                                      |
|  +-------------------------------+   |
|  | Worker                        |   |
|  +-------------------------------+   |
|                                      |
|  +-------------------------------+   |
|  | Kafka                         |   |
|  +-------------------------------+   |
|                                      |
|  +-------------------------------+   |
|  | Redis                         |   |
|  +-------------------------------+   |
|                                      |
|  +-------------------------------+   |
|  | PostgreSQL                    |   |
|  +-------------------------------+   |
+--------------------------------------+
```

## Production (MVP)

All services are deployed on a single VPS using Docker Compose.

Future versions may use Kubernetes for orchestration and horizontal scaling.