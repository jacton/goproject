#encoding:utf-8

# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: http://doc.scrapy.org/en/latest/topics/item-pipeline.html

from scrapy.spiders import CrawlSpider,Rule
from scrapy.linkextractors import LinkExtractor
from zhihu.items import ZhihuItem
from scrapy.selector import Selector
from scrapy.http import Request
import sys
import string
import re
sys.stdout=open('output.txt','w')

add=0
class ZhihuuserSpider(CrawlSpider):
     name="zhihuuserspider"
     #allowed_domains=["https://www.zhihu.com"]
     start_urls=[
          'https://www.zhihu.com/people/wen-yi-yang-81'
          ]
     rules=(
          #Rule(LinkExtractor(allow=('people/',)), restrict_xpaths=('//div[@class="review-more"]/a')),callback='parse_item',follow=True),
          Rule(LinkExtractor(allow=('people/',),deny=('people/(.+)/(.+)',)),callback='parse_item',follow=True),
          Rule(LinkExtractor(allow=('people/(.+)/following$',))), 
          #Rule(LinkExtractor(allow=('question/',)), callback='parse_item',follow=True), 
          )
     print 'start'

     #def start_requests(self):  
     #    print 'start_request'
     #    yield Request("https://www.zhihu.com/explore",  
     #          headers={'User-Agent': "Mozilla/5.0 (Windows NT 6.1; rv:50.0) Gecko/20100101 Firefox/50.0"},
     #          callback=self.parse_item,
     #          dont_filter='True'
     #          ) 
     
     def parse1(self,response):
         print response
     def parse_item(self,response):
        url=response.url  
        print url
        nameinfo=response.xpath('//div[@class="ProfileHeader-contentHead"]/h1/span/text()').extract()
        priseinfo = response.xpath('//div[@class="IconGraf"]/text()').extract()
        followinfo = response.xpath('//div[@class="NumberBoard-value"]/text()').extract()
        shoucanginfo = response.xpath('//div[@class="Profile-sideColumnItemValue"]/text()').extract()
        #edu = response.xpath('//svg[@class="Icon--education"]/parent::*/parent::div/text()').extract()         
        sextype=0 if response.xpath('//svg[contains(@class, "Icon--female")]') else 1
        name=""
        desc=""
        prise=""
        strprise=""
        shoucang=""
        strshou=""
        following=""
        follower=""
        nprise=0
        nganxie=0
        nshoucang=0
        nfollowing=0
        nfollower=0
        
        if len(nameinfo) == 2:                 
            name = nameinfo[0].encode('utf-8')
            desc = nameinfo[1].encode('utf-8')
            print name
        if len(priseinfo) == 1:    
            prise = priseinfo[0].encode('utf-8') 
            strprise=re.findall('\d+',prise)
            if strprise is not None:
                nprise=int(strprise[0])
        if len(shoucanginfo) == 1:    
            shoucang = shoucanginfo[0].encode('utf-8')
            nshou=re.findall('\d+',shoucang)
            if nshou is not None:
                nganxie=int(nshou[0])
                nshoucang=int(nshou[1])
        if len(followinfo) == 2:                 
            nfollowing = int(followinfo[0])
            nfollower = int(followinfo[1])     
           
        if len(nameinfo) == 2:  
            item = ZhihuItem()
            item['url'] = url
            item['name'] = name
            item['desc'] = desc
            item['sextype'] = sextype
            item['prise'] = nprise
            item['ganxie'] = nganxie
            item['shoucang'] = nshoucang
            item['following'] = nfollowing
            item['follower'] = nfollower
            yield item

         	 
