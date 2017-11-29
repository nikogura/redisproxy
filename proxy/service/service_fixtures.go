package service

import "log"

func testCapacity() int {
	return 3
}

func testMaxAge() int {
	return 5
}

func testTimeout() int {
	return 5
}

func testRedisAddr() string {
	return ""
}

func testSlice() []string {
	stringSlice := make([]string, 0)
	stringSlice = append(stringSlice, "foo")
	stringSlice = append(stringSlice, "bar")
	stringSlice = append(stringSlice, "baz")

	return stringSlice
}

func testFoo() string {
	return "foo"
}

func testBar() string {
	return "bar"
}

func testWip() string {
	return "wip"
}

func testZoz() string {
	return "zoz"
}

func testTen() int {
	return 10
}

func testCacheData() (info map[string]interface{}) {
	info = make(map[string]interface{})

	info["foo"] = testFoo()
	info["bar"] = testBar()
	info["wip"] = testWip()
	info["zoz"] = testZoz()
	info["ten"] = testTen()
	info["list"] = testSlice()

	return info
}

func integTestFetchFunc(key string, redisAddr string) (value interface{}, err error) {
	log.Printf("in integTestFetchFunc")
	data := testCacheData()

	if elem, ok := data[key]; ok {
		return elem, nil
	}

	return value, err
}
