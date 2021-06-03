#!/bin/bash

curl -fsSL "https://github.com/$1.keys" >> /home/$2/.ssh/authorized_keys
chmod 700 /home/$2/.ssh
chmod 600 /home/$2/.ssh/authorized_keys