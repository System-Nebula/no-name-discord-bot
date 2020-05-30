# system design.

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
