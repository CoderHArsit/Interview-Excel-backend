# InterviewExcel – Backend

Backend service for **InterviewExcel**, a platform that connects students with experts for mock interviews, mentorship, and career guidance.

Built with **Go**, **Gin**, **PostgreSQL**, and **GORM**, following a clean, scalable backend architecture.

---

## 🚀 Tech Stack

* **Language**: Go (Golang)
* **Framework**: Gin
* **Database**: PostgreSQL
* **ORM**: GORM
* **Authentication**: JWT (Access + Refresh tokens)
* **Logging**: Logrus
* **Architecture**: Layered (Controllers → Services → Repositories)
* **API Style**: REST
* **Deployment Ready**: Docker-friendly

---

## 📁 Project Structure

```
interviewexcel-backend-go/
├── config/             # DB, env, and app configuration
├── controllers/        # HTTP handlers (Gin controllers)
├── middlewares/        # Auth, request logging, guards
├── models/             # GORM models & repositories
├── routes/             # API route definitions
├── utils/              # Helpers (JWT, logging, time utils)
├── constants/          # App-wide constants & enums
├── migrations/         # DB migrations (GORM)
├── main.go             # App entry point
└── README.md
```

---

## 🔐 Authentication

* JWT-based authentication
* Separate flows for:

  * **Students**
  * **Experts**
* Tokens:

  * `access_token` (short-lived)
  * `refresh_token` (long-lived)

Authentication middleware injects:

```
user_uuid → gin.Context
```

---

## 👨‍🏫 Expert Features

* Expert profile creation & update
* Expertise, experience, fees, availability
* Weekly availability slot generation
* Slot cancellation & management

---

## 📅 Availability & Slots

### Slot Status Model

Slots use a **status-based lifecycle** instead of booleans.

```
available → booked → completed
available → cancelled
```

### Slot Status Values

```go
available
booked
cancelled
```

### Slot Model

```go
type AvailabilitySlot struct {
	ID        uint
	ExpertID  string
	Date      time.Time
	StartTime time.Time
	EndTime   time.Time
	Status    string
	StudentID *uint
}
```

---

## 🧾 API Highlights

### Generate Weekly Availability

```
POST /expert/availability/generate
```

### Cancel Slot

```
DELETE /expert/availability/:slot_id
```

### Update Expert Profile

```
PUT /expert/profile
```

---

## 🗄️ Database

* PostgreSQL
* GORM used for ORM and migrations
* Transactions used for multi-table updates
* Indexing on:

  * `expert_id`
  * `date`
  * `status`

---

## 📜 Logging

* Centralized logging using **Logrus**
* Caller information (file & line number)
* Structured error logging

Example:

```
level=error msg="failed to update expert profile" user_uuid=abc123 file=controller.go:42
```

---

## ⚙️ Environment Variables

Create a `.env` file:

```env
APP_ENV=development
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=interviewexcel

JWT_SECRET=supersecretkey
JWT_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=7d
```

---

## ▶️ Running Locally

### 1️⃣ Clone repository

```bash
git clone https://github.com/CoderHArsit/interviewexcel-backend-go.git
cd interviewexcel-backend-go
```

### 2️⃣ Install dependencies

```bash
go mod tidy
```

### 3️⃣ Run server

```bash
go run main.go
```

Server runs on:

```
http://localhost:8080
```

---

## 🧪 Testing (Planned)

* Unit tests for services
* Repository tests with test DB
* API tests using Postman / REST client

---

## 🛣️ Roadmap

* [ ] Student booking flow
* [ ] Payment integration
* [ ] Slot expiry automation
* [ ] Video interview (WebRTC – first-party)
* [ ] Notifications (Email / WhatsApp)
* [ ] Admin dashboard

---

## 🤝 Contributing

1. Fork the repo
2. Create a feature branch
3. Commit changes
4. Open a PR

---

## 👨‍💻 Author

**Harshit Saxena**
B.Tech CSE
InterviewExcel – Backend Engineering

---

## ⭐ Notes

* Designed for **scalability & interview readiness**
* Clean separation of concerns
* Easily extensible for microservices (gRPC-ready)

---

Just tell me 🚀

© 2025 Harshit Saxena. All rights reserved. This code is not licensed for use, redistribution, or modification.
