@echo off
SETLOCAL

REM Configure Git to use the hooks directory
git config core.hooksPath .githooks

ECHO Git hooks have been set up.

ENDLOCAL
pause