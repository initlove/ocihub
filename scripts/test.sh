curl localhost:8080/oci/v1/first/tags/list

echo "\n--------"
curl localhost:8080/oci/v1/second/second/tags/list

echo "\n--------"
curl -X PUT -F "filename=@testmanifest" localhost:8080/oci/v1/second/second/manifest/v0.1

echo "\n--------"
curl localhost:8080/oci/v1/second/second/manifest/v0.1

echo "\n--------"
curl -X DELETE localhost:8080/oci/v1/second/second/manifest/v0.1

echo "\n--------"
curl localhost:8080/oci/v1/second/second/manifest/v0.1
