# telegram-middleman-bot

## Deployment

```sh
$ git clone https://github.com/cubarco/telegram-middleman-bot.git
$ cd telegram-middleman-bot
$ now
```

## Telegram Webhook

```sh
$ curl https://api.telegram.org/bot{YOUR_BOT_TOKEN}/setWebhook?url=https://{YOUR_DOMAIN}/api/updates
```

## Usage

```sh
$ curl -v -X POST -d '{"chat_id":"{YOUR_CHAT_ID}","text":"_test_","key":"{YOUR_KEY}"}' https://{YOUR_DOMAIN}/api/messages
```
