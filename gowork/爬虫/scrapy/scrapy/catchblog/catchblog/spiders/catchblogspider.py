#encoding:utf-8

# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: http://doc.scrapy.org/en/latest/topics/item-pipeline.html

from scrapy.spiders import CrawlSpider,Rule
from scrapy.linkextractors import LinkExtractor
from scrapy.selector import Selector
from catchblog.items import CatchblogItem
from scrapy.http import Request
import sys
import string
sys.stdout=open('output.txt','w')

add=0
class CatchblogSpider(CrawlSpider):
     name="catchblogspider"
     allowed_domains=["www.cnblogs.com"]
     start_urls=[
          'http://www.cnblogs.com/huhuuu'
          ]
     rules=(
          Rule(LinkExtractor(allow=(r'huhuuu/default.html\?page\=([\w]+)',)),callback='parse1',follow=True),
          Rule(LinkExtractor(allow=('huhuuu/p/',)), callback='parse_item'), 
          Rule(LinkExtractor(allow=('huhuuu/archive/',)), callback='parse_item'), 
          )
     print 'start'

     def start_requests(self):  
         print 'start_request'
         yield Request("http://www.cnblogs.com/huhuuu",  
               headers={'User-Agent': "Mozilla/5.0 (Windows NT 6.1; rv:50.0) Gecko/20100101 Firefox/50.0"}) 
     
     def parse1(self,response):
         print response
     def parse_item(self,response):
          global add
          print add
          add+=1
          #sel=Selector(response)
          

          #title=str(sel.xpath("/html/head/title/text()").extract()).encode('utf-8')
          url=response
         
          #item=CatchblogItem()
         
          #item['title']=title
          #item['url']=response
          print url

          #yield item
          #print item
          #items.append(item)

          #return items
         	 
