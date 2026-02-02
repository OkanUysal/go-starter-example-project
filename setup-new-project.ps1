# Setup New Project Script
# This script helps you quickly set up a new project based on go-starter-example-project
#
# IMPORTANT: This script modifies the CURRENT directory!
# Usage:
#   1. Copy the template to a new location:
#      Copy-Item -Recurse go-starter-example-project your-new-project
#   2. Navigate to the new directory:
#      cd your-new-project
#   3. Run this script:
#      .\setup-new-project.ps1 -ProjectName "my-app" -ModulePath "github.com/user/my-app"

param(
    [Parameter(Mandatory=$true)]
    [string]$ProjectName,
    
    [Parameter(Mandatory=$true)]
    [string]$ModulePath,
    
    [Parameter(Mandatory=$false)]
    [string]$Description = "A Go REST API project"
)

# Safety check
Write-Host ""
Write-Host "⚠️  WARNING: This will modify the CURRENT directory!" -ForegroundColor Red
Write-Host "Current directory: $PWD" -ForegroundColor Yellow
Write-Host ""
$confirmation = Read-Host "Are you sure you want to continue? (yes/no)"
if ($confirmation -ne "yes") {
    Write-Host "❌ Setup cancelled" -ForegroundColor Red
    exit
}
Write-Host ""

Write-Host "🚀 Setting up new project: $ProjectName" -ForegroundColor Green
Write-Host "📦 Module path: $ModulePath" -ForegroundColor Cyan
Write-Host ""

# Constants
$OldModulePath = "github.com/OkanUysal/go-starter-example-project"
$OldTitle = "Go Starter Example Project"
$OldTablePrefix = "example"
$safeProjectName = $ProjectName.ToLower() -replace '[^a-z0-9_]', '_'
$NewTablePrefix = $safeProjectName

Write-Host "📋 Configuration:" -ForegroundColor Cyan
Write-Host "   Old table prefix: $OldTablePrefix" -ForegroundColor White
Write-Host "   New table prefix: $NewTablePrefix" -ForegroundColor White
Write-Host ""

# Step 1: Update go.mod
Write-Host "📝 Step 1/8: Updating go.mod..." -ForegroundColor Yellow
$goModContent = Get-Content "go.mod" -Raw
$goModContent = $goModContent -replace [regex]::Escape($OldModulePath), $ModulePath
Set-Content "go.mod" -Value $goModContent
Write-Host "✅ go.mod updated" -ForegroundColor Green
Write-Host ""

# Step 2: Update all Go files
Write-Host "📝 Step 2/7: Updating import paths in Go files..." -ForegroundColor Yellow
$goFiles = Get-ChildItem -Recurse -Filter *.go -Exclude "docs.go"
$fileCount = 0
foreach ($file in $goFiles) {
    $content = Get-Content $file.FullName -Raw
    if ($content -match [regex]::Escape($OldModulePath)) {
        $newContent = $content -replace [regex]::Escape($OldModulePath), $ModulePath
        Set-Content $file.FullName -Value $newContent
        $fileCount++
    }
}
Write-Host "✅ Updated $fileCount Go files" -ForegroundColor Green
Write-Host ""

# Step 3: Update migration files
Write-Host "📝 Step 3/8: Updating database migration files..." -ForegroundColor Yellow
$migrationFiles = Get-ChildItem -Path "migrations" -Filter *.sql -ErrorAction SilentlyContinue
$migrationCount = 0
foreach ($file in $migrationFiles) {
    $content = Get-Content $file.FullName -Raw
    if ($content -match "${OldTablePrefix}_") {
        $newContent = $content -replace "${OldTablePrefix}_user", "${NewTablePrefix}_user"
        $newContent = $newContent -replace "${OldTablePrefix}_token_blacklist", "${NewTablePrefix}_token_blacklist"
        Set-Content $file.FullName -Value $newContent
        $migrationCount++
    }
}
Write-Host "✅ Updated $migrationCount migration files" -ForegroundColor Green
Write-Host ""

# Step 4: Update Swagger docs in main.go
Write-Host "📝 Step 4/8: Updating Swagger documentation..." -ForegroundColor Yellow
$mainGoContent = Get-Content "main.go" -Raw
$mainGoContent = $mainGoContent -replace [regex]::Escape("@title           $OldTitle API"), "@title           $ProjectName API"
$mainGoContent = $mainGoContent -replace "@description     REST API for Go starter example project", "@description     $Description"
Set-Content "main.go" -Value $mainGoContent
Write-Host "✅ Swagger docs updated" -ForegroundColor Green
Write-Host ""

# Step 5: Regenerate Swagger docs
Write-Host "📝 Step 5/8: Regenerating Swagger documentation..." -ForegroundColor Yellow
$swagInstalled = Get-Command swag -ErrorAction SilentlyContinue
if ($swagInstalled) {
    swag init
    Write-Host "✅ Swagger docs regenerated" -ForegroundColor Green
} else {
    Write-Host "⚠️  swag command not found. Please install it:" -ForegroundColor Yellow
    Write-Host "   go install github.com/swaggo/swag/cmd/swag@latest" -ForegroundColor Cyan
    Write-Host "   Then run: swag init" -ForegroundColor Cyan
}
Write-Host ""

# Step 6: Update .env.example
Write-Host "📝 Step 6/8: Updating .env.example..." -ForegroundColor Yellow
$envContent = Get-Content ".env.example" -Raw
$envContent = $envContent -replace "SERVICE_NAME=go-starter-example-project", "SERVICE_NAME=$ProjectName"
$envContent = $envContent -replace "USER_TABLE=example_user", "USER_TABLE=${NewTablePrefix}_user"
$envContent = $envContent -replace "TOKEN_BLACKLIST_TABLE=example_token_blacklist", "TOKEN_BLACKLIST_TABLE=${NewTablePrefix}_token_blacklist"
Set-Content ".env.example" -Value $envContent

# Copy to .env if it doesn't exist
if (-not (Test-Path ".env")) {
    Copy-Item ".env.example" ".env"
    Write-Host "✅ Created .env file from template" -ForegroundColor Green
} else {
    Write-Host "ℹ️  .env file already exists, skipping" -ForegroundColor Cyan
}
Write-Host ""

# Step 7: Update README
Write-Host "📝 Step 7/8: Updating README.md..." -ForegroundColor Yellow
$readmeContent = Get-Content "README.md" -Raw
$readmeContent = $readmeContent -replace [regex]::Escape("# Go Starter Example Project"), "# $ProjectName"
$readmeContent = $readmeContent -replace "A production-ready Go REST API starter template.*", "$Description"
Set-Content "README.md" -Value $readmeContent
Write-Host "✅ README.md updated" -ForegroundColor Green
Write-Host ""

# Step 8: Run go mod tidy
Write-Host "📝 Step 8/8: Running go mod tidy..." -ForegroundColor Yellow
go mod tidy
Write-Host "✅ Dependencies updated" -ForegroundColor Green
Write-Host ""

# Summary
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "✨ Setup Complete!" -ForegroundColor Green
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host ""
Write-Host "📋 Next Steps:" -ForegroundColor Yellow
Write-Host ""
Write-Host "1. Update your .env file with actual values:" -ForegroundColor White
Write-Host "   - DATABASE_URL_LOCAL" -ForegroundColor Cyan
Write-Host "   - JWT_SECRET (generate: openssl rand -base64 32)" -ForegroundColor Cyan
Write-Host "   - GRAFANA_CLOUD_* (if using)" -ForegroundColor Cyan
Write-Host ""
Write-Host "2. Run database migrations:" -ForegroundColor White
Write-Host "   migrate -path migrations -database 'your-db-url' up" -ForegroundColor Cyan
Write-Host ""
Write-Host "3. Start the server:" -ForegroundColor White
Write-Host "   go run main.go" -ForegroundColor Cyan
Write-Host ""
Write-Host "4. Access the API:" -ForegroundColor White
Write-Host "   http://localhost:8080/health" -ForegroundColor Cyan
Write-Host "   http://localhost:8080/swagger/index.html" -ForegroundColor Cyan
Write-Host ""
Write-Host "5. Initialize git (if needed):" -ForegroundColor White
Write-Host "   git init" -ForegroundColor Cyan
Write-Host "   git add ." -ForegroundColor Cyan
Write-Host "   git commit -m 'Initial commit'" -ForegroundColor Cyan
Write-Host ""
Write-Host "Happy Coding! 🚀" -ForegroundColor Green
