# heshi

declare -x http_proxy=xxx
declare -x https_proxy=xx
gb vendor restore
gb build
./bin/heshi_service