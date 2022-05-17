# dispatcher

Create a config file called dispatcher.json:

```json
{
    "overall-interval": "500ms",
    "scripts": [
        {
            "name": "a.sh",
            "interval": "10s"
        },
        {
            "name": "b.sh",
            "interval": "500ms"
        }
    ]
}
```

run the dispatcher in the same folder as the config and presto.