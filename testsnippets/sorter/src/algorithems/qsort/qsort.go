package qsort

func QSort(values []int) {
	qsort(values, 0, len(values)-1)
}

func qsort(values []int, low, high int) {
	if low < high {
		mid := qsortpartition(values, low, high)
		qsort(values, low, mid-1)
		qsort(values, mid+1, high)
	}
}

func qsortpartition(values []int, low, high int) int {
	piv := values[low]
    for low < high {
        for low < high && values[high] >= piv {
        	high--	
    	}
        values[low] = values[high]

        for low < high && values[low] <= piv {
        	low	++
    	}
        values[high] = values[low]            
    }
    values[low] = piv
    return low
}