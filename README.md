<div align="center">
    <h1>Pulse</h1>
    <p>
        <a href="https://pulse.octara.xyz/api/discord"><img src="https://img.shields.io/discord/1330635185590632489?color=5865F2&logo=discord&logoColor=white" alt="Discord server" /></a>
        <a href="https://github.com/octarahq/pulse/commits/main"><img alt="Last commit" src="https://img.shields.io/github/last-commit/octarahq/pulse?logo=github&logoColor=ffffff" /></a>
        <a href="https://octara.xyz"><img src="https://img.shields.io/badge/Website-octara.xyz-blue.svg" alt="Website" /></a>
    </p>
</div>

## About

`Pulse` is a minimalist, modular grid-based TUI (Terminal User Interface) and Web dashboard generator (comming soon) designed for home servers and self-hosted environments. Built entirely in pure **Go**, it allows users to arrange customizable widgets inside a local coordinate grid with smart automatic border merging, high performance, and an ultra-low memory footprint.

You can run the service locally as a single lightweight binary, check your stats via SSH, or host the minimalist layout web-view on any tablet or secondary display.

> For an optimal experience, customize your layout using absolute positioning coordinates in a simple TOML configuration file.

## Features

- **Absolute Grid Layout**: Position your widgets precisely using `x`, `y`, `width`, and `height` properties. Use `-1` to let a widget automatically fill the remaining screen space.
- **Smart Border Merging**: No double-borders cluttering your terminal. Overlapping edge characters automatically merge into clean junction symbols (`┬`, `├`, `┼`, etc.).
- **Automatic Plugin Registry**: An extensible modular architecture. Adding new widgets is fully automated using Go factory registries—no complex switch cases required.
- **Native Components**: Comes packaged with out-of-the-box widgets including dynamic text displays, time-zone aware clocks, and a full-size grid calendar.
- **Zero Dependencies**: Compiles down into a single, standalone binary file with no external runtime, system tools, or database dependencies.

## TODO:
- [ ] Web dashboard frontend with live updates and interactive widgets.
- [ ] Additional widget types: graphs, system monitors, RSS feeds, etc.
- [ ] Theme support and customizable color schemes.
- [ ] User-friendly CLI for generating and validating configuration files.
- [ ] Made possible the creation of custom plugins by the community with a simple API and documentation.

## Configuration

Pulse is entirely configured through a simple `config.toml` file placed next to the execution binary. Here is an example setup:

```toml
dashname = "My Server Dashboard"
theme = "dracula" // Themes are coming soon!

[[widgets]]
type = "display"
x = 0
y = 0
width = 30
height = 3
value = "System Status: OK"

[[widgets]]
type = "calendar"
x = 0
y = 2
width = 26
height = 11
format = "full"
timezone = "Europe/Paris"

```

## Installation

Go 1.22 or newer is required to compile the source code. The engine utilizes pure Go data mapping libraries.

```sh
# Clone the repository
git clone https://github.com/octarahq/pulse.git
cd pulse

# Download and tidy up Go modules
go mod tidy

# Run the live application
go run .

# Build a standalone optimized production binary
go build -ldflags="-s -w" -o pulse .
./pulse

```

## Links

* [GitHub Repository](https://github.com/octarahq/pulse)
* [Discord support](https://octara.xyz/api/discord)
* [Octara HQ](https://octara.xyz)

## Help

Need help setting up your grid configuration or developing your own widget plugin? Ask your questions on our [Discord support server!](https://octara.xyz/api/discord)
