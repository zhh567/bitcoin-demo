# 交易结构 *P2SH*

1. 交易id
2. 交易输出input，由历史中某个output转换而来（可有多个）
    1. 引用的交易id
    2. 此交易中对应output的索引
    3. **解锁脚本**（发起者的签名、公钥）
3. 交易输出output， 表明流向（可有多个）
    1. **锁定脚本**（收款人的地址，可反推出公钥哈希）
    2. 转账金额
4. 时间戳




## 实现流程

1. 获取交易ID
    对交易做hash
2. 创建挖矿交易
    没有input，只有一个output
    挖矿交易需要能识别出来，因为没input，所以不需要签名
3. 使用 Transaction 改写程序
4. 获取挖矿人金额
5. 创建普通交易
6. 转账


## 收集用户余额

1. 从后向前遍历链，寻找属于用户的output（最后一个区块的output可直接加入余额）
2. 遍历的同时将所有用户消耗的input放入容器
3. 当遍历相关output时判断其是否已经被消耗（TxId和索引是否与某个input的相同），没有被消耗则加入余额。


## 创建普通交易

> from, to, amount

1. 遍历账本，找到关于from的utxo集合，返回总金额
2. 金额不足，创建失败
3. 拼接 inputs
    遍历utxo集合，每个output都转换为一个input
4. 拼接 outputs
    创建一个属于to的output
    如果总金额大于转账金额，给from创建找零output
5. 设置hash

## 转账操作

1. 每次send时，添加一个区块
2. 创建挖矿交易，创建普通交易，并添加
3. 执行 addBlock


## IsCoinbase

遍历inputs时，判断是否是挖矿交易。
如果是则去除。签名时也要去除。


```bash
# useful command
help
getnewaddress ( "account" "address_type" )
getblock BLCCK_ADDRESS
generatetoaddress nblocks address (maxtries) //指定挖矿⼈

gettransaction "txid" ( include_watchonly )      # 通过交易id获取交易信息
decoderawtransaction "hexstring" ( iswitness )   # hexstring为交易信息的hex字段，
decodescript "hexstring"

```

test1:
bcrt1qfncen2fuzkpm6fmtcv4y60vnhga87e7dnuytw3
bcrt1qvcmp37pm580n9szhhkuqplnvujk6lw8mcx5q2m
new:
bcrt1qh0xsncgdtzdxu3zxfp6dgqs5xj35qurp8z7npd


test2:
bcrt1qy2v0jnpsqq68vmqqu35lgj54ghe4wkd77vqzy8

```sh
== Blockchain ==
getbestblockhash //最后⼀个区块的哈希
getblock "blockhash" ( verbosity )
getblockchaininfo
getblockcount
getblockhash height
getblockheader "hash" ( verbose )
getchaintips
getchaintxstats ( nblocks blockhash ) //统计区块数量，交易数量
getdifficulty
getmempoolancestors txid (verbose) //必须是在内存中的交易才有效
getmempooldescendants txid (verbose) //TODO
getmempoolentry txid //TODO
getmempoolinfo //当前内存中的交易数据，可读，别和getmemoryinfo混淆了
getrawmempool ( verbose ) //返回交易id
gettxout "txid" n ( include_mempool ) //引⽤交易id内的第n个output
gettxoutproof ["txid",...] ( blockhash ) //TODO
gettxoutsetinfo //统计utxo
preciousblock "blockhash" //TODO
pruneblockchain //TODO，必须在prune mode才能删除旧区块
savemempool //将内存交易保存到磁盘中，当前内存中有2笔交易，为何返回null TODO
verifychain ( checklevel nblocks ) //返回true
verifytxoutproof "proof" //TODO

== Control ==
getmemoryinfo ("mode") //TODO
help ( "command" )
logging ( <include> <exclude> ) //会返回⼀个值全0的json，不知道如何产⽣数据
stop //直接退出当前BitCoin Core，⼿残
uptime //当前客户端启动多久了

== Generating ==
generate nblocks ( maxtries ) //⼿动执⾏挖矿，每个块奖励50BTC，由默认账户挖矿
generatetoaddress nblocks address (maxtries) //指定挖矿⼈

== Mining ==
getblocktemplate ( TemplateRequest ) //TODO，得联⽹？？
getmininginfo
getnetworkhashps ( nblocks height ) //算⼒？？8.913976854892111e-07
prioritisetransaction <txid> <dummy value> <fee delta>
submitblock "hexdata" ( "dummy" ) //TODO

== Network ==
addnode "node" "add|remove|onetry"
clearbanned
disconnectnode "[address]" [nodeid]
getaddednodeinfo ( "node" )
getconnectioncount
getnettotals
getnetworkinfo
getpeerinfo
listbanned
ping
setban "subnet" "add|remove" (bantime) (absolute)
setnetworkactive true|false

== Rawtransactions == //TODO
combinerawtransaction ["hexstring",...]
createrawtransaction [{"txid":"id","vout":n},...] {"address":amount,"data":"hex",...} ( locktime ) ( eplaceable )
decoderawtransaction "hexstring" ( iswitness )
decodescript "hexstring"
fundrawtransaction "hexstring" ( options iswitness )
getrawtransaction "txid" ( verbose "blockhash" )
sendrawtransaction "hexstring" ( allowhighfees )
signrawtransaction "hexstring" ( [{"txid":"id","vout":n,"scriptPubKey":"hex","redee
mScript":"hex"},...] ["privatekey1",...] sighashtype )

== Util ==
createmultisig nrequired ["key",...]
estimatefee nblocks
estimatesmartfee conf_target ("estimate_mode")
signmessagewithprivkey "privkey" "message"
validateaddress "address"
verifymessage "address" "signature" "message"

== Wallet ==
abandontransaction "txid"
abortrescan
addmultisigaddress nrequired ["key",...] ( "account" "address_type" )
backupwallet "destination"
bumpfee "txid" ( options )
dumpprivkey "address"
dumpwallet "filename"
encryptwallet "passphrase"
getaccount "address" //user1
getaccountaddress "account"
getaddressesbyaccount "account"
getbalance ( "account" minconf include_watchonly )
getnewaddress ( "account" "address_type" )
getrawchangeaddress ( "address_type" )
getreceivedbyaccount "account" ( minconf )
getreceivedbyaddress "address" ( minconf )
gettransaction "txid" ( include_watchonly )
getunconfirmedbalance
getwalletinfo
importaddress "address" ( "label" rescan p2sh )
importmulti "requests" ( "options" )
importprivkey "privkey" ( "label" ) ( rescan )
importprunedfunds
importpubkey "pubkey" ( "label" rescan )
importwallet "filename"
keypoolrefill ( newsize )
listaccounts ( minconf include_watchonly)
listaddressgroupings
listlockunspent
listreceivedbyaccount ( minconf include_empty include_watchonly)
listreceivedbyaddress ( minconf include_empty include_watchonly)
listsinceblock ( "blockhash" target_confirmations include_watchonly include_removed)
listtransactions ( "account" count skip include_watchonly)
listunspent ( minconf maxconf ["addresses",...] [include_unsafe] [query_options])
listwallets
lockunspent unlock ([{"txid":"txid","vout":n},...])
move "fromaccount" "toaccount" amount ( minconf "comment" )
removeprunedfunds "txid"
rescanblockchain ("start_height") ("stop_height")
sendfrom "fromaccount" "toaddress" amount ( minconf "comment" "comment_to" )
sendmany "fromaccount" {"address":amount,...} ( minconf "comment" ["address",...] r
eplaceable conf_target "estimate_mode")
sendtoaddress "address" amount ( "comment" "comment_to" subtractfeefromamount repla
ceable conf_target "estimate_mode")
setaccount "address" "account"
settxfee amount
signmessage "address" "message"
walletlock
walletpassphrase "passphrase" timeout
walletpassphrasechange "oldpassphrase" "newpassphrase"
```