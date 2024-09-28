# chip8

This is my implementation of the [CHIP-8 interpreter](https://en.wikipedia.org/wiki/CHIP-8).

## Usage

```
Usage of chip8:
  -c string
        color of pixels (default "white")
  -k string
        type of keyboard (dvorak, qwerty) (default "dvorak")
  -l    color fill pixels (default true)
  -r int
        frame refresh rate (default 4)
```
## Examples

```bash
# set color to grey
$ chip8 -c=grey ./roms/ibm.ch8
```
![](/examples/ibm-logo.png?raw=true "IBM Logo")


```bash
# set keypad to qwerty and color to green
$ chip8 -k=qwerty -c=green ./roms/pong.ch8
```
![](/examples/pong.png?raw=true "Pong")

```bash
# no options
$ chip8 ./roms/ttt.ch8
```
![](/examples/ttt.png?raw=true "Tic-Tac-Toe")

```bash
# set color to red, pixel fill to false, and refresh rate to 8
$ chip8 -c=red -l=false -r=8 ./roms/tetris.ch8
```
![](/examples/tetris.png?raw=true "Tetris")
