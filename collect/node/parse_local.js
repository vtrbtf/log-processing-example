const readline = require('readline');
const fs = require('fs');
const sh = require('shorthash')
const moment = require('moment')

exports.load = (id, file) => {
    return new Promise((resolve) => {
        const rl = readline.createInterface({
          input: fs.createReadStream(file)
        });
        var lineCounter = 0;
        var lineIndex = [];
        var useridIndex = {};
        rl.on('line', (line) => {
            var lineIndex = ++lineCounter
            let userid = line.match(/userid\=([\w-]+)/)[1]
            let timestamp = moment(line.match(/- - \[([\w:\s\d-\/]+)\]/)[1], "DD/MMM/YYYY:HH:mm:ss Z").unix()
            //lineIndex.push(sh.unique(userid))
            if (!useridIndex[userid]) {
                useridIndex[userid] = [{t: timestamp, i: lineIndex}]
            } else {
                useridIndex[userid].push([{t: timestamp, i: lineIndex}])
            }
        }).on('close', () => resolve({serverId: id, filePath: file, useridIndex: useridIndex, lines: lineIndex}))
    })
}
