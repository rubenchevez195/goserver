package main

import (
    "fmt"
    //"errors"
    "net/http"
    "encoding/json"
    "encoding/base64"
    "strings"
    "strconv"
    "io/ioutil"
    // "image"
    // "image/color"
    // "golang.org/x/image/bmp"
    // "github.com/kr/pretty"
    // "log"
    "github.com/buger/jsonparser"
    "googlemaps.github.io/maps"
    "golang.org/x/net/context"

    // "golang.org/x/image/bmp"
    // "image"
    // "image/color"
    // "image/png"
    //"bytes"
    // "bufio"
    "encoding/binary"
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
type Coordinates struct{
  Lat   float64  `json:"lat"`
  Lng   float64  `json:"lng"`
}
type Routes struct{
  Locations []Coordinates `json:"ruta"`
}
type RouteParams struct{
  Origin string  `json:"origen"`
  Destin string  `json:"destino"`
}
type RestaurantParams struct{
  Origin string  `json:"origen"`
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
func Ejercicio1(w http.ResponseWriter, body []byte) error {
    var routes RouteParams
    err := json.Unmarshal(body, &routes)
    if err != nil {  return err  }
    c, _ := maps.NewClient(maps.WithAPIKey("AIzaSyAUOMk8n8nhxUiSTEXq06jNth_kiV_s55E"))
    r := &maps.DirectionsRequest{
        Origin:      routes.Origin,
        Destination: routes.Destin,
    }
    resp, _, _ := c.Directions(context.Background(), r)  
    body, er := json.Marshal(resp)
    if(er != nil){ return err }
    var route Routes
    jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        fmt.Println( "-------------------------------------" )
        var coordinate Coordinates
        var coordinate2 Coordinates
        lat,  _, _, _ := jsonparser.Get(value, "start_location", "lat")
        lng,  _, _, _  := jsonparser.Get(value, "start_location", "lng")
        lat2, _, _, _ := jsonparser.Get(value, "end_location", "lat")
        lng2, _, _, _ := jsonparser.Get(value, "end_location", "lng")
        f1, _ := strconv.ParseFloat( string(lat) , 64)
        f2, _ := strconv.ParseFloat( string(lng) , 64)
        f3, _ := strconv.ParseFloat( string(lat2) , 64)
        f4, _ := strconv.ParseFloat( string(lng2) , 64)
        coordinate.Lat = f1
        coordinate.Lng = f2
        coordinate2.Lat = f3
        coordinate2.Lng = f4
        route.Locations = append(route.Locations, coordinate)
        route.Locations = append(route.Locations, coordinate)
    }, "[0]","legs", "[0]", "steps" )
    sendBody, _ := json.Marshal(route)
    w.Header().Set("Content-Type", "application/json")
    w.Write(sendBody)
    return err
}
func Ejercicio2(w http.ResponseWriter, body []byte) error {
    var info RestaurantParams
    err := json.Unmarshal(body, &info)
    if( err != nil){ return err }
    c, _ := maps.NewClient(maps.WithAPIKey("AIzaSyAUOMk8n8nhxUiSTEXq06jNth_kiV_s55E"))
    r := &maps.TextSearchRequest{
        Query: "restaurants in "+info.Origin,
    }
    resp, _ := c.TextSearch(context.Background(), r)
    jsonBody, _ := json.Marshal(resp)
    var route Routes
    jsonparser.ArrayEach(jsonBody, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        var coordinate Coordinates
        fmt.Println( "-------------------------------------" )
        lat,  _, _, _ := jsonparser.Get(value, "geometry","location", "lat")
        lng,  _, _, _ := jsonparser.Get(value, "geometry","location", "lng")
        f1, _ := strconv.ParseFloat( string(lat) , 64)
        f2, _ := strconv.ParseFloat( string(lng) , 64)
        coordinate.Lat = f1
        coordinate.Lng = f2
        route.Locations = append(route.Locations, coordinate)
        fmt.Println(coordinate)
    }, "Results" )
    sendBody, _ := json.Marshal(route)
    w.Header().Set("Content-Type", "application/json")
    w.Write(sendBody)
    return err
}

func Ejercicio3(w http.ResponseWriter, body []byte) error {
    var img bmpImg
    err := json.Unmarshal(body, &img)
    if err != nil {
        return err
    }
    data := decodeImage(img.Data) 

    var h Header
    var ih InfoHeader
    h, ih =  assignHeaders(h, ih, data)
    var headerOffset = 14 + 40 + 4*ih.ncolours
    fmt.Println(h,ih,  ih.bits,headerOffset )
    var tmpData []byte
    for i := 0; ((i + int(ih.bits) < len(data)) ); i++ {
        if( i >= h.offset ){
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
    sendBody, errr := json.Marshal(img)
    if errr != nil {
        return errr
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(sendBody)
    fileErr2 := ioutil.WriteFile(img.Nombre, tmpData, 0644)
    if(fileErr2 != nil ){
        return fileErr2
    }
    return err
}
func Ejercicio4(w http.ResponseWriter, body []byte) error{
    var img bmpImg
    err := json.Unmarshal(body, &img)
    if err != nil  {
        return err
    }
    fmt.Println("Tamaño: ", img.Size)
    if img.Size.Ancho == 0  {
        sendError("Tamaño Deseadi Incorrecto, Asegurese de mandar un tamaño mayor a (cero, cero)", w)
        return err
    }
    data := decodeImage(img.Data) 
    var h Header
    var ih InfoHeader
    h, ih =  assignHeaders(h, ih, data)
    if(ih.width == 1){
        ih.width = 256
    }
    if(ih.height == 1){
        ih.height = 256   
    }
    fmt.Println( ih.width , " < ", img.Size.Ancho )
    if( ih.width < img.Size.Ancho ){
        sendError("Pides un tamaño mayor ", w)
        return err
    }
    fmt.Println(h,ih,  ih.bits )

    var sendData []byte
    var factor int = 2
    var currentW int = int(ih.width * int(ih.bits))  / 2
    for( (currentW/int(ih.bits)) > img.Size.Ancho ) {
        var tmpData []byte  
        var contx int = 0
        var conty int = 0
        var oldH int = int(ih.height * int(ih.bits)) 
        var oldW int = int(ih.width * int(ih.bits)) 
        var newH int = oldH / factor
        var newW int = oldW / factor
        fmt.Println("new W ",newW ,"new H " ,newH ,"GENERATED: ", currentW/int(ih.bits), " WANTED: ",  img.Size.Ancho, " factor: ", factor )
        currentW = oldW / factor
        
        input := make([][]byte, oldH )
        for i := range input {
            input[i] = make([]byte,  oldW )
        }

        output := make([][]byte, newH)
        for i := range output {
            output[i] = make([]byte,newW)
        }

        for i := 0; i < len(data) - 1 ; i++ {
            if( i < h.offset ){
                if(i ==  18){
                    b := make([]byte, 4)
                    binary.PutVarint(b,  int64(img.Size.Ancho))
                    //fmt.Println("Agregando Ancho: ", b)
                    for ii := 0; ii < len(b); ii++ {
                        tmpData  = append(tmpData , b[ii])    
                    }
                    i += 3
                }else if(i ==  22){
                    b := make([]byte, 4)
                    binary.PutVarint(b,  int64(img.Size.Alto))
                    //fmt.Println("Agregando Alto: ", b)
                    for ii := 0; ii < len(b); ii++ {
                        tmpData  = append(tmpData , b[ii])    
                    }
                    i += 3
                }else{
                    tmpData  = append(tmpData , data[i])
                }
            }else{
                if( contx < newW && conty < newH ){
                    //fmt.Printf(string(data[i]))
                    input[contx][conty]  = data[i]
                    //fmt.Println("input[contx][conty] ", int(input[contx][conty]) )
                    if(contx <  oldW ){
                        contx += 1
                    }else{
                        contx = 0
                        conty += 1
                        //fmt.Printf("")
                    }
                }
            }
        }

        for i := 0; i+1 < newH ; i++ {
            var k int = int(ih.bits) - 1
            for j := 0; j+1 < newW ; j++ {
                var newPosx int = i/2
                var newPosy int = j/2
                var prom int = ( int(input[i][j]) + int(input[i+k][j]) + int(input[i][j+k]) + int(input[i+k][j+k]) ) ;                          
                output[newPosx][newPosy]   = byte( prom / 4 );                     
                fmt.Println("prom ", prom,"output ", int(output[newPosx][newPosy])," [i][j] ", int(input[i][j])," [i+1][j] ", int(input[i+1][j])," [i][j+1] ", int(input[i][j+1])," [i+1][j+1] ", int(input[i+1][j+1]))
                j+= 1
            }
            i+= 1
        }
        
        for i := 0; i < len(tmpData); i++ {
            sendData = append(sendData, tmpData[i] )
        }
        for i := 0; i < newH; i++ {
            for j := 0; j < newW; j++ {
                sendData = append(sendData, output[i][j] )
            }
        }
        break
        
        factor += 2

    }
   
    arr1 := strings.Split(img.Nombre, ".")
    img.Nombre = arr1[0]+"(Blanco Y Negro)."+arr1[1]
    uEnc := base64.URLEncoding.EncodeToString(sendData)
    img.Data = uEnc
    sendBody, errorss := json.Marshal(img)
    if errorss != nil {
        return errorss
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(sendBody)
    fileErr2 := ioutil.WriteFile(img.Nombre, sendData, 0644)
    if(fileErr2 != nil ){
        return fileErr2
    }
    return err
}
func handler(w http.ResponseWriter, r *http.Request) {
    //fmt.Fprintf(w, r.URL.Path[1:])
    defer r.Body.Close()
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
    }
    switch r.URL.Path[1:] {
      case "ejercicio1":
        err := Ejercicio1(w, body)
        if(err != nil){
            sendError("Error en Ejercicio 1 Revise el formato de la solicitud este bien escrita", w)
        }
      case "ejercicio2":
        err := Ejercicio2(w, body)
        if(err != nil){
            sendError("Error en Ejercicio 2 Revise el formato de la solicitud este bien escrita", w)
        }
      case "ejercicio3":
        err := Ejercicio3(w, body)
        if(err != nil){
            sendError("Error en Ejercicio 3 Revise el formato de la solicitud este bien escrita", w)
        }
      case "ejercicio4":
        err := Ejercicio4(w, body)
        if(err != nil){
            sendError("Error en Ejercicio 3 Revise el formato de la solicitud este bien escrita", w)
        }
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

// "https://maps.googleapis.com/maps/api/directions/json?origin="+address1+"&destination="+address2+"&key=AIzaSyAUOMk8n8nhxUiSTEXq06jNth_kiV_s55E"