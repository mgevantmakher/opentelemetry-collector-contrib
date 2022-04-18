# Solace Receiver

Solace receiver receives traces from Solace PubSub+ broker.

Supported pipeline types: traces

## Getting Started

### Configuration

Settings:

- `brokers`: Non-empty list of Solace brokers is expected, if more than one entry is present.
  Receiver prepends 'amqps://' to each configured broker connection. Space can be used as a separator if more than one
  broker is configured. Only TLS secured AMQP connections are supported. Format for each entry is: \<host>:<secure AMQP
  port\>.
- `queue` (default = queue://#TRACE): The name of the Solace queue to get span trace messages from
- `max_unacknowledged` (default = 10): The maximum number of unacknowledged messages the Solace broker can transmit
- `transport`
    - `tls` TLS configuration, non tls connections are not permitted
        - `server_name_override`: ServerName indicates the name of the server requested by the client in order to
          support virtual hosting.
        - `insecure_skip_verify` (default = false): Whether to skip verifying the certificate from a broker or not.
        - `ca_file`: path to the CA cert trust store. For a client this verifies the server certificate. Should only be
          used if `insecure_skip_verify` is set to true.
        - `cert_file`: path to the TLS cert for client cert authentication
        - `key_file`: path to the TLS key for client cert authentication
- `auth`
    - `sasl_plain_text` SASL Plain authentication.
        - `username`: The username to use.
        - `password`: The password to use
    - `sasl_xauth2` SASL XOauth2 authentication.
        - `username`: The username to use.
        - `bearer`: The bearer token in plain text
    - `sasl_external`: SASL External required to be used for TLS client cert authentication.

Examples:

Example 1: connecting to a single locally running Solace broker (amqps://localhost:5671) using 'default' vpn with
configured client user with username 'otel' and a password 'otel01$'

```yaml
receivers:
  solace:
    auth:
      plain_text:
        username: otel
        password: otel01$
```

Example 2: connecting to a single Solace broker (amqps://myHost:5671) with configured client user with username 'otel'
and a password 'otel01$'

```yaml
receivers:
  solace:
    brokers: myHost:5671
    auth:
      sasl_plain_text:
        username: otel
        password: otel01$
```

