@echo off
setlocal
set "KAFKA_HOME=C:\Users\user\apps\kafka"
set "KAFKA_CONFIG=%KAFKA_HOME%\config\portfolio-lens-local.properties"

if not exist "%KAFKA_HOME%\bin\windows\kafka-server-stop.bat" (
  echo Ne naiden Kafka v %KAFKA_HOME%
  exit /b 1
)

"%KAFKA_HOME%\bin\windows\kafka-server-stop.bat" "%KAFKA_CONFIG%"
