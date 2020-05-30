# system design.

## general.

a discord bot based on go-discord. plugins will provide the functionality of the application. we will initially ship with a `base` plugin. all configs MUST be dynamically read on change or at interval. plugins should be sucked into the application dynamically. bot should emit useful logs. bot MUST never log user chat.

## auth levels.

1. owner
1. admin
1. moderator
1. boost user (nitro)

## config

written in TOML, config file will contain:

- mapping of role to auth level.
- plugins.
  - mapping command to action to auth level.

## plugins.

### base.

`.echo` - echo's back a message.

`.whoami` - sends information about a user.

### antispam

protect users from spammy messages.

this plugin will require being able to read all message events sent to the server / observed by the bot.

### twitter.

idk.

## managing secrets.

mozilla sops? AH?!
