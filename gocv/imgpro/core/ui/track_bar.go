package ui

import "gocv.io/x/gocv"

type Trackbar struct {
	bar      *gocv.Trackbar
	name     string
	value    int
	maxValue int
	onChange func(int)
}

func (w *ProcessingWindow) AddTrackbar(name string, maxValue int, onChange func(int)) {
	value := 0
	bar := w.window.CreateTrackbarWithValue(name, &value, maxValue)

	t := &Trackbar{
		bar:      bar,
		name:     name,
		maxValue: maxValue,
		value:    value,
		onChange: onChange,
	}

	w.trackbars = append(w.trackbars, t)
}

func (w *ProcessingWindow) updateTrackbars() bool {
	updated := false
	for _, t := range w.trackbars {
		newValue := t.bar.GetPos()
		if newValue != t.value {
			t.value = newValue
			if t.onChange != nil {
				t.onChange(newValue)
			}
			updated = true
		}
	}
	return updated
}
