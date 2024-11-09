//go:build !solution

package hotelbusiness

import (
	"sort"
)

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func ComputeLoad(guests []Guest) []Load {

	// create the answer slice
	ans := []Load{}

	// go through all the guests and put check in and check out dates into separate maps to keep track of them
	// also keep all the dates in two slices, so they can be sorted later

	checkInDates := make(map[int]int)

	checkOutDates := make(map[int]int)

	for _, guest := range guests {
		checkInDates[guest.CheckInDate] += 1
		checkOutDates[guest.CheckOutDate] += 1
	}

	// sort the dates into one slice
	datesSorted := make([]int, 0)

	for date := range checkInDates {
		datesSorted = append(datesSorted, date)
	}

	for date := range checkOutDates {
		datesSorted = append(datesSorted, date)
	}

	sort.Ints(datesSorted)

	// go through the sorted dates and fill in the ans with hproperly filled loads for each date

	guestCount := 0
	currentDate := 0

	for i := 0; i < len(datesSorted); i++ {
		currentDate = datesSorted[i]
		guestCount += checkInDates[currentDate]
		guestCount -= checkOutDates[currentDate]

		// only add load to the ans if the change is non-zero
		if checkInDates[currentDate]-checkOutDates[currentDate] != 0 {
			ans = append(ans, Load{currentDate, guestCount})
		}
	}

	// return computed answer
	return ans
}

