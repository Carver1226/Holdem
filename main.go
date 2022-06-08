package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"time"
)

func main() {
	//matches, err := ReadJson("match_result.json")
	//matches, err := ReadJson("match_result.json")
	//if err != nil {
	//	panic(err)
	//}
	//
	//for _, match := range matches.Matches {
	//	result := Compare(match.Alice, match.Bob)
	//	if result != match.Result {
	//		fmt.Printf("match_result: %+v\t%v\n", match, result)
	//	}
	//}
	//matches, err = ReadJson("seven_cards_with_ghost.json")
	//if err != nil {
	//	panic(err)
	//}
	//for _, match := range matches.Matches {
	//	result := Compare(match.Alice, match.Bob)
	//	if result != match.Result {
	//		fmt.Printf("seven_cards_with_ghost: %+v\t%v\n", match, result)
	//	}
	//}
	start := time.Now() // 获取当前时间
	matches, err := ReadJson("seven_cards_with_ghost.result.json")
	if err != nil {
		panic(err)
	}
	for _, match := range matches.Matches {
		result := Compare(match.Alice, match.Bob)
		if result != match.Result {
			fmt.Printf("seven_cards_with_ghost.result: %+v\t%v\n", match, result)
		}
	}
	elapsed := time.Since(start)
	fmt.Println("耗时：", elapsed)
}

type Match struct {
	Alice  string `json:"alice"`
	Bob    string `json:"bob"`
	Result int    `json:"result"'`
}

type Matches struct {
	Matches []Match `json:"matches"'`
}

type MyArr struct {
	Arr []int
}

func (arr MyArr) Len() int {
	return len(arr.Arr)
}

func (arr MyArr) Less(i, j int) bool {
	return arr.Arr[i]%100 < arr.Arr[j]%100
}

func (arr MyArr) Swap(i, j int) {
	arr.Arr[i], arr.Arr[j] = arr.Arr[j], arr.Arr[i]
}

//比较两幅牌的大小,a大返回1, b大返回2, 相等返回0
func Compare(a, b string) int {
	aCards := ConvertCards(a)
	bCards := ConvertCards(b)
	aHand, aWeight := Judge(aCards)
	bHand, bWeight := Judge(bCards)
	if aHand > bHand {
		return 1
	}
	if aHand < bHand {
		return 2
	}
	if aHand == bHand {
		if aWeight > bWeight {
			return 1
		}
		if aWeight < bWeight {
			return 2
		}
	}
	return 0
}

//判断牌型 返回牌型和对应的权值
func Judge(cards []int) (int, int) {
	if i, b := isStraightFlush(cards); b {
		return 9, i
	}
	if i, b := isShi(cards); b {
		return 8, i
	}
	if i, b := isGourd(cards); b {
		return 7, i
	}
	if i, b := isSame(cards); b {
		return 6, i
	}
	if i, b := isStraight(cards); b {
		return 5, i
	}
	if i, b := isSan(cards); b {
		return 4, i
	}
	if i, b := isTwoPairs(cards); b {
		return 3, i
	}
	if i, b := isOnePairs(cards); b {
		return 2, i
	}
	return 1, Leaflet(cards)
}

//将字符串转换成卡牌
func ConvertCards(s string) []int {
	cards := make([]int, len(s)/2)
	bytes := []byte(s)
	//下标
	var i int
	for j := range cards {
		switch string(bytes[i]) {
		case "T":
			cards[j] = 10
		case "J":
			cards[j] = 11
		case "Q":
			cards[j] = 12
		case "K":
			cards[j] = 13
		case "A":
			cards[j] = 14
		case "X":
			cards[j] = 99
		default:
			cards[j], _ = strconv.Atoi(string(bytes[i]))
		}
		i++
		switch string(bytes[i]) {
		case "s":
			cards[j] += 100
		case "h":
			cards[j] += 200
		case "d":
			cards[j] += 300
		case "c":
			cards[j] += 400
		}
		i++
	}
	sort.Sort(MyArr{Arr: cards})
	return cards
}

//判断是否连续 返回最大值和是否连续
func isStraight(cards []int) (int, bool) {
	//连续的头号牌
	num := 0
	//连续的最大值
	max := 0
	//赖子个数
	reply := 0
	//添加A
	cards = append([]int{0}, cards...)
	//判断是否有A
	if cards[len(cards)-1]%100 == 14 || cards[len(cards)-2]%100 == 14 {
		cards[0] = 1
	}
	//判断是否有赖子
	if cards[len(cards)-1] == 99 {
		reply++
	}
	//没有赖子的情况下 通过双指针判断是否连续
	if reply == 0 {
		for i, j, repeat := 1, 0, 0; i < len(cards); i++ {
			for j < i && ((cards[i]%100-cards[i-1]%100) != 1 && cards[i]%100 != cards[i-1]%100) {
				repeat = 0
				j++
			}
			//判断前一个是否为重复数
			if cards[i]%100 == cards[i-1]%100 {
				repeat++
			}
			if (i - j + 1 - repeat) > max {
				max = i - j + 1 - repeat
				num = cards[i]
			}
		}
	}
	//若存在一个赖子,则使用队列来进行判断
	if reply > 0 {
		//生成队列
		l := list.New()
		//判断赖子是否被使用过
		var isUse bool
		for i, j := 0, 0; j < len(cards); {
			/*for e := l.Front(); e != nil; e = e.Next() {
				fmt.Printf("%v ", e.Value)
			}
			fmt.Println()*/
			if l.Len() == 0 {
				if cards[j] == 0 {
					j++
					continue
				}
				l.PushBack(cards[j])
				j++
				isUse = false
				continue
			}
			if l.Len() >= max || (l.Len() >= 5 && num%100 < l.Back().Value.(int)%100) {
				max = l.Len()
				num = l.Back().Value.(int)
			}
			lastCard := l.Back().Value.(int)
			if cards[j]%100-lastCard%100 == 1 {
				l.PushBack(cards[j])
				j++
				continue
			}
			if cards[j]%100 == lastCard%100 {
				j++
				continue
			}
			if !isUse {
				l.PushBack(lastCard + 1)
				isUse = true
				continue
			}
			var next *list.Element
			for e := l.Front(); e != nil; e = next {
				next = e.Next()
				l.Remove(e)
			}
			i++
			j = i
		}
	}
	if max >= 5 {
		if num%100 == 14 {
			return num%100, true
		}
		return num%100, true
	}
	return num%100, false
}

//判断是否同花 返回权值和是否同花
func isSame(cards []int) (int, bool) {
	arr := make([]int, 5)
	reply := 0
	for _, v := range cards {
		if v == 99 {
			reply++
			continue
		}
		arr[v/100]++
	}
	for i, v := range arr {
		if v + reply >= 5 {
			result := 0
			num := 0 //记录同花的数量
			//将同花最大的五张牌列出来
			for j := len(cards) - 1; j >= 0; j-- {
				if cards[j]/100 == i || cards[j] == 99{
					result = result*100 + cards[j]%100
					num++
				}
				if num == 5 {
					break
				}
 			}
 			return result, true
		}
	}

	return 0, false
}

//判断是否同花顺 返回最大值和是否同花顺
func isStraightFlush(cards []int) (int, bool) {
	//连续的头号牌
	num := 0
	//连续的最大值
	max := 0
	//赖子个数
	reply := 0
	//添加A
	cards = append([]int{0}, cards...)
	//判断是否有A
	if cards[len(cards)-1]%100 == 14{
		cards[0] = 1 + cards[len(cards)-1]/100*100
	}
	if cards[len(cards)-2]%100 == 14{
		cards[0] = 1 + cards[len(cards)-2]/100*100
	}
	if cards[len(cards)-1] == 99 {
		reply++
	}
	//没有赖子的情况下, 通过队列判断是否连续且同花
	if reply == 0 {
		//生成队列
		l := list.New()
		for i := 0; i < len(cards); i++{
			if l.Len() == 0 {
				if cards[i] == 0 {
					continue
				}
				l.PushBack(cards[i])
				continue
			}
			lastCard := l.Back().Value.(int)
			if cards[i]%100 - lastCard%100 == 1 && cards[i]/100 == lastCard/100 {
				l.PushBack(cards[i])
				if i == len(cards) - 1 {
					if l.Len() >= max {
						max = l.Len()
						num = l.Back().Value.(int)
					}
				}
				continue
			}
			if lastCard%100 == 1 {
				l.Remove(l.Front())
				l.PushBack(cards[i])
				continue
			}
			//判断是否是重复
			if cards[i]%100 == lastCard%100 || cards[i]%100 - lastCard%100 == 1 {
				continue
			}
			if l.Len() >= max {
				max = l.Len()
				num = l.Back().Value.(int)
			}
			if l.Len() == 1 {
				l.Remove(l.Front())
				l.PushBack(cards[i])
				continue
			}
			var next *list.Element
			for e := l.Front(); e != nil; e = next {
				next = e.Next()
				l.Remove(e)
			}
		}
	}
	//存在赖子的情况下,通过队列进行判断
	if reply > 0 {
		//生成队列
		l := list.New()
		//判断赖子是否被使用过
		var isUse bool
		for i, j := 0, 0; j < len(cards); {
			/*for e := l.Front(); e != nil; e = e.Next() {
				fmt.Printf("%v ", e.Value)
			}
			fmt.Println()*/
			if l.Len() == 0 {
				if cards[j] == 0 {
					j++
					continue
				}
				l.PushBack(cards[j])
				j++
				isUse = false
				continue
			}
			if l.Len() >= max || (l.Len() >= 5 && num%100 < l.Back().Value.(int)%100) {
				max = l.Len()
				num = l.Back().Value.(int)
			}
			lastCard := l.Back().Value.(int)
			if cards[j]%100 - lastCard%100 == 1 && cards[j]/100 == lastCard/100 {
				l.PushBack(cards[j])
				j++
				continue
			}
			if cards[j]%100 == lastCard%100 || cards[j]%100 - lastCard%100 == 1 {
				j++
				continue
			}
			if !isUse {
				l.PushBack(lastCard + 1)
				isUse = true
				continue
			}
			var next *list.Element
			for e := l.Front(); e != nil; e = next {
				next = e.Next()
				l.Remove(e)
			}
			i++
			j = i
		}
	}
	if max >= 5 {
		if num%100 == 14 {
			return num%100, true
		}
		return num%100, true
	}
	return num%100, false
}

//判断是否葫芦 返回权值
func isGourd(cards []int) (int, bool) {
	var a, b int
	arr := make([]int, 15)
	reply := 0
	for _, v := range cards {
		if v == 99 {
			reply++
			continue
		}
		arr[v%100]++
	}

	//判断没有赖子的情况下
	if reply == 0 {
		for i := 14; i > 0; i-- {
			if arr[i] >= 3 && a < i {
				a = i
				continue
			}
			if arr[i] >= 2 && b < i {
				b = i
			}
		}
	}

	//判断有一个赖子的情况下
	if reply > 0 {
		var a1, b1, a2, b2 int
		//第一种情况 将赖子用于3的情况下
		for i := 14; i > 0; i-- {
			if arr[i]+reply >= 3 && a1 < i {
				a1 = i
				continue
			}
			if arr[i] >= 2 && b1 < i {
				b1 = i
			}
		}
		//第二种情况 将赖子用于2的情况下
		for i := 14; i > 0; i-- {
			if arr[i] >= 3 && a2 < arr[i] {
				a2 = i
				continue
			}
			if arr[i] + reply >= 2 && b2 < arr[i] {
				b2 = i
			}
		}

		//判断赖子用于第一种情况还是第二种情况
		if a1 != 0 && b1 != 0 {
			a = a1
			b = b1
		}
		if a2 != 0 && b2 != 0 {
			if (a*100 + b) < (a2*100 + b2) {
				a = a2
				b = b2
			}
		}
	}

	if a != 0 && b != 0 {
		return a * 100 + b, true
	}
	return 0, false
}

//判断是否四条 返回权值
func isShi(cards []int) (int, bool) {
	var a, b int
	arr := make([]int, 15)
	reply := 0
	for _, v := range cards {
		if v == 99 {
			reply++
			continue
		}
		arr[v%100]++
	}
	//没有赖子的情况下
	if reply == 0 {
		for i := 14; i > 0; i-- {
			if arr[i] >= 4 && a < i {
				a = i
				continue
			}
			if arr[i] >= 1 && b < i {
				b = i
			}
		}
	}
	//有一个赖子的情况下
	if reply > 0 {
		var a1, b1, a2, b2 int
		//第一种情况 赖子用于四条上
		for i := 14; i > 0; i-- {
			if arr[i] + reply >= 4 && a1 < i {
				a1 = i
				continue
			}
			if arr[i] >= 1 && b1 < i {
				b1 = i
			}
		}
		//第二种情况 赖子用于单牌上
		b2 = 14
		for i := 14; i > 0; i-- {
			if arr[i] >= 4 && a2 < i {
				a2 = i
				break
			}
		}
		//判断赖子用于哪种情况
		if a1 != 0 && b1 != 0 {
			a = a1
			b = b1
		}
		if a2 != 0 && b2 != 0 {
			if (a*100 + b) < (a2*100 + b2) {
				a = a2
				b = b2
			}
		}
	}
	if a != 0 && b != 0 {
		return a*100 + b, true
	}
	return 0, false
}

//判断是否三条 返回权值
func isSan(cards []int) (int, bool) {
	var a, b, c int
	arr := make([]int, 15)
	reply := 0
	for _, v := range cards {
		if v == 99 {
			reply++
			continue
		}
		arr[v%100]++
	}
	//判断没赖子的情况下
	if reply == 0 {
		for i := 14; i > 0; i-- {
			if arr[i] >= 3 && a < i {
				a = i
				continue
			}
			if arr[i] >= 1 && b < i {
				b = i
				continue
			}
			if arr[i] >= 1 && c < i {
				c = i
			}
		}
	}
	//有一张赖子的情况
	if reply > 0 {
		var a1, b1, c1, a2, b2, c2 int
		//第一种情况,赖子用于对子上
		for i := 14; i > 0; i-- {
			if arr[i] + reply >= 3 && a1 < i {
				a1 = i
				continue
			}
			if arr[i] >= 1 && b1 < i {
				b1 = i
				continue
			}
			if arr[i] >= 1 && c1 < i {
				c1 = i
			}
		}
		//第二种情况,赖子用于单牌上
		b2 = 14
		for i := 14; i > 0; i-- {
			if arr[i] >= 3 && a2 < i {
				a2 = i
				continue
			}
			if arr[i] >= 1 && c2 < i {
				c2 = i
				continue
			}
		}
		//判断赖子用于哪种情况
		if a1 != 0 && b1 != 0 && c1 != 0 {
			a = a1
			b = b1
			c = c1
		}
		if a2 != 0 && b2 != 0 && c2 != 0 {
			if (a*10000 + b*100 + c) < (a2*10000 + b2*100 + c2) {
				a = a2
				b = b2
				c = c2
			}
		}
	}
	if a != 0 && b != 0 && c != 0 {
		return a*10000 + b*100 + c, true
	}
	return 0, false
}

//判断是否两对 返回权值
func isTwoPairs(cards []int) (int, bool) {
	var a, b, c int
	arr := make([]int, 15)
	reply := 0
	for _, v := range cards {
		if v == 99 {
			reply++
			continue
		}
		arr[v%100]++
	}
	//没赖子的情况下
	if reply == 0 {
		for i := 14; i > 0; i-- {
			if arr[i] >= 2 && a < i {
				a = i
				continue
			}
			if arr[i] >= 2 && b < i {
				b = i
				continue
			}
			if arr[i] >= 1 && c < i {
				c = i
			}
		}
	}
	//有一个赖子的情况下
	if reply > 0 {
		var a1, b1, c1, a2, b2, c2, a3, b3, c3 int
		//第一种情况 赖子用于第一对
		for i := 14; i > 0; i-- {
			if arr[i] + reply >= 2 && a1 < i {
				a1 = i
				continue
			}
			if arr[i] >= 2 && b1 < i {
				b1 = i
				continue
			}
			if arr[i] >= 1 && c1 < i {
				c1 = i
			}
		}
		//第二种情况 赖子用于第二对
		for i := 14; i > 0; i-- {
			if arr[i] >= 2 && a2 < i {
				a2 = i
				continue
			}
			if arr[i] + reply >= 2 && b2 < i {
				b2 = i
				continue
			}
			if arr[i] >= 1 && c2 < i {
				c2 = i
			}
		}
		//第三种情况 赖子用于单牌
		c3 = 14
		for i := 14; i > 0; i-- {
			if arr[i] >= 2 && a3 < i {
				a3 = i
				continue
			}
			if arr[i] >= 2 && b3 < i {
				b3 = i
				continue
			}
		}
		//判断赖子应该用于哪种情况
		if a1 != 0 && b1 != 0 && c1 != 0 {
			a = a1
			b = b1
			c = c1
		}
		if a2 != 0 && b2 != 0 && c2 != 0 {
			if (a*10000 + b*100 + c) < (a2*10000 + b2*100 + c2) {
				a = a2
				b = b2
				c = c2
			}
		}
		if a3 != 0 && b3 != 0 && c3 != 0 {
			if (a*10000 + b*100 + c) < (a3*10000 + b3*100 + c3) {
				a = a3
				b = b3
				c = c3
			}
		}
	}
	if a != 0 && b != 0 && c != 0 {
		return a*10000 + b*100 + c, true
	}
	return 0, false
}

//判断是否一对 返回权值
func isOnePairs(cards []int) (int, bool) {
	var a, b, c, d int
	arr := make([]int, 15)
	reply := 0
	for _, v := range cards {
		if v == 99 {
			reply++
			continue
		}
		arr[v%100]++
	}
	//没赖子的情况
	if reply == 0 {
		for i := 14; i > 0; i-- {
			if arr[i] >= 2 && a < i {
				a = i
				continue
			}
			if arr[i] >= 1 && b < i {
				b = i
				continue
			}
			if arr[i] >= 1 && c < i {
				c = i
				continue
			}
			if arr[i] >= 1 && d < i {
				d = i
			}
		}
	}
	//有一张赖子的情况
	if reply > 0 {
		var a1, b1, c1, d1, a2, b2, c2, d2 int
		//第一种情况,赖子用于对子上
		for i := 14; i > 0; i-- {
			if arr[i] + reply >= 2 && a1 < i {
				a1 = i
				continue
			}
			if arr[i] >= 1 && b1 < i {
				b1 = i
				continue
			}
			if arr[i] >= 1 && c1 < i {
				c1 = i
				continue
			}
			if arr[i] >= 1 && d1 < i {
				d1 = i
			}
		}
		//第二种情况,赖子用于单牌上
		b2 = 14
		for i := 14; i > 0; i-- {
			if arr[i] >= 2 && a2 < i {
				a2 = i
				continue
			}
			if arr[i] >= 1 && c2 < i {
				c2 = i
				continue
			}
			if arr[i] >= 1 && d2 < i {
				d2 = i
			}
		}
		//判断赖子用于哪种情况
		if a1 != 0 && b1 != 0 && c1 != 0 && d1 != 0 {
			a = a1
			b = b1
			c = c1
			d = d1
		}
		if a2 != 0 && b2 != 0 && c2 != 0 && d2 != 0 {
			if (a*1000000 + b*10000 + c*100 + d) < (a2*1000000 + b2*10000 + c2*100 + d2) {
				a = a2
				b = b2
				c = c2
				d = d2
			}
		}
	}
	if a != 0 && b != 0 && c != 0 {
		return a*1000000 + b*10000 + c*100 + d, true
	}
	return 0, false
}

//单张的情况下,返回权值
func Leaflet(cards []int) int {
	result := 0
	//计算权值
	for i := len(cards) - 1; i > len(cards) - 6; i-- {
		result = result*100 + cards[i]%100
	}
	return result
}

func ReadJson(fileName string) (*Matches, error) {
	byte, err := ioutil.ReadFile("./" + fileName)
	var matches *Matches
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(byte, &matches)
	if err != nil {
		return nil, err
	}
	return matches, nil
}