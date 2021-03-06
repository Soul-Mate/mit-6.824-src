package mapreduce

import (
	"encoding/json"
	"hash/fnv"
	"io/ioutil"
	"os"
)

// doMap manages one map task: it reads one of the input files
// (inFile), calls the user-defined map function (mapF) for that file's
// contents, and partitions the output into nReduce intermediate files.
func doMap(
	jobName string,    // the name of the MapReduce job
	mapTaskNumber int, // which map task this is
	inFile string,
	nReduce int, // the number of reduce task that will be run ("R" in the paper)
	mapF func(file string, contents string) []KeyValue,
) {
	//
	// You will need to write this function.
	// 你需要编写此函数
	//
	// The intermediate output of a map task is stored as multiple
	// files, one per destination reduce task. The file name includes
	// both the map task number and the reduce task number. Use the
	// filename generated by reduceName(jobName, mapTaskNumber, r) as
	// the intermediate file for reduce task r. Call ihash() (see below)
	// on each key, mod nReduce, to pick r for a key/value pair.
	//
	// map 任务的中间输出存储为多个文件, 文件名应该包含 map 任务号和 reduce 任务号.
	// 使用reduceName(jobName, mapTaskNumber, r) 函数为中间输出生成文件名. 每个 key
	// 调用 ihash() (见下文) mod nReduce,  去选择 r 作为键值对写入到中间文件
	//
	//
	// mapF() is the map function provided by the application. The first
	// argument should be the input file name, though the map function
	// typically ignores it. The second argument should be the entire
	// input file contents. mapF() returns a slice containing the
	// key/value pairs for  reduce; see common.go for the definition of
	// KeyValue.
	//
	// mapF() 是应用程序提供的函数. 第一个参数是输入文件名称, 虽然通常会忽略它.
	// 第二个参数是文件内容. mapF() 返回片用于 reduce 的包含key/value 的 slice.
	// 参阅 common.go 以了解 keyValue 的定义.
	//
	// Look at Go's ioutil and os packages for functions to read
	// and write files.
	//
	// 查看 Go 的 ioutil 和 os 包以获取读写文件的功能.
	//
	// Coming up with a scheme for how to format the key/value pairs on
	// disk can be tricky, especially when taking into account that both
	// keys and values could contain newlines, quotes, and any other
	// character you can think of.
	//
	// 提出一个如何格式化磁盘上 key/value 对的方案可能很棘手, 特别考虑到 key 和 value
	// 两者都可能包含换行符, 引号, 和其他你能想到的字符.
	//
	// One format often used for serializing data to a byte stream that the
	// other end can correctly reconstruct is JSON. You are not required to
	// use JSON, but as the output of the reduce tasks *must* be JSON,
	// familiarizing yourself with it here may prove useful. You can write
	// out a data structure as a JSON string to a file using the commented
	// code below. The corresponding decoding functions can be found in
	// common_reduce.go.
	//
	//
	//   enc := json.NewEncoder(file)
	//   for _, kv := ... {
	//     err := enc.Encode(&kv)
	//
	// Remember to close the file after you have written all the values!
	//

	// 读取输入的文件
	b, err := ioutil.ReadFile(inFile)
	if err != nil {
		panic(err)
	}

	// 通过预定义的 map 函数生成中间的 key/value pair
	kvs := mapF(inFile, string(b))

	// 对每个 key 通过分区函数写入到不同的 r 文件
	for _, kv := range kvs {

		// ihash 可以帮助我们对 key 进行 hash
		reduceTaskNumber := ihash(kv.Key) % nReduce

		// reduceName 可以帮助我们创建中间文件名称
		fileName := reduceName(jobName, mapTaskNumber, reduceTaskNumber)

		// 注意不要用 os.Create, 因为我们需要对文件进行追加写
		w, _ := os.OpenFile(fileName, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0644)

		// 将中间 key/value pair 写入到 r 文件
		enc := json.NewEncoder(w)

		if err := enc.Encode(&kv); err != nil {
			panic(err)
		}

		// 记得关闭文件
		if err := w.Close(); err != nil {
			panic(err)
		}
	}

}

func ihash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32() & 0x7fffffff)
}
