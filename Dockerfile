# -------Stage 1 : build the binary ------
FROM golang:1.26 AS builder

WORKDIR /app

# Copy dependecy files first and download - this layer is called,
# so deps only re-download when go.mod/go.sum change
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source and compile to a single binary
COPY .  .
RUN CGO_ENABLED=0 GOOS=linux go build -o notesapi .

# ------- Stage 2 : tiny final image --------
FROM alpine:latest

WORKDIR /app
# copy only the binary from stage 1
COPY --from=builder /app/notesapi . 

EXPOSE 8080
CMD ["./notesapi"]