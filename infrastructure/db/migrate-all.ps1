# migrate-all.ps1
# Скрипт для запуска миграций всех микросервисов

param(
    [string]$DbHost = "localhost",
    [int]$DbPort = 5432,
    [string]$DbUser = "postgres",
    [string]$DbPassword = "postgres"
)

$ErrorActionPreference = "Stop"

Write-Host "========================================" -ForegroundColor Green
Write-Host "  Running database migrations for all services" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green

# Получаем корень проекта
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptDir

# Конфигурация сервисов
$services = @(
    @{
        Name = "auth-server"
        DbName = "users"
        MigrationsPath = "$ProjectRoot\services\auth-server\cmd\db\migrations"
    },
    @{
        Name = "order-service"
        DbName = "orders"
        MigrationsPath = "$ProjectRoot\services\order-service\cmd\migrations"
    }
)

# Функция для запуска миграций
function Run-Migrations {
    param(
        [string]$ServiceName,
        [string]$DbName,
        [string]$MigrationsPath
    )
    
    Write-Host "`n>>> Running migrations for $ServiceName (database: $DbName)" -ForegroundColor Yellow
    
    if (-not (Test-Path $MigrationsPath)) {
        Write-Host "   Error: Migrations path not found: $MigrationsPath" -ForegroundColor Red
        return $false
    }
    
    $connectionString = "postgres://${DbUser}:${DbPassword}@${DbHost}:${DbPort}/${DbName}?sslmode=disable"
    
    & migrate -path $MigrationsPath -database $connectionString up
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "   ✓ Migrations completed for $ServiceName" -ForegroundColor Green
        return $true
    } else {
        Write-Host "   ✗ Migrations failed for $ServiceName" -ForegroundColor Red
        return $false
    }
}

Write-Host "`nDatabase connection: ${DbHost}:${DbPort}`n"

# Запускаем миграции для каждого сервиса
foreach ($service in $services) {
    $success = Run-Migrations -ServiceName $service.Name -DbName $service.DbName -MigrationsPath $service.MigrationsPath
    
    if (-not $success) {
        Write-Host "`nMigration failed for $($service.Name)" -ForegroundColor Red
        exit 1
    }
}

Write-Host "`n========================================" -ForegroundColor Green
Write-Host "  All migrations completed successfully!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
