# nefs - nostr encrypted file sharing
Just an experiment of sharing encrypted files over Nostr.

### how?
Again, it's just an experiment so you may find it dumb.
1. To send a file: split a file in chunks, encrypt each chunk using [nip-44](https://github.com/nostr-protocol/nips/blob/master/44.md), and upload each chunk to a blossom server.
2. To receive the file: step 1 will produce an event ID, use this to know what chunks you need to download. Download each chunk, decrypt each chunk, and put the file back together.

### limitations
1. only one private key is able to decrypt the chunks

### usage
- coming soon™
