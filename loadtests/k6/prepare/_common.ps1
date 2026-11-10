$ErrorActionPreference = "Stop"

# Коды clientCode/firmCode/secCode в API не длиннее 12 символов.
# Собираем prefix + часть runId + индекс; если короче 12 — дополняем нулями, если длиннее — обрезаем.
function Format-LoadCode {
    param(
        [Parameter(Mandatory = $true)][string]$Prefix,
        [Parameter(Mandatory = $true)][string]$RunId,
        [Parameter(Mandatory = $true)][int]$Index,
        [int]$IndexWidth = 4
    )

    $runPartMax = 12 - $Prefix.Length - $IndexWidth
    if ($runPartMax -lt 1) {
        throw "Format-LoadCode: prefix '$Prefix' and index width $IndexWidth do not fit in 12 chars"
    }

    $runPart = ($RunId -replace '[^A-Za-z0-9]', '').ToUpper()
    if ($runPart.Length -gt $runPartMax) {
        $runPart = $runPart.Substring(0, $runPartMax)
    }

    $idxPart = $Index.ToString(("D{0}" -f $IndexWidth))
    $code = "$Prefix$runPart$idxPart"
    if ($code.Length -gt 12) {
        return $code.Substring(0, 12)
    }
    return $code.PadRight(12, '0')
}

function Get-PrepareRoot {
    $root = Resolve-Path (Join-Path $PSScriptRoot "..")
    $stateDir = Join-Path $root "state"
    $fixturesDir = Join-Path $root "fixtures"
    New-Item -ItemType Directory -Path $stateDir -Force | Out-Null
    return @{
        Root         = $root
        StateDir     = $stateDir
        FixturesDir  = $fixturesDir
        DefaultFirmsFixture = Join-Path $fixturesDir "firms.json"
    }
}

function Resolve-FirmsFixturePath {
    param([string]$FirmsFixturePath = "")

    $paths = Get-PrepareRoot
    if ([string]::IsNullOrWhiteSpace($FirmsFixturePath)) {
        return (Resolve-Path $paths.DefaultFirmsFixture).Path
    }
    if (-not (Test-Path $FirmsFixturePath)) {
        throw "Firms fixture not found: $FirmsFixturePath"
    }
    return (Resolve-Path $FirmsFixturePath).Path
}

function Read-FirmsFixtureFile {
    param([Parameter(Mandatory = $true)][string]$FixturePath)

    $raw = Get-Content -Path $FixturePath -Raw -Encoding UTF8 | ConvertFrom-Json
    $firms = @($raw.firms)
    if ($firms.Count -eq 0) {
        throw "Firms fixture is empty: $FixturePath"
    }
    foreach ($f in $firms) {
        if ([string]::IsNullOrWhiteSpace($f.firmCode)) {
            throw "Each firm in fixture must have firmCode: $FixturePath"
        }
    }
    return $firms
}

function Resolve-FirmFromFixtureEntry {
    param(
        [Parameter(Mandatory = $true)]$Entry,
        $ApiFirmsList = $null
    )

    $code = [string]$Entry.firmCode
    if ($null -ne $Entry.id -and "$($Entry.id)" -ne "") {
        return @{
            Id       = [int]$Entry.id
            FirmCode = $code
            FirmName = [string]$Entry.firmName
        }
    }

    if ($null -eq $ApiFirmsList) {
        throw "firmCode $code has no id in fixture; reference-data GET /firms/ was not loaded"
    }

    $found = $ApiFirmsList | Where-Object { $_.firmCode -eq $code } | Select-Object -First 1
    if ($null -eq $found) {
        throw "firmCode '$code' from fixture not found in reference-data; add firm to dbo.firms or set id in firms.json"
    }

    return @{
        Id       = [int]$found.id
        FirmCode = [string]$found.firmCode
        FirmName = [string]$found.firmName
    }
}

function Import-QuikFixtureFirms {
    param(
        [string]$ReferenceDataBaseUrl = "http://127.0.0.1:6030",
        [string]$FirmsFixturePath = "",
        [int]$MaxCount = 0,
        [switch]$CodesOnly
    )

    $fixturePath = Resolve-FirmsFixturePath -FirmsFixturePath $FirmsFixturePath
    $entries = Read-FirmsFixtureFile -FixturePath $fixturePath
    $useCount = if ($MaxCount -gt 0) { [Math]::Min($MaxCount, $entries.Count) } else { $entries.Count }

    $resolved = New-Object System.Collections.Generic.List[object]
    if ($CodesOnly) {
        for ($i = 0; $i -lt $useCount; $i++) {
            $entry = $entries[$i]
            [void]$resolved.Add(@{
                Id       = 0
                FirmCode = [string]$entry.firmCode
                FirmName = [string]$entry.firmName
            })
        }
        Write-Host "Fixture firms (firmCode only, no reference-data): $fixturePath (used $useCount of $($entries.Count))"
        return $resolved
    }

    $needsLookup = @($entries | Select-Object -First $useCount | Where-Object { $null -eq $_.id -or "$($_.id)" -eq "" })
    $apiFirms = $null
    if ($needsLookup.Count -gt 0) {
        $apiFirms = @(Invoke-JsonRequest -Method Get -Uri "$ReferenceDataBaseUrl/api/v1/quik/firms/")
    }

    for ($i = 0; $i -lt $useCount; $i++) {
        $firm = Resolve-FirmFromFixtureEntry -Entry $entries[$i] -ApiFirmsList $apiFirms
        [void]$resolved.Add($firm)
    }

    Write-Host "Fixture firms: $fixturePath (used $useCount of $($entries.Count))"
    return $resolved
}

function Invoke-JsonRequest {
    param(
        [Parameter(Mandatory = $true)][ValidateSet("Get", "Post", "Patch", "Put")][string]$Method,
        [Parameter(Mandatory = $true)][string]$Uri,
        [hashtable]$Body = $null,
        [string]$ApiKey = ""
    )

    $params = @{
        Method  = $Method
        Uri     = $Uri
        Headers = @{ Accept = "application/json" }
    }

    if (-not [string]::IsNullOrWhiteSpace($ApiKey)) {
        $params.Headers["X-API-KEY"] = $ApiKey
    }

    if ($null -ne $Body) {
        $params.ContentType = "application/json"
        $params.Body = ($Body | ConvertTo-Json -Compress)
    }

    try {
        return Invoke-RestMethod @params
    }
    catch {
        $resp = $_.Exception.Response
        if ($null -eq $resp) {
            throw "HTTP request failed (service unreachable?): $Method $Uri - $($_.Exception.Message)"
        }
        $reader = New-Object System.IO.StreamReader($resp.GetResponseStream())
        $detail = $reader.ReadToEnd()
        throw "HTTP $($resp.StatusCode.value__) $Uri : $detail"
    }
}

function New-LoadFirm {
    param(
        [Parameter(Mandatory = $true)][string]$ReferenceDataBaseUrl,
        [Parameter(Mandatory = $true)][string]$RunId,
        [Parameter(Mandatory = $true)][int]$Index
    )

    $firmCode = Format-LoadCode -Prefix "LD" -RunId $RunId -Index $Index
    $firmName = "Load Firm $Index $RunId"

    try {
        $created = Invoke-JsonRequest -Method Post -Uri "$ReferenceDataBaseUrl/api/v1/quik/firms/" -Body @{
            firmCode = $firmCode
            firmName = $firmName
        }
        return @{
            Id       = [int]$created.id
            FirmCode = [string]$created.firmCode
            FirmName = [string]$created.firmName
            Created  = $true
        }
    }
    catch {
        if ($_.Exception.Message -match "409") {
            $list = Invoke-JsonRequest -Method Get -Uri "$ReferenceDataBaseUrl/api/v1/quik/firms/"
            $existing = $list | Where-Object { $_.firmCode -eq $firmCode } | Select-Object -First 1
            if ($null -ne $existing) {
                return @{
                    Id       = [int]$existing.id
                    FirmCode = [string]$existing.firmCode
                    FirmName = [string]$existing.firmName
                    Created  = $false
                }
            }
        }
        throw
    }
}

$script:SettleCodes = @("T0", "T1", "T2", "Tx")

function Get-DefaultLimitCodesFixturePath {
    $paths = Get-PrepareRoot
    return Join-Path $paths.FixturesDir "limit_codes.json"
}

function Read-LimitCodesFixture {
    param([string]$FixturePath = "")

    $path = if ([string]::IsNullOrWhiteSpace($FixturePath)) {
        Get-DefaultLimitCodesFixturePath
    } else {
        $FixturePath
    }
    if (-not (Test-Path $path)) {
        throw "Limit codes fixture not found: $path"
    }
    $raw = Get-Content -Path $path -Raw -Encoding UTF8 | ConvertFrom-Json
    $money = @($raw.moneyCurrencies)
    $sec = @($raw.secCodes)
    $otc = @($raw.otcSecCodes)
    if ($money.Count -eq 0) {
        throw "limit_codes.json: moneyCurrencies is empty"
    }
    if ($sec.Count -eq 0) {
        throw "limit_codes.json: secCodes is empty"
    }
    if ($otc.Count -eq 0) {
        throw "limit_codes.json: otcSecCodes is empty"
    }
    return @{
        MoneyCurrencies = $money
        SecCodes        = $sec
        OtcSecCodes     = $otc
        FixturePath     = (Resolve-Path $path).Path
    }
}

function Test-WriterSeedKindEnabled {
    param(
        [Parameter(Mandatory = $true)][string]$Kind,
        [Parameter(Mandatory = $true)][string[]]$EnabledKinds
    )
    if ($EnabledKinds -contains "All") {
        return $true
    }
    return $EnabledKinds -contains $Kind
}

function New-WriterMoneyLimit {
    param(
        [Parameter(Mandatory = $true)][string]$WriterBaseUrl,
        [Parameter(Mandatory = $true)][string]$FirmCode,
        [Parameter(Mandatory = $true)][string]$ClientCode,
        [string]$Currency = "RUB",
        [string]$SettleCode = "Tx",
        [string]$ApiKey = ""
    )

    Invoke-JsonRequest -Method Put -Uri "$WriterBaseUrl/api/v1/limits" -ApiKey $ApiKey -Body @{
        limitType    = "money"
        clientCode   = $ClientCode
        ticker       = $Currency
        positionCode = "EQTV"
        settleCode   = $SettleCode
        firmCode     = $FirmCode
        balance      = "100.00"
    } | Out-Null
}

function New-WriterSecurityLimit {
    param(
        [Parameter(Mandatory = $true)][string]$WriterBaseUrl,
        [Parameter(Mandatory = $true)][string]$FirmCode,
        [Parameter(Mandatory = $true)][string]$ClientCode,
        [Parameter(Mandatory = $true)][string]$SecCode,
        [string]$SettleCode = "Tx",
        [string]$AcquisitionCurrencyCode = "RUB",
        [string]$ApiKey = ""
    )

    Invoke-JsonRequest -Method Put -Uri "$WriterBaseUrl/api/v1/limits" -ApiKey $ApiKey -Body @{
        limitType               = "security"
        clientCode              = $ClientCode
        ticker                  = $SecCode
        tradeAccount            = "TRADEACC01"
        settleCode              = $SettleCode
        firmCode                = $FirmCode
        balance                 = "10"
        acquisitionCurrencyCode = $AcquisitionCurrencyCode
    } | Out-Null
}

function New-WriterSecurityLimitOtc {
    param(
        [Parameter(Mandatory = $true)][string]$WriterBaseUrl,
        [Parameter(Mandatory = $true)][string]$FirmCode,
        [Parameter(Mandatory = $true)][string]$ClientCode,
        [Parameter(Mandatory = $true)][string]$SecCode,
        [string]$SettleCode = "Tx",
        [string]$AcquisitionCurrencyCode = "RUB",
        [string]$ApiKey = ""
    )

    Invoke-JsonRequest -Method Put -Uri "$WriterBaseUrl/api/v1/limits" -ApiKey $ApiKey -Body @{
        limitType               = "security_otc"
        clientCode              = $ClientCode
        ticker                  = $SecCode
        settleCode              = $SettleCode
        firmCode                = $FirmCode
        balance                 = "5"
        acquisitionCurrencyCode = $AcquisitionCurrencyCode
    } | Out-Null
}

function Build-QuikPrepareState {
    param(
        [string]$PortfolioBaseUrl = "http://127.0.0.1:5030",
        [string]$WriterBaseUrl = "http://127.0.0.1:5060",
        [string]$ReferenceDataBaseUrl = "http://127.0.0.1:6030",
        [int]$TargetRps = 1000,
        [int]$PrepareFirmCount = 0,
        [int]$PrepareClientCount = 20,
        [string]$FirmsFixturePath = "",
        [ValidateSet("Money", "Securities", "SecuritiesOtc", "All")]
        [string[]]$WriterSeedKinds = @("All"),
        [string]$WriterApiKey = "",
        [switch]$SkipWriterLimits,
        [switch]$CreateMissingFirms,
        [switch]$IncludeDate,
        [switch]$FixtureFirmsCodesOnly
    )

    $runId = [guid]::NewGuid().ToString("N").Substring(0, 8)
    $date = (Get-Date).ToString("yyyy-MM-dd")

    $firmIds = New-Object System.Collections.Generic.List[int]
    $firmCodes = New-Object System.Collections.Generic.List[string]
    $clientCodes = New-Object System.Collections.Generic.List[string]

    if ($CreateMissingFirms) {
        Write-Host "Prepare run_id=$runId mode=create_missing_firms (legacy)"
        $firmCount = [Math]::Max(1, $(if ($PrepareFirmCount -gt 0) { $PrepareFirmCount } else { 5 }))
        for ($i = 1; $i -le $firmCount; $i++) {
            $firm = New-LoadFirm -ReferenceDataBaseUrl $ReferenceDataBaseUrl -RunId $runId -Index $i
            [void]$firmIds.Add($firm.Id)
            [void]$firmCodes.Add($firm.FirmCode)
            Write-Host "  firm id=$($firm.Id) code=$($firm.FirmCode) created=$($firm.Created)"
        }
    }
    else {
        $loaded = Import-QuikFixtureFirms `
            -ReferenceDataBaseUrl $ReferenceDataBaseUrl `
            -FirmsFixturePath $FirmsFixturePath `
            -MaxCount $PrepareFirmCount `
            -CodesOnly:$FixtureFirmsCodesOnly
        foreach ($firm in $loaded) {
            if (-not $FixtureFirmsCodesOnly) {
                [void]$firmIds.Add($firm.Id)
            }
            [void]$firmCodes.Add($firm.FirmCode)
            if ($FixtureFirmsCodesOnly) {
                Write-Host "  firm code=$($firm.FirmCode) (fixture codes-only)"
            }
            else {
                Write-Host "  firm id=$($firm.Id) code=$($firm.FirmCode) (fixture)"
            }
        }
    }

    if ($firmCodes.Count -eq 0) {
        throw "No firms loaded for prepare; check firms.json"
    }

    $clientCount = [Math]::Max(1, $PrepareClientCount)
    Write-Host "Prepare run_id=$runId firms=$($firmCodes.Count) clients=$clientCount"

    $enabledKinds = @($WriterSeedKinds)
    if ($enabledKinds.Count -eq 0) {
        $enabledKinds = @("All")
    }
    $seedKindsLabel = $enabledKinds -join ","
    Write-Host "Prepare writer_seed_kinds=$seedKindsLabel"

    $limitCodes = Read-LimitCodesFixture
    Write-Host "Prepare limit_codes=$($limitCodes.FixturePath)"

    if (-not $SkipWriterLimits) {
        for ($i = 0; $i -lt $clientCount; $i++) {
            $firmCode = $firmCodes[$i % $firmCodes.Count]
            $seq = $i + 1
            $clientCode = Format-LoadCode -Prefix "LDCL" -RunId $runId -Index $seq
            $secCode = $limitCodes.SecCodes[$i % $limitCodes.SecCodes.Count]
            $otcSecCode = $limitCodes.OtcSecCodes[$i % $limitCodes.OtcSecCodes.Count]
            $settleCode = $script:SettleCodes[$i % $script:SettleCodes.Count]
            $currency = $limitCodes.MoneyCurrencies[$i % $limitCodes.MoneyCurrencies.Count]
            $acqCcy = if ($currency -eq "USD") { "USD" } else { "RUB" }

            if (Test-WriterSeedKindEnabled -Kind "Money" -EnabledKinds $enabledKinds) {
                New-WriterMoneyLimit -WriterBaseUrl $WriterBaseUrl -FirmCode $firmCode -ClientCode $clientCode -Currency $currency -SettleCode $settleCode -ApiKey $WriterApiKey
                Write-Host "  kind=Money client=$clientCode firm=$firmCode currency=$currency settle=$settleCode"
            }
            if (Test-WriterSeedKindEnabled -Kind "Securities" -EnabledKinds $enabledKinds) {
                New-WriterSecurityLimit -WriterBaseUrl $WriterBaseUrl -FirmCode $firmCode -ClientCode $clientCode -SecCode $secCode -SettleCode $settleCode -AcquisitionCurrencyCode $acqCcy -ApiKey $WriterApiKey
                Write-Host "  kind=Securities client=$clientCode firm=$firmCode secCode=$secCode settle=$settleCode"
            }
            if (Test-WriterSeedKindEnabled -Kind "SecuritiesOtc" -EnabledKinds $enabledKinds) {
                New-WriterSecurityLimitOtc -WriterBaseUrl $WriterBaseUrl -FirmCode $firmCode -ClientCode $clientCode -SecCode $otcSecCode -SettleCode $settleCode -AcquisitionCurrencyCode $acqCcy -ApiKey $WriterApiKey
                Write-Host "  kind=SecuritiesOtc client=$clientCode firm=$firmCode secCode=$otcSecCode settle=$settleCode"
            }
            [void]$clientCodes.Add($clientCode)
        }
    }

    $primaryClientCode = ""
    if ($clientCodes.Count -gt 0) {
        $primaryClientCode = $clientCodes[0]
    }

    $dateMode = if ($IncludeDate) { "include" } else { "omit" }

    return @{
        portfolio_base_url    = $PortfolioBaseUrl.TrimEnd("/")
        writer_base_url       = $WriterBaseUrl.TrimEnd("/")
        reference_base_url    = $ReferenceDataBaseUrl.TrimEnd("/")
        writer_api_key        = $WriterApiKey
        target_rps            = $TargetRps
        date                  = $date
        date_mode             = $dateMode
        include_date          = [bool]$IncludeDate
        run_id                = $runId
        firm_ids              = @($firmIds)
        firm_codes            = @($firmCodes)
        client_codes          = @($clientCodes)
        primary_client_code   = $primaryClientCode
        page_limit            = 50
        seed_client_count     = $clientCount
        writer_seed_kinds     = @($enabledKinds)
        fixture_firms_codes_only = [bool]$FixtureFirmsCodesOnly
        prepare_firm_count    = $firmCodes.Count
        prepare_client_count  = $clientCount
        firms_fixture_path    = (Resolve-FirmsFixturePath -FirmsFixturePath $FirmsFixturePath)
        sec_codes             = @($limitCodes.SecCodes)
        otc_sec_codes         = @($limitCodes.OtcSecCodes)
        money_currencies      = @($limitCodes.MoneyCurrencies)
    }
}

function Write-QuikState {
    param(
        [Parameter(Mandatory = $true)][hashtable]$State,
        [Parameter(Mandatory = $true)][string]$Scenario
    )

    $paths = Get-PrepareRoot
    $stateFile = Join-Path $paths.StateDir "$Scenario.json"
    $State | ConvertTo-Json -Depth 6 | Set-Content -Path $stateFile -Encoding UTF8
    Write-Host "State written: $stateFile"
}
