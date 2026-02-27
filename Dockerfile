# -----------------------------
# BUILD STAGE
# -----------------------------
FROM golang:1.25 AS builder

WORKDIR /app

# After WORKDIR /app
# /
# ├── ...
# └── app/        <-- current working directory (empty)


# -----------------------------------
# Copy go module files
# -----------------------------------
COPY go.mod go.sum ./

# Now filesystem looks like:
# /app
# ├── go.mod
# └── go.sum


# -----------------------------------
# Download dependencies
# -----------------------------------
RUN go mod download

# -----------------------------------
# Copy application source folders
# -----------------------------------
COPY cmd ./cmd/
COPY internal ./internal/

# Now filesystem looks like:
# /app
# ├── go.mod
# ├── go.sum
# ├── cmd/
# │   ├── main.go
# │   └── (other files if any)
# └── internal/
#     ├── orders/
#     ├── products/
#     └── adapters/
#         └── postgresql/
#             └── sqlc/


# -----------------------------------
# Build the application
# -----------------------------------
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping ./cmd/
# /
# ├── docker-gs-ping   <-- compiled static binary
# └── app/
#     ├── go.mod
#     ├── go.sum
#     ├── cmd/
#     └── internal/


# =====================================================
# FINAL STAGE (LIGHTWEIGHT RUNTIME IMAGE)
# =====================================================
FROM alpine:latest AS final

WORKDIR /app

# After WORKDIR
# /
# └── app/   <-- empty


# -----------------------------------
# Copy ONLY the compiled binary
# -----------------------------------
COPY --from=builder /docker-gs-ping .

# Final filesystem now:
# /app
# └── docker-gs-ping
#
# ❌ No go.mod
# ❌ No go.sum
# ❌ No cmd/
# ❌ No internal/
# ❌ No source code
#
# Only the compiled binary exists.


EXPOSE 8080

CMD ["./docker-gs-ping"]

# docker build -t ecom-go-api-project .

# docker run -d --env-file .env ecom-go-api-project

# docker ps

# docker inspect <container_id> | grep IPAddress

# use ip on postmant / apiDog