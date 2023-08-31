package main

import "fmt"

// PersonShow 抽象方法
type PersonShow interface {
	show()
}

// Person 对象
type Person struct {
	name string
}

func (p *Person) show() {
	fmt.Println(p.name)
}

// PersonDecorator 抽象类
type PersonDecorator struct {
	personShow PersonShow
}

func (pd *PersonDecorator) show() {
	if pd.personShow != nil {
		pd.personShow.show()
	}
}

// 设置装饰器
func (pd *PersonDecorator) SetPersonShow(ps PersonShow) {
	pd.personShow = ps
}

// PersonDecorator1 装饰器1
type PersonDecorator1 struct {
	personDecorator PersonDecorator
}

func (pd *PersonDecorator1) show() {
	pd.personDecorator.show()       // 基本方法
	fmt.Println("PersonDecorator1") // 装饰方法
}

// PersonDecorator2 装饰器2
type PersonDecorator2 struct {
	personDecorator PersonDecorator
}

func (pd *PersonDecorator2) show() {
	pd.personDecorator.show()       // 基本方法
	fmt.Println("PersonDecorator2") // 装饰方法
}

func DecorativePattern() {
	person := &Person{name: "lp"}
	pd1 := &PersonDecorator1{}
	pd1.personDecorator.SetPersonShow(person)

	pd2 := &PersonDecorator2{}
	pd2.personDecorator.SetPersonShow(pd1)

	pd2.show()
}
