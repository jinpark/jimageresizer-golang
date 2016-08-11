package main

import (
    "fmt"
    "net/http"
    "time"
    "net/url"
    "strconv"
    "bytes"
    "log"
    "image"
    "image/jpeg"
    _ "image/png"
    _ "image/gif"

    "github.com/disintegration/imaging"
)

var client = http.Client{Timeout: time.Duration(5 * time.Second)}

func cutterHandler(res http.ResponseWriter, req *http.Request) {

    queryParams, _ := url.ParseQuery(req.URL.RawQuery)
    reqImg, err := client.Get(fmt.Sprint(queryParams["url"][0]))
    if err != nil {
        fmt.Fprintf(res, "Sorry, the URL you provided is not valid")
        log.Println("Error getting image %d", err)
        return
    }

    m, format, err := image.Decode(reqImg.Body)
    if err != nil {
        fmt.Fprintf(res, "Sorry, the image you provided is not valid")
        // log.Printf(res, "Error decoding image %d", err)
        return
    }
    fmt.Println("format %s", format)
    g := m.Bounds()

    // Get height and width
    origHeight := g.Dy()
    origWidth := g.Dx()


    width := 0 
    height := 0 
    if len(queryParams["width"]) > 0 {
        width, err = strconv.Atoi(queryParams["width"][0])
        if err != nil || width < 0 {
            fmt.Fprintf(res, "Sorry, the width you provided is not valid")
            // log.Printf(res, "Error decoding image %d", err)
            return
        }
    }
    if len(queryParams["height"]) > 0 {
        height, err = strconv.Atoi(queryParams["height"][0])
        if err != nil || height < 0 {
            fmt.Fprintf(res, "Sorry, the height you provided is not valid")
            // log.Printf(res, "Error decoding image %d", err)
            return
        }
    }
    if width == 0 && height == 0 {
        width = origWidth
        height = origHeight
    }
    fmt.Println("widthheight", width, height)
    resizedImg := imaging.Resize(m, width, height, imaging.Box)

    buffer := new(bytes.Buffer)
    if err := jpeg.Encode(buffer, resizedImg, nil); err != nil {
        fmt.Println("unable to encode image.")
    }


    fmt.Printf("res %d %d", height, width)
    res.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
    res.Header().Set("Content-Type", reqImg.Header.Get("Content-Type"))

    if _, err := res.Write(buffer.Bytes()); err != nil {
        fmt.Println("unable to write image.")
    }
    reqImg.Body.Close()
}

func main() {
    http.HandleFunc("/cut", cutterHandler)
    http.ListenAndServe(":8080", nil) /* TODO Configurable */
}
