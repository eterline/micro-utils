// Copyright (c) 2025 EterLine (Andrew)
// This file is part of My-Go-Project.
// Licensed under the MIT License. See the LICENSE file for details.


package seeip

import "fmt"

type GeoSeer interface {
	Latitude() float64
	Longitude() float64
}

func GoogleMap(s GeoSeer) string {
	return fmt.Sprintf(
		"https://www.google.com/maps?q=%f,%f",
		s.Latitude(),
		s.Longitude(),
	)
}

func OpenStreetMap(s GeoSeer) string {
	return fmt.Sprintf(
		"https://www.openstreetmap.org/?mlat=%f&mlon=%f#map=16/%f/%f",
		s.Latitude(),
		s.Longitude(),
		s.Latitude(),
		s.Longitude(),
	)
}

func YandexMap(s GeoSeer) string {
	return fmt.Sprintf(
		"https://yandex.ru/maps/?ll=%f%%2C%f&z=16",
		s.Longitude(),
		s.Latitude(),
	)
}

func TwoGIS(s GeoSeer) string {
	return fmt.Sprintf(
		"https://2gis.ru/?query=%f,%f",
		s.Latitude(),
		s.Longitude(),
	)
}
