@echo off

:: Define the current
set "current_path=%cd%"

:: Define the full path to the executable
set "executable_path=%current_path%\independent-signer_windows_amd64.exe"

:: Define the name of the scheduled task
set "task_name=IndependentSigner"

:: Create a scheduled task to run the program every 30 minutes
schtasks /create /tn "%task_name%" /tr "\"%executable_path%\"" /sc minute /mo 30 /np /F


:: Run the program and save its output to a temporary file
echo.
"%executable_path%" > temp.txt 2>&1

:: Find the line containing "Public Key:" and extract the fixed-size value
findstr /C:"Public Key:" temp.txt > extracted.txt

:: Read the extracted line and skip the first 32 characters
set /p publickey=<extracted.txt
set "publickey=%publickey:~32%"

:: Copy the extracted Public Key value to the clipboard
echo %publickey% | clip

:: Save the Public Key value to a file
echo %publickey% > public-key

:: Clean up temporary files (optional)
del temp.txt extracted.txt