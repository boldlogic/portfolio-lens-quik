param(
  [string]$DriverName = "",
  [string]$Server = "localhost,1433",
  [string]$DbName = "portfolio_lens_quik",
  [string]$Dsn64 = "QuikPortfolioLocal_64",
  [string]$SqlUser = "quik_portfolio_app",
  [string]$SqlPassword = "",
  [string]$EnvFile = "",
  [switch]$Force,
  [switch]$UseCredentialsFromEnv,
  [switch]$PromptPassword,
  [switch]$UseTrustedConnection
)

$ErrorActionPreference = "Stop"

function Test-IsAdmin {
  $id = [Security.Principal.WindowsIdentity]::GetCurrent()
  $p = New-Object Security.Principal.WindowsPrincipal($id)
  return $p.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

function Get-MssqlPasswordFromEnvFile {
  param([string]$Path)
  if (-not (Test-Path -LiteralPath $Path)) {
    throw "Env file not found: $Path"
  }
  foreach ($line in Get-Content -LiteralPath $Path) {
    $t = $line.Trim()
    if ($t -match '^\s*#' -or $t -eq "") { continue }
    if ($t -match '^\s*MSSQL_SA_PASSWORD\s*=\s*(.*)$') {
      return $Matches[1].Trim()
    }
  }
  throw "MSSQL_SA_PASSWORD not found in $Path"
}

function Get-SqlPassword {
  if ($env:MSSQL_SA_PASSWORD) {
    return $env:MSSQL_SA_PASSWORD
  }
  if ($EnvFile) {
    return (Get-MssqlPasswordFromEnvFile -Path $EnvFile)
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

$resolvedDriverName = Resolve-DriverName
$DriverName = $resolvedDriverName

if ($UseTrustedConnection -and ($UseCredentialsFromEnv -or $PromptPassword -or -not [string]::IsNullOrEmpty($SqlPassword))) {
  throw "UseTrustedConnection cannot be combined with SQL login/password options."
}

$sqlPwd = $SqlPassword
if ($UseCredentialsFromEnv -and [string]::IsNullOrEmpty($sqlPwd)) {
  $sqlPwd = Get-SqlPassword
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
  throw "For SQL auth provide -SqlPassword, or -PromptPassword, or -UseCredentialsFromEnv."
}

Write-Host "Creating 64-bit System DSN for QUIK ODBC export..."
Write-Host "  Driver:  $DriverName"
Write-Host "  Server:  $Server"
Write-Host "  DB:      $DbName"
Write-Host "  DSN:     $Dsn64"
Write-Host "  Encrypt: No (local server)"
if ($UseTrustedConnection) {
  Write-Host "  Auth:    Trusted_Connection=Yes"
} else {
  Write-Host "  Auth:    SQL login ($SqlUser)"
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
Write-Host "Check: C:\Windows\System32\odbcad32.exe -> System DSN -> $Dsn64 -> Test"
Write-Host "Requires 64-bit QUIK (or other 64-bit ODBC client). 32-bit DSN is not created by this script."
