package desgin

import "fmt"

type Simplefactory interface {
	Say(name string)
	Set(name string)
}

type LiShi struct {
	sName string
	sType string
}

type YuWen struct {
	yName     string
	yFileName string
}

func NewLiShi() Simplefactory {
	return &LiShi{}
}
func NewYuWen() Simplefactory {
	return &YuWen{}
}
func (p *LiShi) Say(name string) {
	fmt.Println("LiShi=>", p.sName, name)

}
func (p *LiShi) Set(name string) {
	p.sName = name
}

func (b *YuWen) Say(name string) {
	fmt.Println("YuWen=>", b.yName, name)

}
func (b *YuWen) Set(name string) {
	b.yName = name
}
