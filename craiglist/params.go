package craiglist

import "strconv"

type MinBedsParam struct {
	MinBeds uint
}

func (b MinBedsParam) String() string {
	return ("min_bedrooms=" + strconv.Itoa(int(b.MinBeds)))
}

type MaxBedsParam struct {
	MaxBeds uint
}

func (b MaxBedsParam) String() string {
	return ("max_bedrooms=" + strconv.Itoa(int(b.MaxBeds)))
}

type HasPictureParam struct {
	HasPics bool
}

func (h HasPictureParam) String() string {
	if h.HasPics {
		return "hasPic=1"
	}
	return "hasPic=0"
}

type MaxBathParam struct {
	MaxBaths uint
}

func (b MaxBathParam) String() string {
	return ("max_bathrooms=" + strconv.Itoa(int(b.MaxBaths)))
}

type MinBathParam struct {
	MinBaths uint
}

func (b MinBathParam) String() string {
	return ("min_bathrooms=" + strconv.Itoa(int(b.MinBaths)))
}

type MinPriceParam struct {
	MinPrice uint
}

func (b MinPriceParam) String() string {
	return ("min_price=" + strconv.Itoa(int(b.MinPrice)))
}

type MaxPriceParam struct {
	MaxPrice uint
}

func (b MaxPriceParam) String() string {
	return ("max_price=" + strconv.Itoa(int(b.MaxPrice)))
}
