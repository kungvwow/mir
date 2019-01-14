package mir

type site struct {
	group    Group `mir:"v1"`
	index    Get   `mir:"/index/"`
	articles Get   `mir:"//{subdomain}.domain.com/articles/{category}/{id:[0-9]+}?{filter}&{pages}#GetArticles"`
}

func (h *site) Index() string {
	return "Index"
}

func (h *site) GetArticles() string {
	return "GetArticles"
}

type handlerFunc func() string

type simpleEngine struct {
	pathHandler map[string]handlerFunc
}

func (e *simpleEngine) Register(entries []*TagMir) error {
	for _, entry := range entries {
		for _, field := range entry.Fields {
			e.pathHandler[field.Path] = field.Handler.(func() string)
		}
	}
	return nil
}

type blog struct {
	group    Group `mir:"v1#Log&Count"`
	index    Get   `mir:"/index/"`
	articles Get   `mir:"//{subdomain}.domain.com/articles/{category}/{id:[0-9]+}?{filter}&{pages}#GetArticles"`
}

func (b *blog) Index() string {
	return "Index"
}

func (b *blog) GetArticles() string {
	return "GetArticles"
}

func (b *blog) Log() string {
	return "Log"
}

func (b *blog) Count() string {
	return "Count"
}
