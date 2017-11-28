package redisproxy

import "log"

// you could argue that these are silly, but I really don't want to hard code test values in the tests.

// The test is itself code, and we should always separate code from data.  Plus, reading this you are introduced to the type of testdata I'm expecting.

// In fact, This fixtures file becomes a spec of sorts for the expected inputs.

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

func unitTestFetchFunc(key string) (value interface{}, err error) {
	log.Printf("In test fetch func\n")
	data := testCacheData()

	if elem, ok := data[key]; ok {
		log.Printf("Fetcher returning value: %s", elem)
		return elem, nil
	}

	log.Printf("Fetcher Returning nil\n")
	return value, err
}
