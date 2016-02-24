import azure_demo_create_vm_from_ami as azureapi
import random
import time
import sys
import urllib2
from bs4 import BeautifulSoup
from datetime import datetime


storage_account = "zhouchelp2storage"
lg_dns = 'loadgenerator9795264vm.eastus.cloudapp.azure.com'
dc_dns = 'datacenter162122vm.eastus.cloudapp.azure.com'

#url = lg_dns
#conn = httplib.HTTPConnection(url)

#conn.request('GET', '/password?passwd=2sxlZfOZMO0ruePcHNVcg3vrSYg39n23&andrewid=zhouchel')
#response = conn.getresponse()
#html = response.read()
#print html


#response = urllib2.urlopen('http://' + lg_dns + '/password?passwd=2sxlZfOZMO0ruePcHNVcg3vrSYg39n23&andrewid=zhouchel')
#html = response.read()
#print html

#params = urllib.urlencode({'dns': dc_dns})
#headers = {"Content-type": "application/x-www-form-urlencoded", "Accept": "text/plain"}
#conn.request('POST', '/test/horizontal', params, headers)
#response = conn.getresponse()
#html = response.read()
#print html

#response = urllib2.urlopen('http://' + lg_dns + '/test/horizontal/add?dns='+ dc_dns)
#u = urllib2.urlopen('http://' + lg_dns + '/test/horizontal', params)
#conn.request('POST', '/test/horizontal', params, headers)
#response = conn.getresponse()
#html = response.read()
#print html

    # get log url of this test
#soup = BeautifulSoup(html, 'html.parser')
#log = soup.find('a', href=True)['href']

log = '/log?name=test.1454360059306.log'
html = urllib2.urlopen('http://' + lg_dns + log)
content = html.read()
content = content.split('\n')
qps = 0
for i in range(len(content)-1, -1, -1):
    if content[i].startswith('[Minute'):
        break
    ans = content[i].split('=')
    if len(ans) == 2:
        qps += float(ans[1])
if qps == 0:
    print 'ERROR'
else:
    print datetime.now(), 'qps='+str(qps)
