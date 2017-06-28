Minecraft Server Wrapper – `mcserv`
============================================

A wrapper around running an (especially modded) Minecraft server.

Why?
-----

Most people when running a dedicated (Linux) Minecraft server usually do so
either wrapped in a  terminal multiplexer like `tmux` or `screen`, or just on
its own. Maybe with an init script to keep it alive should it crash.

But doing it either of these ways pose some disadvantages:

- Running in a multiplexer makes it difficult to add to an init system, as the
session may not always exit with a bad code when the underlying script dies.

- Running in a multiplexer can also be annoying from a scripting perspective:
dealing with byzantine `tmux send-keys` et al,and escape sequences, and having
no good way of capturing output.

- Running on its own makes it difficult to administrate, as there is no way
to control the server aside from feeding commands via `stdin` or in-game chat.

The `mcserv` tool solves these issues by wrapping the actual server script and
exposing an RPC socket which can be controlled via a simple JSON interface. It
can be configured to exit with a bad code when the underlying server does, or
attempt to keep it alive itself.

Development
--------------

### Requirements

- Make
- [Glide](https://github.com/Masterminds/glide)
- [Gox](https://github.com/mitchellh/gox)


### Building

By default, uses Gox to build for linux and macOS 64-bit architectures. This can
be configured with the `XC_OS` and `XC_ARCH` variables at compile time:

```bash
$ XC_OS="windows darwin" XC_ARCH="amd64 386" make
```

Compiled binaries can be found in the `build/<OS>/<ARCH>/` dirs.

License
---------

mcserv
Copyright (C) 2017 Joshua Lindsey

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.

