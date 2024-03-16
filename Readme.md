# back serve call, not auth

goctl api go -api ./define/wowfish.api -dir ./api --style goZero

abigen --abi=./define/wowfish.abi --pkg=wowfish  --type WowFishToken --out=api/contract/wowfish.go
abigen --abi=./define/NftToken.abi --pkg=wowfish  --type WowFishNft --out=api/contract/wowfishNft.go
abigen --abi=./define/wowfishBank.abi --pkg=wowfish  --type WowFishBank --out=api/contract/wowfishBank.go


# doc
goctl api plugin -plugin goctl-swagger="swagger" -api ./define/wowfish.api -dir ./doc/

# sql
goctl model mysql ddl -src="./define/db_ton_wallet.sql" -dir="./api/internal/model" --style goZero
