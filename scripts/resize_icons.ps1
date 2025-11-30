param(
    [string]$SourceFile = "..\web\Traefik_auth_manager.png",
    [string]$TargetDir = "..\web"
)

Add-Type -AssemblyName System.Drawing

$sourcePath = Resolve-Path $SourceFile
$targetPath = Resolve-Path $TargetDir
$sourceImage = [System.Drawing.Image]::FromFile($sourcePath)

function Resize-Image {
    param(
        [System.Drawing.Image]$Image,
        [int]$Width,
        [int]$Height,
        [string]$OutputPath
    )

    $rect = New-Object System.Drawing.Rectangle(0, 0, $Width, $Height)
    $destImage = New-Object System.Drawing.Bitmap($Width, $Height)
    
    $destImage.SetResolution($Image.HorizontalResolution, $Image.VerticalResolution)

    $graphics = [System.Drawing.Graphics]::FromImage($destImage)
    $graphics.CompositingMode = [System.Drawing.Drawing2D.CompositingMode]::SourceCopy
    $graphics.CompositingQuality = [System.Drawing.Drawing2D.CompositingQuality]::HighQuality
    $graphics.InterpolationMode = [System.Drawing.Drawing2D.InterpolationMode]::HighQualityBicubic
    $graphics.SmoothingMode = [System.Drawing.Drawing2D.SmoothingMode]::HighQuality
    $graphics.PixelOffsetMode = [System.Drawing.Drawing2D.PixelOffsetMode]::HighQuality

    $graphics.DrawImage($Image, $rect, 0, 0, $Image.Width, $Image.Height, [System.Drawing.GraphicsUnit]::Pixel)

    $destImage.Save($OutputPath, [System.Drawing.Imaging.ImageFormat]::Png)
    
    $graphics.Dispose()
    $destImage.Dispose()
    
    Write-Host "Generated $OutputPath"
}

# Ensure target directory exists
if (-not (Test-Path $TargetDir)) {
    New-Item -ItemType Directory -Force -Path $TargetDir | Out-Null
}

# Generate icons
Resize-Image -Image $sourceImage -Width 192 -Height 192 -OutputPath (Join-Path $targetPath "icon-192.png")
Resize-Image -Image $sourceImage -Width 512 -Height 512 -OutputPath (Join-Path $targetPath "icon-512.png")
Resize-Image -Image $sourceImage -Width 180 -Height 180 -OutputPath (Join-Path $targetPath "apple-touch-icon.png")

$sourceImage.Dispose()
Write-Host "Done."
