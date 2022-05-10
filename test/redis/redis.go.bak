package main

import (
	"github.com/gomodule/redigo/redis"
)

var RedisPool *redis.Pool

func main() {
	conn, err := redis.Dial("tcp", "192.168.10.33:6379")
	if err != nil {
		panic(err)
	}
	//读取
	//redis.Strings：返回多个
	//redis.String：返回一个
	//redis.int:返回统计的数字

	//Do方法 发送一条命令到redis服务器端执行并返回结果
	//字符串(String)
	/*conn.Do("SET", "name", "hi,baby!")
	name, err := conn.Do("GET", "name")
	fmt.Println(name) //[104 101 108 108 111]
	if err != nil {
		fmt.Println("redis get error:", err)
	} else {
		fmt.Printf("Get name: %s \n", name) //Get name: hello
	}*/

	//列表(List)
	/*_, err = conn.Do("LPUSH", "taglist", "ele1", "ele2", "ele3")
	if err != nil {
		fmt.Println("redis mset error:", err)
	}
	for i := 0; i < 3; i++ {
		res, err := redis.String(conn.Do("LPOP", "taglist"))
		if err != nil {
			fmt.Println("redis POP error:", err)
		} else {
			res_type := reflect.TypeOf(res)
			fmt.Printf("res type : %s \n", res_type)
			fmt.Printf("res  : %s \n", res)
		}
	}*/

	//散列(Hash)
	/*_, err = conn.Do("HSET", "student", "name", "wd", "age", 22)
	if err != nil {
		fmt.Println("redis mset error:", err)
	}
	res, err := redis.Int64(conn.Do("HGET", "student", "age"))
	if err != nil {
		fmt.Println("redis HGET error:", err)
	} else {
		res_type := reflect.TypeOf(res)
		fmt.Printf("res type : %s \n", res_type)
		fmt.Printf("res  : %d \n", res)
	}*/

	//集合(Set)
	/*_, err = conn.Do("MSET", "name", "wd", "sa", "da", "age", 22, 23, 43)
	if err != nil {
		fmt.Println("redis mset error:", err)
	}
	res, err := redis.Strings(conn.Do("MGET", "name", "age", "sa"))
	if err != nil {
		fmt.Println("redis get error:", err)
	} else {
		res_type := reflect.TypeOf(res)
		fmt.Printf("res type : %s \n", res_type)
		fmt.Printf("MGET name: %s \n", res)
		fmt.Println(len(res))
	}
	//结果：
	//res type : []string
	//MGET name: [wd 22]
	//2*/

	//事务
	/*conn.Send("MULTI")
	conn.Send("INCR", "foo")
	conn.Send("INCR", "bar")
	r, err := conn.Do("EXEC")
	fmt.Println(r)*/

	//有序集合(Sorted Set)
	//关闭链接
	defer conn.Close()
}

type tag struct {
	Name      string `json:"name"`
	State     int    `json:"state"`
	CreatedBy int    `json:"createdBy"`
}
