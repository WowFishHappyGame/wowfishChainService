# back serve call, not auth

goctl api go -api ./define/wowfish.api -dir ./api --style goZero

abigen --abi=./define/wowfish.abi --pkg=wowfish  --type WowFishToken --out=api/contract/wowfish.go
abigen --abi=./define/NftToken.abi --pkg=wowfish  --type WowFishNft --out=api/contract/wowfishNft.go