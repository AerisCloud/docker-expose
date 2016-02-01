# docker-expose

Small internal tool at Wizcorp that allows developers to expose their docker
containers to the internet for testing.

## CLI

```
NAME:
   docker-expose - Expose docker containers on aeris.cd

USAGE:
   docker-expose [global options] command [command options] [arguments...]

VERSION:
   0.1.0

AUTHOR(S):
   Wizcorp Inc

COMMANDS:
   add		Add a container to the expose server
   list		List currently running containers and their expose status
   rm		Remove a container from the expose server
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h		show help
   --version, -v	print the version
```

### docker-expose add CONTAINER_NAME/ID

Add a container to the expose server

#### Options

`--port, -p "80"	The port of the container to expose, default to 80`

### docker-expose list

List currently running containers and their expose status

### docker-expose rm CONTAINER_NAME/ID

Remove a container from the expose server

## License

MIT
