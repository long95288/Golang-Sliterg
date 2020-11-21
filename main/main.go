package main

import (
    "encoding/json"
    "fmt"
    "github.com/therecipe/qt/core"
    "github.com/therecipe/qt/gui"
    "github.com/therecipe/qt/widgets"
    "io/ioutil"
    "os"
    "time"
)


var (
    DIRECT_UP Direction = Direction{
        x: -1,
        y: 0,
    }
    DIRECT_DOWN Direction = Direction{
        x: 1,
        y: 0,
    }
    DIRECT_LEFT Direction = Direction{
        x: 0,
        y: -1,
    }
    DIRECT_RIGHT Direction = Direction{
        x: 0,
        y: 1,
    }
)
const (
    TYPE_SNAKE_HEAD int8 = 1
    TYPE_SNAKE_BODY int8 = 2
    TYPE_SNAKE_TAIL int8 = 3
    TYPE_FOOD int8 = 4
    
    
)

type config struct {
    BgImage string `json:"bgImage"`
    AppIcon string `json:"appIcon"`
    
}
func (c *config) init() {
    data,err := ioutil.ReadFile("config.json")
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
    fmt.Println("====>>>config<<<====")
    fmt.Println(c)
    fmt.Println("====>>>>>><<<<<<====")
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

type Food struct {
    x int
    y int
    t int
}


var (
    rows int = 30
    column int = 20
)
type Direction struct {
    x int
    y int
}

var (
    app *widgets.QMainWindow
    gameMap [][]int8
    isGameOver bool
    configuration config
    bgPixMap *gui.QPixmap
    bgPalette *gui.QPalette
    brush *gui.QBrush
    layout *widgets.QWidget
    snake Snake
    wHeight int
    wWidth int
    currentDirection Direction
    food []Food
)
func resetSnake(){
    currentDirection = DIRECT_LEFT
    body := SnakeNode{
        x:     rows / 2,
        y:     column / 2,
        color: 0,
        image: "",
        next:  nil,
    }
    snake = Snake{
        head: &body,
        tail: &body,
    }
}
// TODO
func resetFood() {
    food = append(food, Food{
        x: 3,
        y: 1,
        t: 0,
    })
}
func refreshSnake(){
    for !isGameOver {
        // 新节点
        newNodeX := snake.head.x + currentDirection.x
        newNodeY := snake.head.y + currentDirection.y
        // TODO 吃食物
        // 判断撞墙
        if newNodeX < 0 || newNodeX >= rows || newNodeY < 0 || newNodeY >= column{
            isGameOver = true
            break
        }
        snake.AppendNodeToHeader(&SnakeNode{
            x:     newNodeX,
            y:     newNodeY,
            color: 0,
            image: "",
            next:  nil,
        })
        snake.DeleteTail()
        time.Sleep(1 * time.Second)
    }
}
func layoutPaintHandle(event *gui.QPaintEvent){
    painter := gui.NewQPainter2(layout)
    // 刷新,绘制食物,绘制小蛇
    for i := 0;i < rows;i ++{
        for j := 0; j < column; j ++{
            gameMap[i][j] = 0
        }
    }
    // 放置食物
    for i := 0;i < len(food);i ++ {
        gameMap[food[i].x][food[i].y] = TYPE_FOOD
    }
    // 放置小蛇身体
    for cur := snake.tail;cur != nil; cur = cur.next{
        gameMap[cur.x][cur.y] = TYPE_SNAKE_BODY
    }
    gameMap[(snake.tail.x + rows) % rows][(snake.tail.y + column) % column] = TYPE_SNAKE_TAIL
    gameMap[(snake.head.x + rows) % rows][(snake.head.y + column) % column] = TYPE_SNAKE_HEAD
    
    for i := 0;i < rows;i ++{
        for j := 0;j < column;j ++{
            blockType := gameMap[i][j]
            switch blockType {
            case TYPE_SNAKE_TAIL:
                painter.FillRect7(j*40+2,i*40+2,37,37,core.Qt__lightGray)
                break
            case TYPE_SNAKE_HEAD:
                painter.FillRect7(j*40+2,i*40+2,37,37,core.Qt__black)
                break
            case TYPE_SNAKE_BODY:
                painter.FillRect7(j*40+2,i*40+2,37,37,core.Qt__yellow)
                break
            case TYPE_FOOD:
                painter.FillRect7(j*40+2,i*40+2,37,37,core.Qt__red)
                break
            }
        }
    }
    painter.End()
    event.Accept()
}

func pressBtnHandle(event *gui.QKeyEvent){
    key := event.Key()
    switch key {
    case int(core.Qt__Key_Up):
        fmt.Println("up")
        currentDirection = DIRECT_UP
        break
    case int(core.Qt__Key_Down):
        fmt.Println("down")
        currentDirection = DIRECT_DOWN
        break
    case int(core.Qt__Key_Left):
        fmt.Println("left")
        currentDirection = DIRECT_LEFT
        break
    case int(core.Qt__Key_Right):
        fmt.Println("right")
        currentDirection = DIRECT_RIGHT
        break
    default:
        fmt.Printf("unknown press key %d\n",key)
    }
    event.Accept()
}
func InitUI() *widgets.QMainWindow {
    configuration = config{}
    configuration.init()
    gameMap = make([][]int8,rows)
    for i:=0;i<rows;i++{
        gameMap[i] = make([]int8,column)
    }
    resetSnake()
    resetFood()
    app = widgets.NewQMainWindow(nil, 0)
    app.SetWindowTitle("贪吃蛇")
    layout = widgets.NewQWidget(app, core.Qt__Widget)
    app.SetCentralWidget(layout)
    app.SetWindowIcon(gui.NewQIcon5(configuration.AppIcon))
    app.SetMinimumSize2(rows * 40, column * 40)
    
    var gameActions []*widgets.QAction
    // 开始
    startAction := widgets.NewQAction2("开始", app)
    startAction.ConnectTriggered(func(checked bool) {
        isGameOver = false
        resetSnake()
        go func() {
            for !isGameOver {
                // 29fps => 1秒刷新30次
                layout.Repaint()
                time.Sleep(33 * time.Microsecond)
            }
        }()
        go refreshSnake()
    })
    gameActions = append(gameActions, startAction)
    
    // 退出
    exitAction := widgets.NewQAction3(gui.NewQIcon5("exit.png"),"&exit", app)
    exitAction.SetShortcut(gui.NewQKeySequence2("Ctrl+Q", gui.QKeySequence__NativeText))
    exitAction.SetToolTip("退出游戏")
    exitAction.ConnectTriggered(func(checked bool) {
        app.Close()
    })
    gameActions = append(gameActions, exitAction)
    
    menuBar := app.MenuBar()
    gameMenu := menuBar.AddMenu2("游戏")
    gameMenu.AddActions(gameActions)
    
    bgPixMap = gui.NewQPixmap3(configuration.BgImage, "", core.Qt__AutoColor)
    bgPalette = gui.NewQPalette()
    brush = gui.NewQBrush()
    layout.ConnectPaintEvent(layoutPaintHandle)
    app.ConnectKeyPressEvent(pressBtnHandle)
    setStyle(app)
    return app
}
func setStyle(app *widgets.QMainWindow) {
    //
    app.ConnectPaintEvent(func(event *gui.QPaintEvent) {
        // 重设图片宽高以适应应用大小
        if !(wHeight == app.Height() && wWidth == app.Width()) {
            bgPixMapTmp := bgPixMap.Scaled2(app.Width(),app.Height(),core.Qt__IgnoreAspectRatio,core.Qt__SmoothTransformation)
            brush.SetTexture(bgPixMapTmp)
            bgPalette.SetBrush(gui.QPalette__Background,brush)
            app.SetPalette(bgPalette)
            wHeight = app.Height()
            wWidth = app.Width()
        }
        event.Accept()
    })
}

func main() {
    widgets.NewQApplication(len(os.Args),os.Args)
    app := InitUI()
    app.Show()
    widgets.QApplication_Exec()
    
}
