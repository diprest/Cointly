$TradingUrl = "http://localhost:8082/api/v1"
$PortfolioUrl = "http://localhost:8083/api/v1"

function Print-Section ($text) {
    Write-Host "`n==========================================" -ForegroundColor Cyan
    Write-Host " $text" -ForegroundColor Cyan
    Write-Host "==========================================" -ForegroundColor Cyan
}
function Print-Pass ($text) { Write-Host "[PASS] $text" -ForegroundColor Green }
function Print-Fail ($text) { Write-Host "[FAIL] $text" -ForegroundColor Red }
function Print-Info ($text) { Write-Host "[INFO] $text" -ForegroundColor Gray }

function Set-Price ($price) {
    $json = '{"symbol": "BTCUSDT", "price": "' + $price + '"}'
    $json | docker exec -i kafka kafka-console-producer --bootstrap-server kafka:29092 --topic market_prices
    Print-Info ">> Kafka Price Update: $price (Triggering Matcher)"
    Start-Sleep -Seconds 3
}

function Get-Balance ($uid, $asset) {
    try {
        $res = Invoke-RestMethod -Uri "$PortfolioUrl/portfolio/balance?user_id=$uid&asset=$asset" -ErrorAction Stop
        return $res
    } catch {
        return @{amount=0; locked_bal=0}
    }
}

Print-Section "SETUP: Resetting System"
docker stop market-data | Out-Null
docker exec -i go8_project-postgres-1 psql -U user -d trading_db -c "TRUNCATE TABLE orders RESTART IDENTITY;" | Out-Null
Print-Info "Reseting Balances: User 1 gets 20,000 USDT & 1.0 BTC"
$sql = "TRUNCATE TABLE balances; 
INSERT INTO balances (user_id, asset, amount, locked_bal, total_cost) VALUES (1, 'USDT', 20000, 0, 0); 
INSERT INTO balances (user_id, asset, amount, locked_bal, total_cost) VALUES (1, 'BTC', 1.0, 0, 0);"
docker exec -i go8_project-postgres-1 psql -U user -d portfolio_db -c $sql | Out-Null
Start-Sleep -Seconds 2

Print-Section "TEST 1: Market Buy 0.1 BTC (Quantity)"
$body = @{user_id=1; symbol="BTCUSDT"; side="BUY"; type="MARKET"; amount=0.1; price=55000} | ConvertTo-Json
Invoke-RestMethod -Uri "$TradingUrl/orders" -Method Post -Body $body -ContentType "application/json" | Out-Null
Print-Info "1. Order Placed: Buy 0.1 BTC (Pending...)"
Set-Price "50000.00"
$usdt = Get-Balance 1 "USDT"
if ([double]$usdt.amount -eq 15000) { 
    Print-Pass "Balance OK: 15000 USDT" 
} else { 
    Print-Fail "Balance WRONG: $($usdt.amount) (Expected 15000)" 
}

Print-Section "TEST 2: Market Sell 0.5 BTC"
$body = @{user_id=1; symbol="BTCUSDT"; side="SELL"; type="MARKET"; amount=0.5; price=0} | ConvertTo-Json
Invoke-RestMethod -Uri "$TradingUrl/orders" -Method Post -Body $body -ContentType "application/json" | Out-Null
Print-Info "1. Order Placed: Sell 0.5 BTC (Pending...)"
Set-Price "60000.00"
$usdt = Get-Balance 1 "USDT"
if ([double]$usdt.amount -eq 45000) { 
    Print-Pass "Balance OK: 45000 USDT" 
} else { 
    Print-Fail "Balance WRONG: $($usdt.amount) (Expected 45000)" 
}

Print-Section "TEST 3: Market Buy for 5000 USDT (Quote)"
$body = @{user_id=1; symbol="BTCUSDT"; side="BUY"; type="MARKET"; quote_amount=5000} | ConvertTo-Json
Invoke-RestMethod -Uri "$TradingUrl/orders" -Method Post -Body $body -ContentType "application/json" | Out-Null
Print-Info "1. Order Placed: Spend 5000 USDT (Pending...)"
Set-Price "50000.00"
$usdt = Get-Balance 1 "USDT"
$btc = Get-Balance 1 "BTC"
if ([double]$usdt.amount -eq 40000) { Print-Pass "USDT OK: 40000" } else { Print-Fail "USDT WRONG: $($usdt.amount)" }
if ([math]::Round([double]$btc.amount, 2) -eq 0.7) { 
    Print-Pass "BTC OK: 0.7" 
} else { 
    Print-Fail "BTC WRONG: $($btc.amount)" 
}

Print-Section "ALL TESTS COMPLETED"