#!/bin/bash





_formatLoad(){
cat <<-_EOF_
{"name":"new-$RANDOM-$RANDOM-$RANDOM","address":"address here::$RANDOM","floors":["floor-1-$RANDOM","floor-2-$RANDOM"]}
_EOF_
}


echo $(date) start
repeat=100
for ((i=0; i<= $repeat; i++))
{
	curl -X POST    'http://127.0.0.1:8989/v1/api/building' -d "$(_formatLoad)" &
}
wait
echo $(date) done
