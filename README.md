# dispatcher

Create a config file called dispatcher.json:

```json
{
    "overall-interval": "500ms",
    "scripts": [
        {
            "name": "10s-script.sh",
            "interval": "10s"
        },
        {
            "name": "500ms-script.sh",
            "interval": "500ms"
        }
    ]
}
```

run the dispatcher in the same folder as the config and presto.
