param(
    [Parameter(Mandatory = $true)][string]$Scenario,
    [string]$PortfolioBaseUrl = "http://127.0.0.1:5030",
    [string]$WriterBaseUrl = "http://127.0.0.1:5060",
    [Parameter(Mandatory = $true)][int]$TargetRps,
    [int]$DurationSec = 30,
    [int]$PrepareClientCount = 20,
    [string]$FirmsFixturePath = "",
    [switch]$SkipPrepare,
    [switch]$IncludeDate
)

$runArgs = @{
    Scenario           = $Scenario
    PortfolioBaseUrl   = $PortfolioBaseUrl
    WriterBaseUrl      = $WriterBaseUrl
    HttpTargetRps      = $TargetRps
    DurationSec        = $DurationSec
    PrepareClientCount = $PrepareClientCount
    FirmsFixturePath   = $FirmsFixturePath
    SkipPrepare        = $SkipPrepare
}
if ($IncludeDate) {
    $runArgs.IncludeDate = $true
}

& (Join-Path $PSScriptRoot "run.ps1") @runArgs
