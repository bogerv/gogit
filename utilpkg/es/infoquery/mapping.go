package infoquery

const InfoMapping = `
{
    "mappings": {
        "info": {
            "properties": {
                "id": {
                    "type": "text",
                    "fields": {
                        "keyword": {
                            "type": "keyword",
                            "ignore_above": 256
                        }
                    }
                },
                "title": {
                    "type": "text",
                    "analyzer": "ik_max_word",
                    "search_analyzer": "ik_smart",
                    "fields": {
                        "keyword": {
                            "type": "keyword",
                            "ignore_above": 256
                        }
                    }
                },
                "type": {
                    "type": "long"
                },
                "language": {
                    "type": "text",
                    "fields": {
                        "keyword": {
                            "type": "keyword",
                            "ignore_above": 20
                        }
                    }
                },
                "content": {
                    "type": "text",
                    "fields": {
                        "keyword": {
                            "type": "keyword",
                            "ignore_above": 256
                        }
                    }
                },
                "isOnline": {
                    "type": "long"
                },
                "uuid": {
                    "type": "text"
                },
                "contentDetail": {
                    "type": "text",
                    "analyzer": "ik_max_word",
                    "search_analyzer": "ik_smart"
                },
                "checkTime": {
                    "type": "long"
                },
                "isDel": {
                    "type": "long"
                },
                "uptime": {
                    "type": "long"
                },
                "createTime": {
                    "type": "long"
                },
                "updateTime": {
                    "type": "long"
                }
            }
        }
    }
}
`
