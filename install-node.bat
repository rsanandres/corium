@echo off
echo Installing Node.js and npm...

:: Download Node.js installer
curl -o node-installer.msi https://nodejs.org/dist/v20.11.0/node-v20.11.0-x64.msi

:: Install Node.js
msiexec /i node-installer.msi /qn

:: Refresh environment variables
setx PATH "%PATH%;C:\Program Files\nodejs" /M

:: Verify installation
echo.
echo Verifying installation...
node -v
npm -v

echo.
echo Installation complete! You can now use Node.js and npm.
pause 