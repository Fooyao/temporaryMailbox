package main

import (
	"github.com/gin-gonic/gin"
	"math/rand"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")
var htmlIndex = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>随机邮箱生成器</title>
    <style>
        body {
            font-family: 'Segoe UI', 'Roboto', sans-serif;
            background-color: #f2f4f8;
            margin: 0;
            padding: 20px;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            height: 100vh;
        }

        h1 {
            color: #204056;
        }

        .buttons {
            margin-top: 20px;
            display: flex;
            flex-direction: row;
            gap: 10px;
        }

        button {
            padding: 10px 20px;
            font-size: 16px;
            cursor: pointer;
            background-color: #4CAF50;
            color: white;
            border: none;
            border-radius: 5px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.2);
        }

        button:hover {
            background-color: #368f3a;
        }

        #domainSelect, #username {
            padding: 10px;
            border-radius: 5px;
            border: 1px solid #ccc;
            font-size: 16px;
        }

        #randomAddress, #result {
            margin-top: 20px;
            padding: 10px;
            border: 1px solid #bbb;
            border-radius: 5px;
            background-color: #fff;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
            width: 100%;
            max-width: 600px;
            word-break: break-word;
        }
    </style>
</head>
<body>
    <h1><a href="https://hdd.cm">号多多hdd.cm|推特低至2毛</a></h1>
	<div class="buttons">
        <select id="domainSelect">
            <!-- 域名选项将动态添加到这里 -->
        </select>
        <input type="text" id="username" placeholder="自定义用户名">
        <button id="getRandom">获取邮箱</button>
        <button id="getMailList">获取邮件列表</button>
        <button id="getMail">获取邮件</button>
    </div>
    <div id="randomAddress"></div>
    <div id="result"></div>

    <script>
        let randomMailAddress = ''; // 保存随机邮件地址

 	window.onload = function() {
            fetch('` + mainDomain + `/getAllowedDomains')
                .then(response => response.json())
                .then(data => {
                    const select = document.getElementById('domainSelect');
                    data.allowedDomains.forEach(domain => {
                        const option = document.createElement('option');
                        option.value = domain;
                        option.innerText = domain;
                        select.appendChild(option);
                    });
                })
                .catch(error => console.error('Error:', error));
        };

        // 获取随机邮件地址
        document.getElementById('getRandom').addEventListener('click', function() {
	    const selectedDomain = document.getElementById('domainSelect').value;
     	    const username = document.getElementById('username').value;
     		
            randomMailAddress = username + '@' + selectedDomain; // 保存随机邮件地址
            document.getElementById('randomAddress').innerText = '邮件地址: ' + randomMailAddress;
        });

        // 获取邮件列表
        document.getElementById('getMailList').addEventListener('click', function() {
            if (!randomMailAddress) {
                alert('请先获取随机邮件地址！');
                return;
            }

            fetch('` + mainDomain + `/getMailList/${randomMailAddress}')
                .then(response => response.json())
                .then(data => {
                    document.getElementById('result').innerHTML = '';
                    if(data.mails.length==0){
                        document.getElementById('result').innerText = '邮件还没收到，耐心等一下哦';
                    }
else{
		data.mails.forEach(mail => {
                        const mailElement = document.createElement('div');
                        mailElement.innerText = ` + "`发件人: ${mail.from}, 标题: ${mail.title}`" + `;
                        document.getElementById('result').appendChild(mailElement);
                    });}

                })
                .catch(error => console.error('Error:', error));
        });

        // 获取邮件内容
        document.getElementById('getMail').addEventListener('click', function() {
            if (!randomMailAddress) {
                alert('请先获取随机邮件地址！');
                return;
            }

            // 调用后端接口，这里假设使用 fetch
            fetch('` + mainDomain + `/getMail/${randomMailAddress}')
                .then(response => response.json())
                .then(data => {
                    // 假设展示邮件内容
                    document.getElementById('result').innerText = '邮件内容: ' + data.mail.content;
                })
                .catch(error => console.error('Error:', error));
        });
    </script>
</body>
</html>
`

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
func startHTTPServer(allowedDomains []string) {
	gin.SetMode(gin.ReleaseMode)
	httpsrv := gin.Default()
	httpsrv.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, htmlIndex)
	})
	httpsrv.GET("/getAllowedDomains", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"allowedDomains": allowedDomains,
		})
	})
	httpsrv.GET("/getAddress", func(c *gin.Context) {
		tmp := RandStringRunes(8)
		randomIndex := rand.Intn(len(allowedDomains))
		Domain := allowedDomains[randomIndex]
		c.JSON(200, gin.H{
			"random":  tmp,
			"address": tmp + "@" + Domain,
			"domain":  "@" + Domain,
		})
	})
	httpsrv.GET("/getMailList/:randomString", func(c *gin.Context) {
		mailHead := c.Param("randomString")
		if mailBox[mailHead] == nil {
			c.JSON(200, gin.H{
				"mails": make([]string, 0),
			})
		} else {
			mails := make([]gin.H, len(mailBox[mailHead]))
			for i, v := range mailBox[mailHead] {
				mails[i] = gin.H{"from": v.from, "title": v.title}
			}
			c.JSON(200, gin.H{
				"mails": mails,
			})
		}
	})
	httpsrv.GET("/getMail/:randomString", func(c *gin.Context) {
		mailHead := c.Param("randomString")
		if mailBox[mailHead] == nil {
			c.JSON(200, gin.H{
				"mail": "没有邮件",
			})
		} else if len(mailBox[mailHead]) == 0 {
			c.JSON(200, gin.H{
				"mail": "没有邮件",
			})
		} else {
			tmpMail := mailBox[mailHead][len(mailBox[mailHead])-1]
			mailBox[mailHead] = mailBox[mailHead][0 : len(mailBox[mailHead])-1]
			c.JSON(200, gin.H{
				"mail": gin.H{
					"from":    tmpMail.from,
					"title":   tmpMail.title,
					"content": tmpMail.content,
				},
			})
		}
	})
	err := httpsrv.RunTLS(":443", "./certs/server.pem", "./certs/server.key")
	if err != nil {
		return
	}
	err = httpsrv.Run(":80")
	if err != nil {
		return
	}
}
