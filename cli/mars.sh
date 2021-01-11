#!/bin/bash

echo ""

# 取得執行檔的所在位置，若是Link模式時必須取得實際位置
SELF_PATH=${BASH_SOURCE[0]}
if [[ -L "$SELF_PATH" ]]; then
    DIR=$(dirname $(readlink ${SELF_PATH}))
else
    DIR=$(dirname ${SELF_PATH})
fi

if [ $DIR == "." ]; then
    DIR=$PWD
fi

source $DIR/tools/utils.sh

cd $DIR

cmd=$1

case "$cmd" in
env)
    cat <<EOF
PATH: $DIR
LINK: ${BASH_SOURCE[0]}
PWD: $PWD

EOF
    exit 0
    ;;
start)
    repo=$2
    if [[ "$(docker images -q godev_${repo} 2> /dev/null)" == "" ]]; then
        echo "build dev images godev_${repo}"
        docker build -t godev_${repo} -f "${repo}.Dockerfile" .
    fi
    
    cd $PWD/../$repo
    docker-compose up -d

    echo "Initail Schema..."
    for sqlScript in $DIR/migration/mysql/schema/*.sql; do
        cat "$sqlScript" | docker exec -i db mysql -uroot -pmypassword mysql 
        echo "$sqlScript done"
    done
    
    echo "${repo} DONE"
    echo "${repo} Public Addr:${RESTORE}"
    showPublicDomain
    ;;
stop)
    repo=$2
    cd $PWD/../$repo
    docker-compose down 
    ;;
godev)
    repo=$2
    doTelepresence "go" "$repo"
    ;;
link)
    rm /usr/local/bin/mars
    ln -s $PWD/mars.sh /usr/local/bin/mars
    ;;
help)
    marshelp
    ;;
*)
    echo "${LRED}該指令不存在唷～${RESTORE}"
    ;;
esac
