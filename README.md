# tfgrid monitoring bot

This is a bot to monitor the balance in accounts and send warnings if it is under some limit.

## How to start

- Create a new [telegram bot](README.md#create-a-bot-if-you-dont-have) if you don't have.
- Create a new env file `.env`, for example:

```env
TESTNET_MNEMONIC=<your mainnet mnemonic>
MAINNET_MNEMONIC=<your testnet mnemonic>
TFTS_LIMIT=70000
BOT_TOKEN=<your token>
CHAT_ID=<your chat ID>
MINS=<number of minutes between each message>
```

- Create a new json file `wallets.json` and add the list of addresses you want to monitor, for example:

```json
{ 
    "mainnet": ["<your address>"],
    "testnet": ["<your address>"] 
}
```

- Run
  
```bash
make build
bin/tfgridmon -e .env -w wallets.json
```

## Create a bot if you don't have

- Open telegram app
- Create a new bot
  
```ordered
1. Find telegram bot named "@botfarther"
2. Type /newbot
```

- Get the bot token
  
```ordered
1. In the same bot named "@botfarther"
2. Type /token
3. Choose your bot
```

- Get your chat ID

```ordered
1. Search for @RawDataBot and select Telegram Bot Raw from the drop-down list.
2. In the json returned, you will find it in section message -> chat -> id
```

## Test

```bash
make test
```
