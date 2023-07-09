#!/bin/bash

SERVICES=`find . -maxdepth 2 -type d -name proto | cut -f 2 -d / | sort`

cat << EOF
package services

import (
	"micro.dev/v4/service/client"
EOF

for service in ${SERVICES[@]}; do
	echo -e "\t\"github.com/micro/services/${service}/proto\""
done

cat << EOF
)
EOF

cat << EOF

type Client struct { 
EOF

for service in ${SERVICES[@]}; do
	echo -e "\t${service^} ${service}.${service^}Service"
done

cat << EOF
}

EOF

cat << EOF
func NewClient(c client.Client) *Client {
	return &Client{
EOF

for service in ${SERVICES[@]}; do
	echo -e "\t\t${service^}: ${service}.New${service^}Service(\"${service}\", c),"
done

cat << EOF
	}
}
EOF
