package datasource

import (
	"testing"
	"time"
)

func TestDataSource_Query(t *testing.T) {
	ds := New()

	begin := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 1, 1, 0, 0, 0, time.UTC)

	cursor, err := ds.Query(begin, end)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	count := 0
	for {
		dp, ok := cursor.Next()
		if !ok {
			break
		}
		count++
		if dp.Timestamp.Before(begin) || dp.Timestamp.After(end) {
			t.Fatalf("DataPoint out of range: %v", dp.Timestamp)
		}
	}

	if count == 0 {
		t.Fatalf("No data points returned")
	}
}
