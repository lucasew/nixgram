# nixgram ![cachix:lucasew-github](https://img.shields.io/github/workflow/status/lucasew/nixgram/Cachix?label=cachix%3Alucasew-github&logo=github&style=flat-square)

Run custom system commands as a telegram bot

When you run nixgram it needs the bot token and the user id of the administrator.

To discover your user id you can pass any number and send it a message. It will show your user id in the logs.

It is engineered with safety in mind. Only the administrator can send commands and all commands are prefixed with `nixgram-`.

Ex: `hey` looks for `nixgram-hey` in its `$PATH`.

It was developed with nixos in mind, you can setup a new command using `writeShellScriptBin` but also write manually scripts to a folder that is in `$PATH`.

There is nothing that isn't cross platform in the code so non-linux platforms should work as well.
