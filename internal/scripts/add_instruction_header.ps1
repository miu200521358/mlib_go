# Adds a fixed instruction header to newly added Go files detected by git.
# Usage: ./internal/scripts/add_instruction_header.ps1 [-DryRun]

param(
    [switch]$DryRun
)

$ErrorActionPreference = 'Stop'

function Get-RepoRoot {
    $root = & git rev-parse --show-toplevel 2>$null
    if (-not $root) {
        throw "git repository not found."
    }
    return $root.Trim()
}

$repoRoot = Get-RepoRoot

Push-Location $repoRoot
Write-Host "Repository root: $repoRoot"
try {
    $statusLines = & git status --porcelain -uall
} finally {
    Pop-Location
}

if (-not $statusLines) {
    Write-Host "No changes detected."
    return
}
Write-Host "Status lines: $($statusLines)"

$addedPaths = New-Object System.Collections.Generic.List[string]
foreach ($line in $statusLines) {
    if ($line -match '^\?\?\s+(.+)$') {
        $addedPaths.Add($Matches[1])
        continue
    }
    if ($line.Length -ge 3 -and $line[0] -eq 'A') {
        $addedPaths.Add($line.Substring(3))
        continue
    }
}

$expandedPaths = New-Object System.Collections.Generic.List[string]
foreach ($path in $addedPaths) {
    $fullPath = Join-Path $repoRoot $path
    if (Test-Path -Path $fullPath -PathType Container) {
        $files = Get-ChildItem -Path $fullPath -Recurse -File -Filter *.go
        foreach ($file in $files) {
            $relPath = $file.FullName.Substring($repoRoot.Length).TrimStart('\', '/')
            $relPath = $relPath -replace '\\', '/'
            $expandedPaths.Add($relPath)
        }
        continue
    }
    $expandedPaths.Add($path)
}

$addedGoPaths = $expandedPaths |
    Where-Object { $_ -like '*.go' } |
    Sort-Object -Unique

if (-not $addedGoPaths -or $addedGoPaths.Count -eq 0) {
    Write-Host "No added .go files detected."
    return
}

$instruction = '// ' + [char]0x6307 + [char]0x793A + ': miu200521358'
$utf8NoBom = New-Object System.Text.UTF8Encoding($false)

foreach ($relPath in $addedGoPaths) {
    $fullPath = Join-Path $repoRoot $relPath
    if (-not (Test-Path $fullPath)) {
        Write-Host "Skip (missing): $relPath"
        continue
    }

    $content = Get-Content -Raw -Path $fullPath
    if ($content -match [regex]::Escape($instruction)) {
        Write-Host "Skip (already has instruction): $relPath"
        continue
    }

    $eol = if ($content -match "`r`n") { "`r`n" } else { "`n" }
    $hasTrailingNewline = $content.EndsWith($eol)
    $lines = $content -split "`r?`n", -1

    $insertIndex = 0
    $sawBuildTag = $false
    while ($insertIndex -lt $lines.Length) {
        $line = $lines[$insertIndex]
        if ($line -match '^\s*//go:build\b') {
            $sawBuildTag = $true
            $insertIndex++
            continue
        }
        if ($line -match '^\s*// \+build\b') {
            $sawBuildTag = $true
            $insertIndex++
            continue
        }
        if ($sawBuildTag -and $line -match '^\s*$') {
            $insertIndex++
            continue
        }
        break
    }

    $before = if ($insertIndex -gt 0) { @($lines[0..($insertIndex - 1)]) } else { @() }
    $after = if ($insertIndex -lt $lines.Length) { @($lines[$insertIndex..($lines.Length - 1)]) } else { @() }
    $newLines = @()
    $newLines += $before
    $newLines += $instruction
    $newLines += $after
    $newContent = $newLines -join $eol
    if ($hasTrailingNewline) {
        $newContent += $eol
    }

    if ($DryRun) {
        Write-Host "DryRun: $relPath"
        continue
    }

    [System.IO.File]::WriteAllText($fullPath, $newContent, $utf8NoBom)
    Write-Host "Updated: $relPath"
}
