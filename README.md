# airthings-exporter

Prometheus exporter for [Airthings](https://www.airthings.com/) API

## Intro

This is a tool for exporting air quality information from
[Airthings](https://www.airthings.com) to
[Prometheus](https://prometheus.io) so Prometheus can be used for
monitoring and graphing air conditions.


With an [Airthings View Plus](https://www.airthings.com/view-plus),
this tool can monitor CO<sub>2</sub>, humidity, PM1 and PM2.5
particulate levels, barometric pressure, Radon, Temperature, and VOC
concentrations.  Other models are currently untested, but should
support a subset of these values.

This uses [Airthings's
API](https://developer.airthings.com/docs/api-getting-started/index.html),
rather than needing direct Bluetooth access to individual Airthings
devices.  Airthings provides rate-limited access for free for consumer
accounts, but only allows 120 requests per hour.  Individual devices
only publish updates every 2.5-5 minutes, so this should support
full-resolution queries for 5-10 devices using a free account.  This
should also work with paid commercial accounts, but it is untested.

## Building

You'll need to have a recent Go compiler installed.  Then do something
like this:

```
$ go checkout https://github.com/scottlaird/airthings-exporter.git
$ cd airthings-exporter
$ go get
$ go build cmd/airthings-exporter/airthings-exporter.go
```

This will leave a `airthings-exporter` binary in the current
directory.

## Running

First, make sure that you have already registered your Airthings
devices on their
[dashboard](https://dashboard.airthings.com/devices).  Then create an
API client ID via
https://dashboard.airthings.com/integrations/api-integration.  You'll
need to know both the client ID and the client secret; both are
accessible from that URL.

Then run `./airthings-exporter --client-id XXXX --client-secret YYYY`,
replacing XXXX and YYYY with the ID and secret values from Airthings'
dashboard.  If the client ID and secret are valid, then it'll listen
on port 8080 and export results via `/metrics`.

This will get a registered port number eventually, and startup scripts
will be provided.  Feel free to send a PR.