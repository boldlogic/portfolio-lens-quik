param(
    [ValidateSet(
        "mixed",
        "get_positions",
        "get_money_positions",
        "get_security_positions",
        "get_otc_security_positions",
        "get_money_limits_by_client",
        "get_money_limits_paged",
        "get_security_limits_by_client",
        "get_security_limits_paged",
        "get_security_limits_otc_by_client",
        "get_security_limits_otc_paged",
        "get_limits",
        "post_money_limit",
        "post_security_limit",
        "post_security_limit_otc",
        "write_limits_mixed",
        "get_firms",
        "get_firm",
        "post_firm",
        "patch_firm"
    )]
    [string]$Scenario = "mixed",

    [string]$PortfolioBaseUrl = "http://127.0.0.1:5030",
    [string]$WriterBaseUrl = "http://127.0.0.1:5060",
    [string]$ReferenceDataBaseUrl = "http://127.0.0.1:6030",
    [string]$WriterApiKey = "",
    [int]$TargetRps = 1000,
    [int]$DurationSec = 30,
    [int]$HttpTargetRps = 0,
    [int]$PrepareFirmCount = 0,
    [int]$PrepareClientCount = 20,
    [string]$FirmsFixturePath = "",
    [switch]$SkipPrepare,
    [switch]$CreateMissingFirms,
    [switch]$IncludeDate
)

$ErrorActionPreference = "Stop"

. (Join-Path $PSScriptRoot "lib\Resolve-WriterConfig.ps1")
$writerResolved = Resolve-WriterK6Params -WriterBaseUrl $WriterBaseUrl -WriterApiKey $WriterApiKey
$WriterBaseUrl = $writerResolved.WriterBaseUrl
$WriterApiKey = $writerResolved.WriterApiKey

if ($DurationSec -lt 6) {
    throw "DurationSec dolzhen byt >= 6 (razogrev do 5 s + uderzhanie >= 1 s)"
}

$warmupSec = [Math]::Min(5, $DurationSec - 1)
$holdSec = $DurationSec - $warmupSec

$script:LimitHttpScenarioConfig = @{
    get_money_limits_by_client       = @{ Path = "/money-limits";           Mode = "by_client" }
    get_money_limits_paged           = @{ Path = "/money-limits";           Mode = "paged" }
    get_security_limits_by_client    = @{ Path = "/security-limits";        Mode = "by_client" }
    get_security_limits_paged        = @{ Path = "/security-limits";        Mode = "paged" }
    get_security_limits_otc_by_client = @{ Path = "/otc-security-limits";    Mode = "by_client" }
    get_security_limits_otc_paged    = @{ Path = "/otc-security-limits";    Mode = "paged" }
    get_limits                       = @{ Path = "/money-limits";           Mode = "triple_read" }
}

function Get-LimitProbeQueryString {
    param(
        [string]$Mode,
        [int]$PageLimit,
        [string]$ClientCode,
        [bool]$IncludeDate,
        [string]$Date
    )
    $parts = @()
    if ($Mode -eq "paged") {
        $parts += "limit=$PageLimit"
        $parts += "offset=0"
        $parts += "includeTotalCount=true"
    } elseif ($Mode -eq "by_client" -or $Mode -eq "triple_read") {
        $parts += "clientCodes=$([uri]::EscapeDataString($ClientCode))"
    }
    if ($IncludeDate -and $Date) {
        $parts += "loadDate=$([uri]::EscapeDataString($Date))"
    }
    return ($parts -join "&")
}

function Invoke-PortfolioLimitProbe {
    param(
        [string]$PortfolioBaseUrl,
        [string]$LimitsPath,
        [string]$Query,
        [int]$TimeoutSec = 15
    )
    $base = $PortfolioBaseUrl.TrimEnd("/")
    $suffix = if ($Query) { "$LimitsPath`?$Query" } else { $LimitsPath }
    $uri = "$base/api/v1/quik$suffix"
    try {
        return Invoke-RestMethod -Uri $uri -Method Get -TimeoutSec $TimeoutSec
    } catch {
        throw "Portfolio ne otvechaet (probe, test ne zapuskaetsya): $uri. $($_.Exception.Message)"
    }
}

function Resolve-LimitK6RatesFromHttpProbe {
    param(
        [int]$HttpTargetRps,
        [string]$Mode,
        [object]$ProbeResult,
        [int]$PageLimit
    )
    if ($Mode -eq "paged") {
        $total = [int]$ProbeResult.totalCount
        $pages = [Math]::Max(1, [Math]::Ceiling($total / $PageLimit))
        $iter = [double]$HttpTargetRps / $pages
        Write-Host "Probe totalCount=$total pageLimit=$PageLimit pagesPerIter=$pages HTTP target=${HttpTargetRps}/s k6 iter/s=$iter"
        return @{
            IterPerSec   = $iter
            TotalCount   = $total
            PagesPerIter = $pages
        }
    }
    if ($Mode -eq "triple_read") {
        $httpPerIter = 3
        $iter = [double]$HttpTargetRps / $httpPerIter
        Write-Host "Probe OK mode=triple_read HTTP target=${HttpTargetRps}/s k6 iter/s=$iter (3 GET/iter)"
        return @{
            IterPerSec   = $iter
            TotalCount   = 0
            PagesPerIter = $httpPerIter
        }
    }
    Write-Host "Probe OK mode=$Mode HTTP target=${HttpTargetRps}/s k6 iter/s=$HttpTargetRps (1 GET/iter)"
    return @{
        IterPerSec   = [double]$HttpTargetRps
        TotalCount   = 0
        PagesPerIter = 1
    }
}

$k6Root = $PSScriptRoot
$repoRoot = Resolve-Path (Join-Path $k6Root "..\..")
$prepareScript = Join-Path $k6Root "prepare\$Scenario.ps1"
$scenarioScript = Join-Path $k6Root "scenarios\$Scenario.js"
$reportsDir = Join-Path $repoRoot "docs\reports\k6"
$stateDir = Join-Path $k6Root "state"

if (-not (Test-Path $prepareScript)) {
    throw "Prepare script not found: $prepareScript"
}
if (-not (Test-Path $scenarioScript)) {
    throw "Scenario script not found: $scenarioScript"
}

New-Item -ItemType Directory -Path $reportsDir -Force | Out-Null
New-Item -ItemType Directory -Path $stateDir -Force | Out-Null

Write-Host ""
if ($HttpTargetRps -gt 0) {
    Write-Host "========== $Scenario : TargetRps (HTTP GET/s) = $HttpTargetRps ==========" -ForegroundColor Cyan
} else {
    Write-Host "========== $Scenario : TargetRps (k6 iter/s) = $TargetRps ==========" -ForegroundColor Cyan
}
Write-Host ""

Write-Host "Scenario: $Scenario"
Write-Host "DurationSec: $DurationSec (${warmupSec}s razogrev + ${holdSec}s uderzhanie)"
Write-Host "PortfolioBaseUrl: $PortfolioBaseUrl"
Write-Host "WriterBaseUrl: $WriterBaseUrl"
Write-Host "ReferenceDataBaseUrl: $ReferenceDataBaseUrl"

$limitCfg = $null
if ($HttpTargetRps -gt 0) {
    if (-not $script:LimitHttpScenarioConfig.ContainsKey($Scenario)) {
        throw "-HttpTargetRps tolko dlya limit read: $($script:LimitHttpScenarioConfig.Keys -join ', ')"
    }
    $limitCfg = $script:LimitHttpScenarioConfig[$Scenario]
    Write-Host "HttpTargetRps (HTTP GET/s): $HttpTargetRps"
} else {
    Write-Host "TargetRps (k6 iter/s): $TargetRps"
}

$prepareTargetRps = if ($HttpTargetRps -gt 0) { $HttpTargetRps } else { $TargetRps }

if (-not $SkipPrepare) {
    $prepareArgs = @{
        PortfolioBaseUrl       = $PortfolioBaseUrl
        WriterBaseUrl          = $WriterBaseUrl
        ReferenceDataBaseUrl   = $ReferenceDataBaseUrl
        WriterApiKey           = $WriterApiKey
        TargetRps              = $prepareTargetRps
        PrepareFirmCount       = $PrepareFirmCount
        PrepareClientCount     = $PrepareClientCount
        FirmsFixturePath       = $FirmsFixturePath
        CreateMissingFirms     = $CreateMissingFirms
    }
    if ($IncludeDate) {
        $prepareArgs.IncludeDate = $true
    }
    & $prepareScript @prepareArgs
}

$k6IterRps = [double]$TargetRps
$k6HttpRps = 0
$stateFile = Join-Path $stateDir "$Scenario.json"
if (Test-Path $stateFile) {
    $stateObj = Get-Content -Path $stateFile -Raw -Encoding UTF8 | ConvertFrom-Json
    if ($HttpTargetRps -gt 0) {
        $pageLimit = 50
        if ($limitCfg.Mode -eq "paged") {
            $pageLimit = [int]($stateObj.page_limit)
            if ($pageLimit -lt 1) { $pageLimit = 50 }
        }
        $includeDate = ($stateObj.include_date -eq $true)
        $date = [string]$stateObj.date
        $clientCode = [string]($stateObj.primary_client_code)
        if (($limitCfg.Mode -eq "by_client" -or $limitCfg.Mode -eq "triple_read") -and -not $clientCode) {
            throw "state $Scenario.json: net primary_client_code dlya probe by_client/triple_read"
        }
        $query = Get-LimitProbeQueryString `
            -Mode $limitCfg.Mode `
            -PageLimit $pageLimit `
            -ClientCode $clientCode `
            -IncludeDate:$includeDate `
            -Date $date
        $probeResult = Invoke-PortfolioLimitProbe `
            -PortfolioBaseUrl $PortfolioBaseUrl `
            -LimitsPath $limitCfg.Path `
            -Query $query
        $rates = Resolve-LimitK6RatesFromHttpProbe `
            -HttpTargetRps $HttpTargetRps `
            -Mode $limitCfg.Mode `
            -ProbeResult $probeResult `
            -PageLimit $pageLimit
        $k6IterRps = $rates.IterPerSec
        $k6HttpRps = $HttpTargetRps
        $stateObj | Add-Member -NotePropertyName target_http_rps -NotePropertyValue $k6HttpRps -Force
        $stateObj | Add-Member -NotePropertyName target_iter_rps -NotePropertyValue $k6IterRps -Force
        $stateObj | Add-Member -NotePropertyName target_rps -NotePropertyValue $k6HttpRps -Force
        $stateObj | Add-Member -NotePropertyName pages_per_iter -NotePropertyValue $rates.PagesPerIter -Force
        if ($rates.TotalCount -gt 0) {
            $stateObj | Add-Member -NotePropertyName probe_total_count -NotePropertyValue $rates.TotalCount -Force
        }
    } else {
        $stateObj | Add-Member -NotePropertyName target_rps -NotePropertyValue $TargetRps -Force
    }
    $stateObj | Add-Member -NotePropertyName test_duration_sec -NotePropertyValue $DurationSec -Force
    $stateObj | Add-Member -NotePropertyName warmup_duration_sec -NotePropertyValue $warmupSec -Force
    $stateObj | Add-Member -NotePropertyName hold_duration_sec -NotePropertyValue $holdSec -Force
    if (-not [string]::IsNullOrWhiteSpace($WriterApiKey)) {
        $stateObj | Add-Member -NotePropertyName writer_api_key -NotePropertyValue $WriterApiKey -Force
    }
    if (-not [string]::IsNullOrWhiteSpace($WriterBaseUrl)) {
        $stateObj | Add-Member -NotePropertyName writer_base_url -NotePropertyValue $WriterBaseUrl.TrimEnd("/") -Force
    }
    $stateObj | ConvertTo-Json -Depth 8 | Set-Content -Path $stateFile -Encoding UTF8
}

$iterRpsStr = $k6IterRps.ToString([System.Globalization.CultureInfo]::InvariantCulture)
if ($k6HttpRps -gt 0) {
    Write-Host "Effective k6 iter/s: $iterRpsStr (HTTP GET/s: $k6HttpRps)"
} else {
    Write-Host "Effective k6 target (iterations/s): $iterRpsStr"
}

$sidecarHttp = Join-Path $k6Root ".last-http-target-rps"
$sidecarIter = Join-Path $k6Root ".last-k6-iter-rps"
if ($k6HttpRps -gt 0) {
    Set-Content -Path $sidecarHttp -Value $k6HttpRps -Encoding ASCII -NoNewline
}
Set-Content -Path $sidecarIter -Value $iterRpsStr -Encoding ASCII -NoNewline
$sidecarDuration = Join-Path $k6Root ".last-test-duration-sec"
Set-Content -Path $sidecarDuration -Value $DurationSec -Encoding ASCII -NoNewline
Write-Host "Sidecar dlya k6: $sidecarHttp = $k6HttpRps; $sidecarIter = $iterRpsStr; $sidecarDuration = $DurationSec"

$reportsDirUnix = ($reportsDir -replace '\\', '/')

Push-Location $repoRoot
try {
    $k6Env = @(
        "-e", "K6_LAUNCHED_FROM_RUN_PS1=1",
        "-e", "K6_REPORTS_DIR=$reportsDirUnix",
        "-e", "K6_TARGET_RPS=$iterRpsStr",
        "-e", "K6_TEST_DURATION_SEC=$DurationSec"
    )
    if ($k6HttpRps -gt 0) {
        $k6Env += @("-e", "K6_HTTP_TARGET_RPS=$k6HttpRps")
    }
    k6 run @k6Env $scenarioScript
    $lastReport = Join-Path $reportsDir "last_$Scenario.html"
    if (Test-Path $lastReport) {
        Write-Host "HTML otchyot: $lastReport"
    }
    $newest = Get-ChildItem -Path $reportsDir -Filter "load_${Scenario}__*.html" -ErrorAction SilentlyContinue |
        Sort-Object LastWriteTime -Descending |
        Select-Object -First 1
    if ($newest) {
        Write-Host "Polnyj otchyot: $($newest.FullName)"
    }
}
finally {
    Pop-Location
}
