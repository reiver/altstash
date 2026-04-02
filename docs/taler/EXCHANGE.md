# Taler Exchange REST API

This document describes the GNU Taler exchange REST API as it relates to altstash's integration needs: fetching denomination keys, managing reserves, and withdrawing coins.

## Sources

- https://docs.taler.net/core/api-exchange.html — Exchange RESTful API
- https://docs.taler.net/core/api-common.html — Common API types and conventions
- https://docs.taler.net/design-documents/004-wallet-withdrawal-flow.html — Wallet Withdrawal Flow design document

## Development Exchange

altstash develops against the demo exchange at `demo.taler.net`.
The exchange base URL is `https://exchange.demo.taler.net/`.

## Common Types and Encoding Conventions

### Crockford Base32

All cryptographic keys, hashes, and signatures in the Taler API are encoded as Crockford Base32 strings.

### Amounts

In JSON wire format, amounts are strings in the format `"CURRENCY:DECIMAL"`:

```
"KUDOS:10.50"
"EUR:1.23"
"KUDOS:0.00000001"
```

In binary format (used inside signed structures), amounts are:

```c
struct TALER_AmountNBO {
    uint64_t value;              // network byte order, max 2^52
    uint32_t fraction;           // in units of 10^-8
    uint8_t  currency_code[12];  // zero-padded
};
```

The fractional part has 8 decimal digits of precision.
The smallest representable amount is 10^-8 of the base currency.

Note: altstash's internal coin file format uses the structured `{currency, value, fraction}` representation (from wallet-core's `AmountJson`), not the `"CURRENCY:DECIMAL"` wire string format.
Conversion between the two is needed when talking to the exchange API.

### Timestamps

```json
{"t_s": 1700000000}
{"t_s": "never"}
```

`t_s` is seconds since the Unix epoch.
The string `"never"` represents an unbounded timestamp.

### Key and Value Types

| Type             | Size     | Description                                         |
|------------------|----------|-----------------------------------------------------|
| `EddsaPublicKey` | 32 bytes | Ed25519 public key                                  |
| `EddsaSignature` | 64 bytes | Ed25519 signature (R + S)                           |
| `RsaPublicKey`   | variable | RSA public key                                      |
| `Cs25519Point`   | 32 bytes | Curve25519 point (CS denomination key)              |
| `HashCode`       | 64 bytes | SHA-512 hash                                        |
| `Timestamp`      | object   | `{"t_s": <seconds>}` or `{"t_s": "never"}`          |
| `RelativeTime`   | object   | `{"d_us": <microseconds>}` or `{"d_us": "forever"}` |

All cryptographic types are encoded as Crockford Base32 strings in JSON.

Taler uses **SHA-512** (64 bytes) as its primary hash function throughout, including for HKDFs.
SHA-256 (32 bytes, `ShortHashCode`) exists but is used only in specific places.

---

## 1. Fetching Denomination Keys

### Endpoint

```
GET /keys
```

### Query Parameters

| Parameter         | Type   | Required | Description |
|-------------------|--------|----------|-------------|
| `last_issue_date` | uint64 | No       | Unix timestamp (seconds). Returns only keys changed since this date. If the value does not exactly match a `stamp_start`, all keys are returned. |

### Response: `ExchangeKeysResponse`

Top-level fields:

| Field                    | Type                  | Description                                                |
|--------------------------|-----------------------|------------------------------------------------------------|
| `version`                | string                | Protocol version in libtool format (current:revision:age)  |
| `base_url`               | string                | Canonical exchange base URL                                |
| `currency`               | string                | Currency served by this exchange (e.g. `"EUR"`, `"KUDOS"`) |
| `currency_specification` | CurrencySpecification | Display details for the currency                           |
| `master_public_key`      | EddsaPublicKey        | Exchange's offline master signing key                      |
| `asset_type`             | string                | `"fiat"`, `"crypto"`, `"regional"`, or `"stock"`           |
| `denominations`          | DenomGroup[]          | **The denomination keys**                                  |
| `signkeys`               | SignKey[]             | Online signing keys (see below)                            |
| `auditors`               | AuditorKeys[]         | Auditor information (see below)                            |
| `recoup`                 | RecoupDenoms[]        | Revoked denominations eligible for recoup (see below)      |
| `accounts`               | ExchangeWireAccount[] | Bank accounts for funding reserves (see below)             |
| `wire_fees`              | object                | Wire fees by method (see below)                            |
| `global_fees`            | GlobalFees[]          | Global fee schedule (see below)                            |
| `reserve_closing_delay`  | RelativeTime          | How long before unused reserves are closed                 |
| `list_issue_date`        | Timestamp             | When this key list was issued                              |
| `exchange_sig`           | EddsaSignature        | Signature over this response                               |
| `exchange_pub`           | EddsaPublicKey        | Online signing key used for `exchange_sig`                 |

### Denomination Groups

The `denominations` array contains `DenomGroup` objects.
Each group represents denominations that share the same cipher type, face value, and fee schedule.
The group contains a `denoms` array with the individual denomination keys.

`DenomGroup` is a discriminated union on the `cipher` field:

| Variant                      | `cipher` value         | Extra group field   | Per-denom key field         |
|------------------------------|------------------------|---------------------|-----------------------------|
| `DenomGroupRsa`              | `"RSA"`                | —                   | `rsa_pub: RsaPublicKey`     |
| `DenomGroupCs`               | `"CS"`                 | —                   | `cs_pub: Cs25519Point`      |
| `DenomGroupRsaAgeRestricted` | `"RSA+age_restricted"` | `age_mask: AgeMask` | `rsa_pub: RsaPublicKey`     |
| `DenomGroupCsAgeRestricted`  | `"CS+age_restricted"`  | `age_mask: AgeMask` | `cs_pub: Cs25519Point`      |

`AgeMask` is an integer bitmask indicating which age groups the exchange supports for age-restricted denominations.
Each set bit represents an age threshold boundary (e.g., a mask with bits set at positions 8, 14, and 18 means the age groups are: under 8, 8-13, 14-17, and 18+).

All variants share `DenomGroupCommon` fields (at the group level):

| Field          | Type   | Description                                         |
|----------------|--------|-----------------------------------------------------|
| `value`        | Amount | Coin face value (e.g. `"KUDOS:5.00"`)               |
| `fee_withdraw` | Amount | Fee charged to withdraw a coin of this denomination |
| `fee_deposit`  | Amount | Fee charged to deposit (spend)                      |
| `fee_refresh`  | Amount | Fee charged to refresh (make change)                |
| `fee_refund`   | Amount | Fee charged for refund                              |

Each element in the `denoms` array has `DenomCommon` fields plus the cipher-specific public key:

| Field                   | Type               | Description                                                 |
|-------------------------|--------------------|-------------------------------------------------------------|
| `master_sig`            | EddsaSignature     | Master key signature over `TALER_DenominationKeyValidityPS` |
| `stamp_start`           | Timestamp          | When this denomination becomes valid                        |
| `stamp_expire_withdraw` | Timestamp          | Last date to withdraw coins of this denomination            |
| `stamp_expire_deposit`  | Timestamp          | Last date to deposit coins of this denomination             |
| `stamp_expire_legal`    | Timestamp          | Legal dispute settlement deadline                           |
| `lost`                  | boolean (optional) | `true` if the denomination's private key was compromised    |

### Example: RSA Denomination Group

```json
{
    "cipher": "RSA",
    "value": "KUDOS:5.00",
    "fee_withdraw": "KUDOS:0.01",
    "fee_deposit": "KUDOS:0.01",
    "fee_refresh": "KUDOS:0.01",
    "fee_refund": "KUDOS:0.01",
    "denoms": [
        {
            "rsa_pub": "<Crockford-Base32-encoded RSA public key>",
            "master_sig": "<Crockford-Base32-encoded EdDSA signature>",
            "stamp_start": {"t_s": 1700000000},
            "stamp_expire_withdraw": {"t_s": 1700100000},
            "stamp_expire_deposit": {"t_s": 1700200000},
            "stamp_expire_legal": {"t_s": 1700300000}
        }
    ]
}
```

### Currency Specification

The `currency_specification` field describes how to display the currency in a UI:

| Field                                  | Type              | Description                                                      |
|----------------------------------------|-------------------|------------------------------------------------------------------|
| `name`                                 | string            | Human-readable name (e.g. `"US Dollar"`)                         |
| `num_fractional_input_digits`          | Integer           | Digits a user may enter after the decimal point                  |
| `num_fractional_normal_digits`         | Integer           | Digits to render in normal font                                  |
| `num_fractional_trailing_zero_digits`  | Integer           | Digits to always render (padded with zeros if needed)            |
| `alt_unit_names`                       | {log10: string}   | Map of powers of 10 to symbols (e.g. `{"0": "EUR", "-2": "ct"}`) |
| `common_amounts`                       | Amount[]          | Shortcut amounts for UI buttons                                  |

### Wire Accounts

The `accounts` array contains the exchange's bank accounts where wire transfers should be sent to fund reserves.

Each `ExchangeWireAccount`:

| Field                    | Type                 | Description                                                       |
|--------------------------|----------------------|-------------------------------------------------------------------|
| `payto_uri`              | string               | Full `payto://` URI identifying the account and wire method       |
| `credit_restrictions`    | AccountRestriction[] | Restrictions on accounts sending funds to the exchange            |
| `debit_restrictions`     | AccountRestriction[] | Restrictions on accounts receiving funds from the exchange        |
| `master_sig`             | EddsaSignature       | Signature over `TALER_MasterWireDetailsPS`                        |
| `bank_label`             | string (optional)    | Display label for wallets                                         |
| `priority`               | Integer (optional)   | Display priority (0 if missing)                                   |
| `conversion_url`         | string (optional)    | URI to convert amounts from/to the currency                       |
| `open_banking_gateway`   | string (optional)    | Open banking gateway URL                                          |
| `wire_transfer_gateway`  | string (optional)    | Wire transfer gateway URL                                         |

The `payto_uri` is the key field — it tells the wallet where to send a wire transfer to fund a reserve.

`AccountRestriction` is a discriminated union on the `type` field:

| Variant                     | `type` value | Key field                                              |
|-----------------------------|--------------|--------------------------------------------------------|
| `RegexAccountRestriction`   | `"regex"`    | `payto_regex`: regex pattern that accounts must match  |
| `DenyAllAccountRestriction` | `"deny"`     | — (denies all accounts)                                |

### Online Signing Keys

The `signkeys` array contains the exchange's online signing keys, which are rotated periodically.
The exchange uses these keys to sign responses (like the `/keys` response itself via `exchange_sig`/`exchange_pub`).

Each `SignKey`:

| Field          | Type           | Description                                                              |
|----------------|----------------|--------------------------------------------------------------------------|
| `key`          | EddsaPublicKey | The exchange's online EdDSA signing public key                           |
| `stamp_start`  | Timestamp      | When this signing key becomes valid                                      |
| `stamp_expire` | Timestamp      | When the exchange stops using this key (may overlap with next key)       |
| `stamp_end`    | Timestamp      | When all signatures made by this key expire (for legal dispute purposes) |
| `master_sig`   | EddsaSignature | Master key signature over `TALER_ExchangeSigningKeyValidityPS`           |

### Auditors

The `auditors` array lists auditors that have certified denomination keys on this exchange.
Auditors independently verify the exchange's financial integrity.

Each `AuditorKeys`:

| Field               | Type                      | Description                                              |
|---------------------|---------------------------|----------------------------------------------------------|
| `auditor_pub`       | EddsaPublicKey            | The auditor's EdDSA signing public key                   |
| `auditor_url`       | string                    | The auditor's URL                                        |
| `auditor_name`      | string                    | The auditor's human-readable name                        |
| `denomination_keys` | AuditorDenominationKey[]  | Denomination keys the auditor affirms with its signature |

Each `AuditorDenominationKey`:

| Field         | Type           | Description                                            |
|---------------|----------------|--------------------------------------------------------|
| `denom_pub_h` | HashCode       | Hash of the denomination public key                    |
| `auditor_sig` | EddsaSignature | Auditor's signature over `TALER_ExchangeKeyValidityPS` |

### Global Fees

The `global_fees` array defines account-level and purse-related fees, each valid for a date range.

Each `GlobalFees`:

| Field                 | Type           | Description                                                                 |
|-----------------------|----------------|-----------------------------------------------------------------------------|
| `start_date`          | Timestamp      | When this fee schedule takes effect (inclusive)                             |
| `end_date`            | Timestamp      | When this fee schedule ends (exclusive)                                     |
| `history_fee`         | Amount         | Fee charged to retrieve account/reserve history                             |
| `account_fee`         | Amount         | Annual fee for maintaining an open account                                  |
| `purse_fee`           | Amount         | Fee charged for abandoned purses beyond the free limit                      |
| `history_expiration`  | RelativeTime   | How long the exchange preserves account history                             |
| `purse_account_limit` | Integer        | Number of concurrent purses allowed without paying `purse_fee`              |
| `purse_timeout`       | RelativeTime   | How long the exchange keeps a purse after expiry or merge                   |
| `master_sig`          | EddsaSignature | Master key signature over `TALER_GlobalFeesPS`                              |

### Wire Fees

The `wire_fees` field maps wire method names (e.g. `"iban"`) to arrays of fee schedules:

```
wire_fees: { [method: string]: AggregateTransferFee[] }
```

Each `AggregateTransferFee`:

| Field         | Type           | Description                                        |
|---------------|----------------|----------------------------------------------------|
| `wire_fee`    | Amount         | Per-transfer wire fee                              |
| `closing_fee` | Amount         | Per-transfer reserve closing fee                   |
| `start_date`  | Timestamp      | When this fee takes effect (inclusive)             |
| `end_date`    | Timestamp      | When this fee ends (exclusive)                     |
| `sig`         | EddsaSignature | Master key signature over `TALER_MasterWireFeePS`  |

### Recoup

The `recoup` array lists denominations whose private keys have been compromised and are being revoked under the emergency protocol.

Each `RecoupDenoms`:

| Field         | Type     | Description                                              |
|---------------|----------|----------------------------------------------------------|
| `h_denom_pub` | HashCode | Hash of the denomination public key being revoked        |

Coins of these denominations should not be accepted for payment.
Holders of such coins can recover their value through the `/recoup` endpoint by proving ownership (requires the coin's blinding key and reserve information).

The exchange does not sign this list — because the primary recoup scenario involves the exchange having lost its signing keys, such a signature would be meaningless.

---

## 2. Reserve Management

### How Reserves Work

A reserve is a balance held at the exchange, created by wire transfer.
The wallet generates an EdDSA key pair (reserve private key / reserve public key) before initiating the transfer.
The wire transfer subject line must contain the reserve public key encoded as Crockford Base32.
The exchange's Taler Wire Gateway monitors incoming transfers and credits the matching reserve automatically.

After funding, the wallet can withdraw coins from the reserve until the balance is exhausted or the reserve expires.

### Funding a Reserve on `demo.taler.net`

The demo environment at `demo.taler.net` provides a test bank at `bank.demo.taler.net` that can be used to fund reserves without real money.

To fund a reserve:

1. Generate an EdDSA key pair for the reserve
2. Look up the exchange's bank account from the `accounts` field in the `GET /keys` response (the `payto_uri`)
3. Initiate a wire transfer via the demo bank to the exchange's account, with the reserve public key (Crockford Base32) as the wire transfer subject
4. The exchange detects the incoming transfer and credits the reserve

The demo bank provides a web interface and API for creating test accounts and initiating transfers.

### Checking Reserve Status

```
GET /reserves/$RESERVE_PUB
```

`$RESERVE_PUB` is the reserve's EdDSA public key, Crockford Base32 encoded.

#### Query Parameters

| Parameter    | Type   | Required | Description                                                                                    |
|--------------|--------|----------|------------------------------------------------------------------------------------------------|
| `timeout_ms` | uint64 | No       | Long-poll: exchange waits up to this many milliseconds for incoming funds before returning 404 |

#### Response (200 OK): `ReserveSummary`

| Field                | Type               | Description                                           |
|----------------------|--------------------|-------------------------------------------------------|
| `balance`            | Amount             | Remaining funds in the reserve                        |
| `last_origin`        | string (optional)  | `payto://` URI of the most recent funding source      |
| `reserve_expiration` | Timestamp          | When residual value is returned to the origin account |
| `maximum_age_group`  | Integer (optional) | If set, age-restricted withdrawal is required         |

#### Response (404 Not Found)

Reserve public key is unknown to the exchange (not yet funded, or expired and purged).

---

## 3. The Withdrawal Protocol

### Overview

Withdrawal is the process of converting reserve balance into Taler coins held locally.
It uses Chaumian blind signatures: the exchange signs the coin without seeing its identity, so it cannot later link the coin back to the withdrawal (this is the privacy guarantee).

### Complete Flow (RSA Denominations)

#### Step 1: Select Denominations

Choose which denominations to withdraw based on:

* The reserve balance
* Denomination face values and withdrawal fees
* Denomination validity (must be between `stamp_start` and `stamp_expire_withdraw`)
* Denomination status (not `lost`, not in `recoup` list)

The goal is to maximize the value withdrawn while minimizing fees.

#### Step 2: Generate Planchets (Client-Side)

For each coin to withdraw:

1. Generate a fresh Ed25519 key pair: `coin_priv` / `coin_pub`
2. Generate a random RSA blinding factor `beta`
3. Compute the SHA-512 hash of `coin_pub` (this becomes the "planchet")
4. Blind the planchet: `coin_ev = blind(SHA512(coin_pub), beta, denom_rsa_pub)` — produces a `BlindedCoinEv`
5. Persist `coin_priv`, `beta`, and denomination info to disk (critical for recovery if the process is interrupted)

#### Step 3: Compute the Reserve Signature (Client-Side)

Sign a `TALER_WithdrawRequestPS` structure with the reserve's private key.
The signed structure contains:

```c
struct TALER_WithdrawRequestPS {
    // Purpose: TALER_SIGNATURE_WALLET_RESERVE_WITHDRAW
    struct GNUNET_CRYPTO_EccSignaturePurpose purpose;
    // Total coin value (excluding fees)
    struct TALER_Amount amount;
    // Total withdrawal fee
    struct TALER_Amount fee;
    // Running SHA-512 hash over all BlindedCoinHashP values
    // (each captures the blinded coin hash AND the denomination key hash)
    struct TALER_HashPlanchetsP h_planchets;
    // Master seed for CS blinding (all zeros if no CS denominations)
    struct TALER_BlindingMasterSecretP blinding_seed;
    // Maximum age group (0 if no age restriction)
    uint32_t max_age_group;
    // Age mask from exchange config (0 if no age restriction)
    struct TALER_AgeMask mask;
};
```

#### Step 4: Send the Withdrawal Request

```
POST /withdraw
```

##### Request Body: `WithdrawRequest`

| Field           | Type                          | Description                                            |
|-----------------|-------------------------------|--------------------------------------------------------|
| `cipher`        | string                        | `"ED25519"` — reserve signature cipher                 |
| `reserve_pub`   | EddsaPublicKey                | The reserve's public key                               |
| `denoms_h`      | HashCode[]                    | Array of denomination public key hashes (one per coin) |
| `coin_evs`      | CoinEnvelope[]                | Array of blinded coin envelopes (one per coin)         |
| `reserve_sig`   | EddsaSignature                | Signature of `TALER_WithdrawRequestPS`                 |
| `blinding_seed` | BlindingMasterSeed (optional) | Required if any CS denominations are included          |
| `max_age`       | Integer (optional)            | Required if age restriction applies                    |

`CoinEnvelope` is a discriminated union by cipher:

**RSA variant:**
```json
{
    "cipher": "RSA",
    "rsa_blinded_planchet": "<Crockford-Base32-encoded blinded planchet>"
}
```

**Clause-Schnorr variant:**
```json
{
    "cipher": "CS",
    "cs_nonce": "<Crockford-Base32>",
    "cs_blinded_c0": "<Crockford-Base32>",
    "cs_blinded_c1": "<Crockford-Base32>"
}
```

##### Response (200 OK): `WithdrawResponse`

| Field     | Type                           | Description                                             |
|-----------|--------------------------------|---------------------------------------------------------|
| `ev_sigs` | BlindedDenominationSignature[] | Blinded signatures, same order as `coin_evs` in request |

Each `BlindedDenominationSignature` is a discriminated union:

**RSA variant:**
```json
{
    "cipher": "RSA",
    "blinded_rsa_signature": "<Crockford-Base32>"
}
```

**Clause-Schnorr variant:**
```json
{
    "cipher": "CS",
    "b": 0,
    "s": "<Crockford-Base32 Cs25519Scalar>"
}
```

##### Error Responses

| Status | Error Code                                             | Meaning                                       |
|--------|--------------------------------------------------------|-----------------------------------------------|
| 400    | —                                                      | Malformed request                             |
| 403    | `TALER_EC_EXCHANGE_WITHDRAW_RESERVE_SIGNATURE_INVALID` | Invalid reserve signature                     |
| 404    | —                                                      | Reserve unknown or denomination unknown       |
| 409    | `TALER_EC_EXCHANGE_WITHDRAW_INSUFFICIENT_FUNDS`        | Insufficient funds in reserve                 |
| 410    | —                                                      | Denomination expired or revoked               |
| 412    | —                                                      | Denomination not yet valid                    |
| 451    | —                                                      | KYC (Know Your Customer) requirements not met |
| 501    | —                                                      | Unsupported cipher                            |

#### Step 5: Unblind the Signatures (Client-Side)

**For RSA:**

1. Unblind: `signature = unblind(blinded_rsa_signature, beta, denom_rsa_pub)`
2. Verify: `RSA_verify(SHA512(coin_pub), signature, denom_rsa_pub)`
3. If valid, the coin is ready to use

**For Clause-Schnorr:**

Use the `b` (bit) and `s` (scalar) values along with the wallet's blinding secrets to compute the final `(R, s)` Schnorr signature.

#### Step 6: Persist the Coin

Store the completed coin with all fields needed for spending and recovery.
In altstash, this is a `.talercoin` file in `$XDG_DATA_HOME/altstash/talercoins/`.

### RSA Blind Signature Cryptography

The RSA blind signature scheme (Chaum's scheme) works as follows:

1. **Wallet picks blinding factor:** `b` uniformly at random from Z_n (where `n` is the RSA modulus from the denomination key)
2. **Wallet blinds the message:** `m' = m * b^e mod n` (where `e` is the RSA public exponent, `m` is the SHA-512 hash of the coin public key)
3. **Exchange signs blindly:** `s' = (m')^d mod n` (using the denomination private key `d`)
4. **Wallet unblinds:** `s = s' * b^(-1) mod n`
5. **Verification:** `s^e mod n == m` — anyone can verify using the denomination's public key

The exchange never sees `m` (the actual coin identity), only `m'` (the blinded version), so it cannot link the signed coin back to the withdrawal.
This is the core of Taler's privacy guarantee.

### Idempotency

The withdrawal endpoint is idempotent: repeating exactly the same request yields the same response.
If the network fails during withdrawal, the wallet can safely retry with identical parameters without risking double-spending the reserve or losing coins.

---

## 4. Clause-Schnorr (CS) Withdrawal — Extra Step

CS denominations require an additional round-trip before the main withdrawal request.

### Blinding Preparation

```
POST /blinding-prepare
```

#### Request Body: `BlindingPrepareRequestCS`

| Field       | Type                     | Description                                 |
|-------------|--------------------------|---------------------------------------------|
| `cipher`    | string                   | `"CS"`                                      |
| `operation` | string                   | `"withdraw"` or `"melt"`                    |
| `seed`      | BlindingMasterSeed       | Unique per request — MUST NOT be reused     |
| `nks`       | BlindingInputParameter[] | One per CS-denomination coin                |

Each `BlindingInputParameter`:

| Field            | Type     | Description                            |
|------------------|----------|----------------------------------------|
| `coin_offset`    | Integer  | Position in the list of fresh coins    |
| `denom_pub_hash` | HashCode | Hash of the CS denomination public key |

#### Response (200 OK): `BlindingPrepareResponseCS`

| Field    | Type                     | Description                                             |
|----------|--------------------------|---------------------------------------------------------|
| `cipher` | string                   | `"CS"`                                                  |
| `r_pubs` | [CSRPublic, CSRPublic][] | Array of pairs of Curve25519 points, one pair per input |

The wallet uses these `r_pubs` pairs to compute the blinded challenges (`cs_blinded_c0`, `cs_blinded_c1`) that go into the `CSCoinEnvelope` in the subsequent `POST /withdraw` request.
The same `seed` must be provided in both the `/blinding-prepare` and `/withdraw` requests.

---

## 5. Relevance to altstash

### What altstash needs for the TODO items

**"Connect to a Taler exchange — fetch denomination keys":**

* `GET /keys` on `https://exchange.demo.taler.net/`
* Parse the `ExchangeKeysResponse`, extract denomination groups
* Store/cache the denomination keys for use during withdrawal
* Need to handle both RSA and CS cipher types

**"Withdraw coins from a funded reserve":**

* Full withdrawal flow: select denominations → generate planchets → blind → sign → `POST /withdraw` → unblind → persist
* Requires implementing RSA blind signatures (and optionally CS blind signatures)
* Requires EdDSA signing for the reserve signature
* Reserve must be pre-funded via wire transfer (for `demo.taler.net`, the demo bank at `bank.demo.taler.net` can be used)

**"Auto-create `talercoins/` directory if it doesn't exist":**

* No API interaction needed — purely local filesystem operation
* `os.MkdirAll` on the talercoins path before writing coin files

### Complexity Assessment

The denomination key fetch (`GET /keys`) is straightforward HTTP + JSON parsing.

The withdrawal protocol is significantly more complex due to:

1. Cryptographic operations (RSA blinding/unblinding, EdDSA signing)
2. Binary structure signing (`TALER_WithdrawRequestPS` with specific byte layout)
3. Crockford Base32 encoding/decoding
4. Denomination selection algorithm (optimizing value vs. fees)
5. Recovery/persistence of intermediate state (planchets before exchange responds)

For MVP, focusing on RSA denominations only is reasonable — CS support can be added later.
