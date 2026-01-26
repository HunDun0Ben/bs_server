package origin

const IFNAMSIZ = 16

type myIfreq struct {
	Name [IFNAMSIZ]byte
	Flag uint16
	_    [22]byte
}
