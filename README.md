# Fullstack Reservation System

This project is a fullstack microservices-based reservation system, including a frontend and multiple backend services (`reserva`, `pagamento`, `bilhete`, `marketing`, `session`), along with RabbitMQ and MongoDB for messaging and data persistence.

---

## 🧰 Services

- **Frontend**: React + Vite
- **Backend Microservices**:
  - `reserva` (Reservation service)
  - `pagamento` (Payment service)
  - `bilhete` (Ticket service)
  - `marketing` (Marketing service)
  - `session` (Session management service)
- **RabbitMQ**: Message broker
- **MongoDB**: Database

---

## 🚀 Getting Started

To run the entire system, make sure you have [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/) installed.

### 🏁 Run the project

```bash
docker compose up --build
```

This command will:

- Build all Docker images
- Start all services (frontend, microservices, RabbitMQ, MongoDB)
- Map the necessary ports to your local machine

---

### 🌐 Access the frontend

Once all containers are up and running, access the application at:

👉 [http://localhost:5173](http://localhost:5173)

---

### 🐰 Access RabbitMQ dashboard

RabbitMQ provides a management interface that you can access at:

👉 [http://localhost:15672](http://localhost:15672)

**Default credentials:**

- **Username**: `guest`
- **Password**: `guest`
