# Step 1: Using an official image of Go
FROM golang:1.19-alpine AS build

# Step 2: Installing dependencies for working with the PostgreSQL database
RUN apk add --no-cache postgresql-client

# Step 3: Creating directory of the app in the container
RUN mkdir /app

# Step 4: Configuring the working directory
WORKDIR /app

# Step 5: Copying source code into the container
COPY go.mod go.sum ./
COPY cmd/gophermart/main.go ./cmd/gophermart/

# Step 6: Installing the project dependencies
RUN go mod download

# Step 6: Building the app
RUN go build -o /app/cmd/gophermart ./cmd/gophermart

# Step 7: Using an official image of PostgreSQL
FROM postgres:13-alpine

# Step 8: Copying the binary file of the app into the container
COPY --from=build /app/main ./cmd/gophermart

# Step 9: Configuring enviromental variables for the database
ENV POSTGRES_HOST localhost
ENV POSTGRES_PORT 5432
ENV POSTGRES_USER gopher
ENV POSTGRES_PASSWORD gopher
ENV POSTGRES_DB gophermart

# Step 10: Launching the app at the start of the container
CMD ["/app/cmd/gophermart"]
