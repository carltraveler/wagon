for((i = 0; i < 1000; i++));do go run main.go a.out; (($? != 0)) && break; done
