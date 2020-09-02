你好！
很冒昧用这样的方式来和你沟通，如有打扰请忽略我的提交哈。我是光年实验室（gnlab.com）的HR，在招Golang开发工程师，我们是一个技术型团队，技术氛围非常好。全职和兼职都可以，不过最好是全职，工作地点杭州。
我们公司是做流量增长的，Golang负责开发SAAS平台的应用，我们做的很多应用是全新的，工作非常有挑战也很有意思，是国内很多大厂的顾问。
如果有兴趣的话加我微信：13515810775  ，也可以访问 https://gnlab.com/，联系客服转发给HR。
# halite_2016_go

My entry to the [halite.io](https://halite.io) AI programming competition. 

![game](game.png)

In Halite, you have to move your pieces to expand to take over neutral territory (and other players). Each square you control adds a specified number of pieces a turn. There can only be 255 pieces in a square.

My bot is pretty simple, each cell examines the area immediately around it and tries to make the best move.

After all the moves are queued up, I loop through them and remove any moves that would clobber (add up to more than 255 troops) and cancel them.

# Results
My entry ended up 248 out of 1592 which put me in the "Silver Tier".

![results](results.png)

