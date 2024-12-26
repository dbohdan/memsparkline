complete -c memsparkline -s h -l help -d 'Print help message and exit'
complete -c memsparkline -s v -l version -d 'Print version number and exit'
complete -c memsparkline -s d -l dump -d 'File to append memory usage history to' -r
complete -c memsparkline -s l -l length -d 'Sparkline length' -r -a "20 40 60 80"
complete -c memsparkline -s m -l mem-format -d 'Format string for memory amounts' -r -a "%.0f %.1f %.2f"
complete -c memsparkline -s n -l newlines -d 'Print new sparkline on new line'
complete -c memsparkline -s o -l output -d 'Output file to append to' -r
complete -c memsparkline -s q -l quiet -d 'Do not print sparklines'
complete -c memsparkline -s r -l record -d 'How frequently to record memory usage in ms' -r -a "100 200 500 1000"
complete -c memsparkline -s s -l sample -d 'How frequently to sample memory usage in ms' -r -a "100 200 500 1000"
complete -c memsparkline -s t -l time-format -d 'Format string for run time' -r -a "%d:%02d:%04.1f"
complete -c memsparkline -s w -l wait -d 'Set sample and record time simultaneously' -r -a "100 200 500 1000"

# Complete with available commands.
complete -c memsparkline -n __fish_is_first_token -a "(__fish_complete_command)"
