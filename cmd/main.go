package main

import (
	"BolshiGoLang/internal/pkg/server"
	"BolshiGoLang/internal/pkg/storage"
)

func main() {
	r, err := storage.NewStorage()
	if err != nil {
		panic(err)
	}
	s := server.NewServer(":8090", &r)

	s.Start()
}

/*
func main() {
	s, err := fileutils.FileRead()
	if err != nil {
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("RWrite the number of operation 1 - Set key, value; 2 - Get key; 2 times Enter to exit")
		var i string
		i, _ = reader.ReadString('\n')
		i = i[:len(i)-1]
		if i == "1" {
			input, err := reader.ReadString('\n')
			if err != nil || len(input) == 1 {
				break
			}
			input = input[:len(input)-1]
			var key, value string
			_, err = fmt.Sscanf(input, "%s %s", &key, &value)
			//fmt.Printf("%s, %s, %s\n", input, key, value)
			if err != nil {
				continue
			}
			s.Set(key, value)
			fmt.Println("OK")
		} else {
			input, err := reader.ReadString('\n')
			if err != nil || len(input) == 1 {
				break
			}
			input = input[:len(input)-1]
			var key string
			_, err = fmt.Sscanf(input, "%s", &key)
			fmt.Println(s.Get(key))
		}
	}

	err = fileutils.FileWrite(s)
	if err != nil {
		return
	}
}
*/
