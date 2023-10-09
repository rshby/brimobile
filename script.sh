# run docker jaeger
docker container create --name jaeger --publish 6831:6831/udp --publish 16686:16686 --network brimove-network jaegertracing/all-in-one:latest