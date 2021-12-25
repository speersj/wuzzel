# wuzzel

wuzzel turns fuzzel into a window picker.

## Description

[fuzzel](https://codeberg.org/dnkl/fuzzel) is an application launcher for wlroots based Wayland compositors, similar to rofi's `drun` mode. I liked it immediately and wished it had a window picker of sorts, so I hacked one together real quick. Go is not my primary language so stuff may look funny, feel free to patch it up!

## Dependencies

* [sway](https://swaywm.org/) window manager
* [fuzzel](https://codeberg.org/dnkl/fuzzel) application launcher
* [go compiler](https://go.dev/) to build

### Executing program

wuzzel does not take any arguments of its own, but rather passes all arguments along to fuzzel. See fuzzel's man page for details.

Simplest example: `wuzzel`

More complex example with custom font and [dracula colors](https://draculatheme.com/):

```
wuzzel -f "JetBrainsMono Nerd Font":size=14 \
    -b 282a36df -t bd93f9ff -s 44475aff -S 8be9fdff -m ff79c6ff -C 50fa7bff \
    -w 120 -B 4 --line-height=24

```
You can bind wuzzel to a keyboard shortcut by modifying `$HOME/.config/sway/config`. For example, to bind to `mod+slash` add the following line to your sway config (arguments optional):

```
bindsym $mod+slash exec wuzzel -f "JetBrainsMono Nerd Font":size=14 \
    -b 282a36df -t bd93f9ff -s 44475aff -S 8be9fdff -m ff79c6ff -C 50fa7bff \
    -w 120 -B 4 --line-height=24
```

## License

This project is licensed under the MIT License - see LICENSE.md for details.
