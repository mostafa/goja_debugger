var generate = require('generate-source-map');
var fs = require('fs');
var process = require('process');

var file = {
  source: fs.readFileSync(process.argv[2]),
  sourceFile: process.argv[2]
};

var map = generate(file);

console.log(map.toString());

