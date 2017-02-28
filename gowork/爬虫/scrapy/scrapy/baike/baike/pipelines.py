# -*- coding: utf-8 -*-

# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: http://doc.scrapy.org/en/latest/topics/item-pipeline.html

from twisted.enterprise import adbapi

import MySQLdb
import MySQLdb.cursors

class BaikePipeline(object):
    def process_item(self, item, spider):
        query=self.dbpool.runIteraction(self._conditional_insert,item)
        query.addErrback(self.handle_error)
        
    def __init__(self):
        self.dbpool=adbapi.ConnectionPool('MysQLdb',
                                          db='test',
                                          user='root',
                                          passwd='123456',
                                          cursorclass=MySQLdb.cursors.DictCursor,
                                          charset='utf8',
                                          use_unicode=False
                                          )
    def _conditional_insert(self,tx,item):
        for i in range(len(item['movie_star'])):
            movie_name=str(item['movie_name'][i]).replace(' ','').replace('\n','').replace('/','')
            movie_star=item['movie_star'][i]
            movie_url=item['movie_url'][i]      

            tx.execute( \
                "insert into baike(id,title,url,content) \
                 values(null,%s,%s,%s)",
                (movie_name,movie_url,movie_star)
                )
                                          
                                
