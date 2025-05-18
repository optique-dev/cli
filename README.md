# Optique

## Installation

```bash
go install github.com/optique-dev/cli@latest
```

## Usage

Optique is a command line tool that allows you to build microservices at a lightning fast pace. The philosophy behind Optique is based on the [clean architecture principles](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).
Basically, your Optique project is a collection of :
- applications : protocols that are used to communication between microservices (e.g. HTTP, gRPC, AMQP, Kafka, etc.)
- infrastructures : backends that are used to store and process data (e.g. PostgreSQL, Redis, MongoDB, etc.)
- services : business logic that is used to process data (e.g. a user service, a product service, etc.) -> this is the part that you must implement yourself

Optique is made for people that wrote 10 000 times the same REST APIs and now want to be able to reuse there code in a more scalable way.

### Create a new project

```bash
optique init <project-name>
```
