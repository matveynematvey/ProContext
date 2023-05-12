package lib

import (
	"ProContext/valute"
	"bufio"
	"fmt"
	"os"
	"time"
)

func FormDate(date time.Time) string {
	return fmt.Sprintf("%02d/%02d/%d", date.Day(), int(date.Month()), date.Year())
}

func PrintResult(start, end time.Time, daysLimit int64, valutes valute.Valutes) {
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fprintf(out, "Exchange rate in the period from %s to %s\n", FormDate(end), FormDate(start))
	for key, val := range valutes {
		fmt.Fprintf(out, "Valute: %s\nMax: %f, Date: %s\nMin: %f, Date: %s\nMean: %f\n\n",
			key, val.Max.Value, val.Max.Date, val.Min.Value, val.Min.Date, val.GetMean(daysLimit))
	}
}
