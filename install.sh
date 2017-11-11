#/bin/bash

curl -fsSL get.docker.com -o get-docker.sh && sh get-docker.sh
sudo apt-get install -y unzip curl

rm -rf log-processor/ master.zip && wget https://github.com/vtrbtf/log-processing-example/archive/master.zip && unzip master.zip -d log-processor
cd log-processor/log-processing-example-master
cd input && sudo docker build -t accesslog-generator . && cd ..
cd collect && sudo docker build -t accesslog-collector . && cd ..
