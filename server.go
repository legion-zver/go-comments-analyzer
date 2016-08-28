package main

import (
    "fmt"  
    "strconv"  
    "io/ioutil"
    "encoding/json"
    "path/filepath"    
    "github.com/kataras/iris"
    "github.com/kardianos/osext"    
    "github.com/legion-zver/shield"    
    "github.com/iris-contrib/middleware/cors"
    irisConfig "github.com/kataras/iris/config"
)


// Config - config app
type Config struct {
    Port        int64   `json:"port"`
    CurrentDir  string  `json:"-"`
}

func main()  {
    exeFile, _ := osext.Executable()
    currentDir, err := filepath.Abs(filepath.Dir(exeFile))
    if err != nil {
        fmt.Println(err)
        return
    }
    config := Config{Port: 1234}    
    if data, err := ioutil.ReadFile(currentDir+"/config.json"); err == nil {
        err = json.Unmarshal(data, &config)
        if err != nil {
            fmt.Println(err)
        }
    }
    config.CurrentDir = currentDir
    // Init NN
    sh := shield.New(shield.NewRussianTokenizer(), shield.NewLevelDBStore(currentDir+"/db"))
    // Config iris
    sessionsConfig := irisConfig.DefaultSessions()
    sessionsConfig.DisableSubdomainPersistence = true
    irisConfig := irisConfig.Iris{DisableBanner: true, IsDevelopment: true, Sessions: sessionsConfig }
    // Create API Iris                    
    api := iris.New(irisConfig) 
    api.Use(cors.Default())
    api.UseFunc(func (c *iris.Context){            
        c.Set("config", &config)
        c.Set("shield", sh)
        c.Next()
    })
    // 404
    api.OnError(iris.StatusNotFound, func(c *iris.Context){
        c.JSON(iris.StatusNotFound, iris.Map{"error":NewSimpleAPIError(404)})
    })
    // InitAPI    
    InitAPI(api.Party("/api"))
    // Run
    api.Listen(":"+strconv.FormatInt(int64(config.Port),10))
}