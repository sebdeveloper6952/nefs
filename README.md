### nefs

Let's try building an encrypted file sharing CLI tool over nostr.

- Use local relay for now.
- Let's build it with NIP-44 primitives for now.

### flow
1. choose recipient pubkey
2. generate conversation key
3. read file
4. encode into base64
5. split into parts of MAX_SIZE (65536 bytes)
6. encrypt
7. generate one event per part
8. event n references event n+1 in tag
9. send events

### another flow
1. read file, split into base64 chunks according to MAX_SIZE (nip44)
2. generate conversation key with recipient pubkey
3. encrypt each baes64 chunk
4. upload each chunk to a blossom server
5. build the event kind XXX that has one tag per file chunk like this
```
id: "<eid>",
kind: 79999,
pubkey: "<pk>",
tags: [
    ["chunk", "<chunk_hash>", "0"],
    ["chunk", "<chunk_hash>", "1"],
    ["chunk", "<chunk_hash>", "2"],
    ["chunk", "<chunk_hash>", "3"],
    ["chunk", "<chunk_hash>", "4"],
]
```
6. to receive file, pull event with ID `<eid>`
7. get `<pk>` CDN list
8. pull encrypted file chunks from CDN(s)
9. when all chunks are ready, decrypt and assemble
