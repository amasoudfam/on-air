# Viper Package

The Viper package is a popular configuration management library for Go. It provides a convenient and flexible way to work with configuration files and settings in your applications.

## Installation

To install the Viper package, use the following command:

```bash
go get -u github.com/spf13/viper
```

### Usage

To use the Viper package in your Go project, import it as follows:

go
Copy code

```go
import "github.com/spf13/viper"
```

### Features

The Viper package offers the following features:

* Loading configuration from multiple sources such as JSON, YAML, TOML, and more
* Reading configuration values using a hierarchical key system
* Setting default configuration values
* Automatic environment variable binding
* Support for command-line flags
* Handling of different data types for configuration values
* Watching and reloading configuration files
* Encryption and decryption of sensitive configuration values

## Why we use viper package

We use Viper in this project for easy and flexible configuration management. Viper simplifies the process of loading and accessing configuration values, allowing us to configure our application using various sources such as files, environment variables, and more. With Viper, we can easily retrieve configuration values and manage different configurations for different environments, making our application more adaptable and easier to maintain.
