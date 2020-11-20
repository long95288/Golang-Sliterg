package main

import (
    "encoding/json"
    "fmt"
    "github.com/therecipe/qt/widgets"
    "io/ioutil"
    "os"
)

type config struct {
    BgImage string `json:"bgImage"`
    AppIcon string `json:"appIcon"`
    
}
func (c *config) init() {
    data,err := ioutil.ReadFile("conf.json")
    if err != nil {
        c.BgImage = "bg.jpg"
        c.AppIcon="app.png"
        return
    }
    err = json.Unmarshal(data,c)
    if err != nil {
        c.BgImage = "bg.jpg"
        c.AppIcon = "app.png"
    }
}
type SnakeNode struct{
    x int
    y int
    color int
    image string
    next *SnakeNode
}
type Snake struct {
    head *SnakeNode
    tail *SnakeNode
}

func (this *Snake) AppendNodeToHeader(node *SnakeNode) {
    this.head.next = node
    this.head = node
}
func (this *Snake) DeleteTail(){
    if nil != this.tail.next {
        temp := this.tail
        temp.next = nil
        this.tail = this.tail.next
    }
}

var (
    app *widgets.QMainWindow
    configuration config
)

func InitUI() *widgets.QMainWindow {
    app = widgets.NewQMainWindow(nil, 0)
    app.SetWindowTitle("贪吃蛇")
    
    return app
}
func main() {
    //
    fmt.Println("Hello world")
    widgets.NewQApplication(len(os.Args),os.Args)
    app := InitUI()
    app.Show()
    widgets.QApplication_Exec()
    
}
