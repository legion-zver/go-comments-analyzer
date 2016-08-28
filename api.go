package main

import (
    "strings"
    "regexp"    
    "github.com/kataras/iris"
    "github.com/legion-zver/shield"
)

var rxStripTags = regexp.MustCompile("<[a-zA-Z/][^>]*>")

// LearnRequest - learn request
type LearnRequest struct {
    Class string `form:"class" json:"class"`
    Text  string `form:"text" json:"text"`
}

// Learn - обучение
func Learn(c *iris.Context) {
    var sh shield.Shield
    if c.Get("shield") != nil {
        sh = c.Get("shield").(shield.Shield)
    }
    if sh != nil {
        var request LearnRequest
        err := c.ReadJSON(&request)
        if err != nil {
            err = c.ReadForm(&request)
            if err != nil {
                c.JSON(iris.StatusBadRequest, iris.Map{"error":NewAPIError(400,[]APIErrorCause{APIErrorCause{Target:"request",Сause:err.Error()}})})
                return
            }
        }
        request.Class = strings.TrimSpace(strings.ToLower(request.Class))
        request.Text  = strings.TrimSpace(rxStripTags.ReplaceAllString(request.Text, " "))
        if len(request.Text) > 0 {            
            err = sh.Learn(request.Class, request.Text)        
            if err != nil {
                c.JSON(iris.StatusInternalServerError, iris.Map{"error":NewAPIError(500,[]APIErrorCause{APIErrorCause{Target:"learn",Сause:err.Error()}})})
                return
            }
            // ok
            c.JSON(iris.StatusOK, iris.Map{"result":"success"})
            return
        }
        c.JSON(iris.StatusBadRequest, iris.Map{"error":NewAPIError(400,[]APIErrorCause{APIErrorCause{Target:"request.text",Сause:"Empty"}})})
        return
    }    
    c.JSON(iris.StatusServiceUnavailable, iris.Map{"error":NewSimpleAPIError(503)})
} 

// ScoreRequest - score request
type ScoreRequest struct {    
    Text  string `form:"text" json:"text"`
}

// Score - обучение
func Score(c *iris.Context) {
    var sh shield.Shield
    if c.Get("shield") != nil {
        sh = c.Get("shield").(shield.Shield)
    }
    if sh != nil {
        text := strings.TrimSpace(c.URLParam("text"))
        if len(text) < 1 {
            var request ScoreRequest
            err := c.ReadJSON(&request)
            if err != nil {
                err = c.ReadForm(&request)
                if err != nil {
                    c.JSON(iris.StatusBadRequest, iris.Map{"error":NewAPIError(400,[]APIErrorCause{APIErrorCause{Target:"request",Сause:err.Error()}})})
                    return
                }
            }
            text = strings.TrimSpace(request.Text)
        }
        if len(text) > 0 {
            m, err := sh.Score(text)
            if err != nil {
                c.JSON(iris.StatusInternalServerError, iris.Map{"error":NewAPIError(500,[]APIErrorCause{APIErrorCause{Target:"learn",Сause:err.Error()}})})
                return
            }
            // ok
            c.JSON(iris.StatusOK, iris.Map{"result":m})
            return
        }
        c.JSON(iris.StatusBadRequest, iris.Map{"error":NewAPIError(400,[]APIErrorCause{APIErrorCause{Target:"request.text",Сause:"Empty"}})})
        return
    }
    c.JSON(iris.StatusServiceUnavailable, iris.Map{"error":NewSimpleAPIError(503)})
}

// InitAPI - init API
func InitAPI(api iris.MuxAPI) {
    if api != nil {
        // Learn
        api.Get("/learn", Learn)
        api.Post("/learn", Learn)
        // Score
        api.Get("/score", Score)
        api.Post("/score", Score)
    }
}