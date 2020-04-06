#!/usr/bin/env bash

source $(dirname $0)/var.sh

URL=$HOST/v1/orders
MSG_TYPE=Order
PROTO_FILE=model/pb/order.proto

user_id=43b8a2640bc945a0f6b311e3d626d942
game_id=951eae8d883ef8f4744f524c816f0242
bet=10

req=$(
cat << EOF
    user_id: "$user_id"
    game_id: "$game_id"
    bet: $bet
EOF
)

echo $req \
    | protoc --encode=$MSG_TYPE $PROTO_FILE \
    | curl \
        -H "Content-Type: application/protobuf" \
        -H "Authorization: Bearer $TOKEN" \
        --data-binary @- \
        $URL \
    | protoc --decode=$MSG_TYPE $PROTO_FILE