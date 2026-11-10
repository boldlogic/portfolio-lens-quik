$ErrorActionPreference = "Stop"

function Get-WriterK6ConfigFromYaml {
    param(
        [string]$ConfigPath = ""
    )

    if ([string]::IsNullOrWhiteSpace($ConfigPath)) {
        $repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..\..\..")
        $ConfigPath = Join-Path $repoRoot "configs\quik-portfolio-writer-config.yaml"
    }

    if (-not (Test-Path $ConfigPath)) {
        return @{
            ApiKey  = ""
            BaseUrl = ""
            Path    = $ConfigPath
        }
    }

    $content = Get-Content -Path $ConfigPath -Raw -Encoding UTF8
    $apiKey = ""
    $port = 0

    if ($content -match '(?m)^api_key:\s*"?([^"\r\n#]+)"?\s*$') {
        $apiKey = $Matches[1].Trim()
    }
    if ($content -match '(?m)^\s*listen_port:\s*(\d+)\s*$') {
        $port = [int]$Matches[1]
    }

    $baseUrl = ""
    if ($port -gt 0) {
        $baseUrl = "http://127.0.0.1:$port"
    }

    return @{
        ApiKey  = $apiKey
        BaseUrl = $baseUrl
        Path    = (Resolve-Path $ConfigPath).Path
    }
}

function Resolve-WriterK6Params {
    param(
        [string]$WriterBaseUrl,
        [string]$WriterApiKey
    )

    $fromYaml = Get-WriterK6ConfigFromYaml
    $resolvedKey = $WriterApiKey
    if ([string]::IsNullOrWhiteSpace($resolvedKey)) {
        $resolvedKey = [string]$env:WRITER_API_KEY
    }
    if ([string]::IsNullOrWhiteSpace($resolvedKey)) {
        $resolvedKey = $fromYaml.ApiKey
    }

    $resolvedUrl = $WriterBaseUrl
    if ($resolvedUrl -eq "http://127.0.0.1:5060" -and -not [string]::IsNullOrWhiteSpace($fromYaml.BaseUrl)) {
        $resolvedUrl = $fromYaml.BaseUrl
    }

    if (-not [string]::IsNullOrWhiteSpace($resolvedKey) -and (Test-Path $fromYaml.Path)) {
        Write-Host "Writer config: $($fromYaml.Path) (api_key loaded, base=$resolvedUrl)"
    } elseif ([string]::IsNullOrWhiteSpace($resolvedKey)) {
        Write-Host "Writer api_key: not set (pass -WriterApiKey, env WRITER_API_KEY, or configs/quik-portfolio-writer-config.yaml)"
    }

    return @{
        WriterBaseUrl = $resolvedUrl
        WriterApiKey  = $resolvedKey
    }
}
