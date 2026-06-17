@echo off
setlocal EnableDelayedExpansion
set "KAFKA_HOME=C:\Users\user\apps\kafka"
set "KAFKA_CONFIG=%KAFKA_HOME%\config\portfolio-lens-local.properties"
set "KAFKA_DATA=%KAFKA_HOME%\data\kraft-combined-logs"
set "CLUSTER_ID_FILE=%KAFKA_HOME%\data\cluster.id"
set "KAFKA_LOG=%KAFKA_HOME%\logs\server.log"

if not exist "%KAFKA_HOME%\bin\windows\kafka-server-start.bat" (
  echo Kafka not found in %KAFKA_HOME%
  echo Download: https://archive.apache.org/dist/kafka/3.9.1/kafka_2.13-3.9.1.tgz
  exit /b 1
)

where java >nul 2>&1
if errorlevel 1 (
  echo Java not found in PATH. Install JDK 17+ and reopen terminal.
  exit /b 1
)

if not exist "%KAFKA_HOME%\data" mkdir "%KAFKA_HOME%\data"
if not exist "%KAFKA_DATA%" mkdir "%KAFKA_DATA%"
if not exist "%KAFKA_HOME%\logs" mkdir "%KAFKA_HOME%\logs"

if not exist "%KAFKA_DATA%\meta.properties" (
  "%KAFKA_HOME%\bin\windows\kafka-storage.bat" random-uuid > "%CLUSTER_ID_FILE%"
  set /p CLUSTER_ID=<"%CLUSTER_ID_FILE%"
  echo Formatting KRaft storage, cluster id: !CLUSTER_ID!
  call "%KAFKA_HOME%\bin\windows\kafka-storage.bat" format -t !CLUSTER_ID! -c "%KAFKA_CONFIG%"
  if errorlevel 1 (
    echo Storage format failed. See output above.
    exit /b 1
  )
)

netstat -ano | findstr ":9092" | findstr LISTENING >nul
if not errorlevel 1 (
  echo Kafka already listening on localhost:9092
  exit /b 0
)

echo Starting local Kafka on localhost:9092 ...
cd /d "%KAFKA_HOME%"
start "portfolio-lens-kafka" /MIN cmd /c ""bin\windows\kafka-server-start.bat" "config\portfolio-lens-local.properties" >> "logs\server.log" 2>&1"

set /a WAIT_SEC=0
:wait_loop
netstat -ano | findstr ":9092" | findstr LISTENING >nul
if not errorlevel 1 (
  echo Kafka started. Bootstrap: localhost:9092
  echo Log: %KAFKA_LOG%
  exit /b 0
)
set /a WAIT_SEC+=2
if !WAIT_SEC! GEQ 60 (
  echo Kafka did not open port 9092 within 60 seconds.
  echo Check log: %KAFKA_LOG%
  exit /b 1
)
timeout /t 2 /nobreak >nul 2>&1
if errorlevel 1 ping -n 3 127.0.0.1 >nul
goto wait_loop
