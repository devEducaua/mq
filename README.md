
# mq

mq (music queue)

## build

just run:
```sh
make
```

for installing:
```sh
# this need root, use doas or sudo.
make install
```

## usage

basic usage:
```sh
# for seeing the current queue.
mq list

# to add a new song.
mq add example/song.mp3

# to see a directory in the mpd music directory.
mq see example

# write's the albumart to a file, and run a notification script everytime the song changes.
mq albumart --notify

# opens man page.
mq --help
```
for more information on usage consult the manpage.

# configuration

```conf
# the mpd address.
addr: localhost:6600

# the music directory of mpd
base-path: ~/music

# the default command that mq will run when no args are provided.
default-command: list

# where to save the currentsong album art.
cover-path: /tmp/cover.jpg

# the notify script to run
notify-path: ~/.config/mq/notify
```

for config and script example see ./examples.

