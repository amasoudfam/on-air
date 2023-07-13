# Breaker Package

The Breaker package is a powerful implementation of the circuit breaker pattern for Go. It provides a simple and flexible way to improve the resilience and reliability of distributed systems by preventing cascading failures.

## Installation

To install the Breaker package, use the following command:

```
go get -u github.com/eapache/go-resiliency/breaker
```

## Usage

To use the `Breaker` type in your Go project, import the package as follows:

```go
import "github.com/eapache/go-resiliency/breaker"
```

## Features

The Breaker package offers the following features:

- Implementation of the circuit breaker pattern
- Configurable parameters such as failure threshold, timeout duration, and backoff strategy
- Reporting of successful and failed requests
- Check the current state of the circuit breaker (i.e., open, closed, or half-open)
- Support for retrying failed requests

## Why we use Breaker package

We use the Breaker package in this project for improving the resilience and reliability of our distributed systems. The circuit breaker pattern is a powerful design pattern that helps prevent cascading failures by detecting and handling errors that might occur when a service or resource fails. With the `Breaker` type, we can easily wrap a function or method call and monitor its success and failure rates over time. The package's configurable parameters, such as failure threshold and timeout duration, allow us to fine-tune the circuit breaker's behavior to match our application's needs. Using Breaker package helps us improve the robustness of our distributed systems and prevent cascading failures.