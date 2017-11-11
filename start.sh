#1/bin/bash
LOGPROCESSOR_NUM_OF_SERVERS=${LOGPROCESSOR_NUM_OF_SERVERS:-4}
LOGPROCESSOR_NUM_OF_LINES_FOR_EACH_FILE=${LOGPROCESSOR_NUM_OF_LINES_FOR_EACH_FILE:-1000}

sudo rm -rf "$HOME/data/" && sudo mkdir "$HOME/data/"
counter=1
collector_args=""
while [ $counter -le $LOGPROCESSOR_NUM_OF_SERVERS ]
do
    collector_args="${collector_args} --logfile /go/data/server$counter/generated_access_log.log "
    echo "Generating logfiles for server$counter "
    sudo docker run -t -v "$HOME/data:/usr/src/app/data" -i accesslog-generator -n "${LOGPROCESSOR_NUM_OF_LINES_FOR_EACH_FILE}" --prefix "./data/server$counter/generated"
    counter=$((counter+1))
done

collector_args="${collector_args//\'/}"
sudo docker run -t -v "$HOME/data:/go/data/" -v "/tmp:/tmp" -i accesslog-collector ${collector_args}