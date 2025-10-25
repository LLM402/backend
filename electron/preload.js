const { contextBridge } = require('electron');

function getDataDirPath() {
  if (process.env.ELECTRON_DATA_DIR) {
    return process.env.ELECTRON_DATA_DIR;
  }
}

contextBridge.exposeInMainWorld('electron', {
  isElectron: true,
  version: process.versions.electron,
  platform: process.platform,
  versions: process.versions,
  dataDir: getDataDirPath()
});