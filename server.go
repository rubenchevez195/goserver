package main

import (
    "fmt"
    "net/http"
    "encoding/json"
    "encoding/base64"
    //"encoding/binary"
    "strings"
    //"bytes"
    //"strconv"
    "io/ioutil"
    // "image"
    // "image/color"
    // "golang.org/x/image/bmp"
)
// import 
//import "errors"
// import "regexp"
type nSize struct{
   Alto  int     `json:"alto"`
   Ancho int     `json:"ancho"`
}
type bmpImg struct{
  Nombre string  `json:"nombre"`
  Data   string  `json:"data"`
  Size   nSize   `json:"tamaño"`
}
type Header struct {                //Area Size 14
    tipo,reserved1, reserved2     uint16
    size, offset                  int    
}
type InfoHeader struct {            //Area Size 40
    size, width,height                                        int 
    planes, bits                                              uint16
    compression, imagesize                                    int 
    xresolution, yresolution, ncolours, importantcolours      int
    
}
type ErrorMsg struct{
  Status string `json:"status"`
  Error  string `json:"error"`
}
func read_int(data []byte) (int) {
    var value int
    for _ ,element := range data {
        value |= int(element)
    }
    return value
}
func read_uint16(data []byte) (uint16) {
    var value uint16
    for _ ,element := range data {
        value |= uint16(element)
    }
    return value
}
func sendError(err string, w http.ResponseWriter){
  var errorMsg ErrorMsg 
  errorMsg.Status = "400"
  errorMsg.Error = err
  body, error := json.Marshal(errorMsg)
  if error != nil {
      panic(err)
  }
  w.WriteHeader(400)
  w.Write(body)

}
func decodeImage(baseImg string) []byte{
    data, err := base64.StdEncoding.DecodeString(baseImg)
    if err != nil {
        fmt.Println("error:", err)
        return nil
    }
    //fmt.Println("IMAGE DECODED", data)
    return data
}
func assignHeaders(h Header,  ih InfoHeader, data []byte ) (Header, InfoHeader){
    h.tipo      = read_uint16( data[0:4]   )         
    h.size      = read_int(    data[2:6]   )  
    h.reserved1 = read_uint16( data[6:8]   )   
    h.reserved2 = read_uint16( data[8:10]  ) 
    h.offset    = read_int(    data[10:14] ) 
    ih.size              =  read_int(      data[14:18]   ) 
    ih.width             =  read_int(      data[18:22]   ) 
    ih.height            =  read_int(      data[22:26]   )
    ih.planes            =  read_uint16(   data[26:28]   ) 
    ih.bits              =  read_uint16(   data[28:30]   )  
    ih.compression       =  read_int(      data[30:34]   ) 
    ih.imagesize         =  read_int(      data[34:38]   ) 
    ih.xresolution       =  read_int(      data[38:42]   ) 
    ih.yresolution       =  read_int(      data[42:46]   ) 
    ih.ncolours          =  read_int(      data[46:50]   ) 
    ih.importantcolours  =  read_int(      data[50:54]   ) 
    return h, ih
}
func Ejercicio3(w http.ResponseWriter, body []byte) {

    var img bmpImg
    err := json.Unmarshal(body, &img)
    if err != nil {
        sendError("ERROR CONVIRTIENDO EL JSON A BYTES", w)
    }
    data := decodeImage(img.Data) 

    fileErr := ioutil.WriteFile(img.Nombre, data, 0644)
    if(fileErr != nil ){
        panic(fileErr)
    }

    var h Header
    var ih InfoHeader
    h, ih =  assignHeaders(h, ih, data)
    fmt.Println(h,ih,  ih.bits)
    var headerOffset = 14 + 40 + 4*ih.ncolours
    var tmpData []byte
    for i := 0; ((i + int(ih.bits) < len(data)) ); i++ {
        if( i >= headerOffset ){
            var prom int
            for j := 0; j < int(ih.bits) ; j++ {
                prom += int(data[i+j])
            }
            prom += prom/ int(ih.bits)
            for j := 0; j < int(ih.bits) ; j++ {
                tmpData  = append(tmpData, byte(prom) )
            }
            i += int(ih.bits) - 1
        }else{
            tmpData  = append(tmpData , data[i])
        }
    }
    //fmt.Println(tmpData)
    arr1 := strings.Split(img.Nombre, ".")
    img.Nombre = arr1[0]+"(Blanco Y Negro)."+arr1[1]
    uEnc := base64.URLEncoding.EncodeToString(tmpData)
    img.Data = uEnc
    sendBody, error := json.Marshal(img)
    if error != nil {
        panic(error)
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(sendBody)
    fileErr2 := ioutil.WriteFile(img.Nombre, tmpData, 0644)
    if(fileErr2 != nil ){
        panic(fileErr2)
    }
}
func Ejercicio4(w http.ResponseWriter, body []byte) {
    var img bmpImg
    err := json.Unmarshal(body, &img)
    if err != nil {
        sendError("ERROR CONVIRTIENDO EL JSON A BYTES", w)
    }
    data := decodeImage(img.Data) 

    fileErr := ioutil.WriteFile(img.Nombre, data, 0644)
    if(fileErr != nil ){
        panic(fileErr)
    }
    var h Header
    var ih InfoHeader
    h, ih =  assignHeaders(h, ih, data)
    if(ih.width == 1){
        ih.width = 256
    }
    if(ih.height == 1){
        ih.height = 256   
    }
    fmt.Println(h,ih,  ih.bits)
    fmt.Println(img.Size)

    var cont int = 0
    output := make([][]byte, ih.height)
    for i := range output {
        output[i] = make([]byte, ih.width)
    }
    for i := 0; i < ih.width; i++ {
        for j := 0; j < ih.height; j++ {
            if(cont < len(data)){
                output[i][j] = data[cont]
            }
            //fmt.Printf( output[i][j] )
            cont++
        }
        //fmt.Println("-------------------")
    }
    var contH int = 0
    var contW int = 0
    tmp := make([][]byte, img.Size.Alto)
    for i := range tmp {
        tmp[i] = make([]byte, img.Size.Ancho)
    }
    for i := 0; i < img.Size.Ancho ; i++ {
        for j := 0; j < img.Size.Alto ; j++ {
            if(contW < ih.width && contH < ih.height){
                tmp[i][j] = (output[contW][contH] + output[contW+1][contH] + output[contW][contH+1] + output[contW+1][contH+1])/4
            }
            contW = contW + 2
            //fmt.Printf("%b",int(tmp[i][j]))
        }
        contH = contH + 2
        //fmt.Println("")
    }
    fmt.Println(tmp)

    //var headerOffset = 14 + 40 + 4*ih.ncolours

    
}


func handler(w http.ResponseWriter, r *http.Request) {
    //fmt.Fprintf(w, r.URL.Path[1:])
    defer r.Body.Close()
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
    }

    switch r.URL.Path[1:] {
      case "ejercicio1":

      case "ejercicio2":

      case "ejercicio3":
        Ejercicio3(w, body)
      case "ejercicio4":
        Ejercicio4(w, body)
    }

}

func main() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/ejercicio1", handler)
    http.HandleFunc("/ejercicio2", handler)
    http.HandleFunc("/ejercicio3", handler)
    http.HandleFunc("/ejercicio4", handler)
    http.ListenAndServe(":8080", nil)
}