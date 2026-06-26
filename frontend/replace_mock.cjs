const fs = require('fs');
const path = require('path');

const srcDir = path.join(__dirname, 'src');

function processDir(dir) {
  const files = fs.readdirSync(dir);
  for (const file of files) {
    const fullPath = path.join(dir, file);
    if (fs.statSync(fullPath).isDirectory()) {
      processDir(fullPath);
    } else if (fullPath.endsWith('.vue')) {
      let content = fs.readFileSync(fullPath, 'utf8');
      
      const original = content;
      content = content.replace(/['"]\.\.\/utils\/mockData['"]/g, "'../../utils/mockData'");
      
      if (content !== original) {
        fs.writeFileSync(fullPath, content);
      }
    }
  }
}

processDir(path.join(srcDir, 'views'));
console.log('Mock replacements done');
