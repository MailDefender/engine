# **engine**

**Simply the heart of the MailDefender project !**

> Note: this documentation is currently being drafted and will be completed in a future version.

## 📦 Prerequisites

- **Golang**
- **Docker** (optional)
- **[Environment Variables](#-configuration)**

## 🚀 Installation

### With Docker (Recommended)

```bash
docker build -t maildefender/engine .
docker run -p 8080:8080 --env-file .env maildefender/engine
```

### Without Docker

1. Clone the repository

```bash
# Clone the repoisitory
git clone https://github.com/MailDefender/engine.git
cd engine

# Install dependencies
go mod download

# Build
go build -o engine

# Run
source .env
./engine
```

## 🏃‍♂️ Usage

This app exposes APIs, so please refer to the Swagger to get more details about its usage.

## 🧪 Tests

Tests will be added soon.

## 🛠 Configuration

Create a .env file in the project root with the following variables:

```shell
# Delay between each loop
LOOP_DELAYS_SECS=5
# Set this field to true to ignore reputation verification and only sort messages
SKIP_REPUTATION_CHECK=false
# Caching period of rules before a refresh is required
RULES_REFRESH_PERIOD_SECS=300

# Database connection address
DATABASE_DNS=postgresql://user:password@host/db?sslmode=disable
# The URL on which the imap-connector component can be reached
IMAP_CONNECTOR_BASE_ENDPOINT=http://imap-connector:8081
# The URL on which the notifier component can be reached
NOTIFIER_BASE_ENDPOINT=http://notifier:8082
# The URL on which the validator component can be reached
VALIDATOR_PUBLIC_BASE_ENDPOINT=http://validator:8083

# Set this flag to true to enable daily recap
ENABLE_DAILY_RECAP=false
#Set this flag to true to move all received (and whitelisted) emails "today" to the "._Daily" maiblox
ENABLE_DAILY_MAILBOX=true
# Recipient of the daily recap
DAILY_RECAP_RECIPIENT=hello@me.com
```

## 📜 License

This project is licensed under MIT.
