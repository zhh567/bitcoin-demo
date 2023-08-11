# 钱包

## Wallet

1. 私钥
2. 公钥
3. 创建密钥对


# 使用密钥地址

1. 使用生成的地址代替 string
2. 使用私钥进行签名、校验


# 对区块签名

1. 每个input都要有一个签名
2. 签名是对当前交易的签名
3. 签名交易需要哪些数据？
    1. 每个输出的value
    2. 每个输出的公钥哈希（收款人）
    3. input 引用的output的公钥哈希（付款人）
4. input签名流程：(创建交易记录时签名)
    1. 将区块复制一份副本
    2. 将副本的 pubKey、Sig 字段置为空nil
    3. 遍历副本的每个input，将当前的pubKey字段设置为引用的output的公钥哈希
    4. 对当前呃交易进行签名，过程等于计算交易id
    5. 得到的sig写入到原区块的Sig字段
    6. 将当前input的pubKey字段、Sig字段置空，然后继续下一个input
5. 矿工校验流程：（添加区块时验证签名）
    1. 矿工将交易复制一份，在做修改
    2. 唯一要还原的是签名的数据
    3. input 引用的公钥的哈希放到pubKey字段，对整体交易哈希，得到付款人签名的原始数据
    4. 使用 tx.sig、tx.pubKey、得到的数据 进行校验


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