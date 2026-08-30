GoBank CI/CD Kubernetes Deployment

🏦 **GoBank**

Cloud-Ready Banking Web Application with CI/CD, Docker, Kubernetes, and Terraform

## 📌 Overview

GoBank is a simple banking web application built with Golang and PostgreSQL.

The project demonstrates how a basic web application can be made deployment-ready using:

- Environment-based configuration
- PostgreSQL database persistence
- Docker containerization
- Docker Compose multi-container setup
- GitHub Actions CI/CD
- GitHub Container Registry
- Kubernetes deployment using Minikube
- Terraform infrastructure-as-code for future Azure AKS deployment

The application allows users to register, log in, deposit money, withdraw money, and view transaction history.

## 🚀 Features

- 👤 User registration
- 🔐 User login
- 🏦 Account dashboard
- 💰 Deposit funds
- 💸 Withdraw funds
- 📜 Transaction history
- 🗄️ PostgreSQL database storage
- 🐳 Dockerized Go application
- 🧩 Docker Compose setup for app + database
- ⚙️ GitHub Actions CI workflow
- 📦 Docker image publishing to GitHub Container Registry
- ☸️ Kubernetes deployment with ConfigMaps and Secrets
- ❤️ Health check endpoint for Kubernetes probes
- 🏗️ Terraform configuration for Azure AKS infrastructure

## 🏗️ System Architecture

```text
Developer
    ↓
GitHub Repository
    ↓
GitHub Actions CI
    ↓
Docker Build and Publish
    ↓
GitHub Container Registry
    ↓
Kubernetes / Minikube
    ↓
GoBank Application Pods
    ↓
PostgreSQL Service
    ↓
PostgreSQL Pod
⚙️ Application Flow
When a user opens the application, the Go server redirects them to the login page.
A new user can create an account by entering their name, date of birth, gender, username, and password. The application validates the form, checks if the username already exists, and stores the user in PostgreSQL with an initial balance of zero.
After logging in, the user is redirected to the dashboard. The dashboard shows user details and current account balance.
Users can deposit or withdraw money. Each transaction updates the account balance and creates a record in the transactions table.
Users can also view their transaction history, which displays deposits and withdrawals ordered by the latest transaction first.
🧰 Tech Stack
- Backend: Golang
- Database: PostgreSQL
- Frontend: HTML, CSS, Go Templates
- Containerization: Docker
- Local Orchestration: Docker Compose
- CI/CD: GitHub Actions
- Container Registry: GitHub Container Registry
- Kubernetes: Minikube, kubectl
- Infrastructure as Code: Terraform
- Cloud Target: Azure AKS configuration
📁 Project Structure
bank/
│── .github/
│   └── workflows/
│       ├── ci.yml
│       └── docker-publish.yml
│
│── db/
│   └── init.sql
│
│── k8s/
│   ├── gobank-configmap.yaml
│   ├── gobank-deployment.yaml
│   ├── gobank-secret.yaml
│   ├── gobank-service.yaml
│   ├── postgres-configmap.yaml
│   ├── postgres-deployment.yaml
│   ├── postgres-secret.yaml
│   └── postgres-service.yaml
│
│── static/
│   └── styles.css
│
│── templates/
│   ├── dashboard.html
│   ├── history.html
│   ├── login.html
│   └── register.html
│
│── terraform/
│   ├── main.tf
│   ├── variables.tf
│   ├── outputs.tf
│   └── versions.tf
│
│── Dockerfile
│── docker-compose.yml
│── go.mod
│── go.sum
│── main.go
│── README.md
🔐 Environment Variables
The application uses environment variables instead of hardcoded configuration.
DB_HOST
DB_PORT
DB_USER
DB_PASSWORD
DB_NAME
APP_PORT
Default local values:
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=Postgresql
DB_NAME=bankapp
APP_PORT=8080
⚡ Setup Instructions
1️⃣ Clone Repository
git clone https://github.com/Srimanvikonda/gobank-ci-cd-cloud.git
cd gobank-ci-cd-cloud
2️⃣ Run With Docker Compose
docker compose up --build
The application will be available at:
http://localhost:8080
Docker Compose starts:
- GoBank application container
- PostgreSQL container
- Database initialization using db/init.sql
- PostgreSQL health check before app startup
3️⃣ Test Health Endpoint
http://localhost:8080/health
Expected response:
OK
🗄️ Database
The PostgreSQL database contains two main tables:
users
transactions
The schema is stored in:
db/init.sql
This file automatically creates the required tables when the PostgreSQL container starts.
🐳 Docker
The project includes a multi-stage Dockerfile.
The Dockerfile:
- uses a Golang image to build the application
- creates a compiled Go binary
- copies the binary into a smaller Alpine image
- copies templates and static files
- runs the GoBank server on port 8080
Build manually:
docker build -t gobank-app .
Run manually:
docker run --name gobank-app `
  -p 8080:8080 `
  -e DB_HOST=host.docker.internal `
  -e DB_PORT=5432 `
  -e DB_USER=postgres `
  -e DB_PASSWORD=Postgresql `
  -e DB_NAME=bankapp `
  -e APP_PORT=8080 `
  gobank-app
⚙️ CI/CD Pipeline
The project uses GitHub Actions for automation.
CI Workflow
The CI workflow runs on push and pull requests.
It performs:
- Checkout repository
- Set up Go
- Download dependencies
- Build Go application
Docker Publish Workflow
The Docker Publish workflow runs on push to the main branch.
It performs:
- Login to GitHub Container Registry
- Build Docker image
- Tag Docker image
- Push image to GHCR
Docker image:
ghcr.io/srimanvikonda/gobank-ci-cd-cloud:latest
☸️ Kubernetes Deployment
The application is deployed locally using Minikube.
The Kubernetes setup includes:
- GoBank Deployment
- GoBank Service
- PostgreSQL Deployment
- PostgreSQL Service
- ConfigMaps
- Secrets
- Readiness probe
- Liveness probe
Apply Kubernetes manifests:
kubectl apply -f k8s/
Check pods:
kubectl get pods
Access the app locally:
kubectl port-forward service/gobank-service 8080:80
Open:
http://localhost:8080
❤️ Health Checks
The application includes a health endpoint:
/health
It returns:
OK
Kubernetes uses this endpoint for:
- Readiness probe
- Liveness probe
This helps Kubernetes decide whether the app is ready to receive traffic and whether it should be restarted.
🏗️ Terraform
The terraform/ folder contains infrastructure-as-code for future Azure AKS deployment.
It includes:
- Azure Resource Group
- Azure Kubernetes Service
- AKS node pool
- Managed identity
- Network profile
⚠️ Note: Terraform files are included for cloud-readiness, but they were not applied to avoid Azure costs.
Do not run:
terraform apply
unless you have an Azure subscription or credits and understand possible charges.
📊 Output
The project demonstrates:
- Running GoBank locally with PostgreSQL
- Running GoBank using Docker Compose
- CI validation through GitHub Actions
- Docker image publishing to GHCR
- Kubernetes deployment using Minikube
- Health monitoring through Kubernetes probes
- Cloud-ready infrastructure design using Terraform
📈 Project Status
Completed:
- ✅ GoBank web application
- ✅ PostgreSQL database integration
- ✅ Environment variable configuration
- ✅ Database initialization script
- ✅ Dockerfile
- ✅ Docker Compose setup
- ✅ GitHub Actions CI
- ✅ Docker image publishing to GHCR
- ✅ Kubernetes manifests
- ✅ Minikube deployment
- ✅ Health endpoint
- ✅ Kubernetes readiness/liveness probes
- ✅ Terraform AKS infrastructure code
🔮 Future Improvements
- Password hashing using bcrypt
- User sessions instead of query parameters
- Input validation improvements
- Unit and integration tests
- Prometheus and Grafana monitoring
- Azure AKS deployment using Terraform
- HTTPS and ingress controller
- Persistent volume for PostgreSQL in Kubernetes
- Helm chart for deployment packaging
👩‍💻 Author
Srimanvi Konda
📜 License
This project is developed for learning, DevOps practice, and portfolio demonstration.
⭐ Summary
GoBank demonstrates how a Golang web application can be transformed into a deployment-ready DevOps project by combining:
- Docker containerization
- PostgreSQL database persistence
- GitHub Actions CI/CD
- Container registry publishing
- Kubernetes orchestration
- Terraform infrastructure-as-code

After saving this as `README.md`, run:

```powershell
cd "D:\vs code\bank"

& "C:\Program Files\Git\cmd\git.exe" add README.md
& "C:\Program Files\Git\cmd\git.exe" commit -m "Add detailed project README"
& "C:\Program Files\Git\cmd\git.exe" push
