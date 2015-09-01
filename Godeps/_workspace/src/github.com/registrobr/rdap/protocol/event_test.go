package protocol

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	eventDate := Date(2015, 9, 1, 10, 20, 30, 0, time.UTC)

	expected := "2015-09-01T10:20:30Z"
	if out := eventDate.Format(time.RFC3339); out != expected {
		t.Errorf("Expected “%s” and got “%s”", expected, out)
	}
}

func TestNewEventDate(t *testing.T) {
	eventDate := NewEventDate(time.Date(2015, 9, 1, 10, 20, 30, 0, time.UTC))

	expected := "2015-09-01T10:20:30Z"
	if out := eventDate.Format(time.RFC3339); out != expected {
		t.Errorf("Expected “%s” and got “%s”", expected, out)
	}
}

func TestEventDateUnmarshalJSON(t *testing.T) {
	data := []struct {
		description   string
		data          []byte
		expected      EventDate
		expectedError error
	}{
		{
			description: "it should import a RFC3339 correctly",
			data:        []byte(`"2015-08-31T16:12:52Z"`),
			expected: EventDate{
				Time: time.Date(2015, 8, 31, 16, 12, 52, 0, time.UTC),
			},
		},
		{
			description: "it should import a partial RFC3339 correctly (date only)",
			data:        []byte(`"2015-08-31"`),
			expected: EventDate{
				Time: time.Date(2015, 8, 31, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			description: "it should import a partial RFC3339 correctly (no timezone)",
			data:        []byte(`"2015-08-31T16:12:52"`),
			expected: EventDate{
				Time: time.Date(2015, 8, 31, 16, 12, 52, 0, time.UTC),
			},
		},
		{
			description:   "it should fail for an invalid RFC3339",
			data:          []byte(`"31/8/2015"`),
			expectedError: fmt.Errorf(`parsing time ""31/8/2015"" as ""2006-01-02T15:04:05"": cannot parse "/2015"" as "2006"`),
		},
	}

	for i, item := range data {
		var eventDate EventDate
		err := eventDate.UnmarshalJSON(item.data)

		if item.expectedError != nil {
			if fmt.Sprintf("%v", item.expectedError) != fmt.Sprintf("%v", err) {
				t.Errorf("[%d] %s: expected error “%s”, got “%s”", i, item.description, item.expectedError, err)
			}

		} else if err != nil {
			t.Errorf("[%d] %s: unexpected error “%s”", i, item.description, err)

		} else {
			if !reflect.DeepEqual(item.expected, eventDate) {
				t.Errorf("[%d] %s: unexpected event date returned. Expected “%#v” and got “%#v”", i, item.description, item.expected.String(), eventDate.String())
			}
		}
	}
}

func TestEventDateUnmarshalText(t *testing.T) {
	data := []struct {
		description   string
		data          []byte
		expected      EventDate
		expectedError error
	}{
		{
			description: "it should import a RFC3339 correctly",
			data:        []byte("2015-08-31T16:12:52Z"),
			expected: EventDate{
				Time: time.Date(2015, 8, 31, 16, 12, 52, 0, time.UTC),
			},
		},
		{
			description: "it should import a partial RFC3339 correctly (date only)",
			data:        []byte("2015-08-31"),
			expected: EventDate{
				Time: time.Date(2015, 8, 31, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			description: "it should import a partial RFC3339 correctly (no timezone)",
			data:        []byte("2015-08-31T16:12:52"),
			expected: EventDate{
				Time: time.Date(2015, 8, 31, 16, 12, 52, 0, time.UTC),
			},
		},
		{
			description:   "it should fail for an invalid RFC3339",
			data:          []byte("31/8/2015"),
			expectedError: fmt.Errorf(`parsing time "31/8/2015" as "2006-01-02T15:04:05": cannot parse "/2015" as "2006"`),
		},
	}

	for i, item := range data {
		var eventDate EventDate
		err := eventDate.UnmarshalText(item.data)

		if item.expectedError != nil {
			if fmt.Sprintf("%v", item.expectedError) != fmt.Sprintf("%v", err) {
				t.Errorf("[%d] %s: expected error “%s”, got “%s”", i, item.description, item.expectedError, err)
			}

		} else if err != nil {
			t.Errorf("[%d] %s: unexpected error “%s”", i, item.description, err)

		} else {
			if !reflect.DeepEqual(item.expected, eventDate) {
				t.Errorf("[%d] %s: unexpected event date returned. Expected “%#v” and got “%#v”", i, item.description, item.expected.String(), eventDate.String())
			}
		}
	}
}
