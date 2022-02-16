package stm

import (
	"log"
	"runtime"
	"sync"
)

// NewSitemap returns the created the Sitemap's pointer
func NewSitemap(maxProc int) *Sitemap {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	if maxProc < 1 || maxProc > runtime.NumCPU() {
		maxProc = runtime.NumCPU()
	}
	//log.Printf("Max processors %d\n", maxProc)
	runtime.GOMAXPROCS(maxProc)

	sm := &Sitemap{
		opts: NewOptions(),
	}
	return sm
}

// Sitemap provides interface for create sitemap xml file and that has convenient interface.
// And also needs to use first this struct if it wants to use this package.
type Sitemap struct {
	opts  *Options
	bldr  []Builder
	bldrs Builder
}

// SetDefaultHost is your website's host name
func (sm *Sitemap) SetDefaultHost(host string) {
	sm.opts.SetDefaultHost(host)
}

// SetSitemapsHost is the remote host where your sitemaps will be hosted
func (sm *Sitemap) SetSitemapsHost(host string) {
	sm.opts.SetSitemapsHost(host)
}

// SetSitemapsPath sets this to a directory/path if you don't
// want to upload to the root of your `SitemapsHost`
func (sm *Sitemap) SetSitemapsPath(path string) {
	sm.opts.SetSitemapsPath(path)
}

// SetPublicPath is the directory to write sitemaps to locally
func (sm *Sitemap) SetPublicPath(path string) {
	sm.opts.SetPublicPath(path)
}

// SetAdapter can switch output file storage.
// We have S3Adapter and FileAdapter (default: FileAdapter)
func (sm *Sitemap) SetAdapter(adp Adapter) {
	sm.opts.SetAdapter(adp)
}

// SetVerbose can switch verbose output to console.
func (sm *Sitemap) SetVerbose(verbose bool) {
	sm.opts.SetVerbose(verbose)
}

// SetCompress can switch compress for the output file.
func (sm *Sitemap) SetCompress(compress bool) {
	sm.opts.SetCompress(compress)
}

// SetPretty option allows pretty formating to the output files.
func (sm *Sitemap) SetPretty(pretty bool) {
	sm.opts.SetPretty(pretty)
}

// SetFilename can apply any name in this method if you wants to change output file name
func (sm *Sitemap) SetFilename(filename string) {
	sm.opts.SetFilename(filename)
}

// Create method must be that calls first this method in that before call to Add method on this struct.
func (sm *Sitemap) Create() *Sitemap {
	sm.bldrs = NewBuilderIndexfile(sm.opts, sm.opts.IndexLocation())
	return sm
}

// NewBuilderFile returns a BuilderFile with options copied from the Sitemap.
func (sm *Sitemap) NewBuilderFile(filename string) *BuilderFile {
	opts := sm.opts.Clone()
	opts.SetFilename(filename)
	bldr := NewBuilderFile(opts, opts.Location())
	sm.bldr = append(sm.bldr, bldr)
	sm.bldrs.Add(bldr)
	return bldr
}

// XMLContent returns the XML content of the sitemap
func (sm *Sitemap) XMLContent() []byte {
	if sm != nil && len(sm.bldr) > 0 {
		return sm.bldr[0].XMLContent()
	}
	return nil
}

func bldrWrite(bldr Builder, wg *sync.WaitGroup) {
	defer wg.Done()
	bldr.Write()
}

// Finalize writes sitemap and index files if it had some
// specific condition in BuilderFile struct.
func (sm *Sitemap) Finalize() *Sitemap {
	count := 1
	if sm.bldr != nil {
		count += len(sm.bldr)
	}

	wg := &sync.WaitGroup{}
	wg.Add(count)

	for _, bldr := range sm.bldr {
		go bldrWrite(bldr, wg)
	}

	go bldrWrite(sm.bldrs, wg)

	wg.Wait()

	sm.bldr = nil
	return sm
}

// PingSearchEngines requests some ping server.
// It also has that includes PingSearchEngines function.
func (sm *Sitemap) PingSearchEngines(urls ...string) {
	PingSearchEngines(sm.opts, urls...)
}
