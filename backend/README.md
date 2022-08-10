# Materialize Postgres Mock Service

## Building the service

Depending on your preferences, you can either build the service directly using Golang or use the existing [Dockerfile](./Dockerfile) to build an image:

#### Using GoLang:

```bash
go build
```

Then run the service:

```bash
./backend
```

#### Using Docker:

```bash
docker build -t mz-mock .
```

Run the container:

```bash
docker run -p 6875:6875 -d mz-mock
```