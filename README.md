# Smithy Backend [WIP]

> Smithy is an admin dashboard written in Go. It is designed to support multiple existed architectures and databases.

## Table of Contents

- Prerequisites
- Installation
- Quick start

### Prerequisites

**Disclaimer**: smithy works best on macOS and Linux.

- git should be installed
- docker and docker-compose must been installed
- golang version >= 1.10

### Basic Installation

##### Manual Installation

**1. Clone the repository**

    git clone https://github.com/dwarvesf/smithy

**2. Load environment variables**

    make local-env

**3. Build**

    make build

**4. Start database or clear data and permisstion**

    make local-db

**5. Start agent ( PORT 3000 )**

    make up-agent

**6. Set permission for each table**

    bin/smithy generate user

Note: If user existed. Try it

    bin/smithy generate user -f

**6. Start dashboard ( PORT 2999 )**

    make up-dashboard

**Note:** From here. You just need make up-dashboard to start server. Because agent's data have been saved into your local PC

### Quick start

**Start swagger API**

    http://localhost:2999/swaggerui/

Note: You need to follow step by step on swagger interface

**Step 1: Login**

When you enter your username and password to login. You will get a token.

**Step 2: Fill your token into the header**

    Authorization : BEARER "your token here"

**Step 3: Synchronized**

Run agent-sync endpoint to sync agent's data for dashboard

**Step 4: Run**

Run CRUD, config version with available form
