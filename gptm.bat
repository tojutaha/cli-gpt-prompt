@echo off

if "%2" == "" (
    gptm-executable.exe "%~1" | glow
    exit /b
)

gptm-executable.exe "%~1" "%~2" | glow
