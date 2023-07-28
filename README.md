# memsparkline

Track the RAM usage ([resident set size](https://en.wikipedia.org/wiki/Resident_set_size)) of a process, its children, its children's children, etc. in real time with a Unicode text [sparkline](https://en.wikipedia.org/wiki/Sparkline).  See the average and the maximum usage after the process exits, as well as the runtime.


## Examples

```none
> memsparkline -- chromium-browser --incognito http://localhost:8081/
▁▁▁▁▄▇▇▇█ 789.53
 avg: 371.04
 max: 789.53
time: 0:00:12.345
```

```none
> memsparkline -o foo command &
> tail -f foo
```


## Dependencies

Python 3.7 or later, [psutil](https://github.com/giampaolo/psutil).

### Debian/Ubuntu

```sh
sudo apt install python3-psutil
```

### DragonFly BSD 6.6, FreeBSD 13.1

```sh
sudo pkg install py39-psutil
```

### NetBSD 9.3

```sh
sudo pkgin in py310-psutil
```

### OpenBSD

```sh
doas pkg_add py3-psutil
```


## License

MIT.
