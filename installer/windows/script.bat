@echo off

:: Define the current
set "current_path=%cd%"

:: Define the full path to the executable
set "executable_path=%current_path%\independent-signer_windows_amd64.exe"

:: Define the name of the scheduled task
set "task_name=IndependentSigner"

:: Create a scheduled task to run the program every 30 minutes
schtasks /create /tn "%task_name%" /tr "%executable_path%" /sc minute /mo 30 /np /F


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

:: Display a message indicating that the Public Key has been copied to the clipboard
echo The Public Key has been copied to the clipboard: 
echo %publickey%

:: Display instructions for the user
echo.
echo Navigate to https://coiin.io/console/verificationnodes on the Network Validation Layer Nodes page of the Coiin Console.
echo Paste the Public Key printed in the terminal window into the "Enter Public Key" text box and click the "Register Node" button.
echo.

:: Clean up temporary files (optional)
del temp.txt extracted.txt

:: Pause the script to allow the user to read the output
pause