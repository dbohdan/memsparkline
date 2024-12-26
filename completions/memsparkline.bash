_memsparkline() {
    local cur prev opts

    COMPREPLY=()
    cur=${COMP_WORDS[COMP_CWORD]}
    prev=${COMP_WORDS[COMP_CWORD - 1]}
    opts='-h --help -v --version -d --dump -l --length -m --mem-format -n --newlines -o --output -q --quiet -r --record -s --sample -t --time-format -w --wait'

    case "${prev}" in
    -d | --dump | -o | --output)
        COMPREPLY=($(compgen -f -- "${cur}"))
        return 0
        ;;
    -l | --length)
        COMPREPLY=($(compgen -W "20 40 60 80" -- "${cur}"))
        return 0
        ;;
    -r | --record | -s | --sample | -w | --wait)
        COMPREPLY=($(compgen -W "100 200 500 1000" -- "${cur}"))
        return 0
        ;;
    -m | --mem-format)
        COMPREPLY=($(compgen -W "%.0f %.1f %.2f" -- "${cur}"))
        return 0
        ;;
    -t | --time-format)
        COMPREPLY=($(compgen -W "%d:%02d:%04.1f" -- "${cur}"))
        return 0
        ;;
    esac

    if [[ ${cur} == -* ]]; then
        COMPREPLY=($(compgen -W "${opts}" -- "${cur}"))
        return 0
    fi

    # Complete with commands if no other completion applies.
    COMPREPLY=($(compgen -c -- "${cur}"))
}

complete -F _memsparkline memsparkline
