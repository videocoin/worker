[![Build Status](https://ci.videocoin.io/api/badges/VideoCoin/go-videocoin/status.svg)](https://ci.videocoin.io/VideoCoin/go-videocoin)

# go-videocoin

Official golang implementation of the VideoCoin client.

## build

Building `vid` requires both a Go (version 1.9 or later) and a C compiler.

```bash
$ make vid
```

## run

Going through all the possible command line flags is out of scope here (please consult Ethereum's
[CLI Wiki page](https://github.com/ethereum/go-ethereum/wiki/Command-Line-Options)).

As an alternative to passing the numerous flags to the `vid` binary, you can also pass a configuration file via:

```bash
$ vid --config /path/to/your_config.toml
```

To get an idea how the file should look like you can use the `dumpconfig` subcommand to export your existing configuration:

```bash
$ vid dumpconfig
```

## Versioning

This repository uses [Semantic Versioning 2.0.0](https://semver.org/) to determine when and how the version changes.
We're currently in the initial development phase (0.y.z).
