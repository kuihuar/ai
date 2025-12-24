options=$(/usr/bin/getopt -o hln: --long help,list,name: -- "$@")
if [ $? -ne 0 ]; then
  /usr/bin/echo "Error: Invalid options provided." >&2
  usage
  exit 1
fi
eval set -- "$options"
while true; do
    case "$1" in
        -h|--help)
            HELP=true
            shift
            ;;
        -l|--list)
            LIST=true
            shift
            ;;
            --)
            shift
            break
            ;;
        *)
            echo "Unknown option: $1" >&2
            exit 1
            ;;
    esac
done

