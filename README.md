# memsparkline

Track the RAM usage ([resident set size](https://en.wikipedia.org/wiki/Resident_set_size)) of a process, its children, its children's children, etc. in real time with a Unicode text [sparkline](https://en.wikipedia.org/wiki/Sparkline).  See the average and the maximum usage after the process exits.


## Examples

```none
> ./memsparkline chromium-browser --incognito http://localhost:8081/ 
▁▁▁▁▄▇▇▇█ 789.53
avg: 371.04
max: 789.53
```

```none
> ./memsparkline -o foo command &
> tail -f foo
```


## Dependencies

Python 3.5 or later, [psutil](https://github.com/giampaolo/psutil).

### Debian/Ubuntu

```sh
sudo apt install python3-psutil
```

### FreeBSD 12

```sh
sudo pkg install py37-psutil
```


## License

MIT.
