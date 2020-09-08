package goutils

//拆分4份
//fmt.Println("println:",SplitArray(arr1,4))
//数组平分
func SplitArray(arr []interface{}, num int) [][]interface{} {
	max := int(len(arr))
	if max < num {
		return nil
	}
	var segmens = make([][]interface{}, 0)
	quantity := max / num
	end := int(0)
	for i := int(1); i <= num; i++ {
		qu := i * quantity
		if i != num {
			segmens = append(segmens, arr[i-1+end:qu])
		} else {
			segmens = append(segmens, arr[i-1+end:])
		}
		end = qu - i
	}
	return segmens
}
