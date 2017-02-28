# -*- coding: utf-8 -*-

# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: http://doc.scrapy.org/en/latest/topics/item-pipeline.html
import MySQLdb
import datetime
class ZhihuPipeline(object):
    def __init__(self):
        self.conn = MySQLdb.connect(user = "root", passwd = "123456", db = "test", host = "127.0.0.1", charset = 'utf8', use_unicode = True)
        self.cursor = self.conn.cursor()
        # 清空表
        # self.cursor.execute('truncate table weather;')
        # self.conn.commit()

    def process_item(self, item, spider):
        #curTime = datetime.datetime.now()
        try:
            self.cursor.execute(
                """INSERT IGNORE INTO users (url,name,desc,sextype,prise,ganxie,shoucang,following,follower)
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s)""",
                (
                    item['url'],
                    item['name'], 
                    item['desc'], 
                    item['sextype'],                    
                    item['prise'], 
                    item['ganxie'], 
                    item['shoucang'], 
                    item['following'], 
                    item['follower'],                     
                )
            )
            self.conn.commit()
        except MySQLdb.Error, e:
            print 'Error %d %s' % (e.args[0], e.args[1])

        return item
