<div align="right">

[**English**](./README.md)  |  [**فارسی**](./README.fa.md)

</div>

# Developer Guide

Welcome! This guide explains how to set up your environment, build, and run the **bgscan** project locally.

> **Before you start:** please read the [Contributing Guide](./contributing.md) first.

## Table of Contents

- [Project Architecture](#project-architecture)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
  - [1. Clone the repository](#1-clone-the-repository)
  - [2. Create a branch](#2-create-a-branch)
  - [3. Install dependencies](#3-install-dependencies)
  - [4. Run the project](#4-run-the-project)
- [Building Releases](#building-releases)
  - [Install the builder](#install-the-builder)
  - [Build commands](#build-commands)
  - [Building for Android](#building-for-android)
- [bgscan-builder reference](#bgscan-builder-reference)

---

## Project Architecture

```
.
├── assets
│   └── xray
│       └── outbounds
├── bgscan-builder
├── build.sh
├── cmd
│   └── bgscan
│       └── main.go
├── docs
│   ├── developer
│   ├── images
│   └── user
├── go.mod
├── images
├── internal
│   ├── core
│   │   ├── config
│   │   ├── dns
│   │   ├── fileutil
│   │   ├── ip
│   │   ├── iplist
│   │   ├── process
│   │   ├── result
│   │   ├── scanner
│   │   └── xray
│   ├── logger
│   ├── startup
│   └── ui
│       ├── components
│       ├── main
│       ├── shared
│       └── theme
├── ips
├── LICENSE
├── README.md
├── scripts
│   ├── install-builder.ps1
│   ├── install-builder.sh
│   ├── install-deps.ps1
│   ├── install-deps.sh
│   ├── install.ps1
│   ├── install.sh
│   ├── publish.sh
│   └── release.sh
└── settings
```

A quick overview of the main directories:

| Path | Description |
|---|---|
| `cmd/bgscan` | Application entry point (`main.go`). |
| `internal/core` | Core logic: config, DNS, file utilities, IP handling, process management, results, scanner, and Xray integration. |
| `internal/logger` | Logging utilities. |
| `internal/startup` | Startup health checks (config, DNS, logger, Xray). |
| `internal/ui` | UI components, main views, shared elements, and theming. |
| `assets/xray/outbounds` | Bundled Xray outbound assets. |
| `ips` | Default IP range lists per provider (Cloudflare, AWS, Azure, etc.). |
| `settings` | Default `.toml` settings files. |
| `scripts` | Install/build/release helper scripts. |
| `docs` | User and developer documentation. |

---

## Prerequisites

- [Go](https://go.dev/) (see `go.mod` for the required version)
- Git
- For Android builds: [Android NDK](https://developer.android.com/ndk)

---

## Getting Started

### 1. Clone the repository

```bash
git clone <repo-url>
cd bgscan
```

### 2. Create a branch

```bash
git checkout -b feature/my-change
```

### 3. Install dependencies

bgscan uses a companion tool called **`bgscan-builder`** to fetch and build the project's dependencies. You don't need to install it manually — the install scripts below will download it for you, place it in the project root, and use it to fetch the correct dependency build for your OS/architecture.

From the project root, run:

**Linux / macOS**
```bash
./scripts/install-deps.sh
```

**Windows**
```powershell
./scripts/install-deps.ps1
```

This script will:
1. Download `bgscan-builder` into the project root.
2. Run `bgscan-builder setup-dev --project-dir <project-root>` to download the correct dependencies for your OS/arch and place them in the right directory.

### 4. Run the project

Once dependencies are installed:

```bash
go mod tidy
go run ./cmd/bgscan/
```

---

## Building Releases

To build release artifacts, you'll also need `bgscan-builder`. If you don't already have it from the dependency step above, install it with:

**Linux / macOS**
```bash
./scripts/install-builder.sh
```

**Windows**
```powershell
./scripts/install-builder.ps1
```

### Build commands

```bash
bgscan-builder release -os linux -arch amd64
bgscan-builder release -os android -arch arm64 -ndk-dir /opt/android-ndk
bgscan-builder release -os all -arch all -dest ./dist
```

### Building for Android

Android builds require the Android NDK. Pass its path with `-ndk-dir`:

```bash
bgscan-builder release -os android -arch arm64 -ndk-dir /opt/android-ndk
```

---

## bgscan-builder reference

Builds release artifacts for one or more OS/architecture combinations.

**Examples:**

```bash
bgscan-builder release -os linux -arch amd64
bgscan-builder release -os android -arch arm64 -ndk-dir /opt/android-ndk
bgscan-builder release -os all -arch all -dest ./dist
```

**Flags:**

| Flag | Description |
|---|---|
| `-arch string` | Target architecture (`amd64`, `arm64`, `arm32`, `amd32`, `all`) |
| `-dep-version string` | Dependencies version tag (default `"v1.0"`) |
| `-dest string` | Release output directory (default `"./dist"`) |
| `-ndk-dir string` | Android NDK root directory |
| `-os string` | Target operating system (`linux`, `windows`, `macos`, `android`, `all`) |
| `-project-dir string` | Path to the bgscan project |
| `-xray-version string` | Xray version tag (default `"v26.3.27"`) |
