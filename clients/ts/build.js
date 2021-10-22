const chalk = require('chalk');
const path = require('path');
const fs = require('fs');
const rimraf = require('rimraf');
const { ncp } = require('ncp');

function getTmpEsmDirectories() {
  return fs
    .readdirSync('./tmp/esm')
    .filter(file => fs.statSync(`./tmp/esm/${file}`).isDirectory());
}

function log(text) {
  console.log(`${chalk.cyan('M3O JS:')} ${text}`);
}

function writeModulePackageJsonFile(location) {
  fs.writeFileSync(
    `${location}/package.json`,
    `{"module": "./esm/index.js"}`,
    'utf8'
  );
}

function deleteDirectory(directory) {
  return new Promise(resolve => {
    rimraf(directory, err => {
      resolve();
    });
  });
}

function copyAllTmpFolders() {
  return new Promise((resolve, reject) => {
    // Now copy to root level
    ncp(path.join(__dirname, 'tmp'), __dirname, err => {
      if (err) {
        reject(err);
      } else {
        resolve();
      }
    });
  });
}

function moveToLocalEsmFolders() {
  return new Promise((resolve, reject) => {
    const esmDirs = getTmpEsmDirectories();

    // Move the files around in tmp...
    esmDirs.forEach(dir => {
      const currentPath = path.join(__dirname, 'tmp/esm', dir);

      fs.readdirSync(currentPath).forEach(async file => {
        const currentFilePath = path.join(currentPath, file);
        const newFilePath = path.join(__dirname, 'tmp', dir, 'esm', file);
        const esmFolderLocation = path.join(__dirname, 'tmp', dir, 'esm');

        try {
          if (!fs.existsSync(esmFolderLocation)) {
            fs.mkdirSync(esmFolderLocation);
          }

          fs.renameSync(currentFilePath, newFilePath);
          writeModulePackageJsonFile(`./tmp/${dir}`);
          await deleteDirectory(`./tmp/esm/${dir}`);
        } catch (err) {
          reject(err);
        }
      });
    });

    log('Moved local esm folders');
    resolve();
  });
}

async function build() {
  log('Moving to correct folders');

  try {
    await moveToLocalEsmFolders();
    await copyAllTmpFolders();
    writeModulePackageJsonFile('./tmp/esm');
    await deleteDirectory('./tmp');
  } catch (e) {
    console.log(e);
  }
}

build();