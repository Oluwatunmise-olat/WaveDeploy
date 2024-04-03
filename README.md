### WaveDeploy

Wave Deploy is a proof of concept that demonstrates the steps involved in building a very simple PAAS.

[![Watch the video](https://www.loom.com/share/2c07aca6d5424afdb853f052358256a5?sid=39d7ac25-c6ab-480b-9af3-6d22b4eb26f1)](https://www.loom.com/share/2c07aca6d5424afdb853f052358256a5?sid=39d7ac25-c6ab-480b-9af3-6d22b4eb26f1)

## Documentation

For detailed documentation and usage instructions, please visit [here](https://oluwatunmise.gitbook.io/wave-deploy/).

## Prerequisites

Before getting started, make sure you have the following prerequisites installed on your system:

- [Go programming language](https://golang.org/doc/install)
- MySQL

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/Oluwatunmise-olat/WaveDeploy
```

### 2. Configure Environment Variables (Create a github app)

Create a .env file in the root directory of the project and mirror .env.example:

```
MYSQL_CONNECTION_STRING=
ENVIRONMENT=<PRODUCTION, STAGING>
APP_KEY=<ANY KEY>
GITHUB_APP_ID=
GITHUB_APP_PUBLIC_LINK="https://github.com/apps/<YOUR APP>"
CONNECT_TO_GITHUB_URL=
GITHUB_CLIENT_ID=
GITHUB_CLIENT_SECRET=
PORT=
GITHUB_APP_WEBHOOK_SECRET=
```

### 3. Run Migrations

```bash
make db-migrate-up mysql_username=<YOUR USERNAME> mysql_password=<YOUR PASSWORD>
```

### 4. Up and Running

```bash
# Cli
go run main.go
# Start HTTP Server to receive webhooks
go run main.go -serve-http true (Configure url on github)
```

## TODO

- [ ] Domain Name Mapping in Caddy Configuration (_priority 3_).
- [ ] Optimize deployment for varying deployment type (SPA AND API) (_priority 1_).
- [ ] Auto-Scale Based on metrics from prometheus (would be a background worker) (_priority 5_).
- [ ] PR Preview Link (_priority 5_).

## How to Contribute

If you would like to make any updates/improvements to the application, it is greatly appreciated and welcome. Here's how you can contribute:

1. Fork the repository and create a new branch for your changes.
2. Work on the task(s) you'd like to help with.
3. Once you're done, submit a pull request with your changes.
4. Feel free to reach out to us if you have any questions or need further clarification.

Thank you for your support and contributions!
