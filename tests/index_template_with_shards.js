PUT
_index_template/foo_template
{
  "index_patterns": [
    "foo*"
  ],
  "template": {
    "settings": {
      "index.number_of_shards": {{shards}},
      "index.number_of_replicas": 1
    }
  }
}