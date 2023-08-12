@echo off

if "%2" == "" (
    gpt-executable.exe "%~1" | glow
    exit /b
)

gpt-executable.exe "%~1" "%~2" | glow
