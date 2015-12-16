<?php

use Sunra\PhpSimple\HtmlDomParser;

class TestTemplateService extends TestCase
{
    private $service;

    public function __construct($name = null, array $data = [], $dataName = '')
    {
        parent::__construct($name, $data, $dataName);

        $app = $this->createApplication(); 
        
        $this->service = $app->make(newsletters\Services\TemplateService::class); 
    }

    /**
     * @param $html
     * @param $url
     * @param $expected
     *
     * @dataProvider providerAppendImageToDom
     */ 
    public function testAppendImageToDom($html, $url, $expected)
    { 
        $dom = HtmlDomParser::str_get_html($html);

        $newDom = $this->callMethod($this->service, 'appendImg', [$dom, $url]);

        $this->assertEquals($newDom->outertext, $expected);
    }

    /**
     * @param $html
     * @param array $tags
     * @param $expected
     *
     * @dataProvider providerReplaceTagsInTemplate
     */
    public function testReplaceTagsInTemplate($html, $tags, $expected)
    {
        $dom = HtmlDomParser::str_get_html($html);

        $newDom = $this->callMethod($this->service, 'replaceTagsInTemplate', [$dom, $tags]);

        $this->assertEquals($newDom->outertext, $expected);
    }

    public function providerAppendImageToDom()
    { 
        return [
            ['<html><body></body></html>', 'http://test.com', "<html><body><img src=\"http://test.com\"/>\n</body></html>"],
            ['<html><body><div>something here</div></body></html>', 'http://test.com', "<html><body><div>something here</div><img src=\"http://test.com\"/>\n</body></html>"],
            ['<body></body>', 'http://test.com', "<body></body><img src=\"http://test.com\"/>\n"],
            ['<html></html>', 'http://test.com', "<html></html><img src=\"http://test.com\"/>\n"],
            ['<p>something</p>', 'http://test.com', "<p>something</p><img src=\"http://test.com\"/>\n"],
        ];
    }

    public function providerReplaceTagsInTemplate()
    {
        return [
            [
                '<html><body><p>Dear *|Name|*</p><div>*|cOUnTry|*</div></body></html>',
                ['/\*\|Name\|\*/i' => 'John Doe', '/\*\|country\|\*/i' => 'mk'],
                '<html><body><p>Dear John Doe</p><div>mk</div></body></html>'
            ]
        ];
    }
}
