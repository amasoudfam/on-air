# Cobra Package

The Cobra package is a powerful command-line application framework for Go. It provides a simple and elegant way to create command-line interfaces (CLIs) with support for commands, flags, arguments, and more.

## Installation

To install the Cobra package, use the following command:

```bash

go  get  -u  github.com/spf13/cobra

```

### Usage

To use the Cobra package in your Go project, import it as follows:

go

Copy code

```go

import  "github.com/spf13/cobra"

```

### Features

The Cobra package offers the following features:

* Easy creation of CLI applications

* Support for single commands and nested commands

* Flag parsing and support for various flag types

* Arguments parsing

* Command and flag help text generation

* Automatic generation of help and usage information

* Support for subcommands and command aliases

* Extensibility through hooks and plugins

## Why we use cobra package

We use Cobra in this project for command-line interface (CLI) development. Cobra is a powerful and easy-to-use CLI library for Go. It provides a simple and intuitive way to define and handle commands, flags, and arguments in our application. With Cobra, we can quickly create robust CLI applications with features like subcommands, flags with default values, and help documentation. Using Cobra makes it easier for users to interact with our application through the command line and simplifies the development of CLI-based functionalities.

## What we do with this package

We utilize the Cobra package to create a command-line interface (CLI) tool with three commands: migrate, seed, serve, and worker.

* [Migrate Command](#migrate-command)

* [Seed Command](#seed-command)

* [Serve Command](#serve-command)

* [Worker Command](#worker-command)

## migrate command

The migrate command is used to handle database migrations in the application. Here's an explanation of how the migrate command works based on the provided code:

The migrateCmd variable represents the migrate command and is defined using the Cobra package. It has a short description and a long description explaining its purpose.

In the init() function, the migrateCmd command is added to the root command (rootCmd) so that it can be executed as a subcommand.

The migrateCmd command has two flags: state and steps. The state flag specifies the desired migration state (e.g., "up" or "down"), and the steps flag indicates the number of steps to migrate.

The Run function of the migrateCmd command is responsible for executing the migration logic. It retrieves the values of the state and steps flags and calls the migrateDB function to perform the migrations.

Inside the migrateDB function, the configuration is initialized using the provided configPath. Then, a Postgres database instance is initialized using the configuration.

The migrate.NewWithDatabaseInstance function creates a new migration instance with the specified migration files path ("file://migrations") and the Postgres database driver.

Based on the value of the state flag, the migrateDB function performs different migration actions. If the state is "up", it applies the migrations with mig.Up(). If the state is "down", it rolls back the migrations with mig.Down(). If the state is "drop", it drops all applied migrations with mig.Drop(). If the state is "steps", it migrates the specified number of steps with mig.Steps(steps).

After each migration action, a corresponding log message is printed to indicate the completion of the migration.

The seed command allows us to populate the database with initial or sample data. This is useful for setting up a development or testing environment.

## Seed command

The seed command is used to populate the database with initial data. Let's break down how the seed command works based on the provided code:

The seedCmd variable represents the seed command and is defined using the Cobra package. It has a short description and a long description explaining its purpose.

In the init() function, the seedCmd command is added to the root command (rootCmd) so that it can be executed as a subcommand.

The seedCmd command has one flag: fake. This flag is optional and can be used to indicate whether to fill tables with fake data for testing purposes.

The Run function of the seedCmd command is responsible for executing the seeding logic. It retrieves the value of the fake flag and calls the seed function to populate the database.

Inside the seed function, the configuration is initialized using the provided configPath, and a Postgres database instance is initialized based on the configuration.

Overall, the seed command allows you to populate the database with initial data, including user, country, city, flight, passenger, and ticket records. The optional fake flag can be used to populate tables with fake data for testing purposes.

## Serve command

The serve command is used to start the project's server and serve the application. Let's examine how the serve command works based on the provided code:  

The serveCmd variable represents the serve command and is defined using the Cobra package. It has a short description and a long description explaining its purpose.  

In the init() function, the serveCmd command is added to the root command (rootCmd) so that it can be executed as a subcommand.  

The serveCmd command has one flag: port. This flag is used to specify the port number on which the server should listen.  

The Run function of the serveCmd command is responsible for executing the server startup logic. It retrieves the value of the port flag and calls the startServer function to start the server.  

Inside the startServer function, the configuration is initialized using the provided configPath, and both a Postgres database instance and a Redis instance are initialized based on the configuration.  

If the port parameter is empty, the function retrieves the port number from the configuration.  

The server.SetupServer function is called with the initialized configuration, database instance, Redis instance, and the specified port. This function sets up the server and returns a http.Server instance.  

The log.Fatal function is then called with the http.Server instance to start the server and handle any potential errors. If an error occurs, it is logged and the program terminates.  

Overall, the serve command allows you to start the project's server and serve the application. It initializes the necessary configurations, creates database and Redis connections, and starts the server listening on the specified port or the default port from the configuration.  

## Worker command

The worker command is used to run a worker process that handles pending requests. Let's examine how the worker command works based on the provided code:  

The workerCmd variable represents the worker command and is defined using the Cobra package. It has a short description and a long description explaining its purpose.  

In the Run function of the workerCmd command, the SetupWorker function is called, passing the configFlag to initialize the worker.  

Inside the SetupWorker function, the configuration is initialized using the provided configPath. If the worker is disabled in the configuration, a log message is printed, and the function returns.  

If the worker is enabled, the Run function is called in a goroutine. It receives the worker configuration, a context, and a database connection.  

The Run function performs the main logic of the worker. It starts by initializing an APIMockClient (not shown in the provided code) and a ticker that triggers the worker's iteration based on the configured interval.  

A counter is used to track the number of iterations, and a loop runs continuously until the desired number of iterations is reached or the process is terminated.  

Within the loop, the GetExpiredTickets function is called to retrieve a list of expired tickets from the database using the repository package.  

For each expired ticket, a transaction is started using the db.Transaction function. The FindFlightById function is called to retrieve the associated flight information. Then, a refund request is made to the apiMock client (assuming it is properly implemented), passing the flight number and ticket count.  

If the refund operation encounters an error, the ticket status and payment status are updated to indicate expiration using the ChangeTicketStatus and ChangePaymentStatus functions from the repository package.

After processing each expired ticket, the loop continues until the ticker triggers the next iteration or the context is canceled.  

Once the loop is finished, a log message is printed to indicate the completion of the worker's execution.  

Overall, the worker command sets up and runs a worker process that handles pending requests. It initializes the necessary configurations and database connection, performs periodic iterations based on the configured interval, and processes expired tickets by refunding them if possible.
