# Gohjong: Mahjong Library written in Go

A library for mahjong

### Usage

```go
h = "1112224577799m"
sw, _ := gohjong.ShowWaiting(h)
for _, s := range sw {
    fmt.Println(h, "is waiting", s) //=> 1112224577799m is waiting 3m-6m
}
```

### License

MIT
