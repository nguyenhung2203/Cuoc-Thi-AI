const fs = require('fs');
const path = require('path');

const srcDir = path.join(__dirname, 'src');

const replacements = [
  { from: /import\s+(\w+)\s+from\s+['"]\.\.\/components\/Button\.vue['"]/g, to: 'import $1 from \'../../components/common/AppButton.vue\'' },
  { from: /import\s+(\w+)\s+from\s+['"]\.\.\/components\/Card\.vue['"]/g, to: 'import $1 from \'../../components/common/AppCard.vue\'' },
  { from: /import\s+(\w+)\s+from\s+['"]\.\.\/components\/Input\.vue['"]/g, to: 'import $1 from \'../../components/common/AppInput.vue\'' },
  { from: /import\s+(\w+)\s+from\s+['"]\.\.\/components\/Modal\.vue['"]/g, to: 'import $1 from \'../../components/common/AppModal.vue\'' },
  { from: /import\s+(\w+)\s+from\s+['"]\.\.\/components\/Table\.vue['"]/g, to: 'import $1 from \'../../components/common/AppTable.vue\'' },
  { from: /import\s+(\w+)\s+from\s+['"]\.\.\/components\/Toast\.vue['"]/g, to: 'import $1 from \'../../components/common/AppToast.vue\'' },
  { from: /import\s+(\w+)\s+from\s+['"]\.\.\/components\/Badge\.vue['"]/g, to: 'import $1 from \'../../components/common/AppBadge.vue\'' },
  
  // For layout component which is in src/components/layout/
  { from: /import\s+(\w+)\s+from\s+['"]\.\.\/components\/Button\.vue['"]/g, to: 'import $1 from \'../common/AppButton.vue\'' }, // wait, this regex is the same, so I should handle layout separately.
];

function processDir(dir, isLayoutDir = false) {
  const files = fs.readdirSync(dir);
  for (const file of files) {
    const fullPath = path.join(dir, file);
    if (fs.statSync(fullPath).isDirectory()) {
      processDir(fullPath, file === 'layout' || isLayoutDir);
    } else if (fullPath.endsWith('.vue')) {
      let content = fs.readFileSync(fullPath, 'utf8');
      
      if (isLayoutDir) {
        content = content.replace(/import\s+(\w+)\s+from\s+['"]\.\.\/components\/(\w+)\.vue['"]/g, (match, p1, p2) => {
          return `import ${p1} from '../common/App${p2}.vue'`;
        });
      } else {
        content = content.replace(/import\s+(\w+)\s+from\s+['"]\.\.\/components\/(\w+)\.vue['"]/g, (match, p1, p2) => {
          return `import ${p1} from '../../components/common/App${p2}.vue'`;
        });
      }

      fs.writeFileSync(fullPath, content);
    }
  }
}

processDir(srcDir);
console.log('Replacement done');
