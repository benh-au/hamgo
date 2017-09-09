# Installing hamgo

## 1. Copy files

    cp hamgo /usr/local/bin/
    cp hamgo.sample.json /etc/hamgo.json
    cp hamgo.service /etc/systemd/system/

    # Edit /etc/hamgo.json

    systemctl daemon-reload
    systemctl enable hamgo
    systemctl start hamgo
