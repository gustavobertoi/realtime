# Realtime Service

This is a service, written in Golang, designed to provide high-performance real-time capabilities using [WebSockets](https://developer.mozilla.org/pt-BR/docs/Web/API/WebSockets_API) or [Server Sent Events.](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events)

Leverage the power of channels to publish and subscribe to messages with maximum efficiency and high availability.

Overview:

- [Warning](#warning)
- [System design](#system-design-architecture)
- [Features](#features)
- [How to use with docker](#how-to-use-with-docker)

---

## Warning

This service has not yet been fully tested as it should and is still under development.

I don't advise anyone to use it in production yet.

---

## System Design (Architecture)

![system-design-image](./docs/realtime-service.png)

---

## Features

1. WebSocket Connectivity: Connect to channels and engage in two-way communication. Produce and consume messages in real-time through an intuitive WebSocket protocol interface.

2. Server-Sent Events (SSE): For those looking for a read-only subscription model, leverage the SSE protocol. This ensures a seamless and constant flow of updates without the need for continuous requests.

3. Integrations: By integrating with external services as Redis, this service can redistribute messages across multiple service instances. This not only provides persistence but ensures that scalability doesn't come at the expense of performance.

4. Backend Message Push: Equip your application or backend to send messages directly to the connected clients on a channel. This paves the way for dynamic and interactive user experiences.

5. Customization via YAML: Tailor the service to your needs using a YAML configuration file. This ensures that the service can adapt to a variety of use-cases without the need for hard-coded changes.

6. Docker & Kubernetes Friendly: Whether you're a fan of Docker or Kubernetes, this service has got you covered. Attach it using a volume, and you're set for a containerized, scalable realtime experience.

---

## How to use with docker

### Using memory as pub/sub

1. Write your yaml config (/configs/memory.config.yaml)

```yaml
server:
  allow_create_new_channels: true
  allow_push_server_messages: true
  render_chat_html: true
  render_notifications_html: true

pubsub:
  driver: "MEMORY"

channels:
  my-chat-example:
    id: "742fc7fe-1527-4184-8945-10b30bf01347"
    type: "WS"
    max_of_channel_connections: 100

  my-notification-example:
    id: "f71a6fb5-8962-48f4-a0dd-5cf60c42a70f"
    type: "SSE"
    max_of_channel_connections: 100
```

2. Write your docker-compose

```yaml
version: "3"

services:
  realtime-service:
    image: ghcr.io/open-source-cloud/realtime:latest
    container_name: realtime-service
    environment:
      APP_DEBUG: "1"
      CONFIG_FOLDER_PATH: "configs/memory.config.yaml"
    ports:
      - 8080:8080
    volumes:
      - ${PWD}/configs/memory.config.yaml:/configs/memory.config.yaml
    networks:
      - realtime-network

volumes:
  redis-volume:
    driver: local

networks:
  realtime-network:
    external: true
    driver: bridge
```

### Using redis as pub/sub

1. Write your yaml config (/configs/redis.config.yaml)

```yaml
server:
  allow_create_new_channels: true
  allow_push_server_messages: true
  render_chat_html: true
  render_notifications_html: true

pubsub:
  driver: "REDIS"
  redis:
    url: "redis://default:realtime@localhost:6379"

channels:
  my-chat-example:
    id: "742fc7fe-1527-4184-8945-10b30bf01347"
    type: "WS"
    max_of_channel_connections: 100

  my-notification-example:
    id: "f71a6fb5-8962-48f4-a0dd-5cf60c42a70f"
    type: "SSE"
    max_of_channel_connections: 100
```

2. Write your docker-compose

```yaml
version: "3"

services:
  redis:
    image: "bitnami/redis:latest"
    environment:
      REDIS_PASSWORD: toor
    ports:
      - 6379:6379
    volumes:
      - redis-volume:/bitnami/redis/data
    networks:
      - realtime-network

  realtime-service:
    image: ghcr.io/open-source-cloud/realtime:latest
    container_name: realtime-service
    depends_on:
      - redis
    environment:
      APP_DEBUG: "1"
      CONFIG_FOLDER_PATH: "configs/redis.config.yaml"
    ports:
      - 8080:8080
    volumes:
      - ${PWD}/configs/redis.config.yaml:/configs/redis.config.yaml
    networks:
      - realtime-network

volumes:
  redis-volume:
    driver: local

networks:
  realtime-network:
    external: true
    driver: bridge
```

---
