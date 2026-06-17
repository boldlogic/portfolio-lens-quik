@echo off
setlocal
set "KAFKA_HOME=C:\Users\user\apps\kafka"
set "BOOTSTRAP=localhost:9092"
set "TOPIC=quik.currencies"

if not exist "%KAFKA_HOME%\bin\windows\kafka-topics.bat" (
  echo Kafka not found in %KAFKA_HOME%
  exit /b 1
)

echo Creating topic %TOPIC% (configs/quik-currency-config.yaml) ...
"%KAFKA_HOME%\bin\windows\kafka-topics.bat" --bootstrap-server %BOOTSTRAP% ^
  --create --if-not-exists ^
  --topic %TOPIC% ^
  --partitions 1 ^
  --replication-factor 1

"%KAFKA_HOME%\bin\windows\kafka-topics.bat" --bootstrap-server %BOOTSTRAP% --describe --topic %TOPIC%
