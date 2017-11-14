#/bin/bash

which docker
if [ $? -ne 0 ]
then
    curl -fsSL get.docker.com -o get-docker.sh && sh get-docker.sh
fi

sudo apt-get install unzip curl -y < "/dev/null"

rm -rf log-processor/ master.zip
wget https://github.com/vtrbtf/log-processing-example/archive/master.zip 
unzip master.zip -d log-processor

mv log-processor/log-processing-example-master/* .
cd input && sudo docker build -t accesslog-generator . && cd ..
cd collect && sudo docker build -t accesslog-collector . && cd ..
