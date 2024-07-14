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

