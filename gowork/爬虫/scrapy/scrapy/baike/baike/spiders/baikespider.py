from scrapy.spiders import Spider
from scrapy.selector import Selector
from scrapy.http import Request
from baike.items import BaikeItem
import scrapy

class BaikeSpider(Spider):
      name="baikespider"

      start_urls=["https://movie.douban.com/chart"]

      def start_requests(self):
          yield Request("http://movie.douban.com/chart/",headers={'User-Agent': "Mozilla/5.0 (Windows NT 6.1; rv:50.0) Gecko/20100101 Firefox/50.0"})         
      def parse_item(self, response):
           
           item=BaikeItem()
           sel = Selector(response)
           
           movie_name=sel.xpath("//div[@class='pl2']/a/text()").extract()
           movie_url=sel.xpath("//div[@class='pl2']/a/@href").extract()
           movie_star=sel.xpath("//div[@class='pl2']/div/span[@class='rating_nums']/text()").extract()      

           item['movie_name']=[n.encode('utf-8') for n in movie_name]	
           item['movie_star']=[n for n in movie_star]
           item['movie_url']=[n for n in movie_url]
   
           yield item
 
    
