
# mq

mq (music queue)

# build

just run:
```sh
make
```

for installing:
```sh
# this need root, use doas or sudo.
make install
```

# usage

basic usage:
```sh
# for seeing the current queue.
mq list

# to add a new song.
mq add example/song.mp3

# to see a directory in the mpd music directory.
mq see example

# opens man page.
mq --help
```

for more information on usage consult the manpage.
