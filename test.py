import pymysql.cursors
import pymysql,redis


'''conn = pymysql.connect( host='0.0.0.0',port=3307,user='root',
                        password='root12345',db='entries',charset='utf8mb4',
                        cursorclass=pymysql.cursors.DictCursor)

sql = 'select  word,definition from entries '

cursor = conn.cursor()

cursor.execute(sql)

result = cursor.fetchall()
d={}

for el in result:
    #print(type(el))
    d[  el["word"].lower()  ] = el["definition"].lower()
'''

r = redis.Redis( host='localhost',port=6379,password="")

r.set("algebraically","relating to, involving, or according to the laws of algebra")

r.close()