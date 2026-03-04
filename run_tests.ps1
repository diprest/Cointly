$ServiceDirs = @(
    "services/portfolio-service"
    "services/trading-service"
)

$AllTestsPassed = $true

function Run-Tests-For-Service {
    param(
        [string]$Path
    )
    
    Write-Host "`n------------------------------------------" -ForegroundColor Yellow
    Write-Host "Running tests for module: $Path" -ForegroundColor Green
    Write-Host "------------------------------------------" -ForegroundColor Yellow

    $TestCommand = "go test ./internal/service -v -cover"
    Write-Host "Executing: $TestCommand" -ForegroundColor Cyan
    $TestOutput = Invoke-Expression $TestCommand 2>&1
    
    if ($TestOutput -match "FAIL") {
        Write-Host "FAILED: Service tests failed in $Path" -ForegroundColor Red
        Write-Host $TestOutput
        return $false
    }
    
    $CoverageMatch = $TestOutput | Select-String "coverage: (\d+\.\d+)%"
    if ($CoverageMatch) {
        $Coverage = [double]$CoverageMatch.Matches[0].Groups[1].Value
        if ($Coverage -ge 30.0) {
            Write-Host "PASSED: Coverage is $($Coverage)%" -ForegroundColor Green
        } else {
            Write-Host "FAILED: Coverage is too low ($($Coverage)%), minimum required is 30%" -ForegroundColor Red
            return $false
        }
    } else {
        Write-Host "PASSED: Unit tests succeeded." -ForegroundColor Green
    }
    return $true
}

$OriginalPath = Get-Location

foreach ($Dir in $ServiceDirs) {
    Set-Location "$OriginalPath\$Dir"
    
    $Result = Run-Tests-For-Service $Dir
    
    if (-not $Result) {
        $AllTestsPassed = $false
    }
}

Set-Location $OriginalPath

Write-Host "`n==========================================" -ForegroundColor Magenta
if ($AllTestsPassed) {
    Write-Host "ALL REQUIRED TESTS PASSED!" -ForegroundColor Green
} else {
    Write-Host "BUILD FAILED: One or more test suites failed." -ForegroundColor Red
}
Write-Host "==========================================" -ForegroundColor Magenta