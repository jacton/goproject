# -*- coding: utf-8 -*-
import scrapy
import os
import time
from zhihu.items import UserItem
from zhihu.myconfig import UsersConfig
from scrapy.selector import Selector
import sys
sys.stdout=open('output.txt','w')
add=0
class UsersSpider(scrapy.Spider):
    name = 'usersspider'
    dm = 'https://www.zhihu.com'
    domain = 'https://www.zhihu.com/people/ganganray/followers'
    headers = {
        "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
        "Accept-Language": "zh-CN,zh;q=0.8",
        "Connection": "keep-alive",
        "Host": "www.zhihu.com",
        "Upgrade-Insecure-Requests": "1",
        "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.109 Safari/537.36"
    }

    def start_requests(self):
        print "start"
        yield scrapy.Request(
            url = self.domain,
            headers = self.headers,
            callback = self.user_start,
            dont_filter = True
        )
    def user_item(self, response):   
        url=response.url       
        nameinfo = response.xpath('//div[@class="ProfileHeader-contentHead"]/h1/span/text()').extract()
        sextype=0 if response.xpath('//svg[contains(@class, "Icon--female")]') else 1
        #company = response.xpath('//svg[@class="Icon--company"]/parent::*/parent::div/text()').extract() 
        #edu = response.xpath('//svg[@class="Icon--education"]/parent::*/parent::div/text()').extract() 
        #sextype = response.xpath('//div[@class="ProfileHeader-infoItem"]/div[@class="ProfileHeader-iconWrapper"]/svg/@class').extract()
        #strcompany=""
        #stredu=""
        if len(nameinfo) == 2:
            name=nameinfo[0].encode('utf-8')
            desc=nameinfo[1].encode('utf-8')
            print name+"  "+desc
            print sextype
            print url
            #for n in company:
            #    strcompany+=n.encode("utf-8")
            #for n in edu:
            #    stredu+=n.encode("utf-8")
            #print "company:"+strcompany
            #print "edu:"+stredu
            #item = UserItem()

            #item['url'] = url
            #item['name'] = name
            #yield item
    def user_start(self, response):
        global add
        sel=Selector(response)
        sel_following_url = sel.xpath('//div[@class="UserItem-title"]/span/div/div/a/@href').extract()
        sel_following_name = sel.xpath('//div[@class="UserItem-title"]/span/div/div/a/text()').extract()
        # 判断关注列表是否为空
        self.user_item(response)
        if len(sel_following_url):
            for people_url in sel_following_url:
                #print add
                #add+=1
                #print self.dm+people_url
                #yield scrapy.Request(
                #    url = self.dm+people_url,
                #    headers = self.headers,
                #    callback = self.user_item,
                #    dont_filter = True
                #)
                yield scrapy.Request(
                    url = self.dm+people_url + '/following',
                    headers = self.headers,
                    callback = self.user_start,
                    dont_filter = True
                )
                #yield scrapy.Request(
                #    url = self.dm+people_url + '/followers',
                #    headers = self.headers,
                #    callback = self.user_start,
                #    dont_filter = True
                #)


         
