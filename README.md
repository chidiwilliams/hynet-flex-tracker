# hynet-flex-tracker

Tracks data usage on the [MTN HyNetflex](https://mtnbusiness.com.ng/hynetflex) device.

The project contains two packages: the daemon and client.

The daemon records the data usage on a file. (It would eventually be configurable as a background process (on startup?). But for now, the binary must be run manually.)

The client reads the file and displays the tracked data.

## Installation

Pre-built binaries are not yet available on this package, so Go is required to build and install.

To build:

```shell script
go build -o hftrackerd github.com/chidiwilliams/hynet-flex-tracker/daemon

go build -o hftracker github.com/chidiwilliams/hynet-flex-tracker/client
```

## Usage

- Run the built binary:

```shell script
./hftrackerd
```

- Enter your device password and the tracking frequency.

The device password is the password you use to log in to the `192.168.0.1` interface.

The default and recommended value for the tracking frequency is "1d". See [time.ParseDuration](https://golang.org/pkg/time/#ParseDuration) for all valid values.

To read the recorded data, run `./hftracker`.

## Notes

- Sometimes the USSD requests temporarily fail with the following error: `Oops, looks like the code you used was incorrect. Please check and try again.` However, it usually works correctly again in the next run.