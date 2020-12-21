package wordFilter

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

/*
敏感词数据结构：
{  王：{
            isEnd: false
            八：{
                    isEnd：false
                    蛋：{
                              isEnd：true
                       }
                 }
       }
}
*/

//敏感词结构体
type SensitiveMap struct {
	sensitiveNode map[string]interface{} //敏感词节点
	isEnd         bool                   //是否为敏感词最后一个字
}

var s *SensitiveMap

//初始化
func InitMap() *SensitiveMap {
	if s == nil {
		dictionary := IninFilePath()      //获取敏感词文件路径
		s = InitDictionary(s, dictionary) //构造词典树
	}
	return s
}

//初始化，敏感词库文件
func IninFilePath() (dictionaryPaths []string) {
	currentPath, _ := filepath.Abs(filepath.Dir(os.Args[0])) //程序运行路径绝对地址
	//查找文件
	fmt.Println("敏感词库地址：", currentPath)
	//自定义配置文件夹下的所有敏感词文件
	projectPaths, _ := filepath.Glob(currentPath + "/*")
	dictionaryPaths = append(dictionaryPaths, projectPaths...)
	return
}

//初始化敏感词词典结构体
func initSensitiveMap() *SensitiveMap {
	return &SensitiveMap{
		sensitiveNode: make(map[string]interface{}),
		isEnd:         false,
	}
}

//多文件词典,读取所有敏感词
func readDictionarys(paths []string) (dicts []string) {
	for _, row := range paths { //遍历配置文件路径
		dict := readDictionary(row)
		dicts = append(dicts, dict...) //将所有敏感词追加到同一切片
	}
	return
}

//读取词典文件,返回该文件下的敏感词汇切片
func readDictionary(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		return []string{}
	}
	defer file.Close()
	str, err := ioutil.ReadAll(file)
	dictionary := strings.Fields(string(str)) //按空格分隔文本
	return dictionary
}

//初始化敏感词词典，根据DFA算法构建trie
func InitDictionary(s *SensitiveMap, dictionaryPath []string) *SensitiveMap {
	s = initSensitiveMap() //构造结构体
	//多词典，读取所有敏感词
	dictionary := readDictionarys(dictionaryPath)
	//循环敏感词
	for _, words := range dictionary {
		sMapTmp := s
		w := []rune(words)
		wordsLength := len(w)
		for i := 0; i < wordsLength; i++ { //遍历敏感词得到组成字
			t := string(w[i])
			isEnd := false
			//如果是词汇最后一个字，则为true
			if i == (wordsLength - 1) {
				isEnd = true
			}

			//创造词典树
			func(tx string) {
				if _, ok := sMapTmp.sensitiveNode[tx]; !ok { //字在该层索引中不存在，创建新的层级
					sMapTemp := new(SensitiveMap)
					sMapTemp.sensitiveNode = make(map[string]interface{})
					sMapTemp.isEnd = isEnd

					sMapTmp.sensitiveNode[tx] = sMapTemp //新的节点
				}
				sMapTmp = sMapTmp.sensitiveNode[tx].(*SensitiveMap) //进入下一层级
			}(t)

		}
	}
	return s
}

//检测是否含有敏感词，返回第一个敏感词
func (s *SensitiveMap) CheckSensitive(text string) (string, bool) {
	content := []rune(text)
	contentLength := len(content)
	result := false //敏感状态
	ta := ""        //敏感词
	//遍历文本
	for index := range content {
		sMapTmp := s
		target := ""
		in := index
		for { //遍历词典进行比对
			wo := string(content[in])
			target += wo
			if _, ok := sMapTmp.sensitiveNode[wo]; ok { //存在该节点
				if sMapTmp.sensitiveNode[wo].(*SensitiveMap).isEnd { //节点为结束状态
					result = true //敏感词存在
					break
				}
				if in == contentLength-1 { //遍历到最后一个字
					break
				}
				sMapTmp = sMapTmp.sensitiveNode[wo].(*SensitiveMap) //进入到下一层级
				in++
			} else { //节点不存在
				break
			}
		}
		if result {
			ta = target
			break
		}
	}
	return ta, result
}

//检测所有敏感词
type Target struct {
	Indexes []int //敏感词在文本中的位置
	Len     int   //敏感词长度
}

func (s *SensitiveMap) FindAllSensitive(text string) map[string]*Target {
	content := []rune(text)
	contentLength := len(content)
	result := false                //敏感状态
	ta := make(map[string]*Target) //敏感词
	//遍历文本
	for index := range content {
		sMapTmp := s
		target := ""
		in := index
		result = false

		for { //遍历词典进行比对
			wo := string(content[in])
			target += wo
			if _, ok := sMapTmp.sensitiveNode[wo]; ok { //存在该节点
				if sMapTmp.sensitiveNode[wo].(*SensitiveMap).isEnd { //节点为结束状态
					result = true //敏感词存在
					break
				}
				if in == contentLength-1 { //遍历到最后一个字
					break
				}
				sMapTmp = sMapTmp.sensitiveNode[wo].(*SensitiveMap) //进入到下一层级
				in++
			} else { //节点不存在
				break
			}
		}
		if result {
			if _, targetInta := ta[target]; targetInta { //敏感词已检测为存在
				ta[target].Indexes = append(ta[target].Indexes, index) //添加下标位置
			} else {
				ta[target] = &Target{ //新建
					Indexes: []int{index},
					Len:     len([]rune(target)),
				}
			}
		}
	}
	return ta
}

//封装 ,返回待检测内容中所有敏感词
func CheckoutWordAll(text string) (strs []string, isBool bool) {
	isBool = false
	// 应该放在项目启动的时候进行初始化,避免多线程调用导致初始化问题
	s := InitMap()
	if s == nil {
		return
	}
	target := s.FindAllSensitive(text)
	if len(target) > 0 {
		isBool = true
		for key := range target {
			strs = append(strs, key)
		}

	}
	return
}

//检查，返回内容中，第一条敏感词
func CheckWord(text string) (str string, isBool bool) {
	isBool = false
	s := InitMap()
	str, isBool = s.CheckSensitive(text)
	return
}
