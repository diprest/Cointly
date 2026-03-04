Write-Host "Starting Go8 Project for Windows..." -ForegroundColor Cyan

if (!(Get-Process "Docker Desktop" -ErrorAction SilentlyContinue)) {
    Write-Warning "Docker Desktop might not be running. Please ensure Docker is started."
}

docker-compose -f docker-compose.yml -f docker-compose.windows.yml up --build -d

if ($LASTEXITCODE -eq 0) {
    Write-Host "`nProject started successfully!" -ForegroundColor Green
    Write-Host "Frontend: http://localhost:3000"
    Write-Host "API Gateway: http://localhost:8000"
} else {
    Write-Host "`nFailed to start project." -ForegroundColor Red
}
