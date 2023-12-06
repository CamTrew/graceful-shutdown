# Graceful Shutdown Example

## Introduction
This repository provides a solution to gracefully shutting down an API and its dependencies. For this example, we use Gin as our web framework, Rabbit and Mongo as our dependencies, and Docker Compose to deploy.

## How To Run
Run the application with the following command:
```
docker-compose up --build -d
```

## Further Reading
To read more about graceful shutdowns in Go and this code example, visit the blog post https://thegodev.com/graceful-shutdown
