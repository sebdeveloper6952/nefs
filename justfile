encrypt:
    go run . encrypt --sk 6bd65bac933303eeb3d9bb16bcca12dd29ff7e3406c2da9adfdef858c799b316 --pk 0cc3b3f7c99922c233bb73a1069b48cbf183c42b1e00464c61c95221b104f567 --file id_back.png

decrypt:
    go run . decrypt --sk eb1fd85a7266d1afa7b07cc6db2126cd44dcf5881313a5e27455a59ca4d6d67a --pk 3762b6a5d037e5d2f906066a8c90d5bc7735d3108ba8dc96778a2b2dae3a0f14 --file encrypted
