#encoding:utf-8

# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: http://doc.scrapy.org/en/latest/topics/item-pipeline.html

from scrapy.spiders import CrawlSpider,Rule
from scrapy.linkextractors import LinkExtractor
from scrapy.selector import Selector
from zhihu.items import ZhihuItem
from scrapy.http import Request,FormRequest
import sys
import string
import json
sys.stdout=open('output.txt','w')

add=0
class ZhihuuserSpider(CrawlSpider):
     name="zhihuuserspider"
     allowed_domains=["www.zhihu.com"]
     start_urls=[
          'https://www.zhihu.com/people/wen-yi-yang-81/following'
          ]
     rules=(
          Rule(LinkExtractor(allow = ('/people/(\d+)$', )),callback='parse_page', follow = True), 
          )
     headers = {
     "Accept": "*/*",
     "Accept-Encoding": "gzip,deflate",
     "Accept-Language": "en-US,en;q=0.8,zh-TW;q=0.6,zh;q=0.4",
     "Connection": "keep-alive",
     "Content-Type":" application/x-www-form-urlencoded; charset=UTF-8",
     "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.2125.111 Safari/537.36",
     "Referer": "http://www.zhihu.com/"
    }
     print 'start'
     
     def add_cookie(self,request):
         print 'add_cookie'
         return request
     def start_requests(self):  
         print 'start_requests'
         return [Request("https://www.zhihu.com/people/wen-yi-yang-81/following",headers = self.headers)]
         #return [Request("https://www.zhihu.com/#signin",headers = self.headers, meta = {'cookiejar' : 1},callback = self.post_login)]
     
     def post_login(self,response):
         print 'prelogin'
         sel=Selector(response)
         xsrf=sel.xpath('//input[@name="_xsrf"]/@value').extract()[0]
         print "xsrf="+xsrf 
         return [FormRequest(url = 'http://www.zhihu.com/login/email',  
                            meta = {'cookiejar' : response.meta['cookiejar']},
                            headers = self.headers,
                            formdata = {
                            '_xsrf': xsrf,
                            'email': 'hlm123yy@163.com',
                            'password': 'hlm123tt',
                            'captcha':'abcd'
                            },
                            callback = self.check_login,
                            dont_filter = True
                            )]
     def check_login(self,response):
         print 'check_login'
         print json.loads(response.body)
         if json.loads(response.body)['r'] == 0:
            print "login success"
            yield Request('http://www.zhihu.com', 
                           meta = {'cookiejar' : response.meta['cookiejar']},
                           headers=self.headers, 
                           callback=self.parse_page,
                           dont_filter=True,    
                         )     
     def after_login(self,response):
         print "after_login"
         for url in self.start_urls:
             print "parse url:"+url
             yield self.make_requests_from_url(url)
     def parse_page(self, response):
         global add
         print add
         add+=1
         print response
         #print "parse_page"
         #problem = Selector(response)
         #item = ZhihuItem()
         #item['url'] = response.url
         #item['name'] = problem.xpath('//span[@class="name"]/text()').extract()
         #print item['name']
         #item['title'] = problem.xpath('//h2[@class="zm-item-title zm-editable-content"]/text()').extract()
         #item['description'] = problem.xpath('//div[@class="zm-editable-content"]/text()').extract()
         #item['answer']= problem.xpath('//div[@class=" zm-editable-content clearfix"]/text()').extract()
         #return item
