# Mnm (Minio Minimal)

> This project is a work in progress.

Minio Minimal provides a simple GET/PUT API to provide aggregated view of Minio instances.

## Running mnm

```
mnm --config-dir ~/.mnm/config.json --address localhost:8000
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

Clients can upload to a special bucket "mnm" using PUT request. i.e

```
curl -X PUT http://mnmserver.com/mnm/photo.jpg --data @/home/coder/photo.jpg
```

The response of PUT will be the URL that should be used for any subsequent download to this object. i.e
```
http://mnmserver.com/mnm/e1930b4927e6b6d92d120c7c1bba3421/photo.jpg
```

The hash `e1930b4927e6b6d92d120c7c1bba3421` will be internally used to decide the minio server from which the object needs to be fetched.

The client would need to store this URL (in databse) and use it to fetch the object.

