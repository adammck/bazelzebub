MANIFEST = {
    "prefix": "hello-world",
    "environments": [
        {
            "name": "stg",
            "datacenters": [
                {"name_short": "us1", "name_long": "us1.staging.company"}
            ],
        },
        {
            "name": "prd",
            "datacenters": [
                {
                    "name_short": "us1",
                    "name_long": "us1.prod.company",
                    "clusters": [
                        {
                            "name": "org9876",
                            "kafka_cluster": "km-1234",
                            "kafka_topic": "pts-org9876",
                            "shards": [
                                {"name": "S1", "partitions": [1, 2, 3, 4]},
                                {"name": "S2", "partitions": [5, 6]},
                                {"name": "S3", "partitions": [7, 8]},
                            ],
                        }
                    ],
                },
                {
                    "name_short": "eu1",
                    "name_long": "eu1.prod.company",
                    "clusters": [],
                },
            ],
        },
    ],
}
