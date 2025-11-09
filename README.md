# MC Run

`mcrun` is a command-line interface (CLI) utility for creating Minecraft servers. 
The servers are run inside Docker containers and use images built by the [`minecraft-server-docker`](https://github.com/dm0275/minecraft-server-docker) repository.

## Features

- **Easy Server Creation**: Quickly spin up Minecraft servers using Docker images.
- **Customizable**: Pass in server configurations and manage multiple instances with ease.
- **HTTP API**: Launch an API layer and drive server lifecycle actions from other tools or a UI.
- **Managed Mods**: Continue dropping mods into `modsDir` manually or have `mcrun` fetch them straight from CurseForge.

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

```
mcrun is a command-line interface (CLI) utility for creating Minecraft servers.

Usage:
  mcrun [command]

Available Commands:
  api         Start the HTTP API server.
  fabric      Configures a Minecraft Fabric server instance.
  forge       Configures a Minecraft Forge server instance.
  help        Help about any command
  servers     Manage Minecraft servers.
  setup       Setup Minecraft server directory structure Forge server
  vanilla     Configures a Minecraft server (vanilla) instance.
  version     Display the current version of the mcrun CLI.
```

### Listing servers

Use the CLI to see which worlds exist under the mcrun home directory:

```bash
./mcrun servers list
```

The output shows runtime status (running/stopped), type, version, mod count, compose status, and filesystem path for every world. Each server directory contains a `metadata.json` file that stores this information and is reused by the API listing endpoint.

## HTTP API

Launch the API server with:

```bash
./mcrun api --host 127.0.0.1 --port 8080
```

Once running, the following endpoints are available:

### `GET /healthz`

Simple liveness probe.

```bash
curl http://127.0.0.1:8080/healthz
```

```json
{"status":"ok"}
```

### `POST /servers`

Create (or start) a Minecraft server. All JSON fields are optional besides `worldName`; unspecified values fall back to the defaults for the selected `type` (`vanilla`, `forge`, or `fabric`).

```bash
curl -X POST http://127.0.0.1:8080/servers \
  -H 'Content-Type: application/json' \
  -d '{
        "worldName": "my-forge-world",
        "type": "forge",
        "version": "forge-1.20.1",
        "maxMemory": "4G",
        "enableCmdBlock": true,
        "mountDirs": ["resourcepacks"]
      }'
```

Typical response:

```json
{
  "status": "starting",
  "type": "forge",
  "worldName": "my-forge-world"
}
```

Common error responses:

| HTTP Code | Description                            |
|-----------|----------------------------------------|
| 400       | Invalid JSON or unsupported server type |
| 500       | Directory/compose generation failures   |

### `POST /servers/{worldName}/stop`

Stops a running server without deleting any files.

```bash
curl -X POST http://127.0.0.1:8080/servers/my-forge-world/stop
```

Responses:

```json
{
  "status": "stopping",
  "worldName": "my-forge-world"
}
```

Error responses include:

| HTTP Code | Description                    |
|-----------|--------------------------------|
| 404       | Compose file not found for world |
| 500       | Docker compose stop failure     |

### `DELETE /servers/{worldName}`

Deletes the server identified by `worldName`. The API attempts to stop the server first and then removes the entire world directory along with the generated compose file.

```bash
curl -X DELETE http://127.0.0.1:8080/servers/my-forge-world
```

Responses:

```json
{
  "status": "deleted",
  "worldName": "my-forge-world"
}
```

Error responses include:

| HTTP Code | Description                    |
|-----------|--------------------------------|
| 404       | Server directory not found      |
| 500       | Stop/delete failure             |

### `POST /servers/{worldName}/start`

Starts an existing server (using its Docker compose file).

```bash
curl -X POST http://127.0.0.1:8080/servers/my-forge-world/start
```

Responses:

```json
{
  "status": "starting",
  "worldName": "my-forge-world"
}
```

Error responses include:

| HTTP Code | Description                    |
|-----------|--------------------------------|
| 404       | Compose file not found for world |
| 500       | Docker compose start failure    |

### `GET /servers`

Lists the directories detected under `~/.mcrun` along with metadata such as runtime status, server type, mods declared, and whether a Compose file exists.

```bash
curl http://127.0.0.1:8080/servers
```

Response:

```json
[
  {
    "worldName": "my-forge-world",
    "path": "/Users/me/.mcrun/my-forge-world",
    "status": "running",
    "hasCompose": true,
    "metadata": {
      "worldName": "my-forge-world",
      "type": "forge",
      "version": "forge-1.20.1",
      "port": "25565",
      "mods": [
        {
          "source": "curseforge",
          "curseforge": { "projectId": 238222, "gameVersion": "1.20.1" }
        }
      ],
      "createdAt": "2024-07-15T17:00:00Z",
      "updatedAt": "2024-07-15T17:05:00Z"
    }
  }
]
```

## Mod Management

You can keep copying `.jar` files directly into the server’s `mods` folder, or you can ask `mcrun` to download them from CurseForge during provisioning.

1. Create a CurseForge API key and set it as an environment variable:

   ```bash
   export CURSEFORGE_API_KEY=<your-key>
   ```

2. Include one or more `projectID[@gameVersion]` entries when starting servers from the CLI. If you omit `@gameVersion`, `mcrun` infers it from the server version (e.g., `forge-1.20.1` ⇒ `1.20.1`) and selects the newest compatible file. 

   ```bash
   ./mcrun forge start \
     --world-name my-forge-world \
     --curseforge-mod 238222@1.20.1 \
     --curseforge-mod 306612
   ```

3. Or supply mods via the API:

   ```json
   {
     "worldName": "my-forge-world",
     "type": "forge",
     "mods": [
       {
         "source": "curseforge",
         "curseforge": { "projectId": 238222, "gameVersion": "1.20.1" }
       }
     ]
   }
   ```

Existing mods inside the `mods` directory are left untouched; remote downloads are only added if the target file is missing. Downloaded files are cached under `~/.mcrun/cache/mods` and hard-linked (or copied) into each world, so repeated server creations reuse the cached artifacts.

## Web UI (React)

A starter React + Vite frontend lives in `ui/`. It proxies API requests to the Go server and provides forms to create, stop, or delete servers.

```bash
cd ui
npm install   # installs dependencies
npm run dev   # starts Vite on http://localhost:5173
```

Make sure the API (`./mcrun api --host 127.0.0.1 --port 8080`) is running so the UI can reach it via the dev server proxy (`/api` → `127.0.0.1:8080`).

### Building from Source

#### Prerequisites

Building the CLI from source, you will need the following installed:

- [Golang](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Gradle](https://gradle.org/install/)

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
