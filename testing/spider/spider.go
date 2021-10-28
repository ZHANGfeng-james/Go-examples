package spider

type Spider interface {
	GetBody() string
}

func CreateGoVersionSpider() Spider {
	return nil
}
