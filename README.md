# memsparkline

Display the memory usage (resident set size) of a process and its children in real time.

## Example

```none
> ./memsparkline chromium-browser --incognito http://localhost:8081/ 
▁▁▁▁▄▇▇▇█ 789.53
max: 789.53
avg: 371.04
```

## Dependencies

### Debian/Ubuntu

```sh
sudo apt install python3-psutil
```

## License

MIT.
