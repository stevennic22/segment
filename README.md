**Twilio Segment Webserver**
==========================

A simple GoLang web server designed to receive events from Segment
Intended to be used as a quick destination for events

**Requirements**
---------------

* GoLang 1.24 (tested version)

**Usage**
-----

### Running the Server

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

**Configuring Environment Variables**
-------------------------------

Copy `.env.example` to a directory named `static` in the `/path/to/your/content` directory referenced above

- `listen_port` is used to configure the port inside the container that the server will listen to.
  - Update this value to match your `docker run` command
- `db_enable` configures whether a connection to a database will be made
  - Currently only configured for the Allowed Origins/CORS configuration

**Event Format**
--------------

The server expects events in JSON format. The event structure is as follows:
```json
{
  "id": "",
  "event_type": "",
  "user_agent": "",
  "ip": ""
}
```
**Example Use Cases**
--------------------

* Act as a Segment data source
* Allow review of received events

**Troubleshooting**
---------------

If you encounter any issues while running the server, check the following:

* Check the server logs for errors or warnings.
* Make sure to restart the server if necessary.
