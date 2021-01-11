RESTORE=$(echo -en '\033[0m')

RED=$(echo -en '\033[00;31m')
GREEN=$(echo -en '\033[00;32m')
YELLOW=$(echo -en '\033[00;33m')
BLUE=$(echo -en '\033[00;34m')
MAGENTA=$(echo -en '\033[00;35m')
PURPLE=$(echo -en '\033[00;35m')
CYAN=$(echo -en '\033[00;36m')

LIGHTGRAY=$(echo -en '\033[00;37m')
LRED=$(echo -en '\033[01;31m')
LGREEN=$(echo -en '\033[01;32m')
LYELLOW=$(echo -en '\033[01;33m')
LBLUE=$(echo -en '\033[01;34m')
LMAGENTA=$(echo -en '\033[01;35m')
LPURPLE=$(echo -en '\033[01;35m')
LCYAN=$(echo -en '\033[01;36m')

WHITE=$(echo -en '\033[01;37m')

function marshelp() {
    cat <<EOF

${CYAN}指令: galaxy {指令} [參數]${RESTORE}

env...............基本資訊(ex. galaxy 安裝路徑)
start.............啟動開發環境
     ${REPO_NAME}...........參數：mars start${REPO_NAME} 
stop..............關閉開發環境
godev.............Golang 開發模式，可以直接在 k8s 內部開發，不需要重新build image 之後還要用 kubectl 更新 pod
     ${MAGENTA}.............參數：galaxy godev${RESTORE} ${LGREEN}{REPO_NAME}${RESTORE}
     ${MAGENTA}.............範例：galaxy godev${RESTORE} ${LGREEN}crater${RESTORE}
     ${MAGENTA}.............產出：進入一個 Container，自動掛載專案，可以直接修改程式，在 Container 內執行${RESTORE} ${LGREEN}go build && ./REPO_NAME${RESTORE}
link.........軟連結
*.................顯示幫助

EOF
    exit 0
}

function addDomainToHost() {
    IFS=$'\r\n' GLOBIGNORE='*' command eval 'ETC_HOSTS=($(cat /etc/hosts))'
    IFS=$'\r\n' GLOBIGNORE='*' command eval 'GALAXY_HOSTS=($(cat $DIR/tools/hosts))'
    for domain in "${GALAXY_HOSTS[@]}"; do
        IS_EXIST_IN_ETC_HOSTS=$([[ ${ETC_HOSTS[*]} =~ (^|[[:space:]])"$domain"($|[[:space:]]) ]] && echo 'Y' || echo 'N')
        case $IS_EXIST_IN_ETC_HOSTS in
        N)
            echo "Insert '${LYELLOW}$domain${RESTORE}' to /etc/hosts"
            sudo echo $domain | sudo tee -a /etc/hosts 1>/dev/null
            ;;
        esac
    done
}

function removeDomainFromHost() {
    IFS=$'\r\n' GLOBIGNORE='*' command eval 'ETC_HOSTS=($(cat /etc/hosts))'
    IFS=$'\r\n' GLOBIGNORE='*' command eval 'GALAXY_HOSTS=($(cat $DIR/tools/hosts))'
    for domain in "${GALAXY_HOSTS[@]}"; do
        IS_EXIST_IN_ETC_HOSTS=$([[ ${ETC_HOSTS[*]} =~ (^|[[:space:]])"$domain"($|[[:space:]]) ]] && echo 'Y' || echo 'N')
        case $IS_EXIST_IN_ETC_HOSTS in
        Y)
            echo "Remove '${LYELLOW}$domain${RESTORE}' in /etc/hosts"
            sudo sed -i "" "/${domain}/d" /etc/hosts
            ;;
        esac
    done
}

function showPublicDomain() {
    IFS=$'\r\n' GLOBIGNORE='*' command eval 'GALAXY_HOSTS=($(cat $DIR/tools/hosts))'
    for domain in "${GALAXY_HOSTS[@]}"; do
        arr=($(echo $domain))
        host=${arr[1]}
        echo "http://$host"
    done
}

function doTelepresence() {
    lang=$1
    repo=$2

    # docker stop galaxy_telepresence
    # docker rm -f galaxy_telepresence

    # rm -f $DIR/tools/kubectl_config
    ## 必須將 host 的 kubectl config 載入 container 內， container kubectl 才能正確操作 host 的 k8s
    # kubectl config view --raw >$DIR/tools/kubectl_config

    case "$lang" in
    go)
        # docker run --rm --network host --name galaxy_telepresence --privileged -v $DIR/tools/kubectl_config:/root/.kube/config -v $DIR/../$repo:/go/src/code -ti gitlab-new01.vir777.com:5001/galaxy-lib/telepresence:golang telepresence --namespace galaxy-local --swap-deployment backend-$repo
        # telepresence --namespace galaxy-local --swap-deployment backend-$repo --docker-run --rm -it -v $DIR/../$repo:/go/src/code gitlab-new01.vir777.com:5001/galaxy-lib/telepresence:golang
        cd $DIR/../$repo
        docker-compose stop app
        docker run --rm --network ${repo}_default -p 80:80 --name godev -it -v $DIR/../$repo:/go/src/code godev_${repo} 
        ;;
    easyswoole)
        # docker run --rm --network host --name galaxy_telepresence --privileged -v $DIR/tools/kubectl_config:/root/.kube/config -v $DIR/../$repo:/easyswoole -ti gitlab-new01.vir777.com:5001/galaxy-lib/telepresence:easyswoole telepresence --namespace galaxy-local --swap-deployment backend-$repo
        telepresence --namespace galaxy-local --swap-deployment backend-$repo --docker-run --rm -it -v $DIR/../$repo:/easyswoole gitlab-new01.vir777.com:5001/galaxy-lib/telepresence:easyswoole bash
        ;;
    esac
}

function checkCommandTool() {
    requiredCommand=(
        galaxy
        curl
        git
        docker
        kubectl
        helm
    )

    for c in "${requiredCommand[@]}"; do
        command -v $c >/dev/null 2>&1 || {
            echo >&2 "${LRED}無法找到 $c 指令，請確認是否有正確安裝${RESTORE}"
            exit 1
        }
    done
}

function checkFortiConnect() {
    if [ "$isCheckFortiConnect" == "N" ]; then
        return
    fi

    resp=$(curl -s https://gitlab-new01.vir777.com --max-time 10)

    if [ "$resp" == "" ]; then
        echo "${LRED}無法與 Gitlab 進行連線，請確認是否已連接 Forti，或是指令加上 --no-check 跳過檢查機制${RESTORE}"
        exit 1
    fi

    git remote update >/dev/null 2>&1
    git fetch -f --prune --prune-tags >/dev/null 2>&1
}

function doUpgrade() {
    checkFortiConnect
    current=$(git describe --tags)
    latest=$(git describe --tags $(git rev-list --tags --max-count=1))

    if [ "$current" != "$latest" ]; then
        git checkout ${latest} >/dev/null 2>&1

        echo "${LMAGENTA}更新完畢，請重新執行你的指令。必要時刻需要重新啟動整個環境。${RESTORE}"
    else
        echo "${LMAGENTA}已經是最新版本。${RESTORE}"
    fi

    echo
}
