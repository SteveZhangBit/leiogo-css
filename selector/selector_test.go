package selector

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	if doc, err := Parse(
		`<div id="post">
			<div class="cls links">
				<a href="http://www.baidu.com">baidu</a>
				<a href="http://www.zhihu.com">zhihu</a>
			</div>
		</div>
		<div>
			<div class="cls2 links">
				<img src="/images/1.jpg">
				<img src="/images/2.png">
			</div>
		</div>`,
	); err != nil {
		t.Error(err.Error())
	} else {
		el := doc.Find(`#post a`)
		fmt.Println(el.Attr("href"))
	}
}
