const { app, BrowserWindow, dialog, Tray, Menu, shell } = require('electron');
const { spawn } = require('child_process');
const path = require('path');
const http = require('http');
const fs = require('fs');

let mainWindow;
let serverProcess;
let tray = null;
let serverErrorLogs = [];
const PORT = 3000;
const DEV_FRONTEND_PORT = 5173; // Vite dev server port

function saveAndOpenErrorLog() {
  try {
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
    const logFileName = `new-api-crash-${timestamp}.log`;
    const logDir = app.getPath('logs');
    const logFilePath = path.join(logDir, logFileName);
    
    if (!fs.existsSync(logDir)) {
      fs.mkdirSync(logDir, { recursive: true });
    }
    
    const logContent = `New API Crash Log  
Generation Time: ${new Date().toLocaleString('zh-CN')}  
Platform: ${process.platform}  
Architecture: ${process.arch}  
Application Version: ${app.getVersion()}  

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  

Complete Error Log:  

${serverErrorLogs.join('\\n')}  

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  

Log File Location: ${logFilePath}`;
    
    fs.writeFileSync(logFilePath, logContent, 'utf8');
    
    
    shell.openPath(logFilePath).then((error) => {
      if (error) {
        console.error('Failed to open log file:', error);
        shell.showItemInFolder(logFilePath);
      }
    });
    
    return logFilePath;
  } catch (err) {
    console.error('Failed to save error log:', err);
    return null;
  }
}

function analyzeError(errorLogs) {
  const allLogs = errorLogs.join('\n');
  
  if (allLogs.includes('failed to start HTTP server') || 
      allLogs.includes('bind: address already in use') ||
      allLogs.includes('listen tcp') && allLogs.includes('bind: address already in use')) {
    return {
      type: 'Port is occupied',
      title: 'port' + PORT + 'Occupied',
      message: 'Unable to start the server, the port is occupied by another program.',
      solution: `Possible solutions:\\n\\n1. Close other programs occupying port ${PORT}\\n2. Check if another New API instance is already running\\n3. Use the following command to find the process occupying the port:\\n   Mac/Linux: lsof -i :${PORT}\\n   Windows: netstat -ano | findstr :${PORT}\\n4. Restart the computer to free the port`
    };
  }
  
  if (allLogs.includes('database is locked') || 
      allLogs.includes('unable to open database')) {
    return {
      type: 'The data file is in use.',
      title: 'Unable to access the data file',
      message: 'The application data file is being used by another program.',
      solution: 'Possible solutions:\\n\\n1. Check if another New API window is already open\\n   - Look for other New API icons in the taskbar/Dock\\n   - Check for New API icons in the system tray (Windows) or menu bar (Mac)\\n\\n2. If you just closed the application, please wait 10 seconds before trying again\\n\\n3. Restart your computer to free up occupied files\\n\\n4. If the problem persists, you can try:\\n   - Exit all New API instances\\n   - Delete temporary files in the data directory (.db-shm and .db-wal)\\n   - Restart the application'
    };
  }
  
  if (allLogs.includes('permission denied') || 
      allLogs.includes('access denied')) {
    return {
      type: 'Permission error',
      title: 'Insufficient permissions',
      message: 'The program does not have sufficient permissions to perform the operation.',
      solution: 'Possible solutions:\\n\\n1. Run the program with administrator/root privileges\\n2. Check the read and write permissions of the data directory\\n3. Check the permissions of the executable file\\n4. On Mac, check the security and privacy settings'
    };
  }
  
  if (allLogs.includes('network is unreachable') || 
      allLogs.includes('no such host') ||
      allLogs.includes('connection refused')) {
    return {
      type: 'Network error',
      title: 'Network connection failed',
      message: 'Unable to establish a network connection',
      solution: 'Possible solutions:\\n\\n1. Check if the network connection is normal\\n2. Check the firewall settings\\n3. Check the proxy configuration\\n4. Confirm that the target server address is correct'
    };
  }
  
  if (allLogs.includes('invalid configuration') || 
      allLogs.includes('failed to parse config') ||
      allLogs.includes('yaml') || allLogs.includes('json') && allLogs.includes('parse')) {
    return {
      type: 'Configuration error',
      title: 'Configuration file error',
      message: 'The configuration file format is incorrect or contains invalid configurations.',
      solution: 'Possible solutions:\\n\\n1. Check if the configuration file format is correct\\n2. Restore default configuration\\n3. Delete the configuration file to let the program regenerate\\n4. Refer to the documentation for the correct configuration format'
    };
  }
  
  if (allLogs.includes('out of memory') || 
      allLogs.includes('cannot allocate memory')) {
    return {
      type: 'Insufficient memory',
      title: 'Insufficient system memory',
      message: 'Insufficient memory during program execution',
      solution: 'Possible solutions:\\n\\n1. Close other memory-consuming programs\\n2. Increase available system memory\\n3. Restart the computer to free up memory\\n4. Check for memory leaks'
    };
  }
  
  if (allLogs.includes('no such file or directory') || 
      allLogs.includes('cannot find the file')) {
    return {
      type: 'File missing',
      title: 'Cannot find the required file',
      message: 'Missing files required for the program to run.',
      solution: 'Possible solutions:\\n\\n1. Reinstall the application\\n2. Check if the installation directory is complete\\n3. Ensure all dependency files are present\\n4. Check if the file path is correct'
    };
  }
  
  return null;
}

function getBinaryPath() {
  const isDev = process.env.NODE_ENV === 'development';
  const platform = process.platform;

  if (isDev) {
    const binaryName = platform === 'win32' ? 'new-api.exe' : 'new-api';
    return path.join(__dirname, '..', binaryName);
  }

  let binaryName;
  switch (platform) {
    case 'win32':
      binaryName = 'new-api.exe';
      break;
    case 'darwin':
      binaryName = 'new-api';
      break;
    case 'linux':
      binaryName = 'new-api';
      break;
    default:
      binaryName = 'new-api';
  }

  return path.join(process.resourcesPath, 'bin', binaryName);
}

// Check if a server is available with retry logic
function checkServerAvailability(port, maxRetries = 30, retryDelay = 1000) {
  return new Promise((resolve, reject) => {
    let currentAttempt = 0;
    
    const tryConnect = () => {
      currentAttempt++;
      
      if (currentAttempt % 5 === 1 && currentAttempt > 1) {
        console.log(`Attempting to connect to port ${port}... (attempt ${currentAttempt}/${maxRetries})`);
      }
      
      const req = http.get({
        hostname: '127.0.0.1', // Use IPv4 explicitly instead of 'localhost' to avoid IPv6 issues
        port: port,
        timeout: 10000
      }, (res) => {
        // Server responded, connection successful
        req.destroy();
        console.log(`✓ Successfully connected to port ${port} (status: ${res.statusCode})`);
        resolve();
      });

      req.on('error', (err) => {
        if (currentAttempt >= maxRetries) {
          reject(new Error(`Failed to connect to port ${port} after ${maxRetries} attempts: ${err.message}`));
        } else {
          setTimeout(tryConnect, retryDelay);
        }
      });

      req.on('timeout', () => {
        req.destroy();
        if (currentAttempt >= maxRetries) {
          reject(new Error(`Connection timeout on port ${port} after ${maxRetries} attempts`));
        } else {
          setTimeout(tryConnect, retryDelay);
        }
      });
    };
    
    tryConnect();
  });
}

function startServer() {
  return new Promise((resolve, reject) => {
    const isDev = process.env.NODE_ENV === 'development';

    const userDataPath = app.getPath('userData');
    const dataDir = path.join(userDataPath, 'data');
    
    process.env.ELECTRON_DATA_DIR = dataDir;
    
    if (isDev) {
      console.log('Development mode: skipping server startup');
      console.log('Please make sure you have started:');
      console.log('  1. Go backend: go run main.go (port 3000)');
      console.log('  2. Frontend dev server: cd web && bun dev (port 5173)');
      console.log('');
      console.log('Checking if servers are running...');
      
      // First check if both servers are accessible
      checkServerAvailability(DEV_FRONTEND_PORT)
        .then(() => {
          console.log('✓ Frontend dev server is accessible on port 5173');
          resolve();
        })
        .catch((err) => {
          console.error(`✗ Cannot connect to frontend dev server on port ${DEV_FRONTEND_PORT}`);
          console.error('Please make sure the frontend dev server is running:');
          console.error('  cd web && bun dev');
          reject(err);
        });
      return;
    }
    const env = { ...process.env, PORT: PORT.toString() };

    if (!fs.existsSync(dataDir)) {
      fs.mkdirSync(dataDir, { recursive: true });
    }

    env.SQLITE_PATH = path.join(dataDir, 'new-api.db');
    
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
    console.log('📁 Your data storage location:');
    console.log('   ' + dataDir);
    console.log('💡 Backup Tip: Copy this directory to back up all data.');
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');

    const binaryPath = getBinaryPath();
    const workingDir = process.resourcesPath;
    
    console.log('Starting server from:', binaryPath);

    serverProcess = spawn(binaryPath, [], {
      env,
      cwd: workingDir
    });

    serverProcess.stdout.on('data', (data) => {
      console.log(`Server: ${data}`);
    });

    serverProcess.stderr.on('data', (data) => {
      const errorMsg = data.toString();
      console.error(`Server Error: ${errorMsg}`);
      serverErrorLogs.push(errorMsg);
      if (serverErrorLogs.length > 100) {
        serverErrorLogs.shift();
      }
    });

    serverProcess.on('error', (err) => {
      console.error('Failed to start server:', err);
      reject(err);
    });

    serverProcess.on('close', (code) => {
      console.log(`Server process exited with code ${code}`);
      if (code !== 0 && code !== null) {
        const errorDetails = serverErrorLogs.length > 0 
          ? serverErrorLogs.slice(-20).join('\n') 
          : 'No error log captured';
        
        const knownError = analyzeError(serverErrorLogs);
        
        let dialogOptions;
        if (knownError) {
          dialogOptions = {
            type: 'error',
            title: knownError.title,
            message: knownError.message,
            detail: `${knownError.solution}\\n\\n━━━━━━━━━━━━━━━━━━━━━━\\n\\nExit code: ${code}\\n\\nError type: ${knownError.type}\\n\\nRecent error log:\\n${errorDetails}`,
            buttons: ['Exit the application', 'View full log'],
            defaultId: 0,
            cancelId: 0
          };
        } else {
          dialogOptions = {
            type: 'error',
            title: 'Server crash',
            message: 'The server process exited abnormally.',
            detail: `Exit code: ${code}\\n\\nRecent error message:\\n${errorDetails}`,
            buttons: ['Exit the application', 'View full log'],
            defaultId: 0,
            cancelId: 0
          };
        }
        
        dialog.showMessageBox(dialogOptions).then((result) => {
          if (result.response === 1) {
            const logPath = saveAndOpenErrorLog();
          
            const confirmMessage = logPath 
              ? `The log has been saved to:\\n${logPath}\\n\\nThe log file has been opened in the default text editor.\\n\\nClick "Exit" to close the application.`
              : 'Log saving failed, but output has been sent to the console.\\n\\nClick "Exit" to close the application.';
            
            dialog.showMessageBox({
              type: 'info',
              title: 'Log has been saved.',
              message: confirmMessage,
              buttons: ['Exit'],
              defaultId: 0
            }).then(() => {
              app.isQuitting = true;
              app.quit();
            });
            
            
            console.log('=== Complete Error Log ===');
            console.log(serverErrorLogs.join('\n'));
          } else {

            app.isQuitting = true;
            app.quit();
          }
        });
      } else {
      
        if (mainWindow && !mainWindow.isDestroyed()) {
          mainWindow.close();
        }
      }
    });

    checkServerAvailability(PORT)
      .then(() => {
        console.log('✓ Backend server is accessible on port 3000');
        resolve();
      })
      .catch((err) => {
        console.error('✗ Failed to connect to backend server');
        reject(err);
      });
  });
}

function createWindow() {
  const isDev = process.env.NODE_ENV === 'development';
  const loadPort = isDev ? DEV_FRONTEND_PORT : PORT;
  
  mainWindow = new BrowserWindow({
    width: 1080,
    height: 720,
    webPreferences: {
      preload: path.join(__dirname, 'preload.js'),
      nodeIntegration: false,
      contextIsolation: true
    },
    title: 'New API',
    icon: path.join(__dirname, 'icon.png')
  });

  mainWindow.loadURL(`http://127.0.0.1:${loadPort}`);
  
  console.log(`Loading from: http://127.0.0.1:${loadPort}`);

  if (isDev) {
    mainWindow.webContents.openDevTools();
  }

  // Close to tray instead of quitting
  mainWindow.on('close', (event) => {
    if (!app.isQuitting) {
      event.preventDefault();
      mainWindow.hide();
      if (process.platform === 'darwin') {
        app.dock.hide();
      }
    }
  });

  mainWindow.on('closed', () => {
    mainWindow = null;
  });
}

function createTray() {
  // Use template icon for macOS (black with transparency, auto-adapts to theme)
  // Use colored icon for Windows
  const trayIconPath = process.platform === 'darwin'
    ? path.join(__dirname, 'tray-iconTemplate.png')
    : path.join(__dirname, 'tray-icon-windows.png');

  tray = new Tray(trayIconPath);

  const contextMenu = Menu.buildFromTemplate([
    {
      label: 'Show New API',
      click: () => {
        if (mainWindow === null) {
          createWindow();
        } else {
          mainWindow.show();
          if (process.platform === 'darwin') {
            app.dock.show();
          }
        }
      }
    },
    { type: 'separator' },
    {
      label: 'Quit',
      click: () => {
        app.isQuitting = true;
        app.quit();
      }
    }
  ]);

  tray.setToolTip('New API');
  tray.setContextMenu(contextMenu);

  // On macOS, clicking the tray icon shows the window
  tray.on('click', () => {
    if (mainWindow === null) {
      createWindow();
    } else {
      mainWindow.isVisible() ? mainWindow.hide() : mainWindow.show();
      if (mainWindow.isVisible() && process.platform === 'darwin') {
        app.dock.show();
      }
    }
  });
}

app.whenReady().then(async () => {
  try {
    await startServer();
    createTray();
    createWindow();
  } catch (err) {
    console.error('Failed to start application:', err);
    
    
    const knownError = analyzeError(serverErrorLogs);
    
    if (knownError) {
      dialog.showMessageBox({
        type: 'error',
        title: knownError.title,
        message: `Startup failed: ${knownError.message}`,
        detail: `${knownError.solution}\\n\\n━━━━━━━━━━━━━━━━━━━━━━\\n\\nError message: ${err.message}\\n\\nError type: ${knownError.type}`,
        buttons: ['Exit', 'View full log'],
        defaultId: 0,
        cancelId: 0
      }).then((result) => {
        if (result.response === 1) {
          
          const logPath = saveAndOpenErrorLog();
          
          const confirmMessage = logPath 
            ? `The log has been saved to:\\n${logPath}\\n\\nThe log file has been opened in the default text editor.\\n\\nClick "Exit" to close the application.`
            : 'Log saving failed, but output has been sent to the console.\\n\\nClick "Exit" to close the application.';
          
          dialog.showMessageBox({
            type: 'info',
            title: 'Log has been saved.',
            message: confirmMessage,
            buttons: ['Exit'],
            defaultId: 0
          }).then(() => {
            app.quit();
          });
          
          console.log('=== Complete Error Log ===');
          console.log(serverErrorLogs.join('\n'));
        } else {
          app.quit();
        }
      });
    } else {
      dialog.showMessageBox({
        type: 'error',
        title: 'Startup failed',
        message: 'Unable to start the server',
        detail: `Error message: ${err.message}\\n\\nPlease check the logs for more information.`,
        buttons: ['Exit', 'View full log'],
        defaultId: 0,
        cancelId: 0
      }).then((result) => {
        if (result.response === 1) {
          
          const logPath = saveAndOpenErrorLog();
          
          const confirmMessage = logPath 
            ? `The log has been saved to:\\n${logPath}\\n\\nThe log file has been opened in the default text editor.\\n\\nClick "Exit" to close the application.`
            : 'Log saving failed, but output has been sent to the console.\\n\\nClick "Exit" to close the application.';
          
          dialog.showMessageBox({
            type: 'info',
            title: 'Log has been saved.',
            message: confirmMessage,
            buttons: ['Exit'],
            defaultId: 0
          }).then(() => {
            app.quit();
          });
          
          console.log('=== Complete Error Log ===');
          console.log(serverErrorLogs.join('\n'));
        } else {
          app.quit();
        }
      });
    }
  }
});

app.on('window-all-closed', () => {
  // Don't quit when window is closed, keep running in tray
  // Only quit when explicitly choosing Quit from tray menu
});

app.on('activate', () => {
  if (BrowserWindow.getAllWindows().length === 0) {
    createWindow();
  }
});

app.on('before-quit', (event) => {
  if (serverProcess) {
    event.preventDefault();

    console.log('Shutting down server...');
    serverProcess.kill('SIGTERM');

    setTimeout(() => {
      if (serverProcess) {
        serverProcess.kill('SIGKILL');
      }
      app.exit();
    }, 5000);

    serverProcess.on('close', () => {
      serverProcess = null;
      app.exit();
    });
  }
});