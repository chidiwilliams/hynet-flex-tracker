# hynet-flex-tracker

Tracks data usage on the [MTN HyNetflex](https://mtnbusiness.com.ng/hynetflex) device.

## Installation

Pre-built binaries are not yet available on this package, so Go is required to build and install.

```shell script
go build ./...
```

## Usage

Run the built binary:

```shell script
./hynet-flex-tracker
```

Enter your device password and the tracking frequency.

The device password is the password you use to log in to the `192.168.0.1` interface.

The default and recommended value for the tracking frequency is "1d". See [time.ParseDuration](https://golang.org/pkg/time/#ParseDuration) for all valid values.  

## Notes

- Sometimes the USSD requests temporarily fail with the following error: `Oops, looks like the code you used was incorrect. Please check and try again.` However, it usually works correctly again in the next run.