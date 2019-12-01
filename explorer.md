### 公链浏览器服务接口

###### 基础url:

###### 测试服：http://192.168.9.127:48080

###### 请求方法:POST

###### 请示数据格式:json

###### 公共请求数据：

| 字段 | 类型   | 说明                                             |
| ---- | ------ | ------------------------------------------------ |
| data | json   | 请求的具体数据                                   |
| sign | string | 具体数据的签名值，暂时先不用签名，需要的时候再加 |

公共响应数据：

| 字段     | 类型   | 说明                                             |
| -------- | ------ | ------------------------------------------------ |
| err_code | json   | 请求的具体数据                                   |
| err_msg  | string | 具体数据的签名值，暂时先不用签名，需要的时候再加 |
| data     | json   | 响应的具体数据                                   |

​	

1.首页接口：/explore/index

###### 请求数据：

| 字段 | 类型   | 说明             |
| ---- | ------ | ---------------- |

示例：

```json
{
	"data":{
		"addr":"bx0001"
	},
	"sign":"signedstring"
}
```



###### 响应数据：

| 字段          | 类型     | 说明        |
| ------------- | -------- | ----------- |
| chain_id     | string   | 圈子chainId  |
| block_count     | uint64   | 区块总数  |
| address_count     | uint64   | 活跃地址总数  |
| main_coin_count     | float64   | 积分币总数  |
| tx_count     | uint64   | 交易总数  |
| cross_max     | float64   | 可跨链转账总数  |
| gas_cost_count     | float64   | 全网消耗gas总数  |
| txs     | []*Transaction   | 最新交易列表  |
| blocks     | []*Block   | 最新区块列表  |

###### Transaction数据结构

| 字段         | 类型   | 说明                     |
| ------------ | ------ | ------------------------ |
| tx_type      | int8   | 交易类型，暂时为0        |
| addr_from    | string | 交易的发起者地址(出款人) |
| addr_to      | string | 交易的接收者地址(收款人) |
| amount       | string | 交易额                   |
| miner_fee    | string | 手续费                   |
| tx_hash      | string | 交易唯一标识             |
| block_height | ing64  | 交易块高                 |
| tx_time      | int64  | 交易时间                 |
| memo         | string | 备注信息                 |

###### Transaction数据结构

| 字段          | 类型   | 说明             |
| ------------- | ------ | ---------------- |
| height        | string | 区块高度          |
| Hash        | string | 区块hash           |
| TxCount        | string | 交易数量           |
| BlockTime        | string | 区块打包时间          |
| Signers        | []string | 区块签名者，多个           |

示例：

```json
{
    "err_code": 10001,
    "err_msg": "成功",
    "data": {
        "chain_id": "chainId10011",
        "block_count": 10068,
        "address_count": 22322,
        "main_coin_count": 72774,
        "tx_count": 8842,
        "cross_max": 100000,
        "gas_cost_count": 29929229,
        "txs": [
            {
                "tx_type": 0,
                "addr_from": "0x3bD6361959306B1b50797D3ff82B9A43541c3e47",
                "addr_to": "0xf6865694766c7681431E4F0C6D61AF5Ea7E37B6D",
                "amount": "100",
                "miner_fee": "21000",
                "tx_hash": "0xaf4a11b7bf6dc0d734de90b04306af3198a197adacd8d2efcd0929d2b7d5b200",
                "block_height": 5,
                "tx_time": 1575534403734,
                "memo": ""
            },
            {
                "tx_type": 0,
                "addr_from": "0x3bD6361959306B1b50797D3ff82B9A43541c3e47",
                "addr_to": "0xf6865694766c7681431E4F0C6D61AF5Ea7E37B6D",
                "amount": "100",
                "miner_fee": "21000",
                "tx_hash": "0x7bc2fd5fd7c36eef11a3b5e0c5ac6674aafda78467115b073defc2ed41647bae",
                "block_height": 4,
                "tx_time": 1573711652370,
                "memo": ""
            },
            {
                "tx_type": 0,
                "addr_from": "0x3bD6361959306B1b50797D3ff82B9A43541c3e47",
                "addr_to": "0xf6865694766c7681431E4F0C6D61AF5Ea7E37B6D",
                "amount": "100",
                "miner_fee": "21000",
                "tx_hash": "0xcaa1e2c4ceda179b4a92675967344403dd83eae3c94140fe9830c3c54a40f4ab",
                "block_height": 3,
                "tx_time": 1573711642135,
                "memo": ""
            },
            {
                "tx_type": 0,
                "addr_from": "0x3bD6361959306B1b50797D3ff82B9A43541c3e47",
                "addr_to": "0xf6865694766c7681431E4F0C6D61AF5Ea7E37B6D",
                "amount": "100",
                "miner_fee": "21000",
                "tx_hash": "0x250fb43c0a76d9f8cdbde67c0c97ffa285d9f5622ea7a7d6397c85eecc8a28d3",
                "block_height": 2,
                "tx_time": 1573711631898,
                "memo": ""
            },
            {
                "tx_type": 0,
                "addr_from": "0x3bD6361959306B1b50797D3ff82B9A43541c3e47",
                "addr_to": "0xf6865694766c7681431E4F0C6D61AF5Ea7E37B6D",
                "amount": "100",
                "miner_fee": "21000",
                "tx_hash": "0x7a4ddce8ac9be67627c2582fde3fbdd61f44b31ecf42cdb7513e9539322e9c91",
                "block_height": 1,
                "tx_time": 1573711621663,
                "memo": ""
            }
        ],
        "blocks": [
            {
                "height": 0,
                "Hash": "hash0001",
                "TxCount": 10,
                "BlockTime": 150202033823,
                "Signers": ["address001","address002","address003"]
            }
        ]
    }
}
```

2.搜索接口：/explore/search

###### 请求数据：

| 字段 | 类型   | 说明             |
| ---- | ------ | ---------------- |
| content|string|搜索内容，可以是块高or交易hashor地址 |


示例：

```json
{
	"data":{
		"content":"0"
	},
	"sign":"signedstring"
}
```



###### 响应数据【区块】：

| 字段          | 类型   | 说明             |
| ------------- | ------ | ---------------- |
| height        | string | 区块高度          |
| Hash        | string | 区块hash           |
| TxCount        | string | 交易数量           |
| BlockTime        | string | 区块打包时间          |
| Signers        | []string | 区块签名者，多个           |

###### 响应数据【交易】

| 字段         | 类型   | 说明                     |
| ------------ | ------ | ------------------------ |
| tx_type      | int8   | 交易类型，暂时为0        |
| addr_from    | string | 交易的发起者地址(出款人) |
| addr_to      | string | 交易的接收者地址(收款人) |
| amount       | string | 交易额                   |
| miner_fee    | string | 手续费                   |
| tx_hash      | string | 交易唯一标识             |
| block_height | ing64  | 交易块高                 |
| tx_time      | int64  | 交易时间                 |
| memo         | string | 备注信息                 |

# 响应数据【地址】 

###### 响应数据【地址】：

| 字段          | 类型   | 说明             |
| ------------- | ------ | ---------------- |
| name        | string | 名称          |
| contract        | string | 合约地址           |
| type        | string | 类型           |
| symbol        | string | 简称          |
| quantity        | int64 | 数量           |

示例：

```json
{
    "err_code": 10001,
    "err_msg": "成功",
    "data": {
        "asset_list": [
            {
                "name": "酒财币",
                "contract": "0xaaafffbbbccceeed0002223",
                "type": "积分币",
                "symbol": "JCB",
                "quantity": 23323244
            },
            {
                "name": "二哈币",
                "contract": "0xaaafffbbbccceeed0002224",
                "type": "BRCn",
                "symbol": "RHB",
                "quantity": 23323245
            }
        ]
    }
}
```

3.交易列表接口：/explore/txlist

###### 请求数据：

| 字段 | 类型   | 说明             |
| ---- | ------ | ---------------- |
| page|uint32|第几页，从1开始 |
| page_size|uint32|每页的数量 |

示例：

```json
{
	"data":{
		"content":"0"
	},
	"sign":"signedstring"
}
```



###### 响应数据：

| 字段          | 类型   | 说明             |
| ------------- | ------ | ---------------- |
| txs        | []Transaction | 区块高度          |



示例：

```json
{
    "err_code": 10001,
    "err_msg": "成功",
    "data": {
        "txs": [
            {
                "tx_type": 0,
                "addr_from": "0x3bD6361959306B1b50797D3ff82B9A43541c3e47",
                "addr_to": "0xf6865694766c7681431E4F0C6D61AF5Ea7E37B6D",
                "amount": "100",
                "miner_fee": "21000",
                "tx_hash": "0xaf4a11b7bf6dc0d734de90b04306af3198a197adacd8d2efcd0929d2b7d5b200",
                "block_height": 5,
                "tx_time": 1575534403734,
                "memo": ""
            },
            {
                "tx_type": 0,
                "addr_from": "0x3bD6361959306B1b50797D3ff82B9A43541c3e47",
                "addr_to": "0xf6865694766c7681431E4F0C6D61AF5Ea7E37B6D",
                "amount": "100",
                "miner_fee": "21000",
                "tx_hash": "0x7bc2fd5fd7c36eef11a3b5e0c5ac6674aafda78467115b073defc2ed41647bae",
                "block_height": 4,
                "tx_time": 1573711652370,
                "memo": ""
            },
            {
                "tx_type": 0,
                "addr_from": "0x3bD6361959306B1b50797D3ff82B9A43541c3e47",
                "addr_to": "0xf6865694766c7681431E4F0C6D61AF5Ea7E37B6D",
                "amount": "100",
                "miner_fee": "21000",
                "tx_hash": "0xcaa1e2c4ceda179b4a92675967344403dd83eae3c94140fe9830c3c54a40f4ab",
                "block_height": 3,
                "tx_time": 1573711642135,
                "memo": ""
            },
            {
                "tx_type": 0,
                "addr_from": "0x3bD6361959306B1b50797D3ff82B9A43541c3e47",
                "addr_to": "0xf6865694766c7681431E4F0C6D61AF5Ea7E37B6D",
                "amount": "100",
                "miner_fee": "21000",
                "tx_hash": "0x250fb43c0a76d9f8cdbde67c0c97ffa285d9f5622ea7a7d6397c85eecc8a28d3",
                "block_height": 2,
                "tx_time": 1573711631898,
                "memo": ""
            },
            {
                "tx_type": 0,
                "addr_from": "0x3bD6361959306B1b50797D3ff82B9A43541c3e47",
                "addr_to": "0xf6865694766c7681431E4F0C6D61AF5Ea7E37B6D",
                "amount": "100",
                "miner_fee": "21000",
                "tx_hash": "0x7a4ddce8ac9be67627c2582fde3fbdd61f44b31ecf42cdb7513e9539322e9c91",
                "block_height": 1,
                "tx_time": 1573711621663,
                "memo": ""
            }
        ]
    }
}
```

4.区块列表接口：/explore/getblock

###### 请求数据：

| 字段 | 类型   | 说明             |
| ---- | ------ | ---------------- |
| page|uint32|第几页，从1开始 |
| page_size|uint32|每页的数量 |

示例：

```json
{
	"data":{
		"content":"0"
	},
	"sign":"signedstring"
}
```



###### 响应数据：

| 字段          | 类型   | 说明             |
| ------------- | ------ | ---------------- |
| blocks        | []Block | 区块们   |



示例：

```json
{
    "err_code": 10001,
    "err_msg": "成功",
    "data": {
        "blocks": [
            {
                "height": 0,
                "Hash": "hash0001",
                "TxCount": 10,
                "BlockTime": 150202033823,
                "Signers": null
            }
        ]
    }
}
```

5.圈子资产列表：/explore/getassets

###### 请求数据：

| 字段 | 类型   | 说明             |
| ---- | ------ | ---------------- |
| page|uint32|第几页，从1开始 |
| page_size|uint32|每页的数量 |

示例：

```json
{
	"data":{
		"content":"0"
	},
	"sign":"signedstring"
}
```



###### 响应数据：

| 字段          | 类型   | 说明             |
| ------------- | ------ | ---------------- |
| asset_list        | []AssetInfo | 区块高度   |


###### AssetInfo结构

| 字段          | 类型   | 说明             |
| ------------- | ------ | ---------------- |
| name        | string | 名称          |
| contract        | string | 合约地址           |
| type        | string | 类型           |
| symbol        | string | 简称 |
| quantity        | []string | 数量|

示例：

```json
{
    "err_code": 10001,
    "err_msg": "成功",
    "data": {
        "asset_list": [
            {
                "name": "酒财币",
                "contract": "0xaaafffbbbccceeed0002223",
                "type": "积分币",
                "symbol": "JCB",
                "quantity": 23323244
            },
            {
                "name": "二哈币",
                "contract": "0xaaafffbbbccceeed0002224",
                "type": "BRCn",
                "symbol": "RHB",
                "quantity": 23323245
            }
        ]
    }
}
```
6.获取指定区块：/explore/getblockbyid

###### 请求数据：

| 字段 | 类型   | 说明             |
| ---- | ------ | ---------------- |
| height|uint64|区块高度 |


示例：

```json
{
	"data":{
		"height":"0"
	},
	"sign":"signedstring"
}
```



###### 响应数据：

| 字段          | 类型   | 说明             |
| ------------- | ------ | ---------------- |
| hash|uint64|区块hash |
| tx_count|uint64|区块中交易的数量 |
| block_time|int64|区块时间 |
| signers|[]string|区块签名者们 |



示例：

```json
{
    "err_code": 10001,
    "err_msg": "成功",
    "data": {
        "height": 0,
        "hash": "hash0001",
        "tx_count": 10,
        "block_time": 150202033823,
        "signers": null
    }
}
```

7.获取指定交易：/explore/gettxbyhash

###### 请求数据：

| 字段 | 类型   | 说明             |
| ---- | ------ | ---------------- |
| txnash|string|交易hash |



示例：

```json
{
	"data":{
		"txnash":"0xaaaaabbbb"
	},
	"sign":"signedstring"
}
```



###### 响应数据：

| 字段          | 类型   | 说明             |
| ------------- | ------ | ---------------- |
| tx_type      | int8   | 交易类型，暂时为0        |
| addr_from    | string | 交易的发起者地址(出款人) |
| addr_to      | string | 交易的接收者地址(收款人) |
| amount       | string | 交易额                   |
| miner_fee    | string | 手续费                   |
| tx_hash      | string | 交易唯一标识             |
| block_height | ing64  | 交易块高                 |
| tx_time      | int64  | 交易时间                 |
| memo         | string | 备注信息                 |



示例：

```json
{
    "err_code": 10001,
    "err_msg": "成功",
    "data": {
        "tx_type": 0,
        "addr_from": "0x3bD6361959306B1b50797D3ff82B9A43541c3e47",
        "addr_to": "0xf6865694766c7681431E4F0C6D61AF5Ea7E37B6D",
        "amount": "100",
        "miner_fee": "21000",
        "tx_hash": "0x250fb43c0a76d9f8cdbde67c0c97ffa285d9f5622ea7a7d6397c85eecc8a28d3",
        "block_height": 2,
        "tx_time": 1573711631898,
        "memo": ""
    }
}
```



##### 错误码

| 错误码 | 描述           |
| ------ | -------------- |
| 10001  | 成功           |
| 10002  | json参数错误   |
| 10003  | 服务端查询错误 |
| 10004  | 公链错误       |

