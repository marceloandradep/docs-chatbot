Ingest
```
curl -X POST localhost:8080/ingest \
  -H 'content-type: application/json' \
  -d '{"doc_id":"readme","paths":["./docs/README.md"]}'
```

Ask
```
curl -X POST localhost:8080/ask \
  -H 'content-type: application/json' \
  -d '{"query":"What is this about?","doc_id":"readme","top_k":6}'
```