package main

import "fmt"
//import "sort"
type ByteSlice []byte
type Sequence []int


func main() {
	myslice1 := ByteSlice {'a', 'b', 'c'}
	myslice2 := make([]byte, 5)
	myslice3 := make([]byte, 5, 10)

	fmt.Printf("info of slice1: len = %d, cap = %d", len(myslice1), cap(myslice1))
	fmt.Printf("info of slice2: len = %d, cap = %d", len(myslice2), cap(myslice2))
	fmt.Printf("info of slice3: len = %d, cap = %d", len(myslice3), cap(myslice3))

	outPutSlice(myslice1)
	outPutSlice(myslice2)
	outPutSlice(myslice3)
	//testSlice := make(ByteSlice, 5)

	testSlice := myslice1.myAppend(myslice2)

	outPutSlice(testSlice)

	mysq1 := Sequence{3, 5, 6, 7,8,8,9,9,6,543}
	fmt.Println("content of sq1 = ", mysq1)
}


func outPutSlice(slice []byte) {
	fmt.Printf("info of slice: len = %d, cap = %d\n", len(slice), cap(slice))
	for i, v := range slice{
		fmt.Println("slice [" , i , "] = ", v)
	}
}

func (slice ByteSlice) myAppend(data []byte) []byte {
	//slice := *p
	l := len(slice)

	if l + len(data) > cap(slice) {
		newSlice := make([]byte, l+len(data)*2)
		copy(newSlice, slice)
		slice = newSlice
	}

	slice = slice[0:l + len(data)]

	for i, v := range data {
		slice[l+i] = v
	}

	//*p = slice
	return slice

}

func init() {

	fmt.Println("call from init");
}

func init() {

	fmt.Println("call from init2");
}

// Methods required by sort.Interface.
// sort.Interface 所需的方法。
func (s Sequence) Len() int {
    return len(s)
}
func (s Sequence) Less(i, j int) bool {
    return s[i] < s[j]
}
func (s Sequence) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

// Method for printing - sorts the elements before printing.
// 用于打印的方法 - 在打印前对元素进行排序。
// func (s Sequence) String() string {
//     sort.Sort(s)
//     str := "["
//     for i, elem := range s {
//         if i > 0 {
//             str += " "
//         }
//         str += fmt.Sprint(elem)
//     }
//     return str + "]"
// }