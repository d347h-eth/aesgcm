A simple CLI tool to encrypt and decrypt data with password using AES-GCM cryptographic algorithm.
PBKDF2 is used for key derivation with SHA-512 as a hash function.
Key derivation parameters, nonce and ciphertext are packed into JSON and by default the final output is Base64 encoded. It's also possible to render a QR code of the final output.

## Usage
```bash
make build
./build/aesgcm
```

Use `--help` to see instructions. All default parameters and behaviors can be altered with the relevant flags.

Salt and nonce are randomly generated in the runtime (by default 128 and 12 bytes accordingly).

Default output paths are: "INPUT_FILENAME.aes" for encryption and "INPUT_FILENAME.txt" for decryption.

By default, the minimal password length is required to be at least 8 characters.
