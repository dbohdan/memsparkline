# memsparkline

Display the memory usage (resident set size) of a process and its children in real time.  Print the average and maximum usage after it exits.


## Examples

```none
> ./memsparkline chromium-browser --incognito http://localhost:8081/ 
▁▁▁▁▄▇▇▇█ 789.53
max: 789.53
avg: 371.04
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


## License

MIT.
