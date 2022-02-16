package stm

import (
	"bytes"
	"errors"
	"time"
)

// NewBuilderIndexfile returns the created the BuilderIndexfile's pointer
func NewBuilderIndexfile(opts *Options, loc *Location) *BuilderIndexfile {
	return &BuilderIndexfile{opts: opts, loc: loc}
}

// BuilderIndexfile provides implementation for the Builder interface.
type BuilderIndexfile struct {
	opts    *Options
	loc     *Location
	content []byte
	linkcnt int
}

// Add method joins old bytes with creates bytes by it calls from Sitemap.Finalize method.
func (b *BuilderIndexfile) Add(link interface{}) BuilderError {
	bldr, ok := link.(*BuilderFile)
	if !ok {
		return &builderFileError{error: errors.New("link is not a *BuilderFile"), full: true}
	}

	smu := NewSitemapIndexURL(b.opts, URL{
		{"loc", bldr.loc.URL()},
		{"lastmod", time.Now()},
	})

	b.content = append(b.content, smu.XML()...)
	b.linkcnt++

	return nil
}

// Content and BuilderFile.Content are almost the same behavior.
func (b *BuilderIndexfile) Content() []byte {
	return b.content
}

// XMLContent and BuilderFile.XMLContent share almost the same behavior.
func (b *BuilderIndexfile) XMLContent() []byte {
	c := bytes.Join(bytes.Fields(IndexXMLHeader), []byte(" "))
	c = append(append(c, b.Content()...), IndexXMLFooter...)

	return c
}

// Write and Builderfile.Write are almost the same behavior.
func (b *BuilderIndexfile) Write() {
	c := b.XMLContent()

	b.loc.Write(c, b.linkcnt)
}
