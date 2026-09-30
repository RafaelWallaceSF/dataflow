const fs = require('fs');
const filePath = '/opt/vikunja-luah-src/frontend/src/services/projectFile.ts';
let content = fs.readFileSync(filePath, 'utf8');
content = content.replace(
  'transformRequest: formData => formData,',
  `headers: {
\t\t\t\t\t'Content-Type': 'multipart/form-data',
\t\t\t\t},
\t\t\t\ttransformRequest: formData => formData,`
);
fs.writeFileSync(filePath, content);
console.log('Fixed projectFile.ts');
