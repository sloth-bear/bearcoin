# Bearcoin

## Usage
```
go run main.go --mode=<choose 'html' or 'rest'>[required] --port=<port>[optional]
```


## Goals
1. Understand blockchain
2. Create blockchain
3. What makes up hash
4. What makes up mining
5. What makes up block
6. Build explorer of blocks
7. Create a JSON API to explore blockchain to interact with blockchain 
8. How we can have coins in blockchain
9. Why blockchains are so good if you want to have coins
10. Why, and how we can protect internet money using a blockchain + why it is good idea
11. Learn about transactions (transaction input, transaction output) with owner those coins
12. Learn about wallet modules
13. Security, validations 
14. How we can communicate between nodes in our network to decentralize peer to peer network for sharing blockchain to people

And this coin is useless and valueless, it is only for understanding coin and blockchain technology.


## Language
- [Golang](https://go.dev/)


## Database
- [Bolt](https://github.com/boltdb/bolt)


## Concepts
### Mempool
Mempool(Memory pool) is a place where we put all the unconfirmed transactions. The transactions are confirmed when miners the people making blocks include the transactions into the block.

1. The miners find a block. (mining)
2. The miners go to mempool, and they put unconfirmed transactions into the block.
3. Transactions are confirmed.


## Local Documentation
```
godoc -http=:6060
```


## Testing
```
go test ./... -v
```