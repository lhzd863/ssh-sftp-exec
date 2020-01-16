go run rmt-ssh-sftp.go --op="top" \
--host="122.51.161.53" \
--username="k8s" \
--userpasswd="xxx" \
--srcf="/home/k8s/Go/ws/src/github.com/lhzd863/sftp-cmd/t.txt" \
--tard="/home/k8s/Go/ws/src/github.com/lhzd863/sftp-cmd/tmp" \
--cmd="vi /home/k8s/Go/ws/src/github.com/lhzd863/ssh-sftp/t.txt"
