# ASCII Art Converter

Convert images to ASCII art using Go and a simple web interface.

## Features
- Upload images (JPEG, PNG)
- Adjust output width (20-500 characters)
- Choose from two ASCII gradients (default, reverse)
- Copy ASCII art to clipboard
- Pure HTML frontend (no JavaScript required)

## Technologies
- **Backend**: Go (standard library)
- **Frontend**: HTML
- **Image Processing**: Grayscale conversion, resizing

## How to Run
1. Clone the repository:
```bash
git clone https://github.com/your-username/ascii-art-converter.git
```
2. Navigate to the project directory:
```bash
cd ascii-art-converter
```
3. Run the server:
```bash
go run main.go
```
4. Open `http://localhost:8080` in your browser.

## API Endpoints
- **POST /convert**: Convert an image to ASCII art.
    
    - Parameters:
        
        - `image`: The image file to convert.
            
        - `width`: The desired width of the ASCII art (20-500).
            
        - `style`: The gradient style (`default` or `reverse`).
            
    - Response: Plain text ASCII art.
