RentalFlow: Peer-to-Peer Rental Platform
RentalFlow is a distributed, microservices-based P2P rental platform designed for high scalability and real-time user interaction. The system supports renting diverse items such as vehicles and equipment through a decoupled, event-driven architecture.




🚀 Project Overview
This repository contains both the frontend and backend components of the RentalFlow system.


Frontend: A responsive React 18 application built with TypeScript and Tailwind CSS.


Backend: A suite of Go (Golang) microservices handling authentication, inventory, bookings, payments, and notifications .



Infrastructure: Orchestrated via Docker Compose with MongoDB, Redis, and RabbitMQ .


🏗️ System Architecture
1. Backend Microservices (Go)

API Gateway: The single entry point for request routing, load balancing, and JWT validation .


Auth Service: Manages user identities, secure password hashing, and role-based access control .


Inventory Service: Handles CRUD operations and advanced searching for rental items, optimized with Redis caching .


Booking Service: Manages the booking lifecycle and publishes events to RabbitMQ .


Payment Service: Integrates with the Chapa Payment Gateway for transaction processing and webhooks .


Notification Service: Consumes RabbitMQ events to deliver real-time email and in-app notifications .

2. Frontend (React + TS)

Architecture: Component-based design using React Context for global state and React Router for navigation .


Key Features: Kanban-style booking dashboard, interactive item browsing, and real-time notification panels .

🛠️ Technology Stack

Languages: Go, TypeScript, HTML/CSS.


Databases: MongoDB (Primary), Redis (Caching/Sessions).


Messaging: RabbitMQ for asynchronous Pub/Sub communication.



Media: Cloudinary for optimized image storage and delivery .


Documentation: OpenAPI 3.0 (Swagger) .

🚦 Getting Started
Prerequisites
Docker and Docker Compose.

Node.js (for local frontend development).

Go 1.21+ (for local backend development).

Installation & Deployment
Clone the Repository:

Bash

git clone https://github.com/xlxuxs/RentalFlowForSCM.git
cd RentalFlowForSCM
Environment Setup: Copy the .env.example file to .env and fill in your secrets (JWT, API keys, MongoDB URLs) .

Run with Docker:

Bash

docker-compose up --build
This command starts all microservices, databases, and the message broker simultaneously .

🧪 Testing & Quality Assurance

Integration Tests: Run the test_api.sh script in /tests/integration-tests/ to validate API contracts.


E2E Validation: Postman collections are available in the /tests folder to simulate full user flows (Registration → Booking → Payment) .


CI/CD: GitHub Actions automate testing on every pull request .

📄 Documentation
Detailed SCM documentation, including the Configuration Item (CI) Register and SCM Plan, can be found in the /docs directory. For API specifics, refer to the Swagger UI hosted at /docs/openapi.yaml
