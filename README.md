# Telegram Bot for Checking URL

This repository contains source code for a Go-based Telegram bot that performs checks for URL, SSL certificates, HTTP(S) headers, and redirects in websites.

## Usage

1. Clone the repository

```sh
https://github.com/Vivirinter/telegram-bot-url.git
```
2. Navigate to the source code directory

```sh
cd telegram-bot-url
```
3. Build Docker image

```sh 
docker build -t telegram-bot-url .
```

4. Run Docker container

```sh 
docker run -p 8080:8080 -e TELEGRAM_TOKEN=<TOKEN> telegram-bot-url
```

## Environment Variables

* TELEGRAM_TOKEN: The Telegram Bot Token, required to use the bot in the Telegram application. Get your token by creating a bot using the BotFather in Telegram.

## Commands


You can use the following commands within the Telegram bot:
- `/start:` This command starts the bot.
- `/help:` This command shows a help message with the list of commands that can be used.
- `/check URL:` This command can be used to perform HTTPS and certificate checks for a specific URL. Replace URL with the actual URL to be checked.


## Contributions

Contributions are welcome. Please fork the project and submit a pull request with your changes.