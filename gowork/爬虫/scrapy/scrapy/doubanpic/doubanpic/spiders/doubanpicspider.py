from scrapy.spiders import Spider
from scrapy.selector import Selector
from doubanpic.items import DoubanpicItem
import scrapy

class DoubanNewMovieSpider(Spider):
     name="doubanpicspider"
     start_urls=[
          'http://tieba.baidu.com/p/4023230951'
          ]

     #def start_requests(self):  
         #yield scrapy.FormRequest("http://movie.douban.com/chart/",  
               #headers={'User-Agent': "Mozilla/5.0 (Windows NT 6.1; rv:50.0) Gecko/20100101 Firefox/50.0"}) 
     
     def parse(self,response):
          sel=Selector(response)

          image_url=sel.xpath("//div[@id='post_content_75283192143']/img[@class='BDE_Image']/@src").extract()
          print 'the urls:/n'
          print image_url
          print '/n'
        
          item=DoubanpicItem()
         
          item['image_url']=image_url

          yield item
          
