**Twilio Segment Webserver**
==========================

A simple GoLang web server designed to receive events from Segment
Intended to be used as a quick destination for events

**Requirements**
---------------

* GoLang 1.24 (tested version)


**Configuring Environment Variables**
-------------------------------

Copy `.env.example` to a directory named `static` in the `/path/to/your/content` directory referenced above

- `listen_port` is used to configure the port inside the container that the server will listen to.
  - Update this value to match your `docker run` command
  - Defaults to `8080`


**Example Use Cases**
--------------------

* Act as a Segment data source
* Allow review of received events


**Usage**
-----

### Running the Server (in a container)

To build the docker container, navigate to the project directory and execute:
```bash
docker build -f Dockerfile -t stevennic/segment-server:latest ./
```

Once built, run the container with the following:
```bash
docker run \
--name=segment-server-test \
--volume=/etc/timezone:/etc/timezone:ro \
--env=TZ=America/New_York \
--volume=/path/to/your/content:/usr/src/app/content \
--workdir=/usr/src/app \
-p 8085:8080 -p 8085:8080/udp --expose=8080 \
--restart=unless-stopped \
--detach=true \
stevennic/segment-server:latest
```

### Running the Server (locally)

Access directory
Example:
```
cd ~/path/to/Github/repo
go run ./
```

### Accessible resources

`/events/create` is where events are ingested

`/events/list` will show a list of events since the server has been up

`/content/*` is an exposed directory. Files stored under the configured local directory `/content/direct/` will be available here. Useful for testing Segment source scripts.

`/ws/client` exposes a websocket client to connect to live events

`/ws` endpoint to connect to websocket

`/ws/health` and `/ws/clients` give a JSON output of websocket health and connected clients.


**Example Directory Structure**
--------------------
*Default directory in container is `/usr/src/app/content`*

*Fallsback to `./content` if not found/outside of a container*

* content/
  * logs/
    * segment-*env*.log
    * ...
  * static/
    * .env
  * direct/
    * ajs.html
    * ...


**Troubleshooting**
---------------

If you encounter any issues while running the server, check the following:

* Check the server logs for errors or warnings.
* Make sure to restart the server if necessary.
