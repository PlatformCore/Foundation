param(
  [Parameter(Mandatory=$true)][string]$ModulePath,
  [Parameter(Mandatory=$true)][string]$Version
)

$ErrorActionPreference = "Stop"
if ($Version -notmatch '^v\d+\.\d+\.\d+$') {
  throw "Version must be semver format, for example v1.2.3"
}

$manifestPath = "tools/module_versions.json"
if (!(Test-Path $manifestPath)) {
  throw "Missing $manifestPath"
}

$json = Get-Content $manifestPath -Raw | ConvertFrom-Json
$found = $false
foreach ($m in $json.modules) {
  if ($m.path -eq $ModulePath) {
    $m.version = $Version
    $found = $true
    break
  }
}

if (-not $found) {
  throw "Module path not found in manifest: $ModulePath"
}

$json.updated_at = (Get-Date).ToString("yyyy-MM-dd")
$json | ConvertTo-Json -Depth 8 | Set-Content $manifestPath -Encoding utf8

Write-Host "Updated $ModulePath => $Version in $manifestPath"
Write-Host "Next: create tag $ModulePath/$Version and push."
