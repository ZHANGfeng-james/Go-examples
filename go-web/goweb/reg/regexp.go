package reg

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func IsIP(ip string) {
	if m, _ := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$", ip); !m {
		fmt.Println("false")
	}
	fmt.Println("true")
}

func handleValues() {
	resp, err := http.Get("http://www.baidu.com")
	if err != nil {
		fmt.Println("http get error.")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("http read error")
		return
	}

	src := string(body)

	// 将 HTML 标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	// 去除 STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	// 去除 SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	// 去除所有尖括号内的 HTML 代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")

	// 去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

	fmt.Println(strings.TrimSpace(src))
}

func regexpTest() {
	a := "I am learning Go language"

	re, _ := regexp.Compile("[a-z]{2,4}")

	// 查找符合正则的第一个
	one := re.Find([]byte(a))
	// Find: am
	fmt.Println("Find:", string(one))

	// 查找符合正则的所有 slice, n 小于 0 表示返回全部符合的字符串，不然就是返回指定的长度
	all := re.FindAll([]byte(a), -1)
	// FindAll:[am lear ning lang uage].
	fmt.Printf("FindAll:%s.\n", all)

	// 查找符合条件的 index 位置, 开始位置和结束位置
	index := re.FindIndex([]byte(a))
	fmt.Println("FindIndex", index)

	// 查找符合条件的所有的 index 位置，n 同上
	allindex := re.FindAllIndex([]byte(a), -1)
	fmt.Println("FindAllIndex", allindex)

	re2, _ := regexp.Compile("am(.*)lang(.*)")

	// 查找 Submatch, 返回数组，第一个元素是匹配的全部元素，第二个元素是第一个 () 里面的，第三个是第二个 () 里面的
	// 下面的输出第一个元素是 "am learning Go language"
	// 第二个元素是 " learning Go "，注意包含空格的输出
	// 第三个元素是 "uage"
	submatch := re2.FindSubmatch([]byte(a))
	fmt.Printf("FindSubmatch:%s\n", submatch)
	for _, v := range submatch {
		fmt.Println(string(v))
	}

	// 定义和上面的 FindIndex 一样
	submatchindex := re2.FindSubmatchIndex([]byte(a))
	fmt.Println(submatchindex)

	// FindAllSubmatch, 查找所有符合条件的子匹配
	submatchall := re2.FindAllSubmatch([]byte(a), -1)
	fmt.Printf("submatchall:%s\n", submatchall)

	// FindAllSubmatchIndex, 查找所有字匹配的 index
	submatchallindex := re2.FindAllSubmatchIndex([]byte(a), -1)
	fmt.Println(submatchallindex)
}

func regExpExpand() {
	src := []byte(`
        call hello alice
        hello bob
        call hello eve
    `)
	pat := regexp.MustCompile(`(?m)(call)\s+(?P<cmd>\w+)\s+(?P<arg>.+)\s*$`)
	res := []byte{}

	slice := pat.FindAllSubmatchIndex(src, -1)
	fmt.Println(slice)
	for _, s := range slice {
		res = pat.Expand(res, []byte("$cmd('$arg')\n"), src, s)
	}
	fmt.Println(string(res))
}
