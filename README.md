# MC Run

`mcrun` is a command-line interface (CLI) utility for creating Minecraft servers. 
The servers are run inside Docker containers and use images built by the [`minecraft-server-docker`](https://github.com/dm0275/minecraft-server-docker) repository.

## Features

- **Easy Server Creation**: Quickly spin up Minecraft servers using Docker images.
- **Customizable**: Pass in server configurations and manage multiple instances with ease.

## Prerequisites
* [Docker](https://docs.docker.com/get-docker/)

## Installation & Usage
To use `mcrun`, you can download the built binaries from the latest release.

1. Download the binary:

   ```bash
   curl -L -o mcrun https://github.com/dm0275/mcrun/releases/download/<version>/mcrun-<os>-<arch>
   
   # Example
   curl -L -o mcrun https://github.com/dm0275/mcrun/releases/download/v0.0.2/mcrun-darwin-amd64
   ```

2. Make the binary executable:

   ```bash
   chmod +x mcrun
   ```

3. Run the CLI to create a new Minecraft server:

   ```bash
   ./mcrun forge start --version forge-1.20.1 --world-name <world_name>
   ```

This command will pull the appropriate Docker image and create a new Minecraft server instance.


### Available Commands
TODO

### Building from Source

#### Prerequisites

Building the CLI from source, you will need the following installed:

- [Golang](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Gradle](https://gradle.org/install/) (with Kotlin DSL support)

1. Clone the repository:

   ```bash
   git clone https://github.com/dm0275/mcrun.git 
   ```

2. Build the CLI using Gradle:

   ```bash
   ./gradlew build
   ```

3. The built CLI will be available in the `build` directory:

   ```bash
   ./build/mcrun-<os>-<arch>
   ```
