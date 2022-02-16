package stm

import (
	"reflect"
	"strings"
	"testing"

	"github.com/clbanning/mxj"
)

func TestSitemapGenerator(t *testing.T) {
	buf := BufferAdapter{}

	sm := NewSitemap(0)
	sm.SetPretty(true)
	sm.SetVerbose(false)
	sm.SetSitemapsPath("sitemap")
	sm.SetAdapter(&buf)

	sm.Create()
	bldr := sm.NewBuilderFile("test")
	for i := 1; i <= 1; i++ {
		bldr.Add(URL{{"loc", "home"}, {"changefreq", "always"}, {"mobile", true}, {"lastmod", "2018-10-28T17:56:02+09:00"}})
		bldr.Add(URL{{"loc", "readme"}, {"lastmod", "2018-10-28T17:56:02+09:00"}})
		bldr.Add(URL{{"loc", "aboutme"}, {"priority", 0.1}, {"lastmod", "2018-10-28T17:56:02+09:00"}})
	}
	sm.Finalize()

	buffers := buf.Bytes()

	//data := buffers[len(buffers)-1]
	expects := [][]byte{
		[]byte(`
	<?xml version="1.0" encoding="UTF-8"?>
	<sitemapindex xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/siteindex.xsd" xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	  <sitemap>
		<loc>http://www.example.com/sitemaps/test.xml.gz</loc>
		<lastmod>2018-10-28T17:37:21+09:00</lastmod>
	  </sitemap>
	</sitemapindex>`),
		[]byte(`
	<?xml version="1.0" encoding="UTF-8"?> <urlset xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd" xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:image="http://www.google.com/schemas/sitemap-image/1.1" xmlns:video="http://www.google.com/schemas/sitemap-video/1.1" xmlns:geo="http://www.google.com/geo/schemas/sitemap/1.0" xmlns:news="http://www.google.com/schemas/sitemap-news/0.9" xmlns:mobile="http://www.google.com/schemas/sitemap-mobile/1.0" xmlns:pagemap="http://www.google.com/schemas/sitemap-pagemap/1.0" xmlns:xhtml="http://www.w3.org/1999/xhtml" ><url>
	  <loc>http://www.example.com/home</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>always</changefreq>
	  <priority>0.5</priority>
	  <mobile:mobile/>
	</url>
	<url>
	  <loc>http://www.example.com/readme</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>weekly</changefreq>
	  <priority>0.5</priority>
	</url>
	<url>
	  <loc>http://www.example.com/aboutme</loc>
	  <lastmod>2018-10-28T17:56:02+09:00</lastmod>
	  <changefreq>weekly</changefreq>
	  <priority>0.1</priority>
	</url>
	</urlset>`),
	}

	for _, data := range buffers {
		mdata, _ := mxj.NewMapXml(data)
		if !strings.Contains(string(data), "sitemapindex") {
			continue
		}
		mdata.Remove("sitemapindex.sitemap.lastmod")
		matchFound := false
		for _, expect := range expects {
			if !strings.Contains(string(expect), "sitemapindex") {
				continue
			}
			mexpect, _ := mxj.NewMapXml(expect)
			mexpect.Remove("sitemapindex.sitemap.lastmod")
			if reflect.DeepEqual(mdata, mexpect) {
				matchFound = true
			}
			if !matchFound {
				t.Errorf(`Failed to match: Expect: %s, Received: %s`, string(expect), string(data))
			}
		}

	}

}
