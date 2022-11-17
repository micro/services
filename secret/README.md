Encrypted secret storage

# Secret Service

Store secrets, tokens, passwords and key-value config in encrypted data storage. 
Data is encrypted using AES-256.

Keys are individually allocated spaces for values. Values are encrypted at rest. 
Values can be strings or JSON with path based lookup if the latter.
