# Ming (Minio Gateway)

Minio gateway provides an aggregated view of Minio instances.

## Running ming

```
ming --config-dir ~/.ming/config.json --address localhost:8000
```

### Configuration file

```
{
  "hosts": [
    {
      "url": "http://localhost:9001",
      "accessKey": "SXO8VW2OFKKP2OG7AC85",
      "secretKey": "CKWSSgrUgvfUMTaNBkB63exet4WW+uNhQvi91Bc3"
    },
    {
      "url": "http://localhost:9002",
      "accessKey": "SXO8VW2OFKKP2OG7AC85",
      "secretKey": "CKWSSgrUgvfUMTaNBkB63exet4WW+uNhQvi91Bc3"
    }
  ],
}
```

### Upload and Download

Clients can upload to a special bucket "ming" using PUT request. i.e

```
curl -X PUT http://mingserver.com/ming/photo.jpg --data @/home/coder/photo.jpg
```

The response of PUT will be the URL that should be used for any subsequent download to this object. i.e
```
http://mingserver.com/ming/e1930b4927e6b6d92d120c7c1bba3421/photo.jpg
```

The hash `e1930b4927e6b6d92d120c7c1bba3421` will be internally used to decide the minio server from which the object needs to be fetched.

The client would need to store this URL (in databse) and use it to fetch the object.

