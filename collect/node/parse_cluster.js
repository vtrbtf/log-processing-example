var parser = require('./parse_local.js')
var fs = require('fs')

var loadp1 = parser.load('server1', process.argv[2])
var loadp2 = parser.load('server2', process.argv[3])
var loadp3 = parser.load('server3', process.argv[4])
var loadp4 = parser.load('server4', process.argv[5])

console.time("load")
Promise.all([loadp1, loadp2, loadp3, loadp4]).then((values) => {
    fs.writeFileSync('result.json', JSON.stringify(values))
    console.timeEnd("load")
})


