param(
  [string]$DriverName = "",
  [string]$Server = "",
  [string]$DbName = "",
  [string]$Dsn64 = "",
  [string]$SqlUser = "",
  [string]$SqlPassword = "",
  [string]$EnvFile = "",
  [switch]$Force,
  [switch]$PromptPassword,
  [switch]$UseTrustedConnection
)

$ErrorActionPreference = "Stop"

$RepoRoot = Split-Path $PSScriptRoot -Parent
if ([string]::IsNullOrEmpty($EnvFile)) {
  $EnvFile = Join-Path $RepoRoot ".env"
}

function Test-IsAdmin {
  $id = [Security.Principal.WindowsIdentity]::GetCurrent()
  $p = New-Object Security.Principal.WindowsPrincipal($id)
  return $p.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

function Normalize-EnvValue {
  param([string]$Value)
  $v = $Value.Trim().Trim("`r")
  if (($v.Length -ge 2) -and (($v.StartsWith('"') -and $v.EndsWith('"')) -or ($v.StartsWith("'") -and $v.EndsWith("'")))) {
    return $v.Substring(1, $v.Length - 2)
  }
  return $v
}

function Read-DotEnv {
  param([string]$Path)

  $result = @{}
  if (-not (Test-Path -LiteralPath $Path)) {
    return $result
  }

  foreach ($line in Get-Content -LiteralPath $Path -Encoding UTF8) {
    $t = $line.Trim()
    if ($t -match '^\s*#' -or $t -eq "") { continue }
    if ($t -notmatch '^\s*([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(.*)$') { continue }
    $result[$Matches[1]] = (Normalize-EnvValue -Value $Matches[2])
  }

  return $result
}

function Get-DotEnvValue {
  param(
    [hashtable]$DotEnv,
    [string[]]$Keys
  )
  foreach ($key in $Keys) {
    if ($DotEnv.ContainsKey($key) -and -not [string]::IsNullOrEmpty($DotEnv[$key])) {
      return $DotEnv[$key]
    }
  }
  return $null
}

function Get-DriverDllPath {
  param([string]$DriverRegKey, [string]$DriverDisplayName)

  if (-not (Test-Path $DriverRegKey)) {
    throw "ODBC driver not found in registry: '$DriverRegKey'. Install x64 '$DriverDisplayName'."
  }
  $p = Get-ItemProperty -Path $DriverRegKey
  if (-not $p.Driver) {
    throw "ODBC driver registry key exists, but 'Driver' value is empty: '$DriverRegKey'."
  }
  return $p.Driver
}

function Resolve-DriverName {
  if ($DriverName) {
    return $DriverName
  }

  $candidates = @(
    "ODBC Driver 18 for SQL Server",
    "ODBC Driver 17 for SQL Server"
  )

  foreach ($candidate in $candidates) {
    $key = "HKLM:\Software\ODBC\ODBCINST.INI\$candidate"
    if (Test-Path $key) {
      return $candidate
    }
  }

  throw "ODBC Driver 18/17 for SQL Server (x64) not found. Install driver or pass -DriverName."
}

function Add-AuthToDsn {
  param(
    [string]$DsnKeyPath,
    [bool]$TrustedConnection,
    [string]$User,
    [string]$PlainPassword
  )
  if (-not (Test-Path $DsnKeyPath)) {
    return
  }

  if ($TrustedConnection) {
    New-ItemProperty -Path $DsnKeyPath -Name "Trusted_Connection" -Value "Yes" -PropertyType String -Force | Out-Null
    Remove-ItemProperty -Path $DsnKeyPath -Name "UID" -ErrorAction SilentlyContinue
    Remove-ItemProperty -Path $DsnKeyPath -Name "PWD" -ErrorAction SilentlyContinue
    return
  }

  New-ItemProperty -Path $DsnKeyPath -Name "Trusted_Connection" -Value "No" -PropertyType String -Force | Out-Null
  New-ItemProperty -Path $DsnKeyPath -Name "UID" -Value $User -PropertyType String -Force | Out-Null
  New-ItemProperty -Path $DsnKeyPath -Name "PWD" -Value $PlainPassword -PropertyType String -Force | Out-Null
}

function New-SystemDsn {
  param(
    [string]$BaseKey,
    [string]$DriverRegKey,
    [string]$DsnName
  )

  $dsnKey = "$BaseKey\$DsnName"
  $dsnsListKey = "$BaseKey\ODBC Data Sources"

  if ((Test-Path $dsnKey) -and (-not $Force)) {
    throw "DSN '$DsnName' already exists at '$dsnKey'. Re-run with -Force to overwrite."
  }

  $driverDll = Get-DriverDllPath -DriverRegKey $DriverRegKey -DriverDisplayName $DriverName

  New-Item -Path $dsnKey -Force | Out-Null
  New-ItemProperty -Path $dsnKey -Name "Driver" -Value $driverDll -PropertyType String -Force | Out-Null
  New-ItemProperty -Path $dsnKey -Name "Server" -Value $Server -PropertyType String -Force | Out-Null
  New-ItemProperty -Path $dsnKey -Name "Database" -Value $DbName -PropertyType String -Force | Out-Null
  New-ItemProperty -Path $dsnKey -Name "Encrypt" -Value "No" -PropertyType String -Force | Out-Null

  New-Item -Path $dsnsListKey -Force | Out-Null
  New-ItemProperty -Path $dsnsListKey -Name $DsnName -Value $DriverName -PropertyType String -Force | Out-Null

  return @{
    DsnName   = $DsnName
    BaseKey   = $BaseKey
    DriverDll = $driverDll
    Server    = $Server
    Database  = $DbName
    DsnKey    = $dsnKey
  }
}

if (-not (Test-IsAdmin)) {
  throw "Run this script as Administrator (System DSN writes to HKLM)."
}

$dotEnv = Read-DotEnv -Path $EnvFile
$envLoaded = $dotEnv.Count -gt 0

if (-not $PSBoundParameters.ContainsKey('Server')) {
  $fromEnv = Get-DotEnvValue -DotEnv $dotEnv -Keys @('QUIK_ODBC_SERVER', 'MSSQL_SERVER', 'DB_SERVER')
  if ($fromEnv) { $Server = $fromEnv }
}
if ([string]::IsNullOrEmpty($Server)) {
  $Server = "localhost,1433"
}

if (-not $PSBoundParameters.ContainsKey('DbName')) {
  $fromEnv = Get-DotEnvValue -DotEnv $dotEnv -Keys @('QUIK_ODBC_DB', 'QUIK_ODBC_DATABASE', 'MSSQL_DB', 'DB_NAME')
  if ($fromEnv) { $DbName = $fromEnv }
}
if ([string]::IsNullOrEmpty($DbName)) {
  $DbName = "portfolio_lens_quik"
}

if (-not $PSBoundParameters.ContainsKey('Dsn64')) {
  $fromEnv = Get-DotEnvValue -DotEnv $dotEnv -Keys @('QUIK_ODBC_DSN')
  if ($fromEnv) { $Dsn64 = $fromEnv }
}
if ([string]::IsNullOrEmpty($Dsn64)) {
  $Dsn64 = "QuikPortfolioLocal_64"
}

if (-not $PSBoundParameters.ContainsKey('SqlUser')) {
  $fromEnv = Get-DotEnvValue -DotEnv $dotEnv -Keys @('QUIK_ODBC_USER', 'MSSQL_USER')
  if ($fromEnv) { $SqlUser = $fromEnv }
}
if ([string]::IsNullOrEmpty($SqlUser)) {
  $SqlUser = "quik_odbc_writer"
}

$resolvedDriverName = Resolve-DriverName
$DriverName = $resolvedDriverName

if ($UseTrustedConnection -and ($PromptPassword -or -not [string]::IsNullOrEmpty($SqlPassword))) {
  throw "UseTrustedConnection cannot be combined with SQL login/password options."
}

$sqlPwd = $SqlPassword
if ([string]::IsNullOrEmpty($sqlPwd)) {
  $fromEnv = Get-DotEnvValue -DotEnv $dotEnv -Keys @('QUIK_ODBC_PASSWORD', 'MSSQL_PASSWORD')
  if ($fromEnv) { $sqlPwd = $fromEnv }
}
if ([string]::IsNullOrEmpty($sqlPwd) -and $env:QUIK_ODBC_PASSWORD) {
  $sqlPwd = (Normalize-EnvValue -Value $env:QUIK_ODBC_PASSWORD)
}
if ($PromptPassword -and [string]::IsNullOrEmpty($sqlPwd)) {
  $secure = Read-Host "Enter SQL password for user $SqlUser" -AsSecureString
  $bstr = [Runtime.InteropServices.Marshal]::SecureStringToBSTR($secure)
  try {
    $sqlPwd = [Runtime.InteropServices.Marshal]::PtrToStringBSTR($bstr)
  } finally {
    [Runtime.InteropServices.Marshal]::ZeroFreeBSTR($bstr)
  }
}

if (-not $UseTrustedConnection -and [string]::IsNullOrEmpty($sqlPwd)) {
  throw "Укажите QUIK_ODBC_PASSWORD в $EnvFile или -SqlPassword."
}

Write-Host "Creating 64-bit System DSN for QUIK ODBC export..."
if ($envLoaded) {
  Write-Host "  Env:     $EnvFile"
} else {
  Write-Host "  Env:     not found ($EnvFile), defaults and -Server/-SqlPassword"
}
Write-Host "  Driver:  $DriverName"
Write-Host "  Server:  $Server"
Write-Host "  DB:      $DbName"
Write-Host "  DSN:     $Dsn64"
Write-Host "  Encrypt: No (local server)"
if ($UseTrustedConnection) {
  Write-Host "  Auth:    Trusted_Connection=Yes"
} else {
  Write-Host "  Auth:    SQL login ($SqlUser), write only to schema quik"
}

$res64 = New-SystemDsn `
  -BaseKey "HKLM:\Software\ODBC\ODBC.INI" `
  -DriverRegKey "HKLM:\Software\ODBC\ODBCINST.INI\$resolvedDriverName" `
  -DsnName $Dsn64

Add-AuthToDsn -DsnKeyPath $res64.DsnKey -TrustedConnection $UseTrustedConnection.IsPresent -User $SqlUser -PlainPassword $sqlPwd

Write-Host ""
Write-Host "Done."
Write-Host ("  DSN: {0} -> {1} / {2} (Driver DLL: {3})" -f $res64.DsnName, $res64.Server, $res64.Database, $res64.DriverDll)
Write-Host ""
Write-Host "Prerequisite: scripts/sql/bootstrap/create_quik_odbc_user.sql and migration 008 (grants on quik.*)."
Write-Host "Check: C:\Windows\System32\odbcad32.exe -> System DSN -> $Dsn64 -> Test"
Write-Host "Requires 64-bit QUIK (or other 64-bit ODBC client). 32-bit DSN is not created by this script."
